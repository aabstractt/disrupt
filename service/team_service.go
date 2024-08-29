package service

import (
	"context"
	"errors"
	"github.com/bitrule/hcteams/startup"
	"github.com/bitrule/hcteams/startup/message"
	"github.com/bitrule/hcteams/team"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"sync"
)

var IDKey = "_id"

type TeamService struct {
	conf startup.TeamsConfig
	col  *mongo.Collection

	teamsMu sync.RWMutex
	teams   map[string]team.Team

	teamIdsMu sync.RWMutex
	teamIds   map[string]string

	membersMu sync.RWMutex
	members   map[string]string
}

// LookupByMember looks up a team by a member's XUID.
func (s *TeamService) LookupByMember(xuid string) *team.PlayerTeam {
	s.membersMu.RLock()
	defer s.membersMu.RUnlock()

	id, ok := s.members[xuid]
	if !ok {
		return nil
	}

	return s.LookupById(id).(*team.PlayerTeam)
}

// LookupByName looks up a team by its name.
func (s *TeamService) LookupByName(name string) team.Team {
	s.teamsMu.RLock()
	defer s.teamsMu.RUnlock()

	s.teamIdsMu.RLock()
	defer s.teamIdsMu.RUnlock()

	if id, ok := s.teamIds[name]; ok {
		return s.teams[id]
	}

	return nil
}

// LookupById looks up a team by its ID.
func (s *TeamService) LookupById(id string) team.Team {
	s.teamsMu.RLock()
	defer s.teamsMu.RUnlock()

	if t, ok := s.teams[id]; ok {
		return t
	}

	return nil
}

// Delete deletes a team by its ID.
func (s *TeamService) Delete(id string) {
	s.teamsMu.Lock()
	defer s.teamsMu.Unlock()

	s.teamIdsMu.Lock()
	defer s.teamIdsMu.Unlock()

	if t, ok := s.teams[id]; ok {
		delete(s.teams, id)
		delete(s.teamIds, strings.ToLower(t.Tracker().Name()))
	}
}

func (s *TeamService) DeleteMember(xuid string) {
	s.membersMu.Lock()
	defer s.membersMu.Unlock()

	if _, ok := s.members[xuid]; ok {
		delete(s.members, xuid)
	}
}

// CacheMember caches a member's team ID.
func (s *TeamService) CacheMember(xuid, teamId string) {
	s.membersMu.Lock()
	s.members[xuid] = teamId
	s.membersMu.Unlock()
}

// cache caches a team.
func (s *TeamService) cache(t team.Team) {
	s.teamsMu.Lock()
	s.teams[t.Tracker().Id()] = t
	s.teamsMu.Unlock()

	s.teamIdsMu.Lock()
	s.teamIds[strings.ToLower(t.Tracker().Name())] = t.Tracker().Id()
	s.teamIdsMu.Unlock()
}

// Create creates a team.
// Use this function into a goroutine to prevent blocking the main thread.
func (s *TeamService) Create(p *player.Player, name, teamType string) {
	if t := team.Empty(p.XUID(), name, teamType); t == nil {
		p.Message(text.Red + "Failed to create the team: Team is nil")
	} else if err := s.Save(t); err != nil {
		p.Message(text.Red + "Failed to create the team: " + err.Error())
	} else {
		// Store the team in the service.
		s.cache(t)

		p.Message(message.SuccessSelfTeamCreated.Build(name))

		_, ok := t.(*team.PlayerTeam)
		if !ok {
			return
		}

		s.CacheMember(p.XUID(), t.Tracker().Id())

		if _, err = chat.Global.WriteString(message.SuccessTeamCreated.Build(p.Name(), name)); err != nil {
			p.Message(team.Prefix + text.Red + "Failed to broadcast team creation: " + err.Error())
		}
	}
}

// Invite invites a player to a team.
// This function will add the player to the team's invites and broadcast a message to the player and the team.
// The target can accept the invite by using the '/team accept' command or decline it by using the '/team decline' command.
// I think this function should be removed because its only used one time in the codebase.
func (s *TeamService) Invite(t *team.PlayerTeam, p *player.Player) error {
	if t.Member(p.XUID()) == team.Undefined {
		return errors.New(message.ErrPlayerAlreadyMember.Build(p.Name()))
	} else if s.LookupByMember(p.XUID()) != nil {
		return errors.New(message.ErrPlayerAlreadyInTeam.Build(p.Name()))
	} else if t.HasInvite(p.XUID()) {
		return errors.New(message.ErrPlayerAlreadyInvited.Build(p.Name()))
	}

	t.AddInvite(p.XUID())

	p.Message(message.SuccessTeamInviteReceived.Build(p.Name(), t.Tracker().Name()))
	t.Broadcast(message.SuccessBroadcastTeamInviteSent.Build(p.Name(), t.Tracker().Name()))

	return nil
}

// Disband disbands a team
// This function will delete the team from the repository and broadcast a message to all the members.
// also, it will delete all the members from the team and the service.
// Use this function into a goroutine to prevent blocking the main thread.
func (s *TeamService) Disband(t *team.PlayerTeam) error {
	if s.col == nil {
		return errors.New("missing repository")
	} else if u := userService.LookupByXUID(t.Ownership()); u == nil {
		return errors.New("leader not found")
	} else if r, err := s.col.DeleteOne(context.TODO(), bson.M{IDKey: t.Tracker().Id()}); err != nil {
		return err
	} else if r.DeletedCount == 0 {
		return errors.New("team not found into our database")
	} else {
		t.Broadcast(message.SuccessTeamDisband.Build(u.Name(), t.Tracker().Name()))

		for xuid := range t.Members() {
			s.DeleteMember(xuid)
		}
	}

	return nil
}

func (s *TeamService) Save(t team.Team) error {
	if s.col == nil {
		return errors.New("missing repository")
	}

	r, err := s.col.UpdateOne(context.TODO(), bson.M{IDKey: t.Tracker().Id()}, bson.M{"$set": t.Marshal()})
	if err != nil {
		return errors.New("failed to save the team: " + err.Error())
	}

	if r.MatchedCount == 0 && r.UpsertedCount == 0 {
		return errors.New("failed to save the team: no documents matched the filter")
	}

	return nil
}

// DoTick ticks all the system teams.
// This function should be called every tick.
func (s *TeamService) DoTick() {
	s.teamsMu.RLock()
	defer s.teamsMu.RUnlock()

	for _, t := range s.teams {
		if st, ok := t.(*team.SystemTeam); ok {
			st.DoTick()
		}
	}
}

func (s *TeamService) Config() startup.TeamsConfig {
	return s.conf
}

// Team returns the team service.
func Team() *TeamService {
	return teamService
}

var teamService = &TeamService{
	teams:   make(map[string]team.Team),
	teamIds: make(map[string]string),
	members: make(map[string]string),
}

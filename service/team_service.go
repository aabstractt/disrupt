package service

import (
	"errors"
	"github.com/bitrule/hcteams/common/message"
	"github.com/bitrule/hcteams/repository"
	"github.com/bitrule/hcteams/team"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
	"sync"
)

type TeamService struct {
	repository repository.Repository[team.Team]

	valuesMu sync.RWMutex
	values   map[string]team.Team

	identifiersMu sync.RWMutex
	identifiers   map[string]string

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
	s.valuesMu.RLock()
	defer s.valuesMu.RUnlock()

	s.identifiersMu.RLock()
	defer s.identifiersMu.RUnlock()

	if id, ok := s.identifiers[name]; ok {
		return s.values[id]
	}

	return nil
}

// LookupById looks up a team by its ID.
func (s *TeamService) LookupById(id string) team.Team {
	s.valuesMu.RLock()
	defer s.valuesMu.RUnlock()

	if t, ok := s.values[id]; ok {
		return t
	}

	return nil
}

// Delete deletes a team by its ID.
func (s *TeamService) Delete(id string) {
	s.valuesMu.Lock()
	defer s.valuesMu.Unlock()

	s.identifiersMu.Lock()
	defer s.identifiersMu.Unlock()

	if t, ok := s.values[id]; ok {
		delete(s.values, id)
		delete(s.identifiers, strings.ToLower(t.Tracker().Name()))
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
	s.valuesMu.Lock()
	s.values[t.Tracker().Id()] = t
	s.valuesMu.Unlock()

	s.identifiersMu.Lock()
	s.identifiers[strings.ToLower(t.Tracker().Name())] = t.Tracker().Id()
	s.identifiersMu.Unlock()
}

// Create creates a team.
// Use this function into a goroutine to prevent blocking the main thread.
func (s *TeamService) Create(p *player.Player, name, teamType string) {
	t := team.Empty(p.XUID(), name, teamType)
	if t == nil {
		p.Message(text.Red + "Failed to create the team: Team is nil")

		return
	}

	r, err := s.repository.Insert(t)
	if err != nil {
		p.Message(text.Red + "Failed to create the team: " + err.Error())

		return
	}

	if r.ModifiedCount > 0 {
		p.Message(text.DarkRed + "An error occurred while creating the team.")
		p.Message(text.Red + "Error: " + text.DarkRed + "Team already exists.")

		return
	}

	if _, ok := t.(*team.PlayerTeam); ok {
		if _, err = chat.Global.WriteString(message.SuccessTeamCreated.Build(p.Name(), name)); err != nil {
			p.Message(team.Prefix + text.Red + "Failed to broadcast team creation: " + err.Error())

			return
		}

		s.CacheMember(p.XUID(), t.Tracker().Id())
	}

	p.Message(message.SuccessSelfTeamCreated.Build(name))

	// Store the team in the service.
	s.cache(t)
}

// Invite invites a player to a team.
func (s *TeamService) Invite(t *team.PlayerTeam, p *player.Player) error {
	if t.HasMember(p.XUID()) {
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

func (s *TeamService) Disband(t *team.PlayerTeam) error {
	if s.repository == nil {
		return errors.New("missing repository")
	}

	u := userService.LookupByXUID(t.Ownership())
	if u == nil {
		return errors.New("leader not found")
	}

	if r, err := s.repository.Delete(t.Tracker().Id()); err != nil {
		return err
	} else if r.DeletedCount == 0 {
		return errors.New("team not found into our database")
	}

	t.Broadcast(message.SuccessBroadcastTeamDisbanded.Build(u.Name(), t.Tracker().Name()))

	for xuid := range t.Members() {
		s.DeleteMember(xuid)
	}

	return nil
}

// Team returns the team service.
func Team() *TeamService {
	return teamService
}

var teamService = &TeamService{
	values:      make(map[string]team.Team),
	identifiers: make(map[string]string),
	members:     make(map[string]string),
}

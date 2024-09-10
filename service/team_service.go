package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitrule/disrupt"
	"github.com/bitrule/disrupt/config"
	"github.com/bitrule/disrupt/message"
	"github.com/bitrule/disrupt/team"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"strings"
	"sync"
	"time"
)

var IDKey = "_id"

type TeamService struct {
	col *mongo.Collection // Repository

	teamsMu sync.RWMutex         // Protects teams
	teams   map[string]team.Team // Team ID -> Team

	teamIdsMu sync.RWMutex      // Protects teamIds
	teamIds   map[string]string // Team name as lower case -> Team ID

	teamsPerChunkMu sync.RWMutex                           // Protects teamsPerChunk
	teamsPerChunk   map[string]map[world.ChunkPos][]string // World name -> Chunk position -> Team ID

	membersMu sync.RWMutex      // Protects members
	members   map[string]string // XUID -> Team ID
}

// LookupByMember looks up a team by a member's XUID.
func (s *TeamService) LookupByMember(xuid string) *team.PlayerTeam {
	s.membersMu.RLock()
	defer s.membersMu.RUnlock()

	id, ok := s.members[xuid]
	if !ok {
		return nil
	}

	if t, ok := s.LookupById(id).(*team.PlayerTeam); ok {
		return t
	}

	panic(fmt.Sprintf("team '%s' is not an instance of '*team.PlayerTeam'", id))
}

// LookupByName looks up a team by its name.
func (s *TeamService) LookupByName(name string) team.Team {
	s.teamIdsMu.RLock()
	defer s.teamIdsMu.RUnlock()

	if id, ok := s.teamIds[strings.ToLower(name)]; ok {
		return s.LookupById(id)
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

// LookupByChunk looks up teams by a world and a Vec3.
func (s *TeamService) LookupByChunk(w *world.World, vec3 mgl64.Vec3) []team.Team {
	s.teamsPerChunkMu.RLock()
	defer s.teamsPerChunkMu.RUnlock()

	chunks, ok := s.teamsPerChunk[w.Name()]
	if !ok || chunks == nil {
		return nil
	}

	teamIds, ok := chunks[world.ChunkPos{int32(math.Floor(vec3[0])) >> 4, int32(math.Floor(vec3[2])) >> 4}]
	if !ok || teamIds == nil {
		return nil
	}

	teams := make([]team.Team, 0, len(teamIds))
	for _, id := range teamIds {
		if t := s.LookupById(id); t != nil {
			teams = append(teams, t)
		}
	}

	return teams
}

// LookupAt looks up a team by a Vec3.
// First it looks up the teams in the chunk the Vec3 is in, then it checks if the Vec3 is within any of the
// teams' bounding boxes. Also, see LookupByChunk.
func (s *TeamService) LookupAt(w *world.World, vec3 mgl64.Vec3) team.Team {
	teamsPerChunk := s.LookupByChunk(w, vec3)
	if teamsPerChunk == nil || len(teamsPerChunk) == 0 {
		return nil
	}

	for _, t := range teamsPerChunk {
		bBoxes, ok := t.Tracker().Cuboids()[w.Name()]
		if !ok || bBoxes == nil {
			continue
		}

		for _, bbox := range bBoxes {
			if bbox.Vec3Within(vec3) {
				return t
			}
		}
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

// DisplayName returns the display name of a team.
// This function will return the display name of a team based on the player's role in the team.
func (s *TeamService) DisplayName(p *player.Player, t team.Team) string {
	if v, ok := t.Tracker().Option(team.DisplayNameKeyOption).(string); ok {
		return v
	}

	var colour string
	if pt, ok := t.(*team.PlayerTeam); ok {
		colour = s.DisplayColour(p, pt)
	} else {
		colour = text.Red
	}

	return colour + t.Tracker().Name()
}

// DisplayColour returns the display colour of a team.
// This function will return the display colour of a team based on the player's role in the team.
func (s *TeamService) DisplayColour(p *player.Player, t *team.PlayerTeam) string {
	// If the player is member, his role never will be undefined.
	if t.Member(p.XUID()) != team.Undefined {
		return config.TeamConfig().Display.FriendlyColour
	} else if t.HasInvite(p.XUID()) {
		return config.TeamConfig().Display.InvitedColour
	}

	return config.TeamConfig().Display.EnemyColour
}

// Create creates a team.
// Use this function into a goroutine to prevent blocking the main thread.
func (s *TeamService) Create(p *player.Player, t team.Team) {
	if err := s.Save(t); err != nil {
		p.Message(text.Red + "Failed to create the team: " + err.Error())
	} else {
		// Store the team in the service.
		s.cache(t)

		p.Message(message.SuccessSelfTeamCreated.Build(t.Tracker().Name()))

		if _, ok := t.(*team.PlayerTeam); !ok {
			return
		}

		s.CacheMember(p.XUID(), t.Tracker().Id())

		if _, err = chat.Global.WriteString(message.SuccessTeamCreated.Build(p.Name(), t.Tracker().Name())); err != nil {
			p.Message(team.Prefix + text.Red + "Failed to broadcast team creation: " + err.Error())
		}
	}
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
		for xuid := range t.Members() {
			s.DeleteMember(xuid)
		}

		t.Broadcast(message.SuccessTeamDisband.Build(u.Name(), t.Tracker().Name()))

		s.Delete(t.Tracker().Id())
	}

	return nil
}

func (s *TeamService) Save(t team.Team) error {
	if s.col == nil {
		return errors.New("missing repository")
	}

	r, err := s.col.UpdateOne(context.TODO(), bson.M{IDKey: t.Tracker().Id()}, bson.M{"$set": t.Marshal()})
	if err != nil {
		return err
	}

	if r.MatchedCount == 0 && r.UpsertedCount == 0 {
		return errors.New("no documents matched the filter")
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

func (s *TeamService) Hook() error {
	if s.col != nil {
		return errors.New("repository already set")
	}

	s.col = disrupt.Mongo.Database(config.DBConfig().DBName).Collection("teams")

	cur, err := s.col.Find(context.TODO(), bson.M{})
	if err != nil {
		return errors.Join(errors.New("failed to hook the repository: "), err)
	}

	for cur.Next(context.Background()) {
		var body map[string]interface{}
		if err = cur.Decode(&body); err != nil {
			return errors.Join(errors.New("failed to decode the body of team: "), err)
		}

		t, err := team.Unmarshal(body)
		if err != nil {
			return errors.Join(errors.New("failed to unmarshal the team: "), err)
		}

		s.cache(t)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err = cur.Close(ctx); err != nil {
		return errors.Join(errors.New("failed to close the cursor: "), err)
	}

	if err = ctx.Err(); err != nil {
		return errors.Join(errors.New("failed to check the context: "), err)
	}

	return nil
}

// Shutdown saves all the teams in the service.
func (s *TeamService) Shutdown() error {
	if s.col == nil {
		return errors.New("missing repository")
	}

	s.teamsMu.RLock()
	defer s.teamsMu.RUnlock()

	systemCount := 0
	playerCount := 0

	for _, t := range s.teams {
		if err := s.Save(t); err != nil {
			return errors.New("failed to save the team '" + t.Tracker().Name() + "': " + err.Error())
		}

		if _, ok := t.(*team.SystemTeam); ok {
			systemCount++
		} else {
			playerCount++
		}
	}

	disrupt.Log.Infof("Successfully saved %d 'player' team(s) in our database", playerCount)
	disrupt.Log.Infof("Successfully saved %d 'admin' team(s) in our database", systemCount)

	return nil
}

// Team returns the team service.
func Team() *TeamService {
	return teamService
}

var teamService = &TeamService{
	teams:         make(map[string]team.Team),
	teamIds:       make(map[string]string),
	members:       make(map[string]string),
	teamsPerChunk: make(map[string]map[world.ChunkPos][]string),
}

package service

import (
	"github.com/bitrule/hcteams/team"
	"strings"
	"sync"
)

type TeamService struct {
	valuesMu sync.RWMutex
	values   map[string]team.Team

	identifiersMu sync.RWMutex
	identifiers   map[string]string

	membersMu sync.RWMutex
	members   map[string]string
}

// LookupByMember looks up a team by a member's XUID.
func (s *TeamService) LookupByMember(xuid string) team.Team {
	s.membersMu.RLock()
	defer s.membersMu.RUnlock()

	if id, ok := s.members[xuid]; ok {
		return s.LookupById(id)
	}

	return nil
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

func Team() *TeamService {
	return nil
}

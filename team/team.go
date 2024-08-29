package team

import (
	"github.com/bitrule/hcteams/repository"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
	"sync"
)

var (
	repo repository.Repository[Team]

	teamsIdMu sync.RWMutex
	teamsId   = make(map[string]string)

	teamsMu sync.RWMutex
	teams   = make(map[string]Team)

	Prefix = text.Blue + "[" + text.Yellow + "HCTeams" + text.Blue + "] " + text.Reset

	SystemTeamType = "System"
	PlayerTeamType = "Player"
)

type Team interface {

	// Tracker returns the team's tracker
	Tracker() *Tracker
	// Unmarshal unmarshals the team's tracker from a map
	Unmarshal(prop map[string]interface{}) error
	// Marshal returns the team's tracker as a map
	Marshal() (map[string]interface{}, error)
}

// Lookup returns the team with the given ID
func Lookup(id string) Team {
	teamsMu.RLock()
	defer teamsMu.RUnlock()

	return teams[id]
}

// Store registers the team into the cache
func Store(t Team) {
	// Lock the teams map to prevent deadlocks
	teamsMu.Lock()
	teams[t.Tracker().Id()] = t
	teamsMu.Unlock()

	// Lock the teamsId map to prevent deadlocks
	teamsIdMu.Lock()
	teamsId[strings.ToLower(t.Tracker().Name())] = t.Tracker().Id()
	teamsIdMu.Unlock()

	// Lock the membersId map to prevent deadlocks
	membersMu.Lock()
	for xuid := range t.(*PlayerTeam).Members() {
		membersId[xuid] = t.Tracker().Id()
	}
	membersMu.Unlock()
}

// Delete removes the team from the cache
func Delete(id string) {
	teamsMu.Lock()
	defer teamsMu.Unlock()

	if t, ok := teams[id]; ok {
		teamsIdMu.Lock()
		delete(teamsId, strings.ToLower(t.Tracker().Name()))
		teamsIdMu.Unlock()

		delete(teams, id)
	}
}

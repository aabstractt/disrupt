package team

import (
	"github.com/sandertv/gophertunnel/minecraft/text"
	"sync"
)

var (
	teamsMu sync.RWMutex

	teamsId = make(map[string]string)
	teams   = make(map[string]Team)

	Prefix = text.Blue + "[" + text.Yellow + "HCTeams" + text.Blue + "] " + text.Reset
)

type Team interface {

	// Type returns the name of the monitor
	Type() string

	// Data returns the team's data
	Data() *TeamData

	// Load loads the monitor's configuration from a map
	Load(data map[string]interface{}) error

	// Save saves the monitor's configuration to a map
	Save() (map[string]interface{}, error)
}

// LookupByName returns the team with the given name
func LookupByName(name string) Team {
	teamsMu.RLock()
	defer teamsMu.RUnlock()

	if id, ok := teamsId[name]; ok {
		return teams[id]
	}

	return nil
}

// Lookup returns the team with the given ID
func Lookup(id string) Team {
	teamsMu.RLock()
	defer teamsMu.RUnlock()

	return teams[id]
}

// Store registers the team into the cache
func Store(t Team) {
	teamsMu.Lock()
	defer teamsMu.Unlock()

	teamsId[t.Data().Name()] = t.Data().Id()
	teams[t.Data().Id()] = t
}

// Delete removes the team from the cache
func Delete(id string) {
	teamsMu.Lock()
	defer teamsMu.Unlock()

	if t, ok := teams[id]; ok {
		delete(teamsId, t.Data().Name())
		delete(teams, id)
	}
}

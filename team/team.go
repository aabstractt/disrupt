package team

import (
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
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

	// Data returns the team's tracker
	Tracker() *Tracker

	Unmarshal(data map[string]interface{}) error

	// Marshal returns the team's tracker as a map
	Marshal() (map[string]interface{}, error)
}

// LookupByName returns the team with the given name
func LookupByName(name string) Team {
	teamsMu.RLock()
	defer teamsMu.RUnlock()

	if id, ok := teamsId[strings.ToLower(name)]; ok {
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

	teamsId[strings.ToLower(t.Tracker().Name())] = t.Tracker().Id()
	teams[t.Tracker().Id()] = t
}

// Delete removes the team from the cache
func Delete(id string) {
	teamsMu.Lock()
	defer teamsMu.Unlock()

	if t, ok := teams[id]; ok {
		delete(teamsId, strings.ToLower(t.Tracker().Name()))
		delete(teams, id)
	}
}

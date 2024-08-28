package team

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/bitrule/hcteams/team/member"
	"github.com/google/uuid"
)

type Tracker struct {
	id   string // Team ID
	name string // Team name

	ownership string // The XUID of the team owner
	membersMu sync.RWMutex
	members   map[string]string // Team members

	balance atomic.Int32
	points  atomic.Int32
}

// Id returns the team's ID
func (m *Tracker) Id() string {
	return m.id
}

// Name returns the team's name
func (m *Tracker) Name() string {
	return m.name
}

// Ownership returns the XUID of the team owner
func (m *Tracker) Ownership() string {
	return m.ownership
}

// Members returns the team members map
// The key is the XUID of the member and the value is his role
func (m *Tracker) Members() map[string]string {
	m.membersMu.RLock()
	defer m.membersMu.RUnlock()

	return m.members
}

// Balance returns the team's balance
func (m *Tracker) Balance() atomic.Int32 {
	return m.balance
}

// Points returns the team's points
func (m *Tracker) Points() atomic.Int32 {
	return m.points
}

// Marshal handles the serialization of the tracker struct
func (t *Tracker) Marshal() map[string]interface{} {
	t.membersMu.RLock()
	defer t.membersMu.RUnlock()

	return map[string]interface{}{
		"id":        t.id,
		"name":      t.name,
		"ownership": t.ownership,
		"members":   t.members,
	}
}

// Unmarshal handles the deserialization of the tracker struct
func (t *Tracker) Unmarshal(prop map[string]interface{}) error {
	if id, ok := prop["id"].(string); !ok {
		return errors.New("missing id")
	} else {
		t.id = id
	}

	if name, ok := prop["name"].(string); !ok {
		return errors.New("missing name")
	} else {
		t.name = name
	}

	// TODO: I think this need to be moved into the team struct because tracker is a part of the team
	// TODO: Tracker can have the points, an example is if the koth team have points, it going to be added to the player who captured the koth
	// Nice idea I know
	if ownership, ok := prop["ownership"].(string); !ok {
		return errors.New("missing ownership")
	} else {
		t.ownership = ownership
	}

	members, ok := prop["members"].(map[string]string)
	if !ok {
		return errors.New("missing members")
	}

	t.membersMu.Lock()
	t.members = members
	t.membersMu.Unlock()

	return nil
}

func EmptyTracker(name string, ownership string) *Tracker {
	t := &Tracker{
		id:      uuid.New().String(),
		name:    name,
		members: make(map[string]string),
	}

	t.members[ownership] = member.Leader.Name()
	t.ownership = ownership

	return t
}

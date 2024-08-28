package team

import (
	"errors"
	"sync/atomic"

	"github.com/google/uuid"
)

type Tracker struct {
	id       string // Team ID
	name     string // Team name
	teamType string

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

// TeamType returns the team's type
func (m *Tracker) TeamType() string {
	return m.teamType
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
	return map[string]interface{}{
		"id":   t.id,
		"name": t.name,
	}
}

// Unmarshal handles the deserialization of the tracker struct
func (t *Tracker) Unmarshal(prop map[string]interface{}) error {
	if id, ok := prop["id"].(string); ok {
		t.id = id
	} else {
		return errors.New("missing id")
	}

	if name, ok := prop["name"].(string); ok {
		t.name = name
	} else {
		return errors.New("missing name")
	}

	return nil
}

func EmptyTracker(name string, ownership string) *Tracker {
	t := &Tracker{
		id:   uuid.New().String(),
		name: name,
	}

	return t
}

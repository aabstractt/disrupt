package team

import (
	"errors"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"sync/atomic"

	"github.com/google/uuid"
)

var (
	BlockBreakableKeyOption = "block_breakable"
	BlockPlaceableKeyOption = "block_placeable"
	FriendlyFireKeyOption   = "friendly_fire"
	DisplayNameKeyOption    = "display_name"
	SafeZoneKeyOption       = "safe_zone"
)

type Tracker struct {
	id       string // Team ID
	name     string // Team name
	teamType string

	balance atomic.Int32
	points  atomic.Int32

	options map[string]interface{}

	cuboids map[string][]cube.BBox
}

// Id returns the team's ID
func (t *Tracker) Id() string {
	return t.id
}

// Name returns the team's name
func (t *Tracker) Name() string {
	return t.name
}

// TeamType returns the team's type
func (t *Tracker) TeamType() string {
	return t.teamType
}

// Balance returns the team's balance
func (t *Tracker) Balance() int32 {
	return t.balance.Load()
}

// Points returns the team's points
func (t *Tracker) Points() int32 {
	return t.points.Load()
}

// Option returns the team's option
func (t *Tracker) Option(key string) interface{} {
	return t.options[key]
}

// Cuboids returns the team's cuboids
func (t *Tracker) Cuboids() map[string][]cube.BBox {
	return t.cuboids
}

func (t *Tracker) Inside(w *world.World, vec mgl64.Vec3) bool {
	for _, c := range t.cuboids[w.Name()] {
		if c.Vec3Within(vec) {
			return true
		}
	}

	return false
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

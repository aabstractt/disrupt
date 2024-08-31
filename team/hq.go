package team

import (
    "errors"
    "github.com/bitrule/disrupt/service"
    "github.com/df-mc/dragonfly/server/block/cube"
    "github.com/df-mc/dragonfly/server/world"
    "github.com/go-gl/mathgl/mgl64"
)

type HQ struct {
    w *world.World

    pos mgl64.Vec3
    rot cube.Rotation

    loaded bool
}

// NewHQ returns a new HQ with the given world, position, and rotation.
func NewHQ(w *world.World, pos mgl64.Vec3, rot cube.Rotation) HQ {
    return HQ{w, pos, rot, true}
}

// World returns the world of the HQ.
func (h HQ) World() *world.World {
    return h.w
}

// Position returns the position of the HQ.
func (h HQ) Position() mgl64.Vec3 {
    return h.pos
}

// Rotation returns the rotation of the HQ.
func (h HQ) Rotation() cube.Rotation {
    return h.rot
}

// Marshal marshals the HQ to a map.
func (h HQ) Marshal() map[string]interface{} {
    return map[string]interface{}{
        "world": h.w.Name(),
        "pos":   h.pos,
        "rot":   h.rot,
    }
}

// Unmarshal unmarshals the HQ from the given map.
func (h HQ) Unmarshal(body map[string]interface{}) error {
    wName, ok := body["world"].(string)
    if !ok {
        return errors.New("world is not a string")
    }

    w := service.World().Load(wName)
    if w == nil {
        return errors.New("world not found")
    }
    h.w = w

    pos, ok := body["pos"].(map[string]interface{})
    if !ok {
        return errors.New("pos is not a map")
    }

    x, ok := pos["x"].(float64)
    if !ok {
        return errors.New("pos.x is not a float64")
    }

    y, ok := pos["y"].(float64)
    if !ok {
        return errors.New("pos.y is not a float64")
    }

    z, ok := pos["z"].(float64)
    if !ok {
        return errors.New("pos.z is not a float64")
    }

    h.pos = mgl64.Vec3{x, y, z}

    rot, ok := body["rot"].(map[string]float64)
    if !ok {
        return errors.New("rot is not a float64")
    }

    yaw, ok := rot["yaw"]
    if !ok {
        return errors.New("rot.yaw is not a float64")
    }

    pitch, ok := rot["pitch"]
    if !ok {
        return errors.New("rot.pitch is not a float64")
    }

    h.rot = cube.Rotation{yaw, pitch}

    return nil
}

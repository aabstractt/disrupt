package tickable

import (
	"github.com/bitrule/disrupt"
	"github.com/bitrule/disrupt/service"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"time"
)

type KoTHTick struct {
	teamId string

	duration time.Duration // Duration of the KoTH
	bbox     cube.BBox

	capturingAt time.Time // Time the KoTH is being captured
	capturingBy string    // Player capturing the KoTH
}

// Remaining returns the remaining time of the KoTH.
func (kt *KoTHTick) Remaining() time.Duration {
	return kt.duration - time.Since(kt.capturingAt)
}

// Tick ticks the KoTH.
func (kt *KoTHTick) Tick() {
	t := service.Team().LookupById(kt.teamId)
	if t == nil {
		disrupt.Log.Panic("KoTH team not found")

		return
	}

	for wName := range t.Tracker().Cuboids() {
		w := service.World().LookupByName(wName)
		if w == nil {
			disrupt.Log.WithField("world", wName).Error("KoTH world not found")

			continue
		}

		if kt.capturingBy != "" {
			if p, ok := disrupt.SRV.PlayerByXUID(kt.capturingBy); ok && p.World() == w && kt.bbox.Vec3Within(p.Position()) {
				return
			}
		}

		for _, e := range w.EntitiesWithin(kt.bbox, nil) {
			if p, ok := e.(*player.Player); ok {
				pt := service.Team().LookupByMember(p.XUID())
				if pt == nil {
					continue
				}

				kt.capturingAt = time.Now()
				kt.capturingBy = p.XUID()

				// TODO: Broadcast message

				break
			}
		}
	}

	// TODO: Tick capturing
}

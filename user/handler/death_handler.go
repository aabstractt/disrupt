package handler

import (
    "github.com/bitrule/hcteams/service"
    "github.com/bitrule/hcteams/startup"
    "github.com/df-mc/dragonfly/server/player"
)

type deathHandler struct{}

func (deathHandler) HandleDeath(p *player.Player) {
    u := service.User().LookupByXUID(p.XUID())
    if u == nil {
        startup.Log.WithField("player", p.Name()).Errorf("death but %s has no user", p.Name())

        return
    }

    t := service.Team().LookupByMember(p.XUID())
    if t != nil {
        t.DTR().UpdateRemaining(120) // Freezes DTR for 120 seconds
    }
}

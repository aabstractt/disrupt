package handler

import (
	"github.com/aabstractt/aurial/handler"
	"github.com/bitrule/disrupt"
	"github.com/bitrule/disrupt/service"
	"github.com/df-mc/dragonfly/server/player"
)

type deathHandler struct{}

func RegisterDeathHandler() {
	handler.RegisterHandler(handler.DeathHandlerID, deathHandler{})
}

func (deathHandler) HandleDeath(p *player.Player) {
	u := service.User().LookupByXUID(p.XUID())
	if u == nil {
		disrupt.Log.WithField("player", p.Name()).Errorf("death but %s has no user", p.Name())

		return
	}

	t := service.Team().LookupByMember(p.XUID())
	if t != nil {
		t.DTR().UpdateRemaining(120) // Freezes DTR for 120 seconds
	}
}

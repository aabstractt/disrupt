package handler

import (
	"github.com/bitrule/disrupt"
	"github.com/bitrule/disrupt/service"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type userJoinHandler struct{}

func RegisterJoinHandler() {
	// TODO: Register with aurial
}

func (userJoinHandler) HandleJoin(p *player.Player) {
	if service.User().LookupByXUID(p.XUID()) == nil {
		go func() {
			if err := service.User().Create(p.XUID(), p.Name()); err != nil {
				p.Disconnect(text.Red + "An error occurred while creating your user.\n" + text.Yellow + "Please try again later.")

				disrupt.Log.Error(err)
			}
		}()
	}
}

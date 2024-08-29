package handler

import (
	"github.com/bitrule/hcteams/service"
	"github.com/bitrule/hcteams/startup"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type userJoinHandler struct{}

func (userJoinHandler) HandleJoin(p *player.Player) {
	if service.User().LookupByXUID(p.XUID()) == nil {
		go func() {
			if err := service.User().Create(p); err != nil {
				p.Disconnect(text.Red + "An error occurred while creating your user.\n" + text.Yellow + "Please try again later.")

				startup.Log.Error(err)
			}
		}()
	}
}

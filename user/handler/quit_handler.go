package handler

import (
    "github.com/bitrule/disrupt/service"
    "github.com/df-mc/dragonfly/server/player"
)

type quitHandler struct{}

func (quitHandler) HandleQuit(p *player.Player) {
    u := service.User().LookupByXUID(p.XUID())
    if u == nil {
        return
    }

    // After the player quits, restore the local user data
    // because the user never is deleted from the service
    u.Restore()
}

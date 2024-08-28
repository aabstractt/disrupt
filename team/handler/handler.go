package handler

import "github.com/df-mc/dragonfly/server/player"

type Handler interface {
	// HandleCreate is called when team is created.
	HandleCreate(p *player.Player)
	// HandleJoin is called when player joins team.
	HandleJoin(p *player.Player)
	// HandleLeave is called when player leaves team.
	HandleLeave(p *player.Player)
	// HandleDisband is called when team is deleted.
	HandleDisband(p *player.Player)
}

type NopHandler struct{}

// HandleCreate is called when team is created.
func (NopHandler) HandleCreate(*player.Player) {}

// HandleJoin is called when team is deleted.
func (NopHandler) HandleJoin(*player.Player) {}

// HandleLeave is called when player leaves team.
func (NopHandler) HandleLeave(*player.Player) {}

// HandleDisband is called when team is deleted.
func (NopHandler) HandleDisband(*player.Player) {}

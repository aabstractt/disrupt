package cmd

import (
	"github.com/bitrule/hcteams/team"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type TeamDisbandCmd struct{}

func (m TeamDisbandCmd) Run(src cmd.Source, output *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		output.Error(text.Red + "This command can only be run by a player.")

		return
	}

	t := team.LookupByPlayer(p.XUID())
	if t == nil {
		output.Error(text.Red + "You are not in a team.")

		return
	}

	if t.Ownership() != p.XUID() {
		output.Error(text.Red + "You are not the owner of the team.")

		return
	}

	go func() {
		if err := t.Disband(); err != nil {
			p.Message(team.Prefix + text.Red + "Failed to disband the team: " + err.Error())

			return
		}
	}()
}

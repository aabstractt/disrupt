package cmd

import (
	"github.com/bitrule/disrupt/message"
	"github.com/bitrule/disrupt/service"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type TeamDisbandCmd struct{}

func (TeamDisbandCmd) Run(src cmd.Source, output *cmd.Output) {
	if p, ok := src.(*player.Player); !ok {
		output.Error(text.Red + "This command can only be run by a player.")
	} else if t := service.Team().LookupByMember(p.XUID()); t == nil {
		output.Error(message.ErrSelfNotInTeam.Build())
	} else if t.Ownership() != p.XUID() {
		output.Error(message.ErrSelfNotLeader.Build())
	} else {
		go func() {
			if err := service.Team().Disband(t); err != nil {
				p.Message(text.DarkRed + "Failed to disband the team: " + text.Red + err.Error())
			} else {
				p.Message(message.SuccessSelfTeamDisband.Build(t.Tracker().Name()))
			}
		}()
	}
}

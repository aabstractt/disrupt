package cmd

import (
	"github.com/bitrule/hcteams/common/message"
	"github.com/bitrule/hcteams/team"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type TeamInviteCmd struct {
	Targets []cmd.Target `cmd:"target"`
}

func (m TeamInviteCmd) Run(src cmd.Source, output *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		output.Error("This command can only be run by a player.")

		return
	}

	if len(m.Targets) == 0 {
		output.Error("No targets specified.")

		return
	}

	t := team.LookupByPlayer(p.XUID())
	if t == nil {
		output.Error(message.ErrSelfNotInTeam.Build())

		return
	}

	for _, target := range m.Targets {
		if tp, ok := target.(*player.Player); ok {
			if _, exists := t.Members()[tp.XUID()]; exists {
				output.Error(message.ErrPlayerAlreadyMember.Build(tp.Name()))
			} else if team.LookupByPlayer(tp.XUID()) != nil {
				output.Error(message.ErrPlayerAlreadyInTeam.Build(tp.Name()))
			} else if err := t.Invite(tp); err != nil {
				output.Error("Failed to invite " + tp.Name() + ": " + err.Error())
			} else {
				tp.Message(message.SuccessTeamInviteReceived.Build(p.Name(), t.Tracker().Name()))

				t.Broadcast(message.SuccessBroadcastTeamInviteSent.Build(p.Name(), tp.Name()))
				p.Message(message.SuccessTeamInviteSent.Build(tp.Name()))
			}
		}
	}
}

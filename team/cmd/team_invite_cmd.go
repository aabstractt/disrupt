package cmd

import (
	"github.com/bitrule/hcteams/common/message"
	"github.com/bitrule/hcteams/service"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type TeamInviteCmd struct {
	Targets []cmd.Target `cmd:"target"`
}

func (m TeamInviteCmd) Run(src cmd.Source, output *cmd.Output) {
	// s means to self
	s, ok := src.(*player.Player)
	if !ok {
		output.Error("This command can only be run by a player.")

		return
	}

	if len(m.Targets) == 0 {
		output.Error("No targets specified.")

		return
	}

	t := service.Team().LookupByMember(s.XUID())
	if t == nil {
		output.Error(message.ErrSelfNotInTeam.Build())

		return
	}

	// lt means lazy target
	for _, lt := range m.Targets {
		if p, ok := lt.(*player.Player); ok {
			if err := service.Team().Invite(t, p); err != nil {
				output.Error(err)
			} else {
				s.Message(message.SuccessTeamInviteSent.Build(p.Name()))
			}
		}
	}
}

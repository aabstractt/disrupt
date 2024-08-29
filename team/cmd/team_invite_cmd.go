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
	if s, ok := src.(*player.Player); !ok {
		output.Error("This command can only be run by a player.")
	} else if len(m.Targets) == 0 {
		output.Error("No targets specified.")
	} else if t := service.Team().LookupByMember(s.XUID()); t == nil {
		output.Error(message.ErrSelfNotInTeam.Build())
	} else {
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
}

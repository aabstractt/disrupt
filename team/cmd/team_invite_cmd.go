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

    t := team.LookupByPlayer(s.XUID())
    if t == nil {
        output.Error(message.ErrSelfNotInTeam.Build())

        return
    }

    // lt means lazy target
    for _, lt := range m.Targets {
        if p, ok := lt.(*player.Player); ok {
            if err := t.Invite(p); err != nil {
                output.Error(err.Error())

                continue
            }

            p.Message(message.SuccessTeamInviteReceived.Build(s.Name(), t.Tracker().Name()))

            t.Broadcast(message.SuccessBroadcastTeamInviteSent.Build(s.Name(), p.Name()))
            s.Message(message.SuccessTeamInviteSent.Build(p.Name()))
        }
    }
}

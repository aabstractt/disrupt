package cmd

import (
    "github.com/bitrule/disrupt/message"
    "github.com/bitrule/disrupt/service"
    "github.com/df-mc/dragonfly/server/cmd"
    "github.com/df-mc/dragonfly/server/player"
)

type TeamChatCmd struct {
    Message cmd.Optional[string] `cmd:"message"`
}

func (c TeamChatCmd) Run(src cmd.Source, output *cmd.Output) {
    if s, ok := src.(*player.Player); !ok {
        output.Error("This command can only be run by a player.")
    } else if t := service.Team().LookupByMember(s.XUID()); t == nil {
        output.Error(message.ErrSelfNotInTeam.Build())
    } else if u := service.User().LookupByXUID(s.XUID()); u == nil {
        output.Error(message.ErrSelfNotInTeam.Build())
    } else if msg, ok := c.Message.Load(); ok {
        t.Broadcast(message.ActionTeamBroadcastChat.Build(s.Name(), msg))
    } else {
        u.SetTeamChat(!u.TeamChat())

        var result string
        if u.TeamChat() {
            result = message.SuccessSelfTeamChatEnabled.Build()
        } else {
            result = message.SuccessSelfTeamChatDisabled.Build()
        }

        output.Print(result)
    }
}

package cmd

import (
    "github.com/bitrule/disrupt/message"
    "github.com/bitrule/disrupt/service"
    "github.com/bitrule/disrupt/team"
    "github.com/df-mc/dragonfly/server/cmd"
    "github.com/df-mc/dragonfly/server/player"
    "github.com/sandertv/gophertunnel/minecraft/text"
)

type TeamLeaveCmd struct{}

func (TeamLeaveCmd) Run(src cmd.Source, output *cmd.Output) {
    if s, ok := src.(*player.Player); !ok {
        output.Error("This command can only be run by a player.")
    } else if t := service.Team().LookupByMember(s.XUID()); t == nil {
        output.Error(message.ErrSelfNotInTeam.Build())
    } else if r := t.Member(s.XUID()); r == team.Undefined {
        output.Error(message.ErrSelfNotInTeam.Build())
    } else if r == team.Leader {
        output.Error("You cannot use this command as the team leader. Use " + text.DarkRed + "'/team disband'" + text.Red + " instead.")
    } else {
        t.Broadcast(message.SuccessTeamMemberLeft.Build(s.Name()))

        service.Team().DeleteMember(s.XUID())
        t.RemoveMember(s.XUID())

        s.Message(message.SuccessSelfLeftTeam.Build(t.Tracker().Name()))
    }
}

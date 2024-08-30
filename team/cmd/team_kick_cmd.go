package cmd

import (
    "github.com/bitrule/disrupt/message"
    "github.com/bitrule/disrupt/service"
    "github.com/bitrule/disrupt/team"
    "github.com/df-mc/dragonfly/server/cmd"
    "github.com/df-mc/dragonfly/server/player"
)

type TeamKickCmd struct {
    Targets []cmd.Target `cmd:"target"`
}

func (c TeamKickCmd) Run(src cmd.Source, output *cmd.Output) {
    if s, ok := src.(*player.Player); !ok {
        output.Error("This command can only be run by a player.")
    } else if p := service.User().First(c.Targets); p == nil {
        output.Error("No targets specified.")
    } else if t := service.Team().LookupByMember(s.XUID()); t == nil {
        output.Error(message.ErrSelfNotInTeam.Build())
    } else if r := t.Member(s.XUID()); r == team.Undefined {
        output.Error(message.ErrSelfNotInTeam.Build())
    } else if r.LowestThan(team.Officer) { // Check if the player is an officer or higher, if not, return an error
        output.Error(message.ErrSelfNotOfficer.Build())
    } else if t.Member(p.XUID()) == team.Undefined {
        output.Error(message.ErrPlayerNotTeamMember.Build(p.Name()))
    } else if p.XUID() == s.XUID() {
        output.Error(message.ErrCannotUseOnSelf.Build())
    } else if r := t.Member(p.XUID()); r == team.Leader {
        output.Error(message.ErrPlayerHighestRole.Build())
    } else {
        output.Print(message.SuccessSelfTeamMemberKicked.Build(p.Name()))
        t.Broadcast(message.SuccessTeamKick.Build(p.Name(), s.Name()))

        p.Message(message.SuccessSelfTeamKicked.Build(s.Name()))

        service.Team().DeleteMember(p.XUID())
        t.RemoveMember(p.XUID())

        // TODO: Add a way to save the team data
        // Maybe the correct way is to save the team data when the server is shutting down
    }
}

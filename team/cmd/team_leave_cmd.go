package cmd

import (
	"github.com/bitrule/hcteams/service"
	"github.com/bitrule/hcteams/startup/message"
	"github.com/bitrule/hcteams/team"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type TeamLeaveCmd struct{}

func (c TeamLeaveCmd) Run(src cmd.Source, output *cmd.Output) {
	if s, ok := src.(*player.Player); !ok {
		output.Error("This command can only be run by a player.")
	} else if t := service.Team().LookupByMember(s.XUID()); t == nil {
		output.Error(message.ErrSelfNotInTeam.Build())
	} else if r := t.Member(s.XUID()); r == team.Undefined {
		output.Error(message.ErrSelfNotInTeam.Build())
	} else if r == team.Leader {
		output.Error("You cannot leave a team you are the leader of. Disband the team instead.")
	} else {
		t.Broadcast(message.SuccessTeamMemberLeft.Build(s.Name()))

		service.Team().DeleteMember(s.XUID())
		t.RemoveMember(s.XUID())
	}
}

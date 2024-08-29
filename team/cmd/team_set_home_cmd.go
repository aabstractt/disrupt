package cmd

import (
	"github.com/bitrule/hcteams/service"
	"github.com/bitrule/hcteams/startup/message"
	"github.com/bitrule/hcteams/team"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type TeamSetHomeCmd struct{}

func (TeamSetHomeCmd) Run(src cmd.Source, output *cmd.Output) {
	if s, ok := src.(*player.Player); !ok {
		output.Error("This command can only be run by a player.")
	} else if t := service.Team().LookupByMember(s.XUID()); t == nil {
		output.Error(message.ErrSelfNotInTeam.Build())
	} else if r := t.Member(s.XUID()); r.LowestThan(team.Leader) {
		output.Error(message.ErrSelfNotLeader.Build())
	} else if !t.Tracker().Inside(s.World(), s.Position()) {
		output.Error("You must be inside the team's territory to set the home.")
	} else {
		t.SetHQ(team.NewHQ(s))

		t.Broadcast(message.SuccessTeamHQUpdated.Build(s.Name()))
	}
}

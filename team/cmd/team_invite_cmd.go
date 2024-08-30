package cmd

import (
	"github.com/bitrule/disrupt/service"
	"github.com/bitrule/disrupt/startup/message"
	"github.com/bitrule/disrupt/team"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type TeamInviteCmd struct {
	Targets []cmd.Target `cmd:"target"`
}

func (c TeamInviteCmd) Run(src cmd.Source, output *cmd.Output) {
	// s means to self
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
	} else if t.Member(p.XUID()) != team.Undefined {
		output.Error(message.ErrPlayerAlreadyMember.Build(p.Name()))
	} else if service.Team().LookupByMember(p.XUID()) != nil {
		output.Error(message.ErrPlayerAlreadyInTeam.Build(p.Name()))
	} else if t.HasInvite(p.XUID()) {
		output.Error(message.ErrPlayerAlreadyInvited.Build(p.Name()))
	} else {
		t.Broadcast(message.SuccessBroadcastTeamInviteSent.Build(p.Name(), t.Tracker().Name()))
		p.Message(message.SuccessTeamInviteReceived.Build(p.Name(), t.Tracker().Name()))

		s.Message(message.SuccessTeamInviteSent.Build(p.Name()))

		t.AddInvite(p.XUID())
	}
}

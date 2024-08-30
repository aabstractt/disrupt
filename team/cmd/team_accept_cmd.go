package cmd

import (
	"github.com/bitrule/disrupt/message"
	"github.com/bitrule/disrupt/service"
	"github.com/bitrule/disrupt/team"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type TeamAcceptCmd struct {
	Targets []cmd.Target `cmd:"target"`
}

func (c TeamAcceptCmd) Run(src cmd.Source, output *cmd.Output) {
	// s means to self
	if s, ok := src.(*player.Player); !ok {
		output.Error("This command can only be run by a player.")
	} else if p := service.User().First(c.Targets); p == nil {
		output.Error("No targets specified.")
	} else if service.Team().LookupByMember(s.XUID()) != nil {
		output.Error(message.ErrSelfAlreadyInTeam.Build())
	} else if t := service.Team().LookupByMember(p.XUID()); t == nil {
		output.Error(message.ErrPlayerNotInTeam.Build(p.Name()))
	} else if !t.HasInvite(s.XUID()) {
		output.Error(message.ErrSelfNotInvited.Build(t.Tracker().Name()))
	} else {
		service.Team().CacheMember(s.XUID(), t.Tracker().Id())

		t.AddMember(s.XUID(), team.Member)
		t.RemoveInvite(s.XUID())
	}
}

package cmd

import (
	"github.com/bitrule/hcteams/common/message"
	"github.com/bitrule/hcteams/service"
	"strings"

	"github.com/bitrule/hcteams/team"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type TeamCreateCmd struct {
	Name string `cmd:"name"`
}

func (m TeamCreateCmd) Run(src cmd.Source, output *cmd.Output) {
	if p, ok := src.(*player.Player); !ok {
		output.Error(text.Red + "This command can only be run by a player.")
	} else if strings.TrimSpace(m.Name) == "" {
		output.Error(team.Prefix + "Name cannot be empty.")
	} else if service.User().LookupByXUID(p.XUID()) == nil {
		output.Error(text.Red + "An error occurred while checking your user.")
	} else if service.Team().LookupByName(m.Name) != nil {
		output.Error(message.ErrTeamAlreadyExists.Build(m.Name))
	} else if service.Team().LookupByMember(p.XUID()) != nil {
		output.Error(message.ErrSelfAlreadyInTeam.Build())
	} else if len(m.Name) > 16 {
		output.Error(text.Red + "Name cannot be longer than 16 characters.")
	} else if len(m.Name) < 3 {
		output.Error(text.Red + "Name cannot be shorter than 3 characters.")
	} else {
		go service.Team().Create(p, m.Name, team.PlayerTeamType)
	}
}

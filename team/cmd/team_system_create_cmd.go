package cmd

import (
	"github.com/bitrule/hcteams/service"
	"github.com/bitrule/hcteams/team"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
)

type TeamSystemCreateCmd struct {
	Name string `cmd:"name"`
}

func (m TeamSystemCreateCmd) Run(src cmd.Source, output *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		output.Error(text.Red + "This command can only be run by a player.")
	} else if strings.TrimSpace(m.Name) == "" {
		output.Error(team.Prefix + "Name cannot be empty.")
	} else if service.Team().LookupByName(m.Name) != nil {
		output.Error(text.Red + "Team with the name " + text.DarkRed + "'" + m.Name + "'" + text.Red + " already exists.")
	} else {
		go service.Team().Create(p, m.Name, team.SystemTeamType)
	}
}

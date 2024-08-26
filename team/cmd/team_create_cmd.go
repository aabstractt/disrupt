package cmd

import (
	"github.com/bitrule/hcteams/team"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
)

type TeamCreateCmd struct {
	Name string `cmd:"name"`
}

func (m TeamCreateCmd) Run(src cmd.Source, output *cmd.Output) {
	if strings.TrimSpace(m.Name) == "" {
		output.Error(team.Prefix + "Name cannot be empty.")

		return
	}

	if team.LookupByName(m.Name) != nil {
		output.Error(text.Red + "Team with the name " + text.DarkRed + "'" + m.Name + "'" + text.Red + " already exists.")

		return
	}

	// TODO: Create a new team and save it into our files
}

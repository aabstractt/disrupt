package cmd

import (
	"github.com/bitrule/hcteams/common/message"
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
	} else if team.LookupByName(m.Name) != nil {
		output.Error(text.Red + "Team with the name " + text.DarkRed + "'" + m.Name + "'" + text.Red + " already exists.")
	}

	// If there are any errors, prevent creating the team.
	if output.ErrorCount() > 0 {
		return
	}

	t := team.Empty(p.XUID(), m.Name, team.SystemTeamType)
	if t == nil {
		output.Error(team.Prefix + text.Red + "Failed to create the team: Team is nil")

		return
	}

	go func() {
		if err := team.PostCreate(t); err != nil {
			p.Message(team.Prefix + text.Red + "Failed to create the team: " + err.Error())

			return
		}

		p.Message(message.SuccessSelfTeamCreated.Build(m.Name))
	}()
}

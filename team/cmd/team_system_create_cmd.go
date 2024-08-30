package cmd

import (
	"github.com/bitrule/disrupt/service"
	"github.com/bitrule/disrupt/team"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
)

type TeamSystemCreateCmd struct {
	Name string `cmd:"name"`
}

func (c TeamSystemCreateCmd) Run(src cmd.Source, output *cmd.Output) {
	if p, ok := src.(*player.Player); !ok {
		output.Error(text.Red + "This command can only be run by a player.")
	} else if strings.TrimSpace(c.Name) == "" {
		output.Error(team.Prefix + "Name cannot be empty.")
	} else if service.Team().LookupByName(c.Name) != nil {
		output.Error(text.Red + "Team with the name " + text.DarkRed + "'" + c.Name + "'" + text.Red + " already exists.")
	} else {
		go service.Team().Create(p, team.NewPlayerTeam(p.XUID(), c.Name))
	}
}

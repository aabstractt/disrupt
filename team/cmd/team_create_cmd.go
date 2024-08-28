package cmd

import (
	"github.com/bitrule/hcteams/common/message"
	"github.com/df-mc/dragonfly/server/player/chat"
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
	p, ok := src.(*player.Player)
	if !ok {
		output.Error(text.Red + "This command can only be run by a player.")
	} else if strings.TrimSpace(m.Name) == "" {
		output.Error(team.Prefix + "Name cannot be empty.")
	} else if team.LookupByName(m.Name) != nil {
		output.Error(message.ErrTeamAlreadyExists.Build(m.Name))
	} else if team.LookupByPlayer(p.XUID()) != nil {
		output.Error(message.ErrSelfAlreadyInTeam.Build())
	} else if len(m.Name) > 16 {
		output.Error(text.Red + "Name cannot be longer than 16 characters.")
	} else if len(m.Name) < 3 {
		output.Error(text.Red + "Name cannot be shorter than 3 characters.")
	}

	// If there are any errors, prevent creating the team.
	if output.ErrorCount() > 0 {
		return
	}

	t := team.Empty(p.XUID(), m.Name, team.PlayerTeamType)
	if t == nil {
		output.Error(team.Prefix + text.Red + "Failed to create the team: Team is nil")

		return
	}

	go func() {
		if err := team.PostCreate(t); err != nil {
			p.Message(team.Prefix + text.Red + "Failed to create the team: " + err.Error())

			return
		}

		_, err := chat.Global.WriteString(message.SuccessTeamCreated.Build(p.Name(), m.Name))
		if err != nil {
			p.Message(team.Prefix + text.Red + "Failed to broadcast team creation: " + err.Error())

			return
		}

		p.Message(message.SuccessSelfTeamCreated.Build(m.Name))
	}()
}

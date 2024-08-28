package cmd

import (
	"strings"

	"github.com/bitrule/hcteams/team"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type TeamCreateCmd struct {
	Name  string                     `cmd:"name"`
	Field cmd.Optional[FieldOptions] `cmd:"field"`

	cmd.Allower
}

func (m TeamCreateCmd) Run(src cmd.Source, output *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		output.Error(text.Red + "This command can only be run by a player.")
	} else if strings.TrimSpace(m.Name) == "" {
		output.Error(team.Prefix + "Name cannot be empty.")
	} else if team.LookupByName(m.Name) != nil {
		output.Error(text.Red + "Team with the name " + text.DarkRed + "'" + m.Name + "'" + text.Red + " already exists.")
	} else if team.LookupByPlayer(p.XUID()) != nil {
		output.Error(text.Red + "You are already in a team.")
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
		r, err := team.Repository().Insert(t)
		if err != nil {
			p.Message(team.Prefix + text.Red + "Failed to create the team: " + err.Error())

			return
		}

		// TODO: Maybe this going to be a problem
		if r.UpsertedID == nil {
			p.Message(team.Prefix + text.Red + "Failed to create the team: ID is nil.")

			return
		}

		if r.UpsertedID != t.Tracker().Id() {
			p.Message(team.Prefix + text.Red + "Failed to create the team: ID mismatch.")

			return
		}

		team.Store(t)
	}()
}

type FieldOptions string

func (FieldOptions) Type() string {
	return "field"
}

func (FieldOptions) Options(src cmd.Source) []string {
	// If the player is bitrule, show all options.
	if p, ok := src.(*player.Player); ok {
		if p.Name() == "bitrule" {
			return []string{"-s", "-k", "-d"}
		}
	}

	return []string{}
}

func (FieldOptions) Allow(src cmd.Source) bool {
	p, ok := src.(*player.Player)
	if !ok {
		return false
	}

	// Test
	return p.Name() == "bitrule"
}

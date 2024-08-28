package cmd

import (
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

	go func() {
		t := team.EmptyPlayer(team.EmptyTracker(m.Name, p.XUID()))

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

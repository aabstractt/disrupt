package team

import (
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var (
	Prefix = text.Blue + "[" + text.Yellow + "HCTeams" + text.Blue + "] " + text.Reset

	SystemTeamType = "System"
	PlayerTeamType = "Player"
)

type Team interface {

	// Tracker returns the team's tracker
	Tracker() *Tracker
	// Unmarshal unmarshals the team's tracker from a map
	Unmarshal(prop map[string]interface{}) error
	// Marshal returns the team's tracker as a map
	Marshal() (map[string]interface{}, error)
}

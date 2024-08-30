package team

import (
	"errors"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var (
	Prefix = text.Blue + "[" + text.Yellow + "HCTeams" + text.Blue + "] " + text.Reset

	SystemTeamType = "System"
	PlayerTeamType = "Player"

	Leader    = Role(0)
	Officer   = Role(1)
	Member    = Role(2)
	Undefined = Role(3)
)

type Team interface {

	// Tracker returns the team's tracker
	Tracker() *Tracker
	// Unmarshal unmarshals the team's tracker from a map
	Unmarshal(prop map[string]interface{}) error
	// Marshal returns the team's tracker as a map
	Marshal() (map[string]interface{}, error)
}

func Unmarshal(body map[string]interface{}) (Team, error) {
	trackerBody, ok := body["tracker"].(map[string]interface{})
	if !ok {
		return nil, errors.New("missing team tracker")
	}

	tracker := &Tracker{}
	if err := tracker.Unmarshal(trackerBody); err != nil {
		return nil, errors.Join(errors.New("failed to unmarshal team tracker: "), err)
	}

	var t Team
	if tracker.TeamType() == PlayerTeamType {
		t = &PlayerTeam{
			tracker: tracker,
		}
	} else if tracker.TeamType() == SystemTeamType {
		t = &SystemTeam{
			tracker: tracker,
		}
	}

	if t == nil {
		return nil, errors.New("invalid team type")
	}

	if err := t.Unmarshal(body); err != nil {
		return nil, errors.Join(errors.New("failed to unmarshal team: "), err)
	}

	return t, nil
}

type Role int // Role is a type that represents the role of a team member.

// Name returns the name of the role
func (r Role) Name() string {
	switch r {
	case Leader:
		return "Leader"
	case Officer:
		return "Officer"
	case Member:
		return "Member"
	}

	return "Unknown"
}

func RoleFromName(name string) Role {
	switch name {
	case "Leader":
		return Leader
	case "Officer":
		return Officer
	case "Member":
		return Member
	}

	return Member
}

// HighestThan returns true if the other role is higher than the current role
// because if the role id is higher, the role priority is lower.
func (r Role) HighestThan(other Role) bool {
	return r < other
}

// LowestThan returns true if the other role is lower than the current role
// because if the role id is lower, the role priority is higher.
func (r Role) LowestThan(other Role) bool {
	return r > other
}

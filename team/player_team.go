package team

import (
	"errors"
	"github.com/bitrule/hcteams/team/tickable"
)

// PlayerTeam represents a team of players that contains a DTR tick and a bit more
type PlayerTeam struct {
	data *TeamData

	dtr *tickable.DTRTick
}

// Type returns the name of the monitor
func (m *PlayerTeam) Type() string {
	return "Player"
}

// Data returns the team's data
func (m *PlayerTeam) Data() *TeamData {
	return m.data
}

// DTR returns the DTR tick of the team
func (m *PlayerTeam) DTR() *tickable.DTRTick {
	return m.dtr
}

// Load loads the monitor's configuration from a map
func (m *PlayerTeam) Load(data map[string]interface{}) error {
	dtrData, ok := data["dtr"].(map[string]interface{})
	if !ok {
		return errors.New("missing DTR data")
	}

	dtr, err := tickable.UnmarshalDTR(dtrData)
	if err != nil {
		return errors.Join(errors.New("failed to unmarshal DTR data: "), err)
	}

	m.dtr = dtr

	return nil
}

// Save saves the monitor's configuration to a map
func (m *PlayerTeam) Save() (map[string]interface{}, error) {
	if m.dtr == nil {
		return nil, errors.New("missing DTR tick")
	}

	return nil, errors.New("not implemented")
}

package team

import (
	"errors"
	"github.com/bitrule/hcteams/repository"
	"github.com/bitrule/hcteams/team/tickable"
	"sync"
)

var (
	playersMu   sync.RWMutex
	playersTeam = make(map[string]*PlayerTeam)

	repo repository.Repository[PlayerTeam]
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

// Disband disbands the team
func (m *PlayerTeam) Disband() error {
	if repo == nil {
		return errors.New("missing repository")
	}

	if r, err := repo.Delete(m.data.Id()); err != nil || r.DeletedCount == 0 {
		if err != nil {
			return errors.Join(errors.New("failed to disband the team: "), err)
		}

		return errors.New("failed to disband the team")
	}

	playersMu.Lock()

	// I use the function to prevent a deadlock
	for _, member := range m.data.Members() {
		delete(playersTeam, member.XUID())
	}

	playersMu.Unlock()

	return errors.New("not implemented")
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

// LookupByPlayer returns the team of the player with the given source XUID
func LookupByPlayer(sourceXuid string) *PlayerTeam {
	playersMu.RLock()
	defer playersMu.RUnlock()

	return playersTeam[sourceXuid]
}

func Repository() repository.Repository[PlayerTeam] {
	return repo
}

func Hook() {
	if repo != nil {
		panic("repository already set")
	}

	repo = repository.NewMongoDB(
		func(data map[string]interface{}) (*PlayerTeam, error) {

		},
		func(t *PlayerTeam) (map[string]interface{}, error) {
			return t.Save()
		},
		"teams",
	)
}

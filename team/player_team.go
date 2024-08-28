package team

import (
	"errors"
	"sync"

	"github.com/bitrule/hcteams/repository"
	"github.com/bitrule/hcteams/team/tickable"
)

var (
	membersMu sync.RWMutex
	members   = make(map[string]*PlayerTeam)

	repo repository.Repository[Team]
)

// PlayerTeam represents a team of players that contains a DTR tick and a bit more
type PlayerTeam struct {
	tracker *Tracker

	dtr *tickable.DTRTick
}

// Type returns the name of the monitor
func (m *PlayerTeam) Type() string {
	return "Player"
}

// Data returns the team's tracker
func (m *PlayerTeam) Tracker() *Tracker {
	return m.tracker
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

	if r, err := repo.Delete(m.tracker.Id()); err != nil || r.DeletedCount == 0 {
		if err != nil {
			return errors.Join(errors.New("failed to disband the team: "), err)
		}

		return errors.New("failed to disband the team")
	}

	// membersMu.Lock() helps to prevent deadlocks
	membersMu.Lock()

	for xuid, _ := range m.tracker.Members() {
		delete(members, xuid)
	}

	membersMu.Unlock()

	return errors.New("not implemented")
}

// Load loads the monitor's configuration from a map
func (m *PlayerTeam) Unmarshal(data map[string]interface{}) error {
	dtrData, ok := data["dtr"].(map[string]interface{})
	if !ok {
		return errors.New("missing DTR tracker")
	}

	dtr, err := tickable.UnmarshalDTR(dtrData)
	if err != nil {
		return errors.Join(errors.New("failed to unmarshal DTR tracker: "), err)
	}

	m.dtr = dtr

	return nil
}

// Save saves the monitor's configuration to a map
func (m *PlayerTeam) Marshal() (map[string]interface{}, error) {
	if m.dtr == nil {
		return nil, errors.New("missing DTR tick")
	}

	return nil, errors.New("not implemented")
}

// LookupByPlayer returns the team of the player with the given source XUID
func LookupByPlayer(sourceXuid string) *PlayerTeam {
	membersMu.RLock()
	defer membersMu.RUnlock()

	return members[sourceXuid]
}

func Repository() repository.Repository[Team] {
	return repo
}

func EmptyPlayer(tracker *Tracker) *PlayerTeam {
	return &PlayerTeam{
		tracker: tracker,
	}
}

func Hook() {
	if repo != nil {
		panic("repository already set")
	}

	repo = repository.NewMongoDB(
		func(data map[string]interface{}) (Team, error) {
			trackData, ok := data["tracker"].(map[string]interface{})
			if !ok {
				return nil, errors.New("missing team tracker")
			}

			track := &Tracker{}
			if err := track.Unmarshal(trackData); err != nil {
				return nil, errors.Join(errors.New("failed to unmarshal team tracker: "), err)
			}

			t := &PlayerTeam{
				tracker: track,
			}
			if err := t.Unmarshal(data); err != nil {
				return nil, errors.Join(errors.New("failed to unmarshal player team: "), err)
			}

			return t, nil
		},
		func(t Team) (map[string]interface{}, error) {
			return t.Marshal()
		},
		"teams",
	)
}

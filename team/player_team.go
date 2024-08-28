package team

import (
	"errors"
	"github.com/bitrule/hcteams/common"
	"github.com/bitrule/hcteams/team/handler"
	"github.com/bitrule/hcteams/team/member"
	"github.com/google/uuid"
	"sync"
	"sync/atomic"

	"github.com/bitrule/hcteams/repository"
	"github.com/bitrule/hcteams/team/tickable"
)

var (
	membersMu sync.RWMutex
	membersId = make(map[string]string)
)

// PlayerTeam represents a team of players that contains a DTR tick and a bit more
type PlayerTeam struct {
	tracker *Tracker
	handler handler.Handler

	ownership string

	membersMu sync.RWMutex
	members   map[string]string

	dtr *tickable.DTRTick
}

// Tracker returns the team's tracker
func (m *PlayerTeam) Tracker() *Tracker {
	return m.tracker
}

// Handle sets the team's handler
func (m *PlayerTeam) Handle(h handler.Handler) {
	m.handler = h
}

// Handler returns the team's handler
func (m *PlayerTeam) Handler() handler.Handler {
	return m.handler
}

func (m *PlayerTeam) Ownership() string {
	return m.ownership
}

func (m *PlayerTeam) Members() map[string]string {
	m.membersMu.RLock()
	defer m.membersMu.RUnlock()

	return m.members
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

	for xuid := range m.Members() {
		delete(membersId, xuid)
	}

	membersMu.Unlock()

	return errors.New("not implemented")
}

// Unmarshal loads the monitor's configuration from a map
func (m *PlayerTeam) Unmarshal(prop map[string]interface{}) error {
	dtrData, ok := prop["dtr"].(map[string]interface{})
	if !ok {
		return errors.New("missing DTR tracker")
	}

	dtr := &tickable.DTRTick{}
	if err := dtr.Unmarshal(dtrData); err != nil {
		return errors.Join(errors.New("failed to unmarshal DTR tracker: "), err)
	}

	m.dtr = dtr

	return nil
}

// Marshal saves the monitor's configuration to a map
func (m *PlayerTeam) Marshal() (map[string]interface{}, error) {
	if m.dtr == nil {
		return nil, errors.New("missing DTR tick")
	}

	prop := make(map[string]interface{})
	if dtrData, err := m.dtr.Marshal(); err != nil {
		return nil, errors.Join(errors.New("failed to marshal DTR tracker: "), err)
	} else {
		prop["dtr"] = dtrData
	}

	return prop, nil
}

func Empty(ownership, name, teamType string) Team {
	tracker := &Tracker{
		id:       uuid.New().String(),
		name:     name,
		teamType: teamType,

		balance: atomic.Int32{},
		points:  atomic.Int32{},
	}

	if teamType == PlayerTeamType {
		return &PlayerTeam{
			tracker:   tracker,
			ownership: ownership,
			members: map[string]string{
				ownership: member.Leader.Name(),
			},
		}
	} else if teamType == SystemTeamType {
		return &SystemTeam{
			tracker: tracker,
		}
	}

	return nil
}

// LookupByPlayer returns the team of the player with the given source XUID
func LookupByPlayer(sourceXuid string) *PlayerTeam {
	membersMu.RLock()
	defer membersMu.RUnlock()

	id, ok := membersId[sourceXuid]
	if !ok {
		return nil
	}

	t, ok := Lookup(id).(*PlayerTeam)
	if !ok {
		return nil
	}

	return t
}

func Hook() {
	if repo != nil {
		common.Log.Panic("repository for teams already exists")
	}

	repo = repository.NewMongoDB(
		func(data map[string]interface{}) (Team, error) {
			trackProp, ok := data["tracker"].(map[string]interface{})
			if !ok {
				return nil, errors.New("missing team tracker")
			}

			tracker := &Tracker{}
			if err := tracker.Unmarshal(trackProp); err != nil {
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

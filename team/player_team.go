package team

import (
	"errors"
	"github.com/bitrule/hcteams/startup"
	"github.com/bitrule/hcteams/team/member"
	"github.com/google/uuid"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/bitrule/hcteams/repository"
	"github.com/bitrule/hcteams/team/tickable"
)

// PlayerTeam represents a team of players that contains a DTR tick and a bit more
type PlayerTeam struct {
	tracker *Tracker

	ownership string

	membersMu sync.RWMutex
	members   map[string]string

	invitesMu sync.RWMutex
	invites   []string

	dtr *tickable.DTRTick
}

// Tracker returns the team's tracker
func (t *PlayerTeam) Tracker() *Tracker {
	return t.tracker
}

func (t *PlayerTeam) Ownership() string {
	return t.ownership
}

func (t *PlayerTeam) Members() map[string]string {
	t.membersMu.RLock()
	defer t.membersMu.RUnlock()

	return t.members
}

func (t *PlayerTeam) HasMember(xuid string) bool {
	t.membersMu.RLock()
	defer t.membersMu.RUnlock()

	if _, ok := t.members[xuid]; ok {
		return true
	}

	return false
}

func (t *PlayerTeam) Role(xuid string) Role {
	t.membersMu.RLock()
	defer t.membersMu.RUnlock()

	if role, ok := t.members[xuid]; ok {
		return Role(role)
	}

	return Member
}

// Broadcast sends a message to all the team members
func (t *PlayerTeam) Broadcast(message string) {
	for xuid := range t.Members() {
		if p, ok := startup.SRV.PlayerByXUID(xuid); ok {
			p.Message(message)
		}
	}
}

// DTR returns the DTR tick of the team
func (t *PlayerTeam) DTR() *tickable.DTRTick {
	return t.dtr
}

// AddInvite adds an invitation to the team
func (t *PlayerTeam) AddInvite(xuid string) {
	t.invitesMu.Lock()
	t.invites = append(t.invites, xuid)
	t.invitesMu.Unlock()
}

// RemoveInvite removes an invite from the team
func (t *PlayerTeam) RemoveInvite(xuid string) {
	t.invitesMu.Lock()
	defer t.invitesMu.Unlock()

	if i := slices.Index(t.invites, xuid); i != -1 {
		t.invites = append(t.invites[:i], t.invites[i+1:]...)
	}
}

// HasInvite checks if the team has an invitation for a player
func (t *PlayerTeam) HasInvite(xuid string) bool {
	t.invitesMu.RLock()
	defer t.invitesMu.RUnlock()

	if i := slices.Index(t.invites, xuid); i != -1 {
		return true
	}

	return false
}

// Unmarshal loads the monitor's configuration from a map
func (t *PlayerTeam) Unmarshal(prop map[string]interface{}) error {
	invites, ok := prop["invites"].([]string)
	if !ok {
		return errors.New("missing invites")
	}
	t.invites = invites

	dtrProp, ok := prop["dtr"].(map[string]interface{})
	if !ok {
		return errors.New("missing DTR tracker")
	}

	dtr := &tickable.DTRTick{}
	if err := dtr.Unmarshal(dtrProp); err != nil {
		return errors.Join(errors.New("failed to unmarshal DTR tracker: "), err)
	}

	t.dtr = dtr

	return nil
}

// Marshal saves the monitor's configuration to a map
func (t *PlayerTeam) Marshal() (map[string]interface{}, error) {
	if t.dtr == nil {
		return nil, errors.New("missing DTR tick")
	}

	prop := make(map[string]interface{})
	prop["tracker"] = t.tracker.Marshal()
	prop["ownership"] = t.ownership

	t.invitesMu.RLock()
	prop["invites"] = t.invites
	t.invitesMu.RUnlock()

	t.membersMu.RLock()
	prop["members"] = t.members
	t.membersMu.RUnlock()

	if dtrData, err := t.dtr.Marshal(); err != nil {
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

func Hook() {
	// TODO: Optimize this a many bit
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

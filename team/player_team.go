package team

import (
	"errors"
	"github.com/bitrule/disrupt"
	"github.com/google/uuid"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/bitrule/disrupt/team/tickable"
)

// PlayerTeam represents a team of players that contains a DTR tick and a bit more
type PlayerTeam struct {
	tracker *Tracker

	ownership string
	hq        HQ

	membersMu sync.RWMutex
	members   map[string]Role

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

func (t *PlayerTeam) HQ() HQ {
	return t.hq
}

func (t *PlayerTeam) SetHQ(hq HQ) {
	t.hq = hq
}

func (t *PlayerTeam) Members() map[string]Role {
	t.membersMu.RLock()
	defer t.membersMu.RUnlock()

	return t.members
}

// AddMember adds a member to the team
func (t *PlayerTeam) AddMember(xuid string, role Role) {
	t.membersMu.Lock()
	t.members[xuid] = role
	t.membersMu.Unlock()
}

// RemoveMember removes a member from the team
func (t *PlayerTeam) RemoveMember(xuid string) {
	t.membersMu.Lock()

	if _, ok := t.members[xuid]; ok {
		delete(t.members, xuid)
	}

	t.membersMu.Unlock()
}

// Member returns the role of a player in the team
func (t *PlayerTeam) Member(xuid string) Role {
	t.membersMu.RLock()
	defer t.membersMu.RUnlock()

	if role, ok := t.members[xuid]; ok {
		return role
	}

	return Undefined
}

// Broadcast sends a message to all the team members
func (t *PlayerTeam) Broadcast(message string) {
	for xuid := range t.Members() {
		if p, ok := disrupt.SRV.PlayerByXUID(xuid); ok {
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

// RemoveInvite removes an invitation from the team
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

	body := make(map[string]interface{})
	body["tracker"] = t.tracker.Marshal()

	if t.hq.loaded {
		body["hq"] = t.hq.Marshal()
	}

	body["ownership"] = t.ownership

	t.invitesMu.RLock()
	body["invites"] = t.invites
	t.invitesMu.RUnlock()

	t.membersMu.RLock()
	body["members"] = t.members
	t.membersMu.RUnlock()

	if dtrData, err := t.dtr.Marshal(); err != nil {
		return nil, errors.Join(errors.New("failed to marshal DTR tracker: "), err)
	} else {
		body["dtr"] = dtrData
	}

	return body, nil
}

func NewPlayerTeam(ownership, name string) *PlayerTeam {
	return &PlayerTeam{
		tracker: &Tracker{
			id:       uuid.New().String(),
			name:     name,
			teamType: PlayerTeamType,

			balance: atomic.Int32{},
			points:  atomic.Int32{},
		},
		ownership: ownership,
		members: map[string]Role{
			ownership: Leader,
		},
	}
}

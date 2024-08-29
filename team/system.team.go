package team

import (
	"errors"
	"github.com/bitrule/hcteams/team/tickable"
)

type SystemTeam struct {
	tracker *Tracker

	tick tickable.Tick
}

func (t *SystemTeam) Tracker() *Tracker {
	return t.tracker
}

func (t *SystemTeam) DoTick() {
	t.tick.DoTick()
}

func (t *SystemTeam) Marshal() (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (t *SystemTeam) Unmarshal(prop map[string]interface{}) error {
	return errors.New("not implemented")
}

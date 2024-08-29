package team

import "errors"

type SystemTeam struct {
	tracker *Tracker
}

func (t *SystemTeam) Tracker() *Tracker {
	return t.tracker
}

func (t *SystemTeam) Marshal() (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (t *SystemTeam) Unmarshal(prop map[string]interface{}) error {
	return errors.New("not implemented")
}

package user

import (
	"errors"
)

// New creates an empty user
func New(xuid, name string) *User {
	return &User{
		xuid: xuid,
		name: name,
		tracker: &Tracker{
			kills:   0,
			deaths:  0,
			assists: 0,
		},
	}
}

type User struct {
	xuid string
	name string

	tracker *Tracker
}

// XUID returns the user's XUID
func (u *User) XUID() string {
	return u.xuid
}

// Name returns the user's name
func (u *User) Name() string {
	return u.name
}

// Tracker returns the user's tracker
func (u *User) Tracker() *Tracker {
	return u.tracker
}

// Unmarshal unmarshals the user from a map
func (u *User) Unmarshal(prop map[string]interface{}) error {
	xuid, ok := prop["xuid"].(string)
	if !ok {
		return errors.New("missing user XUID")
	}
	u.xuid = xuid

	name, ok := prop["name"].(string)
	if !ok {
		return errors.New("missing user name")
	}
	u.name = name

	trackerData, ok := prop["tracker"].(map[string]interface{})
	if !ok {
		return errors.New("missing user tracker")
	}

	tracker := &Tracker{}
	if err := tracker.Unmarshal(trackerData); err != nil {
		return errors.Join(errors.New("failed to unmarshal user tracker: "), err)
	}

	u.tracker = tracker

	return nil
}

// Marshal returns the user as a map
func (u *User) Marshal() (map[string]interface{}, error) {
	trackMarshal, err := u.tracker.Marshal()
	if err != nil {
		return nil, errors.Join(errors.New("failed to marshal user tracker: "), err)
	}

	return map[string]interface{}{
		"_id":     u.xuid,
		"name":    u.name,
		"tracker": trackMarshal,
	}, nil
}

package user

import (
	"errors"
)

// New creates an empty user
func New(xuid, name string) *User {
	return &User{
		xuid:    xuid,
		name:    name,
		tracker: &Tracker{},
	}
}

type User struct {
	xuid string
	name string

	teamAt string

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

// TeamAt returns the team the user is at
func (u *User) TeamAt() string {
	return u.teamAt
}

// SetTeamAt sets the team the user is at
func (u *User) SetTeamAt(team string) {
	u.teamAt = team
}

// Tracker returns the user's tracker
func (u *User) Tracker() *Tracker {
	return u.tracker
}

// Unmarshal unmarshals the user from a map
func (u *User) Unmarshal(body map[string]interface{}) error {
	xuid, ok := body["xuid"].(string)
	if !ok {
		return errors.New("missing user XUID")
	}
	u.xuid = xuid

	name, ok := body["name"].(string)
	if !ok {
		return errors.New("missing user name")
	}
	u.name = name

	trackerData, ok := body["tracker"].(map[string]interface{})
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

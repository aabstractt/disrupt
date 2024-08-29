package user

import (
	"errors"
	"github.com/bitrule/hcteams/repository"
	"github.com/bitrule/hcteams/startup"
	"sync"
)

var (
	usersIdMu sync.RWMutex
	usersMu   sync.RWMutex

	usersId = make(map[string]string)
	users   = make(map[string]*User)

	repo repository.Repository[User]
)

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

// Hook hooks the repository
func Hook() {
	if repo != nil {
		startup.Log.Panic("repository for users already exists")
	}

	repo = repository.NewMongoDB(
		func(data map[string]interface{}) (User, error) {
			u := User{}
			if err := u.Unmarshal(data); err != nil {
				return u, err
			}

			return u, nil
		},
		func(u User) (map[string]interface{}, error) {
			return u.Marshal()
		},
		"users",
	)

	values, err := repo.FindAll()
	if err != nil {
		startup.Log.Panic("failed to find all users: ", err)
	}

	for _, u := range values {
		// Make u as a pointer and store it
		Store(&u)
	}

	startup.Log.Infof("Successfully loaded %d user(s)", len(values))
}

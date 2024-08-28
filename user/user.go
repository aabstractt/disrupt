package user

import (
    "errors"
    "github.com/bitrule/hcteams/common"
    "github.com/bitrule/hcteams/repository"
    "strings"
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
	trackerData, err := u.tracker.Marshal()
	if err != nil {
		return nil, errors.Join(errors.New("failed to marshal user tracker: "), err)
	}

	return map[string]interface{}{
		"_id":     u.xuid,
		"name":    u.name,
		"tracker": trackerData,
	}, nil
}

// LookupByName returns the user with the given name
func LookupByName(name string) *User {
	usersMu.RLock()
	defer usersMu.RUnlock()

	usersIdMu.RLock()
	defer usersIdMu.RUnlock()

	if id, ok := usersId[strings.ToLower(name)]; ok {
		return users[id]
	}

	return nil
}

// Lookup returns the user with the given XUID
func Lookup(xuid string) *User {
	usersMu.RLock()
	defer usersMu.RUnlock()

	return users[xuid]
}

// Store registers the user into the cache
func Store(u *User) {
	// Lock the users map to prevent deadlocks
	usersMu.Lock()
	users[u.XUID()] = u
	usersMu.Unlock()

	// Lock the usersId map to prevent deadlocks
	usersIdMu.Lock()
	usersId[strings.ToLower(u.Name())] = u.XUID()
	usersIdMu.Unlock()
}

// Delete removes the user from the cache
func Delete(xuid string) {
	usersMu.Lock()
	defer usersMu.Unlock()

	if u, ok := users[xuid]; ok {
		usersIdMu.Lock()
		delete(usersId, strings.ToLower(u.Name()))
		usersIdMu.Unlock()
	}
}

// Hook hooks the repository
func Hook() {
	if repo != nil {
		common.Log.Panic("repository for users already exists")
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
		common.Log.Panic("failed to find all users: ", err)
	}

	for _, u := range values {
		// Make u as a pointer and store it
		Store(&u)
	}

	common.Log.Infof("Successfully loaded %d user(s)", len(values))
}

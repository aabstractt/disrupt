package service

import (
	"errors"
	"github.com/bitrule/hcteams/repository"
	"github.com/bitrule/hcteams/startup"
	"github.com/bitrule/hcteams/user"
	"github.com/df-mc/dragonfly/server/player"
	"sync"
)

type UserService struct {
	repository repository.Repository[*user.User]

	usersMu sync.RWMutex
	users   map[string]*user.User

	xuidsMu sync.RWMutex
	xuids   map[string]string
}

// LookupByXUID looks up a user by their XUID.
func (s *UserService) LookupByXUID(xuid string) *user.User {
	s.usersMu.RLock()
	defer s.usersMu.RUnlock()

	if u, ok := s.users[xuid]; ok {
		return u
	}

	return nil
}

// LookupByName looks up a user by their name.
func (s *UserService) LookupByName(name string) *user.User {
	s.xuidsMu.RLock()
	defer s.xuidsMu.RUnlock()

	if xuid, ok := s.xuids[name]; ok {
		return s.LookupByXUID(xuid)
	}

	return nil
}

// Cache caches a user in the repository.
func (s *UserService) cache(u *user.User) {
	s.usersMu.Lock()
	s.users[u.XUID()] = u
	s.usersMu.Unlock()

	s.xuidsMu.Lock()
	s.xuids[u.Name()] = u.XUID()
	s.xuidsMu.Unlock()
}

// Unload unloads a user from the repository.
func (s *UserService) Unload(p *player.Player) {
	s.usersMu.Lock()
	delete(s.users, p.XUID())
	s.usersMu.Unlock()

	s.xuidsMu.Lock()
	delete(s.xuids, p.Name())
	s.xuidsMu.Unlock()
}

// Save saves a user to the repository.
func (s *UserService) Save(xuid string) error {
	u := s.LookupByXUID(xuid)
	if u == nil {
		return errors.New("user not found")
	}

	if s.repository == nil {
		return errors.New("missing repository")
	}

	r, err := s.repository.Insert(u)
	if err != nil {
		return errors.Join(errors.New("failed to save the user: "), err)
	}

	if r.MatchedCount == 0 && r.UpsertedCount == 0 {
		return errors.New("failed to save the user: no documents matched the filter")
	}

	return nil
}

// Create creates a user.
func (s *UserService) Create(p *player.Player) error {
	if s.repository == nil {
		return errors.New("missing repository")
	}

	u := user.Empty(p.XUID(), p.Name())

	r, err := s.repository.Insert(u)
	if err != nil {
		return errors.Join(errors.New("failed to create the user: "), err)
	}

	if r.MatchedCount == 0 && r.UpsertedCount == 0 {
		return errors.New("failed to create the user: no documents matched the filter")
	}

	s.cache(u)

	return nil
}

// Hook hooks the repository to the service.
func (s *UserService) Hook() error {
	if s.repository != nil {
		return errors.New("repository already hooked")
	}

	users, err := s.repository.FindAll()
	if err != nil {
		return errors.Join(errors.New("failed to hook the repository: "), err)
	}

	for _, u := range users {
		s.cache(u)
	}

	startup.Log.Infof("Successfully loaded %d user(s)", len(s.users))

	return nil
}

var userService = &UserService{
	users: make(map[string]*user.User),
	xuids: make(map[string]string),
}

func User() *UserService {
	return userService
}

package service

import (
	"context"
	"errors"
	"github.com/bitrule/disrupt"
	"github.com/bitrule/disrupt/config"
	"github.com/bitrule/disrupt/user"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type UserService struct {
	col *mongo.Collection

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

// Hook hooks the repository to the service.
func (s *UserService) Hook() error {
	if s.col != nil {
		return errors.New("repository already hooked")
	}

	if disrupt.Mongo == nil {
		return errors.New("missing mongo client")
	}

	s.col = disrupt.Mongo.Database(config.DBConfig().DBName).Collection("users")

	cur, err := s.col.Find(context.Background(), nil)
	if err != nil {
		return errors.Join(errors.New("failed to hook the repository: "), err)
	}

	for cur.Next(context.Background()) {
		var body map[string]interface{}
		if err := cur.Decode(&body); err != nil {
			return errors.Join(errors.New("failed to decode the user: "), err)
		}

		u := &user.User{}
		if err := u.Unmarshal(body); err != nil {
			return errors.Join(errors.New("failed to unmarshal the user: "), err)
		}

		s.cache(u)
	}

	if err := cur.Close(context.TODO()); err != nil {
		return errors.Join(errors.New("failed to close the cursor: "), err)
	}

	disrupt.Log.Infof("Successfully loaded %d user(s)", len(s.users))

	return nil
}

// Save saves a user to the repository.
func (s *UserService) Save(u *user.User) error {
	if s.col == nil {
		return errors.New("missing repository")
	}

	r, err := s.col.UpdateOne(context.Background(), bson.M{IDKey: u.XUID()}, bson.M{"$set": u.Marshal()})
	if err != nil {
		return errors.Join(errors.New("failed to save the user: "), err)
	}

	if r.MatchedCount == 0 && r.UpsertedCount == 0 {
		return errors.New("failed to save the user: no documents matched the filter")
	}

	return nil
}

// Create creates a user.
func (s *UserService) Create(xuid, name string) error {
	u := user.New(xuid, name)
	if err := s.Save(u); err != nil {
		return err
	}

	s.cache(u)

	return nil
}

func (s *UserService) First(targets []cmd.Target) *player.Player {
	// Why this have more than one target?
	if len(targets) > 1 {
		return nil
	}

	for _, lt := range targets {
		if p, ok := lt.(*player.Player); ok {
			return p
		}
	}

	return nil
}

var userService = &UserService{
	users: make(map[string]*user.User),
	xuids: make(map[string]string),
}

func User() *UserService {
	return userService
}

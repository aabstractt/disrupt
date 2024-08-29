package service

import (
	"context"
	"errors"
	"github.com/bitrule/hcteams/startup"
	"github.com/bitrule/hcteams/user"
	"github.com/df-mc/dragonfly/server/player"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// Save saves a user to the repository.
func (s *UserService) Save(xuid string) error {
	u := s.LookupByXUID(xuid)
	if u == nil {
		return errors.New("user not found")
	}

	if s.col == nil {
		return errors.New("missing repository")
	}

	r, err := s.col.UpdateOne(context.Background(), bson.M{IDKey: u.XUID()}, bson.M{"$set": u})
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
	if s.col == nil {
		return errors.New("missing mongo collection")
	}

	u := user.New(p.XUID(), p.Name())

	r, err := s.col.UpdateOne(
		context.TODO(),
		bson.M{IDKey: u.XUID()},
		bson.M{"$set": u.Marshal()},
		options.Update().SetUpsert(true),
	)
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
func (s *UserService) Hook(dbname string) error {
	if s.col != nil {
		return errors.New("repository already hooked")
	}

	if startup.Mongo == nil {
		return errors.New("missing mongo client")
	}

	s.col = startup.Mongo.Database(dbname).Collection("users")

	cur, err := s.col.Find(context.Background(), nil)
	if err != nil {
		return errors.Join(errors.New("failed to hook the repository: "), err)
	}

	for cur.Next(context.Background()) {
		var prop map[string]interface{}
		if err := cur.Decode(&prop); err != nil {
			return errors.Join(errors.New("failed to decode the user: "), err)
		}

		u, err := s.decode(prop)
		if err != nil {
			return errors.Join(errors.New("failed to decode the user: "), err)
		}

		s.cache(u)
	}

	if err := cur.Close(context.TODO()); err != nil {
		return errors.Join(errors.New("failed to close the cursor: "), err)
	}

	startup.Log.Infof("Successfully loaded %d user(s)", len(s.users))

	return nil
}

func (s *UserService) decode(prop map[string]interface{}) (*user.User, error) {
	u := &user.User{}
	if err := u.Unmarshal(prop); err != nil {
		return nil, errors.Join(errors.New("failed to decode the user: "), err)
	}

	return u, nil
}

var userService = &UserService{
	users: make(map[string]*user.User),
	xuids: make(map[string]string),
}

func User() *UserService {
	return userService
}

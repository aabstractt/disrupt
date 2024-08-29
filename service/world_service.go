package service

import (
	"github.com/df-mc/dragonfly/server/world"
	"sync"
)

type WorldService struct {
	worldsMu sync.RWMutex
	worlds   map[string]*world.World
}

// LookupByName looks up a world by its name.
func (s *WorldService) LookupByName(name string) *world.World {
	s.worldsMu.RLock()
	defer s.worldsMu.RUnlock()

	if w, ok := s.worlds[name]; ok {
		return w
	}

	return nil
}

func (s *WorldService) Load(name string) *world.World {
	// TODO: Create world struct and load it
	w := world.New(name)
	s.cache(w)
	return w
}

// Unload unloads a world from the repository.
func (s *WorldService) Unload(w *world.World) {
	s.worldsMu.Lock()
	delete(s.worlds, w.Name())
	s.worldsMu.Unlock()
}

// Cache caches a world in the repository.
func (s *WorldService) cache(w *world.World) {
	s.worldsMu.Lock()
	s.worlds[w.Name()] = w
	s.worldsMu.Unlock()
}

func (s *WorldService) Hook() error {

}

func World() *WorldService {
	return worldService
}

var worldService = &WorldService{
	worlds: make(map[string]*world.World),
}

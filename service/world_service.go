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

// Cache caches a world in the repository.
func (s *WorldService) cache(w *world.World) {
	s.worldsMu.Lock()
	s.worlds[w.Name()] = w
	s.worldsMu.Unlock()
}

// Unload unloads a world from the repository.
func (s *WorldService) Unload(w *world.World) {
	s.worldsMu.Lock()
	delete(s.worlds, w.Name())
	s.worldsMu.Unlock()
}

func World() *WorldService {
	return &worldService
}

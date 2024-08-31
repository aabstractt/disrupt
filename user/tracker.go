package user

import (
    "errors"
    "sync/atomic"
)

type Tracker struct {
    kills  atomic.Int64
    deaths atomic.Int64

    assists atomic.Int64
}

// Kills returns the number of kills the user has
func (t *Tracker) Kills() int64 {
    return t.kills.Load()
}

// IncKills increments the number of kills the user has
func (t *Tracker) IncKills() {
    t.kills.Add(1)
}

// Deaths returns the number of deaths the user has
func (t *Tracker) Deaths() int64 {
    return t.deaths.Load()
}

// IncDeaths increments the number of deaths the user has
func (t *Tracker) IncDeaths() {
    t.deaths.Add(1)
}

// Assists returns the number of assists the user has
func (t *Tracker) Assists() int64 {
    return t.assists.Load()
}

// IncAssists increments the number of assists the user has
func (t *Tracker) IncAssists() {
    t.assists.Add(1)
}

// Marshal returns the tracker as a map
func (t *Tracker) Marshal() (map[string]interface{}, error) {
    return map[string]interface{}{
        "kills":   t.kills.Load(),
        "deaths":  t.deaths.Load(),
        "assists": t.assists.Load(),
    }, nil
}

// Unmarshal unmarshals the tracker from a map
func (t *Tracker) Unmarshal(body map[string]interface{}) error {
    kills, ok := body["kills"].(int64)
    if !ok {
        return errors.New("missing user kills")
    }

    deaths, ok := body["deaths"].(int64)
    if !ok {
        return errors.New("missing user deaths")
    }

    assists, ok := body["assists"].(int64)
    if !ok {
        return errors.New("missing user assists")
    }

    t.assists.Store(assists)
    t.deaths.Store(deaths)
    t.kills.Store(kills)

    return nil
}

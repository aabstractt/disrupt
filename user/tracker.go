package user

import "errors"

type Tracker struct {
	kills  int
	deaths int

	assists int
}

// Kills returns the number of kills the user has
func (t *Tracker) Kills() int {
	return t.kills
}

// IncKills increments the number of kills the user has
func (t *Tracker) IncKills() {
	t.kills++
}

// Deaths returns the number of deaths the user has
func (t *Tracker) Deaths() int {
	return t.deaths
}

// IncDeaths increments the number of deaths the user has
func (t *Tracker) IncDeaths() {
	t.deaths++
}

// Assists returns the number of assists the user has
func (t *Tracker) Assists() int {
	return t.assists
}

// IncAssists increments the number of assists the user has
func (t *Tracker) IncAssists() {
	t.assists++
}

// Unmarshal unmarshals the tracker from a map
func (t *Tracker) Unmarshal(prop map[string]interface{}) error {
	kills, ok := prop["kills"].(int)
	if !ok {
		return errors.New("missing user kills")
	}
	t.kills = kills

	deaths, ok := prop["deaths"].(int)
	if !ok {
		return errors.New("missing user deaths")
	}
	t.deaths = deaths

	assists, ok := prop["assists"].(int)
	if !ok {
		return errors.New("missing user assists")
	}
	t.assists = assists

	return nil
}

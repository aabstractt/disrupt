package team

import (
	"github.com/bitrule/hcteams/team/member"
	"sync"
	"sync/atomic"
)

type TeamData struct {
	id   string // Team ID
	name string // Team name

	ownership string // The XUID of the team owner
	membersMu sync.RWMutex
	members   map[string]*member.TeamMember // Team members

	balance atomic.Int32
	points  atomic.Int32
}

// Id returns the team's ID
func (m *TeamData) Id() string {
	return m.id
}

// Name returns the team's name
func (m *TeamData) Name() string {
	return m.name
}

// Ownership returns the XUID of the team owner
func (m *TeamData) Ownership() string {
	return m.ownership
}

func (m *TeamData) Members() map[string]*member.TeamMember {
	m.membersMu.RLock()
	defer m.membersMu.RUnlock()

	return m.members
}

func ParseTeamData(data map[string]interface{}) TeamData {
	id, ok := data["_id"].(string)
	if id == "" || !ok {
		panic("No id")
	}
	
	name, ok := data["name"].(string)
	if name == "" || !ok {
		panic("Exception")
	}

	ownership, ok := data["ownership"].(string)
	members, ok := data["members"].(map[string]string)

	bal, ok := data["balance"].(int64)
	points, ok := data["points"].(int64)

	// TODO: Wrap team data
}

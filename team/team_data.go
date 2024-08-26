package team

import (
	"github.com/bitrule/hcteams/team/member"
	"sync/atomic"
)

type TeamData struct {
	id   string // Team ID
	name string // Team name

	ownership string                        // The XUID of the team owner
	members   map[string]*member.TeamMember // Team members

	balance atomic.Int32
	points  atomic.Int32
}

func (m *TeamData) Id() string {
	return m.id
}

func (m *TeamData) Name() string {
	return m.name
}

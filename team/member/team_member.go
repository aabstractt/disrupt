package member

type TeamMember struct {
	xuid string // Xbox Live User ID
	name string // Gamertag

	kills  int // Number of kills
	deaths int // Number of deaths
}

// XUID returns the Xbox Live User ID of the team member
func (m *TeamMember) XUID() string {
	return m.xuid
}

// Name returns the gamertag of the team member
func (m *TeamMember) Name() string {
	return m.name
}

// Kills returns the number of kills the team member has
func (m *TeamMember) Kills() int {
	return m.kills
}

// AddKill increments the number of kills the team member has
func (m *TeamMember) AddKill() {
	m.kills++
}

// Deaths returns the number of deaths the team member has
func (m *TeamMember) Deaths() int {
	return m.deaths
}

// AddDeath increments the number of deaths the team member has
func (m *TeamMember) AddDeath() {
	m.deaths++
}

func NewTeamMember(xuid, name string, kills, deaths int) *TeamMember {
	return &TeamMember{
		xuid:   xuid,
		name:   name,
		kills:  kills,
		deaths: deaths,
	}
}

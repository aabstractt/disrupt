package member

type TeamMember struct {
	xuid string // Xbox Live User ID
	name string // Gamertag
	role Role   // Role

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

func NewTeamMember(xuid, name string, role Role, kills, deaths int) *TeamMember {
	return &TeamMember{
		xuid: xuid,
		name: name,

		role: role,

		kills:  kills,
		deaths: deaths,
	}
}

var (
	Leader  = Role(0)
	Officer = Role(1)
	Member  = Role(2)
)

type Role int // Role is a type that represents the role of a team member.

// Name returns the name of the role
func (r Role) Name() string {
	switch r {
	case Leader:
		return "Leader"
	case Officer:
		return "Officer"
	case Member:
		return "Member"
	}

	return "Unknown"
}

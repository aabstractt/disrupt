package config

var teamConfig TeamsConfig

type TeamsConfig struct {
	Name struct { // This is the section for the name values
		MinLength int `yaml:"min-length"`
		MaxLength int `yaml:"max-length"`
	} `yaml:"name"`

	Display struct { // This is the section for the display values
		FriendlyColour string `yaml:"friendly-colour"` // Friendly colour means the colour if is member of the team
		InvitedColour  string `yaml:"invited-colour"`  // Invited colour means the colour if is invited to the team
		EnemyColour    string `yaml:"enemy-colour"`    // Enemy colour means the colour if is not member of the team
	} `yaml:"display"`
}

// TeamConfig returns the team configuration.
func TeamConfig() TeamsConfig {
	return teamConfig
}

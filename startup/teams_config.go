package startup

type TeamsConfig struct {
	MongoDB struct {
		URI    string
		DBName string
	}

	Teams struct {
		DisplayColourSameTeam string
		DisplayColourInvited  string
		DisplayColourOpponent string
	}
}

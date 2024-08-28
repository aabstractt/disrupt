package main

import (
	"github.com/bitrule/hcteams/team"
	tcmd "github.com/bitrule/hcteams/team/cmd"
	"github.com/bitrule/hcteams/user"
	"github.com/df-mc/dragonfly/server/cmd"
)

func main() {

	// Hook means to initialize the package
	user.Hook()
	team.Hook()

	cmd.Register(cmd.New(
		"team",
		"Team commands",
		[]string{"t", "f", "faction"},
		tcmd.TeamSystemCreateCmd{},
		tcmd.TeamCreateCmd{},
		tcmd.TeamDisbandCmd{},
	))
}

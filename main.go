package main

import (
	"github.com/bitrule/hcteams/team"
	"github.com/bitrule/hcteams/user"
)

func main() {

	// Hook means to initialize the package
	user.Hook()
	team.Hook()
}

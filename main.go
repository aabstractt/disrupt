package main

import (
	"github.com/bitrule/hcteams/service"
	tcmd "github.com/bitrule/hcteams/team/cmd"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	if err := service.World().Hook(); err != nil {
		log.WithError(err).Panic("failed to hook world service")
	}

	if err := service.User().Hook(); err != nil {
		log.WithError(err).Panic("failed to hook user service")
	}

	cmd.Register(cmd.New(
		"team",
		"Manage your team. Use '/team help' for more information.",
		[]string{"t", "f", "faction"},
		tcmd.TeamSystemCreateCmd{},
		tcmd.TeamCreateCmd{},
		tcmd.TeamInviteCmd{},
		tcmd.TeamDisbandCmd{},
	))
}

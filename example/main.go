package main

import (
    "github.com/aabstractt/aurial/handler"
    "github.com/bitrule/disrupt/service"
    tcmd "github.com/bitrule/disrupt/team/cmd"
    "github.com/df-mc/dragonfly/server"
    "github.com/df-mc/dragonfly/server/cmd"
    "github.com/df-mc/dragonfly/server/player"
    "github.com/sirupsen/logrus"
    "time"
)

func main() {
    now := time.Now()
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
        tcmd.TeamLeaveCmd{},
        tcmd.TeamAcceptCmd{},
    ))

    ticker := time.NewTicker(50 * time.Millisecond)
    go func() {
        for range ticker.C {
            service.Team().DoTick()
            service.User().DoTick()
        }
    }()

    srv := server.New()
    srv.Accept(func(p *player.Player) {
        handler.Hook(p)
    })

    log.Infof("Server shutdown after %s", time.Since(now))

    log.Info("Shutting down services...")

    shutdownAt := time.Now()
    if err := service.Team().Shutdown(); err != nil {
        log.WithError(err).Panic("failed to shutdown team service")
    } else {
        log.Infof("Service for 'teams' has been shutdown in %s", time.Since(shutdownAt))
    }
}

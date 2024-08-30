package HCTeams

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var Log logrus.Logger
var SRV server.Server
var Mongo *mongo.Client

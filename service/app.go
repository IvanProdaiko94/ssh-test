package service

import (
	"github.com/IvanProdaiko94/ssh-test/cfg"
	"github.com/IvanProdaiko94/ssh-test/game"
	"github.com/IvanProdaiko94/ssh-test/persistence/postgres"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var App *Service

func init() {
	config := cfg.ReadEnv()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)

	db, err := postgres.InitWithRetries(config.PostgresConfig, time.Second*5, 3)
	if err != nil {
		panic(err)
	}
	var policy game.Policy
	if config.PolicyFilePath != "" {
		var err error
		policy, err = game.NewDefaultPolicy(config.PolicyFilePath)
		if err != nil {
			panic(err)
		}
	}
	App = New(config, postgres.NewTicTacToe(db), policy)
}

package main

import (
	"tirelease/api"
	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/cron"
)

func main() {
	// Load config
	configs.LoadConfig("config.yaml")

	// Connect database
	database.Connect(configs.Config)

	// Github Client (If Needed: V3 & V4)
	git.Connect(configs.Config.Github.AccessToken)
	git.ConnectV4(configs.Config.Github.AccessToken)

	// Start Cron (If Needed)
	cron.IssueCron()
	cron.PullRequestCron()

	// Start website && REST-API
	router := api.Routers("website/build/")
	router.Run(":8080")
}

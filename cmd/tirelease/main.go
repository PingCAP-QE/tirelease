package main

import (
	"tirelease/api"
	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/cron"
	"tirelease/internal/repository"
	"tirelease/internal/task"
)

func main() {
	// Load config
	configs.LoadConfig("config.yaml")

	// Connect database
	database.Connect(configs.Config)
	repository.InitHrEmployeeDB()

	// Github Client (If Needed: V3 & V4)
	git.Connect(configs.Config.Github.AccessToken)
	git.ConnectV4(configs.Config.Github.AccessToken)

	// Start Cron (If Needed)
	cron.InitCron()

	// Start Task Execution
	go task.StartTaskExecution()

	// Start website && REST-API
	router := api.Routers("website/build/")
	router.Run(":8080")
}

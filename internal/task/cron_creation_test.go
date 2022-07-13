package task

import (
	"testing"
	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/internal/entity"
)

func TestCreateCronTask(t *testing.T) {

	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	task := entity.Task{
		Type: entity.TASK_TYPE_REFRESH_EMPLOYEE,
	}

	// Create the same cron task twice to ensure the cron task is created only once
	CreateCronTask(task)

	CreateCronTask(task)
}

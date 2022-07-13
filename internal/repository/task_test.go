package repository

import (
	"fmt"
	"os"
	"testing"
	"time"
	"tirelease/commons/database"
	"tirelease/internal/entity"
)

func TestSelectAndUpdateFirst(t *testing.T) {

	var config = generateConfig()
	database.Connect(config)

	taskType := entity.TASK_TYPE_REFRESH_EMPLOYEE
	taskStatus := entity.TASK_STATUS_CREATED
	selectOption := entity.TaskOption{
		Type:   &taskType,
		Status: &taskStatus,
	}
	hostname, _ := os.Hostname()
	nowTime := time.Now()
	executingStatus := entity.TASK_STATUS_EXECUTING
	updateOption := entity.TaskOption{
		Status:      &executingStatus,
		Executor:    &hostname,
		ExecuteTime: &nowTime,
	}

	task, _ := SelectAndUpdateFirstTask(selectOption, updateOption)
	fmt.Printf("task: %+v\n", task)
}

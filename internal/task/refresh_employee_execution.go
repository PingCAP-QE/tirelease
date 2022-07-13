package task

import (
	"fmt"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type RefreshEmployeeTask struct {
	TaskExecutionBase
}

func (task RefreshEmployeeTask) fetch() *entity.Task {
	targetType := entity.TASK_TYPE_REFRESH_EMPLOYEE
	targetStatus := entity.TASK_STATUS_CREATED
	selectOption := entity.TaskOption{
		Type:   &targetType,
		Status: &targetStatus,
	}

	executingStatus := entity.TASK_STATUS_EXECUTING
	updateOption := entity.TaskOption{
		Status: &executingStatus,
	}

	targetTask, err := repository.SelectAndUpdateFirstTask(selectOption, updateOption)
	if err != nil {
		fmt.Printf("SelectAndUpdateFirstTask error: %s", err.Error())
		return nil
	}

	return targetTask
}

func (refreshTask RefreshEmployeeTask) process(task *entity.Task) error {
	fmt.Printf("RefreshEmployeeTask process %v \n", task)

	return nil
}

func NewRefreshEmployeeTask() RefreshEmployeeTask {
	task := RefreshEmployeeTask{}
	task.ITaskExecution = interface{}(task).(ITaskExecution)

	return task
}

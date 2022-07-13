package task

import (
	"fmt"
	"tirelease/internal/entity"
)

type RefreshEmployeeTask struct {
	TaskExecutionBase
}

func (refreshTask RefreshEmployeeTask) getTaskType() entity.TaskType {
	return entity.TASK_TYPE_REFRESH_EMPLOYEE
}

func (refreshTask RefreshEmployeeTask) process(task *entity.Task) error {
	// TODO Replace the logic using real refreshing of employee data.
	fmt.Printf("RefreshEmployeeTask process %v \n", task)

	return nil
}

func NewRefreshEmployeeTask() RefreshEmployeeTask {
	task := &RefreshEmployeeTask{}
	task.ITaskExecution = interface{}(task).(ITaskExecution)

	return *task
}

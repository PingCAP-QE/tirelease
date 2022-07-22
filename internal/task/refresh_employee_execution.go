package task

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type RefreshEmployeeTask struct {
	TaskExecutionBase
}

func (refreshTask RefreshEmployeeTask) getTaskType() entity.TaskType {
	return entity.TASK_TYPE_REFRESH_EMPLOYEE
}

func (refreshTask RefreshEmployeeTask) process(task *entity.Task) error {
	hrEmployees, err := repository.SelectAllHrEmployee()
	if err != nil {
		return err
	}

	employees := make([]entity.Employee, 0)
	for _, hrEmployee := range hrEmployees {
		employees = append(employees, hrEmployee.Trans())
	}

	err = repository.BatchCreateOrUpdateEmployees(employees)

	return err
}

func NewRefreshEmployeeTask() RefreshEmployeeTask {
	task := &RefreshEmployeeTask{}
	task.ITaskExecution = interface{}(task).(ITaskExecution)

	return *task
}

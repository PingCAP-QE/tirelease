package task

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

// Execute all the tasks in the task queue
func StartTaskExecution() {
	for {
		randSeed := rand.Intn(600)
		time.Sleep(time.Duration(randSeed) * time.Second)

		refreshEmployeeTask := NewRefreshEmployeeTask()
		refreshEmployeeTask.Execute()

		refreshIssueTask := NewRefreshIssueTask()
		refreshIssueTask.Execute()

		refreshPrTask := NewRefreshPrTask()
		refreshPrTask.Execute()
	}
}

type ITaskExecution interface {
	Execute()
	getTaskType() entity.TaskType
	fetch() *entity.Task
	init(task *entity.Task) error
	process(task *entity.Task) error
	finish(task *entity.Task, message string)
}

// TaskExecution template
// Using template design pattern
// The real task execution is implemented in the sub class
//     which only need to implement the **process and getTaskType** method to do the real work.
type TaskExecutionBase struct {
	ITaskExecution
}

func (execution TaskExecutionBase) Execute() {
	task := execution.ITaskExecution.fetch()
	if task == nil {
		return
	}

	err := execution.ITaskExecution.init(task)
	if err != nil {
		execution.ITaskExecution.finish(task, err.Error())
		return
	}

	err = execution.ITaskExecution.process(task)
	if err != nil {
		execution.ITaskExecution.finish(task, err.Error())
		return
	}

	execution.ITaskExecution.finish(task, "")
}

func (execution TaskExecutionBase) fetch() *entity.Task {
	targetType := execution.ITaskExecution.getTaskType()
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

func (execution TaskExecutionBase) init(task *entity.Task) error {
	id := task.ID
	selectOption := entity.TaskOption{
		ID: &id,
	}

	executeTime := time.Now()
	executor, _ := os.Hostname()
	updateOption := entity.TaskOption{
		ExecuteTime: &executeTime,
		Executor:    &executor,
	}

	updatedTask, err := repository.SelectAndUpdateFirstTask(selectOption, updateOption)
	if err != nil {
		return err
	}

	*task = *updatedTask
	return nil
}

func (execution TaskExecutionBase) finish(task *entity.Task, message string) {
	finishTime := time.Now()
	task.FinishTime = &finishTime

	if message != "" {
		task.Status = entity.TASK_STATUS_FAILED
		repository.UpdateTask(*task)
		return
	}

	task.Status = entity.TASK_STATUS_FINISHED
	repository.UpdateTask(*task)
}

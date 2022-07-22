package task

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
	"tirelease/internal/service"
)

type RefreshPrTask struct {
	TaskExecutionBase
}

func (refreshTask RefreshPrTask) getTaskType() entity.TaskType {
	return entity.TASK_TYPE_REFRESH_PR
}

func (refreshTask RefreshPrTask) process(task *entity.Task) error {
	repoOption := &entity.RepoOption{}
	repos, err := repository.SelectRepo(repoOption)
	if err != nil {
		return err
	}
	params := &service.RefreshPullRequestParams{
		Repos:       repos,
		BeforeHours: -2,
		Batch:       20,
		Total:       500,
	}

	return service.CronRefreshPullRequestV4(params)
}

func NewRefreshPrTask() RefreshPrTask {
	task := &RefreshPrTask{}
	task.ITaskExecution = interface{}(task).(ITaskExecution)

	return *task
}

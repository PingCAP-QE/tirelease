package cron

import (
	"tirelease/internal/task"

	"tirelease/commons/cron"
)

func InitCron() {
	cronTasks := task.GetCronTasks()
	// ATTENTION: the **for loop pattern** can not be 'for _, cronTask...'
	// Because the cron.Create function will always fetch the last cron task to create
	for i := range cronTasks {
		cronTask := cronTasks[i]
		cron.Create(cronTask.CronSchedule,
			func() {
				task.CreateCronTask(cronTask.Task)
			},
		)
	}
}

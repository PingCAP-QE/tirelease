package task

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type CronTask struct {
	entity.Task
	CronSchedule string
}

func GetCronTasks() []CronTask {
	cronTasks := make([]CronTask, 0)

	cronTasks = append(cronTasks, CronTask{
		Task: entity.Task{
			Type: entity.TASK_TYPE_REFRESH_PR,
		},
		CronSchedule: "0 0 */3 * * ?",
	})

	cronTasks = append(cronTasks, CronTask{
		Task: entity.Task{
			Type: entity.TASK_TYPE_REFRESH_ISSUE,
		},
		CronSchedule: "0 0 */1 * * ?",
	})

	cronTasks = append(cronTasks, CronTask{
		Task: entity.Task{
			Type: entity.TASK_TYPE_REFRESH_EMPLOYEE,
		},
		CronSchedule: "0 0 */12 * * ?",
	})

	return cronTasks
}

func CreateCronTask(task entity.Task) error {
	// Use random sleep time to ensure the task creation won't be conflict.
	randSeed := rand.Intn(5000)
	time.Sleep(time.Duration(randSeed) * time.Millisecond)

	hostname, _ := os.Hostname()
	task.HookType = entity.TASK_HOOK_TYPE_CRON
	task.Creator = hostname
	task.UniqueMeta = fmt.Sprintf("%d-%d-%d %d %s %s", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), task.Type, task.HookType)

	return repository.CreateTaskIfNotExist(task)
}

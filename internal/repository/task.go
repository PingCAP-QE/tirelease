package repository

import (
	"fmt"
	"time"
	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func CreateTaskIfNotExist(task entity.Task) error {
	initedTask := initTask(task)

	if err := database.DBConn.DB.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "type"}, {Name: "hook_type"}, {Name: "unique_meta"}},
			DoNothing: true,
		},
	).Create(&initedTask).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create task: %+v failed", initedTask))
	}
	return nil
}

func initTask(task entity.Task) entity.Task {
	if task.CreateTime.IsZero() {
		task.CreateTime = time.Now()
	}
	if task.UpdateTime.IsZero() {
		task.UpdateTime = time.Now()
	}

	task.Status = entity.TASK_STATUS_CREATED
	task.IsDeleted = false

	return task
}

package repository

import (
	"fmt"
	"time"
	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Select and update first task in one transaction.
// Return the updated task and error.
func SelectAndUpdateFirstTask(selectOption, updateOption entity.TaskOption) (*entity.Task, error) {
	result := &entity.Task{}
	updateTime := time.Now()
	updateOption.UpdateTime = &updateTime

	err := database.DBConn.DB.Transaction(func(tx *gorm.DB) error {
		tx = selectOption.Where(tx)

		if err := tx.Model(&entity.Task{}).First(result).Error; err != nil {
			return err
		}

		if err := tx.Model(result).Clauses(clause.Returning{}).Updates(updateOption.UpdateMap()).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, err
}

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

func SelectFirstTask(option entity.TaskOption) (*entity.Task, error) {
	result := &entity.Task{}
	if err := database.DBConn.DB.Where(option).First(result).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find task: %+v failed", option))
	}
	return result, nil
}

func UpdateTask(task entity.Task) error {
	task.UpdateTime = time.Now()
	if err := database.DBConn.DB.Save(&task).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("update task: %+v failed", task))
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

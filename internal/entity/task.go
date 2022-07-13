package entity

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          int64      `json:"id,omitempty"`
	IsDeleted   bool       `json:"is_deleted,omitempty"`
	Creator     string     `json:"creator,omitempty"`
	Executor    string     `json:"executor,omitempty"`
	Type        Type       `json:"type,omitempty"`
	HookType    HookType   `json:"hook_type,omitempty"`
	Status      Status     `json:"status,omitempty"`
	CreateTime  time.Time  `json:"create_time,omitempty"`
	UpdateTime  time.Time  `json:"update_time,omitempty"`
	ExecuteTime *time.Time `json:"execute_time,omitempty"`
	FinishTime  *time.Time `json:"finish_time,omitempty"`
	Message     string     `json:"message,omitempty"`
	UniqueMeta  string     `json:"unique_meta,omitempty"`
}

func (Task) TableName() string {
	return "task"
}

type HookType string

const (
	TASK_HOOK_TYPE_CRON    HookType = "cron"
	TASK_HOOK_TYPE_TRIGGER HookType = "trigger"
)

type Type string

const (
	TASK_TYPE_REFRESH_PR       = Type("REFRESH_PR")
	TASK_TYPE_REFRESH_ISSUE    = Type("REFRESH_ISSUE")
	TASK_TYPE_REFRESH_EMPLOYEE = Type("REFRESH_EMPLOYEE")
)

type Status string

const (
	TASK_STATUS_CREATED   = Status("created")
	TASK_STATUS_EXECUTING = Status("executing")
	TASK_STATUS_FINISHED  = Status("finished")
	TASK_STATUS_FAILED    = Status("failed")
	// STATUS_WAITING
)

type TaskOption struct {
	ID          *int64
	IsDeleted   *bool
	Creator     *string
	Executor    *string
	Type        *Type
	HookType    *HookType
	Status      *Status
	CreateTime  *time.Time
	UpdateTime  *time.Time
	ExecuteTime *time.Time
	FinishTime  *time.Time
	Message     *string
	UniqueMeta  *string
}

func (option TaskOption) Where(tx *gorm.DB) *gorm.DB {

	if option.ID != nil {
		tx = tx.Where("task.id = ?", *option.ID)
	}
	if option.IsDeleted != nil {
		tx = tx.Where("task.is_deletes = ?", *option.IsDeleted)
	}
	if option.Creator != nil {
		tx = tx.Where("task.creator = ?", *option.Creator)
	}
	if option.Executor != nil {
		tx = tx.Where("task.executor = ?", *option.Executor)
	}
	if option.Type != nil {
		tx = tx.Where("task.type = ?", *option.Type)
	}
	if option.HookType != nil {
		tx = tx.Where("task.hook_type = ?", *option.HookType)
	}
	if option.Status != nil {
		tx = tx.Where("task.status = ?", *option.Status)
	}
	if option.CreateTime != nil {
		tx = tx.Where("task.is_deletes = ?", *option.IsDeleted)
	}

	return tx
}

func (option TaskOption) UpdateMap() map[string]interface{} {

	result := make(map[string]interface{})
	if option.ID != nil {
		result["id"] = *option.ID
	}
	if option.IsDeleted != nil {
		result["is_deleted"] = *option.IsDeleted
	}
	if option.Creator != nil {
		result["creator"] = *option.Creator
	}
	if option.Executor != nil {
		result["executor"] = *option.Executor
	}
	if option.Type != nil {
		result["type"] = *option.Type
	}
	if option.HookType != nil {
		result["hook_type"] = *option.HookType
	}
	if option.Status != nil {
		result["status"] = *option.Status
	}
	if option.CreateTime != nil {
		result["create_time"] = *option.CreateTime
	}
	if option.UpdateTime != nil {
		result["update_time"] = *option.UpdateTime
	}
	if option.ExecuteTime != nil {
		result["execute_time"] = *option.ExecuteTime
	}
	if option.FinishTime != nil {
		result["finish_time"] = *option.FinishTime
	}
	if option.Message != nil {
		result["message"] = *option.Message
	}
	if option.UniqueMeta != nil {
		result["unique_meta"] = *option.UniqueMeta
	}

	return result
}

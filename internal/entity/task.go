package entity

import "time"

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
	TASK_STATUS_ERROR     = Status("error")
	// STATUS_WAITING
)

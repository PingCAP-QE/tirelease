package entity

import (
	"time"
)

// Struct of ci_detail
type TestEntity struct {
	ID            int64            `json:"id"`
	CreateTime    time.Time        `json:"create_time"`
	UpdateTime    time.Time        `json:"update_time"`
	JobName       string           `json:"job_name"`
	JobURL        string           `json:"job_url"`
	Repo          string           `json:"repo"`
	Branch        string           `json:"branch"`
	PullRequestID int64            `json:"pull_request_id"`
	CommitID      string           `json:"commit_id"`
	SuiteName     string           `json:"suite_name"`
	CaseName      string           `json:"case_name"`
	CaseClass     string           `json:"case_class"`
	ExecutionTime string           `json:"execution_time"`
	Status        TestEntityStatus `json:"status"`
	ErrorDetail   string           `json:"error_detail"`
	StackTrace    string           `json:"stack_trace"`
}

// Enum type
type TestEntityStatus string

// Enum list...
const (
	TestEntityStatusPassed = TestEntityStatus("passed")
	TestEntityStatusFailed = TestEntityStatus("failed")
	TestEntityStatusSkiped = TestEntityStatus("skipped")
)

// Select Option
type TestEntityOption struct {
	ID int64 `json:"id" form:"id"`
}

// DB-Table
func (TestEntity) TableName() string {
	return "test_entity"
}

/**
mysql --host 172.16.4.36 --port 3306 -u cicd_online -pwGEXq8a4MeCw6G

CREATE TABLE IF NOT EXISTS test_entity (
	id INT(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
	create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
	job_name VARCHAR(255) COMMENT '任务名称',
	job_url VARCHAR(1024) COMMENT '任务链接',
	repo VARCHAR(64) NOT NULL COMMENT '代码仓库',
	branch VARCHAR(64) COMMENT '代码分支',
	pull_request_id INT(11) COMMENT '合入请求ID',
	commit_id VARCHAR(255) COMMENT '代码提交ID',
	suite_name VARCHAR(255) COMMENT '组件名称',
	case_name VARCHAR(64) NOT NULL COMMENT '用例名称',
	case_class VARCHAR(1024) COMMENT '用例类名',
	execution_time VARCHAR(255) COMMENT '执行时长',
	status VARCHAR(32) NOT NULL COMMENT '用例结果',
	error_detail TEXT COMMENT '错误信息',
	stack_trace TEXT COMMENT '错误堆栈',

	PRIMARY KEY (id),
	// UNIQUE KEY uk_xxx (xxx)
	INDEX idx_createtime (create_time),
	INDEX idx_jobname_repo_branch (job_name, repo, branch)
)
ENGINE = INNODB DEFAULT CHARSET = utf8 COMMENT 'test_entity测试表';
**/

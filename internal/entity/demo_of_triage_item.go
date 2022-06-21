package entity

import (
	"time"
)

// Struct of triage_item
type TriageItem struct {
	ID            int64            `json:"id"`
	CreateTime    time.Time        `json:"create_time"`
	UpdateTime    time.Time        `json:"update_time"`
	ProjectName   string           `json:"project_name"`
	Repo          string           `json:"repo"`
	IssueID       int              `json:"issue_id"`
	PullRequestID int              `json:"pull_request_id"`
	Status        TriageItemStatus `json:"status"`
	Comment       string           `json:"comment"`
	IssueUrl      string           `json:"issue_url"`
}

// Enum type
type TriageItemStatus string

// Enum list...
const (
	TriageItemStatusInit   = TriageItemStatus("Init")
	TriageItemStatusPassed = TriageItemStatus("Accepted")
	TriageItemStatusFailed = TriageItemStatus("Won't Fix")
	TriageItemStatusSkiped = TriageItemStatus("Later")
)

// List Option
type TriageItemOption struct {
	ID int64 `json:"id" form:"id"`
}

// DB-Table
func (TriageItem) TableName() string {
	return "triage_item"
}

/**

mysql --host 172.16.4.36 --port 3306 -u cicd_online -pwGEXq8a4MeCw6G

CREATE TABLE IF NOT EXISTS triage_item (
	id INT(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
	create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
	project_name VARCHAR(255) COMMENT '项目名称',
	repo VARCHAR(64) NOT NULL COMMENT '代码仓库',
	issue_id INT(11) COMMENT '需求ID',
	pull_request_id INT(11) COMMENT '合入请求ID',
	status VARCHAR(32) NOT NULL COMMENT 'Triage结果',
	comment VARCHAR(1024) COMMENT '评论',
	issue_url VARCHAR(1024) COMMENT '需求链接',

	PRIMARY KEY (id),
	UNIQUE KEY uk_repo_issueid (repo, issue_id),
	INDEX idx_status (status)
)
ENGINE = INNODB DEFAULT CHARSET = utf8 COMMENT 'triage_item需求分类表';

**/

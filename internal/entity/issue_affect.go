package entity

import (
	"time"
)

type IssueAffect struct {
	ID         int64     `json:"id,omitempty"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`

	IssueID       string             `json:"issue_id,omitempty"`
	AffectVersion string             `json:"affect_version,omitempty"`
	AffectResult  AffectResultResult `json:"affect_result,omitempty"`
}

// Enum type
type AffectResultResult string

const (
	AffectResultResultUnKnown = AffectResultResult("UnKnown")
	AffectResultResultYes     = AffectResultResult("Yes")
	AffectResultResultNo      = AffectResultResult("No")
)

// List Option
type IssueAffectOption struct {
	ID            int64              `json:"id,omitempty" form:"id"`
	IssueID       string             `json:"issue_id,omitempty" form:"issue_id"`
	AffectVersion string             `json:"affect_version,omitempty" form:"affect_version"`
	AffectResult  AffectResultResult `json:"affect_result,omitempty" form:"affect_result"`

	IDs []int64 `json:"ids,omitempty" form:"ids"`

	ListOption
}

// Update Option
type IssueAffectUpdateOption struct {
	ID            int64              `json:"id,omitempty" form:"id"`
	IssueID       string             `json:"issue_id,omitempty" form:"issue_id"`
	AffectVersion string             `json:"affect_version,omitempty" form:"affect_version"`
	AffectResult  AffectResultResult `json:"affect_result,omitempty" form:"affect_result"`
}

// DB-Table
func (IssueAffect) TableName() string {
	return "issue_affect"
}

/**

CREATE TABLE IF NOT EXISTS issue_affect (
	id INT(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
	create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

	issue_id VARCHAR(255) COMMENT 'Issue全局ID',
	affect_version VARCHAR(255) NOT NULL COMMENT '版本号',
	affect_result VARCHAR(32) COMMENT 'Triage状态',

	PRIMARY KEY (id),
	UNIQUE KEY uk_issueid_affectversion (issue_id, affect_version)
)
ENGINE = INNODB DEFAULT CHARSET = utf8 COMMENT 'Issue影响版本信息表';

**/

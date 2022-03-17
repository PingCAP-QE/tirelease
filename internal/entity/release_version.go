package entity

import (
	"time"
)

type ReleaseVersion struct {
	// Columns
	ID int64 `json:"id,omitempty"`

	CreateTime        time.Time  `json:"create_time,omitempty"`
	UpdateTime        time.Time  `json:"update_time,omitempty"`
	PlanReleaseTime   *time.Time `json:"plan_release_time,omitempty"`
	ActualReleaseTime *time.Time `json:"actual_release_time,omitempty"`

	Name        string               `json:"name,omitempty"`
	Description string               `json:"description,omitempty"`
	Owner       string               `json:"owner,omitempty"`
	Type        ReleaseVersionType   `json:"type,omitempty"`
	Status      ReleaseVersionStatus `json:"status,omitempty"`

	ReleaseBranch            string `json:"release_branch,omitempty"`
	FatherReleaseVersionName string `json:"father_release_version_name,omitempty"`

	ReposString  string `json:"repos_string,omitempty"`
	LabelsString string `json:"labels_string,omitempty"`

	// OutPut-Serial
	Repos  *[]string `json:"repos,omitempty" gorm:"-"`
	Labels *[]string `json:"labels,omitempty" gorm:"-"`
}

// Enum status
type ReleaseVersionStatus string

const (
	ReleaseVersionStatusOpen     = ReleaseVersionStatus("Open")
	ReleaseVersionStatusClosed   = ReleaseVersionStatus("Closed")
	ReleaseVersionStatusFrozen   = ReleaseVersionStatus("Frozen")
	ReleaseVersionStatusReleased = ReleaseVersionStatus("Released")
)

// Enum type
type ReleaseVersionType string

const (
	ReleaseVersionTypeMajor = ReleaseVersionType("Major")
	ReleaseVersionTypeMinor = ReleaseVersionType("Minor")
	ReleaseVersionTypePatch = ReleaseVersionType("Patch")
)

// List Option
type ReleaseVersionOption struct {
	ID                       int64                `json:"id" form:"id"`
	Name                     string               `json:"name,omitempty" form:"name"`
	FatherReleaseVersionName string               `json:"father_release_version_name,omitempty" form:"father_release_version_name"`
	Type                     ReleaseVersionType   `json:"type,omitempty" form:"type"`
	Status                   ReleaseVersionStatus `json:"status,omitempty" form:"status"`
}

// DB-Table
func (ReleaseVersion) TableName() string {
	return "release_version"
}

/**

CREATE TABLE IF NOT EXISTS release_version (
	id INT(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
	create_time TIMESTAMP COMMENT '创建时间',
	update_time TIMESTAMP COMMENT '更新时间',
	plan_release_time TIMESTAMP COMMENT '计划发布时间',
	actual_release_time TIMESTAMP COMMENT '实际发布时间',

	name VARCHAR(255) NOT NULL COMMENT '版本号',
	description VARCHAR(1024) COMMENT '版本说明',
	owner VARCHAR(255) COMMENT '版本负责人',
	type VARCHAR(32) COMMENT '版本类型',
	status VARCHAR(32) COMMENT '版本状态',

	release_branch VARCHAR(255) COMMENT '发布分支',
	father_release_version_name VARCHAR(255) COMMENT '父版本号',
	repos_string VARCHAR(1024) COMMENT '代码仓库列表',
	labels_string VARCHAR(1024) COMMENT '标签规则列表',

	PRIMARY KEY (id),
	UNIQUE KEY uk_name (name),
	INDEX idx_createtime (create_time),
	INDEX idx_fathername (father_release_version_name)
)
ENGINE = INNODB DEFAULT CHARSET = utf8 COMMENT '发布版本表';

**/

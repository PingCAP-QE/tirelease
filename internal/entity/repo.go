package entity

import (
	"time"
)

// Struct of Issue
type Repo struct {
	// DataBase Column
	ID         int64     `json:"id,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty"`
	UpdateTime time.Time `json:"update_time,omitempty"`

	Owner       string  `json:"owner,omitempty"`
	Repo        string  `json:"repo,omitempty"`
	FullName    string  `json:"full_name,omitempty"`
	HTMLURL     *string `json:"html_url,omitempty"`
	Description *string `json:"description,omitempty"`
}

// List Option
type RepoOption struct {
	ID       int64  `json:"id" form:"id"`
	Owner    string `json:"owner,omitempty" form:"owner"`
	Repo     string `json:"repo,omitempty" form:"repo"`
	FullName string `json:"full_name,omitempty" form:"full_name"`
}

// DB-Table
func (Repo) TableName() string {
	return "repo"
}

/**

CREATE TABLE IF NOT EXISTS repo (
	id INT(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
	create_time TIMESTAMP COMMENT '创建时间',
	update_time TIMESTAMP COMMENT '更新时间',

	owner VARCHAR(255) COMMENT '仓库所有者',
	repo VARCHAR(255) COMMENT '仓库名称',
	full_name VARCHAR(255) COMMENT '仓库全称',
	html_url VARCHAR(1024) COMMENT '链接',
	description VARCHAR(1024) COMMENT '描述',

	PRIMARY KEY (id),
	UNIQUE KEY uk_owner_repo (owner, repo),
	INDEX idx_fullname (full_name)
)
ENGINE = INNODB DEFAULT CHARSET = utf8 COMMENT 'repo信息表';

**/

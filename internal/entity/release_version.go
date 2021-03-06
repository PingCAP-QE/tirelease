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

	Name     string `json:"name,omitempty"`
	Major    int    `json:"major,omitempty"`
	Minor    int    `json:"minor,omitempty"`
	Patch    int    `json:"patch,omitempty"`
	Addition string `json:"addition,omitempty"`

	Description string               `json:"description,omitempty"`
	Owner       string               `json:"owner,omitempty"`
	Type        ReleaseVersionType   `json:"type,omitempty"`
	Status      ReleaseVersionStatus `json:"status,omitempty"`

	ReleaseBranch string `json:"release_branch,omitempty"`
	ReposString   string `json:"repos_string,omitempty"`
	LabelsString  string `json:"labels_string,omitempty"`

	// OutPut-Serial
	Repos  *[]string `json:"repos,omitempty" gorm:"-"`
	Labels *[]string `json:"labels,omitempty" gorm:"-"`
}

// Enum status
type ReleaseVersionStatus string

const (
	ReleaseVersionStatusPlanned   = ReleaseVersionStatus("planned")
	ReleaseVersionStatusUpcoming  = ReleaseVersionStatus("upcoming")
	ReleaseVersionStatusFrozen    = ReleaseVersionStatus("frozen")
	ReleaseVersionStatusReleased  = ReleaseVersionStatus("released")
	ReleaseVersionStatusCancelled = ReleaseVersionStatus("cancelled")
)

// Enum type
type ReleaseVersionType string

const (
	ReleaseVersionTypeMajor  = ReleaseVersionType("Major")
	ReleaseVersionTypeMinor  = ReleaseVersionType("Minor")
	ReleaseVersionTypePatch  = ReleaseVersionType("Patch")
	ReleaseVersionTypeHotfix = ReleaseVersionType("Hotfix")
)

// Enum short type
type ReleaseVersionShortType string

const (
	ReleaseVersionShortTypeMajor   = ReleaseVersionShortType("%d")
	ReleaseVersionShortTypeMinor   = ReleaseVersionShortType("%d.%d")
	ReleaseVersionShortTypePatch   = ReleaseVersionShortType("%d.%d.%d")
	ReleaseVersionShortTypeHotfix  = ReleaseVersionShortType("%d.%d.%d-%s")
	ReleaseVersionShortTypeUnKnown = ReleaseVersionShortType("unknown")
)

// List Option
type ReleaseVersionOption struct {
	ID       int64                `json:"id" form:"id"`
	Name     string               `json:"name,omitempty" form:"name"`
	Major    int                  `json:"major,omitempty"`
	Minor    int                  `json:"minor,omitempty"`
	Patch    int                  `json:"patch,omitempty"`
	Addition string               `json:"addition,omitempty"`
	Type     ReleaseVersionType   `json:"type,omitempty" form:"type"`
	Status   ReleaseVersionStatus `json:"status,omitempty" form:"status"`

	StatusList []ReleaseVersionStatus `json:"status_list,omitempty" form:"status_list"`

	ShortType ReleaseVersionShortType `json:"short_type,omitempty" form:"short_type"`

	ListOption
}

// DB-Table
func (ReleaseVersion) TableName() string {
	return "release_version"
}

/**

CREATE TABLE IF NOT EXISTS release_version (
	id INT(11) NOT NULL AUTO_INCREMENT COMMENT '??????',
	create_time TIMESTAMP COMMENT '????????????',
	update_time TIMESTAMP COMMENT '????????????',
	plan_release_time TIMESTAMP COMMENT '??????????????????',
	actual_release_time TIMESTAMP COMMENT '??????????????????',

	name VARCHAR(255) NOT NULL COMMENT '?????????',
	major INT(11) NOT NULL COMMENT '????????????',
	minor INT(11) NOT NULL COMMENT '????????????',
	patch INT(11) NOT NULL COMMENT '???????????????',
	addition VARCHAR(255) COMMENT '?????????????????????',

	description VARCHAR(1024) COMMENT '????????????',
	owner VARCHAR(255) COMMENT '???????????????',
	type VARCHAR(32) COMMENT '????????????',
	status VARCHAR(32) COMMENT '????????????',

	release_branch VARCHAR(255) COMMENT '????????????',
	repos_string VARCHAR(1024) COMMENT '??????????????????',
	labels_string VARCHAR(1024) COMMENT '??????????????????',

	PRIMARY KEY (id),
	UNIQUE KEY uk_name (name),
	INDEX idx_major_minor_patch_addition (major, minor, patch, addition),
	INDEX idx_createtime (create_time)
)
ENGINE = INNODB DEFAULT CHARSET = utf8 COMMENT '???????????????';

**/

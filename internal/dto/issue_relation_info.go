package dto

import (
	"tirelease/internal/entity"
)

// IssueRelationInfo Query Struct
type IssueRelationInfoQuery struct {
	// Issue
	ID      int64  `json:"id,omitempty" form:"id"`
	IssueID string `json:"issue_id,omitempty" form:"issue_id"`
	Number  int    `json:"number,omitempty" form:"number"`
	State   string `json:"state,omitempty" form:"state"`
	Owner   string `json:"owner,omitempty" form:"owner"`
	Repo    string `json:"repo,omitempty" form:"repo"`

	SeverityLabel string `json:"severity_label,omitempty" form:"severity_label"`
	TypeLabel     string `json:"type_label,omitempty" form:"type_label"`

	// Filter Option
	AffectVersion string `json:"affect_version,omitempty" form:"affect_version"`
	BaseBranch    string `json:"base_branch,omitempty" form:"base_branch"`
}

// IssueRelationInfo ReturnBack Struct
type IssueRelationInfo struct {
	Issue            *entity.Issue
	IssueAffects     *[]entity.IssueAffect
	IssuePrRelations *[]entity.IssuePrRelation
	PullRequests     *[]entity.PullRequest
	VersionTriages   *[]entity.VersionTriage
}

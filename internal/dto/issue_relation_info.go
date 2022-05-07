package dto

import (
	"time"

	"tirelease/internal/entity"
)

// IssueRelationInfo Query Struct
type IssueRelationInfoQuery struct {
	// Issue
	entity.IssueOption

	CreatedAtStamp int64 `json:"created_at_stamp" form:"created_at_stamp"`
	UpdatedAtStamp int64 `json:"updated_at_stamp" form:"updated_at_stamp"`
	ClosedAtStamp  int64 `json:"closed_at_stamp" form:"closed_at_stamp"`

	// Filter Option
	AffectVersion string                      `json:"affect_version,omitempty" form:"affect_version" uri:"affect_version"`
	AffectResult  entity.AffectResultResult   `json:"affect_result,omitempty" form:"affect_result" uri:"affect_result"`
	BaseBranch    string                      `json:"base_branch,omitempty" form:"base_branch" uri:"base_branch"`
	VersionStatus entity.ReleaseVersionStatus `json:"version_status,omitempty" form:"version_status" uri:"version_status"`
}

// IssueRelationInfo ReturnBack Struct
type IssueRelationInfo struct {
	Issue            *entity.Issue             `json:"issue"`
	IssueAffects     *[]entity.IssueAffect     `json:"issue_affects"`
	IssuePrRelations *[]entity.IssuePrRelation `json:"issue_pr_relations"`
	PullRequests     *[]entity.PullRequest     `json:"pull_requests"`
	VersionTriages   *[]entity.VersionTriage   `json:"version_triages"`
}

// Join IssueRelationInfo
type IssueRelationInfoByJoin struct {
	// issue
	IssueID string `json:"issue_id,omitempty"`

	// issue_affect
	IssueAffectIDs string `json:"issue_affect_ids,omitempty"`
}

func (query *IssueRelationInfoQuery) ParamFill() {
	if query.CreatedAtStamp != 0 {
		query.CreateTime = time.Unix(query.CreatedAtStamp, 0)
	}
	if query.UpdatedAtStamp != 0 {
		query.UpdateTime = time.Unix(query.UpdatedAtStamp, 0)
	}
	if query.ClosedAtStamp != 0 {
		query.CloseTime = time.Unix(query.ClosedAtStamp, 0)
	}
}

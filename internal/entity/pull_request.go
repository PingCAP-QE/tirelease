package entity

import (
	"strings"
	"time"

	"tirelease/commons/git"

	"github.com/google/go-github/v41/github"
)

// Struct of Pull Request
type PullRequest struct {
	// DataBase columns
	ID            int64  `json:"id,omitempty"`
	PullRequestID string `json:"pull_request_id,omitempty"`
	Number        int    `json:"number,omitempty"`
	State         string `json:"state,omitempty"`
	Title         string `json:"title,omitempty"`
	Owner         string `json:"owner,omitempty"`
	Repo          string `json:"repo,omitempty"`
	HTMLURL       string `json:"html_url,omitempty"`
	BaseBranch    string `json:"base_branch,omitempty"`

	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
	MergedAt  *time.Time `json:"merged_at,omitempty"`

	Merged             bool    `json:"merged,omitempty"`
	MergeableState     *string `json:"mergeable_state,omitempty"`
	CherryPickApproved bool    `json:"cherry_pick_approved,omitempty"`
	AlreadyReviewed    bool    `json:"already_reviewed,omitempty"`

	SourcePullRequestID string `json:"source_pull_request_id,omitempty"`

	LabelsString             string `json:"labels_string,omitempty"`
	AssigneesString          string `json:"assignees_string,omitempty"`
	RequestedReviewersString string `json:"requested_reviewers_string,omitempty"`

	// OutPut-Serial
	Labels             *[]github.Label `json:"labels,omitempty" gorm:"-"`
	Assignees          *[]github.User  `json:"assignees,omitempty" gorm:"-"`
	RequestedReviewers *[]github.User  `json:"requested_reviewers,omitempty" gorm:"-"`
}

// List Option
type PullRequestOption struct {
	ID                  int64  `json:"id" form:"id"`
	PullRequestID       string `json:"pull_request_id,omitempty" form:"pull_request_id"`
	Number              int    `json:"number,omitempty" form:"number"`
	State               string `json:"state,omitempty" form:"state"`
	Owner               string `json:"owner,omitempty" form:"owner"`
	Repo                string `json:"repo,omitempty" form:"repo"`
	BaseBranch          string `json:"base_branch,omitempty" form:"base_branch"`
	SourcePullRequestID string `json:"source_pull_request_id,omitempty" form:"source_pull_request_id"`
}

// DB-Table
func (PullRequest) TableName() string {
	return "pull_request"
}

// ComposePullRequestFromV3
func ComposePullRequestFromV3(pullRequest *github.PullRequest) *PullRequest {
	alreadyReviwed := false
	cherryPickApproved := false
	labels := &[]github.Label{}
	for i := range pullRequest.Labels {
		node := pullRequest.Labels[i]
		label := github.Label{
			Name:  node.Name,
			Color: node.Color,
		}
		*labels = append(*labels, label)

		if *label.Name == git.CherryPickLabel {
			cherryPickApproved = true
		}
		if *label.Name == git.LGT2Label {
			alreadyReviwed = true
		}

	}
	assignees := &[]github.User{}
	for i := range pullRequest.Assignees {
		node := pullRequest.Assignees[i]
		user := github.User{
			Login: node.Login,
		}
		*assignees = append(*assignees, user)
	}
	requestedReviewers := &[]github.User{}
	for i := range pullRequest.RequestedReviewers {
		node := pullRequest.RequestedReviewers[i]
		user := github.User{
			Login: node.Login,
		}
		*requestedReviewers = append(*requestedReviewers, user)
	}
	mergeableState := strings.ToLower(*pullRequest.MergeableState)

	return &PullRequest{
		PullRequestID: *pullRequest.NodeID,
		Number:        *pullRequest.Number,
		State:         strings.ToLower(*pullRequest.State),
		Title:         *pullRequest.Title,
		Owner:         *pullRequest.Base.Repo.Owner.Login,
		Repo:          *pullRequest.Base.Repo.Name,
		HTMLURL:       *pullRequest.HTMLURL,
		BaseBranch:    *pullRequest.Base.Ref,

		CreatedAt: *pullRequest.CreatedAt,
		UpdatedAt: *pullRequest.UpdatedAt,
		ClosedAt:  pullRequest.ClosedAt,
		MergedAt:  pullRequest.MergedAt,

		Merged:             *pullRequest.Merged,
		MergeableState:     &mergeableState,
		CherryPickApproved: cherryPickApproved,
		AlreadyReviewed:    alreadyReviwed,

		Labels:             labels,
		Assignees:          assignees,
		RequestedReviewers: requestedReviewers,
	}
}

// ComposePullRequestFromV4
// TODO: v4 implement by tony at 2022/02/14
func ComposePullRequestFromV4(pullRequestField *git.PullRequestField) *PullRequest {
	alreadyReviwed := false
	cherryPickApproved := false
	labels := &[]github.Label{}
	for i := range pullRequestField.Labels.Nodes {
		node := pullRequestField.Labels.Nodes[i]
		label := github.Label{
			Name: github.String(string(node.Name)),
		}
		if node.Color != "" {
			label.Color = github.String(string(node.Color))
		}
		*labels = append(*labels, label)

		if *label.Name == git.CherryPickLabel {
			cherryPickApproved = true
		}
		if *label.Name == git.LGT2Label {
			alreadyReviwed = true
		}
	}
	assignees := &[]github.User{}
	for i := range pullRequestField.Assignees.Nodes {
		node := pullRequestField.Assignees.Nodes[i]
		user := github.User{
			Login: (*string)(&node.Login),
		}
		*assignees = append(*assignees, user)
	}
	requestedReviewers := &[]github.User{}
	for i := range pullRequestField.ReviewRequests.Nodes {
		node := pullRequestField.ReviewRequests.Nodes[i]
		user := github.User{
			Login: (*string)(&node.RequestedReviewer.Login),
		}
		*requestedReviewers = append(*requestedReviewers, user)
	}
	mergeableState := strings.ToLower(string(pullRequestField.Mergeable))

	pr := &PullRequest{
		PullRequestID: pullRequestField.ID.(string),
		Number:        int(pullRequestField.Number),
		State:         strings.ToLower(string(pullRequestField.State)),
		Title:         string(pullRequestField.Title),
		Owner:         string(pullRequestField.Repository.Owner.Login),
		Repo:          string(pullRequestField.Repository.Name),
		HTMLURL:       string(pullRequestField.Url),
		BaseBranch:    string(pullRequestField.BaseRefName),

		CreatedAt: pullRequestField.CreatedAt.Time,
		UpdatedAt: pullRequestField.UpdatedAt.Time,

		Merged:             bool(pullRequestField.Merged),
		MergeableState:     &mergeableState,
		CherryPickApproved: cherryPickApproved,
		AlreadyReviewed:    alreadyReviwed,

		Labels:             labels,
		Assignees:          assignees,
		RequestedReviewers: requestedReviewers,
	}
	if pullRequestField.ClosedAt != nil {
		pr.ClosedAt = &pullRequestField.ClosedAt.Time
	}
	if pullRequestField.MergedAt != nil {
		pr.MergedAt = &pullRequestField.MergedAt.Time
	}
	return pr
}

func ComposePullRequestWithoutTimelineFromV4(withoutTimeline *git.PullRequestFieldWithoutTimelineItems) *PullRequest {
	pullRequestField := &git.PullRequestField{
		PullRequestFieldWithoutTimelineItems: *withoutTimeline,
	}
	return ComposePullRequestFromV4(pullRequestField)
}

/**

CREATE TABLE IF NOT EXISTS pull_request (
	id INT(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
	pull_request_id VARCHAR(255) COMMENT 'Pr全局ID',
	number INT(11) NOT NULL COMMENT '当前库ID',
	state VARCHAR(32) NOT NULL COMMENT '状态',
	title VARCHAR(1024) COMMENT '标题',

	owner VARCHAR(255) COMMENT '仓库所有者',
	repo VARCHAR(255) COMMENT '仓库名称',
	html_url VARCHAR(1024) COMMENT '链接',
	base_branch VARCHAR(255) COMMENT '目标分支',

	closed_at TIMESTAMP COMMENT '关闭时间',
	created_at TIMESTAMP COMMENT '创建时间',
	updated_at TIMESTAMP COMMENT '更新时间',
	merged_at TIMESTAMP COMMENT '合入时间',

	merged BOOLEAN COMMENT '是否已合入',
	mergeable_state VARCHAR(32) COMMENT '可合入状态',
	cherry_pick_approved BOOLEAN COMMENT '是否已进入版本',
	already_reviewed BOOLEAN COMMENT '是否已代码评审',

	source_pull_request_id VARCHAR(255) COMMENT '来源ID',
	labels_string TEXT COMMENT '标签',
	assignees_string TEXT COMMENT '处理人列表',
	requested_reviewers_string TEXT COMMENT '处理人列表',

	PRIMARY KEY (id),
	UNIQUE KEY uk_prid (pull_request_id),
	INDEX idx_state (state),
	INDEX idx_owner_repo (owner, repo),
	INDEX idx_createdat (created_at),
	INDEX idx_sourceprid (source_pull_request_id)
)
ENGINE = INNODB DEFAULT CHARSET = utf8 COMMENT 'pull_request信息表';

**/

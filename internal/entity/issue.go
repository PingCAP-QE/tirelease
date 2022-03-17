package entity

import (
	"strings"
	"time"

	"tirelease/commons/git"

	"github.com/google/go-github/v41/github"
	"github.com/shurcooL/githubv4"
)

// Struct of Issue
type Issue struct {
	// DataBase Column
	ID      int64  `json:"id,omitempty"`
	IssueID string `json:"issue_id,omitempty"`
	Number  int    `json:"number,omitempty"`
	State   string `json:"state,omitempty"`
	Title   string `json:"title,omitempty"`
	Owner   string `json:"owner,omitempty"`
	Repo    string `json:"repo,omitempty"`
	HTMLURL string `json:"html_url,omitempty"`

	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`

	LabelsString    string `json:"labels_string,omitempty"`
	AssigneesString string `json:"assignees_string,omitempty"`

	ClosedByPullRequestID string `json:"closed_by_pull_request_id,omitempty"`
	SeverityLabel         string `json:"severity_label,omitempty"`
	TypeLabel             string `json:"type_label,omitempty"`

	// OutPut-Serial
	Labels    *[]github.Label `json:"labels,omitempty" gorm:"-"`
	Assignees *[]github.User  `json:"assignees,omitempty" gorm:"-"`
}

// List Option
type IssueOption struct {
	ID      int64  `json:"id" form:"id"`
	IssueID string `json:"issue_id,omitempty" form:"issue_id"`
	Number  int    `json:"number,omitempty" form:"number"`
	State   string `json:"state,omitempty" form:"state"`
	Owner   string `json:"owner,omitempty" form:"owner"`
	Repo    string `json:"repo,omitempty" form:"repo"`

	SeverityLabel string `json:"severity_label,omitempty" form:"severity_label"`
	TypeLabel     string `json:"type_label,omitempty" form:"type_label"`
}

// DB-Table
func (Issue) TableName() string {
	return "issue"
}

// ComposeIssueFromV3
func ComposeIssueFromV3(issue *github.Issue) *Issue {
	severityLabel := ""
	typeLabel := ""
	labels := &[]github.Label{}
	for i := range issue.Labels {
		node := issue.Labels[i]
		label := &github.Label{
			Name:  node.Name,
			Color: node.Color,
		}
		*labels = append(*labels, *label)

		if strings.HasPrefix(*label.Name, git.SeverityLabel) {
			severityLabel = *label.Name
		}
		if strings.HasPrefix(*label.Name, git.TypeLabel) {
			typeLabel = *label.Name
		}
	}
	assignees := &[]github.User{}
	for i := range issue.Assignees {
		node := issue.Assignees[i]
		user := &github.User{
			Login: node.Login,
		}
		*assignees = append(*assignees, *user)
	}
	url := strings.Split(*issue.RepositoryURL, "/")
	owner := url[len(url)-2]
	repo := url[len(url)-1]

	return &Issue{
		IssueID: *issue.NodeID,
		Number:  *issue.Number,
		State:   strings.ToLower(*issue.State),
		Title:   *issue.Title,
		Owner:   owner,
		Repo:    repo,
		HTMLURL: *issue.HTMLURL,

		CreatedAt: *issue.CreatedAt,
		UpdatedAt: *issue.UpdatedAt,
		ClosedAt:  issue.ClosedAt,

		Labels:    labels,
		Assignees: assignees,

		SeverityLabel: severityLabel,
		TypeLabel:     typeLabel,
	}
}

// ComposeIssueFromV4
// TODO: v4 implement by tony at 2022/02/14
func ComposeIssueFromV4(issueFiled *git.IssueField) *Issue {
	severityLabel := ""
	typeLabel := ""
	labels := &[]github.Label{}
	for i := range issueFiled.Labels.Nodes {
		node := issueFiled.Labels.Nodes[i]
		label := github.Label{
			Name: github.String(string(node.Name)),
		}
		if node.Color != "" {
			label.Color = github.String(string(node.Color))
		}
		*labels = append(*labels, label)

		if strings.HasPrefix(*label.Name, git.SeverityLabel) {
			severityLabel = *label.Name
		}
		if strings.HasPrefix(*label.Name, git.TypeLabel) {
			typeLabel = *label.Name
		}
	}
	assignees := &[]github.User{}
	for i := range issueFiled.Assignees.Nodes {
		node := issueFiled.Assignees.Nodes[i]
		user := github.User{
			Login: (*string)(&node.Login),
		}
		*assignees = append(*assignees, user)
	}
	closedByPrID := ""
	if issueFiled.State == githubv4.IssueStateClosed {
		for _, edge := range issueFiled.TimelineItems.Edges {
			closer := edge.Node.ClosedEvent.Closer.PullRequest
			if closer.Number != 0 {
				closedByPrID = closer.ID.(string)
			}
		}
	}

	issue := &Issue{
		IssueID: issueFiled.ID.(string),
		Number:  int(issueFiled.Number),
		State:   strings.ToLower(string(issueFiled.State)),
		Title:   string(issueFiled.Title),
		Owner:   string(issueFiled.Repository.Owner.Login),
		Repo:    string(issueFiled.Repository.Name),
		HTMLURL: string(issueFiled.Url),

		CreatedAt: issueFiled.CreatedAt.Time,
		UpdatedAt: issueFiled.UpdatedAt.Time,

		Labels:    labels,
		Assignees: assignees,

		ClosedByPullRequestID: closedByPrID,
		SeverityLabel:         severityLabel,
		TypeLabel:             typeLabel,
	}
	if issueFiled.ClosedAt != nil {
		issue.ClosedAt = &issueFiled.ClosedAt.Time
	}

	return issue
}

/**

CREATE TABLE IF NOT EXISTS issue (
	id INT(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
	issue_id VARCHAR(255) NOT NULL COMMENT 'Issue全局ID',
	number INT(11) NOT NULL COMMENT '当前库ID',
	state VARCHAR(32) NOT NULL COMMENT '状态',
	title VARCHAR(1024) COMMENT '标题',
	owner VARCHAR(255) COMMENT '仓库所有者',
	repo VARCHAR(255) COMMENT '仓库名称',
	html_url VARCHAR(1024) COMMENT '链接',

	closed_at TIMESTAMP COMMENT '关闭时间',
	created_at TIMESTAMP COMMENT '创建时间',
	updated_at TIMESTAMP COMMENT '更新时间',

	labels_string TEXT COMMENT '标签',
	assignees_string TEXT COMMENT '处理人列表',

	closed_by_pull_request_id VARCHAR(255) COMMENT '处理的PR',
	severity_label VARCHAR(255) COMMENT '严重等级',
	type_label VARCHAR(255) COMMENT '类型',

	PRIMARY KEY (id),
	UNIQUE KEY uk_issueid (issue_id),
	INDEX idx_state (state),
	INDEX idx_owner_repo (owner, repo),
	INDEX idx_createdat (created_at),
	INDEX idx_severitylabel (severity_label),
	INDEX idx_typelabel (type_label)
)
ENGINE = INNODB DEFAULT CHARSET = utf8 COMMENT 'issue信息表';

**/

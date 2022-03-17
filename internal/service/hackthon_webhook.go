package service

import (
	"strings"
	"time"
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/google/go-github/v41/github"
	"github.com/shurcooL/githubv4"
)

func UpdatePrAndIssue(webhookPayload WebhookPayload) error {
	if webhookPayload.Issue != nil {
		issue, err := git.ClientV4.GetIssueByNumber(webhookPayload.Repository.Owner.Login, webhookPayload.Repository.Name, webhookPayload.Issue.Number)
		// create issue affect
		for i := range MinorVersionList {
			minorVersion := MinorVersionList[i]
			issueAffect := entity.IssueAffect{
				CreateTime:    time.Now(),
				UpdateTime:    time.Now(),
				IssueID:       issue.ID.(string),
				AffectVersion: minorVersion,
				AffectResult:  entity.AffectResultResultUnKnown,
			}
			if err := repository.CreateOrUpdateIssueAffect(&issueAffect); err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
		if err := repository.CreateOrUpdateIssue(IssueFieldToIssue(issue)); err != nil {
			return err
		}
	}
	if webhookPayload.PullRequest != nil {
		pullRequest, err := git.ClientV4.GetPullRequestsByNumber(webhookPayload.Repository.Owner.Login, webhookPayload.Repository.Name, webhookPayload.PullRequest.Number)
		if err != nil {
			return err
		}
		if err := repository.CreateOrUpdatePullRequest(PullRequestFieldToPullRequest(pullRequest)); err != nil {
			return err
		}
	}
	return nil
}

func InitDB() error {
	issues, err := git.ClientV4.GetIssuesByTimeRangeV4("pingcap", "tidb", []string{"type/bug"}, time.Now().Add(-96*time.Hour), time.Now(), 20, 500)
	if err != nil {
		return err
	}
	for i := range issues {
		issue := issues[i]
		for j := range MinorVersionList {
			minorVersion := MinorVersionList[j]
			issueAffect := entity.IssueAffect{
				CreateTime:    time.Now(),
				UpdateTime:    time.Now(),
				IssueID:       issue.ID.(string),
				AffectVersion: minorVersion,
				AffectResult:  entity.AffectResultResultUnKnown,
			}
			if err := repository.CreateOrUpdateIssueAffect(&issueAffect); err != nil {
				return err
			}
		}
	}
	for _, issueFiled := range issues {
		issue := IssueFieldToIssue(&issueFiled)
		if err := repository.CreateOrUpdateIssue(issue); err != nil {
			return err
		}
	}
	prs, err := git.ClientV4.GetPullRequestsFromV4("pingcap", "tidb", time.Now().Add(-48*time.Hour), 20, 500)
	if err != nil {
		return err
	}
	for _, prNode := range prs {
		pr := PullRequestFieldToPullRequest(&prNode)
		if err := repository.CreateOrUpdatePullRequest(pr); err != nil {
			return err
		}
	}
	return nil
}

func IssueFieldToIssue(issueFiled *git.IssueField) *entity.Issue {
	labels := &[]github.Label{}
	for i := range issueFiled.Labels.Nodes {
		node := issueFiled.Labels.Nodes[i]
		label := github.Label{}
		label.Name = github.String(string(node.Name))
		*labels = append(*labels, label)
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

	resp := &entity.Issue{
		IssueID: issueFiled.ID.(string),
		Number:  int(issueFiled.Number),
		State:   string(issueFiled.State),
		Title:   string(issueFiled.Title),
		Repo:    strings.Join([]string{string(issueFiled.Repository.Owner.Login), string(issueFiled.Repository.Name)}, "/"),
		// ClosedAt:              issueFiled.ClosedAt.Time,
		CreatedAt:             issueFiled.CreatedAt.Time,
		UpdatedAt:             issueFiled.UpdatedAt.Time,
		Labels:                labels,
		Assignees:             assignees,
		ClosedByPullRequestID: closedByPrID,
	}
	// if !issueFiled.ClosedAt.Time.IsZero() {
	// 	resp.ClosedAt = issueFiled.ClosedAt.Time
	// }
	return resp
}

func PullRequestFieldToPullRequest(pullRequestField *git.PullRequestField) *entity.PullRequest {
	labels := &[]github.Label{}
	for i := range pullRequestField.Labels.Nodes {
		node := pullRequestField.Labels.Nodes[i]
		label := github.Label{}
		label.Name = github.String(string(node.Name))
		*labels = append(*labels, label)
	}
	assignees := &[]github.User{}
	for i := range pullRequestField.Assignees.Nodes {
		node := pullRequestField.Assignees.Nodes[i]
		user := github.User{
			Login: (*string)(&node.Login),
		}
		*assignees = append(*assignees, user)
	}
	SourcePrID := ""
	for _, edge := range pullRequestField.TimelineItems.Edges {
		sourcePr := edge.Node.CrossReferencedEvent.Source.PullRequest
		if sourcePr.Number != 0 && strings.HasPrefix(string(pullRequestField.Title), string(sourcePr.Title)) {
			SourcePrID = sourcePr.ID.(string)
		}
	}

	resp := &entity.PullRequest{
		PullRequestID: pullRequestField.ID.(string),
		Number:        int(pullRequestField.Number),
		State:         string(pullRequestField.State),
		Title:         string(pullRequestField.Title),
		Repo:          strings.Join([]string{string(pullRequestField.Repository.Owner.Login), string(pullRequestField.Repository.Name)}, "/"),
		BaseBranch:    string(pullRequestField.BaseRefName),
		// MergedAt:            pullRequestField.MergedAt.Time,
		CreatedAt:           pullRequestField.CreatedAt.Time,
		UpdatedAt:           pullRequestField.UpdatedAt.Time,
		Merged:              bool(pullRequestField.Merged),
		Labels:              labels,
		Assignees:           assignees,
		SourcePullRequestID: SourcePrID,
	}
	// if !pullRequestField.MergedAt.Time.IsZero() {
	// 	resp.MergedAt = pullRequestField.MergedAt.Time
	// }

	return resp
}

type WebhookPayload struct {
	PullRequest *PullRequest `json:"pull_request,omitempty"`
	Issue       *Issue       `json:"issue,omitempty"`
	Repository  *Repository  `json:"repository,omitempty"`
}

type PullRequest struct {
	ID     int64 `json:"id"`
	Number int   `json:"number"`
}

type Issue struct {
	ID     int64 `json:"id"`
	Number int   `json:"number"`
}

type Repository struct {
	Name  string `json:"name"`
	Owner Owner  `json:"owner"`
}

type Owner struct {
	Login string `json:"login"`
}

package service

import (
	"time"

	"tirelease/commons/git"
	"tirelease/internal/entity"

	"github.com/google/go-github/v41/github"
)

// Cron Job
type RefreshIssueParams struct {
	Repos       *[]entity.Repo `json:"repos"`
	BeforeHours int64          `json:"before_hours"`
	Batch       int            `json:"batch"`
	Total       int            `json:"total"`
}

func CronRefreshIssuesV4(params *RefreshIssueParams) error {
	// get repos
	if params.Repos == nil || len(*params.Repos) == 0 {
		return nil
	}

	// multi-batch refresh
	for _, repo := range *params.Repos {
		issues, err := git.ClientV4.GetIssuesByTimeRangeV4(
			repo.Owner, repo.Repo, nil,
			time.Now().Add(time.Duration(params.BeforeHours)*time.Hour), time.Now(),
			params.Batch, params.Total)
		if err != nil {
			return err
		}

		for i := range issues {
			issueRelation, err := ComposeIssueRelationInfoByIssueV4(&(issues[i]))
			if err != nil {
				return err
			}
			err = SaveIssueRelationInfo(issueRelation)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Git Webhook
// Webhook param only support v3 (v4 has no webhook right now)
func WebhookRefreshIssueV4(issue *github.Issue) error {
	// params
	if issue == nil {
		return nil
	}
	issueID := *issue.NodeID

	// handler
	issueRelationInfo, err := GetIssueRelationInfoByIssueIDV4(issueID)
	if err != nil {
		return err
	}
	if issueRelationInfo == nil {
		return nil
	}
	err = SaveIssueRelationInfo(issueRelationInfo)
	if err != nil {
		return err
	}

	return nil
}

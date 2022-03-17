package service

import (
	"time"

	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/google/go-github/v41/github"
)

// Cron Job
type RefreshPullRequestParams struct {
	Repos       *[]entity.Repo `json:"repos"`
	BeforeHours int64          `json:"before_hours"`
	Batch       int            `json:"batch"`
	Total       int            `json:"total"`
}

func CronRefreshPullRequestV4(params *RefreshPullRequestParams) error {
	// get repos
	if params.Repos == nil || len(*params.Repos) == 0 {
		return nil
	}

	// multi-batch refresh
	for _, repo := range *params.Repos {
		prs, err := git.ClientV4.GetPullRequestsFromV4(
			repo.Owner, repo.Repo,
			time.Now().Add(time.Duration(params.BeforeHours)*time.Hour),
			params.Batch, params.Total)
		if err != nil {
			return err
		}

		for i := range prs {
			err = repository.CreateOrUpdatePullRequest(entity.ComposePullRequestFromV4(&(prs[i])))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Git Webhook
// Webhook param only support v3 (v4 has no webhook right now)
func WebhookRefreshPullRequestV3(pr *github.PullRequest) error {
	// params
	if pr == nil {
		return nil
	}

	// handler
	err := repository.CreateOrUpdatePullRequest(entity.ComposePullRequestFromV3(pr))
	if err != nil {
		return err
	}

	return nil
}

package service

import (
	"log"
	"time"

	"tirelease/commons/git"
	"tirelease/internal/dto"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/google/go-github/v41/github"
)

// Cron Job
type RefreshIssueParams struct {
	Repos           *[]entity.Repo           `json:"repos"`
	BeforeHours     int64                    `json:"before_hours"`
	Batch           int                      `json:"batch"`
	Total           int                      `json:"total"`
	IsHistory       bool                     `json:"is_history"`
	ReleaseVersions *[]entity.ReleaseVersion `json:"release_versions"`
	Order           string                   `json:"order"`
}

func CronRefreshIssuesV4(params *RefreshIssueParams) error {
	// get repos
	if params.Repos == nil || len(*params.Repos) == 0 {
		return nil
	}

	// multi-batch refresh
	for _, repo := range *params.Repos {
		request := &git.RemoteIssueRangeRequest{
			Owner:      repo.Owner,
			Name:       repo.Repo,
			Labels:     nil,
			From:       time.Now().Add(time.Duration(params.BeforeHours) * time.Hour),
			To:         time.Now(),
			BatchLimit: params.Batch,
			TotalLimit: params.Total,
			Order:      params.Order,
		}
		issues, err := git.ClientV4.GetIssuesByTimeRangeV4(request)
		if err != nil {
			log.Fatal(err)
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

			if params.IsHistory {
				err := ExportHistoryVersionTriageInfo(issueRelation, params.ReleaseVersions)
				if err != nil {
					return err
				}
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
	err := WebhookRefreshIssueV4ByIssueID(issueID)
	if err != nil {
		return err
	}

	return nil
}

func WebhookRefreshIssueV4ByIssueID(issueID string) error {
	if issueID == "" {
		return nil
	}
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

// Service Function
func RefreshIssueLabel(label string, option *entity.IssueOption) error {
	// select issues
	issues, err := repository.SelectIssue(option)
	if err != nil {
		return err
	}
	if len(*issues) == 0 {
		return nil
	}

	// refresh label
	for _, issue := range *issues {
		_, _, err := git.Client.AddLabel(issue.Owner, issue.Repo, issue.Number, label)
		if err != nil {
			return err
		}
	}
	return nil
}

func RefreshIssueField(option *entity.IssueOption) error {
	// select issues
	issues, err := repository.SelectIssue(option)
	if err != nil {
		return err
	}
	if len(*issues) == 0 {
		return nil
	}

	// refresh info
	for _, issue := range *issues {
		issueID := issue.IssueID
		issueFieldWithoutTimelineItems, err := git.ClientV4.GetIssueWithoutTimelineByID(issueID)
		if err != nil {
			return err
		}
		if issueFieldWithoutTimelineItems == nil {
			continue
		}

		issueField := git.IssueField{
			IssueFieldWithoutTimelineItems: *issueFieldWithoutTimelineItems,
		}
		err = repository.CreateOrUpdateIssue(entity.ComposeIssueFromV4(&issueField))
		if err != nil {
			return err
		}
		log.Printf("id: %d OK", issue.ID)
	}
	return nil
}

func ExportHistoryVersionTriageWithDatabase(option *dto.IssueRelationInfoQuery) error {
	infos, _, err := FindIssueRelationInfo(option)
	if err != nil {
		return err
	}
	releaseVersions, err := repository.SelectReleaseVersion(&entity.ReleaseVersionOption{})
	if err != nil {
		return err
	}

	for i := range *infos {
		info := (*infos)[i]
		if info.PullRequests == nil || len(*info.PullRequests) == 0 {
			continue
		}

		err := ExportHistoryVersionTriageInfo(&info, releaseVersions)
		if err != nil {
			return err
		}
	}
	return nil
}

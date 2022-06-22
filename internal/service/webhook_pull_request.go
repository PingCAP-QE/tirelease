package service

import (
	"fmt"
	"strings"
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
		request := &git.RemoteIssueRangeRequest{
			Owner:      repo.Owner,
			Name:       repo.Repo,
			From:       time.Now().Add(time.Duration(params.BeforeHours) * time.Hour),
			BatchLimit: params.Batch,
			TotalLimit: params.Total,
		}
		prs, err := git.ClientV4.GetPullRequestsFromV4(request)
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

func CronMergeRetryPullRequestV3() error {
	// select no merge PRs
	merged := false
	cherryPickApproved := true
	alreadyReviewed := true
	option := &entity.PullRequestOption{
		State:              git.OpenStatus,
		Merged:             &merged,
		MergeableState:     &git.MergeableStateMergeable,
		CherryPickApproved: &cherryPickApproved,
		AlreadyReviewed:    &alreadyReviewed,
	}
	prs, err := repository.SelectPullRequest(option)
	if err != nil {
		return err
	}

	// retry
	for _, pr := range *prs {
		_, _, err := git.Client.CreateCommentByNumber(pr.Owner, pr.Repo, pr.Number, git.MergeRetryComment)
		if err != nil {
			return err
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

	// handler approve later

	return nil
}

func WebHookRefreshPullRequestRefIssue(pr *github.PullRequest) error {
	// params
	if pr == nil {
		return nil
	}
	pullRequestID := *(pr.NodeID)
	if pullRequestID == "" {
		return nil
	}

	// find close or ref issue numbers
	prV4, err := git.ClientV4.GetPullRequestByID(pullRequestID)
	if err != nil {
		return err
	}
	baseBranch := prV4.BaseRefName
	if baseBranch == "" || !strings.HasPrefix(string(baseBranch), git.ReleaseBranchPrefix) {
		return nil
	}

	doubleCheckRefreshIssuePrRef(prV4)

	return nil
}

func AutoRefreshPrApprovedLabel(pr *github.PullRequest) error {
	prEntity := entity.ComposePullRequestFromV3(pr)

	// Will not change the label temporaly for the stability
	// TODO If there is no other way to refresh the need-cherry-pick label, remove below condition
	if prEntity.CherryPickApproved {
		return nil
	}
	issueNumbers, err := git.ParseIssueNumber(prEntity.Body, prEntity.Owner, prEntity.Repo)

	if err != nil {
		return err
	}

	if len(issueNumbers) == 0 {
		return fmt.Errorf("pullrequest %s does not refer to a issue", prEntity.PullRequestID)
	}

	// Query issues refered by pullrequest
	issues := make([]entity.Issue, 0)
	for _, issueNumber := range issueNumbers {
		issueDOs, err := repository.SelectIssue(
			&entity.IssueOption{
				Number: issueNumber.Number,
				Owner:  issueNumber.Owner,
				Repo:   issueNumber.Repo,
			},
		)
		if err != nil {
			return err
		}
		issues = append(issues, *issueDOs...)
	}

	// Query releaseVersions from pullrequest
	// At most cases the size of releaseVersions is 1
	releaseVersions, err := getPrRelatedReleaseVersions(*prEntity)
	if err != nil {
		return err
	}

	hasFrozen, allApproved, err := checkTriageStatus(releaseVersions, issues)
	if err != nil {
		return err
	}

	// Skip below statuses to save brandwith
	// TODO If the labels need to be consistent with TiRelease, remove the condition
	if hasFrozen || !allApproved {
		return nil
	}

	err = ChangePrApprovedLabel(prEntity.PullRequestID, hasFrozen, allApproved)
	if err != nil {
		return err
	}
	return nil
}

// Check all triage status of versions and issues to see whether there is frozen or unapproved triage.
func checkTriageStatus(versions []entity.ReleaseVersion, issues []entity.Issue) (bool, bool, error) {
	if len(versions) == 0 || len(issues) == 0 {
		return false, false, nil
	}
	hasFrozen := false
	allApproved := true

	for _, version := range versions {
		for _, issue := range issues {
			triages, err := repository.SelectVersionTriage(
				&entity.VersionTriageOption{
					IssueID:     issue.IssueID,
					VersionName: version.Name,
				})

			if err != nil {
				return false, false, err
			}

			// If there is no triage history
			if len(*triages) == 0 {
				allApproved = false
				break
			} else if len(*triages) > 1 {
				return false, false, fmt.Errorf("value exception: There are multiple triage infos for version %s, issue %s", version.Name, issue.IssueID)
			}

			triage := (*triages)[0]
			if triage.TriageResult == entity.VersionTriageResultAcceptFrozen {
				hasFrozen = true
			}
			if triage.TriageResult != entity.VersionTriageResultAccept {
				allApproved = false
			}
		}
	}

	return hasFrozen, allApproved, nil
}

// Double check to ensure the relationship between pullrequest and issues is built.
func doubleCheckRefreshIssuePrRef(prV4 *git.PullRequestField) {
	refreshPrIssueRefByPrContent(prV4)
	refreshPrIssueRefByIssueNumber(prV4)
}

// Build the relation by parsing the exact issue number in the  body of pullrequest.
func refreshPrIssueRefByPrContent(prV4 *git.PullRequestField) error {
	repo := prV4.Repository.Name
	owner := prV4.Repository.Owner.Login
	prID := prV4.ID
	content := prV4.Body

	// Ensure that the pr is already restored.
	pr, err := repository.SelectPullRequestUnique(
		&entity.PullRequestOption{
			PullRequestID: prID.(string),
		},
	)
	if err != nil {
		return err
	}
	if pr == nil {
		return fmt.Errorf("pullrequest %s is not restored", prID)
	}

	issueNumbers, err := git.ParseIssueNumber(string(content), string(owner), string(repo))
	if err != nil {
		return err
	}

	for _, issueNumber := range issueNumbers {
		issue, err := repository.SelectIssueUnique(
			&entity.IssueOption{
				Number: issueNumber.Number,
				Owner:  issueNumber.Owner,
				Repo:   issueNumber.Repo,
			},
		)
		if err != nil || issue == nil {
			continue
		}

		repository.CreateIssuePrRelation(
			&entity.IssuePrRelation{
				PullRequestID: prID.(string),
				IssueID:       issue.IssueID,
			},
		)
	}
	return nil
}

// Build the relation by the blurred issue number in the body of pullrequest.
func refreshPrIssueRefByIssueNumber(prV4 *git.PullRequestField) error {

	issueNumbers, err := GetPullRequestRefIssuesByRegexFromV4(prV4)
	if err != nil {
		return err
	}

	// refresh cross-referenced issue
	if len(issueNumbers) > 0 {
		for _, issueNumber := range issueNumbers {
			issueOption := &entity.IssueOption{
				Number: issueNumber,
			}
			issues, err := repository.SelectIssue(issueOption)
			if err != nil {
				return err
			}
			if len(*issues) == 0 {
				continue
			}

			for _, issue := range *issues {
				err := WebhookRefreshIssueV4ByIssueID(issue.IssueID)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil

}

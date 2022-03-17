package service

import (
	"time"

	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/google/go-github/v41/github"
)

// Operation
func AddLabelByIssueID(issueID, label string) error {
	// select issue by id
	option := &entity.IssueOption{
		IssueID: issueID,
	}
	issue, err := repository.SelectIssueUnique(option)
	if nil != err {
		return err
	}

	// add issue label
	_, _, err = git.Client.AddLabel(issue.Owner, issue.Repo, issue.Number, label)
	if nil != err {
		return err
	}
	return nil
}

func RemoveLabelByIssueID(issueID, label string) error {
	// select issue by id
	option := &entity.IssueOption{
		IssueID: issueID,
	}
	issue, err := repository.SelectIssueUnique(option)
	if nil != err {
		return err
	}

	// remove issue label
	_, err = git.Client.RemoveLabel(issue.Owner, issue.Repo, issue.Number, label)
	if nil != err {
		return err
	}
	return nil
}

// Query Issue From Github And Construct Issue Data Service
func GetIssueByNumberFromV3(owner, repo string, number int) (*entity.Issue, error) {
	issue, _, err := git.Client.GetIssueByNumber(owner, repo, number)
	if nil != err {
		return nil, err
	}
	return entity.ComposeIssueFromV3(issue), nil
}

func GetIssuesByTimeFromV3(owner, repo string, time *time.Time) ([]*entity.Issue, error) {
	var page = 1
	var pageSize = 100
	option := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: pageSize,
		},
		Sort: "updated",
	}
	if nil != time {
		option.Since = *time
	}

	issues := make([]*entity.Issue, 0)
	for {
		gitIssues, _, err := git.Client.GetIssuesByOption(owner, repo, option)
		if nil != err {
			return nil, err
		}
		for i := range gitIssues {
			if nil == gitIssues[i].PullRequestLinks { // V3 considers every pull request an issue, so api return both issues and pull requests in the response
				issues = append(issues, entity.ComposeIssueFromV3(gitIssues[i]))
			}
		}
		page++
		option.ListOptions.Page = page

		if len(gitIssues) < pageSize {
			break
		}
	}

	return issues, nil
}

func GetIssuesByOptionV3(owner, repo string, option *github.IssueListByRepoOptions) ([]*entity.Issue, error) {
	var page = 1
	var pageSize = 20
	option.Page = page
	option.PerPage = pageSize
	issues := make([]*entity.Issue, 0)

	for {
		gitIssues, _, err := git.Client.GetIssuesByOption(owner, repo, option)
		if nil != err {
			return nil, err
		}
		for i := range gitIssues {
			if git.IsIssue(*gitIssues[i].HTMLURL) { // V3 considers every pull request an issue, so api return both issues and pull requests in the response
				issues = append(issues, entity.ComposeIssueFromV3(gitIssues[i]))
			}
		}
		page++
		option.ListOptions.Page = page

		if len(gitIssues) < pageSize {
			break
		}
	}
	return issues, nil
}

func BatchLabelIssues(issues []*entity.Issue, label string) error {
	if nil == issues {
		return nil
	}

	for _, issue := range issues {
		_, _, err := git.Client.AddLabel(issue.Owner, issue.Repo, issue.Number, label)
		if nil != err {
			return err
		}
	}
	return nil
}

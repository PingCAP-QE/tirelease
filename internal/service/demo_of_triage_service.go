package service

import (
	"context"

	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/google/go-github/v41/github"
)

// Collect triageItems
func CollectTriageItemByRepo(owner string, repo string) ([]*entity.TriageItem, error) {
	// Request
	issueListByRepoOptions := github.IssueListByRepoOptions{
		State: "all",
	}
	issueListByRepoOptions.Page = 0
	issueListByRepoOptions.PerPage = 100

	// Remote Search
	issues, _, err := git.Client.Client.Issues.ListByRepo(context.Background(), owner, repo, &issueListByRepoOptions)
	if nil != err {
		return nil, err
	}

	// Transform
	triageItems := transform(issues, owner, repo)
	return triageItems, nil
}

// Save triage_item list
func SavaTriageItems(triageItems []*entity.TriageItem) error {
	for _, triageItem := range triageItems {
		repository.TriageItemInsert(triageItem)
	}
	return nil
}

// Label operations
func AddLabelOfAccept(owner string, repo string, number int, labels []string) error {
	// Edit Labels
	_, _, err := git.Client.Client.Issues.AddLabelsToIssue(context.Background(), "VelocityLight", "tirelease", number, labels)
	if nil != err {
		return err
	}

	// UpdateDB
	return nil
}

func transform(issues []*github.Issue, owner string, repo string) []*entity.TriageItem {
	resp := []*entity.TriageItem{}
	for i := range issues {
		issue := issues[i]
		triageItem := &entity.TriageItem{
			CreateTime:  issue.GetCreatedAt(),
			UpdateTime:  issue.GetUpdatedAt(),
			ProjectName: "v4.0.16",
			Repo:        owner + "/" + repo,
			IssueID:     issue.GetNumber(),
			Status:      entity.TriageItemStatusInit,
			IssueUrl:    *issue.HTMLURL,
		}
		resp = append(resp, triageItem)
	}
	return resp
}

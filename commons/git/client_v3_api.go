package git

import (
	"context"

	"github.com/google/go-github/v41/github"
)

// ============================================================================ Repository
func (client *GithubInfo) GetRepositoriesByUser(user string) ([]*github.Repository, *github.Response, error) {
	return client.Client.Repositories.List(context.Background(), user, nil)
}

func (client *GithubInfo) GetRepositoriesByOrg(org string) ([]*github.Repository, *github.Response, error) {
	return client.Client.Repositories.ListByOrg(context.Background(), org, nil)
}

func (client *GithubInfo) GetRepositoryByOwnerAndName(owner, name string) (*github.Repository, *github.Response, error) {
	return client.Client.Repositories.Get(context.Background(), owner, name)
}

// ============================================================================ Issue
func (client *GithubInfo) GetIssueByNumber(owner, name string, number int) (*github.Issue, *github.Response, error) {
	return client.Client.Issues.Get(context.Background(), owner, name, number)
}

func (client *GithubInfo) GetIssuesByOption(owner, name string, option *github.IssueListByRepoOptions) ([]*github.Issue, *github.Response, error) {
	return client.Client.Issues.ListByRepo(context.Background(), owner, name, option)
}

func (client *GithubInfo) GetIssueCommentsByIssueNumber(owner, name string, number int) ([]*github.IssueComment, *github.Response, error) {
	return client.Client.Issues.ListComments(context.Background(), owner, name, number, &github.IssueListCommentsOptions{})
}

func (client *GithubInfo) GetIssueTimelinesByIssueNumber(owner, name string, number int) ([]*github.Timeline, *github.Response, error) {
	return client.Client.Issues.ListIssueTimeline(context.Background(), owner, name, number, nil)
}

func (client *GithubInfo) GetIssueEventsByIssueNumber(owner, name string, number int) ([]*github.IssueEvent, *github.Response, error) {
	return client.Client.Issues.ListIssueEvents(context.Background(), owner, name, number, nil)
}

// ============================================================================ PullRequest
func (client *GithubInfo) GetPullRequestByNumber(owner, name string, number int) (*github.PullRequest, *github.Response, error) {
	return client.Client.PullRequests.Get(context.Background(), owner, name, number)
}

// ============================================================================ Comment
func (client *GithubInfo) CreateCommentByNumber(owner, name string, number int, comment string) (*github.IssueComment, *github.Response, error) {
	issueComment := &github.IssueComment{
		Body: &comment,
	}
	return client.Client.Issues.CreateComment(context.Background(), owner, name, number, issueComment)
}

// ============================================================================ Label
func (client *GithubInfo) AddLabel(owner, name string, number int, label string) ([]*github.Label, *github.Response, error) {
	return client.Client.Issues.AddLabelsToIssue(context.Background(), owner, name, number, []string{label})
}

func (client *GithubInfo) RemoveLabel(owner, name string, number int, label string) (*github.Response, error) {
	res, err := client.Client.Issues.RemoveLabelForIssue(context.Background(), owner, name, number, label)
	if res.Response.StatusCode == 404 {
		// duplicate remove should return success
		return res, nil
	}
	return res, err
}

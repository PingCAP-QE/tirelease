package git

import (
	"context"

	"github.com/shurcooL/githubv4"
)

// ============================================================================ Issue
func (client *GithubInfoV4) GetIssueByID(id string) (*IssueField, error) {
	var query struct {
		Node struct {
			IssueField `graphql:"... on Issue"`
		} `graphql:"node(id: $id)"`
	}
	params := map[string]interface{}{
		"id": githubv4.ID(id),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}
	return &query.Node.IssueField, nil
}

func (client *GithubInfoV4) GetIssueWithoutTimelineByID(id string) (*IssueFieldWithoutTimelineItems, error) {
	var query struct {
		Node struct {
			IssueFieldWithoutTimelineItems `graphql:"... on Issue"`
		} `graphql:"node(id: $id)"`
	}
	params := map[string]interface{}{
		"id": githubv4.ID(id),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}
	return &query.Node.IssueFieldWithoutTimelineItems, nil
}

func (client *GithubInfoV4) GetIssueByNumber(owner, name string, number int) (*IssueField, error) {
	var query struct {
		Repository struct {
			Issue struct {
				IssueField
			} `graphql:"issue(number: $number)"`
		} `graphql:"repository(name: $name, owner: $owner)"`
	}
	num := githubv4.Int(number)
	params := map[string]interface{}{
		"number": num,
		"name":   githubv4.String(name),
		"owner":  githubv4.String(owner),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}
	return &query.Repository.Issue.IssueField, nil
}

// ============================================================================ PullRequest
func (client *GithubInfoV4) GetPullRequestByID(id string) (*PullRequestField, error) {
	var query struct {
		Node struct {
			PullRequestField `graphql:"... on PullRequest"`
		} `graphql:"node(id: $id)"`
	}
	params := map[string]interface{}{
		"id": githubv4.ID(id),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}
	return &query.Node.PullRequestField, nil
}

func (client *GithubInfoV4) GetPullRequestWithoutTimelineByID(id string) (*PullRequestFieldWithoutTimelineItems, error) {
	var query struct {
		Node struct {
			PullRequestFieldWithoutTimelineItems `graphql:"... on PullRequest"`
		} `graphql:"node(id: $id)"`
	}
	params := map[string]interface{}{
		"id": githubv4.ID(id),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}
	return &query.Node.PullRequestFieldWithoutTimelineItems, nil
}

func (client *GithubInfoV4) GetPullRequestsByNumber(owner, name string, number int) (*PullRequestField, error) {
	var query struct {
		Repository struct {
			PullRequest struct {
				PullRequestField
			} `graphql:"pullRequest(number: $number)"`
		} `graphql:"repository(name: $name, owner: $owner)"`
	}
	num := githubv4.Int(number)
	params := map[string]interface{}{
		"number": num,
		"name":   githubv4.String(name),
		"owner":  githubv4.String(owner),
	}
	if err := client.client.Query(context.Background(), &query, params); err != nil {
		return nil, err
	}
	return &query.Repository.PullRequest.PullRequestField, nil
}

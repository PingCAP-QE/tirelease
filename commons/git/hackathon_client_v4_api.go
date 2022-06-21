package git

import (
	"context"
	"log"
	"time"

	"github.com/shurcooL/githubv4"
)

// ============================================================================ Others
type RemoteIssueRangeRequest struct {
	Owner      string
	Name       string
	Labels     []string
	From       time.Time
	To         time.Time
	BatchLimit int
	TotalLimit int
	Order      string
}

func (client *GithubInfoV4) GetIssuesByTimeRangeV4(request *RemoteIssueRangeRequest) (issues []IssueField, err error) {
	// var query struct {
	// 	Repository struct {
	// 		Issues struct {
	// 			Edges []struct {
	// 				Cursor githubv4.String
	// 				Node   IssueField
	// 			}
	// 		} `graphql:"issues(first: $limit, after: $cursor, orderBy: {field: UPDATED_AT, direction: ASC}, labels: $labels, filterBy: {since: $since})"`
	// 	} `graphql:"repository(name: $name, owner: $owner)"`
	// }
	var query struct {
		Repository struct {
			Issues struct {
				Edges []struct {
					Cursor githubv4.String
					Node   IssueField
				}
			} `graphql:"issues(first: $limit, after: $cursor, orderBy: {field: UPDATED_AT, direction: $order}, filterBy: {since: $since})"`
		} `graphql:"repository(name: $name, owner: $owner)"`
	}

	cursor := (*githubv4.String)(nil)
	total := 0
	ghLabels := make([]githubv4.String, 0, len(request.Labels))
	if len(request.Labels) > 0 {
		for i := range request.Labels {
			ghLabels = append(ghLabels, githubv4.String(request.Labels[i]))
		}
	}

	since := request.From.Add(-1 * time.Minute).Format(time.RFC3339)
	log.Printf("fetching since %s", since)

	for request.TotalLimit != 0 {
		limit := request.BatchLimit
		if request.TotalLimit > 0 && request.TotalLimit < limit {
			limit = request.TotalLimit
		}
		param := map[string]interface{}{
			"name":   githubv4.String(request.Name),
			"owner":  githubv4.String(request.Owner),
			"limit":  githubv4.Int(limit),
			"cursor": cursor,
			"since":  githubv4.DateTime{Time: request.From.Add(-1 * time.Minute)},
			"order":  githubv4.OrderDirection(request.Order),
		}
		if len(ghLabels) > 0 {
			param["labels"] = ghLabels
		}

		err = client.client.Query(context.Background(), &query, param)
		if err != nil {
			log.Println(err)
			return
		}
		edges := query.Repository.Issues.Edges

		for i := range edges {
			issues = append(issues, edges[i].Node)
			log.Printf("%06d %s %s\n", edges[i].Node.Number, edges[i].Node.UpdatedAt.Format(time.RFC3339), edges[i].Node.Title)
		}

		cnt := len(edges)
		if cnt != 0 {
			lastIssue := &edges[cnt-1]
			cursor = &lastIssue.Cursor
			lastUpdated := lastIssue.Node.UpdatedAt.Time
			total += cnt
			request.TotalLimit -= cnt
			log.Println(cnt, "fetced", request.Owner, request.Name, request.Labels, lastUpdated)
			if lastUpdated.After(request.To) {
				break
			}
		}
		if cnt != limit {
			break
		}
	}

	log.Printf("fetched %d issues from %s/%s\n", total, request.Owner, request.Name)
	return
}

func (client *GithubInfoV4) GetPullRequestsFromV4(request *RemoteIssueRangeRequest) (prs []PullRequestField, err error) {
	var query struct {
		Repository struct {
			PullRequests struct {
				Edges []struct {
					Cursor githubv4.String
					Node   PullRequestField
				}
			} `graphql:"pullRequests(first: $limit, after: $cursor, orderBy: {field: UPDATED_AT, direction: DESC})"`
		} `graphql:"repository(name: $name, owner: $owner)"`
	}

	cursor := (*githubv4.String)(nil)
	total := 0

	since := request.From.Add(-1 * time.Minute)
	log.Printf("fetching since %s", since)

	for request.TotalLimit != 0 {
		limit := request.BatchLimit
		if request.TotalLimit > 0 && request.TotalLimit < limit {
			limit = request.TotalLimit
		}
		param := map[string]interface{}{
			"name":   githubv4.String(request.Name),
			"owner":  githubv4.String(request.Owner),
			"limit":  githubv4.Int(limit),
			"cursor": cursor,
		}

		err = client.client.Query(context.Background(), &query, param)
		if err != nil {
			log.Println(err)
			return
		}
		edges := query.Repository.PullRequests.Edges

		for i := range edges {
			prs = append(prs, edges[i].Node)
			log.Printf("%06d %s %s\n", edges[i].Node.Number, edges[i].Node.UpdatedAt.Format(time.RFC3339), edges[i].Node.Title)
		}

		cnt := len(edges)
		if cnt != 0 {
			lastIssue := &edges[cnt-1]
			cursor = &lastIssue.Cursor
			lastUpdated := lastIssue.Node.UpdatedAt.Time
			total += cnt
			request.TotalLimit -= cnt
			log.Println(cnt, "fetced", request.Owner, request.Name, lastUpdated)
			if since.After(lastUpdated) {
				break
			}
		}
		if cnt != limit {
			break
		}
	}

	log.Printf("fetched %d pull requests from %s/%s\n", total, request.Owner, request.Name)
	return
}

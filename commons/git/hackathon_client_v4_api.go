package git

import (
	"context"
	"log"
	"time"

	"github.com/shurcooL/githubv4"
)

// ============================================================================ Others
func (client *GithubInfoV4) GetIssuesByTimeRangeV4(owner, name string, labels []string, from time.Time, to time.Time, batchLimit int, totalLimit int) (issues []IssueField, err error) {
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
			} `graphql:"issues(first: $limit, after: $cursor, orderBy: {field: UPDATED_AT, direction: DESC}, filterBy: {since: $since})"`
		} `graphql:"repository(name: $name, owner: $owner)"`
	}

	cursor := (*githubv4.String)(nil)
	total := 0
	ghLabels := make([]githubv4.String, 0, len(labels))
	if len(labels) > 0 {
		for i := range labels {
			ghLabels = append(ghLabels, githubv4.String(labels[i]))
		}
	}

	since := from.Add(-1 * time.Minute).Format(time.RFC3339)
	log.Printf("fetching since %s", since)

	for totalLimit != 0 {
		limit := batchLimit
		if totalLimit > 0 && totalLimit < limit {
			limit = totalLimit
		}
		param := map[string]interface{}{
			"name":   githubv4.String(name),
			"owner":  githubv4.String(owner),
			"limit":  githubv4.Int(limit),
			"cursor": cursor,
			"since":  githubv4.DateTime{Time: from.Add(-1 * time.Minute)},
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
			totalLimit -= cnt
			log.Println(cnt, "fetced", owner, name, labels, lastUpdated)
			if lastUpdated.After(to) {
				break
			}
		}
		if cnt != limit {
			break
		}
	}

	log.Printf("fetched %d issues from %s/%s\n", total, owner, name)
	return
}

func (client *GithubInfoV4) GetPullRequestsFromV4(owner, name string, from time.Time, batchLimit int, totalLimit int) (prs []PullRequestField, err error) {
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

	since := from.Add(-1 * time.Minute)
	log.Printf("fetching since %s", since)

	for totalLimit != 0 {
		limit := batchLimit
		if totalLimit > 0 && totalLimit < limit {
			limit = totalLimit
		}
		param := map[string]interface{}{
			"name":   githubv4.String(name),
			"owner":  githubv4.String(owner),
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
			totalLimit -= cnt
			log.Println(cnt, "fetced", owner, name, lastUpdated)
			if since.After(lastUpdated) {
				break
			}
		}
		if cnt != limit {
			break
		}
	}

	log.Printf("fetched %d pull requests from %s/%s\n", total, owner, name)
	return
}

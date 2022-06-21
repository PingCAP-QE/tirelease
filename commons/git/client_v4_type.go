// V4 Object: https://docs.github.com/en/graphql/reference/objects

package git

import "github.com/shurcooL/githubv4"

// ============================================================= Struct Of Needed Fields
type IssueField struct {
	IssueFieldWithoutTimelineItems
	TimelineItems IssueTimelineItems `graphql:"timelineItems(first: 50, itemTypes: [CROSS_REFERENCED_EVENT, CLOSED_EVENT] )"`
}

type PullRequestField struct {
	PullRequestFieldWithoutTimelineItems
	TimelineItems PullRequestTimelineItems `graphql:"timelineItems(first: 50, itemTypes: [CROSS_REFERENCED_EVENT, ISSUE_COMMENT] )"`
}

type RepositoryField struct {
	Name  githubv4.String
	Owner struct {
		Login githubv4.String
	}
}

type UserField struct {
	Login githubv4.String
}

type LabelField struct {
	Nodes []struct {
		Name  githubv4.String
		Color githubv4.String
	}
}

type AssigneesFiled struct {
	Nodes []struct {
		UserField `graphql:"... on User"`
	}
}

type ReviewRequestField struct {
	Nodes []struct {
		RequestedReviewer struct {
			UserField `graphql:"... on User"`
		}
	}
}

type IssueFieldWithoutTimelineItems struct {
	ID     githubv4.ID
	Number githubv4.Int
	State  githubv4.IssueState
	Title  githubv4.String
	// Author     UserField
	Repository RepositoryField
	Url        githubv4.String
	// Body       githubv4.String

	CreatedAt githubv4.DateTime
	UpdatedAt githubv4.DateTime
	ClosedAt  *githubv4.DateTime

	Labels    LabelField     `graphql:"labels(first: 30)"`
	Assignees AssigneesFiled `graphql:"assignees(first: 10)"`
}

type PullRequestFieldWithoutTimelineItems struct {
	ID          githubv4.ID
	Number      githubv4.Int
	State       githubv4.PullRequestState
	Title       githubv4.String
	Body        githubv4.String
	Repository  RepositoryField
	Url         githubv4.String
	BaseRefName githubv4.String
	HeadRefName githubv4.String

	CreatedAt githubv4.DateTime
	UpdatedAt githubv4.DateTime
	ClosedAt  *githubv4.DateTime
	MergedAt  *githubv4.DateTime

	Merged    githubv4.Boolean
	Mergeable githubv4.MergeableState

	// MergeCommit struct {
	// 	OID           githubv4.GitObjectID
	// 	CommittedDate githubv4.DateTime
	// }
	// Author UserField

	Labels         LabelField         `graphql:"labels(first: 30)"`
	Assignees      AssigneesFiled     `graphql:"assignees(first: 10)"`
	ReviewRequests ReviewRequestField `graphql:"reviewRequests(first: 10)"`
}

type IssueTimelineItems struct {
	Edges []struct {
		Node struct {
			Typename             string `graphql:"__typename"`
			CrossReferencedEvent struct {
				WillCloseTarget githubv4.Boolean
				Source          struct {
					PullRequest PullRequestFieldWithoutTimelineItems `graphql:"... on PullRequest"`
				}
			} `graphql:"... on CrossReferencedEvent"`
			ClosedEvent struct {
				Actor  UserField
				Closer struct {
					PullRequest PullRequestFieldWithoutTimelineItems `graphql:"... on PullRequest"`
				}
			} `graphql:"... on ClosedEvent"`
		}
	}
}

type PullRequestTimelineItems struct {
	Edges []struct {
		Node struct {
			Typename             string `graphql:"__typename"`
			CrossReferencedEvent struct {
				Source struct {
					PullRequest PullRequestFieldWithoutTimelineItems `graphql:"... on PullRequest"`
				}
			} `graphql:"... on CrossReferencedEvent"`
			IssueComment struct {
				// Author UserField
				Body githubv4.String
			} `graphql:"... on IssueComment"`
		}
	}
}

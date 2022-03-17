package service

import (
	"time"

	"tirelease/commons/git"
	"tirelease/internal/dto"
	"tirelease/internal/entity"
)

// ================================================================ Compose Function From Remote Query
func GetIssueRelationInfoByIssueIDV4(issueID string) (*dto.IssueRelationInfo, error) {
	issue, err := git.ClientV4.GetIssueByID(issueID)
	if nil != err {
		return nil, err
	}

	return ComposeIssueRelationInfoByIssueV4(issue)
}

func GetIssueRelationInfoByIssueNumberV4(owner, repo string, number int) (*dto.IssueRelationInfo, error) {
	issue, err := git.ClientV4.GetIssueByNumber(owner, repo, number)
	if nil != err {
		return nil, err
	}

	return ComposeIssueRelationInfoByIssueV4(issue)
}

func ComposeIssueRelationInfoByIssueV4(issue *git.IssueField) (*dto.IssueRelationInfo, error) {
	issueAffects, err := ComposeIssueAffectWithIssueV4(issue)
	if nil != err {
		return nil, err
	}

	issuePrRelations, pullRequests, err := ComposeIssuePrRelationsByIssueV4(issue)
	if nil != err {
		return nil, err
	}

	triageRelationInfo := &dto.IssueRelationInfo{
		Issue:            entity.ComposeIssueFromV4(issue),
		IssueAffects:     issueAffects,
		IssuePrRelations: issuePrRelations,
		PullRequests:     pullRequests,
	}
	return triageRelationInfo, nil
}

func ComposeIssuePrRelationsByIssueV4(issue *git.IssueField) (*[]entity.IssuePrRelation, *[]entity.PullRequest, error) {
	edges := issue.TimelineItems.Edges
	if nil == edges || len(edges) == 0 {
		return nil, nil, nil
	}

	// Analyze timeline to compose IssuePrRelations & PullRequests
	issuePrRelations := make([]entity.IssuePrRelation, 0)
	pullRequests := make([]entity.PullRequest, 0)
	for i := range edges {
		edge := edges[i]
		if nil == &edge.Node || nil == &edge.Node.CrossReferencedEvent ||
			nil == &edge.Node.CrossReferencedEvent.Source || nil == &edge.Node.CrossReferencedEvent.Source.PullRequest {
			continue
		}
		if git.CrossReferencedEvent != edge.Node.Typename {
			continue
		}
		pr := edge.Node.CrossReferencedEvent.Source.PullRequest
		if nil == pr.ID {
			continue
		}

		var issuePrRelation = &entity.IssuePrRelation{
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
			IssueID:       issue.ID.(string),
			PullRequestID: pr.ID.(string),
		}
		issuePrRelations = append(issuePrRelations, *issuePrRelation)
		pullRequests = append(pullRequests, *(entity.ComposePullRequestWithoutTimelineFromV4(&pr)))
	}

	// Return
	return &issuePrRelations, &pullRequests, nil
}

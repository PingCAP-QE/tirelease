package service

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

func getTriageAndPRsMap(triages []entity.VersionTriage, version string) (map[entity.VersionTriage][]entity.PullRequest, error) {
	// Get all related issueIds
	issueIDs := extractIssueIDsFromTriage(triages)

	issuePrRelations, err := getIssuePrRelation(issueIDs)
	if err != nil {
		return nil, err
	}

	pullrequests, err := getRelatedPullRequests(issuePrRelations)
	if err != nil {
		return nil, err
	}

	releaseVersion, err := repository.SelectReleaseVersionLatest(
		&entity.ReleaseVersionOption{
			Name: version,
		},
	)
	if err != nil {
		return nil, err
	}

	pullrequests = filterPRsByBranch(pullrequests, ComposeVersionBranch(releaseVersion))

	return mapVersionTriagesWithPrs(triages, issuePrRelations, pullrequests), nil
}

func extractIssueIDsFromTriage(triages []entity.VersionTriage) []string {
	issueIDs := make([]string, 0)
	for _, triage := range triages {
		issueIDs = append(issueIDs, triage.IssueID)
	}
	return issueIDs
}

func filterPRsByBranch(prs []entity.PullRequest, branch string) []entity.PullRequest {
	var filteredPRs []entity.PullRequest
	for _, pr := range prs {
		if pr.BaseBranch == branch {
			filteredPRs = append(filteredPRs, pr)
		}
	}
	return filteredPRs
}

// Notes: the version of pullrequests in the params is the same with the triage
// TODO: refactor the model of version triage to contain the related info.
func mapVersionTriagesWithPrs(triages []entity.VersionTriage, issuePrRelation []entity.IssuePrRelation, prs []entity.PullRequest) map[entity.VersionTriage][]entity.PullRequest {
	triagePRMap := make(map[entity.VersionTriage][]entity.PullRequest)
	for _, triage := range triages {
		for _, relation := range issuePrRelation {
			if triage.IssueID != relation.IssueID {
				continue
			}

			for _, pr := range prs {
				if relation.PullRequestID == pr.PullRequestID {
					triagePRMap[triage] = append(triagePRMap[triage], pr)
				}
			}
		}
	}
	return triagePRMap
}

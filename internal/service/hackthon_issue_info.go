package service

import (
	"strings"
	"time"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
	// "github.com/google/go-github/v41/github"
)

var MinorVersionList = []string{"4.0", "5.0", "5.1", "5.2", "5.3", "5.4"}

func ListIssueInfo(state string) ([]*IssueInfo, error) {
	// resp := []*IssueInfo{}
	// status := entity.ReleaseVersionStatusUpcoming
	// typeV := entity.ReleaseVersionTypePatch
	// optionRV := entity.ReleaseVersionOption{
	// 	Status: status,
	// 	Type:   typeV,
	// }
	// patchVersions, err := repository.SelectReleaseVersion(&optionRV)
	// if err != nil {
	// 	return resp, err
	// }
	// minorPatchVersionMap := map[string]string{}
	// for i := range *patchVersions {
	// 	patchVersion := (*patchVersions)[i]
	// 	minorPatchVersionMap[patchVersion.FatherReleaseVersionName] = patchVersion.Name
	// }

	// optionI := entity.IssueOption{State: state}
	// issues, err := repository.SelectIssue(&optionI)
	// if err != nil {
	// 	return resp, err
	// }

	// for i := range *issues {
	// 	issue := (*issues)[i]
	// 	if issue.Labels == nil {
	// 		issue.Labels = &[]github.Label{}
	// 	}
	// 	if issue.Assignees == nil {
	// 		issue.Assignees = &[]github.User{}
	// 	}
	// 	if issue.ClosedAt == nil {
	// 		issue.ClosedAt = &time.Time{}
	// 	}
	// 	affects, err := ListAffected(issue.IssueID, issue.ClosedByPullRequestID, minorPatchVersionMap)
	// 	if err != nil {
	// 		return resp, err
	// 	}

	// 	issueInfo := &IssueInfo{
	// 		IssueID:   issue.IssueID,
	// 		Number:    issue.Number,
	// 		Title:     issue.Title,
	// 		Url:       issue.HTMLURL,
	// 		CreatedAt: issue.CreatedAt,
	// 		ClosedAt:  *issue.ClosedAt,
	// 		Severity:  getIssueSeverity(*issue.Labels),
	// 		State:     issue.State,
	// 		Type:      getIssueType(*issue.Labels),
	// 		Assignee:  getAssignee(*issue.Assignees),
	// 		Labels:    []string{},
	// 		Affects:   affects,
	// 	}
	// 	resp = append(resp, issueInfo)
	// }
	// return resp, nil
	return nil, nil
}

func FilterIssueInfo(minorVersion string) ([]*IssueInfo, error) {
	resp := []*IssueInfo{}
	issues, err := ListIssueInfo("CLOSED")
	if err != nil {
		return resp, err
	}
	for i := range issues {
		for _, affect := range issues[i].Affects {
			if affect.Version == minorVersion {
				if affect.Affect == "unknown" || affect.Affect == "yes" {
					if affect.Release.TriageStatus == "" || affect.Release.TriageStatus == "unknown" || affect.Release.TriageStatus == "accept" {
						resp = append(resp, issues[i])
					}
				}
			}
		}
	}
	return resp, nil
}

func ListAffected(issueID string, closedPrID string, minorPatchVersionMap map[string]string) ([]*Affect, error) {
	resp := []*Affect{}
	issueAffects, err := repository.SelectIssueAffect(&entity.IssueAffectOption{IssueID: issueID})
	if err != nil {
		return resp, err
	}
	versionTraiges, err := repository.SelectVersionTriage(&entity.VersionTriageOption{IssueID: issueID})
	if err != nil {
		return resp, err
	}
	cherryPickToPrs := []entity.PullRequest{}
	if closedPrID != "" {
		results, err := repository.SelectPullRequest(&entity.PullRequestOption{SourcePullRequestID: closedPrID})
		if err != nil {
			return resp, err
		}
		cherryPickToPrs = *results
	}
	for i := range *issueAffects {
		issueAffect := (*issueAffects)[i]
		triagestatus := "unknown"
		for _, versionTraige := range *versionTraiges {
			if versionTraige.VersionName == minorPatchVersionMap[issueAffect.AffectVersion] {
				triagestatus = string(versionTraige.TriageResult)
				break
			}
		}
		patch := ""
		versionsNums := strings.Split(minorPatchVersionMap[issueAffect.AffectVersion], ".")
		if len(versionsNums) == 3 {
			patch = versionsNums[2]
		}

		release := Release{
			BaseVersion:  issueAffect.AffectVersion,
			TriageStatus: strings.ToLower(triagestatus),
			Patch:        patch,
		}
		pr := Pr{}
		for j := range cherryPickToPrs {
			cpr := cherryPickToPrs[j]
			if cpr.MergeTime == nil {
				cpr.MergeTime = &time.Time{}
			}
			if strings.HasSuffix(cpr.BaseBranch, issueAffect.AffectVersion) {
				pr = Pr{
					PrID:        cpr.PullRequestID,
					Repository:  cpr.Repo,
					Number:      cpr.Number,
					Title:       cpr.Title,
					Url:         cpr.HTMLURL,
					State:       cpr.State,
					MergeTarget: cpr.BaseBranch,
					MergedAt:    *cpr.MergeTime,
				}
			}
		}

		affect := &Affect{
			Version: issueAffect.AffectVersion,
			Affect:  strings.ToLower(string((issueAffect.AffectResult))),
			Release: &release,
			PR:      &pr,
		}
		resp = append(resp, affect)
	}
	return resp, nil
}

// func getIssueType(labels []github.Label) string {
// 	resp := ""
// 	for _, label := range labels {
// 		if strings.HasPrefix(*label.Name, "type/") {
// 			resp = strings.ReplaceAll(*label.Name, "type/", "")
// 			return resp
// 		}
// 	}
// 	return resp
// }

// func getIssueSeverity(labels []github.Label) string {
// 	resp := ""
// 	for _, label := range labels {
// 		if strings.HasPrefix(*label.Name, "severity/") {
// 			resp = strings.ReplaceAll(*label.Name, "severity/", "")
// 			return resp
// 		}
// 	}
// 	return resp
// }

// func getAssignee(users []github.User) string {
// 	resp := ""
// 	for _, user := range users {
// 		resp = *user.Name
// 	}
// 	return resp
// }

type IssueInfo struct {
	IssueID   string
	Number    int
	Title     string
	Url       string
	CreatedAt time.Time
	ClosedAt  time.Time
	Severity  string
	State     string
	Type      string
	Assignee  string
	Labels    []string
	Affects   []*Affect
}

type Affect struct {
	Version string
	Affect  string
	PR      *Pr
	Release *Release
}

type Pr struct {
	PrID        string
	Repository  string
	Number      int
	Title       string
	Url         string
	State       string
	MergeTarget string
	MergedAt    time.Time
	Commit      string
	Refs        string
}

type Release struct {
	BaseVersion  string
	TriageStatus string
	Patch        string
}

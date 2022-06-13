package service

import (
	"fmt"
	"strconv"
	"strings"

	"tirelease/commons/git"
	"tirelease/internal/dto"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/pkg/errors"
)

// ============================================================================
// ============================================================================ CURD Of IssueRelationInfo
// Get relation infomations of target issue
// relation infomations include:
// 		1. Issue : Issue basic info
// 		2. IssueAffects : The minor versions affected by the issue
// 		3. IssuePrRelations : The pull requests related to the issue **regardless** of the version**
// 		4. PullRequests	: The pull requests related to the issue **in the version**
// 		5. VersionTriages : The version triage history of the issue
// ============================================================================
// TODO: Decouple the infos of current version from the infos of all versions
//     meta: Issue
//     current version infos: PullRequests
//	   all issue info：IssueAffects, IssuePrRelations, VersionTriages
func FindIssueRelationInfo(option *dto.IssueRelationInfoQuery) (*[]dto.IssueRelationInfo, *entity.ListResponse, error) {
	// select join
	joins, err := repository.SelectIssueRelationInfoByJoin(option)
	if nil != err {
		return nil, nil, err
	}

	count, err := repository.CountIssueRelationInfoByJoin(option)
	if nil != err {
		return nil, nil, err
	}
	response := &entity.ListResponse{
		TotalCount: count,
		Page:       option.IssueOption.Page,
		PerPage:    option.IssueOption.PerPage,
	}
	response.CalcTotalPage()

	// Get all issue ids for further batch select of other entities
	issueIDs := getIssueIDs(*joins)

	// Get all affected minor versions of the issue
	issueAffectAll, err := getIssueAffectVersions(*joins)
	if err != nil {
		return nil, nil, err
	}

	issueAll, err := getIssues(issueIDs)
	if err != nil {
		return nil, nil, err
	}

	// The pull requests related to the issue **regardless** of the version**
	issuePrRelationAll, err := getIssuePrRelation(issueIDs)
	if err != nil {
		return nil, nil, err
	}

	// Get pullrequests whose base branch **regardless** of the version**
	// option.baseBranch
	pullRequestAll, err := getRelatedPullRequests(issuePrRelationAll)
	if err != nil {
		return nil, nil, err
	}

	// Get pullrequests whose base branch **in the version**
	versionPRs := getSameVersionPullRequests(pullRequestAll, option.BaseBranch)

	versionTriageAll, err := getVersionTriages(issueIDs, option.VersionStatus)
	if err != nil {
		return nil, nil, err
	}
	versionTriageAll = appendVersionTriageMergeStatus(versionTriageAll, pullRequestAll, issuePrRelationAll)

	// compose
	issueRelationInfos := composeIssueRelationInfos(issueAll, issueAffectAll, issuePrRelationAll, versionPRs, versionTriageAll)

	return &issueRelationInfos, response, nil
}

func SelectIssueRelationInfoUnique(option *dto.IssueRelationInfoQuery) (*dto.IssueRelationInfo, error) {
	infos, _, err := FindIssueRelationInfo(option)
	if nil != err {
		return nil, err
	}
	if len(*infos) != 1 {
		return nil, errors.New(fmt.Sprintf("more than one issue_relation found: %+v", option))
	}
	return &((*infos)[0]), nil
}

func SaveIssueRelationInfo(issueRelationInfo *dto.IssueRelationInfo) error {

	if issueRelationInfo == nil {
		return nil
	}

	// Save Issue
	if issueRelationInfo.Issue != nil {
		if err := repository.CreateOrUpdateIssue(issueRelationInfo.Issue); nil != err {
			return err
		}
	}

	// Save IssueAffects
	if issueRelationInfo.IssueAffects != nil {
		for _, issueAffect := range *issueRelationInfo.IssueAffects {
			if err := repository.CreateOrUpdateIssueAffect(&issueAffect); nil != err {
				return err
			}
		}
	}

	// Save IssuePrRelations
	if issueRelationInfo.IssuePrRelations != nil {
		for _, issuePrRelation := range *issueRelationInfo.IssuePrRelations {
			if err := repository.CreateIssuePrRelation(&issuePrRelation); nil != err {
				return err
			}
		}
	}

	// Save PullRequests
	if issueRelationInfo.PullRequests != nil {
		for _, pullRequest := range *issueRelationInfo.PullRequests {
			if err := repository.CreateOrUpdatePullRequest(&pullRequest); nil != err {
				return err
			}
		}
	}

	return nil
}

func containVersion(versions *[]entity.ReleaseVersion, name string) bool {
	for _, version := range *versions {
		if version.Name == name {
			return true
		}
	}

	return false
}

func composeIssueRelationInfos(issueAll []entity.Issue, issueAffectAll []entity.IssueAffect,
	issuePrRelationAll []entity.IssuePrRelation, pullRequestAll []entity.PullRequest,
	versionTriageAll []entity.VersionTriage) []dto.IssueRelationInfo {

	// compose
	issueRelationInfos := make([]dto.IssueRelationInfo, 0)
	for index := range issueAll {
		issue := issueAll[index]

		issueRelationInfo := &dto.IssueRelationInfo{}
		issueRelationInfo.Issue = &issue

		issueAffects := make([]entity.IssueAffect, 0)
		if len(issueAffectAll) > 0 {
			for i := range issueAffectAll {
				issueAffect := issueAffectAll[i]
				if issueAffect.IssueID == issue.IssueID {
					issueAffects = append(issueAffects, issueAffect)
				}
			}
		}
		issueRelationInfo.IssueAffects = &issueAffects

		issuePrRelations := make([]entity.IssuePrRelation, 0)
		pullRequests := make([]entity.PullRequest, 0)
		if len(issuePrRelationAll) > 0 {
			for i := range issuePrRelationAll {
				issuePrRelation := issuePrRelationAll[i]
				if issuePrRelation.IssueID != issue.IssueID {
					continue
				}

				issuePrRelations = append(issuePrRelations, issuePrRelation)
				if len(pullRequestAll) > 0 {
					for j := range pullRequestAll {
						pullRequest := pullRequestAll[j]
						if pullRequest.PullRequestID == issuePrRelation.PullRequestID {
							pullRequests = append(pullRequests, pullRequest)
						}
					}
				}
			}
		}
		issueRelationInfo.IssuePrRelations = &issuePrRelations
		issueRelationInfo.PullRequests = &pullRequests

		versionTriages := make([]entity.VersionTriage, 0)
		if len(versionTriageAll) > 0 {
			for i := range versionTriageAll {
				versionTriage := versionTriageAll[i]
				if versionTriage.IssueID == issue.IssueID {
					versionTriages = append(versionTriages, versionTriage)
				}
			}
		}
		issueRelationInfo.VersionTriages = &versionTriages

		issueRelationInfos = append(issueRelationInfos, *issueRelationInfo)
	}

	return issueRelationInfos
}

func getIssueIDs(joins []dto.IssueRelationInfoByJoin) []string {
	issueIDs := make([]string, 0)
	for i := range joins {
		join := joins[i]
		issueIDs = append(issueIDs, join.IssueID)
	}

	return issueIDs
}

func getIssueAffectVersions(joins []dto.IssueRelationInfoByJoin) ([]entity.IssueAffect, error) {
	issueAffectIDs := make([]int64, 0)
	for i := range joins {
		join := (joins)[i]
		ids := strings.Split(join.IssueAffectIDs, ",")
		for _, id := range ids {
			idint, _ := strconv.Atoi(id)
			issueAffectIDs = append(issueAffectIDs, int64(idint))
		}
	}

	issueAffectAll := make([]entity.IssueAffect, 0)

	if len(issueAffectIDs) > 0 {
		issueAffectOption := &entity.IssueAffectOption{
			IDs: issueAffectIDs,
		}
		issueAffectAlls, err := repository.SelectIssueAffect(issueAffectOption)
		if nil != err {
			return nil, err
		}
		issueAffectAll = append(issueAffectAll, (*issueAffectAlls)...)
	}

	return issueAffectAll, nil
}

func getIssues(issueIDs []string) ([]entity.Issue, error) {
	issueAll := make([]entity.Issue, 0)
	if len(issueIDs) > 0 {
		issueOption := &entity.IssueOption{
			IssueIDs: issueIDs,
		}
		issueAlls, err := repository.SelectIssue(issueOption)
		if nil != err {
			return nil, err
		}
		issueAll = append(issueAll, (*issueAlls)...)
	}

	return issueAll, nil
}

func getIssuePrRelation(issueIDs []string) ([]entity.IssuePrRelation, error) {
	issuePrRelationAll := make([]entity.IssuePrRelation, 0)

	if len(issueIDs) > 0 {
		issuePrRelation := &entity.IssuePrRelationOption{
			IssueIDs: issueIDs,
		}
		issuePrRelationAlls, err := repository.SelectIssuePrRelation(issuePrRelation)
		if nil != err {
			return nil, err
		}
		issuePrRelationAll = append(issuePrRelationAll, (*issuePrRelationAlls)...)
	}

	return issuePrRelationAll, nil
}

func getRelatedPullRequests(relatedPrs []entity.IssuePrRelation) ([]entity.PullRequest, error) {
	pullRequestIDs := make([]string, 0)
	pullRequestAll := make([]entity.PullRequest, 0)

	if len(relatedPrs) > 0 {
		for i := range relatedPrs {
			issuePrRelation := relatedPrs[i]
			pullRequestIDs = append(pullRequestIDs, issuePrRelation.PullRequestID)
		}
		pullRequestOption := &entity.PullRequestOption{
			PullRequestIDs: pullRequestIDs,
		}
		pullRequestAlls, err := repository.SelectPullRequest(pullRequestOption)
		if nil != err {
			return nil, err
		}
		pullRequestAll = append(pullRequestAll, (*pullRequestAlls)...)
	}

	return pullRequestAll, nil
}

func getSameVersionPullRequests(pullRequestAll []entity.PullRequest, baseBranch string) []entity.PullRequest {
	if baseBranch == "" {
		return pullRequestAll
	}

	pullRequests := make([]entity.PullRequest, 0)

	if len(pullRequestAll) == 0 {
		return pullRequests
	}

	for i := range pullRequestAll {
		pullRequest := pullRequestAll[i]
		if pullRequest.BaseBranch == baseBranch {
			pullRequests = append(pullRequests, pullRequest)
		}
	}

	return pullRequests
}

func getVersionTriages(issueIDs []string, versionStatus entity.ReleaseVersionStatus) ([]entity.VersionTriage, error) {
	versionTriageAll := make([]entity.VersionTriage, 0)
	if len(issueIDs) > 0 {
		versionTriageOption := &entity.VersionTriageOption{
			IssueIDs: issueIDs,
		}
		versionTriageAlls, err := repository.SelectVersionTriage(versionTriageOption)
		if nil == err && versionStatus == entity.ReleaseVersionStatusUpcoming {
			versionTriageAlls, err = pickUpcomingTriages(versionTriageAlls)
		}
		if nil != err {
			return nil, err
		}

		versionTriageAll = append(versionTriageAll, (*versionTriageAlls)...)
	}

	return versionTriageAll, nil
}

// 只选择对应version状态为“upcoming”的versionTriage
// 由于upcoming状态的version数据量小，因此在查询时不对versionName进行限制
func pickUpcomingTriages(triages *[]entity.VersionTriage) (*[]entity.VersionTriage, error) {
	versionOption := &entity.ReleaseVersionOption{Status: entity.ReleaseVersionStatusUpcoming}
	upcomingVersions, err := repository.SelectReleaseVersion(versionOption)

	if err != nil {
		return nil, err
	}

	upcomingTriages := make([]entity.VersionTriage, 0)

	for i := range *triages {
		triage := (*triages)[i]
		if containVersion(upcomingVersions, triage.VersionName) {
			upcomingTriages = append(upcomingTriages, triage)
		}
	}

	return &upcomingTriages, nil
}

func appendVersionTriageMergeStatus(versionTriages []entity.VersionTriage, pullrequestAll []entity.PullRequest, issuePrRelations []entity.IssuePrRelation) []entity.VersionTriage {
	triages := make([]entity.VersionTriage, 0)
	if len(versionTriages) == 0 {
		return triages
	}

	for i := range versionTriages {
		versionTriage := versionTriages[i]
		prIDs := make([]string, 0)

		for _, relation := range issuePrRelations {
			if relation.IssueID == versionTriage.IssueID {
				prIDs = append(prIDs, relation.PullRequestID)
			}
		}

		major, minor, _, _ := ComposeVersionAtom(versionTriage.VersionName)
		minorVersion := fmt.Sprintf("%d.%d", major, minor)
		baseBranch := git.ReleaseBranchPrefix + minorVersion

		prs := make([]entity.PullRequest, 0)
		for _, pr := range pullrequestAll {
			if pr.BaseBranch != baseBranch {
				continue
			}

			for _, prID := range prIDs {
				if prID == pr.PullRequestID {
					prs = append(prs, pr)
				}
			}
		}

		versionTriage.MergeStatus = ComposeVersionTriageMergeStatus(prs)
		triages = append(triages, versionTriage)
	}

	return triages
}

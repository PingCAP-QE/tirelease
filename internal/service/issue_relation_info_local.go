package service

import (
	"fmt"
	"strconv"
	"strings"

	"tirelease/internal/dto"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/pkg/errors"
)

// ============================================================================
// ============================================================================ CURD Of IssueRelationInfo
func SelectIssueRelationInfo(option *dto.IssueRelationInfoQuery) (*[]dto.IssueRelationInfo, *entity.ListResponse, error) {
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

	// batch select all issue relation info
	issueIDs := make([]string, 0)
	issueAffectIDs := make([]int64, 0)
	pullRequestIDs := make([]string, 0)
	for i := range *joins {
		join := (*joins)[i]
		issueIDs = append(issueIDs, join.IssueID)

		ids := strings.Split(join.IssueAffectIDs, ",")
		for _, id := range ids {
			idint, _ := strconv.Atoi(id)
			issueAffectIDs = append(issueAffectIDs, int64(idint))
		}
	}

	issueAll := make([]entity.Issue, 0)
	issueAffectAll := make([]entity.IssueAffect, 0)
	issuePrRelationAll := make([]entity.IssuePrRelation, 0)
	pullRequestAll := make([]entity.PullRequest, 0)
	versionTriageAll := make([]entity.VersionTriage, 0)

	if len(issueIDs) > 0 {
		issueOption := &entity.IssueOption{
			IssueIDs: issueIDs,
		}
		issueAlls, err := repository.SelectIssue(issueOption)
		if nil != err {
			return nil, nil, err
		}
		issueAll = append(issueAll, (*issueAlls)...)
	}

	if len(issueAffectIDs) > 0 {
		issueAffectOption := &entity.IssueAffectOption{
			IDs: issueAffectIDs,
		}
		issueAffectAlls, err := repository.SelectIssueAffect(issueAffectOption)
		if nil != err {
			return nil, nil, err
		}
		issueAffectAll = append(issueAffectAll, (*issueAffectAlls)...)
	}

	if len(issueIDs) > 0 {
		issuePrRelation := &entity.IssuePrRelationOption{
			IssueIDs: issueIDs,
		}
		issuePrRelationAlls, err := repository.SelectIssuePrRelation(issuePrRelation)
		if nil != err {
			return nil, nil, err
		}
		issuePrRelationAll = append(issuePrRelationAll, (*issuePrRelationAlls)...)
	}

	if len(issuePrRelationAll) > 0 {
		for i := range issuePrRelationAll {
			issuePrRelation := issuePrRelationAll[i]
			pullRequestIDs = append(pullRequestIDs, issuePrRelation.PullRequestID)
		}
		pullRequestOption := &entity.PullRequestOption{
			PullRequestIDs: pullRequestIDs,
		}
		if option.BaseBranch != "" {
			pullRequestOption.BaseBranch = option.BaseBranch
		}
		pullRequestAlls, err := repository.SelectPullRequest(pullRequestOption)
		if nil != err {
			return nil, nil, err
		}
		pullRequestAll = append(pullRequestAll, (*pullRequestAlls)...)
	}

	if len(issueIDs) > 0 {
		versionTriageOption := &entity.VersionTriageOption{
			IssueIDs: issueIDs,
		}
		versionTriageAlls, err := repository.SelectVersionTriage(versionTriageOption)
		if nil == err && option.VersionStatus == entity.ReleaseVersionStatusUpcoming {
			versionTriageAlls, err = pickUpcomingTriages(versionTriageAlls)
		}
		if nil != err {
			return nil, nil, err
		}

		versionTriageAll = append(versionTriageAll, (*versionTriageAlls)...)
	}

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

	// for _, join := range *joins {
	// 	issueRelationInfo := &dto.IssueRelationInfo{}

	// 	// issue
	// 	issueOption := &entity.IssueOption{
	// 		IssueID: join.IssueID,
	// 	}
	// 	issue, err := repository.SelectIssueUnique(issueOption)
	// 	if nil != err {
	// 		return nil, nil, err
	// 	}
	// 	issueRelationInfo.Issue = issue

	// 	// issue_affects
	// 	if join.IssueAffectIDs != "" {
	// 		idList := make([]int64, 0)
	// 		ids := strings.Split(join.IssueAffectIDs, ",")
	// 		for _, id := range ids {
	// 			idint, _ := strconv.Atoi(id)
	// 			idList = append(idList, int64(idint))
	// 		}

	// 		issueAffectOption := &entity.IssueAffectOption{
	// 			IDs: idList,
	// 		}
	// 		issueAffects, err := repository.SelectIssueAffect(issueAffectOption)
	// 		if nil != err {
	// 			return nil, nil, err
	// 		}
	// 		issueRelationInfo.IssueAffects = issueAffects
	// 	}

	// 	// issue_pr_relations
	// 	issuePrRelationOption := &entity.IssuePrRelationOption{
	// 		IssueID: issue.IssueID,
	// 	}
	// 	issuePrRelations, err := repository.SelectIssuePrRelation(issuePrRelationOption)
	// 	if nil != err {
	// 		return nil, nil, err
	// 	}
	// 	issueRelationInfo.IssuePrRelations = issuePrRelations

	// 	// prs
	// 	if issuePrRelations != nil && len(*issuePrRelations) > 0 {
	// 		pullRequestIDs := make([]string, 0)
	// 		for _, issuePrRelation := range *issuePrRelations {
	// 			pullRequestIDs = append(pullRequestIDs, string(issuePrRelation.PullRequestID))
	// 		}

	// 		pullRequestOption := &entity.PullRequestOption{
	// 			PullRequestIDs: pullRequestIDs,
	// 		}
	// 		if option.BaseBranch != "" {
	// 			pullRequestOption.BaseBranch = option.BaseBranch
	// 		}
	// 		pullRequests, err := repository.SelectPullRequest(pullRequestOption)
	// 		if nil != err {
	// 			return nil, nil, err
	// 		}
	// 		issueRelationInfo.PullRequests = pullRequests
	// 	}

	// 	// version_triage
	// 	versionTriageOption := &entity.VersionTriageOption{
	// 		IssueID: issue.IssueID,
	// 	}
	// 	versionTriages, err := repository.SelectVersionTriage(versionTriageOption)
	// 	if nil != err {
	// 		return nil, nil, err
	// 	}
	// 	issueRelationInfo.VersionTriages = versionTriages

	// 	// return
	// 	issueRelationInfos = append(issueRelationInfos, *issueRelationInfo)
	// }

	return &issueRelationInfos, response, nil
}

func SelectIssueRelationInfoUnique(option *dto.IssueRelationInfoQuery) (*dto.IssueRelationInfo, error) {
	infos, _, err := SelectIssueRelationInfo(option)
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

// 只选择对应version状态为“upcoming”的versionTriage
// 由于upcoming状态的version数据量小，因此在查询时不对versionName进行限制
func pickUpcomingTriages(triages *[]entity.VersionTriage) (*[]entity.VersionTriage, error) {
	versionOption := &entity.ReleaseVersionOption{Status: "upcoming"}
	upcomingVersions, err := repository.SelectReleaseVersion(versionOption)

	if err != nil {
		return nil, err
	}

	upcomingTriages := make([]entity.VersionTriage, 0)

	for _, triage := range *triages {
		if containVersion(upcomingVersions, triage.VersionName) {
			upcomingTriages = append(upcomingTriages, triage)
		}
	}

	return &upcomingTriages, nil
}

func containVersion(versions *[]entity.ReleaseVersion, name string) bool {
	for _, version := range *versions {
		if version.Name == name {
			return true
		}
	}

	return false
}

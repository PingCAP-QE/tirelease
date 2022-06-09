package service

import (
	"strings"

	"tirelease/commons/git"
	"tirelease/internal/dto"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/pkg/errors"
)

func CreateOrUpdateVersionTriageInfo(versionTriage *entity.VersionTriage, updatedVars ...entity.VersionTriageUpdatedVar) (*dto.VersionTriageInfo, error) {
	// find version
	releaseVersion, err := SelectReleaseVersionActive(versionTriage.VersionName)
	if err != nil {
		return nil, err
	}
	releaseBranch := releaseVersion.ReleaseBranch
	versionTriage.VersionName = releaseVersion.Name

	// basic info
	issueRelationInfo, err := SelectIssueRelationInfoUnique(&dto.IssueRelationInfoQuery{
		IssueOption: entity.IssueOption{
			IssueID: versionTriage.IssueID,
		},
		BaseBranch: releaseBranch,
	})
	if err != nil {
		return nil, err
	}

	storedVersionTriage, err := repository.SelectVersionTriage(&entity.VersionTriageOption{
		IssueID:     versionTriage.IssueID,
		VersionName: releaseVersion.Name,
	})

	if err != nil {
		return nil, err
	}
	// create or update
	var isFrozen bool = releaseVersion.Status == entity.ReleaseVersionStatusFrozen
	var isAccept bool = versionTriage.TriageResult == entity.VersionTriageResultAccept
	if isFrozen && isAccept {
		versionTriage.TriageResult = entity.VersionTriageResultAcceptFrozen
		updatedVars = append(updatedVars, entity.VersionTriageUpdatedVarTriageResult)
	}

	// 当issue为critical，且数据库中没有数据或BlockVersionRelease字段为空时，设置默认值为Block
	if issueRelationInfo.Issue.SeverityLabel == git.SeverityCriticalLabel && (len(*storedVersionTriage) == 0 || (*storedVersionTriage)[0].BlockVersionRelease == "") {
		versionTriage.BlockVersionRelease = entity.BlockVersionReleaseResultBlock
		updatedVars = append(updatedVars, entity.VersionTriageUpdatedVarBlockRelease)
	}
	if err := repository.CreateOrUpdateVersionTriage(versionTriage, updatedVars...); err != nil {
		return nil, err
	}

	// remote operation
	if issueRelationInfo != nil && issueRelationInfo.PullRequests != nil && len(*issueRelationInfo.PullRequests) > 0 {
		for i := range *issueRelationInfo.PullRequests {
			pr := (*issueRelationInfo.PullRequests)[i]
			err := ChangePrApprovedLabel(pr.PullRequestID, isFrozen, isAccept)
			if err != nil {
				return nil, err
			}
		}
	}

	// return
	return &dto.VersionTriageInfo{
		ReleaseVersion: releaseVersion,
		IsFrozen:       isFrozen,
		IsAccept:       isAccept,

		VersionTriage:            versionTriage,
		VersionTriageMergeStatus: ComposeVersionTriageMergeStatus(issueRelationInfo),

		IssueRelationInfo: issueRelationInfo,
	}, nil
}

// Hook github api to change the `approve` related labels.
func ChangePrApprovedLabel(prId string, isFrozen, isAccept bool) error {
	if !isFrozen && isAccept {
		err := RemoveLabelByPullRequestID(prId, git.NotCheryyPickLabel)
		if err != nil {
			return err
		}

		err = AddLabelByPullRequestID(prId, git.CherryPickLabel)
		if err != nil {
			return err
		}
	} else {
		err := RemoveLabelByPullRequestID(prId, git.CherryPickLabel)
		if err != nil {
			return err
		}

		err = AddLabelByPullRequestID(prId, git.NotCheryyPickLabel)
		if err != nil {
			return err
		}
	}

	return nil
}

func SelectVersionTriageInfo(query *dto.VersionTriageInfoQuery) (*dto.VersionTriageInfoWrap, *entity.ListResponse, error) {

	// dependency
	releaseVersion, err := repository.SelectReleaseVersionLatest(&entity.ReleaseVersionOption{
		Name: query.Version,
	})
	if err != nil {
		return nil, nil, err
	}

	// compose
	var versionTriages []entity.VersionTriage
	if releaseVersion.Status == entity.ReleaseVersionStatusUpcoming || releaseVersion.Status == entity.ReleaseVersionStatusFrozen {
		versionTriages, err = ComposeVersionTriageUpcomingList(query.Version)
		if err != nil {
			return nil, nil, err
		}
	} else {
		versionTriagesPoint, err := repository.SelectVersionTriage(&query.VersionTriageOption)
		if err != nil {
			return nil, nil, err
		}
		versionTriages = *versionTriagesPoint
	}

	// detail
	issueIDs := make([]string, 0)
	for i := range versionTriages {
		versionTriage := versionTriages[i]
		issueIDs = append(issueIDs, versionTriage.IssueID)
	}
	versionTriageInfos := make([]dto.VersionTriageInfo, 0)
	if len(issueIDs) > 0 {
		infoOption := &dto.IssueRelationInfoQuery{
			IssueOption: entity.IssueOption{
				IssueIDs: issueIDs,
			},
			BaseBranch: releaseVersion.ReleaseBranch,
		}
		issueRelationInfos, _, err := SelectIssueRelationInfo(infoOption)
		if err != nil {
			return nil, nil, err
		}

		for i := range versionTriages {

			versionTriage := versionTriages[i]
			versionTriageInfo := dto.VersionTriageInfo{}
			versionTriageInfo.VersionTriage = &versionTriage
			versionTriageInfo.ReleaseVersion = releaseVersion

			for j := range *issueRelationInfos {
				issueRelationInfo := (*issueRelationInfos)[j]
				if issueRelationInfo.Issue.IssueID == versionTriage.IssueID {
					versionTriageInfo.IssueRelationInfo = &issueRelationInfo
					versionTriageInfo.VersionTriageMergeStatus = ComposeVersionTriageMergeStatus(&issueRelationInfo)
					fillVersionTriageDefaultValue(&issueRelationInfo, &versionTriage)
					break
				}
			}

			versionTriageInfos = append(versionTriageInfos, versionTriageInfo)
		}
	}

	// versionTriageInfos := make([]dto.VersionTriageInfo, 0)
	// for i := range versionTriages {
	// 	versionTriage := versionTriages[i]
	// 	issueRelationInfo, err := SelectIssueRelationInfoUnique(&dto.IssueRelationInfoQuery{
	// 		IssueOption: entity.IssueOption{
	// 			IssueID: versionTriage.IssueID,
	// 		},
	// 		BaseBranch: releaseVersion.ReleaseBranch,
	// 	})
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}

	// 	versionTriageInfo := dto.VersionTriageInfo{}
	// 	versionTriageInfo.VersionTriage = &versionTriage
	// 	versionTriageInfo.IssueRelationInfo = issueRelationInfo
	// 	versionTriageInfo.ReleaseVersion = releaseVersion
	// 	versionTriageInfo.VersionTriageMergeStatus = ComposeVersionTriageMergeStatus(issueRelationInfo)
	// 	versionTriageInfos = append(versionTriageInfos, versionTriageInfo)
	// }

	// return
	wrap := &dto.VersionTriageInfoWrap{
		ReleaseVersion:     releaseVersion,
		VersionTriageInfos: &versionTriageInfos,
	}
	response := &entity.ListResponse{
		Page:    query.VersionTriageOption.Page,
		PerPage: query.VersionTriageOption.PerPage,
	}
	response.CalcTotalPage()
	return wrap, response, nil
}

func UpdateVersionTriage(versionTriage *entity.VersionTriage) error {
	if versionTriage.ID == 0 {
		versionTriage.TriageResult = entity.VersionTriageResultUnKnown
		err := repository.CreateVersionTriage(versionTriage)
		if err != nil {
			return err
		}
	} else {
		err := repository.UpdateVersionTriage(versionTriage)
		if err != nil {
			return err
		}
	}
	return nil
}

func InheritVersionTriage(fromVersion string, toVersion string) error {
	// Select
	versionTriageOption := &entity.VersionTriageOption{
		VersionName: fromVersion,
	}
	versionTriages, err := repository.SelectVersionTriage(versionTriageOption)
	if err != nil {
		return err
	}
	if len(*versionTriages) == 0 {
		return nil
	}

	// Migrate
	for i := range *versionTriages {
		versionTriage := (*versionTriages)[i]
		switch versionTriage.TriageResult {
		case entity.VersionTriageResultAccept:
			versionTriage.TriageResult = entity.VersionTriageResultReleased
		case entity.VersionTriageResultUnKnown:
			versionTriage.VersionName = toVersion
		case entity.VersionTriageResultLater, entity.VersionTriageResultAcceptFrozen:
			versionTriage.VersionName = toVersion
			versionTriage.TriageResult = entity.VersionTriageResultAccept
		}
		if err := repository.CreateOrUpdateVersionTriage(&versionTriage); err != nil {
			return err
		}
	}

	return nil
}

func ComposeVersionTriageMergeStatus(issueRelationInfo *dto.IssueRelationInfo) entity.VersionTriageMergeStatus {
	if issueRelationInfo.PullRequests == nil || len(*issueRelationInfo.PullRequests) == 0 {
		return entity.VersionTriageMergeStatusPr
	}
	allMerge := true
	closeNums := 0
	for _, pr := range *issueRelationInfo.PullRequests {
		if pr.State == "closed" {
			closeNums++
			continue
		}

		//TODO: 当前存在approve成功hook到git，但是数据库中状态不一致的问题
		// 这里先兼容该情况，认为merge后的pr都是已approve过的，待重新设计状态机后修改逻辑
		if pr.Merged {
			continue
		}

		if !pr.CherryPickApproved {
			return entity.VersionTriageMergeStatusApprove
		} else if !pr.AlreadyReviewed {
			return entity.VersionTriageMergeStatusReview
		} else if !pr.Merged {
			allMerge = false
		}
	}
	if closeNums == len(*issueRelationInfo.PullRequests) {
		return entity.VersionTriageMergeStatusPr
	}
	if allMerge {
		return entity.VersionTriageMergeStatusMerged
	} else {
		return entity.VersionTriageMergeStatusCITesting
	}
}

func ComposeVersionTriageUpcomingList(version string) ([]entity.VersionTriage, error) {
	// select all issue which may affect this minor version
	major, minor, _, _ := ComposeVersionAtom(version)
	minorVersion := ComposeVersionMinorNameByNumber(major, minor)
	affectOption := &entity.IssueAffectOption{
		AffectVersion: minorVersion,
		AffectResult:  entity.AffectResultResultYes,
	}
	issueAffects, err := repository.SelectIssueAffect(affectOption)
	if err != nil {
		return nil, err
	}

	// select all triaged list under this minor version
	versionOption := &entity.ReleaseVersionOption{
		Major:     major,
		Minor:     minor,
		ShortType: entity.ReleaseVersionShortTypeMinor,
	}
	releaseVersions, err := repository.SelectReleaseVersion(versionOption)
	if err != nil {
		return nil, err
	}
	versions := make([]string, 0)
	for i := range *releaseVersions {
		versions = append(versions, (*releaseVersions)[i].Name)
	}

	versionTriageOption := &entity.VersionTriageOption{
		VersionNameList: versions,
	}
	versionTriageData, err := repository.SelectVersionTriage(versionTriageOption)
	if err != nil {
		return nil, err
	}

	// compose: version_triage = affected - triaged
	versionTriages := make([]entity.VersionTriage, 0)
	for i := range *issueAffects {
		issueAffect := (*issueAffects)[i]
		find := false
		for j := range *versionTriageData {
			versionTriage := (*versionTriageData)[j]
			if issueAffect.IssueID != versionTriage.IssueID {
				continue
			}
			find = true

			if versionTriage.TriageResult == entity.VersionTriageResultReleased ||
				versionTriage.TriageResult == entity.VersionTriageResultWontFix ||
				versionTriage.TriageResult == entity.VersionTriageResultLater {
				if version != versionTriage.VersionName {
					continue
				}
			}
			versionTriages = append(versionTriages, versionTriage)
		}
		if !find {
			versionTriage := entity.VersionTriage{
				IssueID:      issueAffect.IssueID,
				VersionName:  version,
				TriageResult: entity.VersionTriageResultUnKnown,
			}
			versionTriages = append(versionTriages, versionTriage)
		}
	}
	return versionTriages, nil
}

// Export history data (Only database operation, no remote operation)
func ExportHistoryVersionTriageInfo(info *dto.IssueRelationInfo, releaseVersions *[]entity.ReleaseVersion) error {
	// param check
	if info == nil || releaseVersions == nil {
		return errors.New("ExportHistoryVersionTriageInfo params invalid")
	}
	if info.PullRequests == nil || len(*info.PullRequests) == 0 {
		return nil
	}

	// insert version triage
	for i := range *info.PullRequests {
		pr := (*info.PullRequests)[i]
		if !pr.Merged || !pr.CherryPickApproved ||
			!strings.HasPrefix(pr.BaseBranch, git.ReleaseBranchPrefix) {
			continue
		}
		releaseBranch := string(pr.BaseBranch)
		branchVersion := strings.Replace(pr.BaseBranch, git.ReleaseBranchPrefix, "", -1)
		major, minor, _, _ := ComposeVersionAtom(branchVersion)

		// search version in time section
		// release version is already sorted desc
		for i := len(*releaseVersions) - 1; i >= 0; i-- {
			releaseVersion := (*releaseVersions)[i]
			if releaseVersion.Status != entity.ReleaseVersionStatusReleased {
				continue
			}
			if releaseVersion.Major != major || releaseVersion.Minor != minor || releaseVersion.ReleaseBranch != releaseBranch {
				continue
			}
			if releaseVersion.ActualReleaseTime.After(*(pr.MergeTime)) {
				versionTriage := &entity.VersionTriage{
					IssueID:      info.Issue.IssueID,
					VersionName:  releaseVersion.Name,
					TriageResult: entity.VersionTriageResultReleased,
					CreateTime:   *(pr.MergeTime),
					UpdateTime:   *(pr.MergeTime),
				}
				if err := repository.CreateOrUpdateVersionTriage(versionTriage); err != nil {
					return err
				}

				issueAffect := &entity.IssueAffect{
					IssueID:       info.Issue.IssueID,
					AffectVersion: branchVersion,
					AffectResult:  entity.AffectResultResultYes,
					CreateTime:    *(pr.MergeTime),
					UpdateTime:    *(pr.MergeTime),
				}
				if err := repository.CreateOrUpdateIssueAffect(issueAffect); err != nil {
					return err
				}

				break
			}
		}
	}

	return nil
}

func fillVersionTriageDefaultValue(issueRelationInfo *dto.IssueRelationInfo, versionTriage *entity.VersionTriage) {
	if len(*issueRelationInfo.VersionTriages) == 0 {
		if issueRelationInfo.Issue.SeverityLabel == git.SeverityCriticalLabel {
			versionTriage.BlockVersionRelease = entity.BlockVersionReleaseResultBlock
		}
	}
}

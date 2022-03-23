package service

import (
	"fmt"

	// "tirelease/commons/git"
	"tirelease/commons/git"
	"tirelease/internal/dto"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/pkg/errors"
)

func CreateOrUpdateVersionTriageInfo(versionTriage *entity.VersionTriage) (*dto.VersionTriageInfo, error) {
	// check
	shortType := ComposeVersionShortType(versionTriage.VersionName)
	major, minor, patch, _ := ComposeVersionAtom(versionTriage.VersionName)
	releaseVersionOption := &entity.ReleaseVersionOption{}
	if shortType == entity.ReleaseVersionShortTypeMinor {
		releaseVersionOption.Major = major
		releaseVersionOption.Minor = minor
		releaseVersionOption.Status = entity.ReleaseVersionStatusUpcoming
	} else if shortType == entity.ReleaseVersionShortTypePatch || shortType == entity.ReleaseVersionShortTypeHotfix {
		releaseVersionOption.Major = major
		releaseVersionOption.Minor = minor
		releaseVersionOption.Patch = patch
	} else {
		return nil, errors.New(fmt.Sprintf("CreateOrUpdateVersionTriageInfo params invalid: %+v failed", versionTriage))
	}
	releaseVersion, err := CheckReleaseVersion(releaseVersionOption)
	if err != nil {
		return nil, err
	}
	releaseBranch := releaseVersion.ReleaseBranch
	versionTriage.VersionName = releaseVersion.Name

	// basic info
	issueRelationInfo, err := SelectIssueRelationInfoUnique(&dto.IssueRelationInfoQuery{
		IssueID:    versionTriage.IssueID,
		BaseBranch: releaseBranch,
	})
	if err != nil {
		return nil, err
	}

	// create or update
	var isFrozen bool = releaseVersion.Status == entity.ReleaseVersionStatusFrozen
	var isAccept bool = versionTriage.TriageResult == entity.VersionTriageResultAccept
	if isFrozen && isAccept {
		versionTriage.TriageResult = entity.VersionTriageResultAcceptFrozen
	}
	if issueRelationInfo.Issue.SeverityLabel == git.SeverityCriticalLabel {
		versionTriage.BlockVersionRelease = entity.BlockVersionReleaseResultBlock
	}
	if err := repository.CreateOrUpdateVersionTriage(versionTriage); err != nil {
		return nil, err
	}

	// remote operation
	if len(*issueRelationInfo.PullRequests) == 1 {
		pr := (*issueRelationInfo.PullRequests)[0]
		if !isFrozen && isAccept {
			err := RemoveLabelByPullRequestID(pr.PullRequestID, git.NotCheryyPickLabel)
			if err != nil {
				return nil, err
			}

			err = AddLabelByPullRequestID(pr.PullRequestID, git.CherryPickLabel)
			if err != nil {
				return nil, err
			}
		}
		var isRelease bool = releaseVersion.Status == entity.ReleaseVersionStatusReleased
		if !isAccept && !isRelease {
			err := RemoveLabelByPullRequestID(pr.PullRequestID, git.CherryPickLabel)
			if err != nil {
				return nil, err
			}

			err = AddLabelByPullRequestID(pr.PullRequestID, git.NotCheryyPickLabel)
			if err != nil {
				return nil, err
			}
		}
	}

	// Return
	return &dto.VersionTriageInfo{
		ReleaseVersion: releaseVersion,

		VersionTriage: versionTriage,
		IsFrozen:      isFrozen,
		IsAccept:      isAccept,

		IssueRelationInfo: issueRelationInfo,
	}, nil
}

func SelectVersionTriageInfo(query *dto.VersionTriageInfoQuery) (*dto.VersionTriageInfoWrap, error) {
	// Option
	versionTriageOption := &entity.VersionTriageOption{
		ID:           query.ID,
		IssueID:      query.IssueID,
		VersionName:  query.Version,
		TriageResult: query.TriageResult,
	}

	// Select
	versionTriages, err := repository.SelectVersionTriage(versionTriageOption)
	if err != nil {
		return nil, err
	}

	releaseVersion, err := repository.SelectReleaseVersionLatest(&entity.ReleaseVersionOption{
		Name: query.Version,
	})
	if err != nil {
		return nil, err
	}

	// Compose
	versionTriageInfos := make([]dto.VersionTriageInfo, 0)
	for i := range *versionTriages {
		versionTriage := (*versionTriages)[i]
		issueRelationInfos, err := SelectIssueRelationInfo(&dto.IssueRelationInfoQuery{
			IssueID:    versionTriage.IssueID,
			BaseBranch: releaseVersion.ReleaseBranch,
		})
		if err != nil {
			return nil, err
		}

		versionTriageInfo := dto.VersionTriageInfo{}
		versionTriageInfo.VersionTriage = &versionTriage
		if len(*issueRelationInfos) == 1 {
			versionTriageInfo.IssueRelationInfo = &(*issueRelationInfos)[0]
		}
		versionTriageInfo.ReleaseVersion = releaseVersion
		versionTriageInfos = append(versionTriageInfos, versionTriageInfo)
	}

	return &dto.VersionTriageInfoWrap{
		ReleaseVersion:     releaseVersion,
		VersionTriageInfos: &versionTriageInfos,
	}, nil
}

func CheckReleaseVersion(option *entity.ReleaseVersionOption) (*entity.ReleaseVersion, error) {
	if option == nil || option.Name == "" {
		return nil, errors.New(fmt.Sprintf("CheckReleaseVersion params invalid: %+v failed", option))
	}
	releaseVersion, err := repository.SelectReleaseVersionLatest(option)
	if err != nil {
		return nil, err
	}
	if releaseVersion.Status == entity.ReleaseVersionStatusReleased || releaseVersion.Status == entity.ReleaseVersionStatusCancelled {
		return nil, errors.Wrap(err, fmt.Sprintf("find release version is already released or cancelled: %+v failed", releaseVersion))
	}
	return releaseVersion, nil
}

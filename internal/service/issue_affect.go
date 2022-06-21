package service

import (
	"fmt"
	"strings"

	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

// ============================================================================
// ============================================================================ Restful API(From UI) Handler
func CreateOrUpdateIssueAffect(issueAffect *entity.IssueAffect) error {
	// create or update
	err := repository.CreateOrUpdateIssueAffect(issueAffect)
	if err != nil {
		return err
	}

	// result operation
	err = OperateIssueAffectResult(issueAffect)
	if err != nil {
		return err
	}

	return nil
}

func OperateIssueAffectResult(issueAffect *entity.IssueAffect) error {
	// param protection
	if issueAffect.AffectResult == "" {
		return nil
	}

	// operate git label
	affectLabel := fmt.Sprintf(git.AffectsLabel, issueAffect.AffectVersion)
	mayAffectLabel := fmt.Sprintf(git.MayAffectsLabel, issueAffect.AffectVersion)
	if issueAffect.AffectResult == entity.AffectResultResultUnKnown {
		err := RemoveLabelByIssueID(issueAffect.IssueID, affectLabel)
		if err != nil {
			return err
		}

		err = AddLabelByIssueID(issueAffect.IssueID, mayAffectLabel)
		if err != nil {
			return err
		}
	}
	if issueAffect.AffectResult == entity.AffectResultResultYes {
		err := RemoveLabelByIssueID(issueAffect.IssueID, mayAffectLabel)
		if err != nil {
			return err
		}

		err = AddLabelByIssueID(issueAffect.IssueID, affectLabel)
		if err != nil {
			return err
		}
	}
	if issueAffect.AffectResult == entity.AffectResultResultNo {
		err := RemoveLabelByIssueID(issueAffect.IssueID, mayAffectLabel)
		if err != nil {
			return err
		}

		err = RemoveLabelByIssueID(issueAffect.IssueID, affectLabel)
		if err != nil {
			return err
		}
	}

	// operate cherry-pick record: select latest version & insert cherry-pick
	// if issueAffect.AffectResult == entity.AffectResultResultYes {
	// 	major, minor, _, _ := ComposeVersionAtom(issueAffect.AffectVersion)
	// 	versionTriage := &entity.VersionTriage{
	// 		IssueID:      issueAffect.IssueID,
	// 		VersionName:  ComposeVersionMinorName(&entity.ReleaseVersion{Major: major, Minor: minor}),
	// 		TriageResult: entity.VersionTriageResultUnKnown,
	// 	}
	// 	_, err := CreateOrUpdateVersionTriageInfo(versionTriage)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// ============================================================================
// ============================================================================ Compose From DataBase & Return To UI
func ComposeIssueAffectWithIssueID(issueID string, releaseVersions *[]entity.ReleaseVersion) (*[]entity.IssueAffect, error) {
	// Select Exist Issue Affect
	issueAffectOption := &entity.IssueAffectOption{
		IssueID: issueID,
	}
	issueAffects, err := repository.SelectIssueAffect(issueAffectOption)
	if err != nil {
		return nil, err
	}

	// Implement New Issue Affect
	if releaseVersions == nil {
		releaseVersionOption := &entity.ReleaseVersionOption{
			Status: entity.ReleaseVersionStatusUpcoming,
		}
		releaseVersions, err = repository.SelectReleaseVersion(releaseVersionOption)
		if nil != err {
			return nil, err
		}
	}
	for i := range *releaseVersions {
		releaseVersion := (*releaseVersions)[i]
		var isExist bool = false
		for _, issueAffect := range *issueAffects {
			major, minor, _, _ := ComposeVersionAtom(issueAffect.AffectVersion)
			if major == releaseVersion.Major && minor == releaseVersion.Minor {
				isExist = true
				break
			}
		}

		if !isExist {
			newAffect := &entity.IssueAffect{
				IssueID:       issueID,
				AffectVersion: ComposeVersionMinorName(&releaseVersion),
				AffectResult:  entity.AffectResultResultUnKnown,
			}
			(*issueAffects) = append((*issueAffects), *newAffect)
		}
	}
	return issueAffects, nil
}

// ============================================================================
// ============================================================================ Compose From Webhook (Save To DataBase)
func ComposeIssueAffectWithIssueV4(issue *git.IssueField) (*[]entity.IssueAffect, error) {
	if nil == issue || len(issue.Labels.Nodes) == 0 {
		return nil, nil
	}

	issueAffects := make([]entity.IssueAffect, 0)
	// github affect label: set Yes or UnKnown
	for i := range issue.Labels.Nodes {
		label := issue.Labels.Nodes[i]
		labelName := string(label.Name)
		if strings.HasPrefix(labelName, git.AffectsPrefixLabel) {
			version := strings.Replace(labelName, git.AffectsPrefixLabel, "", -1)
			issueAffect := entity.IssueAffect{
				IssueID:       issue.ID.(string),
				AffectVersion: version,
				AffectResult:  entity.AffectResultResultYes,
			}
			issueAffects = append(issueAffects, issueAffect)
		}
		if strings.HasPrefix(labelName, git.MayAffectsPrefixLabel) {
			version := strings.Replace(labelName, git.MayAffectsPrefixLabel, "", -1)
			issueAffect := entity.IssueAffect{
				IssueID:       issue.ID.(string),
				AffectVersion: version,
				AffectResult:  entity.AffectResultResultUnKnown,
			}
			issueAffects = append(issueAffects, issueAffect)
		}
	}

	// github del affect label: set No
	issueAffectOption := &entity.IssueAffectOption{
		IssueID: issue.ID.(string),
	}
	oldAffects, err := repository.SelectIssueAffect(issueAffectOption)
	if err != nil {
		return nil, err
	}
	for i := range *oldAffects {
		oldAffect := (*oldAffects)[i]
		var isExist bool = false
		for _, issueAffect := range issueAffects {
			if issueAffect.AffectVersion == oldAffect.AffectVersion {
				isExist = true
				break
			}
		}
		if !isExist {
			issueAffect := entity.IssueAffect{
				IssueID:       issue.ID.(string),
				AffectVersion: oldAffect.AffectVersion,
				AffectResult:  entity.AffectResultResultNo,
			}
			issueAffects = append(issueAffects, issueAffect)
		}
	}

	return &issueAffects, nil
}

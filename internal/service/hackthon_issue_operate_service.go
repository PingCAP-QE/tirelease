package service

import (
	// "fmt"
	// "strings"
	// "time"

	"tirelease/internal/entity"
	// "tirelease/internal/repository"
)

func IssueAffectOperate(updateOption *entity.IssueAffectUpdateOption) error {
	// Params Check
	// option := &entity.IssueAffectOption{}
	// if updateOption.ID != 0 {
	// 	option.ID = updateOption.ID
	// } else if updateOption.IssueID != "" && updateOption.AffectVersion != "" {
	// 	option.IssueID = updateOption.IssueID
	// 	option.AffectVersion = updateOption.AffectVersion
	// } else {
	// 	return fmt.Errorf("update option params key is empty: %+v failed", updateOption)
	// }
	// if updateOption.AffectResult == "" {
	// 	return fmt.Errorf("update option params values is empty: %+v failed", updateOption)
	// }

	// // Select
	// issueAffects, err := repository.SelectIssueAffect(option)
	// if nil != err {
	// 	return err
	// }
	// if len(*issueAffects) != 1 {
	// 	return fmt.Errorf("select issue affect is empty or many: %+v failed", *issueAffects)
	// }

	// // Update issue_affect
	// (*issueAffects)[0].AffectResult = updateOption.AffectResult
	// err = repository.CreateOrUpdateIssueAffect(&(*issueAffects)[0])
	// if nil != err {
	// 	return err
	// }
	// releaseType := entity.ReleaseVersionTypePatch
	// status := entity.ReleaseVersionStatusOpen
	// versions, err := repository.SelectReleaseVersion(&entity.ReleaseVersionOption{Type: releaseType, Status: status})
	// if err != nil {
	// 	return err
	// }
	// patchVersion := ""
	// for _, version := range *versions {
	// 	if version.FatherReleaseVersionName == option.AffectVersion {
	// 		patchVersion = version.Name
	// 	}
	// }

	// // Insert version_triage
	// if string((*issueAffects)[0].AffectResult) == strings.ToLower(string(entity.AffectResultResultYes)) {
	// 	versionTriage := &entity.VersionTriage{
	// 		VersionName:  patchVersion,
	// 		IssueID:      (*issueAffects)[0].IssueID,
	// 		TriageResult: entity.VersionTriageResultUnKnown,
	// 		CreateTime:   time.Now(),
	// 		UpdateTime:   time.Now(),
	// 	}
	// 	err = repository.CreateVersionTriage(versionTriage)
	// 	if nil != err {
	// 		return err
	// 	}
	// }

	return nil
}

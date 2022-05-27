package repository

import (
	"encoding/json"
	"fmt"

	"tirelease/commons/database"
	"tirelease/internal/entity"
	"tirelease/internal/service/component"

	"github.com/google/go-github/v41/github"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func SelectIssue(option *entity.IssueOption) (*[]entity.Issue, error) {
	sql := "select * from issue where 1=1" + IssueWhere(option) + IssueOrderBy(option) + IssueLimit(option)

	// 查询
	var issues []entity.Issue
	if err := database.DBConn.RawWrapper(sql, option).Find(&issues).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("select issue by raw failed, option: %+v", option))
	}

	// 加工
	for i := 0; i < len(issues); i++ {
		unSerializeIssue(&issues[i])
	}

	return &issues, nil
}

func CountIssue(option *entity.IssueOption) (int64, error) {
	sql := "select count(*) from issue where 1=1" + IssueWhere(option)

	var count int64
	if err := database.DBConn.RawWrapper(sql, option).Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("count issue by raw failed, option: %+v", option))
	}

	return count, nil
}

func SelectIssueUnique(option *entity.IssueOption) (*entity.Issue, error) {
	// 查询
	issues, err := SelectIssue(option)
	if err != nil {
		return nil, err
	}

	// 校验
	if len(*issues) == 0 {
		return nil, errors.New(fmt.Sprintf("issue not found: %+v", option))
	}
	if len(*issues) > 1 {
		return nil, errors.New(fmt.Sprintf("more than one issue found: %+v", option))
	}

	return &((*issues)[0]), nil
}

func CreateOrUpdateIssue(issue *entity.Issue) error {
	// 加工
	serializeIssue(issue)

	// 存储
	if err := database.DBConn.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Omit("Labels", "Assignees").Create(&issue).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create or update issue: %+v failed", issue))
	}

	return nil
}

// func DeleteIssue(issue *entity.Issue) error {
// 	if err := database.DBConn.DB.Delete(issue).Error; err != nil {
// 		return errors.Wrap(err, fmt.Sprintf("delete issue: %+v failed", issue))
// 	}
// 	return nil
// }

// func CreateIssue(issue *entity.Issue) error {
// 	// 加工
// 	serializeIssue(issue)

// 	// 存储
// 	if err := database.DBConn.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&issue).Error; err != nil {
// 		return errors.Wrap(err, fmt.Sprintf("create issue: %+v failed", issue))
// 	}
// 	return nil
// }

// 序列化和反序列化
func serializeIssue(issue *entity.Issue) {
	if nil != issue.Assignees {
		assigneesString, _ := json.Marshal(issue.Assignees)
		issue.AssigneesString = string(assigneesString)
	}
	if nil != issue.Labels {
		labelsString, _ := json.Marshal(issue.Labels)
		issue.LabelsString = string(labelsString)
	}
}

// TODO Add a layer between repository and the service and move below logic to the layer.
func unSerializeIssue(issue *entity.Issue) {
	if issue.AssigneesString != "" {
		var assignees []github.User
		json.Unmarshal([]byte(issue.AssigneesString), &assignees)
		issue.Assignees = &assignees
	}
	if issue.LabelsString != "" {
		var labels []github.Label
		json.Unmarshal([]byte(issue.LabelsString), &labels)
		issue.Labels = &labels
	}
	issue.Components = component.GetComponents(issue.Owner, issue.Repo, issue.LabelsString)
}

func IssueWhere(option *entity.IssueOption) string {
	sql := ""

	if option.ID != 0 {
		sql += " and issue.id = @ID"
	}
	if option.IssueID != "" {
		sql += " and issue.issue_id = @IssueID"
	}
	if option.Number != 0 {
		sql += " and issue.number = @Number"
	}
	if option.State != "" {
		sql += " and issue.state = @State"
	}
	if option.Owner != "" {
		sql += " and issue.owner = @Owner"
	}
	if option.Repo != "" {
		sql += " and issue.repo = @Repo"
	}
	if option.SeverityLabel != "" {
		sql += " and issue.severity_label = @SeverityLabel"
	}
	if option.TypeLabel != "" {
		sql += " and issue.type_label = @TypeLabel"
	}
	if !option.CreateTime.IsZero() {
		sql += " and issue.create_time > @CreateTime"
	}
	if !option.UpdateTime.IsZero() {
		sql += " and issue.update_time > @UpdateTime"
	}
	if !option.CloseTime.IsZero() {
		sql += " and issue.close_time > @CloseTime"
	}
	if option.IssueIDs != nil && len(option.IssueIDs) > 0 {
		sql += " and issue.issue_id in @IssueIDs"
	}
	if option.SeverityLabels != nil && len(option.SeverityLabels) > 0 {
		sql += " and issue.severity_label in @SeverityLabels"
	}
	if option.NotSeverityLabels != nil && len(option.NotSeverityLabels) > 0 {
		sql += " and issue.severity_label not in @NotSeverityLabels"
	}

	return sql
}

func IssueOrderBy(option *entity.IssueOption) string {
	sql := ""

	if option.OrderBy != "" {
		sql += " order by " + option.OrderBy
	}
	if option.Order != "" {
		sql += " " + option.Order
	}
	if option.OrderBy == "" && option.Order == "" {
		sql += " order by issue.update_time desc, issue.id desc"
	}

	return sql
}

func IssueLimit(option *entity.IssueOption) string {
	sql := ""

	if option.Page != 0 && option.PerPage != 0 {
		option.ListOption.CalcOffset()
		sql += " limit @Offset,@PerPage"
	}

	return sql
}

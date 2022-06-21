package repository

import (
	"fmt"
	"time"

	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func CreateIssuePrRelation(issuePrRelation *entity.IssuePrRelation) error {
	if issuePrRelation.CreateTime.IsZero() {
		issuePrRelation.CreateTime = time.Now()
	}
	if issuePrRelation.UpdateTime.IsZero() {
		issuePrRelation.UpdateTime = time.Now()
	}
	// 存储
	if err := database.DBConn.DB.Clauses(
		clause.OnConflict{DoNothing: true}).Create(&issuePrRelation).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create issue_pr_relation: %+v failed", issuePrRelation))
	}
	return nil
}

func SelectIssuePrRelation(option *entity.IssuePrRelationOption) (*[]entity.IssuePrRelation, error) {
	sql := "select * from issue_pr_relation where 1=1" + IssuePrRelationWhere(option) + IssuePrRelationOrderBy(option) + IssuePrRelationLimit(option)
	// 查询
	var issuePrRelations []entity.IssuePrRelation
	if err := database.DBConn.RawWrapper(sql, option).Find(&issuePrRelations).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find issue_pr_relation: %+v failed", option))
	}
	return &issuePrRelations, nil
}

// func DeleteIssuePrRelation(issuePrRelation *entity.IssuePrRelation) error {
// 	if err := database.DBConn.DB.Delete(issuePrRelation).Error; err != nil {
// 		return errors.Wrap(err, fmt.Sprintf("delete issue_pr_relation: %+v failed", issuePrRelation))
// 	}
// 	return nil
// }

func IssuePrRelationWhere(option *entity.IssuePrRelationOption) string {
	sql := ""

	if option.ID != 0 {
		sql += " and issue_pr_relation.id = @ID"
	}
	if option.IssueID != "" {
		sql += " and issue_pr_relation.issue_id = @IssueID"
	}
	if option.PullRequestID != "" {
		sql += " and issue_pr_relation.pull_request_id = @PullRequestID"
	}
	if option.IssueIDs != nil && len(option.IssueIDs) > 0 {
		sql += " and issue_pr_relation.issue_id in @IssueIDs"
	}

	return sql
}

func IssuePrRelationOrderBy(option *entity.IssuePrRelationOption) string {
	sql := ""

	if option.OrderBy != "" {
		sql += " order by " + option.OrderBy
	}
	if option.Order != "" {
		sql += " " + option.Order
	}

	return sql
}

func IssuePrRelationLimit(option *entity.IssuePrRelationOption) string {
	sql := ""

	if option.Page != 0 && option.PerPage != 0 {
		option.ListOption.CalcOffset()
		sql += " limit @Offset,@PerPage"
	}

	return sql
}

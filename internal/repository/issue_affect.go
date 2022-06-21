package repository

import (
	"fmt"
	"time"

	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func SelectIssueAffect(option *entity.IssueAffectOption) (*[]entity.IssueAffect, error) {
	sql := "select * from issue_affect where 1=1" + IssueAffectWhere(option) + IssueAffectOrderBy(option) + IssueAffectLimit(option)

	// search
	var issueAffects []entity.IssueAffect
	if err := database.DBConn.RawWrapper(sql, option).Find(&issueAffects).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find issue affect: %+v failed", option))
	}
	return &issueAffects, nil
}

func CreateOrUpdateIssueAffect(issueAffect *entity.IssueAffect) error {
	// update
	if issueAffect.CreateTime.IsZero() {
		issueAffect.CreateTime = time.Now()
	}
	if issueAffect.UpdateTime.IsZero() {
		issueAffect.UpdateTime = time.Now()
	}
	if err := database.DBConn.DB.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"update_time", "affect_result"}),
	}).Create(&issueAffect).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create or update issue affect: %+v failed", issueAffect))
	}
	return nil
}

// func CreateIssueAffect(issueAffect *entity.IssueAffect) error {
// 	issueAffect.CreateTime = time.Now()
// 	issueAffect.UpdateTime = time.Now()
// 	// 存储
// 	if err := database.DBConn.DB.Clauses(
// 		clause.OnConflict{DoNothing: true}).Create(&issueAffect).Error; err != nil {
// 		return errors.Wrap(err, fmt.Sprintf("create issue affect: %+v failed", issueAffect))
// 	}
// 	return nil
// }

func IssueAffectWhere(option *entity.IssueAffectOption) string {
	sql := ""

	if option.ID != 0 {
		sql += " and issue_affect.id = @ID"
	}
	if option.IssueID != "" {
		sql += " and issue_affect.issue_id = @IssueID"
	}
	if option.AffectVersion != "" {
		sql += " and issue_affect.affect_version = @AffectVersion"
	}
	if option.AffectResult != "" {
		sql += " and issue_affect.affect_result = @AffectResult"
	}
	if option.IDs != nil && len(option.IDs) > 0 {
		sql += " and issue_affect.id in @IDs"
	}

	return sql
}

func IssueAffectOrderBy(option *entity.IssueAffectOption) string {
	sql := ""

	if option.OrderBy != "" {
		sql += " order by " + option.OrderBy
	}
	if option.Order != "" {
		sql += " " + option.Order
	}

	return sql
}

func IssueAffectLimit(option *entity.IssueAffectOption) string {
	sql := ""

	if option.Page != 0 && option.PerPage != 0 {
		option.ListOption.CalcOffset()
		sql += " limit @Offset,@PerPage"
	}

	return sql
}

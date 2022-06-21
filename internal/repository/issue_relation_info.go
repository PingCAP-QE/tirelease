package repository

import (
	"fmt"

	"tirelease/commons/database"
	"tirelease/internal/dto"
	"tirelease/internal/entity"

	"github.com/pkg/errors"
)

func SelectIssueRelationInfoByJoin(option *dto.IssueRelationInfoQuery) (*[]dto.IssueRelationInfoByJoin, error) {
	sql := "select issue.issue_id, group_concat(issue_affect.id) as issue_affect_ids from "
	sql += ComposeIssueRelationInfoByJoin(option, true)

	// 查询
	var issueRelationInfoJoin []dto.IssueRelationInfoByJoin
	if err := database.DBConn.RawWrapper(sql, option).Find(&issueRelationInfoJoin).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("select issue_relation by raw by join failed, option: %+v", option))
	}

	return &issueRelationInfoJoin, nil
}

func CountIssueRelationInfoByJoin(option *dto.IssueRelationInfoQuery) (int64, error) {
	sql := "select count(*) from "
	sql += ComposeIssueRelationInfoByJoin(option, false)

	// 查询
	var count int64
	if err := database.DBConn.RawWrapper(sql, option).Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("count issue_relation by raw by join failed, option: %+v", option))
	}
	return count, nil
}

func ComposeIssueRelationInfoByJoin(option *dto.IssueRelationInfoQuery, isLimit bool) string {
	issueAffectOption := &entity.IssueAffectOption{
		AffectVersion: option.AffectVersion,
		AffectResult:  option.AffectResult,
	}
	isAffectFilter := (option.AffectVersion != "" || option.AffectResult != "")

	sql := ""
	sql += " ( "
	sql += "select * from issue where 1=1 " + IssueWhere(&option.IssueOption) + IssueOrderBy(&option.IssueOption)
	if !isAffectFilter && isLimit {
		sql += IssueLimit(&option.IssueOption)
	}
	sql += " ) as issue "
	sql += "left join "
	sql += " ( "
	sql += "select * from issue_affect where 1=1 " + IssueAffectWhere(issueAffectOption)
	sql += " ) as issue_affect "
	sql += "on issue.issue_id = issue_affect.issue_id "
	if isAffectFilter {
		sql += "where issue_affect.issue_id is not null "
	}
	sql += "group by issue.issue_id "
	sql += IssueOrderBy(&option.IssueOption)
	if isAffectFilter && isLimit {
		sql += IssueLimit(&option.IssueOption)
	}

	return sql
}

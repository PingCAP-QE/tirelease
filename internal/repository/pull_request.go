package repository

import (
	"encoding/json"
	"fmt"

	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/google/go-github/v41/github"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func SelectPullRequest(option *entity.PullRequestOption) (*[]entity.PullRequest, error) {
	sql := "select * from pull_request where 1=1" + PullRequestWhere(option) + PullRequestOrderBy(option) + PullRequestLimit(option)

	// 查询
	var prs []entity.PullRequest
	if err := database.DBConn.RawWrapper(sql, option).Find(&prs).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find pull request: %+v failed", option))
	}

	// 加工
	for i := 0; i < len(prs); i++ {
		unSerializePullRequest(&prs[i])
	}
	return &prs, nil
}

func SelectPullRequestUnique(option *entity.PullRequestOption) (*entity.PullRequest, error) {
	// 查询
	prs, err := SelectPullRequest(option)
	if err != nil {
		return nil, err
	}

	// 校验
	if len(*prs) == 0 {
		return nil, errors.New(fmt.Sprintf("pull request not found: %+v", option))
	}
	if len(*prs) > 1 {
		return nil, errors.New(fmt.Sprintf("more than one pull request found: %+v", option))
	}
	return &((*prs)[0]), nil
}

func CreateOrUpdatePullRequest(pullRequest *entity.PullRequest) error {
	// 加工
	serializePullRequest(pullRequest)

	// 存储
	if err := database.DBConn.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Omit("Labels", "Assignees", "RequestedReviewers").Create(&pullRequest).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create or update pull request: %+v failed", pullRequest))
	}
	return nil
}

// func CreatePullRequest(pullRequest *entity.PullRequest) error {
// 	// 加工
// 	serializePullRequest(pullRequest)

// 	// 存储
// 	if err := database.DBConn.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&pullRequest).Error; err != nil {
// 		return errors.Wrap(err, fmt.Sprintf("create pull request: %+v failed", pullRequest))
// 	}
// 	return nil
// }

// func DeletePullRequest(pullRequest *entity.PullRequest) error {
// 	if err := database.DBConn.DB.Delete(pullRequest).Error; err != nil {
// 		return errors.Wrap(err, fmt.Sprintf("delete pull request: %+v failed", pullRequest))
// 	}
// 	return nil
// }

// 序列化和反序列化
func serializePullRequest(pullRequest *entity.PullRequest) {
	if nil != pullRequest.Assignees {
		assigneesString, _ := json.Marshal(pullRequest.Assignees)
		pullRequest.AssigneesString = string(assigneesString)
	}
	if nil != pullRequest.Labels {
		labelsString, _ := json.Marshal(pullRequest.Labels)
		pullRequest.LabelsString = string(labelsString)
	}
	if nil != pullRequest.RequestedReviewers {
		requestedReviewersString, _ := json.Marshal(pullRequest.RequestedReviewers)
		pullRequest.RequestedReviewersString = string(requestedReviewersString)
	}
}

func unSerializePullRequest(pullRequest *entity.PullRequest) {
	if pullRequest.AssigneesString != "" {
		var assignees []github.User
		json.Unmarshal([]byte(pullRequest.AssigneesString), &assignees)
		pullRequest.Assignees = &assignees
	}
	if pullRequest.LabelsString != "" {
		var labels []github.Label
		json.Unmarshal([]byte(pullRequest.LabelsString), &labels)
		pullRequest.Labels = &labels
	}
	if pullRequest.RequestedReviewersString != "" {
		var requestedReviewers []github.User
		json.Unmarshal([]byte(pullRequest.RequestedReviewersString), &requestedReviewers)
		pullRequest.RequestedReviewers = &requestedReviewers
	}
}

func PullRequestWhere(option *entity.PullRequestOption) string {
	sql := ""

	if option.ID != 0 {
		sql += " and pull_request.id = @ID"
	}
	if option.PullRequestID != "" {
		sql += " and pull_request.pull_request_id = @PullRequestID"
	}
	if option.Number != 0 {
		sql += " and pull_request.number = @Number"
	}
	if option.State != "" {
		sql += " and pull_request.state = @State"
	}
	if option.Owner != "" {
		sql += " and pull_request.owner = @Owner"
	}
	if option.Repo != "" {
		sql += " and pull_request.repo = @Repo"
	}
	if option.BaseBranch != "" {
		sql += " and pull_request.base_branch = @BaseBranch"
	}
	if option.SourcePullRequestID != "" {
		sql += " and pull_request.source_pull_request_id = @SourcePullRequestID"
	}
	if option.Merged != nil {
		sql += " and pull_request.merged = @Merged"
	}
	if option.MergeableState != nil {
		sql += " and pull_request.mergeable_state = @MergeableState"
	}
	if option.CherryPickApproved != nil {
		sql += " and pull_request.cherry_pick_approved = @CherryPickApproved"
	}
	if option.AlreadyReviewed != nil {
		sql += " and pull_request.already_reviewed = @AlreadyReviewed"
	}
	if option.PullRequestIDs != nil && len(option.PullRequestIDs) > 0 {
		sql += " and pull_request.pull_request_id in @PullRequestIDs"
	}

	return sql
}

func PullRequestOrderBy(option *entity.PullRequestOption) string {
	sql := ""

	if option.OrderBy != "" {
		sql += " order by " + option.OrderBy
	}
	if option.Order != "" {
		sql += " " + option.Order
	}

	return sql
}

func PullRequestLimit(option *entity.PullRequestOption) string {
	sql := ""

	if option.Page != 0 && option.PerPage != 0 {
		option.ListOption.CalcOffset()
		sql += " limit @Offset,@PerPage"
	}

	return sql
}

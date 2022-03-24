package repository

import (
	"fmt"
	"time"

	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func CreateVersionTriage(versionTriage *entity.VersionTriage) error {
	if versionTriage.CreateTime.IsZero() {
		versionTriage.CreateTime = time.Now()
	}
	if versionTriage.UpdateTime.IsZero() {
		versionTriage.UpdateTime = time.Now()
	}
	// 存储
	if err := database.DBConn.DB.Clauses(
		clause.OnConflict{DoNothing: true}).Create(&versionTriage).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create version triage: %+v failed", versionTriage))
	}
	return nil
}

func SelectVersionTriage(option *entity.VersionTriageOption) (*[]entity.VersionTriage, error) {
	// 查询
	var versionTriages []entity.VersionTriage
	if err := database.DBConn.DB.Where(option).Find(&versionTriages).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find version triage: %+v failed", option))
	}
	return &versionTriages, nil
}

func UpdateVersionTriage(versionTriage *entity.VersionTriage) error {
	// 更新
	if versionTriage.UpdateTime.IsZero() {
		versionTriage.UpdateTime = time.Now()
	}
	if err := database.DBConn.DB.Omit("CreateTime").Save(&versionTriage).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("update version triage: %+v failed", versionTriage))
	}
	return nil
}

func CreateOrUpdateVersionTriage(versionTriage *entity.VersionTriage) error {
	// 存储
	if versionTriage.CreateTime.IsZero() {
		versionTriage.CreateTime = time.Now()
	}
	if versionTriage.UpdateTime.IsZero() {
		versionTriage.UpdateTime = time.Now()
	}
	if err := database.DBConn.DB.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"update_time", "triage_owner", "triage_result", "block_version_release", "due_time", "comment"}),
	}).Create(&versionTriage).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create or update version triage: %+v failed", versionTriage))
	}
	return nil
}

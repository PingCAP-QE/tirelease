package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/pkg/errors"
)

func CreateReleaseVersion(version *entity.ReleaseVersion) error {
	// 加工
	if version.CreateTime.IsZero() {
		version.CreateTime = time.Now()
	}
	if version.UpdateTime.IsZero() {
		version.UpdateTime = time.Now()
	}
	serializeReleaseVersion(version)

	// 存储
	if err := database.DBConn.DB.Omit("Repos", "Labels").Create(&version).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create release version: %+v failed", version))
	}
	return nil
}

func UpdateReleaseVersion(version *entity.ReleaseVersion) error {
	// 加工
	if version.UpdateTime.IsZero() {
		version.UpdateTime = time.Now()
	}
	serializeReleaseVersion(version)

	// 更新
	if err := database.DBConn.DB.Omit("CreateTime", "Repos", "Labels").Save(&version).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("update release version: %+v failed", version))
	}
	return nil
}

func SelectReleaseVersion(option *entity.ReleaseVersionOption) (*[]entity.ReleaseVersion, error) {
	// 查询
	var releaseVersions []entity.ReleaseVersion
	if err := database.DBConn.DB.Where(option).Order("major desc, minor desc, patch desc, create_time desc").Find(&releaseVersions).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("find release version: %+v failed", option))
	}

	// 加工
	for i := 0; i < len(releaseVersions); i++ {
		unSerializeReleaseVersion(&releaseVersions[i])
	}
	return &releaseVersions, nil
}

func SelectReleaseVersionLatest(option *entity.ReleaseVersionOption) (*entity.ReleaseVersion, error) {
	// 查询
	releaseVersions, err := SelectReleaseVersion(option)

	// 校验
	if err != nil {
		return nil, err
	}
	length := len(*releaseVersions)
	if length == 0 {
		return nil, errors.New(fmt.Sprintf("find release version unique is nil: %+v failed", option))
	}
	return &((*releaseVersions)[0]), nil
}

// 序列化和反序列化
func serializeReleaseVersion(version *entity.ReleaseVersion) {
	if nil != version.Repos {
		reposString, _ := json.Marshal(version.Repos)
		version.ReposString = string(reposString)
	}
	if nil != version.Labels {
		labelsString, _ := json.Marshal(version.Labels)
		version.LabelsString = string(labelsString)
	}
}

func unSerializeReleaseVersion(version *entity.ReleaseVersion) {
	if version.ReposString != "" {
		var repos []string
		json.Unmarshal([]byte(version.ReposString), &repos)
		version.Repos = &repos
	}
	if version.LabelsString != "" {
		var labels []string
		json.Unmarshal([]byte(version.LabelsString), &labels)
		version.Labels = &labels
	}
}

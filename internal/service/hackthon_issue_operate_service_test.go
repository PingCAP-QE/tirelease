package service

import (
	"testing"
	"time"

	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestIssueAffectOperate(t *testing.T) {
	// Init
	var config = generateConfig()
	database.Connect(config)

	// Init-Data
	var issueAffect = &entity.IssueAffect{
		CreateTime: time.Now(),
		UpdateTime: time.Now(),

		AffectVersion: "5.5.2",
		IssueID:       "100",
		AffectResult:  entity.AffectResultResultUnKnown,
	}
	err := repository.CreateOrUpdateIssueAffect(issueAffect)
	assert.Equal(t, true, err == nil)

	// Update
	var updateOption = &entity.IssueAffectUpdateOption{
		IssueID:       "100",
		AffectVersion: "5.5.2",
		AffectResult:  entity.AffectResultResultYes,
	}
	err = IssueAffectOperate(updateOption)
	assert.Equal(t, true, err == nil)
}

func TestIssueAffectOperateWeb(t *testing.T) {
	// Init
	var config = generateConfig()
	database.Connect(config)

	// Update
	var updateOption = &entity.IssueAffectUpdateOption{
		IssueID:       "I_kwDOAoCpQc5BYBWZ",
		AffectVersion: "5.3",
		AffectResult:  entity.AffectResultResultYes,
	}
	err := IssueAffectOperate(updateOption)
	assert.Equal(t, true, err == nil)
}

func generateConfig() *configs.ConfigYaml {
	var config = &configs.ConfigYaml{}

	config.Mysql.UserName = "cicd_online"
	config.Mysql.PassWord = "wGEXq8a4MeCw6G"
	config.Mysql.Host = "172.16.4.36"
	config.Mysql.Port = "3306"
	config.Mysql.DataBase = "cicd_online"
	config.Mysql.CharSet = "utf8"
	config.Mysql.TimeZone = "Asia%2FShanghai"

	return config
}

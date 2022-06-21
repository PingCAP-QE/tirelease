package repository

import (
	"testing"

	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestSelectIssueRaw(t *testing.T) {
	// Init
	var config = generateConfig()
	database.Connect(config)

	// Select
	var option = &entity.IssueOption{
		IssueIDs: []string{git.TestIssueNodeID, git.TestIssueNodeID2},

		ListOption: entity.ListOption{
			Page:    1,
			PerPage: 10,

			OrderBy: "id",
			Order:   "desc",
		},
	}
	issues, err := SelectIssue(option)
	count, _ := CountIssue(option)
	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(*issues) > 0)
	assert.Equal(t, true, count > 0)
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

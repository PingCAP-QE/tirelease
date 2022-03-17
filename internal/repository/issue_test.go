package repository

import (
	"testing"
	"time"

	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/google/go-github/v41/github"
	"github.com/stretchr/testify/assert"
)

func TestIssue(t *testing.T) {
	// Init
	var config = generateConfig()
	database.Connect(config)

	// Create
	assignee := &github.User{Login: github.String("jcye")}
	var assignees = &([]github.User{*assignee})
	var issue = &entity.Issue{
		IssueID:   "100",
		Number:    100,
		State:     "open",
		Title:     "first",
		Repo:      "ff",
		HTMLURL:   "json",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Assignees: assignees,
	}
	err := CreateOrUpdateIssue(issue)
	// Assert
	assert.Equal(t, true, err == nil)

	// Select
	var option = &entity.IssueOption{
		IssueID: "100",
	}
	issues, err := SelectIssue(option)
	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(*issues) > 0)
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

package service

import (
	"testing"
	"tirelease/commons/configs"
	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestSelectIssues(t *testing.T) {
	configs.LoadConfig("../../config.yaml")
	config := configs.Config
	database.Connect(config)

	issueOption := &entity.IssueOption{
		IssueIDs: []string{"MDU6SXNzdWU1NDU2Mzg0Njk=", "I_kwDOAuklds47leZS"},
	}
	issues, err := SelectIssues(issueOption)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(*issues))
}

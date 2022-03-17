package service

import (
	"testing"
	"time"

	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/dto"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdateVersionTriageInfo(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())

	versionTriage := &entity.VersionTriage{
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
		VersionName:  "5.4.0",
		IssueID:      git.TestIssueNodeID,
		TriageResult: entity.VersionTriageResultAccept,
	}
	info, err := CreateOrUpdateVersionTriageInfo(versionTriage)

	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, info != nil)
	assert.Equal(t, true, info.IsAccept)
}

func TestSelectVersionTriageInfo(t *testing.T) {
	database.Connect(generateConfig())

	query := &dto.VersionTriageInfoQuery{
		Version: "5.4.1",
	}

	info, err := SelectVersionTriageInfo(query)
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, info != nil)
}

package service

import (
	"testing"

	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/dto"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdateVersionTriageInfo(t *testing.T) {
	t.Skip()
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())

	versionTriage := &entity.VersionTriage{
		VersionName:  "6.0",
		IssueID:      git.TestIssueNodeID2,
		TriageResult: entity.VersionTriageResultUnKnown,
	}
	info, err := CreateOrUpdateVersionTriageInfo(versionTriage)

	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, info != nil)
	assert.Equal(t, true, info.IsAccept)
}

func TestSelectVersionTriageInfo(t *testing.T) {
	database.Connect(generateConfig())

	query := &dto.VersionTriageInfoQuery{
		VersionTriageOption: entity.VersionTriageOption{
			VersionName: "5.2.2",
		},
		Version: "5.2.2",
	}
	info, response, err := SelectVersionTriageInfo(query)
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, response.TotalCount > 0)
	assert.Equal(t, true, info != nil)
}

func TestComposeVersionTriageUpcomingList(t *testing.T) {
	t.Skip()
	database.Connect(generateConfig())

	versionTriages, err := ComposeVersionTriageUpcomingList("5.0.7")
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(versionTriages) > 0)
}

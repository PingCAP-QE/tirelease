package repository

import (
	"testing"

	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestVersionTriage(t *testing.T) {
	t.Skip()
	// Init
	var config = generateConfig()
	database.Connect(config)

	// Create
	var versionTriage = &entity.VersionTriage{
		VersionName:  "5.4.1",
		IssueID:      "100",
		TriageResult: entity.VersionTriageResultUnKnown,
		Comment:      "Patch版本5.4.1",
	}
	err := CreateVersionTriage(versionTriage)
	// Assert
	assert.Equal(t, true, err == nil)

	// Update
	versionTriage.TriageResult = entity.VersionTriageResultReleased
	err = UpdateVersionTriage(versionTriage)
	// Assert
	assert.Equal(t, true, err == nil)

	// Select
	var option = &entity.VersionTriageOption{
		VersionName: "5.4.1",
	}
	versions, err := SelectVersionTriage(option)
	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(*versions) > 0)
}

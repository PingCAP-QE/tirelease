package repository

import (
	"testing"

	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestIssueAffect(t *testing.T) {
	t.Skip()
	// Init
	var config = generateConfig()
	database.Connect(config)

	// Create
	var issueAffect = &entity.IssueAffect{
		AffectVersion: "5.4.1",
		IssueID:       "100",
		AffectResult:  entity.AffectResultResultUnKnown,
	}
	issueAffect.AffectResult = entity.AffectResultResultYes
	err := CreateOrUpdateIssueAffect(issueAffect)
	// Assert
	assert.Equal(t, true, err == nil)

	// Select
	var option = &entity.IssueAffectOption{
		AffectVersion: "5.4.1",
	}
	issueAffects, err := SelectIssueAffect(option)
	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(*issueAffects) > 0)
}

package repository

import (
	"testing"

	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/dto"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestSelectIssueRelationInfoByJoin(t *testing.T) {
	// Init
	var config = generateConfig()
	database.Connect(config)

	// Option
	option := &dto.IssueRelationInfoQuery{
		IssueOption: entity.IssueOption{
			IssueID: git.TestIssueNodeID,
		},
		AffectVersion: "5.4",
	}
	issueRelationInfoJoin, err := SelectIssueRelationInfoByJoin(option)

	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(*issueRelationInfoJoin) > 0)

}

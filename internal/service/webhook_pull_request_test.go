package service

import (
	"testing"

	"tirelease/commons/database"
	"tirelease/commons/git"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestCronRefreshPullRequestV4(t *testing.T) {
	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())
	repo := &entity.Repo{
		Owner: git.TestOwner2,
		Repo:  git.TestRepo2,
	}
	repos := []entity.Repo{*repo}
	params := &RefreshPullRequestParams{
		Repos:       &repos,
		BeforeHours: -4380,
		Batch:       20,
		Total:       3000,
	}

	// detail
	err := CronRefreshPullRequestV4(params)
	assert.Equal(t, true, err == nil)
}

func TestWebhookRefreshPullRequestV3(t *testing.T) {
	// init
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)
	database.Connect(generateConfig())

	// detail
	pr, _, err := git.Client.GetPullRequestByNumber(git.TestOwner2, git.TestRepo2, git.TestPullRequestId2)
	assert.Equal(t, true, err == nil)
	err = WebhookRefreshPullRequestV3(pr)
	assert.Equal(t, true, err == nil)
}

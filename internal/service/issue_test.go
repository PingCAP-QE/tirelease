package service

import (
	"testing"
	"time"

	"tirelease/commons/git"
	"tirelease/internal/entity"

	"github.com/google/go-github/v41/github"
	"github.com/stretchr/testify/assert"
)

func TestGetIssueByNumberFromV3(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)

	issue, err := GetIssueByNumberFromV3(git.TestOwner, git.TestRepo, git.TestIssueId)

	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, issue != nil)
}

func TestGetIssuesByTimeFromV3(t *testing.T) {
	git.Connect(git.TestToken)
	git.ConnectV4(git.TestToken)

	day, _ := time.ParseDuration("-24h")
	time := time.Now().Add(15 * day)
	issues, err := GetIssuesByTimeFromV3(git.TestOwner, git.TestRepo, &time)

	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, len(issues) > 0)
}

func TestGetIssuesByOptionV3(t *testing.T) {
	t.Skip()

	git.Connect(git.TestToken)

	newLabel := "affects-6.0"
	labels := []string{"severity/major", "type/bug"}
	option := github.IssueListByRepoOptions{
		State:  "open",
		Labels: labels,
	}
	repos := []entity.Repo{
		{
			Owner: "VelocityLight",
			Repo:  "tirelease",
		},
		{
			Owner: "pingcap",
			Repo:  "tidb",
		},
		{
			Owner: "pingcap",
			Repo:  "tiflash",
		},
		{
			Owner: "pingcap",
			Repo:  "tidb-binlog",
		},
		{
			Owner: "pingcap",
			Repo:  "br",
		},
		{
			Owner: "pingcap",
			Repo:  "tidb-tools",
		},
		{
			Owner: "pingcap",
			Repo:  "ticdc",
		},
		{
			Owner: "pingcap",
			Repo:  "dumpling",
		},
		{
			Owner: "tikv",
			Repo:  "tikv",
		},
		{
			Owner: "tikv",
			Repo:  "pd",
		},
		{
			Owner: "tikv",
			Repo:  "importer",
		},
	}

	for _, repo := range repos {
		issues, err := GetIssuesByOptionV3(repo.Owner, repo.Repo, &option)
		assert.Equal(t, true, err == nil)
		assert.Equal(t, true, len(issues) > 0)
		err = BatchLabelIssues(issues, newLabel)
		assert.Equal(t, true, err == nil)
	}
}

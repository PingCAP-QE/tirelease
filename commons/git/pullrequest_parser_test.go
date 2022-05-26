package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIssueNumber(t *testing.T) {

	issueNumberDatas, err := ParseIssueNumber(shortIssueNumbers, "pingcap", "tidb")
	assert.Equal(t, 1, len(issueNumberDatas))
	assert.Equal(t, "pingcap", issueNumberDatas[0].Owner)
	assert.Equal(t, "tidb", issueNumberDatas[0].Repo)
	assert.Equal(t, 11111, issueNumberDatas[0].Number)
	assert.Equal(t, nil, err)

	issueNumberDatas, err = ParseIssueNumber(fullIssueNumebrs, "", "")
	assert.Equal(t, 1, len(issueNumberDatas))
	assert.Equal(t, "pingcap", issueNumberDatas[0].Owner)
	assert.Equal(t, "tidb", issueNumberDatas[0].Repo)
	assert.Equal(t, 22222, issueNumberDatas[0].Number)
	assert.Equal(t, nil, err)

	issueNumberDatas, err = ParseIssueNumber(linkIssueNumbers, "", "")
	assert.Equal(t, 1, len(issueNumberDatas))
	assert.Equal(t, "pingcap", issueNumberDatas[0].Owner)
	assert.Equal(t, "tidb", issueNumberDatas[0].Repo)
	assert.Equal(t, 3333, issueNumberDatas[0].Number)
	assert.Equal(t, nil, err)

	issueNumberDatas, err = ParseIssueNumber(multipleIssueNumbers, "pingcap", "tidb")
	assert.Equal(t, 3, len(issueNumberDatas))
	assert.Equal(t, "pingcap", issueNumberDatas[0].Owner)
	assert.Equal(t, "tidb", issueNumberDatas[0].Repo)
	assert.Equal(t, 1, issueNumberDatas[0].Number)
	assert.Equal(t, nil, err)

}

const (
	shortIssueNumbers = "xxxxx Issue Number: close #11111"
	fullIssueNumebrs  = "xxxxx Issue Number: close pingcap/tidb#22222"
	linkIssueNumbers  = "xxxxx Issue Number: close https://github.com/pingcap/tidb/issues/3333"

	multipleIssueNumbers = "xxxxx Issue Number: close #1, close pingcap/tidb#2, close https://github.com/pingcap/tidb/issues/3"
)

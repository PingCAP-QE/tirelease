package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	shortIssueNumbers = "xxxxx Issue Number: close #11111"
	fullIssueNumebrs  = "xxxxx Issue Number: close pingcap/tidb#22222"
	linkIssueNumbers  = "xxxxx Issue Number: close https://github.com/pingcap/tidb/issues/3333"

	multipleIssueNumbers = "xxxxx Issue Number: close #1, close pingcap/tidb#2, close https://github.com/pingcap/tidb/issues/3"
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
	releaseNote = "<!--\r\n\r\nThank you for contributing to TiDB!\r\n\r\nPR Title Format:\r\n1. pkg [, pkg2, pkg3]: what's changed\r\n2. *: what's changed\r\n\r\n-->\r\n\r\n### What problem does this PR solve?\r\n<!--\r\n\r\nPlease create an issue first to describe the problem.\r\n\r\nThere MUST be one line starting with \"Issue Number:  \" and \r\nlinking the relevant issues via the \"close\" or \"ref\".\r\n\r\nFor more info, check https://pingcap.github.io/tidb-dev-guide/contribute-to-tidb/contribute-code.html#referring-to-an-issue.\r\n\r\n-->\r\n\r\nIssue Number: close #35880\r\n\r\nProblem Summary:\r\n\r\n### What is changed and how it works?\r\nsee issue description\r\n\r\n### Check List\r\n\r\nTests <!-- At least one of them must be included. -->\r\n\r\n- [x] Unit test\r\n- [ ] Integration test\r\n- [ ] Manual test (add detailed scripts or steps below)\r\n- [ ] No code\r\n\r\nSide effects\r\n\r\n- [ ] Performance regression: Consumes more CPU\r\n- [ ] Performance regression: Consumes more Memory\r\n- [ ] Breaking backward compatibility\r\n\r\nDocumentation\r\n\r\n- [ ] Affects user behaviors\r\n- [ ] Contains syntax changes\r\n- [ ] Contains variable changes\r\n- [ ] Contains experimental features\r\n- [ ] Changes MySQL compatibility\r\n\r\n### Release note\r\n\r\n<!-- compatibility change, improvement, bugfix, and new feature need a release note -->\r\n\r\nPlease refer to [Release Notes Language Style Guide](https://pingcap.github.io/tidb-dev-guide/contribute-to-tidb/release-notes-style-guide.html) to write a quality release note.\r\n\r\n```release-note\r\nfix connect to tidb when using ipv6 host\r\n```\r\n"
)

func TestParseReleaseNote(t *testing.T) {
	releaseNote, err := ParseReleaseNote(releaseNote)
	assert.Equal(t, nil, err)
	assert.Equal(t, "fix connect to tidb when using ipv6 host", releaseNote.ReleaseNote)
}

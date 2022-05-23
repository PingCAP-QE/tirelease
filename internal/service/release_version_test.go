package service

import (
	"testing"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestComposeVersionAtom(t *testing.T) {
	str := "5.4.1"
	major, minor, patch, _ := ComposeVersionAtom(str)
	short := ComposeVersionShortType(str)
	assert.Equal(t, 5, major)
	assert.Equal(t, 4, minor)
	assert.Equal(t, 1, patch)
	assert.Equal(t, entity.ReleaseVersionShortTypePatch, short)

	str = "5.4"
	major, minor, patch, _ = ComposeVersionAtom(str)
	short = ComposeVersionShortType(str)
	assert.Equal(t, 5, major)
	assert.Equal(t, 4, minor)
	assert.Equal(t, 0, patch)
	assert.Equal(t, entity.ReleaseVersionShortTypeMinor, short)

	str = "5.4-hotfix-1"
	major, minor, patch, addition := ComposeVersionAtom(str)
	short = ComposeVersionShortType(str)
	assert.Equal(t, 5, major)
	assert.Equal(t, 4, minor)
	assert.Equal(t, 0, patch)
	assert.Equal(t, "hotfix-tiflash-patch1", addition)
	assert.Equal(t, entity.ReleaseVersionShortTypeMinor, short)

	str = "5.4.1-hotfix"
	major, minor, patch, addition = ComposeVersionAtom(str)
	short = ComposeVersionShortType(str)
	assert.Equal(t, 5, major)
	assert.Equal(t, 4, minor)
	assert.Equal(t, 1, patch)
	assert.Equal(t, "hotfix", addition)
	assert.Equal(t, entity.ReleaseVersionShortTypeMinor, short)
}

// 注意: 该脚本涉及线上变更，按需使用
func TestLabelAffect(t *testing.T) {
	// git.Connect(git.TestToken)
	// git.ConnectV4(git.TestToken)
	// database.Connect(generateConfig())

	// option := &entity.IssueOption{
	// 	State:          git.OpenStatus,
	// 	SeverityLabels: []string{git.SeverityCriticalLabel, git.SeverityMajorLabel},
	// }
	// // 以下版本号随实际版本变更
	// // label := fmt.Sprintf(git.AffectsLabel, ComposeVersionMinorName(releaseVersion))
	// label := fmt.Sprintf(git.AffectsLabel, "6.1")
	// err := RefreshIssueLabel(label, option)
	// if nil != err {
	// 	fmt.Printf("%v", err)
	// }
}

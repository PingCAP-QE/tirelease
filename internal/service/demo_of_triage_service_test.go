package service

import (
	"strings"
	"testing"
	"time"

	"tirelease/commons/git"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	// Connect
	git.Connect("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")

	// List all repositories for the authenticated user
	triageItems, err := CollectTriageItemByRepo("VelocityLight", "tirelease")

	// Assert
	assert.Equal(t, true, len(triageItems) > 0)
	assert.Equal(t, true, err == nil)
}

func TestLabel(t *testing.T) {
	// Connect
	git.Connect("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")

	labels := []string{"jc_test"}
	err := AddLabelOfAccept("VelocityLight", "tirelease", 1, labels)

	// Assert
	assert.Equal(t, true, err == nil)
}

func TestStrings(t *testing.T) {
	s := "release-5.4"
	replace := strings.Replace(s, git.ReleaseBranchPrefix, "", -1)
	assert.Equal(t, "release-5.4", s)
	assert.Equal(t, "5.4", replace)
}

func TestTimeNil(t *testing.T) {
	triageItem := &entity.TriageItem{}
	isZero := triageItem.CreateTime.IsZero()
	assert.Equal(t, true, isZero)

	var time time.Time
	isZero = time.IsZero()
	assert.Equal(t, true, isZero)
}

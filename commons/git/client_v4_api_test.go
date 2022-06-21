package git

import (
	"context"
	"testing"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
)

//=======================================================================Login
func TestGetLoginV4(t *testing.T) {
	// Connect
	ConnectV4(TestToken)

	// Query
	var query struct {
		Viewer struct {
			Login     githubv4.String
			CreatedAt githubv4.DateTime
		}
	}
	err := ClientV4.client.Query(context.Background(), &query, nil)

	// Assert
	assert.Equal(t, true, err == nil)
}

//=======================================================================Issue
func TestGetIssueV4(t *testing.T) {
	// Connect
	ConnectV4(TestToken)

	// Query
	issue, err := ClientV4.GetIssueByNumber(TestOwner, TestRepo, TestIssueId)

	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, issue != nil)
}

func TestGetIssueByIDV4(t *testing.T) {
	// Connect
	ConnectV4(TestToken)

	// Query
	issue, err := ClientV4.GetIssueByID(TestIssueNodeID)

	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, issue != nil)
}

func TestGetIssueWithoutTimelineByIDV4(t *testing.T) {
	// Connect
	ConnectV4(TestToken)

	// Query
	issue, err := ClientV4.GetIssueWithoutTimelineByID(TestIssueNodeID)

	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, issue != nil)
}

//=======================================================================Pr
func TestGetPullRequestV4(t *testing.T) {
	// Connect
	ConnectV4(TestToken)

	// Query
	pr, err := ClientV4.GetPullRequestsByNumber(TestOwner, TestRepo, TestPullRequestId)

	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, pr != nil)
}

func TestGetPullRequestByIDV4(t *testing.T) {
	// Connect
	ConnectV4(TestToken)

	// Query
	pr, err := ClientV4.GetPullRequestByID(TestPullRequestNodeID)

	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, pr != nil)
}

func TestGetPullRequestWithoutTimelineByIDV4(t *testing.T) {
	// Connect
	ConnectV4(TestToken)

	// Query
	pr, err := ClientV4.GetPullRequestWithoutTimelineByID(TestPullRequestNodeID)

	// Assert
	assert.Equal(t, true, err == nil)
	assert.Equal(t, true, pr != nil)
}

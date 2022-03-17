package controller

import (
	"net/http"

	"tirelease/commons/git"
	"tirelease/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v41/github"
)

func GithubWebhookHandler(c *gin.Context) {
	// parse webhook payload
	payload, err := github.ValidatePayload(c.Request, nil)
	if err != nil {
		c.Error(err)
		return
	}
	event, err := github.ParseWebHook(github.WebHookType(c.Request), payload)
	if err != nil {
		c.Error(err)
		return
	}

	// handle event
	switch event := event.(type) {

	case *github.IssuesEvent:
		err := service.WebhookRefreshIssueV4(event.Issue)
		if err != nil {
			c.Error(err)
			return
		}

	case *github.IssueCommentEvent:
		url := event.Issue.HTMLURL
		if git.IsIssue(*url) {
			err := service.WebhookRefreshIssueV4(event.Issue)
			if err != nil {
				c.Error(err)
				return
			}
		}

	case *github.PullRequestEvent:
		err := service.WebhookRefreshPullRequestV3(event.PullRequest)
		if err != nil {
			c.Error(err)
			return
		}

	default:

	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

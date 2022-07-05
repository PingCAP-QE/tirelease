package controller

import (
	"net/http"
	"strings"

	"tirelease/commons/git"
	"tirelease/internal/service"
	gconsumer "tirelease/internal/service/git_event_consumer"

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
		nodeID := event.Issue.NodeID
		if git.IsIssue(*url) {
			err := service.WebhookRefreshIssueV4(event.Issue)
			if err != nil {
				c.Error(err)
				return
			}
		}
		if git.IsPullRequest(*url) {
			pr := &github.PullRequest{
				NodeID: nodeID,
			}
			err := service.WebHookRefreshPullRequestRefIssue(pr)
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
		baseBranch := event.PullRequest.Base.Ref
		if baseBranch != nil && strings.HasPrefix(*baseBranch, git.ReleaseBranchPrefix) {
			// If the auto refesh operion failed ,just let it go.
			service.AutoRefreshPrApprovedLabel(event.PullRequest)

			err := service.WebHookRefreshPullRequestRefIssue(event.PullRequest)
			if err != nil {
				c.Error(err)
				return
			}
		}
		// notify all pullrequest consumers
		// TODO: Changed to isolated notifier with other consumers.
		prConsumers := gconsumer.GetPREventConsumers()
		for _, prConsumer := range prConsumers {
			if prConsumer.Validate(*event) {
				err := prConsumer.Consume(*event)
				if err != nil {
					c.Error(err)
					return
				}
			}
		}

	default:
		c.JSON(http.StatusAccepted, gin.H{"status": "accepted but not supported"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

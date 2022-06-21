// Tool Url: https://github.com/google/go-github
// Tool Guide: https://docs.github.com/en/rest/reference/webhooks#repository-webhook-configuration

package git

import (
	"context"
	"time"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

type GithubInfo struct {
	Client *github.Client
}

// V3版本Rest-API
var Client = &GithubInfo{}

func Connect(accessToken string) {
	// Outh
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	tc.Timeout = time.Second * 180

	// Github client
	githubClient := github.NewClient(tc)
	Client.Client = githubClient
}

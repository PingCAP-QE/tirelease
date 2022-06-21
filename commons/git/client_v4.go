// Tool Url: https://github.com/shurcooL/githubv4
// Tool Guide: https://docs.github.com/en/graphql
// Web API Explorer: https://docs.github.com/en/graphql/overview/explorer

package git

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GithubInfoV4 struct {
	client *githubv4.Client
}

// V4版本GraphAPI
var ClientV4 = &GithubInfoV4{}

func ConnectV4(accessToken string) {
	// Outh
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	httpClient := oauth2.NewClient(context.Background(), src)
	httpClient.Timeout = time.Second * 180

	// Client
	ClientV4.client = githubv4.NewClient(httpClient)
}

package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

// Operation
func AddLabelByPullRequestID(pullRequestID, label string) error {
	// select issue by id
	option := &entity.PullRequestOption{
		PullRequestID: pullRequestID,
	}
	pr, err := repository.SelectPullRequestUnique(option)
	if nil != err {
		return err
	}

	// add issue label
	_, _, err = git.Client.AddLabel(pr.Owner, pr.Repo, pr.Number, label)
	if nil != err {
		return err
	}
	return nil
}

func RemoveLabelByPullRequestID(pullRequestID, label string) error {
	// select issue by id
	option := &entity.PullRequestOption{
		PullRequestID: pullRequestID,
	}
	pr, err := repository.SelectPullRequestUnique(option)
	if nil != err {
		return err
	}

	// remove issue label
	_, err = git.Client.RemoveLabel(pr.Owner, pr.Repo, pr.Number, label)
	if nil != err {
		return err
	}
	return nil
}

// Query PullRequest From Github And Construct Issue Data Service
func GetPullRequestByNumberFromV3(owner, repo string, number int) (*entity.PullRequest, error) {
	pr, _, err := git.Client.GetPullRequestByNumber(owner, repo, number)
	if nil != err {
		return nil, err
	}
	return entity.ComposePullRequestFromV3(pr), nil
}

func GetPullRequestRefIssuesByRegexFromV4(pr *git.PullRequestField) ([]int, error) {
	issueNumbers := make([]int, 0)
	if pr == nil {
		return issueNumbers, nil
	}

	// from body reference
	if pr.Body != "" {
		bodyIssues, err := RegexReferenceNumbers(string(pr.Body))
		if err != nil {
			return nil, err
		}
		issueNumbers = append(issueNumbers, bodyIssues...)
	}

	// from timeline reference
	edges := pr.TimelineItems.Edges
	if nil != edges && len(edges) > 0 {
		for i := range edges {
			edge := edges[i]
			if nil == &edge.Node || nil == &edge.Node.IssueComment ||
				nil == &edge.Node.IssueComment.Body {
				continue
			}
			if git.IssueComment != edge.Node.Typename {
				continue
			}
			issueComment := string(edge.Node.IssueComment.Body)
			issueCommentNumbers, err := RegexReferenceNumbers(issueComment)
			if err != nil {
				return nil, err
			}
			issueNumbers = append(issueNumbers, issueCommentNumbers...)
		}
	}

	return issueNumbers, nil
}

func RegexReferenceNumbers(text string) ([]int, error) {
	// param protect
	issueNumbers := make([]int, 0)
	if text == "" {
		return issueNumbers, nil
	}

	// issue number regex
	issueStrs := make([]string, 0)
	re := regexp.MustCompile(`[#][0-9]+`)
	for _, match := range re.FindAllString(text, -1) {
		re2 := regexp.MustCompile(`[0-9]+`)
		issueStrs = append(issueStrs, re2.FindAllString(match, -1)...)
	}

	// compose issue number list
	if len(issueStrs) > 0 {
		for _, issueStr := range issueStrs {
			issueNumber, err := strconv.Atoi(issueStr)
			if nil != err {
				return nil, err
			}
			issueNumbers = append(issueNumbers, issueNumber)
		}
	}

	return issueNumbers, nil
}

// Get Pullrequest related release versions by the baseBranch
func getPrRelatedReleaseVersions(pr entity.PullRequest) ([]entity.ReleaseVersion, error) {
	// Disclose the whole inside story
	// If the base branch of pr is not start with "release-" then it will not be triggered
	if !strings.HasPrefix(pr.BaseBranch, git.ReleaseBranchPrefix) {
		return nil, fmt.Errorf("pr %s is not checked out from a released branch", pr.PullRequestID)
	}

	branchVersion := strings.Replace(pr.BaseBranch, git.ReleaseBranchPrefix, "", -1)
	major, minor, _, _ := ComposeVersionAtom(branchVersion)

	// select all triaged list under this minor version
	versionOption := &entity.ReleaseVersionOption{
		Major:     major,
		Minor:     minor,
		ShortType: entity.ReleaseVersionShortTypeMinor,
	}
	releaseVersions, err := repository.SelectReleaseVersion(versionOption)
	if err != nil {
		return nil, err
	}
	return *releaseVersions, err
}

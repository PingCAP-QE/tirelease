package git

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	issueNumberBlockRegexpTemplate = "(?i)%s?\\s*%s(?P<issue_number>[1-9]\\d*)"

	associatePrefixRegexp           = "(?P<associate_prefix>ref|close[sd]?|resolve[sd]?|fix(e[sd])?)"
	orgRegexp                       = "[a-zA-Z0-9][a-zA-Z0-9-]{0,38}"
	repoRegexp                      = "[a-zA-Z0-9-_]{1,100}"
	issueNumberPrefixRegexpTemplate = "(?P<issue_number_prefix>(https|http)://github\\.com/%s/%s/issues/|%s/%s#|#)"
	linkPrefixRegexpTemplate        = "(https|http)://github\\.com/(?P<org>%s)/(?P<repo>%s)/issues/"
	fullPrefixRegexpTemplate        = "(?P<org>%s)/(?P<repo>%s)#"
	shortPrefix                     = "#"

	associatePrefixGroupName   = "associate_prefix"
	issueNumberPrefixGroupName = "issue_number_prefix"
	issueNumberGroupName       = "issue_number"
	orgGroupName               = "org"
	repoGroupName              = "repo"
)

type IssueNumberData struct {
	Owner  string
	Repo   string
	Number int
}

var (
	issueNumberPrefixRegexp = fmt.Sprintf(issueNumberPrefixRegexpTemplate, orgRegexp, repoRegexp, orgRegexp, repoRegexp)
	linkPrefixRegexp        = fmt.Sprintf(linkPrefixRegexpTemplate, orgRegexp, repoRegexp)
	fullPrefixRegexp        = fmt.Sprintf(fullPrefixRegexpTemplate, orgRegexp, repoRegexp)

	pingcapIssueRefBlockRegexp = "(?im)Issue Number:\\s*((,\\s*)?(ref|close[sd]?|resolve[sd]?|fix(e[sd])?)\\s*((https|http)://github\\.com/" + orgRegexp + "/" + repoRegexp + "/issues/|" + orgRegexp + "/" + repoRegexp + "#|#)(?P<issue_number>[1-9]\\d*))+"
)

func ParseIssueNumber(content, currOwner, currRepo string) ([]IssueNumberData, error) {
	issueNumberBlockRegexp := fmt.Sprintf(issueNumberBlockRegexpTemplate, associatePrefixRegexp, issueNumberPrefixRegexp)
	pingcapIssueRefBlock, err := composeIssueBlocks(content, pingcapIssueRefBlockRegexp)

	if err != nil {
		return nil, err
	}

	compile, err := regexp.Compile(issueNumberBlockRegexp)
	if err != nil {
		return nil, err
	}

	allMatches := compile.FindAllStringSubmatch(pingcapIssueRefBlock, -1)
	groupNames := compile.SubexpNames()

	issueNumberDatas := make([]IssueNumberData, 0)

	for _, matches := range allMatches {
		issueNumberPrefix := ""
		issueNumber := 0
		for i, groupName := range groupNames {
			switch groupName {
			case issueNumberPrefixGroupName:
				issueNumberPrefix = strings.ToLower(strings.TrimSpace(matches[i]))
			case issueNumberGroupName:
				issueNumber, err = strconv.Atoi(strings.TrimSpace(matches[i]))
				if err != nil {
					return nil, err
				}
			}
		}

		if b, org, repo, err := isLinkPrefix(issueNumberPrefix); b && err == nil {
			issueNumberDatas = append(issueNumberDatas, IssueNumberData{Owner: org, Repo: repo, Number: issueNumber})
		} else if b, org, repo, err := isFullPrefix(issueNumberPrefix); b && err == nil {
			issueNumberDatas = append(issueNumberDatas, IssueNumberData{Owner: org, Repo: repo, Number: issueNumber})
		} else if isShortPrefix(issueNumberPrefix) {
			issueNumberDatas = append(issueNumberDatas, IssueNumberData{Owner: currOwner, Repo: currRepo, Number: issueNumber})
		}

		if err != nil {
			return nil, err
		}

	}

	return issueNumberDatas, nil

}

func composeIssueBlocks(content, issueRegexp string) (string, error) {

	compile, err := regexp.Compile(issueRegexp)

	if err != nil {
		return "", err
	}

	allMatches := compile.FindAllString(content, -1)

	matchedStrings := make([]string, 0)

	for _, matche := range allMatches {
		matchedStrings = append(matchedStrings, strings.ToLower(strings.TrimSpace(matche)))
	}

	return strings.Join(matchedStrings, " "), nil
}

// isLinkPrefix used to determine whether the prefix style of the issue number is link prefix,
// for example: https://github/com/pingcap/tidb/issues/123.
func isLinkPrefix(prefix string) (bool, string, string, error) {
	compile, err := regexp.Compile(linkPrefixRegexp)
	if err != nil {
		return false, "", "", nil
	}

	matches := compile.FindStringSubmatch(prefix)
	groupNames := compile.SubexpNames()

	if matches == nil {
		return false, "", "", nil
	}

	org := ""
	repo := ""
	for i, match := range matches {
		groupName := groupNames[i]
		if groupName == orgGroupName {
			org = match
		} else if groupName == repoGroupName {
			repo = match
		}
	}

	return true, org, repo, nil
}

// isFullPrefix used to determine whether the prefix style of the issue number is full prefix,
// for example: pingcap/tidb#123.
func isFullPrefix(prefix string) (bool, string, string, error) {
	compile, err := regexp.Compile(fullPrefixRegexp)
	if err != nil {
		return false, "", "", err
	}

	matches := compile.FindStringSubmatch(prefix)
	groupNames := compile.SubexpNames()

	if matches == nil {
		return false, "", "", nil
	}

	org := ""
	repo := ""
	for i, match := range matches {
		groupName := groupNames[i]
		if groupName == orgGroupName {
			org = match
		} else if groupName == repoGroupName {
			repo = match
		}
	}

	return true, org, repo, nil
}

// isShortPrefix used to determine whether the prefix style of the issue number is short prefix,
// for example: #123.
func isShortPrefix(prefix string) bool {
	return prefix == shortPrefix
}

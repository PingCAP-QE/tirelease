package service

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/google/go-github/v41/github"
)

func SelectIssues(option *entity.IssueOption) (*[]entity.Issue, error) {
	issues, err := repository.SelectIssue(option)
	if err != nil {
		return nil, err
	}

	githubUsers := getAllGithubUsers(issues)
	githubLogins := extractLoginsFromGitUsers(githubUsers)

	employees, err := repository.BatchSelectEmployeesByGhLogins(githubLogins)
	if err != nil {
		return nil, err
	}
	loginEmployeeMap := composeLoginEmployeeMap(employees)

	for i := range *issues {
		issue := &(*issues)[i]
		assignedEmployees := composeAssignedEmployees(*issue, loginEmployeeMap)
		issue.AssignedEmployees = &assignedEmployees
	}

	return issues, nil
}

func getAllGithubUsers(issues *[]entity.Issue) (githubUsers []github.User) {
	for _, issue := range *issues {
		githubUsers = append(githubUsers, *issue.Assignees...)
	}
	return
}

func extractLoginsFromGitUsers(githubUsers []github.User) []string {
	var githubLogins []string
	for _, githubUser := range githubUsers {
		if login := githubUser.GetLogin(); login != "" {
			githubLogins = append(githubLogins, login)
		}
	}
	return githubLogins
}

func composeLoginEmployeeMap(employees []entity.Employee) map[string]*entity.Employee {
	var loginEmployeeMap = make(map[string]*entity.Employee)
	for i := range employees {
		employee := &employees[i]
		loginEmployeeMap[employee.GithubId] = employee
	}
	return loginEmployeeMap
}

func composeAssignedEmployees(issue entity.Issue, loginEmployeeMap map[string]*entity.Employee) []entity.Employee {
	assignedEmployees := make([]entity.Employee, 0)
	assignees := issue.Assignees
	for _, assignee := range *assignees {
		employee := loginEmployeeMap[assignee.GetLogin()]
		if employee != nil {
			assignedEmployees = append(assignedEmployees,
				entity.Employee{
					ID:       employee.ID,
					Name:     employee.Name,
					Email:    employee.Email,
					GithubId: employee.GithubId,
					GhName:   employee.GhName,
				})
		} else {
			assignedEmployees = append(assignedEmployees,
				entity.Employee{
					GithubId: assignee.GetLogin(),
				},
			)
		}
	}

	return assignedEmployees
}

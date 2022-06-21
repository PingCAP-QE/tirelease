package cron

import (
	"tirelease/commons/cron"
	"tirelease/internal/entity"
	"tirelease/internal/repository"
	"tirelease/internal/service"
)

func IssueCron() {
	// Cron 表达式及功能方法
	repos, err := repository.SelectRepo(&entity.RepoOption{})
	if err != nil {
		return
	}
	releaseVersions, err := repository.SelectReleaseVersion(&entity.ReleaseVersionOption{})
	if err != nil {
		return
	}
	params := &service.RefreshIssueParams{
		Repos:           repos,
		BeforeHours:     -2,
		Batch:           20,
		Total:           500,
		IsHistory:       true,
		ReleaseVersions: releaseVersions,
		Order:           "DESC",
	}
	cron.Create("0 0 */1 * * ?", func() { service.CronRefreshIssuesV4(params) })
}

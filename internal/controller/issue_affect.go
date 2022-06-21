package controller

import (
	"net/http"

	"tirelease/internal/entity"
	"tirelease/internal/repository"
	"tirelease/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func SelectIssueAffect(c *gin.Context) {
	// Params
	option := entity.IssueAffectOption{}
	if err := c.ShouldBindWith(&option, binding.JSON); err != nil {
		c.Error(err)
		return
	}

	// Action
	issueAffects, err := repository.SelectIssueAffect(&option)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": issueAffects})
}

func SelectIssueAffectResult(c *gin.Context) {
	var enumResult = struct {
		AffectResultResultUnKnown entity.AffectResultResult
		AffectResultResultYes     entity.AffectResultResult
		AffectResultResultNo      entity.AffectResultResult
	}{
		AffectResultResultUnKnown: entity.AffectResultResultUnKnown,
		AffectResultResultYes:     entity.AffectResultResultYes,
		AffectResultResultNo:      entity.AffectResultResultNo,
	}

	c.JSON(http.StatusOK, gin.H{"data": enumResult})
}

func CreateOrUpdateIssueAffect(c *gin.Context) {
	// Params
	issueAffect := entity.IssueAffect{}
	if err := c.ShouldBindWith(&issueAffect, binding.JSON); err != nil {
		c.Error(err)
		return
	}

	// Action
	err := service.CreateOrUpdateIssueAffect(&issueAffect)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

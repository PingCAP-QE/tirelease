package controller

import (
	"net/http"

	"tirelease/internal/dto"
	"tirelease/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func SelectIssueRelationInfos(c *gin.Context) {
	// Params
	option := dto.IssueRelationInfoQuery{}

	if err := c.ShouldBindWith(&option, binding.Form); err != nil {
		c.Error(err)
		return
	}
	if option.Page == 0 {
		option.Page = 1
	}
	if option.PerPage == 0 {
		option.PerPage = 10
	}
	option.ParamFill()

	// Action
	issueRelationInfos, response, err := service.SelectIssueRelationInfo(&option)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": issueRelationInfos, "response": response})
}

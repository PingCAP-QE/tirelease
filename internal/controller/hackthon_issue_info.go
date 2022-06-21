package controller

import (
	"log"
	"net/http"

	"tirelease/internal/service"

	"github.com/gin-gonic/gin"
)

func ListIssueInfo(c *gin.Context) {
	state := c.DefaultQuery("state", "CLOSED")
	issueInfos, err := service.ListIssueInfo(state)
	if err != nil {
		log.Fatal(err)
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, issueInfos)
}

func FilterIssueInfo(c *gin.Context) {
	version := c.DefaultQuery("version", "5.4")
	issueInfos, err := service.FilterIssueInfo(version)
	if err != nil {
		log.Fatal(err)
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, issueInfos)
}

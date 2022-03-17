package controller

import (
	"log"
	"net/http"

	"tirelease/internal/entity"
	"tirelease/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func UpdateIssueAffectAndTriage(c *gin.Context) {
	option := entity.IssueAffectUpdateOption{}
	if err := c.ShouldBindWith(&option, binding.JSON); err != nil {
		log.Fatal(err)
		c.Error(err)
		return
	}
	if err := service.IssueAffectOperate(&option); err != nil {
		log.Fatal(err)
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

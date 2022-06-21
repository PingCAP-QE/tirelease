package controller

import (
	"net/http"

	"tirelease/internal/entity"
	"tirelease/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func UpdateVersionTriage(c *gin.Context) {
	// Params
	versionTriage := entity.VersionTriage{}
	if err := c.ShouldBindWith(&versionTriage, binding.JSON); err != nil {
		c.Error(err)
		return
	}

	// Action
	err := service.UpdateVersionTriage(&versionTriage)
	if nil != err {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

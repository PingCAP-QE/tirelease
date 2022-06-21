package controller

import (
	"net/http"

	"tirelease/internal/dto"
	"tirelease/internal/entity"
	"tirelease/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreateOrUpdateVersionTriage(c *gin.Context) {
	// Params
	versionTriageModifyOption := entity.VersionTriageModifyOption{}
	if err := c.ShouldBindWith(&versionTriageModifyOption, binding.JSON); err != nil {
		c.Error(err)
		return
	}

	// Action
	versionTriageInfo, err := service.CreateOrUpdateVersionTriageInfo(&versionTriageModifyOption.VersionTriage, versionTriageModifyOption.UpdatedVars...)
	if nil != err {
		c.Error(err)
		return
	}

	// Response
	var statusCode int = 200
	if nil != versionTriageInfo && versionTriageInfo.IsFrozen && versionTriageInfo.IsAccept {
		statusCode = 202
	}
	c.JSON(statusCode, gin.H{"status": "ok", "data": versionTriageInfo})
}

func SelectVersionTriageInfo(c *gin.Context) {
	// Params
	versionTriageInfoQuery := dto.VersionTriageInfoQuery{}
	if err := c.ShouldBindUri(&versionTriageInfoQuery); err != nil {
		c.Error(err)
		return
	}
	if versionTriageInfoQuery.Version != "" {
		versionTriageInfoQuery.VersionName = versionTriageInfoQuery.Version
	}

	// Action
	versionTriageInfos, response, err := service.SelectVersionTriageInfo(&versionTriageInfoQuery)
	if nil != err {
		c.Error(err)
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{"data": versionTriageInfos, "response": response})
}

func SelectVersionTriageResult(c *gin.Context) {
	var enumResult = struct {
		VersionTriageResultUnKnown      entity.VersionTriageResult
		VersionTriageResultAccept       entity.VersionTriageResult
		VersionTriageResultAcceptFrozen entity.VersionTriageResult
		VersionTriageResultLater        entity.VersionTriageResult
		VersionTriageResultWontFix      entity.VersionTriageResult
		VersionTriageResultReleased     entity.VersionTriageResult
	}{
		VersionTriageResultUnKnown:      entity.VersionTriageResultUnKnown,
		VersionTriageResultAccept:       entity.VersionTriageResultAccept,
		VersionTriageResultAcceptFrozen: entity.VersionTriageResultAcceptFrozen,
		VersionTriageResultLater:        entity.VersionTriageResultLater,
		VersionTriageResultWontFix:      entity.VersionTriageResultWontFix,
		VersionTriageResultReleased:     entity.VersionTriageResultReleased,
	}

	c.JSON(http.StatusOK, gin.H{"data": enumResult})
}

package controller

import (
	"net/http"

	"tirelease/internal/entity"
	"tirelease/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreateReleaseVersion(c *gin.Context) {
	// Params
	releaseVersion := entity.ReleaseVersion{}
	if err := c.ShouldBindWith(&releaseVersion, binding.JSON); err != nil {
		c.Error(err)
		return
	}

	// Action
	err := service.CreateReleaseVersion(&releaseVersion)
	if nil != err {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func UpdateReleaseVersion(c *gin.Context) {
	// Params
	releaseVersion := entity.ReleaseVersion{}
	if err := c.ShouldBindWith(&releaseVersion, binding.JSON); err != nil {
		c.Error(err)
		return
	}

	// Action
	err := service.UpdateReleaseVersion(&releaseVersion)
	if nil != err {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func SelectReleaseVersion(c *gin.Context) {
	// Params
	option := entity.ReleaseVersionOption{}
	if err := c.ShouldBindWith(&option, binding.Form); err != nil {
		c.Error(err)
		return
	}

	// Action
	releaseVersions, err := service.SelectReleaseVersion(&option)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": releaseVersions})
}

func SelectReleaseVersionMaintained(c *gin.Context) {
	res, err := service.SelectReleaseVersionMaintained()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func SelectReleaseVersionStatus(c *gin.Context) {
	var enumResult = struct {
		ReleaseVersionStatusPlanned   entity.ReleaseVersionStatus
		ReleaseVersionStatusUpcoming  entity.ReleaseVersionStatus
		ReleaseVersionStatusFrozen    entity.ReleaseVersionStatus
		ReleaseVersionStatusReleased  entity.ReleaseVersionStatus
		ReleaseVersionStatusCancelled entity.ReleaseVersionStatus
	}{
		ReleaseVersionStatusPlanned:   entity.ReleaseVersionStatusPlanned,
		ReleaseVersionStatusUpcoming:  entity.ReleaseVersionStatusUpcoming,
		ReleaseVersionStatusFrozen:    entity.ReleaseVersionStatusFrozen,
		ReleaseVersionStatusReleased:  entity.ReleaseVersionStatusReleased,
		ReleaseVersionStatusCancelled: entity.ReleaseVersionStatusCancelled,
	}

	c.JSON(http.StatusOK, gin.H{"data": enumResult})
}

func SelectReleaseVersionType(c *gin.Context) {
	var enumResult = struct {
		ReleaseVersionTypeMajor entity.ReleaseVersionType
		ReleaseVersionTypeMinor entity.ReleaseVersionType
		ReleaseVersionTypePatch entity.ReleaseVersionType
	}{
		ReleaseVersionTypeMajor: entity.ReleaseVersionTypeMajor,
		ReleaseVersionTypeMinor: entity.ReleaseVersionTypeMinor,
		ReleaseVersionTypePatch: entity.ReleaseVersionTypePatch,
	}

	c.JSON(http.StatusOK, gin.H{"data": enumResult})
}

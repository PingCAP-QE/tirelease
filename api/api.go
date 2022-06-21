package api

import (
	"net/http"
	"tirelease/internal/controller"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// Create gin-routers
func Routers(file string) (router *gin.Engine) {
	router = gin.Default()

	// Error
	routeError(router)

	// Cors
	routeCors(router)

	// Root html & Folder static/ for JS/CSS & home/*any for website('url' match to website/src/routes.js)
	routeHtml(router, file)

	// Real REST-API registry
	routeRestAPI(router)

	return router
}

func routeError(router *gin.Engine) {
	router.Use(APIErrorJSONReporter())
}

func routeCors(router *gin.Engine) {
	router.Use(Cors())
}

func routeHtml(router *gin.Engine, file string) {
	router.Use(
		static.Serve("/", static.LocalFile(file, true)),
	)
	router.Use(
		static.Serve("/static", static.LocalFile(file, true)),
	)
	homePages := router.Group("/home")
	{
		homePages.GET("/*any", func(c *gin.Context) {
			c.FileFromFS("/", http.Dir(file))
		})
	}
}

// REST-API Function
func routeRestAPI(router *gin.Engine) {
	// testEntity := router.Group("/testentity")
	// {
	// 	testEntity.GET("/select", controller.TestEntitySelect)
	// 	testEntity.POST("/insert", controller.TestEntityInsert)
	// }

	// triageItem := router.Group("/triage")
	// {
	// 	triageItem.GET("/select", controller.SelectTriageItems)
	// 	triageItem.POST("/insert", controller.InsertTriageItems)
	// 	triageItem.POST("/accept", controller.AddLabelsToIssue)
	// }

	// issue := router.Group("/issue")
	// {
	// 	issue.GET("", controller.ListIssueInfo)
	// 	issue.GET("/init", controller.InitDataForDemo)
	// 	issue.GET("/filter", controller.FilterIssueInfo)
	// 	issue.POST("/affect", controller.UpdateIssueAffectAndTriage)
	// }

	ping := router.Group("/ping")
	{
		ping.GET("", controller.TestPingPongGet)
		ping.GET("/error", controller.TestPingError)
		ping.POST("", controller.TestPingPongPost)
		ping.POST("/:name", controller.TestPingPongPost)
	}

	issue := router.Group("/issue")
	{
		issue.GET("", controller.SelectIssueRelationInfos)

		issue.POST("/:issue_id/cherrypick/:version_name", controller.CreateOrUpdateVersionTriage)
		issue.PATCH("/:issue_id/cherrypick/:version_name", controller.CreateOrUpdateVersionTriage)
		issue.GET("/cherrypick/:version", controller.SelectVersionTriageInfo)
		issue.GET("/cherrypick/result", controller.SelectVersionTriageResult)

		issue.PATCH("/:issue_id/affect/:version_name", controller.CreateOrUpdateIssueAffect)
		issue.GET("/affect/result", controller.SelectIssueAffectResult)
	}

	releaseVersion := router.Group("/version")
	{
		releaseVersion.GET("", controller.SelectReleaseVersion)
		releaseVersion.GET("/maintained", controller.SelectReleaseVersionMaintained)
		releaseVersion.POST("", controller.CreateReleaseVersion)
		releaseVersion.PATCH("", controller.UpdateReleaseVersion)
	}

	versionTriage := router.Group("/version_triage")
	{
		versionTriage.PATCH("", controller.UpdateVersionTriage)
	}

	repo := router.Group("repo")
	{
		repo.GET("", controller.SelectRepo)
	}

	webhook := router.Group("/webhook")
	{
		webhook.POST("", controller.GithubWebhookHandler)
	}

}

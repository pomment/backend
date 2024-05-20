package pomment

import (
	"github.com/gin-gonic/gin"
	"github.com/pomment/pomment/controller"
	"github.com/pomment/pomment/middleware"
)

// InitStandaloneRoutes Initialize base routes for standalone operation
func InitStandaloneRoutes(engine *gin.Engine, prefix string) {
	r := engine

	group := r.Group(prefix)
	group.GET("health", controller.Health)
	group.GET("robots.txt", controller.RobotsTxt)
	group.POST("admin/auth", controller.Auth)

	InitPublicRoutes(group, "public")
	InitWebRoutes(group, "web")
	InitManageRoutes(group, "admin", true)
}

// InitPublicRoutes Initialize routes for public API
func InitPublicRoutes(engine *gin.RouterGroup, prefix string) {
	m := engine.Group(prefix)
	{
		m.GET("thread/meta/:id", controller.ThreadMeta)
		m.POST("thread/meta/byUrl", controller.ThreadMetaByURL)
		m.POST("thread/meta/byUrls", controller.ThreadMetaByURLs)
		m.GET("posts/:id", controller.PostsByID)
		m.POST("posts/byUrl", controller.PostsByURL)
		m.POST("posts/add", controller.PostsAdd)
	}
}

// InitWebRoutes Initialize routes for email interaction web pages
func InitWebRoutes(engine *gin.RouterGroup, prefix string) {
	m := engine.Group(prefix)
	{
		m.GET("unsubscribe/:threadId/:postId/:editKey", controller.Unsubscribe)
		m.POST("unsubscribe/:threadId/:postId/:editKey", controller.UnsubscribeConfirm)
	}
}

// InitManageRoutes Initialize routes for management API
func InitManageRoutes(engine *gin.RouterGroup, prefix string, enableVerify bool) {
	verifyHandler := middleware.NoVerifyToken
	if enableVerify {
		verifyHandler = middleware.VerifyToken
	}
	m := engine.Group(prefix, verifyHandler)
	{
		m.GET("health", controller.Health)
		m.GET("thread/list", controller.ThreadList)
		m.POST("thread/refresh", controller.UpdateAllThreadMeta)
		m.GET("thread/:id", controller.ThreadID)
		m.GET("thread/meta/:id", controller.ThreadMetaID)
		m.PUT("thread/meta", controller.ThreadMetaIDEdit)
		m.GET("posts/:threadId/:postId", controller.PostID)
		m.POST("posts/:id", controller.PostsAddAdmin)
		m.PUT("posts/:id", controller.PostsEdit)
		m.POST("fcm-token", controller.FCMTokenReplace)
		m.DELETE("fcm-token", controller.FCMTokenDelete)
		m.DELETE("cache", controller.CacheDelete)
	}
}

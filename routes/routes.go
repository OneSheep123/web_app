package routes

import (
	"net/http"
	"web_app/controller"
	"web_app/logger"
	"web_app/middlewares"

	swaggerFiles "github.com/swaggo/files"

	gs "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	//if mode == gin.ReleaseMode {
	//	gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	//}
	r := gin.New()
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	v1 := r.Group("/api/v1")
	v1.POST("/login", controller.LoginHandler)
	v1.POST("/signup", controller.SignUpHandler)

	// 业务接口
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		// 获取社区列表
		v1.GET("/community", controller.CommunityHandler)
		// 获取社区详情
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		// 创建帖子
		v1.POST("/post", controller.CreatePost)
		// 获取帖子详情
		v1.GET("/post/:id", controller.GetPostDetail)
		// 获取帖子列表
		v1.GET("/post", controller.GetPostList)

		v1.GET("/post2", controller.PostList2Handler)

		// 创建投票
		v1.POST("/vote", controller.CreateVoteHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}

package router

import (
	"blue_bell/controller"
	"blue_bell/logger"
	"blue_bell/middlewares"
	"blue_bell/settings"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"time"
)

// Setup 注册路由
func Setup(config *settings.Config, mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 模式设置为发布模式，其他为调试模式
	}
	r := gin.New()
	// 使用配置好的日志,中间件的使用
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	v1 := r.Group("/api/v1")

	// 注册业务
	//1.用户注册
	v1.POST("/signup", controller.SignUpHandler)
	// 2. 用户登录
	v1.POST("/login", controller.LoginHandler)

	// 之后的服务需要使用到jwt认证,中间件的使用
	v1.Use(middlewares.JWTAuthMiddleware())

	{ // 社区业务
		//3.社区列表
		v1.GET("/community", controller.CommunityHandler)
		// 4.某个社区的详情
		v1.GET("/community/:id", controller.CommunityDetailHandler)
	}

	{ // post业务   帖子业务
		//5. 创建帖子
		v1.POST("/post", controller.CreatePostHandler)
		// 访问帖子详情
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		// get 获取帖子的列表
		v1.GET("/posts", controller.GetPostListHandler)
		// get 获取帖子的列表 该列表具有按照时间或分数展示帖子功能
		v1.GET("/posts2", controller.GetPostListHandler2)

	}
	{ // 帖子投票功能
		v1.POST("/vote", controller.PostVoteController)
	}

	pprof.Register(r)

	return r
}

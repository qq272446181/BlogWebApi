package main

import (
	"net/http"

	"github.com/qq272446181/BlogWebApi/config"
	"github.com/qq272446181/BlogWebApi/controllers"
	"github.com/qq272446181/BlogWebApi/middleware"
	"github.com/qq272446181/BlogWebApi/models"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// 导入生成的docs包
	_ "github.com/qq272446181/BlogWebApi/docs"
)

// @title 个人博客系统 API
// @version 1.0
// @description 这是一个使用 Go + Gin + GORM 实现的个人博客系统
// @contact.name API Support
// @contact.email 272446181@qq.com
// @license.name MIT
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config.LoadConfig()
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	r := gin.Default()
	r.GET("/swagger.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 根路由
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	initRouter(r)
	r.Run(":8080")
}

func initRouter(r *gin.Engine) {
	//认证路由
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}
	//文章路由
	posts := r.Group("/posts")
	{
		posts.GET("/", controllers.GetPosts)

		posts.GET("/:id", controllers.GetPost)

		authposts := posts.Group("/")
		authposts.Use(middleware.Auth())
		{
			authposts.POST("/", controllers.CreatePost)
			authposts.PUT("/:id", controllers.UpdatePost)

			authposts.DELETE("/:id", controllers.DeletePost)

			authposts.POST("/:id/addcomment", controllers.CreateComment)

			authposts.DELETE("/:id/deletecomment/:commentId", controllers.DeleteComment)
		}
	}
}

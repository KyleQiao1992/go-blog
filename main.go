package main

import (
	"os"

	"go-blog/internal/db"
	"go-blog/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	router := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 健康检查接口，方便测试服务是否启动
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api")

	// 认证相关路由：无需登录即可访问
	authGroup := api.Group("/auth")
	authGroup.POST("/register", handlers.Register)
	authGroup.POST("/login", handlers.Login)

	protected := api.Group("/protected")
	protected.Use(handlers.AuthMiddleware())
	protected.GET("/me", func(c *gin.Context) {
		userID := c.GetUint("userID")
		username := c.GetString("username")

		c.JSON(200, gin.H{
			"userID":   userID,
			"username": username,
		})
	})

	// 启动服务器
	address := ":" + port
	if err := router.Run(address); err != nil {
		panic(err)
	}
}

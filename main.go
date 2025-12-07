package main

import (
	"go-blog/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db.InitDB()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
	r.Run()
}

package handlers

import (
	"net/http"

	"go-blog/internal/auth"
	"go-blog/internal/db"
	"go-blog/internal/logging"
	"go-blog/internal/models"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register handles user registration
// Route: post /api/v1/auth/register
func Register(c *gin.Context) {
	var req RegisterRequest

	//等价写法
	//err := c.ShouldBindJSON(&req)
	//if err != nil {
	// 错误处理
	//}
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数校验失败
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password: " + err.Error(),
		})
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user: " + err.Error(),
		})
		return
	}

	logging.Logger.WithFields(map[string]interface{}{
		"userID":   user.ID,
		"username": user.Username,
	}).Info("user registered successfully")

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

// Login handles user login
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	var user models.User
	if err := db.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid username or password",
		})
		return
	}

	if err := auth.CheckPassword(user.Password, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid username or password",
		})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token: " + err.Error(),
		})
		return
	}

	logging.Logger.WithFields(map[string]interface{}{
		"userID":   user.ID,
		"username": user.Username,
	}).Info("user logged in successfully")

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

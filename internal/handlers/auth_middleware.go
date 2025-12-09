package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-blog/internal/auth"
)

// AuthMiddleware returns a Gin middleware that checks JWT token in Authorization header.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Expect header: Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing Authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid Authorization header format",
			})
			return
		}

		scheme := parts[0]
		tokenString := parts[1]

		if strings.ToLower(scheme) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization scheme must be Bearer",
			})
			return
		}

		// Parse and validate token.
		claims, err := auth.ParseAndValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// Extract user information from claims and put into context.
		// sub is user ID, usr is username.
		userIDValue, okID := claims["sub"]
		usernameValue, okName := claims["usr"]
		if !okID || !okName {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			return
		}

		// JWT 里的数字会以 float64 形式解析，需要转换为 uint。
		floatID, okFloat := userIDValue.(float64)
		username, okString := usernameValue.(string)
		if !okFloat || !okString {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claim types",
			})
			return
		}

		userID := uint(floatID)

		// 保存到 Gin context 中，后续 handler 可以使用。
		c.Set("userID", userID)
		c.Set("username", username)

		// 继续后续处理链
		c.Next()
	}
}

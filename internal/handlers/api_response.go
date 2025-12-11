package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-blog/internal/logging"
)

// ErrorResponse 统一的错误响应结构
type ErrorResponse struct {
	Code    string `json:"code"`    // 用于前端或日志识别的错误码，例如 "AUTH_FAILED"
	Message string `json:"message"` // 返回给用户看的错误信息（简单、可读）
}

// JSONError 用于统一返回错误，并记录日志
func JSONError(c *gin.Context, status int, code string, message string, details map[string]interface{}) {
	// 记录日志，带上路径、方法和 details
	entry := logging.Logger.WithFields(map[string]interface{}{
		"status": status,
		"code":   code,
		"path":   c.FullPath(),
		"method": c.Request.Method,
	})
	for k, v := range details {
		entry = entry.WithField(k, v)
	}
	if status >= http.StatusInternalServerError {
		entry.Error(message)
	} else if status >= http.StatusBadRequest {
		entry.Warn(message)
	} else {
		entry.Info(message)
	}

	// 返回统一格式的 JSON 错误响应
	c.JSON(status, ErrorResponse{
		Code:    code,
		Message: message,
	})
}

package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Logger 是全局日志记录器
var Logger = log.New()

// Init 初始化日志配置
func Init() {
	// 输出到标准输出
	Logger.SetOutput(os.Stdout)

	// 设置日志级别，开发阶段用 Info 或 Debug
	Logger.SetLevel(log.InfoLevel)

	// 设置日志格式
	Logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

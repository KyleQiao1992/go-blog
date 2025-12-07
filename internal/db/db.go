package db

import (
	"go-blog/internal/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	//我调用这个函数，但我故意忽略它的返回值。
	//1️⃣ 加载 .env 文件（读取数据库配置）
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file:", err)
	}

	// 2️⃣ 读取环境变量
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("❌ Missing DB_DSN in .env file")
	}

	// 3️⃣ 用 GORM 连接 MySQL，设置默认字符串长度，避免索引 TEXT 字段报错
	mysqlConfig := mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 191, // 避免 unique/索引字段自动变成 TEXT 导致 1170 错误
	}

	conn, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// 4️⃣ 打印日志提示成功
	log.Println("✅ MySQL connected successfully!")

	conn.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	// 5️⃣ 保存全局数据库连接句柄
	DB = conn
}

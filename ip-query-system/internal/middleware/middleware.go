package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"ip-query-system/internal/config"
)

// Setup 注册中间件
func Setup(r *gin.Engine, cfg *config.Config) {
	// 跨域配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 生产环境替换为具体域名
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 日志中间件（Gin默认）
	r.Use(gin.Logger())

	// 恢复中间件（捕获panic）
	r.Use(gin.Recovery())
}
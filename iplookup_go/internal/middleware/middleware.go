package middleware

import (
	"time"
	"iplookup/iplookup_go/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup 注册中间件
func Setup(r *gin.Engine, cfg *config.Config) {
	// 跨域配置（生产环境建议限制Origin）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 替换为前端实际域名
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 日志和恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
}
package middleware

import (
	"time"
	"iplookup/iplookup_go/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup 注册中间件
func Setup(r *gin.Engine, cfg *config.Config) {
	// 跨域配置（使用配置文件中的允许源）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.API.AllowOrigins,
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
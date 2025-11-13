package middleware

import (
	"net/http"
	"strings"
	"time"
	"iplookup/iplookup_go/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"iplookup/iplookup_go/internal/middleware/ratelimit"
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

	// 限流中间件（使用配置文件中的限流值）
	rl := ratelimit.NewMiddleware(cfg.API.RateLimit, func(c *gin.Context) string {
		// 从X-Forwarded-For或RemoteAddr获取客户端IP
		ip := c.ClientIP()
		if ip == "" {
			ip = c.Request.RemoteAddr
			// 去除端口部分
			parts := strings.Split(ip, ":")
			if len(parts) > 0 {
				ip = parts[0]
			}
		}
		return ip // 使用IP作为限流键
	})
	r.Use(rl.Handler())

	// 日志和恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
}
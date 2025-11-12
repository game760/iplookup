package middleware

import (
	"net/http"
	"strconv"
	"time"
	"iplookup/iplookup_go/internal/config"
	"iplookup/iplookup_go/internal/model"
    "github.com/golang-jwt/jwt"
    "go.uber.org/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

// Setup 注册中间件
func Setup(r *gin.Engine, cfg *config.Config) {
	// 1. 跨域配置（从配置读取允许的源）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     parseAllowOrigins(cfg.API.AllowedOrigins),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 2. 日志和恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 3. 速率限制中间件（公共接口）
	publicLimiter := NewIPRateLimiter(cfg.API.PublicRateLimit)
	r.Use(ratelimit.NewMiddleware(publicLimiter))
}

// 解析允许的源（支持逗号分隔）
func parseAllowOrigins(origins string) []string {
	if origins == "*" {
		return []string{"*"}
	}
	// 简单分割逗号（生产环境可增强为更健壮的解析）
	return strings.Split(origins, ",")
}

// 基于IP的速率限制器
func NewIPRateLimiter(limitStr string) ratelimit.KeyFunc {
	// 解析配置的速率限制（如"30-M"表示每分钟30次）
	limit, period := parseRateLimit(limitStr)
	return func(c *gin.Context) (*rate.Limiter, time.Duration) {
		ip := c.ClientIP()
		limiter := rate.NewLimiter(rate.Limit(limit), int(limit)) // 桶大小=速率
		return limiter, period
	}
}

// 解析速率限制字符串（如"30-M" -> 30次/分钟）
func parseRateLimit(limitStr string) (int, time.Duration) {
	// 简化实现：假设格式为"数字-单位"（M:分钟, H:小时）
	parts := strings.Split(limitStr, "-")
	if len(parts) != 2 {
		return 30, time.Minute // 默认30次/分钟
	}
	num, _ := strconv.Atoi(parts[0])
	switch parts[1] {
	case "H":
		return num, time.Hour
	default:
		return num, time.Minute
	}
}

// JWT认证中间件（用于需要登录的接口）
func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取token
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse("未提供有效令牌"))
			c.Abort()
			return
		}
		tokenStr = tokenStr[7:] // 移除"Bearer "前缀

		// 验证token（此处需补充JWT验证逻辑）
		// 示例：解析token并验证签名和过期时间
		valid, err := verifyJWT(tokenStr, secret)
		if !valid || err != nil {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse("令牌无效或已过期"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// 验证JWT（需实现具体逻辑，依赖jwt库）
func verifyJWT(tokenStr, secret string) (bool, error) {
	// 实际项目中使用github.com/golang-jwt/jwt等库实现
	return true, nil // 仅为示例
}
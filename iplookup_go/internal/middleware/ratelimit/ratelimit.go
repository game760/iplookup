package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// KeyFunc 定义获取限流键的函数类型
type KeyFunc func(c *gin.Context) string

// Middleware 限流中间件结构体
type Middleware struct {
	limiter    *rate.Limiter
	keyFunc    KeyFunc
	limiters   map[string]*rate.Limiter
	mu         sync.RWMutex
	rate       rate.Limit
	burst      int
}

// NewMiddleware 创建新的限流中间件
func NewMiddleware(rateLimit int, keyFunc KeyFunc) *Middleware {
	return &Middleware{
		limiters: make(map[string]*rate.Limiter),
		keyFunc:  keyFunc,
		rate:     rate.Limit(rateLimit), // 每秒允许的请求数
		burst:    rateLimit,             // 突发请求数
	}
}

// 获取或创建限流器
func (m *Middleware) getLimiter(key string) *rate.Limiter {
	m.mu.RLock()
	limiter, exists := m.limiters[key]
	m.mu.RUnlock()

	if !exists {
		m.mu.Lock()
		limiter = rate.NewLimiter(m.rate, m.burst)
		m.limiters[key] = limiter
		m.mu.Unlock()
	}
	return limiter
}

// Handler 实现gin.HandlerFunc接口
func (m *Middleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := m.keyFunc(c)
		if key == "" {
			c.Next()
			return
		}

		limiter := m.getLimiter(key)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
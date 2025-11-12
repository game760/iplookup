package api

import (
	"iplookup/iplookup_go/internal/config"
	"iplookup/iplookup_go/internal/database"
	"iplookup/iplookup_go/internal/ipdb"
	"iplookup/iplookup_go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter 初始化路由
func NewRouter(cfg *config.Config, db *database.DB, ipDB *ipdb.IPDB) *gin.Engine {
	r := gin.Default()
	handler := NewHandler(ipDB)

	// 基础路由组（公共接口）
	v1 := r.Group(cfg.API.Prefix)
	{
		// 无需认证的接口
		v1.GET("/ip/query", IPQueryHandler(ipDB))  // 基础查询
		v1.GET("/ip/my", handler.GetMyIP)          // 本机IP查询

		// 需要认证的接口（应用JWT中间件）
		auth := v1.Group("/")
		auth.Use(middleware.JWTAuthMiddleware(cfg.JWT.Secret))
		{
			auth.GET("/ip/detail", handler.QueryIP) // 详细查询（需登录）
		}
	}

	return r
}
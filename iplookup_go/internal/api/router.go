package api

import (
	"iplookup/iplookup_go/internal/config"
	"iplookup/iplookup_go/internal/ipdb"
	"iplookup/iplookup_go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter 初始化路由
func NewRouter(cfg *config.Config, ipDB *ipdb.IPDB) *gin.Engine {
	r := gin.Default()
	handler := NewHandler(ipDB)

	// 基础路由组（公共接口）
	v1 := r.Group(cfg.API.Prefix)
	{
		// 公共查询接口
		v1.GET("/ip/query/ipv4", handler.QueryIPv4)  // IPv4查询
		v1.GET("/ip/query/ipv6", handler.QueryIPv6)  // IPv6查询
		v1.GET("/ip/query", handler.QueryIP)         // 自动识别IP类型查询
		v1.GET("/ip/my", handler.GetMyIP)            // 本机IP查询
	}

	return r
}
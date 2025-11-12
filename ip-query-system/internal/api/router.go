package api

import (
	"github.com/gin-gonic/gin"
	"ip-query-system/internal/config"
	"ip-query-system/internal/database"
	"ip-query-system/internal/ipdb"
)

// NewRouter 初始化路由
func NewRouter(cfg *config.Config, db *database.DB, ipDB *ipdb.IPDB) *gin.Engine {
	r := gin.Default()
	handler := NewHandler(ipDB) // 初始化处理器

	// 注册路由
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ip/query", IPQueryHandler(ipDB))      // 基础查询接口
		v1.GET("/ip/detail", handler.QueryIP)          // 详细查询接口
		v1.GET("/ip/my", handler.GetMyIP)              // 本机IP查询接口
	}

	return r
}
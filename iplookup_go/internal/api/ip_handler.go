package api

import (
	"net/http"
	"iplookup/iplookup_go/internal/ipdb"
	"iplookup/iplookup_go/internal/model"

	"github.com/gin-gonic/gin"
)

// Handler 处理器结构体
type Handler struct {
	ipDB *ipdb.IPDB
}

// GetDBVersion 获取数据库版本信息
func (h *Handler) GetDBVersion(c *gin.Context) {
	versions := h.ipDB.GetDatabaseVersion()
	c.JSON(http.StatusOK, model.SuccessResponse(versions))
}

// NewHandler 创建处理器实例
func NewHandler(ipDB *ipdb.IPDB) *Handler {
	return &Handler{ipDB: ipDB}
}


// QueryIP 自动识别IP类型查询
func (h *Handler) QueryIP(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("请提供IP地址"))
		return
	}

	ipType := h.ipDB.GetIPType(ip)
	switch ipType {
	case "ipv4":
		resp, err := h.ipDB.QueryIPv4(ip)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse(resp.Message))
			return
		}
		c.JSON(http.StatusOK, resp)
	case "ipv6":
		resp, err := h.ipDB.QueryIPv6(ip)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse(resp.Message))
			return
		}
		c.JSON(http.StatusOK, resp)
	default:
		c.JSON(http.StatusBadRequest, model.ErrorResponse("无效的IP地址"))
	}
}

// GetMyIP 查询本机IP接口
func (h *Handler) GetMyIP(c *gin.Context) {
	ip := c.ClientIP()
	if ip == "" {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse("无法获取本机IP"))
		return
	}

	ipType := h.ipDB.GetIPType(ip)
	var result interface{}
	
	switch ipType {
	case "ipv4":
		resp, err := h.ipDB.QueryIPv4(ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse("查询失败: "+resp.Message))
			return
		}
		result = resp
	case "ipv6":
		resp, err := h.ipDB.QueryIPv6(ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse("查询失败: "+resp.Message))
			return
		}
		result = resp
	default:
		c.JSON(http.StatusInternalServerError, model.ErrorResponse("无法识别IP类型"))
		return
	}

	c.JSON(http.StatusOK, result)
}
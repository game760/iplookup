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

// NewHandler 创建处理器实例
func NewHandler(ipDB *ipdb.IPDB) *Handler {
	return &Handler{ipDB: ipDB}
}

// IPQueryHandler 基础IP查询接口
func IPQueryHandler(ipDB *ipdb.IPDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			IP string `form:"ip" binding:"required"`
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.IPQueryResponse{
				Code:    1,
				Message: "无效的IP参数",
			})
			return
		}

		// 修复：正确处理ipDB.Query返回的单个响应值
		resp := ipDB.Query(req.IP)
		if resp.Code != 0 {
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		// 类型断言提取IPLocation
		location, ok := resp.Data.(*model.IPLocation)
		if !ok {
			c.JSON(http.StatusInternalServerError, model.IPQueryResponse{
				Code:    2,
				Message: "查询结果格式错误",
			})
			return
		}

		// 构建成功响应
		c.JSON(http.StatusOK, model.IPQueryResponse{
			Code:    0,
			Message: "查询成功",
			Data: &model.IPLocation{
				IP:        req.IP,
				Type:      location.Type,
				Country:   location.Country,
				Region:    location.Region,
				City:      location.City,
				Latitude:  location.Latitude,
				Longitude: location.Longitude,
			},
		})
	}
}

// QueryIP 详细IP查询接口
func (h *Handler) QueryIP(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("请提供IP地址"))
		return
	}

	// 修复：正确处理h.ipDB.Query返回值
	resp := h.ipDB.Query(ip)
	if resp.Code != 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(resp.Message))
		return
	}

	location, ok := resp.Data.(*model.IPLocation)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse("查询结果格式错误"))
		return
	}

	// 构建详细响应
	c.JSON(http.StatusOK, model.SuccessResponse(model.IPInfo{
		IP:         ip,
		Type:       location.Type,
		Country:    location.Country,
		Region:     location.Region,
		City:       location.City,
		Latitude:   location.Latitude,
		Longitude:  location.Longitude,
		ZipCode:    "", // 可根据实际IP库补充
		Timezone:   "", // 可根据实际IP库补充
	}))
}

// GetMyIP 查询本机IP接口
func (h *Handler) GetMyIP(c *gin.Context) {
	ip := c.ClientIP()
	if ip == "" {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse("无法获取本机IP"))
		return
	}

	// 修复：正确处理h.ipDB.Query返回值及语法错误（括号匹配）
	resp := h.ipDB.Query(ip)
	if resp.Code != 0 {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse("查询失败: "+resp.Message))
		return
	}

	location, ok := resp.Data.(*model.IPLocation)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse("查询结果格式错误"))
		return
	}

	// 构建本机IP响应
	c.JSON(http.StatusOK, model.SuccessResponse(model.IPInfo{
		IP:         ip,
		Type:       location.Type,
		Country:    location.Country,
		Region:     location.Region,
		City:       location.City,
		Latitude:   location.Latitude,
		Longitude:  location.Longitude,
		ZipCode:    "",
		Timezone:   "",
	}))
}
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ip-query-system/internal/ipdb"
	"ip-query-system/internal/model"
)

// Handler 处理器结构体
type Handler struct {
	ipDB *ipdb.IPDB
}

// NewHandler 创建处理器实例
func NewHandler(ipDB *ipdb.IPDB) *Handler {
	return &Handler{ipDB: ipDB}
}

// IPQueryHandler 处理IP查询请求
func IPQueryHandler(ipDB *ipdb.IPDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.IPQueryRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.IPQueryResponse{
				Code:    1,
				Message: "无效的请求参数: " + err.Error(),
			})
			return
		}

		record, ipType, err := ipDB.Query(req.IP)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.IPQueryResponse{
				Code:    2,
				Message: "查询失败: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, model.IPQueryResponse{
			Code:    0,
			Message: "查询成功",
			Data: &model.IPLocation{
				IP:        req.IP,
				Type:      ipType,
				Country:   record.CountryName,
				Region:    record.RegionName,
				City:      record.CityName,
				Latitude:  record.Latitude,
				Longitude: record.Longitude,
			},
		})
	}
}

func (h *Handler) QueryIP(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("请提供IP地址"))
		return
	}

	record, ipType, err := h.ipDB.Query(ip)
	if err != nil {
		if err.Error() == "无效的IP地址格式" || err.Error() == "不支持的IP地址类型" {
			c.JSON(http.StatusBadRequest, model.ErrorResponse(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse("查询失败: "+err.Error()))
		}
		return
	}

	// 构建响应（注意：根据数据库结构，移除了原有的ISP/ASN等不存在的字段）
	response := model.SuccessResponse(model.IPInfo{
		IP:         ip,
		Type:       ipType,
		Country:    record.CountryName,
		Region:     record.RegionName,
		City:       record.CityName,
		Latitude:   record.Latitude,
		Longitude:  record.Longitude,
		ZipCode:    record.ZipCode,
		Timezone:   record.TimeZone,
	})

	c.JSON(http.StatusOK, response)
}

// GetMyIP 查询本机IP
func (h *Handler) GetMyIP(c *gin.Context) {
	ip := c.ClientIP()
	
	record, ipType, err := h.ipDB.Query(ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse("查询失败: "+err.Error()))
		return
	}

	response := model.SuccessResponse(model.IPInfo{
		IP:         ip,
		Type:       ipType,
		Country:    record.CountryName,
		Region:     record.RegionName,
		City:       record.CityName,
		Latitude:   record.Latitude,
		Longitude:  record.Longitude,
		ZipCode:    record.ZipCode,
		Timezone:   record.TimeZone,
	})

	c.JSON(http.StatusOK, response)
}
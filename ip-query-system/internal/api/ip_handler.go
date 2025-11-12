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

		resp := ipDB.Query(req.IP)
		if resp.Code != 0 {
			c.JSON(http.StatusBadRequest, resp) // 直接返回查询失败的响应
			return
			}
		}

// 从响应数据中提取IPLocation
location, ok := resp.Data.(*model.IPLocation)
if !ok {
    c.JSON(http.StatusInternalServerError, model.IPQueryResponse{
        Code:    2,
        Message: "查询结果格式错误",
    })
    return
}

// 构建成功响应（使用location中的数据）
c.JSON(http.StatusOK, model.IPQueryResponse{
    Code:    0,
    Message: "查询成功",
    Data: &model.IPLocation{
        IP:        req.IP,
        Type:      location.Type,       // 从location获取IP类型
        Country:   location.Country,    // 从location获取国家
        Region:    location.Region,     // 从location获取地区
        City:      location.City,       // 从location获取城市
        Latitude:  location.Latitude,
        Longitude: location.Longitude,
    },
})
	}


func (h *Handler) QueryIP(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("请提供IP地址"))
		return
	}

	// 原错误代码
record, ipType, err := h.ipDB.Query(ip)
if err != nil {
    // ...
}

// 修复后代码
resp := h.ipDB.Query(ip)
if resp.Code != 0 {
    c.JSON(http.StatusBadRequest, model.ErrorResponse(resp.Message))
    return
}

// 从响应数据中提取IPLocation
location, ok := resp.Data.(*model.IPLocation)
if !ok {
    c.JSON(http.StatusInternalServerError, model.ErrorResponse("查询结果格式错误"))
    return
}

// 构建详细信息响应
response := model.SuccessResponse(model.IPInfo{
    IP:         ip,
    Type:       location.Type,       // 从location获取IP类型
    Country:    location.Country,    // 从location获取国家
    Region:     location.Region,     // 从location获取地区
    City:       location.City,       // 从location获取城市
    Latitude:   location.Latitude,
    Longitude:  location.Longitude,
    ZipCode:    "", // 若有需要可从IP库补充
    Timezone:   "", // 若有需要可从IP库补充
})

c.JSON(http.StatusOK, response)
// GetMyIP 查询本机IP
func (h *Handler) GetMyIP(c *gin.Context) {
	ip := c.ClientIP()
	
// 修复后代码
resp := h.ipDB.Query(ip)
if resp.Code != 0 {
    c.JSON(http.StatusInternalServerError, model.ErrorResponse(resp.Message))
    return
}

// 从响应数据中提取IPLocation
location, ok := resp.Data.(*model.IPLocation)
if !ok {
    c.JSON(http.StatusInternalServerError, model.ErrorResponse("查询结果格式错误"))
    return
}

// 构建本机IP信息响应
response := model.SuccessResponse(model.IPInfo{
    IP:         ip,
    Type:       location.Type,       // 从location获取IP类型
    Country:    location.Country,    // 从location获取国家
    Region:     location.Region,     // 从location获取地区
    City:       location.City,       // 从location获取城市
    Latitude:   location.Latitude,
    Longitude:  location.Longitude,
    ZipCode:    "", // 若有需要可从IP库补充
    Timezone:   "", // 若有需要可从IP库补充
})

c.JSON(http.StatusOK, response)
}
package model

// IPQueryRequest IP查询请求参数
type IPQueryRequest struct {
	IP string `form:"ip" binding:"required"`
}

// IPQueryResponse IP查询响应结构
type IPQueryResponse struct {
	Code    int         `json:"code"`   // 0：成功 非0：错误
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`   // 业务数据（IP地理位置信息）
}

// IPLocation IP地理位置信息
type IPLocation struct {
	IP        string  `json:"ip"`
	Type      string  `json:"type"`
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// IPInfo 详细IP信息（扩展字段）
type IPInfo struct {
	IP         string  `json:"ip"`
	Type       string  `json:"type"`
	Country    string  `json:"country"`
	Region     string  `json:"region"`
	City       string  `json:"city"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	ISP        string  `json:"isp"`
	Domain     string  `json:"domain"`
	ASN        string  `json:"asn"`
	ASName     string  `json:"as_name"`
	Timezone   string  `json:"timezone"`
	ZipCode    string  `json:"zip_code"`
	UsageType  string  `json:"usage_type"`
}

// SuccessResponse 通用成功响应（修复：使用传入的data参数）
func SuccessResponse(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    data, // 正确使用函数参数data（小写）
	}
}

// ErrorResponse 通用错误响应（修复：错误响应数据设为nil）
func ErrorResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"code":    1,
		"message": message,
		"data":    nil, // 错误时无业务数据，设为nil
	}
}

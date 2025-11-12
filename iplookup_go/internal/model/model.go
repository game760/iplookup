package model

// IPQueryResponse IP查询通用响应
type IPQueryResponse struct {
	Code    int         `json:"code"`    // 0:成功，非0:错误
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据（成功时返回）
}

// IPInfo 详细IP信息（用于详细查询接口）
type IPInfo struct {
	IP         string  `json:"ip"`
	Type       string  `json:"type"`
	Country    string  `json:"country"`
	Region     string  `json:"region"`
	City       string  `json:"city"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	ZipCode    string  `json:"zip_code,omitempty"`
	Timezone   string  `json:"timezone,omitempty"`
	ISP        string  `json:"isp,omitempty"` // 互联网服务提供商
}

// SuccessResponse 通用成功响应
func SuccessResponse(data interface{}) IPQueryResponse {
	return IPQueryResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

// ErrorResponse 通用错误响应
func ErrorResponse(message string) IPQueryResponse {
	return IPQueryResponse{
		Code:    1,
		Message: message,
		Data:    nil,
	}
}
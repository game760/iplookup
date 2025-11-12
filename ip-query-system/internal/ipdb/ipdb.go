package ipdb

import (

	"ip-query-system/internal/model" // 导入model包（必须）

)

// 假设IP数据库查询的核心逻辑
type IPDB struct {

	DBPath string
}

// NewIPDB 初始化IP数据库实例
func NewIPDB(dbPath string) *IPDB {
	return &IPDB{
		DBPath: dbPath,
	}
}

// Query 查询IP地理位置信息
func (db *IPDB) Query(ip string) model.IPQueryResponse {
	// 模拟无效IP校验
	if ip == "" {
		return model.IPQueryResponse{
			Code:    1,
			Message: "无效的IP地址",
			Data:    nil,
		}
	}

	// 模拟查询逻辑
	location := &model.IPLocation{
		IP:        ip,
		Type:      "IPv4",
		Country:   "中国",
		Region:    "北京",
		City:      "北京",
		Latitude:  39.9042,
		Longitude: 116.4074,
	}

	// 返回成功响应
	return model.IPQueryResponse{
		Code:    0,
		Message: "查询成功",
		Data:    location,
	}
}

// BatchQuery 批量查询IP
func (db *IPDB) BatchQuery(ips []string) model.IPQueryResponse {
	if len(ips) == 0 {
		return model.IPQueryResponse{
			Code:    1,
			Message: "IP列表不能为空",
			Data:    nil,
		}
	}

	var locations []*model.IPLocation
	for _, ip := range ips {
		locations = append(locations, &model.IPLocation{
			IP:      ip,
			Country: "中国",
		})
	}

	return model.IPQueryResponse{
		Code:    0,
		Message: "批量查询成功",
		Data:    locations,  // 切片类型，现在可匹配interface{}
	}
}
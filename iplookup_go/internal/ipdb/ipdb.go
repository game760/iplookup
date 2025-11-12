package ipdb

import (
	"database/sql"
	"errors"
	"net"
	"iplookup/iplookup_go/internal/config"
	"iplookup/iplookup_go/internal/model"
)

// IPLocation 存储IP地理位置信息
type IPLocation struct {
	IP        string  `json:"ip"`
	Type      string  `json:"type"` // "ipv4" 或 "ipv6"
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// IPDB IP数据库查询器
type IPDB struct {
	db        *sql.DB
	ipv4Table string
	ipv6Table string
}

// Init 初始化IP查询器
func Init(db *sql.DB, cfg *config.Config) (*IPDB, error) {
	return &IPDB{
		db:        db,
		ipv4Table: cfg.IPDatabase.IPv4Table,
		ipv6Table: cfg.IPDatabase.IPv6Table,
	}, nil
}

// Query 查询IP信息
func (ipdb *IPDB) Query(ipStr string) (model.IPQueryResponse, error) {
	// 1. 验证IP类型
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return model.IPQueryResponse{
			Code:    1,
			Message: "无效的IP地址",
		}, errors.New("invalid ip")
	}

	// 2. 区分IPv4/IPv6，选择对应表
	var table string
	var ipType string
	if ip.To4() != nil {
		table = ipdb.ipv4Table
		ipType = "ipv4"
	} else if ip.To16() != nil {
		table = ipdb.ipv6Table
		ipType = "ipv6"
	} else {
		return model.IPQueryResponse{
			Code:    2,
			Message: "不支持的IP类型",
		}, errors.New("unsupported ip type")
	}

	// 3. 从数据库查询（假设表结构包含country/region/city/latitude/longitude字段）
	querySQL := `
		SELECT country, region, city, latitude, longitude 
		FROM ` + table + ` 
		WHERE ip_start <= INET6_ATON(?) AND ip_end >= INET6_ATON(?)
		LIMIT 1
	`
	var loc IPLocation
	err := ipdb.db.QueryRow(querySQL, ipStr, ipStr).Scan(
		&loc.Country,
		&loc.Region,
		&loc.City,
		&loc.Latitude,
		&loc.Longitude,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.IPQueryResponse{
				Code:    3,
				Message: "未查询到IP信息",
			}, nil
		}
		return model.IPQueryResponse{
			Code:    4,
			Message: "数据库查询失败",
		}, err
	}

	// 4. 组装结果
	loc.IP = ipStr
	loc.Type = ipType
	return model.IPQueryResponse{
		Code:    0,
		Message: "查询成功",
		Data:    loc,
	}, nil
}

// Close 关闭资源（如果需要）
func (ipdb *IPDB) Close() error {
	return nil
}
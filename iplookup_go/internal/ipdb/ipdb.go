package ipdb

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"iplookup/iplookup_go/internal/config"
	"iplookup/iplookup_go/internal/model"
)

// IPDB 封装IPv4和IPv6数据库查询功能
type IPDB struct {
	v4db       *xdb.Searcher // IPv4数据库
	v6db       *xdb.Searcher // IPv6数据库
	v4Version  int           // IPv4数据库版本
	v6Version  int           // IPv6数据库版本
}

// Init 初始化IPv4和IPv6数据库
func Init(cfg *config.Config) (*IPDB, error) {
	// 加载IPv4数据库
	v4Version, v4Data, err := xdb.LoadContentFromFile(cfg.IPDatabase.IPv4DB)
	if err != nil {
		return nil, fmt.Errorf("无法加载IPv4数据库: %w", err)
	}
	v4Searcher, err := xdb.NewWithBuffer(v4Version, v4Data)
	if err != nil {
		return nil, fmt.Errorf("初始化IPv4查询器失败: %w", err)
	}

	// 加载IPv6数据库
	v6Version, v6Data, err := xdb.LoadContentFromFile(cfg.IPDatabase.IPv6DB)
	if err != nil {
		return nil, fmt.Errorf("无法加载IPv6数据库: %w", err)
	}
	v6Searcher, err := xdb.NewWithBuffer(v6Version, v6Data)
	if err != nil {
		return nil, fmt.Errorf("初始化IPv6查询器失败: %w", err)
	}

	return &IPDB{
		v4db:       v4Searcher,
		v6db:       v6Searcher,
		v4Version:  v4Version,
		v6Version:  v6Version,
	}, nil
}

// Close 关闭数据库连接
func (ipdb *IPDB) Close() error {
	ipdb.v4db.Close()
	ipdb.v6db.Close()
	return nil
}

// GetDatabaseVersion 获取数据库版本信息
func (ipdb *IPDB) GetDatabaseVersion() map[string]int {
	return map[string]int{
		"ipv4": ipdb.v4Version,
		"ipv6": ipdb.v6Version,
	}
}

// 解析ip2region返回格式: 国家|区域|省份|城市|ISP
func parseRegionData(data string) []string {
	// 处理空结果
	if data == "" {
		return []string{"", "", "", "", ""}
	}
	
	parts := strings.SplitN(data, "|", 5)
	// 确保返回5个元素
	if len(parts) < 5 {
		padding := make([]string, 5-len(parts))
		parts = append(parts, padding...)
	}
	return parts
}

// QueryIPv4 查询IPv4地址信息
func (ipdb *IPDB) QueryIPv4(ipStr string) (model.IPv4Response, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil || ip.To4() == nil {
		return model.IPv4Response{
			Code:    1,
			Message: "无效的IPv4地址",
		}, errors.New("invalid ipv4 address")
	}

	// 使用IP对象查询，更高效
	result, err := ipdb.v4db.SearchByIP(ip)
	if err != nil {
		return model.IPv4Response{
			Code:    2,
			Message: "查询失败: " + err.Error(),
		}, err
	}

	parts := parseRegionData(result)
	
	// 注意：ip2region数据库本身不提供经纬度，这里保持0或可考虑其他数据源补充
	return model.IPv4Response{
		Code:    0,
		Message: "查询成功",
		Data: model.IPv4Info{
			IP:          ipStr,
			CountryName: parts[0],
			Region:      parts[1],
			Province:    parts[2],
			City:        parts[3],
			ISP:         parts[4],
			Latitude:    0,
			Longitude:   0,
		},
	}, nil
}

// QueryIPv6 查询IPv6地址信息
func (ipdb *IPDB) QueryIPv6(ipStr string) (model.IPv6Response, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil || ip.To16() == nil || ip.To4() != nil {
		return model.IPv6Response{
			Code:    1,
			Message: "无效的IPv6地址",
		}, errors.New("invalid ipv6 address")
	}

	// 使用IP对象查询，更高效
	result, err := ipdb.v6db.SearchByIP(ip)
	if err != nil {
		return model.IPv6Response{
			Code:    2,
			Message: "查询失败: " + err.Error(),
		}, err
	}

	parts := parseRegionData(result)
	
	return model.IPv6Response{
		Code:    0,
		Message: "查询成功",
		Data: model.IPv6Info{
			IP:          ipStr,
			CountryName: parts[0],
			Region:      parts[1],
			Province:    parts[2],
			City:        parts[3],
			ISP:         parts[4],
			Latitude:    0,
			Longitude:   0,
		},
	}, nil
}

// GetIPType 判断IP类型
func (ipdb *IPDB) GetIPType(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "invalid"
	}
	if ip.To4() != nil {
		return "ipv4"
	}
	if ip.To16() != nil {
		return "ipv6"
	}
	return "unknown"
}
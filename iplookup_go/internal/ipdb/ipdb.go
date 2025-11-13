package ipdb

import (
	"errors"
	"net"
	"strconv"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"iplookup/iplookup_go/internal/config"
	"iplookup/iplookup_go/internal/model"
)

// IPDB 封装IPv4和IPv6数据库查询功能
type IPDB struct {
	v4db *xdb.Searcher // IPv4数据库
	v6db *xdb.Searcher // IPv6数据库
}

// Init 初始化IPv4和IPv6数据库（修复返回值不匹配问题）
func Init(cfg *config.Config) (*IPDB, error) {
	// 加载IPv4数据库（修复：xdb.LoadContentFromFile返回2个值）
	v4Data, err := xdb.LoadContentFromFile(cfg.IPDatabase.IPv4DB)
	if err != nil {
		return nil, errors.New("无法加载IPv4数据库: " + err.Error())
	}
	v4Searcher, err := xdb.NewWithBuffer(v4Data) // 修复：仅传入数据
	if err != nil {
		return nil, errors.New("初始化IPv4查询器失败: " + err.Error())
	}

	// 加载IPv6数据库（修复：xdb.LoadContentFromFile返回2个值）
	v6Data, err := xdb.LoadContentFromFile(cfg.IPDatabase.IPv6DB)
	if err != nil {
		return nil, errors.New("无法加载IPv6数据库: " + err.Error())
	}
	v6Searcher, err := xdb.NewWithBuffer(v6Data) // 修复：仅传入数据
	if err != nil {
		return nil, errors.New("初始化IPv6查询器失败: " + err.Error())
	}

	return &IPDB{
		v4db: v4Searcher,
		v6db: v6Searcher,
	}, nil
}

// Close 关闭数据库连接
func (ipdb *IPDB) Close() error {
	ipdb.v4db.Close()
	ipdb.v6db.Close()
	return nil
}

// 解析ip2region返回格式: 国家|区域|省份|城市|ISP
func parseRegionData(data string) []string {
	parts := make([]string, 5)
	current := ""
	idx := 0
	for _, c := range data {
		if c == '|' {
			parts[idx] = current
			current = ""
			idx++
			if idx >= 5 {
				break
			}
		} else {
			current += string(c)
		}
	}
	if idx < 5 {
		parts[idx] = current
	}
	return parts
}

// QueryIPv4 查询IPv4地址信息（修复方法不存在问题）
func (ipdb *IPDB) QueryIPv4(ipStr string) (model.IPv4Response, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil || ip.To4() == nil {
		return model.IPv4Response{
			Code:    1,
			Message: "无效的IPv4地址",
		}, errors.New("invalid ipv4 address")
	}

	// 修复：使用 SearchByStr 方法（而非 SearchByIP），传入字符串IP
	result, err := ipdb.v4db.SearchByStr(ipStr)
	if err != nil {
		return model.IPv4Response{
			Code:    2,
			Message: "查询失败: " + err.Error(),
		}, err
	}

	parts := parseRegionData(result)
	lat, _ := strconv.ParseFloat("0", 64)
	lng, _ := strconv.ParseFloat("0", 64)

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
			Latitude:    lat,
			Longitude:   lng,
		},
	}, nil
}

// QueryIPv6 查询IPv6地址信息（修复方法不存在问题）
func (ipdb *IPDB) QueryIPv6(ipStr string) (model.IPv6Response, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil || ip.To16() == nil || ip.To4() != nil {
		return model.IPv6Response{
			Code:    1,
			Message: "无效的IPv6地址",
		}, errors.New("invalid ipv6 address")
	}

	// 修复：使用 SearchByStr 方法（而非 SearchByIP），传入字符串IP
	result, err := ipdb.v6db.SearchByStr(ipStr)
	if err != nil {
		return model.IPv6Response{
			Code:    2,
			Message: "查询失败: " + err.Error(),
		}, err
	}

	parts := parseRegionData(result)
	lat, _ := strconv.ParseFloat("0", 64)
	lng, _ := strconv.ParseFloat("0", 64)

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
			Latitude:    lat,
			Longitude:   lng,
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
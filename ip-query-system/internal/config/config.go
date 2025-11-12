package config

import (
	"os"

	"gopkg.in/yaml.v3" // 修正为标准YAML包（原yml.v3为错误依赖）
)

// Config 全局配置结构体
// 标签同时支持 yml 和 yaml，解析时优先使用 yml 标签（yaml.v3会优先匹配第一个有效标签）
type Config struct {
	Server struct {
		Port         string `yml:"port" yaml:"port"`         // 服务器端口
		ReadTimeout  int    `yml:"read_timeout" yaml:"read_timeout"` // 读取超时（秒）
		WriteTimeout int    `yml:"write_timeout" yaml:"write_timeout"` // 写入超时（秒）
	} `yml:"server" yaml:"server"`

	Database struct {
		Host     string `yml:"host" yaml:"host"`     // MySQL地址
		Port     string `yml:"port" yaml:"port"`     // MySQL端口
		User     string `yml:"user" yaml:"user"`     // MySQL用户
		Password string `yml:"password" yaml:"password"` // MySQL密码
		DBName   string `yml:"db_name" yaml:"db_name"`   // 数据库名（默认ip2location）
	} `yml:"database" yaml:"database"`

	IPDatabase struct {
		IPv4Table string `yml:"ipv4_table" yaml:"ipv4_table"` // IPv4数据表名
		IPv6Table string `yml:"ipv6_table" yaml:"ipv6_table"` // IPv6数据表名
	} `yml:"ip_database" yaml:"ip_database"`
}

// Load 从YML/YAML文件加载配置，支持两种标签格式（优先yml）
func Load(path string) (*Config, error) {
	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err // 错误包含文件读取失败信息（如文件不存在）
	}

	var cfg Config
	// 解析YAML内容，yaml.v3会优先识别结构体中的`yml`标签（若存在），同时兼容`yaml`标签
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err // 错误包含解析失败信息（如格式错误）
	}

	// 设置完整默认值（避免空值导致后续业务错误）
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080" // 默认服务器端口
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 30 // 默认读取超时30秒
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 30 // 默认写入超时30秒
	}

	if cfg.Database.Host == "" {
		cfg.Database.Host = "localhost" // 默认MySQL主机
	}
	if cfg.Database.Port == "" {
		cfg.Database.Port = "3306" // 默认MySQL端口
	}
	if cfg.Database.User == "" {
		cfg.Database.User = "root" // 默认MySQL用户（可根据实际场景调整）
	}
	// 密码默认空（通常由用户在配置文件中指定，不强制默认值）

	if cfg.Database.DBName == "" {
		cfg.Database.DBName = "ip2location" // 数据库名默认值
	}

	if cfg.IPDatabase.IPv4Table == "" {
		cfg.IPDatabase.IPv4Table = "ip2location_db11" // IPv4表默认名
	}
	if cfg.IPDatabase.IPv6Table == "" {
		cfg.IPDatabase.IPv6Table = "ip2location_db11_ipv6" // IPv6表默认名
	}

	return &cfg, nil
}
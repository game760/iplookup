package config  // 关键：添加包声明

import (
	"os"
	"gopkg.in/yaml.v3"
)

// Config 全局配置结构体
type Config struct {
	Server struct {
		Port         string `yml:"port" yaml:"port"`
		ReadTimeout  int    `yml:"read_timeout" yaml:"read_timeout"`
		WriteTimeout int    `yml:"write_timeout" yaml:"write_timeout"`
	} `yml:"server" yaml:"server"`

	IPDatabase struct {
		IPv4DB string `yml:"ipv4_db" yaml:"ipv4_db"`
		IPv6DB string `yml:"ipv6_db" yaml:"ipv6_db"`
		IPv4Version string `yml:"ipv4_version" yaml:"ipv4_version"`
		IPv6Version string `yml:"ipv6_version" yaml:"ipv6_version"`
	} `yml:"ip_database" yaml:"ip_database"`

	API struct {
		Prefix       string   `yml:"prefix" yaml:"prefix"`
		AllowOrigins []string `yml:"allow_origins" yaml:"allow_origins"`
	} `yml:"api" yaml:"api"`
}

// Load 从YML/YAML文件加载配置
func Load(path string) (*Config, error) {
	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	// 解析YAML内容
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// 设置默认值
	if cfg.Server.Port == "" {
		cfg.Server.Port = "7895"
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 30
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 30
	}

	if cfg.IPDatabase.IPv4DB == "" {
		cfg.IPDatabase.IPv4DB = "./ipdata/ipv4.xdb"
	}
	if cfg.IPDatabase.IPv6DB == "" {
		cfg.IPDatabase.IPv6DB = "./ipdata/ipv6.xdb"
	}

	if cfg.API.Prefix == "" {
		cfg.API.Prefix = "/api/v1"
	}
	if len(cfg.API.AllowOrigins) == 0 {
		cfg.API.AllowOrigins = []string{"*"}
	}

	return &cfg, nil
}
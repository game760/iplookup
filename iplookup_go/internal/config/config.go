package config

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

	Database struct {
		Host     string `yml:"host" yaml:"host"`
		Port     string `yml:"port" yaml:"port"`
		User     string `yml:"user" yaml:"user"`
		Password string `yml:"password" yaml:"password"`
		DBName   string `yml:"db_name" yaml:"db_name"`
	} `yml:"database" yaml:"database"`

	IPDatabase struct {
		IPv4Table string `yml:"ipv4_table" yaml:"ipv4_table"`
		IPv6Table string `yml:"ipv6_table" yaml:"ipv6_table"`
	} `yml:"ip_database" yaml:"ip_database"`
}
	// 添加API相关配置（解决cfg.API未定义问题）
	API struct {
		JWTSecret   string `yml:"jwt_secret" yaml:"jwt_secret"`
		RateLimit   int    `yml:"rate_limit" yaml:"rate_limit"`
		AllowOrigins []string `yml:"allow_origins" yaml:"allow_origins"`
	} `yml:"api" yaml:"api"`
}

// Load 从YML/YAML文件加载配置，支持两种标签格式（优先yml）
func Load(path string) (*Config, error) {
	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	// 解析YAML内容，yaml.v3会优先识别结构体中的`yml`标签（若存在），同时兼容`yaml`标签
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// 设置完整默认值（避免空值导致后续业务错误）
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 30
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 30
	}

	if cfg.Database.Host == "" {
		cfg.Database.Host = "localhost"
	}
	if cfg.Database.Port == "" {
		cfg.Database.Port = "3306"
	}
	if cfg.Database.User == "" {
		cfg.Database.User = "root"
	}
	// 密码默认空（通常由用户在配置文件中指定，不强制默认值）

	if cfg.Database.DBName == "" {
		cfg.Database.DBName = "ip2location"
	}

	if cfg.IPDatabase.IPv4Table == "" {
		cfg.IPDatabase.IPv4Table = "ip2location_db11"
	}
	if cfg.IPDatabase.IPv6Table == "" {
		cfg.IPDatabase.IPv6Table = "ip2location_db11_ipv6"
	}
// API配置默认值
	if cfg.API.RateLimit == 0 {
		cfg.API.RateLimit = 100
	}
	if len(cfg.API.AllowOrigins) == 0 {
		cfg.API.AllowOrigins = []string{"*"}
	}
	return &cfg, nil
}
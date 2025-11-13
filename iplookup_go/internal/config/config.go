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
	} `yml:"ip_database" yaml:"ip_database"`

	API struct {
		Prefix       string   `yml:"prefix" yaml:"prefix"`
		RateLimit    int      `yml:"rate_limit" yaml:"rate_limit"`
		AllowOrigins []string `yml:"allow_origins" yaml:"allow_origins"`
	} `yml:"api" yaml:"api"`
}
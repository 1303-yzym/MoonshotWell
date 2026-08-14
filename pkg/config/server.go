package config

type ServerConfig struct {
	HTTP HTTPConfig `json:"http" comment:"http配置"`
}

type HTTPConfig struct {
	Port         int    `mapstructure:"port" json:"port" comment:"端口"`
	RouterPrefix string `mapstructure:"router_prefix" json:"router_prefix" comment:"路由前缀"`
}

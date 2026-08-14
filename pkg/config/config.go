package config

type Env string

const (
	EnvDev  Env = "dev"
	EnvTest Env = "test"
	EnvProd Env = "prod"
)

type Config[E any] struct {
	BasicConfig `mapstructure:",squash"`
	Remain      map[string]any `mapstructure:",remain" json:"remain"`
	Exp         E              `mapstructure:"exp" json:"exp" comment:"扩展配置"`
}

type BasicConfig struct {
	Env     Env           `json:"env" comment:"运行环境"`
	Server  ServerConfig  `json:"server" comment:"服务器配置"`
	Storage StorageConfig `json:"storage" comment:"存储配置"`
	// TODO DB Server Redis ……
}

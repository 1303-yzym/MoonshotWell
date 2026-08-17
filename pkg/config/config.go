package config

// Env 运行环境，不同的运行环境应有不同的行为限制。
type Env string

const (
	// EnvDev 开发环境
	EnvDev Env = "dev"
	// EnvTest 测试环境
	EnvTest Env = "test"
	// EnvProd 生产环境
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

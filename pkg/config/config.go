package config

import (
	"github.com/1303-yzym/MoonshotWell/pkg/infra"
	"github.com/1303-yzym/MoonshotWell/pkg/infra/DB"
)

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
	Env         Env             `json:"env" comment:"运行环境"`
	Log         infra.LogConfig `mapstructure:"log" json:"log" comment:"日志存储配置"`
	ServiceName string          `mapstructure:"service_name" json:"service_name" comment:"服务名称"`
	Server      ServerConfig    `json:"server" comment:"服务器配置"`
	DB          db.DBConfig     `mapstructure:"db" json:"db" comment:"数据库配置"`
	// TODO DB Server Redis ……
}

// Validation 配置文件验证
func (c BasicConfig) Validation() error {
	// todo: 配置文件验证
	return nil
}

// IsDev 是否是开发环境
func (c BasicConfig) IsDev() bool {
	return c.Env == EnvDev
}

// IsProd 是否是生产环境
func (c BasicConfig) IsProd() bool {
	return c.Env == EnvProd
}

// IsTest 是否是测试环境
func (c BasicConfig) IsTest() bool {
	return c.Env == EnvTest
}

//type BasicConfig struct {
//	Env           string                    `mapstructure:"env" json:"env" comment:"环境 dev | test | prod"`
//	AppPrivateKey string                    `mapstructure:"app_private_key" json:"app_private_key" comment:"应用加密数据私钥"`
//	DeviceId      int64                     `mapstructure:"device_id" json:"device_id" comment:"服务设备id"`
//	ServiceName   string                    `mapstructure:"service_name" json:"service_name" comment:"服务名称"`
//	SecurityKit   infra.SecurityKit         `mapstructure:"security" json:"security" comment:"安全套件"`
//	Log           infra.LogConfig           `mapstructure:"log" json:"log" comment:"日志存储配置"`
//	Server        ServerConfig              `mapstructure:"server" json:"server" comment:"服务器配置"`
//	JWT           jwts.JWTConf              `mapstructure:"jwt" json:"jwt" comment:"访问令牌配置"`
//	DB            infra.DBConfig            `mapstructure:"db" json:"db" comment:"数据库配置"`
//	Redis         infra.RedisConfig         `mapstructure:"redis" json:"redis" comment:"Redis配置"`
//	ES            infra.ElasticsearchConfig `mapstructure:"es" json:"es" comment:"ES数据库配置"`
//	MQ            infra.MQConfig            `mapstructure:"mq" json:"mq" comment:"消息队列配置"`
//	Storage       StorageConfig             `mapstructure:"storage" json:"storage" comment:"对象存储配置"`
//	Sdk           infra.SDKConfig           `mapstructure:"sdk" json:"sdk" comment:"第三方SDK配置"`
//	Pay           infra.PayConfig           `mapstructure:"pay" json:"pay" comment:"支付"`
//	Otel          infra.OtelConfig          `mapstructure:"otel" json:"otel" comment:"otel"`
//}

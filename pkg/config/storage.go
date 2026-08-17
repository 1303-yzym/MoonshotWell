package config

import (
	"github.com/1303-yzym/MoonshotWell/pkg/infra/DB"
)

type StorageConfig struct {
	DB DB.MySQLConfig `json:"mysql" comment:"数据库配置"`
}

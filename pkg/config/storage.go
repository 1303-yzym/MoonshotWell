package config

import "github.com/1303-yzym/MoonshotWell/pkg/infra"

type StorageConfig struct {
	DB infra.DBConfig `json:"db" comment:"数据库配置"`
}

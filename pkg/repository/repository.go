package repository

import (
	"github.com/1303-yzym/MoonshotWell/models/sqls"
	"github.com/1303-yzym/MoonshotWell/pkg/config"
	"github.com/1303-yzym/MoonshotWell/pkg/contract"
	"github.com/1303-yzym/MoonshotWell/pkg/crud"
)

// Common 所有模型均适用的CRUD
type Common struct {
	// TODO 所有的model均需要在此注册
	User crud.CRUD[sqls.UserModel]
}

func New(storage contract.Storage, env config.Env) *Common {
	return &Common{
		User: crud.New[sqls.UserModel](storage, env),
	}
}

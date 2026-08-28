package repository

import (
	"github.com/1303-yzym/MoonshotWell/pkg/config"
	"github.com/1303-yzym/MoonshotWell/pkg/contract"
	repo "github.com/1303-yzym/MoonshotWell/pkg/repository"
)

type Repository struct {
	// 一般CRUD
	Common *repo.Common

	// TODO 注册特殊CRUD
}

func New(storage contract.Storage, env config.Env) *Repository {
	return &Repository{
		Common: repo.New(storage, env),
	}
}

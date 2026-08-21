package contract

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/elastic/go-elasticsearch/v9/typedapi"
	"github.com/redis/go-redis/v9"
)

type DB interface {
	DB() *gorm.DB
}

type Redis interface {
	RDB() *redis.Client
}

type ES interface {
	ES() *typedapi.API
}

type Storage interface {
	DB
	Redis
	ES
	LOG(ctx context.Context) *zap.Logger
}

type StorageImpl struct {
	log   *zap.Logger
	db    *gorm.DB
	redis *redis.Client
	es    *typedapi.API
}

func NewStorage(db *gorm.DB, redis *redis.Client, es *typedapi.API) Storage {
	return &StorageImpl{db: db, redis: redis, es: es}
}

func (s *StorageImpl) LOG(ctx context.Context) *zap.Logger {
	return LOGWith(ctx, s.log)
}

func (s *StorageImpl) DB() *gorm.DB {
	return s.db
}

func (s *StorageImpl) RDB() *redis.Client {
	return s.redis
}

func (s *StorageImpl) ES() *typedapi.API {
	return s.es
}

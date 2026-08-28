package crud

import (
	"context"
	"errors"

	"github.com/1303-yzym/MoonshotWell/pkg/config"
	"github.com/1303-yzym/MoonshotWell/pkg/contract"
	"github.com/1303-yzym/MoonshotWell/pkg/infra/logger"
	"github.com/1303-yzym/MoonshotWell/pkg/udt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type CRUD[V schema.Tabler] interface {
	WithDB(tx *gorm.DB) CRUD[V]
	G() gorm.Interface[V]
	Add(ctx context.Context, add *V) (err error)
	Upsert(ctx context.Context, entity *V, whereFuncs ...func(db *gorm.Statement)) (err error)
	QueryFirst(ctx context.Context, scopeFuncs ...func(db *gorm.Statement)) (v V, err error)
	QueryFirstNoErr(ctx context.Context, scopeFuncs ...func(db *gorm.Statement)) *V
	QueryById(ctx context.Context, id uint64, scopeFuncs ...func(db *gorm.Statement)) (v V, err error)
	QueryByIdNoErr(ctx context.Context, id uint64, scopeFuncs ...func(db *gorm.Statement)) *V
	QueryPagination(ctx context.Context, py udt.Pagination, scopeFuncs ...func(db *gorm.Statement)) (list []V, count int64, err error)
	QueryAllWithCount(ctx context.Context, scopeFuncs ...func(db *gorm.Statement)) (list []V, count int64, err error)
	Updates(ctx context.Context, whereFuncs func(db *gorm.Statement), values V) (err error)
	UpdatesById(ctx context.Context, id uint64, values V, scopeFuncs ...func(db *gorm.Statement)) (err error)
	Update(ctx context.Context, whereFuncs func(db *gorm.Statement), name string, value any) (err error)
	UpdateById(ctx context.Context, id uint64, name string, value any, scopeFuncs ...func(db *gorm.Statement)) (err error)
	Delete(ctx context.Context, whereFuncs ...func(db *gorm.Statement)) (err error)
	DeleteById(ctx context.Context, id uint64, unscoped bool) (err error)
}

type Crud[V schema.Tabler] struct {
	storage contract.Storage
	tx      *gorm.DB
}

func New[V schema.Tabler](storage contract.Storage, env config.Env) Crud[V] {
	// 根据环境开启自动迁移
	if env == config.EnvDev || env == config.EnvTest {
		var model V
		err := storage.DB().AutoMigrate(&model)
		if err != nil {
			logger.Error().Error("自动迁移失败", zap.String("table", model.TableName()), zap.Error(err))
		}
	}
	// 根据环境开启自动迁移
	return Crud[V]{storage: storage}
}

func (c Crud[V]) WithDB(tx *gorm.DB) CRUD[V] {
	return Crud[V]{
		storage: c.storage,
		tx:      tx,
	}
}

func (c Crud[V]) db() *gorm.DB {
	if c.tx != nil {
		return c.tx
	}

	return c.storage.DB()
}

// G https://github.com/go-gorm/gorm/pull/7424
func (c Crud[V]) G() gorm.Interface[V] {
	return gorm.G[V](c.db())
}

func (c Crud[V]) Add(ctx context.Context, add *V) (err error) {
	if err = c.G().Create(ctx, add); err != nil {
		return err
	}

	return
}

// TODO 存在并发问题，重写
func (c Crud[V]) Upsert(ctx context.Context, entity *V, whereFuncs ...func(db *gorm.Statement)) (err error) {
	// 先尝试查询是否存在
	_, err = c.G().Scopes(whereFuncs...).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在则创建
			return c.Add(ctx, entity)
		}
		// 其他错误直接返回
		return errors.New("no query rows")
	}
	// 存在则更新
	rowsAffected, err := c.G().Scopes(whereFuncs...).Updates(ctx, *entity)
	if err != nil {
		return errors.New("no affected rows")
	}

	if rowsAffected == 0 {
		return errors.New("no affected rows")
	}

	return nil
}

func (c Crud[V]) Delete(ctx context.Context, whereFuncs ...func(db *gorm.Statement)) (err error) {
	rowsAffected, err := c.G().Scopes(whereFuncs...).Delete(ctx)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		err = errors.New("no affected rows")
	}

	return
}

func (c Crud[V]) DeleteById(ctx context.Context, id uint64, unscoped bool) (err error) {
	rowsAffected, err := c.G().Scopes(func(db *gorm.Statement) {
		db.Unscoped = unscoped
	}).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		err = errors.New("no affected rows")
	}

	return
}

func (c Crud[V]) Updates(ctx context.Context, whereFuncs func(db *gorm.Statement), values V) (err error) {
	rowsAffected, err := c.G().Scopes(whereFuncs).Updates(ctx, values)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no affected rows")
	}

	return nil
}

func (c Crud[V]) Update(ctx context.Context, whereFuncs func(db *gorm.Statement), name string, value any) (err error) {
	rowsAffected, err := c.G().Scopes(whereFuncs).Update(ctx, name, value)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no affected rows")
	}

	return nil
}

func (c Crud[V]) UpdatesById(ctx context.Context, id uint64, values V, scopeFuncs ...func(db *gorm.Statement)) (err error) {
	rowsAffected, err := c.G().Scopes(scopeFuncs...).
		Where("id = ?", id).Updates(ctx, values)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no affected rows")
	}

	return nil
}

func (c Crud[V]) UpdateById(ctx context.Context, id uint64, name string, value any, scopeFuncs ...func(db *gorm.Statement)) (err error) {
	rowsAffected, err := c.G().Scopes(scopeFuncs...).
		Where("id = ?", id).Update(ctx, name, value)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no affected rows")
	}

	return nil
}

func (c Crud[V]) QueryFirst(ctx context.Context, scopeFuncs ...func(db *gorm.Statement)) (v V, err error) {
	v, err = c.G().Scopes(scopeFuncs...).First(ctx)
	if err != nil {
		return v, err
	}

	return v, nil
}

func (c Crud[V]) QueryFirstNoErr(ctx context.Context, scopeFuncs ...func(db *gorm.Statement)) *V {
	v, err := c.G().Scopes(scopeFuncs...).First(ctx)
	if err != nil {
		return nil
	}

	return &v
}

func (c Crud[V]) QueryById(ctx context.Context, id uint64, scopeFuncs ...func(db *gorm.Statement)) (v V, err error) {
	v, err = c.G().Scopes(scopeFuncs...).Scopes(func(db *gorm.Statement) {
		db.Where("id = ?", id)
	}).First(ctx)
	if err != nil {
		return v, errors.New("no query rows")
	}

	return v, nil
}

func (c Crud[V]) QueryByIdNoErr(ctx context.Context, id uint64, scopeFuncs ...func(db *gorm.Statement)) *V {
	v, err := c.G().Scopes(scopeFuncs...).Scopes(func(db *gorm.Statement) {
		db.Where("id = ?", id)
	}).First(ctx)
	if err != nil {
		return nil
	}

	return &v
}

func (c Crud[V]) QueryPagination(ctx context.Context, py udt.Pagination, scopeFuncs ...func(db *gorm.Statement)) (list []V, count int64, err error) {
	list, err = c.G().Scopes(scopeFuncs...).Limit(py.Size()).Offset(py.Offset()).Find(ctx)
	if err != nil {
		err = errors.New("no query rows")

		return
	}

	count, err = c.G().Scopes(scopeFuncs...).Count(ctx, "*")
	if err != nil {
		err = errors.New("no query rows")

		return
	}

	return list, count, nil
}

func (c Crud[V]) QueryAllWithCount(ctx context.Context, scopeFuncs ...func(db *gorm.Statement)) (list []V, count int64, err error) {
	list, err = c.G().Scopes(scopeFuncs...).Find(ctx)
	if err != nil {
		err = errors.New("no query rows")

		return
	}

	count, err = c.G().Scopes(scopeFuncs...).Count(ctx, "*")
	if err != nil {
		err = errors.New("no query rows")

		return
	}

	return list, count, nil
}

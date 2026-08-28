package db

import (
	"context"
	"time"

	"github.com/1303-yzym/MoonshotWell/pkg/contract"
	"github.com/1303-yzym/MoonshotWell/pkg/otel/zapgorm"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitMySqlDB(log *zap.Logger, cfg SQLConfig) (db *gorm.DB, err error) {
	// 驱动配置
	mysqlConfig := mysql.Config{
		DSN:                       cfg.DSN(),
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	return initGorm(log, mysql.New(mysqlConfig), cfg)
}

func initGorm(log *zap.Logger, dial gorm.Dialector, cfg SQLConfig) (db *gorm.DB, err error) {
	if db, err = gorm.Open(dial, &gorm.Config{
		CreateBatchSize: 1000,
		Logger:          dbLogger(log),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}); err != nil {
		return db, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}

	var connPool SQLConnPool

	if cfg.SQLConnPool != nil {
		connPool = *cfg.SQLConnPool
	} else {
		connPool = defaultSQLConnPool()
	}

	if err = connPool.Validation(); err != nil {
		return db, err
	}

	onnMaxIdleTime, err := time.ParseDuration(connPool.ConnMaxIdleTime)
	if err != nil {
		return db, err
	}

	connMaxLifetime, err := time.ParseDuration(connPool.ConnMaxLifetime)
	if err != nil {
		return db, err
	}

	sqlDB.SetMaxIdleConns(connPool.MaxIdleConn)
	sqlDB.SetMaxOpenConns(connPool.MaxOpenConn)
	sqlDB.SetConnMaxIdleTime(onnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	return db, nil
}

func dbLogger(log *zap.Logger) zapgorm.Logger {
	opts := []zapgorm.Options{
		zapgorm.WithLogLevel(gormlogger.Warn),
		zapgorm.WithContextFn(func(ctx context.Context) []zap.Field {
			if trace, ok := ctx.(contract.Trace); ok {
				return []zap.Field{zap.String("traceId", trace.TraceId())}
			} else {
				return nil
			}
		}),
		zapgorm.WithIgnoreRecordNotFoundError(true),
	}

	ll := zapgorm.New(log, opts...)
	ll.SetAsDefault()

	return ll
}

func defaultSQLConnPool() SQLConnPool {
	return SQLConnPool{
		MaxIdleConn:     5,
		MaxOpenConn:     20,
		ConnMaxIdleTime: "30m",
		ConnMaxLifetime: "2h",
	}
}

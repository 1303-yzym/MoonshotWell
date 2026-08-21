package db

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DBConfig struct {
	DBType string    `mapstructure:"db_type" json:"db_type" comment:"db类型 mysql | postgresql"`
	MySQL  SQLConfig `mapstructure:"mysql" json:"mysql" comment:"MySQL配置"`
	// Postgresql SQLConfig `mapstructure:"postgresql" json:"postgresql" comment:"MySQL配置"`
	// TODO 其他数据库库
}

type SQLConfig struct {
	Addr        string       `mapstructure:"addr" json:"addr" comment:"地址"`
	Config      string       `mapstructure:"config" json:"config" comment:"附加配置"`
	DBSchema    string       `mapstructure:"db_schema" json:"db_schema" comment:"架构"`
	Username    string       `mapstructure:"username" json:"username" comment:"用户"`
	Password    string       `mapstructure:"password" json:"password" comment:"密码"`
	SQLConnPool *SQLConnPool `mapstructure:"pool" json:"pool" comment:"连接池"`
}

type SQLConnPool struct {
	MaxIdleConn     int    `mapstructure:"max_idle_conn" json:"max_idle_conn" comment:"最大空闲控制器" mock:"5"`
	MaxOpenConn     int    `mapstructure:"max_open_conn" json:"max_open_conn" comment:"最大打开连接数" mock:"20"`
	ConnMaxIdleTime string `mapstructure:"conn_max_idle_time" json:"conn_max_idle_time" comment:"连接最大空闲时间" mock:"30m"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime" comment:"连接最大寿命" mock:"2h"`
}

func (s SQLConnPool) Validation() error {
	if !(s.MaxIdleConn >= 1 && s.MaxOpenConn >= 2 && s.ConnMaxIdleTime != "" && s.ConnMaxLifetime != "") {
		return errors.New("validation err")
	}

	return nil
}

func (c SQLConfig) DSN() string {
	return c.Username + ":" + c.Password + "@tcp(" + c.Addr + ")/" + c.DBSchema + "?" + c.Config
}

func (c SQLConfig) ShowDsn() string {
	return c.Username + ":" + "<mask>" + "@tcp(" + c.Addr + ")/" + c.DBSchema + "?" + c.Config
}

func InitDB(log *zap.Logger, dbCfg DBConfig) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	if dbCfg.DBType == "mysql" {
		db, err = InitMySqlDB(log, dbCfg.MySQL)
		if err != nil {
			zap.L().Panic("连接数据库失败", zap.Error(err))
		}

		zap.L().Sugar().Infof("Mysql [%s]", dbCfg.MySQL.ShowDsn())
	} else {
		zap.L().Panic("数据类型错误 db_type = mysql | postgresql", zap.String("db_type", dbCfg.DBType))
	}

	return db
}

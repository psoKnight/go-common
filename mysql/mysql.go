package mysqlx

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"log"
	"time"
)

type MySQL struct {
	cli *gorm.DB
	cfg *MySQLConfig
}

type MySQLConfig struct {
	Username      string        `json:"username"`       // 用户名
	Password      string        `json:"password"`       // 密码
	Address       string        `json:"address"`        // ip，示例：127.0.0.1:3306
	DatabaseName  string        `json:"database_name"`  // 数据库名称
	MaxOpenConns  int           `json:"max_open_conns"` // 最大连接数
	MaxIdleConns  int           `json:"max_idle_conns"` // 最大空闲连接数
	LogMode       string        `json:"log_mode"`       // 日志模式
	Logger        *log.Logger   `json:"logger"`         // 日志logger
	SlowThreshold time.Duration `json:"slow_threshold"` // 慢查询判定时间
}

// NewMySQL 新建mysql
func NewMySQL(cfg *MySQLConfig) (*MySQL, error) {
	if cfg == nil {
		return nil, errors.New("[mysql]config is nil")
	}

	logLevel := glogger.Silent

	switch cfg.LogMode {
	case "debug", "info":
		logLevel = glogger.Info
	case "warn":
		logLevel = glogger.Warn
	case "error":
		logLevel = glogger.Error
	default:
		logLevel = glogger.Info
	}

	newLogger := glogger.New(
		cfg.Logger,
		glogger.Config{
			SlowThreshold:             cfg.SlowThreshold, // 慢sql 阈值
			Colorful:                  false,             // 禁用彩色打印
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logLevel,
		},
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			cfg.Username, cfg.Password, cfg.Address, cfg.DatabaseName), // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用datetime 精度，MySQL v5.6之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL v5.7 之前的数据库和MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用change 重命名列，MySQL v8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前MySQL 版本自动配置
	}), &gorm.Config{
		Logger: newLogger,
		//SkipDefaultTransaction: true, // 是否默认开启事务：https://gorm.io/docs/transactions.html
	})
	if err != nil {
		return nil, err
	}

	// 获取通用数据库对象sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour * 4)

	return &MySQL{cli: db}, nil
}

// GetClient 获取连接
func (m *MySQL) GetClient() *gorm.DB {
	return m.cli
}

// Close 关闭连接
func (m *MySQL) Close() error {
	db, err := m.cli.DB()
	if err != nil {
		return err
	}
	if err = db.Close(); err != nil {
		return err
	}
	return nil
}

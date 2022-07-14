package mysqlx

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

func NewMySQL(cfg *Config) (*Mysql, error) {
	if cfg == nil {
		return nil, errors.New("[mysql]Config is nil.")
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

	return &Mysql{DB: db}, nil
}

func (d *Mysql) Close() {
	db, _ := d.DB.DB()
	if err := db.Close(); err != nil {
		d.config.Logger.Printf("[mysql]Close database err: %v.", err)
		return
	}
}

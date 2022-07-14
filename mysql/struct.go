package mysqlx

import (
	"gorm.io/gorm"
	"log"
	"time"
)

type Mysql struct {
	*gorm.DB
	config *Config
}

type Config struct {
	Username      string
	Password      string
	Address       string // 示例：127.0.0.1:3306
	DatabaseName  string
	MaxOpenConns  int
	MaxIdleConns  int
	LogMode       string
	Logger        *log.Logger
	SlowThreshold time.Duration
}

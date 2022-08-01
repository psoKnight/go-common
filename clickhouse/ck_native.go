package clickhouse

import (
	"context"
	"errors"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"time"
)

type ClickhouseConfig struct {
	Addrs              []string      `json:"endpoints"`            // clickhouse 地址
	Database           string        `json:"database"`             // 数据库名称
	Username           string        `json:"username"`             // 用户名是用于认证的用户名
	Password           string        `json:"password"`             // 密码是用于认证的密码
	InsecureSkipVerify bool          `json:"insecure_skip_verify"` // 是否建立TLS 建立安全连接
	DialTimeout        time.Duration `json:"dial_timeout"`         // 拨号心跳超时时间
	ConnMaxLifetime    time.Duration `json:"conn_max_lifetime"`    // 连接生存时长
	MaxExecutionTime   int           `json:"max_execution_time"`   // 最大执行时间
	Debug              bool          `json:"debug"`                // 是否开启debug
	MaxIdleConns       int           `json:"max_idle_conns"`       // 设置最大空闲连接数
	MaxOpenConns       int           `json:"max_open_conns"`       // 设置最大连接数
}

type Clickhouse struct {
	cli driver.Conn
	cfg *ClickhouseConfig // 基础配置
}

// NewClickhouse 获取cliskhouse 客户端
func NewClickhouse(cfg *ClickhouseConfig) (*Clickhouse, error) {

	if cfg == nil {
		return nil, errors.New("[ck]clickhouse config is nil")
	} else if len(cfg.Addrs) == 0 {
		return nil, errors.New("[ck]clickhouse addrs is 0")
	}

	client := &Clickhouse{
		cfg: cfg,
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: cfg.Addrs,
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		Debug:           cfg.Debug,
		DialTimeout:     cfg.DialTimeout,
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	client.cli = conn

	return client, nil
}

// GetClient 获取clickhouse client
func (c *Clickhouse) GetClient() driver.Conn {
	return c.cli
}

// ExecCQL 执行CQL
func (c *Clickhouse) ExecCQL(query string, ctx context.Context) error {
	if len(query) == 0 {
		return errors.New("[ck]query is '', please check")
	}
	return c.cli.Exec(ctx, query)
}

// SelectCQL 查询结果
func (c *Clickhouse) SelectCQL(query string, results interface{}, ctx context.Context) error {
	if len(query) == 0 {
		return errors.New("[ck]query is '', please check")
	}
	return c.cli.Select(ctx, results, query)
}

// GetCQLCount 获取CQL 结果总数
func (c *Clickhouse) GetCQLCount(query string, ctx context.Context) (uint64, error) {

	if len(query) == 0 {
		return 0, errors.New("[ck]query is '', please check")
	}

	var cnt uint64 = 0
	err := c.cli.QueryRow(ctx, query).Scan(&cnt)
	if err != nil {
		return 0, err
	}

	return cnt, nil
}

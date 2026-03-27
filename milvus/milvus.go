package milvusx

import (
	"context"
	"errors"
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"time"
)

// Milvus 封装
type Milvus struct {
	cli client.Client
	cfg *MilvusConfig
}

// MilvusConfig 配置
type MilvusConfig struct {
	Address  string `json:"address"`  // 127.0.0.1:19530
	Username string `json:"username"` // 可选
	Password string `json:"password"` // 可选
	Timeout  int    `json:"timeout"`  // 连接超时，单位：秒
	DbName   string `json:"db_name"`  // 默认集合
}

// NewMilvus 创建连接
func NewMilvus(cfg *MilvusConfig) (*Milvus, error) {
	if cfg == nil {
		return nil, errors.New("[milvus]config is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	conf := client.Config{
		Address:  cfg.Address,
		Username: cfg.Username,
		Password: cfg.Password,
		DBName:   cfg.DbName,
	}
	cli, err := client.NewClient(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("connect milvus failed: %w", err)
	}

	return &Milvus{
		cli: cli,
		cfg: cfg,
	}, nil
}

// GetClient 获取客户端
func (m *Milvus) GetClient() client.Client {
	return m.cli
}

// Close 关闭连接
func (m *Milvus) Close() error {
	return m.cli.Close()
}

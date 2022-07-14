package redis

import (
	"github.com/gomodule/redigo/redis"
)

type Hash struct {
	hashName string
	cli      *Redis
}

// NewHash 新建hash，hash key 为："HASH:"+key
func NewHash(hashName string, client *Redis) *Hash {
	return &Hash{hashName: "HASH:" + hashName, cli: client}
}

// HSET HSET 方法
func (h *Hash) HSET(field, value string) error {
	_, err := h.cli.ExecCommand("HSET", h.hashName, field, value)
	return err
}

// HGET HGET 方法
func (h *Hash) HGET(field string) (string, error) {
	return redis.String(h.cli.ExecCommand("HGET", h.hashName, field))
}

// HDEL HDEL 方法
func (h *Hash) HDEL(fields []string) error {
	if len(fields) == 0 {
		return nil
	}

	args := make([]interface{}, 0, len(fields)+1)
	args = append(args, h.hashName)

	for i := range fields {
		args = append(args, fields[i])
	}

	_, err := h.cli.ExecCommand("HDEL", args...)
	return err
}

// Clear 清除hash 表
func (h *Hash) Clear() error {
	_, err := h.cli.ExecCommand("DEL", h.hashName)
	return err
}

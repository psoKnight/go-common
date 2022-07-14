package redis

import (
	"github.com/gomodule/redigo/redis"
)

/**
获取Hash
	hash key 为："Hash:"+key
*/
func NewHash(key string, client *Redis) *Hash {
	return &Hash{name: "Hash:" + key, client: client}
}

/**
Hash HSET 方法
*/
func (h *Hash) HSet(field, value string) error {
	_, err := h.client.ExecCommand("HSET", h.name, field, value)

	return err
}

/**
Hash HGET 方法
*/
func (h *Hash) HGet(field string) (string, error) {
	return redis.String(h.client.ExecCommand("HGET", h.name, field))
}

/**
Hash HDEL 方法
*/
func (h *Hash) HDel(fields []string) error {
	if len(fields) == 0 {
		return nil
	}

	args := make([]interface{}, 0, len(fields)+1)
	args = append(args, h.name)

	for i := range fields {
		args = append(args, fields[i])
	}

	_, err := h.client.ExecCommand("HDEL", args...)
	return err
}

/**
Hash DEL 方法
*/
func (h *Hash) Del() error {
	_, err := h.client.ExecCommand("DEL", h.name)

	return err
}

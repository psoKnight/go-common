package redis

import (
	"testing"
	"time"
)

func TestHash(t *testing.T) {

	// 获取Redis
	var cfg = &RedisConfig{
		Password:       "yZY0G0Dzh5N",
		Address:        []string{"10.171.5.193:6382"},
		DatabaseId:     0,
		MaxIdle:        4,
		MaxActive:      64,
		IdleTimeout:    time.Duration(5000) * time.Millisecond,
		ConnectTimeout: time.Duration(5000) * time.Millisecond,
		ReadTimeout:    time.Duration(5000) * time.Millisecond,
		WriteTimeout:   time.Duration(180) * time.Second,
		IsCluster:      false,
	}

	redisClient, err := NewRedis(cfg)
	if err != nil {
		t.Errorf("Redis connect failed, err: %v.", err)
		return
	}

	// 获取Hash
	hashKey := "test_hash"
	hash := NewHash(hashKey, redisClient)

	// HSET 方法
	if err = hash.HSET("field_hash_1", "value_1"); err != nil {
		t.Errorf("Hash set err: %v.", err)
		return
	}

	// HGET 方法
	v, err := hash.HGET("field_hash_1")
	if err != nil {
		t.Errorf("Hash get err: %v.", err)
		return
	}
	t.Log(v)

	// HDEL 方法
	if err = hash.HDEL([]string{"field_hash_1"}); err != nil {
		t.Errorf("Hash del err: %v.", err)
		return
	}

	// 清空
	if err := hash.Clear(); err != nil {
		t.Errorf("Hash Clear err: %v.", err)
		return
	}
}

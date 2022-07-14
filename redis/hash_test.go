package redis

import (
	"testing"
	"time"
)

func TestHash(t *testing.T) {

	// 获取Redis
	var cfg = &Config{
		Password:       "3E*SWNf153D",
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

	redisClient, errNR := NewRedis(cfg)
	if errNR != nil {
		t.Errorf("Redis connect failed, err: %v.", errNR)
		return
	}

	// 获取Hash
	hashKey := "test_hash"
	hash := NewHash(hashKey, redisClient)

	// HSET 方法
	errS := hash.HSet("field_hash_1", "value_1")
	if errS != nil {
		t.Errorf("Hash set err: %v.", errS)
		return
	}

	// HGET 方法
	v, errG := hash.HGet("field_hash_1")
	if errG != nil {
		t.Errorf("Hash get err: %v.", errG)
		return
	}
	t.Log(v)

	// HDEL 方法
	errD := hash.HDel([]string{"field_hash_1"})
	if errD != nil {
		t.Errorf("Hash del err: %v.", errD)
		return
	}

	// DEL 方法
	errC := hash.Del()
	if errC != nil {
		t.Errorf("Hash Clear err: %v.", errC)
		return
	}

}

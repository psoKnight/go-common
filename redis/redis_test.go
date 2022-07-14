package redis

import (
	"testing"
	"time"
)

func TestRedis(t *testing.T) {

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

	// 关闭redis
	defer redisClient.ClosePool()

	// SET 方法
	set, err := redisClient.SET("field_1", "value_1", SetWithEX(5))
	if err != nil || set != "OK" {
		t.Errorf("Redis set err: %v, reply: %s.", err, set)
		return
	}
	t.Log(set)

	val, err := redisClient.GET("field_1")
	if err != nil {
		t.Errorf("Redis get err: %v.", err)
		return
	}
	t.Log(val)

	del, err := redisClient.DEL("field_1")
	if err != nil {
		t.Errorf("Redis del err: %v, reply: %d.", err, del)
		return
	}
	t.Log(del)
}

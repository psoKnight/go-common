package redis

import (
	"testing"
	"time"
)

func TestRedis(t *testing.T) {

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

	// 关闭redis
	defer redisClient.Close()

	// SET 方法
	set, errS := redisClient.SET("field_1", "value_1", SetWithEx(5))
	if errS != nil || set != "OK" {
		t.Errorf("Redis set err: %v, reply: %s.", errS, set)
		return
	}
	t.Log(set)

	val, errG := redisClient.GET("field_1")
	if errG != nil {
		t.Errorf("Redis get err: %v.", errG)
		return
	}
	t.Log(val)

	del, errD := redisClient.DEL("field_1")
	if errD != nil {
		t.Errorf("Redis del err: %v, reply: %d.", errD, del)
		return
	}
	t.Log(del)
}

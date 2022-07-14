package redis

import (
	"testing"
	"time"
)

func TestQuene(t *testing.T) {

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

	// 获取Quene
	queueKey := "test_queue"
	queue := NewPriorityQueue(queueKey, redisClient)

	// ZADD 方法
	if err = queue.ZADD("field_1", 1); err != nil {
		t.Errorf("Quene zadd field 1 err: %v.", err)
		return
	}

	if err = queue.ZADD("field_2", 2); err != nil {
		t.Errorf("Quene zadd field 2 err: %v.", err)
		return
	}

	if err = queue.ZADD("field_3", 3); err != nil {
		t.Errorf("Quene zadd field 3 err: %v.", err)
		return
	}

	if err = queue.ZADD("field_4", 4); err != nil {
		t.Errorf("Quene zadd field 3 err: %v.", err)
		return
	}

	// ZREM 方法
	if err = queue.ZREM([]string{"field_1"}); err != nil {
		t.Errorf("Queue zrem err: %v.", err)
		return
	}

	// ZSCORE 方法
	zscore, err := queue.ZSCORE("field_2")
	if err != nil {
		t.Errorf("Queue zscore err: %v.", err)
		return
	}
	t.Log(zscore)

	// ZRANGEBYSCORE 方法
	score, err := queue.ZRANGEBYSCORE(-1, 3)
	if err != nil {
		t.Errorf("Queue zrangebyscore err: %v.", err)
		return
	}
	t.Log(score)

	// 获取topk 小
	minScore, err := queue.TopMinScore(2)
	if err != nil {
		t.Errorf("Queue get top min err: %v.", err)
		return
	}
	t.Log(minScore)

	// 获取topk 大
	maxScore, err := queue.TopMaxScore(2)
	if err != nil {
		t.Errorf("Queue get top max err: %v.", err)
		return
	}
	t.Log(maxScore)

	// 清空
	if err = queue.Clear(); err != nil {
		t.Errorf("Queue del err: %v.", err)
		return
	}
}

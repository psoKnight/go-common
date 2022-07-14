package redis

import (
	"testing"
	"time"
)

func TestQuene(t *testing.T) {

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

	// 获取Quene
	queueKey := "test_queue"
	queue := NewPriorityQueue(queueKey, redisClient)

	// ZADD 方法
	errA1 := queue.ZADD("field_1", 1)
	if errA1 != nil {
		t.Errorf("Quene zadd field 1 err: %v.", errA1)
		return
	}

	errA2 := queue.ZADD("field_2", 2)
	if errA2 != nil {
		t.Errorf("Quene zadd field 2 err: %v.", errA2)
		return
	}

	errA3 := queue.ZADD("field_3", 3)
	if errA3 != nil {
		t.Errorf("Quene zadd field 3 err: %v.", errA3)
		return
	}

	errA4 := queue.ZADD("field_4", 4)
	if errA4 != nil {
		t.Errorf("Quene zadd field 3 err: %v.", errA4)
		return
	}

	// ZREM 方法
	errR := queue.ZREM([]string{"field_1"})
	if errR != nil {
		t.Errorf("Queue zrem err: %v.", errR)
		return
	}

	// ZSCORE 方法
	zscore, errS := queue.ZSCORE("field_2")
	if errS != nil {
		t.Errorf("Queue zscore err: %v.", errS)
		return
	}
	t.Log(zscore)

	// ZRANGEBYSCORE 方法
	score, errRS := queue.ZRANGEBYSCORE(-1, 3)
	if errRS != nil {
		t.Errorf("Queue zrangebyscore err: %v.", errRS)
		return
	}
	t.Log(score)

	// 获取topk 小
	minScore, errMin := queue.TopMinScore(2)
	if errMin != nil {
		t.Errorf("Queue get top min err: %v.", errMin)
		return
	}
	t.Log(minScore)

	// 获取topk 大
	maxScore, errMax := queue.TopMaxScore(2)
	if errMax != nil {
		t.Errorf("Queue get top max err: %v.", errMax)
		return
	}
	t.Log(maxScore)

	// DEL 方法
	errD := queue.Del()
	if errD != nil {
		t.Errorf("Queue del err: %v.", errD)
		return
	}

}

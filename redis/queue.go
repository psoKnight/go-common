package redis

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type PriorityQueue struct {
	queueName string
	cli       *Redis
}

type ZSetData struct {
	Value string
	Score float64
}

// NewPriorityQueue 新建队列
func NewPriorityQueue(queueName string, redisClient *Redis) *PriorityQueue {
	return &PriorityQueue{queueName: "QUEUE:" + queueName, cli: redisClient}
}

// ZADD ZADD 方法
func (pq *PriorityQueue) ZADD(elem string, scores int64) error {
	_, err := pq.cli.ExecCommand("ZADD", pq.queueName, scores, elem)
	if err != nil {
		return err
	}

	return nil
}

// ZSCORE ZSCORE 方法
func (pq *PriorityQueue) ZSCORE(elem string) (float64, error) {
	score, err := pq.cli.ExecCommand("ZSCORE", pq.queueName, elem)
	if err != nil {
		return 0, err
	}
	if score != nil {
		scoreString := string(score.([]byte))
		return strconv.ParseFloat(scoreString, 64)
	}
	return 0, nil
}

// ZREM ZREM 方法
func (pq *PriorityQueue) ZREM(elems []string) error {
	interfaceSlice := make([]interface{}, 0, len(elems)+1)
	interfaceSlice = append(interfaceSlice, pq.queueName)

	for _, elem := range elems {
		interfaceSlice = append(interfaceSlice, elem)
	}

	_, err := pq.cli.ExecCommand("ZREM", interfaceSlice...)
	if err != nil {
		return err
	}

	return nil
}

// ZRANGEBYSCORE ZRANGEBYSCORE 方法
func (pq *PriorityQueue) ZRANGEBYSCORE(minScores, maxScores int64) ([]ZSetData, error) {
	datas, err := redis.Strings(pq.cli.ExecCommand("ZRANGEBYSCORE", pq.queueName, minScores, maxScores, "WITHSCORES"))
	if err != nil {
		return make([]ZSetData, 0), err
	}

	dataSet := make([]ZSetData, 0, len(datas)/2)
	for i := 0; i < len(datas); i += 2 {
		score, err := strconv.ParseFloat(datas[i+1], 64)
		if err != nil {
			return make([]ZSetData, 0), err
		}
		dataSet = append(dataSet, ZSetData{
			Value: datas[i],
			Score: score,
		})
	}

	return dataSet, nil
}

// TopMinScore 根据score 获取topk 最小值所对应的数据
func (pq *PriorityQueue) TopMinScore(num int) ([]ZSetData, error) {
	return pq.topMinOrMaxScore(num, 0)
}

// TopMaxScore 根据score 获取topk 最大值所对应的数据
func (pq *PriorityQueue) TopMaxScore(num int) ([]ZSetData, error) {
	return pq.topMinOrMaxScore(num, 1)
}

// Clear 删除 quene
func (pq *PriorityQueue) Clear() error {
	_, err := pq.cli.ExecCommand("DEL", pq.queueName)
	if err != nil {
		return err
	}
	return nil
}

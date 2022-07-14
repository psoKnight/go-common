package redis

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
)

/**
获取Queue
*/
func NewPriorityQueue(queueName string, redisClient *Redis) *PriorityQueue {
	return &PriorityQueue{queueName: "Queue:" + queueName, redisClient: redisClient}
}

/**
Queue ZADD 方法
*/
func (pq *PriorityQueue) ZADD(elem string, scores int64) error {
	_, err := pq.redisClient.ExecCommand("ZADD", pq.queueName, scores, elem)
	if err != nil {
		return err
	}

	return nil
}

/**
Queue ZSCORE 方法
*/
func (pq *PriorityQueue) ZSCORE(elem string) (float64, error) {
	score, err := pq.redisClient.ExecCommand("ZSCORE", pq.queueName, elem)
	if err != nil {
		return 0, err
	}
	if score != nil {
		scoreString := string(score.([]byte))
		return strconv.ParseFloat(scoreString, 64)
	}
	return 0, nil
}

/**
Queue ZREM 方法
*/
func (pq *PriorityQueue) ZREM(elems []string) error {
	interfaceSlice := make([]interface{}, 0, len(elems)+1)
	interfaceSlice = append(interfaceSlice, pq.queueName)

	for _, elem := range elems {
		interfaceSlice = append(interfaceSlice, elem)
	}

	_, err := pq.redisClient.ExecCommand("ZREM", interfaceSlice...)
	if err != nil {
		return err
	}

	return nil
}

/**
Queue ZRANGEBYSCORE 方法
*/
func (pq *PriorityQueue) ZRANGEBYSCORE(minScores, maxScores int64) ([]ZSetData, error) {
	datas, err := redis.Strings(pq.redisClient.ExecCommand("ZRANGEBYSCORE", pq.queueName, minScores, maxScores, "WITHSCORES"))
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

/**
Quene 根据score 获取topk 最小值所对应的数据
*/

func (pq *PriorityQueue) TopMinScore(num int) ([]ZSetData, error) {
	return pq.topMinOrMaxScore(num, 0)
}

/**
Quene 根据score 获取topk 最大值所对应的数据
*/
func (pq *PriorityQueue) TopMaxScore(num int) ([]ZSetData, error) {
	return pq.topMinOrMaxScore(num, 1)
}

/**
Quene 删除
*/
func (pq *PriorityQueue) Del() error {
	_, err := pq.redisClient.ExecCommand("DEL", pq.queueName)
	if err != nil {
		return err
	}

	return nil
}

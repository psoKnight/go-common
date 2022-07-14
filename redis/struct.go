package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

// Set 参数列表
type setArgs []interface{}

// Set 额外参数
type SetOption func(args setArgs) setArgs

type Redis struct {
	redisPool *redis.Pool
}

// Redis
type Config struct {
	Password       string
	Address        []string
	DatabaseId     int
	MaxIdle        int
	MaxActive      int
	IdleTimeout    time.Duration
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration

	IsCluster bool
}

// Hash
type Hash struct {
	name   string
	client *Redis
}

// Queue
type PriorityQueue struct {
	queueName   string
	redisClient *Redis
}

type ZSetData struct {
	Value string
	Score float64
}

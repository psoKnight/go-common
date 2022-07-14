package redis

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
)

/**
设置指定的过期时间，以秒为单位。
>= 2.6.12
*/
func SetWithEx(seconds int) SetOption {
	return func(args setArgs) setArgs {
		args = append(args, "EX", seconds)
		return args
	}
}

/**
设置指定的过期时间，以毫秒为单位。
>= 2.6.12
*/
func SetWithPx(milliseconds int) SetOption {
	return func(args setArgs) setArgs {
		args = append(args, "PX", milliseconds)
		return args
	}
}

/**
以秒为单位设置指定的 Unix 时间，密钥将在该时间到期。
>= 6.2
*/
func SetWithExat(timestampSeconds int) SetOption {
	return func(args setArgs) setArgs {
		args = append(args, "EXAT", timestampSeconds)
		return args
	}
}

/**
设置指定的 Unix 时间，密钥将在该时间到期，以毫秒为单位。
>= 6.2
*/
func SetWithPxat(timestampMilliseconds int) SetOption {
	return func(args setArgs) setArgs {
		args = append(args, "PXAT", timestampMilliseconds)
		return args
	}
}

/**
只设置键，如果它不存在。
>= 2.6.12
*/
func SetWithNx() SetOption {
	return func(args setArgs) setArgs {
		args = append(args, "NX")
		return args
	}
}

/**
仅设置已存在的密钥。
>= 2.6.12
*/
func SetWithXx() SetOption {
	return func(args setArgs) setArgs {
		args = append(args, "XX")
		return args
	}
}

/**
保留与密钥关联的生存时间。
>= 6.0
*/
func SetWithKeepTTL() SetOption {
	return func(args setArgs) setArgs {
		args = append(args, "KEEPTTL")
		return args
	}
}

/**
返回存储在 key 中的旧值，或者当 key 不存在时返回 nil。
>= 6.2
*/
func SetWithGet() SetOption {
	return func(args setArgs) setArgs {
		args = append(args, "GET")
		return args
	}
}

/**
Quene 根据score 获取topk 最小或者最大值所对应的数据
	flag：
		0: 最小/升序
		1: 最大/降序
*/
func (r *PriorityQueue) topMinOrMaxScore(num, flag int) ([]ZSetData, error) {
	var err error
	var datas []string
	if flag == 0 {
		datas, err = redis.Strings(r.redisClient.ExecCommand("ZRANGE", r.queueName, 0, num-1, "WITHSCORES"))
	} else {
		datas, err = redis.Strings(r.redisClient.ExecCommand("ZREVRANGE", r.queueName, 0, num-1, "WITHSCORES"))
	}

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

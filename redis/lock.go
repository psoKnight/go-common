package redis

import "github.com/gomodule/redigo/redis"

var deleteAndRPUSHScript = redis.NewScript(2, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		redis.call("DEL", KEYS[1])
		if redis.call("LLEN", KEYS[2]) == 0 then
			redis.call("RPUSH", KEYS[2], 1)
			redis.call("EXPIRE", KEYS[2], 5)
			return 1
		else
			return 0
		end
	else
		return 0
	end
`)

// TryGetLock 扩展NX 实现锁，该锁是非阻塞的
func (r *Redis) TryGetLock(key, value string, expireTimeSeconds int64) (bool, error) {
	ok, err := r.SET(key, value, SetWithEX(int(expireTimeSeconds)), SetWithNX())
	if ok == "OK" && err == nil {
		return true, nil
	}
	return false, err
}

// WaitForGetLock 等待获取锁
func (r *Redis) WaitForGetLock(waitKey string, expireTimeSeconds int64) (bool, error) {
	_, err := r.BRPOP(waitKey, expireTimeSeconds)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ReleaseLockAndRpush 释放锁并且重新添加
func (r *Redis) ReleaseLockAndRpush(key, waitKey, value string) error {
	rp := r.redisPool.Get()
	defer rp.Close()
	_, err := deleteAndRPUSHScript.Do(rp, key, waitKey, value)
	return err
}

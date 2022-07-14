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

// TODO 后续补充

/**
SET 扩展NX 实现锁，该锁是非阻塞的
*/
func (c *Redis) TryGetLock(key, value string, expireTimeSeconds int64) (bool, error) {
	ok, err := c.SET(key, value, SetWithEx(int(expireTimeSeconds)), SetWithNx())
	if ok == "OK" && err == nil {
		return true, nil
	}
	return false, err
}

func (c *Redis) WaitForGetLock(waitKey string, expireTimeSeconds int64) (bool, error) {
	_, err := c.BRPOP(waitKey, expireTimeSeconds)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *Redis) ReleaseLockAndRpush(key, waitKey, value string) error {
	rp := c.redisPool.Get()
	defer rp.Close()

	_, err := deleteAndRPUSHScript.Do(rp, key, waitKey, value)
	return err
}

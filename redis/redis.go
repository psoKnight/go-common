package redis

import (
	"errors"
	"github.com/FZambia/sentinel"
	"github.com/gomodule/redigo/redis"
	"github.com/vmihailenco/msgpack/v4"
	"reflect"
)

/**
获取Redis 客户端
*/
func NewRedis(cfg *Config) (*Redis, error) {
	if cfg == nil {
		return nil, errors.New("[redis]Cfg is nil.")
	}

	redisPool := &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		IdleTimeout: cfg.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := wrapDial(cfg)
			if err != nil {
				return nil, err
			}

			if _, err := c.Do("AUTH", cfg.Password); err != nil {
				c.Close()
				return nil, err
			}
			_, err = c.Do("SELECT", cfg.DatabaseId)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	return &Redis{redisPool: redisPool}, nil
}

/**
Redis dial
*/
func wrapDial(cfg *Config) (redis.Conn, error) {
	redisAddr := cfg.Address[0]
	if cfg.IsCluster {
		sntnl := &sentinel.Sentinel{
			Addrs:      cfg.Address,
			MasterName: "mymaster",
			Dial: func(addr string) (redis.Conn, error) {
				// dial 是每次轮询的时候查询的dial，所以是否需要添加后续的timeout 暂未可知，后续再整理
				c, err := redis.Dial("tcp", addr, redis.DialConnectTimeout(cfg.ConnectTimeout),
					redis.DialReadTimeout(cfg.ReadTimeout),
					redis.DialWriteTimeout(cfg.WriteTimeout))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		}
		var err error
		redisAddr, err = sntnl.MasterAddr()
		if err != nil {
			return nil, err
		}
	}

	return redis.Dial("tcp", redisAddr, redis.DialConnectTimeout(cfg.ConnectTimeout),
		redis.DialReadTimeout(cfg.ReadTimeout),
		redis.DialWriteTimeout(cfg.WriteTimeout))
}

/**
Redis 关闭客户端
*/
func (r *Redis) Close() {
	if r.redisPool != nil {
		r.redisPool.Close()
	}
}

/**
Redis 执行command 命令
*/
func (r *Redis) ExecCommand(command string, args ...interface{}) (interface{}, error) {
	rc := r.redisPool.Get()
	defer rc.Close()

	return rc.Do(command, args...)
}

/**
Redis 存储struct，操作带有过期时间（毫秒）
*/
func (r *Redis) SetExpWithMP(key string, value interface{}, expireMilliseconds int) error {
	v := reflect.ValueOf(value)
	k := v.Kind()
	if k != reflect.Ptr {
		return errors.New("[redis]Value must be a pointer to a struct.")
	}
	v = v.Elem()
	k = v.Kind()
	if k != reflect.Struct {
		return errors.New("[redis]Value must be a pointer to a struct.")
	}

	b, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	_, err = r.ExecCommand("SET", key, b, "PX", expireMilliseconds)
	return err
}

/**
Redis 根据key 获取struct
*/
func (c *Redis) GetWithMP(key string, value interface{}) error {
	v := reflect.ValueOf(value)
	k := v.Kind()
	if k != reflect.Ptr {
		return errors.New("[redis]Value must be a pointer to a struct.")
	}
	v = v.Elem()
	k = v.Kind()
	if k != reflect.Struct {
		return errors.New("[redis]Value must be a pointer to a struct.")
	}

	b, err := redis.Bytes(c.ExecCommand("GET", key))
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(b, value)
}

/**
Redis MGET 方法
*/
func (c *Redis) MGet(keys []string) ([][]byte, error) {
	cfgs := make([]interface{}, len(keys))
	for i := range keys {
		cfgs[i] = keys[i]
	}

	return redis.ByteSlices(c.ExecCommand("MGET", cfgs...))
}

/**
Redis 获取连接
*/
func (r *Redis) GetCon() redis.Conn {
	return r.redisPool.Get()
}

/**
Redis 关闭从pool 获取的连接
*/
func (r *Redis) CloseCon(rc redis.Conn) {
	rc.Close()
}

/**
Redis SET 方法

SET 设置键来保存字符串值。
如果key 已经包含一个值，则无论其类型如何，它都会被覆盖。
成功的SET 操作将丢弃与密钥关联的任何先前的生存时间。
	[NOTE] 由于SET 命令选项可以替换SETNX、SETEX、PSETEX、GETSET，因此在未来的Redis 版本中，这些命令可能会被弃用并最终被删除。
*/
func (c *Redis) SET(key string, value interface{}, options ...SetOption) (string, error) {
	args := setArgs{key, value}
	for _, f := range options {
		args = f(args)
	}

	return redis.String(c.ExecCommand("SET", args...))
}

/**
Redis DEL 方法

DEL 删除指定的键。如果键不存在，则忽略该键。
*/
func (c *Redis) DEL(keys ...string) (int, error) {
	ikeys := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		ikeys = append(ikeys, key)
	}
	return redis.Int(c.ExecCommand("DEL", ikeys...))
}

/**
Redis EXISTS 方法

EXISTS 如果键存在则返回。
*/
func (c *Redis) EXISTS(key string) (int, error) {
	return redis.Int(c.ExecCommand("EXISTS", key))
}

/**
Redis PEXPIRE 方法

PEXPIRE 此命令的工作方式与EXPIRE 完全相同，但密钥的生存时间以毫秒而不是秒为单位指定。
*/
func (c *Redis) PEXPIRE(key string, milliseconds int) (int, error) {
	return redis.Int(c.ExecCommand("PEXPIRE", key, milliseconds))
}

/**
Redis EXPIRE 方法

EXPIRE 此命令设置密钥的生存时间，以秒为单位。
*/
func (c *Redis) EXPIRE(key string, seconds int) (int, error) {
	return redis.Int(c.ExecCommand("EXPIRE", key, seconds))
}

/**
Redis PEXPIREAT 方法

PEXPIREAT 与EXPIREAT 具有相同的效果和语义，但密钥到期的Unix 时间是以毫秒而不是秒指定的。
*/
func (c *Redis) PEXPIREAT(key string, millisecondsTimestamp int64) (int, error) {
	return redis.Int(c.ExecCommand("PEXPIREAT", key, millisecondsTimestamp))
}

/**
Redis PTTL 方法

PTTL 与TTL 一样，此命令返回设置了过期时间的密钥的剩余生存时间，唯一的区别是TTL 以秒为单位返回剩余时间量，而PTTL 以毫秒为单位返回剩余时间。
如果密钥不存在，该命令将返回-2。如果密钥存在但没有关联的过期，则该命令返回-1。
*/
func (c *Redis) PTTL(key string) (int, error) {
	return redis.Int(c.ExecCommand("PTTL", key))
}

/**
Redis GET 方法

GET 获取键的值。如果键不存在，则返回特殊值nil。如果key 中存储的值不是字符串，则会返回错误，因为GET 仅处理字符串值。
*/
func (c *Redis) GET(key string) (string, error) {
	return redis.String(c.ExecCommand("GET", key))
}

/**
Redis KEYS 方法

KEYS 返回所有匹配模式的键。
[Warning]：将KEYS 视为仅应在生产环境中极其小心地使用的命令。
当它针对大型数据库执行时，它可能会破坏性能。
此命令用于调试和特殊操作，例如更改键空间布局。
不要在常规应用程序代码中使用KEYS。
如果您正在寻找一种在键空间子集中查找键的方法，请考虑使用SCAN 或集合。
支持的glob 样式模式：
     h?llo matches hello, hallo and hxllo
     h*llo matches hllo and heeeello
     h[ae]llo matches hello and hallo, but not hillo
     h[^e]llo matches hallo, hbllo, ... but not hello
     h[a-b]llo matches hallo and hbllo
 Use \ to escape special characters if you want to match them verbatim.
*/
func (c *Redis) KEYS(key string) ([]string, error) {
	return redis.Strings(c.ExecCommand("KEYS", key))
}

/**
Redis BPOP 方法

自 2.0.0 起可用。
时间复杂度：O(N)，其中 N 是提供的键的数量。
BRPOP 是一个阻止列表弹出原语。它是RPOP 的阻塞版本，因为当没有要从任何给定列表中弹出的元素时，它会阻塞连接。
*/
func (c *Redis) BRPOP(key string, expireTimeSeconds int64) ([]string, error) {
	return redis.Strings(c.ExecCommand("BRPOP", key, expireTimeSeconds))
}

/**
Redis LPUSH 方法
*/
func (c *Redis) LPUSH(key, value string) (int, error) {
	return redis.Int(c.ExecCommand("LPUSH", key, value))
}

/**
Redis RPOP 方法
*/
func (c *Redis) RPOP(key string) (string, error) {
	return redis.String(c.ExecCommand("RPOP", key))
}

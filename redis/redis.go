package redis

import (
	"errors"
	"github.com/FZambia/sentinel"
	"github.com/gomodule/redigo/redis"
	"github.com/vmihailenco/msgpack/v4"
	"reflect"
	"time"
)

type RedisConfig struct {
	Password       string        `json:"password"`
	Address        []string      `json:"address"`
	DatabaseId     int           `json:"database_id"`
	MaxIdle        int           `json:"max_idle"`
	MaxActive      int           `json:"max_active"`
	IdleTimeout    time.Duration `json:"idle_timeout"`
	ConnectTimeout time.Duration `json:"connect_timeout"`
	ReadTimeout    time.Duration `json:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout"`

	IsCluster bool `json:"is_cluster"`
}

type Redis struct {
	redisPool *redis.Pool
}

// NewRedis 新建redis
func NewRedis(cfg *RedisConfig) (*Redis, error) {
	if cfg == nil {
		return nil, errors.New("[redis]cfg is nil")
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

// Redis dial
func wrapDial(cfg *RedisConfig) (redis.Conn, error) {
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

// GetClient 获取client
func (r *Redis) GetClient() redis.Conn {
	return r.redisPool.Get()
}

// CloseClient 关闭从pool 获取的client
func (r *Redis) CloseClient(rc redis.Conn) error {
	if err := rc.Close(); err != nil {
		return err
	}
	return nil
}

// GetPool 获取pool
func (r *Redis) GetPool() *redis.Pool {
	return r.redisPool
}

// ClosePool 关闭pool
func (r *Redis) ClosePool() {
	if r.redisPool != nil {
		r.redisPool.Close()
	}
}

// ExecCommand 执行command 命令
func (r *Redis) ExecCommand(command string, args ...interface{}) (interface{}, error) {
	rc := r.redisPool.Get()
	defer rc.Close()

	return rc.Do(command, args...)
}

// SetExpWithMP 存储struct，操作带有过期时间（毫秒）
func (r *Redis) SetExpWithMP(key string, value interface{}, expireMilliseconds int) error {
	v := reflect.ValueOf(value)
	k := v.Kind()
	if k != reflect.Ptr {
		return errors.New("[redis]value must be a pointer to a struct")
	}
	v = v.Elem()
	k = v.Kind()
	if k != reflect.Struct {
		return errors.New("[redis]value must be a pointer to a struct")
	}

	b, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	_, err = r.ExecCommand("SET", key, b, "PX", expireMilliseconds)
	return err
}

// GetWithMP 根据key 获取struct
func (r *Redis) GetWithMP(key string, value interface{}) error {
	v := reflect.ValueOf(value)
	k := v.Kind()
	if k != reflect.Ptr {
		return errors.New("[redis]value must be a pointer to a struct")
	}
	v = v.Elem()
	k = v.Kind()
	if k != reflect.Struct {
		return errors.New("[redis]value must be a pointer to a struct")
	}

	b, err := redis.Bytes(r.ExecCommand("GET", key))
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(b, value)
}

// MGET 方法
func (r *Redis) MGET(keys []string) ([][]byte, error) {
	cfgs := make([]interface{}, len(keys))
	for i := range keys {
		cfgs[i] = keys[i]
	}

	return redis.ByteSlices(r.ExecCommand("MGET", cfgs...))
}

// SET 方法
/**
SET 设置键来保存字符串值。
如果key 已经包含一个值，则无论其类型如何，它都会被覆盖。
成功的SET 操作将丢弃与密钥关联的任何先前的生存时间。
	[NOTE] 由于SET 命令选项可以替换SETNX、SETEX、PSETEX、GETSET，因此在未来的Redis 版本中，这些命令可能会被弃用并最终被删除。
*/
func (r *Redis) SET(key string, value interface{}, options ...SetOption) (string, error) {
	args := setArgs{key, value}
	for _, f := range options {
		args = f(args)
	}

	return redis.String(r.ExecCommand("SET", args...))
}

// DEL 方法
/**
DEL 删除指定的键。如果键不存在，则忽略该键。
*/
func (r *Redis) DEL(keys ...string) (int, error) {
	ikeys := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		ikeys = append(ikeys, key)
	}
	return redis.Int(r.ExecCommand("DEL", ikeys...))
}

// EXISTS 方法
/**
EXISTS 如果键存在则返回。
*/
func (r *Redis) EXISTS(key string) (int, error) {
	return redis.Int(r.ExecCommand("EXISTS", key))
}

// PEXPIRE 方法
/**
PEXPIRE 此命令的工作方式与EXPIRE 完全相同，但密钥的生存时间以毫秒而不是秒为单位指定。
*/
func (r *Redis) PEXPIRE(key string, milliseconds int) (int, error) {
	return redis.Int(r.ExecCommand("PEXPIRE", key, milliseconds))
}

// EXPIRE 方法
/**
EXPIRE 此命令设置密钥的生存时间，以秒为单位。
*/
func (r *Redis) EXPIRE(key string, seconds int) (int, error) {
	return redis.Int(r.ExecCommand("EXPIRE", key, seconds))
}

// PEXPIREAT 方法
/**
PEXPIREAT 与EXPIREAT 具有相同的效果和语义，但密钥到期的Unix 时间是以毫秒而不是秒指定的。
*/
func (r *Redis) PEXPIREAT(key string, millisecondsTimestamp int64) (int, error) {
	return redis.Int(r.ExecCommand("PEXPIREAT", key, millisecondsTimestamp))
}

// PTTL 方法
/**
PTTL 与TTL 一样，此命令返回设置了过期时间的密钥的剩余生存时间，唯一的区别是TTL 以秒为单位返回剩余时间量，而PTTL 以毫秒为单位返回剩余时间。
如果密钥不存在，该命令将返回-2。如果密钥存在但没有关联的过期，则该命令返回-1。
*/
func (r *Redis) PTTL(key string) (int, error) {
	return redis.Int(r.ExecCommand("PTTL", key))
}

// GET 方法
/**
GET 获取键的值。如果键不存在，则返回特殊值nil。如果key 中存储的值不是字符串，则会返回错误，因为GET 仅处理字符串值。
*/
func (r *Redis) GET(key string) (string, error) {
	return redis.String(r.ExecCommand("GET", key))
}

// KEYS 方法
/**
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
func (r *Redis) KEYS(key string) ([]string, error) {
	return redis.Strings(r.ExecCommand("KEYS", key))
}

// BRPOP 方法
/**
自 2.0.0 起可用。
时间复杂度：O(N)，其中 N 是提供的键的数量。
BRPOP 是一个阻止列表弹出原语。它是RPOP 的阻塞版本，因为当没有要从任何给定列表中弹出的元素时，它会阻塞连接。
*/
func (r *Redis) BRPOP(key string, expireTimeSeconds int64) ([]string, error) {
	return redis.Strings(r.ExecCommand("BRPOP", key, expireTimeSeconds))
}

// LPUSH 方法
func (r *Redis) LPUSH(key, value string) (int, error) {
	return redis.Int(r.ExecCommand("LPUSH", key, value))
}

// RPOP 方法
func (r *Redis) RPOP(key string) (string, error) {
	return redis.String(r.ExecCommand("RPOP", key))
}

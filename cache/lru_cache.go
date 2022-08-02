package cache

import (
	"github.com/bluele/gcache"
	"github.com/sirupsen/logrus"
	"time"
)

type Cache struct {
	cache      gcache.Cache  // LRU 缓存
	capacity   int           // Cache 最大缓冲容量
	expiration time.Duration // Cache 缓存过期时间
}

// InitGCache 初始化缓存区
func InitGCache(capacity int, expiration time.Duration) *Cache {
	c := &Cache{}

	c.cache = gcache.New(capacity).LRU().Expiration(expiration).Build()
	c.capacity = capacity
	c.expiration = expiration
	return c
}

// GetClient 获取cache client
func (c *Cache) GetClient() gcache.Cache {
	return c.cache
}

// Get 根据key 获取缓存数据
func (c *Cache) Get(cacheKey interface{}) interface{} {
	cache, err := c.cache.Get(cacheKey)
	if err == nil {
		return cache
	}
	return nil
}

// GetIfExist 根据key 获取缓存数据(如果存在key)
func (c *Cache) GetIfExist(cacheKey interface{}) interface{} {
	present, err := c.cache.GetIFPresent(cacheKey)
	if err == nil {
		return present
	}
	return nil
}

// GetAll 获取全部缓存
/**
checkExpired：是否校验过期key
*/
func (c *Cache) GetAll(checkExpired bool) map[interface{}]interface{} {
	return c.cache.GetALL(checkExpired)
}

// Set 根据key 设置缓存数据
func (c *Cache) Set(cacheKey interface{}, cacheData interface{}) error {
	return c.cache.Set(cacheKey, cacheData)
}

// SetWithExpire 根据key 设置缓存数据（存在过期时间）
func (c *Cache) SetWithExpire(cacheKey interface{}, cacheData interface{}, cacheTime time.Duration) error {
	return c.cache.SetWithExpire(cacheKey, cacheData, cacheTime)
}

// Remove 根据key 删除缓存
func (c *Cache) Remove(cacheKey interface{}) {
	c.cache.Remove(cacheKey)
	logrus.Infof("[gcache]'%v' removed.", cacheKey)
	return
}

// Purge 清除全部缓存
func (c *Cache) Purge() {
	c.cache.Purge()
	logrus.Info("[gcache]purged.")
}

// Keys 获取全部keys
/**
checkExpired：是否校验过期key
*/
func (c *Cache) Keys(checkExpired bool) []interface{} {
	return c.cache.Keys(checkExpired)
}

// Len 获取缓存的长度
/**
checkExpired：是否校验过期key
*/
func (c *Cache) Len(checkExpired bool) int {
	return c.cache.Len(checkExpired)
}

// Has 获取缓存是否存在key
func (c *Cache) Has(cacheKey interface{}) bool {
	return c.cache.Has(cacheKey)
}

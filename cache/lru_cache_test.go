package cache

import (
	"testing"
	"time"
)

func TestGetGcache(t *testing.T) {
	cache := InitGCache(3, time.Duration(3)*time.Second)
	t.Log(cache)
}

func TestSet(t *testing.T) {
	cache := InitGCache(3, time.Duration(3)*time.Second)

	err := cache.Set("key_1", "value_1")
	if err != nil {
		t.Errorf("Cache set err: %v.", err)
		return
	}
}

func TestGet(t *testing.T) {
	cache := InitGCache(3, time.Duration(3)*time.Second)

	err := cache.Set("key_2", "value_2")
	if err != nil {
		t.Errorf("Cache set err: %v.", err)
		return
	}

	value := cache.Get("key_2")
	t.Log(value)

	time.Sleep(time.Duration(3) * time.Second)
	valueAfterSleep := cache.Get("key_2")
	t.Log(valueAfterSleep)
}

func TestSetWithExpire(t *testing.T) {
	cache := InitGCache(3, time.Duration(3)*time.Second)

	err := cache.SetWithExpire("key_3", "value_3", time.Duration(5)*time.Second)
	if err != nil {
		t.Errorf("Cache set with expire err: %v.", err)
		return
	}

	value := cache.Get("key_3")
	t.Log(value)

	time.Sleep(time.Duration(5) * time.Second)
	valueAfterSleep := cache.Get("key_3")
	t.Log(valueAfterSleep)
}

func TestPurge(t *testing.T) {

	cache := InitGCache(100, time.Duration(60)*time.Second)

	err := cache.SetWithExpire("key_4", "value_4", time.Duration(30)*time.Second)
	if err != nil {
		t.Errorf("Cache set with expire err: %v.", err)
		return
	}

	value := cache.Get("key_4")
	t.Log(value)

	cache.Purge()

	valueAfterPurge := cache.Get("key_4")
	t.Log(valueAfterPurge)
}

func TestGetIFPresent(t *testing.T) {
	cache := InitGCache(3, time.Duration(3)*time.Second)

	err := cache.Set("key_5", "value_5")
	if err != nil {
		t.Errorf("Cache set err: %v.", err)
		return
	}

	value := cache.GetIfExist("key_5")
	t.Log(value)

	time.Sleep(time.Duration(3) * time.Second)
	valueAfterSleep := cache.GetIfExist("key_5")
	t.Log(valueAfterSleep)
}

func TestGetAll(t *testing.T) {
	cache := InitGCache(3, time.Duration(3)*time.Second)

	err := cache.Set("key_6", "value_6")
	if err != nil {
		t.Errorf("Cache set err: %v.", err)
		return
	}

	value := cache.GetAll(false)
	t.Log(value)
	valueTrue := cache.GetAll(true)
	t.Log(valueTrue)

	time.Sleep(time.Duration(3) * time.Second)
	valueAfterSleep := cache.GetAll(false)
	t.Log(valueAfterSleep)

	valueAfterSleepTrue := cache.GetAll(true)
	t.Log(valueAfterSleepTrue)
}

func TestRevove(t *testing.T) {
	cache := InitGCache(3, time.Duration(3)*time.Second)

	err := cache.Set("key_7", "value_7")
	if err != nil {
		t.Errorf("Cache set err: %v.", err)
		return
	}

	err = cache.Set("key_8", "value_8")
	if err != nil {
		t.Errorf("Cache set err: %v.", err)
		return
	}

	value := cache.GetAll(true)
	t.Log(value)

	cache.Remove("key_8")

	valueAfterRemoveTrue := cache.GetAll(true)
	t.Log(valueAfterRemoveTrue)
}

func TestKeys(t *testing.T) {
	cache := InitGCache(3, time.Duration(3)*time.Second)

	err := cache.Set("key_9", "value_9")
	if err != nil {
		t.Errorf("Cache set err: %v.", err)
		return
	}

	err = cache.Set("key_10", "value_10")
	if err != nil {
		t.Errorf("Cache set err: %v.", err)
		return
	}

	keys := cache.Keys(true)
	t.Log(keys)

	time.Sleep(time.Duration(3) * time.Second)

	keysAfterSleep := cache.Keys(true)
	t.Log(keysAfterSleep)
}

func TestLen(t *testing.T) {
	cache := InitGCache(3, time.Duration(3)*time.Second)

	err := cache.Set("key_11", "value_11")
	if err != nil {
		t.Errorf("Cache set err: %v.", err)
		return
	}

	i := cache.Len(true)
	t.Log(i)

	time.Sleep(time.Duration(3) * time.Second)

	j := cache.Len(true)
	t.Log(j)
}

func TestHas(t *testing.T) {
	cache := InitGCache(3, time.Duration(3)*time.Second)

	err := cache.Set("key_12", "value_12")
	if err != nil {
		t.Errorf("Cache set err: %v.", err)
		return
	}

	has := cache.Has("key_12")
	t.Log(has)

	time.Sleep(time.Duration(3) * time.Second)

	hasAfterSleep := cache.Has("key_12")
	t.Log(hasAfterSleep)
}

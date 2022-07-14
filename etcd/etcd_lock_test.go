package etcd

import (
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestEtcdLock(t *testing.T) {
	// 获取etcd
	etcdCfg := &EtcdConfig{
		Endpoints:   []string{"10.117.49.69:12379", "10.117.49.69:22379", "10.117.49.69:32379"},
		DialTimeout: time.Duration(15) * time.Second,
		Username:    "",
		Password:    "",
		TTL:         5, // 续约时间为：time.Now().Add((time.Duration(karesp.TTL) * time.Second) / 3.0)
	}

	c := make(chan os.Signal)
	signal.Notify(c)

	lockKey := "/lock"

	go func() {
		t.Log("start go1")

		lock, err := NewEtcdLock(etcdCfg, lockKey)
		if err != nil {
			t.Errorf("New etcd1 lock err: %v.", err)
			return
		}

		if err := lock.LockWithDuration(time.Duration(3) * time.Second); err != nil {
			t.Errorf("go1 get lock with 3 second err: %v", err)
			return
		}

		if err = lock.UnLock(); err != nil {
			t.Errorf("go1 get unlock err: %v", err)
			return
		}
		t.Log("go1 release lock")
	}()

	go func() {
		t.Log("start go2")

		lock, err := NewEtcdLock(etcdCfg, lockKey)
		if err != nil {
			t.Errorf("New etcd2 lock err: %v.", err)
			return
		}

		if err := lock.Lock(); err != nil {
			t.Errorf("go2 get lock err: %v", err)
			return
		}

		time.Sleep(time.Duration(5) * time.Second) // 休眠5s

		if err = lock.UnLock(); err != nil {
			t.Errorf("go2 get unlock err: %v", err)
			return
		}
		t.Log("go2 release lock")
	}()

	go func() {
		t.Log("start go3")

		lock, err := NewEtcdLock(etcdCfg, lockKey)
		if err != nil {
			t.Errorf("New etcd3 lock err: %v.", err)
			return
		}

		for {
			if err := lock.TryLock(); err != nil {
				t.Errorf("go3 try lock err: %v", err)
				time.Sleep(time.Duration(1) * time.Second)
			} else {
				t.Log("go3 try lock success.")
				break
			}
		}

		if err = lock.UnLock(); err != nil {
			t.Errorf("go3 get unlock err: %v", err)
			return
		}
		t.Log("go3 release lock")
	}()

	<-c
}

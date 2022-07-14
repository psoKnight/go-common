package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"sync"
	"time"
)

type EtcdLock struct {
	lockName string
	session  *concurrency.Session
	mutex    *concurrency.Mutex
}

var onceNew sync.Once
var etcdClient *clientv3.Client

// NewEtcdLock 新建etcd 分布式锁
func NewEtcdLock(cfg *EtcdConfig, lockName string) (*EtcdLock, error) {
	onceNew.Do(func() {
		c, err := NewEtcdClient(cfg)
		if err != nil {
			return
		}

		etcdClient = c
	})

	session, err := concurrency.NewSession(etcdClient)
	if err != nil {
		return nil, err
	}

	mutex := concurrency.NewMutex(session, lockName)
	if err != nil {
		return nil, err
	}

	return &EtcdLock{
		lockName: lockName,
		session:  session,
		mutex:    mutex,
	}, nil
}

// LockWithDuration 获取锁（存在过期时间）
func (el *EtcdLock) LockWithDuration(duration time.Duration) error {
	ctx, _ := context.WithTimeout(context.Background(), duration)
	return el.mutex.Lock(ctx)
}

// Lock 获取锁
func (el *EtcdLock) Lock() error {
	return el.mutex.Lock(context.Background())
}

// UnLock 释放锁
func (el *EtcdLock) UnLock() error {
	return el.mutex.Unlock(context.Background())
}

// TryLock 尝试获取锁
func (el *EtcdLock) TryLock() error {
	return el.mutex.TryLock(context.Background())
}

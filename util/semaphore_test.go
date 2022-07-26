package util

import (
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	semaphore := NewSemaphore(1)

	// 获取信号量
	semaphore.Acquire()

	// 尝试获取信号量
	acquire := semaphore.TryAcquire()
	t.Logf("Try acquire: %t.", acquire)

	// 尝试获取信号量（存在超时时间）
	semaphore.TryAcquireOnTime(time.Duration(3) * time.Second)
	t.Logf("Try acquire on time: %t.", acquire)

	// 释放信号量
	semaphore.Release()

	// 再次尝试获取信号量
	afterTryAcquire := semaphore.TryAcquire()
	t.Logf("Try acquire again: %t.", afterTryAcquire)
}

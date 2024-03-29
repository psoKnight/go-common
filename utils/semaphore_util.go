package utils

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Semaphore struct {
	permits int      // 许可数量
	channel chan int // 通道
}

// NewSemaphore 创建信号量
func NewSemaphore(permits int) *Semaphore {
	return &Semaphore{channel: make(chan int, permits), permits: permits}
}

// Acquire 获取许可
func (s *Semaphore) Acquire() {
	s.channel <- 0
	logrus.Infof("Semaphore acquire success, available permits: %d.", s.AvailablePermits())
}

// Release 释放许可
func (s *Semaphore) Release() {
	<-s.channel
	logrus.Infof("Semaphore release success, available permits: %d.", s.AvailablePermits())
}

// TryAcquire 尝试获取许可
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.channel <- 0:
		return true
	default:
		return false
	}
}

// TryAcquireOnTime 尝试指定时间内获取许可
func (s *Semaphore) TryAcquireOnTime(timeout time.Duration) bool {
	for {
		select {
		case s.channel <- 0:
			return true
		case <-time.After(timeout):
			return false
		}
	}
}

// AvailablePermits 当前可用的许可数
func (s *Semaphore) AvailablePermits() int {
	return s.permits - len(s.channel)
}

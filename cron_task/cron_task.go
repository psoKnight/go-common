package cron_task

import (
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

// CronTaskConfig 定义每个任务的配置
type CronTaskConfig struct {
	Name      string       // 任务名称
	Schedule  string       // Cron 表达式
	StartTime time.Time    // 调度开始日期
	EndTime   time.Time    // 调度截止日期
	TaskFunc  func()       // 任务执行的函数
	entryID   cron.EntryID // Cron 任务ID
}

// CronTaskManager 管理多个任务
type CronTaskManager struct {
	cron  *cron.Cron                 // Cron 实例
	tasks map[string]*CronTaskConfig // 任务列表
	mu    sync.Mutex                 // 锁，用于保证线程安全
}

// NewCronTaskManager 创建任务管理器实例
func NewCronTaskManager(opts ...cron.Option) *CronTaskManager {
	c := cron.New(opts...)
	cronTaskManager := &CronTaskManager{
		cron:  c,
		tasks: make(map[string]*CronTaskConfig),
	}
	return cronTaskManager
}

// AddTask 添加任务
func (ctm *CronTaskManager) AddTask(config *CronTaskConfig) error {
	ctm.mu.Lock()
	defer ctm.mu.Unlock()

	// 校验开始/结束时间
	if config.StartTime.After(config.EndTime) {
		return errors.New("start time can‘t be after end time")
	}

	// 校验结束时间
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	now := time.Now().In(loc)
	if now.After(config.EndTime) {
		//log.Printf("[AddTask]task '%s' end time has already passed, ignore add", config.Name)

		// 不注册已过期的任务
		return nil
	}

	// 注册任务
	entryID, err := ctm.cron.AddFunc(config.Schedule, func() {
		l, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			//log.Printf("[AddTask AddFunc]task '%s' failed to load time zone: %v", config.Name, err)
			return
		}
		n := time.Now().In(l)
		if n.Before(config.StartTime) || n.After(config.EndTime) {
			// 任务不在有效期内，跳过执行
			return
		}

		config.TaskFunc()
	})
	if err != nil {
		return errors.Errorf("failed to schedule task '%s': %v", config.Name, err)
	}

	config.entryID = entryID
	ctm.tasks[config.Name] = config

	//log.Printf("[AddTask]task '%s' added successfully", config.Name)
	return nil
}

// RemoveTask 移除任务
func (ctm *CronTaskManager) RemoveTask(taskName string) error {
	ctm.mu.Lock()
	defer ctm.mu.Unlock()

	task, exists := ctm.tasks[taskName]
	if !exists {
		//log.Printf("[RemoveTask]task '%s' not found, skip remove", taskName)
		return nil
	}
	ctm.cron.Remove(task.entryID)
	delete(ctm.tasks, taskName)

	//log.Printf("[RemoveTask]task '%s' removed successfully", taskName)
	return nil
}

// UpdateTask 更新任务（先移除再新增）
func (ctm *CronTaskManager) UpdateTask(config *CronTaskConfig) error {
	if err := ctm.RemoveTask(config.Name); err != nil {
		return err
	}
	if err := ctm.AddTask(config); err != nil {
		return err
	}

	//log.Printf("[UpdateTask]task '%s' updated successfully", config.Name)
	return nil
}

// ListTasks 列出当前所有任务
func (ctm *CronTaskManager) ListTasks() []*CronTaskConfig {
	ctm.mu.Lock()
	defer ctm.mu.Unlock()

	taskList := make([]*CronTaskConfig, 0, len(ctm.tasks))
	for _, task := range ctm.tasks {
		taskList = append(taskList, task)
	}
	return taskList
}

// Start 启动所有任务
func (ctm *CronTaskManager) Start() {
	ctm.cron.Start()
	//log.Println("CronTaskManager started")
}

// Stop 停止所有任务
func (ctm *CronTaskManager) Stop() {
	ctm.cron.Stop()
	//log.Println("CronTaskManager stopped")
}

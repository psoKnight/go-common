package cron_task

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"testing"
	"time"
)

// TestCronTask 完整示例
func TestCronTask(t *testing.T) {
	manager := NewCronTaskManager(
		cron.WithSeconds(), // 需要秒级精度
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)), // 防止任务重叠
	)
	manager.Start()
	defer manager.Stop()

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)

	// 添加任务
	task := &CronTaskConfig{
		Name:      "add",
		Schedule:  "* * * * * *",
		StartTime: now.Add(-1 * time.Second),
		EndTime:   now.Add(10 * time.Second),
		TaskFunc:  func() { fmt.Println("add running...") },
	}
	if err := manager.AddTask(task); err != nil {
		t.Fatal(err)
		return
	}

	time.Sleep(time.Duration(3) * time.Second)

	// 更新任务
	task.TaskFunc = func() { fmt.Println("update running...") }
	task.Schedule = "0/2 * * * * *"
	if err := manager.UpdateTask(task); err != nil {
		t.Fatal(err)
		return
	}

	time.Sleep(time.Duration(3) * time.Second)

	// 查询任务
	t.Log("tasks num:", len(manager.ListTasks()))

	// 移除任务
	if err := manager.RemoveTask(task.Name); err != nil {
		t.Fatal(err)
		return
	}
	t.Log("tasks num:", len(manager.ListTasks()))

}

// TestAddTask 测试 AddTask
func TestAddTask(t *testing.T) {
	manager := NewCronTaskManager()
	defer manager.Stop()
	manager.Start()

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)

	task := &CronTaskConfig{
		Name:      "at",
		Schedule:  "@every 200ms",
		StartTime: now.Add(-1 * time.Second),
		EndTime:   now.Add(5 * time.Second),
		TaskFunc:  func() {}, // 空函数，避免执行副作用
	}

	// 添加任务
	if err := manager.AddTask(task); err != nil {
		t.Fatalf("AddTask failed: %v", err)
		return
	}

	// 查询任务列表
	tasks := manager.ListTasks()
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
		return
	}
	if tasks[0].Name != "at" {
		t.Fatalf("expected task name 'at', got %s", tasks[0].Name)
		return
	}
}

// TestRemoveTask 测试 RemoveTask
func TestRemoveTask(t *testing.T) {
	manager := NewCronTaskManager()
	defer manager.Stop()
	manager.Start()

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	task := &CronTaskConfig{
		Name:      "rt",
		Schedule:  "@every 200ms",
		StartTime: now.Add(-1 * time.Second),
		EndTime:   now.Add(1 * time.Second),
		TaskFunc:  func() {},
	}
	manager.AddTask(task)

	if err := manager.RemoveTask("rt"); err != nil {
		t.Fatalf("RemoveTask failed: %v", err)
	}
	if len(manager.ListTasks()) != 0 {
		t.Fatalf("task still exists after RemoveTask")
	}
}

// TestUpdateTask 测试 UpdateTask
func TestUpdateTask(t *testing.T) {
	manager := NewCronTaskManager()
	defer manager.Stop()
	manager.Start()

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)

	// 记录执行日志
	executed := make(chan string, 1)

	task := &CronTaskConfig{
		Name:      "ut",
		Schedule:  "@every 200ms",
		StartTime: now.Add(-1 * time.Second),
		EndTime:   now.Add(2 * time.Second),
		TaskFunc: func() {
			executed <- "old"
		},
	}

	// 添加任务
	if err := manager.AddTask(task); err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	// 更新任务：打印不同内容
	task.TaskFunc = func() {
		executed <- "new"
	}
	if err := manager.UpdateTask(task); err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}

	// 验证更新后的函数被调用
	timeout := time.After(1 * time.Second)
	for {
		select {
		case v := <-executed:
			if v == "new" {
				t.Fatalf("updated task run new function")
				return // 测试通过
			}
			// 如果拿到 old，就继续等
		case <-timeout:
			t.Fatalf("updated task did not run new function")
		}
	}
}

// TestListTasks 测试 ListTasks
func TestListTasks(t *testing.T) {
	manager := NewCronTaskManager()
	defer manager.Stop()
	manager.Start()

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	task := &CronTaskConfig{
		Name:      "lt",
		Schedule:  "@every 200ms",
		StartTime: now.Add(-1 * time.Second),
		EndTime:   now.Add(1 * time.Second),
		TaskFunc:  func() {},
	}
	manager.AddTask(task)

	tasks := manager.ListTasks()
	if len(tasks) != 1 || tasks[0].Name != "lt" {
		t.Fatalf("ListTasks returned wrong result")
		return
	}
}

// TestStartStop 测试 Start 和 Stop
func TestStartStop(t *testing.T) {
	manager := NewCronTaskManager()

	manager.Start()
	if manager.cron == nil {
		t.Fatalf("cron instance should not be nil after Start")
	}

	manager.Stop()
	// Stop 没有返回值，只要不panic 就算通过
}

// TestAddTaskExpired 测试添加已过期任务
func TestAddTaskExpired(t *testing.T) {
	manager := NewCronTaskManager()
	defer manager.Stop()
	manager.Start()

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	expiredTask := &CronTaskConfig{
		Name:      "e",
		Schedule:  "@every 1s",
		StartTime: now.Add(-2 * time.Second),
		EndTime:   now.Add(-1 * time.Second),
		TaskFunc:  func() { t.Fatalf("expired task should not run") },
	}

	if err := manager.AddTask(expiredTask); err != nil {
		t.Fatalf("AddTask failed for expired task: %v", err)
	}
	if len(manager.ListTasks()) != 0 {
		t.Fatalf("expired task should not be added")
	}
}

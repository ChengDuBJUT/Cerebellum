package task

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"cerebellum/internal/memory"
)

// TaskType 任务类型
type TaskType string

const (
	TaskTypePeriodic TaskType = "periodic"
	TaskTypeOnce     TaskType = "once"
)

// BrainTask 大脑分配的任务
type BrainTask struct {
	ID       string            `json:"id"`
	Type     TaskType          `json:"type"`
	Interval string            `json:"interval,omitempty"`
	Command  string            `json:"command"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// TaskPlan 小脑生成的任务计划
type TaskPlan struct {
	ID        string    `json:"id"`
	Type      TaskType  `json:"type"`
	Command   string    `json:"command"`
	Interval  string    `json:"interval"` // 周期任务的间隔（不能省略）
	CreatedAt time.Time `json:"created_at"`
	NextRun   time.Time `json:"next_run,omitempty"`
	LastRun   time.Time `json:"last_run,omitempty"`
	ExecCount int       `json:"exec_count"`
	Status    string    `json:"status"`
	Result    string    `json:"result,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// TaskResult 完成的任务结果
type TaskResult struct {
	ID      string `json:"id"`
	Result  string `json:"result"`
	Command string `json:"command"`
}

// ChangeType 变化类型
type ChangeType string

const (
	ChangeTypeAdded     ChangeType = "added"
	ChangeTypeCompleted ChangeType = "completed"
	ChangeTypeFailed    ChangeType = "failed"
	ChangeTypeUpdated   ChangeType = "updated"
)

// TaskChange 任务变化
type TaskChange struct {
	Type      ChangeType `json:"type"`
	TaskID    string     `json:"task_id"`
	Timestamp time.Time  `json:"timestamp"`
	OldStatus string     `json:"old_status,omitempty"`
	NewStatus string     `json:"new_status,omitempty"`
	Result    string     `json:"result,omitempty"`
}

// PlanGenerator 生成任务计划（增强版）
type PlanGenerator struct {
	periodicTasks map[string]*TaskPlan
	onceTasks     map[string]*TaskPlan
	changes       []TaskChange
	changesMu     sync.Mutex
	taskCount     int
	lastTaskCount int
	memory        *memory.JSONLMemory
	dataDir       string
	mu            sync.Mutex
}

// NewPlanGenerator 创建计划生成器
func NewPlanGenerator(mem *memory.JSONLMemory) *PlanGenerator {
	return &PlanGenerator{
		periodicTasks: make(map[string]*TaskPlan),
		onceTasks:     make(map[string]*TaskPlan),
		changes:       make([]TaskChange, 0),
		memory:        mem,
	}
}

// GeneratePlan 从大脑任务生成计划
func (g *PlanGenerator) GeneratePlan(tasks []BrainTask) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.lastTaskCount = g.taskCount
	newTaskCount := 0

	for _, task := range tasks {
		if task.Type == TaskTypePeriodic {
			if _, exists := g.periodicTasks[task.ID]; !exists {
				// 确保 Interval 有默认值
				interval := task.Interval
				if interval == "" {
					interval = "30s" // 默认间隔 30 秒
					log.Printf("WARNING: Task %s has empty interval, using default 30s", task.ID)
				}
				nextRun := g.calcNextRun(interval, time.Now())
				g.periodicTasks[task.ID] = &TaskPlan{
					ID:        task.ID,
					Type:      TaskTypePeriodic,
					Command:   task.Command,
					Interval:  interval,
					CreatedAt: time.Now(),
					NextRun:   nextRun,
					Status:    "pending",
				}
				g.recordChange(ChangeTypeAdded, task.ID, "", "pending")
				newTaskCount++

				// 写入记忆
				if g.memory != nil {
					g.memory.Write("task_assigned", task.ID,
						fmt.Sprintf("New periodic task assigned: %s (interval: %s)", task.Command, interval),
						task)
				}
			}
		} else if task.Type == TaskTypeOnce {
			if _, exists := g.onceTasks[task.ID]; !exists {
				g.onceTasks[task.ID] = &TaskPlan{
					ID:        task.ID,
					Type:      TaskTypeOnce,
					Command:   task.Command,
					CreatedAt: time.Now(),
					NextRun:   time.Now(),
					Status:    "pending",
				}
				g.recordChange(ChangeTypeAdded, task.ID, "", "pending")
				newTaskCount++

				if g.memory != nil {
					g.memory.Write("task_assigned", task.ID,
						fmt.Sprintf("New one-time task assigned: %s", task.Command),
						task)
				}
			}
		}
	}

	g.taskCount = len(g.periodicTasks) + len(g.onceTasks)
	return newTaskCount
}

// calcNextRun 计算下次执行时间
func (g *PlanGenerator) calcNextRun(interval string, from time.Time) time.Time {
	d, err := time.ParseDuration(interval)
	if err != nil {
		return from.Add(time.Hour)
	}
	return from.Add(d)
}

// ExecuteTasks 执行待执行的任务
func (g *PlanGenerator) ExecuteTasks(executor func(command string) (string, error)) {
	now := time.Now()

	// 执行一次性任务
	for id, task := range g.onceTasks {
		if task.Status == "pending" || task.Status == "failed" {
			oldStatus := task.Status
			task.Status = "running"
			task.LastRun = now

			result, err := executor(task.Command)
			if err != nil {
				task.Status = "failed"
				task.Error = err.Error()
				g.recordChange(ChangeTypeFailed, id, oldStatus, "failed")

				if g.memory != nil {
					g.memory.Write("task_failed", id,
						fmt.Sprintf("Task failed: %v", err),
						nil)
				}
			} else {
				task.Status = "completed"
				task.Result = result
				task.ExecCount++
				g.recordChange(ChangeTypeCompleted, id, oldStatus, "completed")

				if g.memory != nil {
					g.memory.Write("task_completed", id,
						fmt.Sprintf("Task completed: %s", result),
						nil)
				}
			}
		}
	}

	// 执行周期任务
	for id, task := range g.periodicTasks {
		if now.Equal(task.NextRun) || now.After(task.NextRun) {
			oldStatus := task.Status
			task.Status = "running"
			task.LastRun = now
			task.ExecCount++

			// 确保 Interval 有值
			if task.Interval == "" {
				log.Printf("WARNING: Task %s has empty Interval, using default 30s", id)
				task.Interval = "30s"
			}

			result, err := executor(task.Command)
			if err != nil {
				task.Status = "failed"
				task.Error = err.Error()
			} else {
				task.Status = "completed"
				task.Result = result
			}

			task.NextRun = g.calcNextRun(task.Interval, now)
			g.recordChange(ChangeTypeUpdated, task.ID, oldStatus, task.Status)

			if g.memory != nil {
				g.memory.Write("task_executed", task.ID,
					fmt.Sprintf("Periodic task executed: %s", result),
					nil)
			}
		} else {
			task.Status = "pending"
		}
	}
}

// recordChange 记录任务变化
func (g *PlanGenerator) recordChange(changeType ChangeType, taskID, oldStatus, newStatus string) {
	g.changesMu.Lock()
	defer g.changesMu.Unlock()

	g.changes = append(g.changes, TaskChange{
		Type:      changeType,
		TaskID:    taskID,
		Timestamp: time.Now(),
		OldStatus: oldStatus,
		NewStatus: newStatus,
	})
}

// HasSignificantChanges 检查是否有显著变化（变化数 > 1）
func (g *PlanGenerator) HasSignificantChanges() bool {
	g.changesMu.Lock()
	defer g.changesMu.Unlock()

	return len(g.changes) > 1
}

// GetAndClearChanges 获取并清空变化列表
func (g *PlanGenerator) GetAndClearChanges() []TaskChange {
	g.changesMu.Lock()
	defer g.changesMu.Unlock()

	changes := g.changes
	g.changes = make([]TaskChange, 0)
	return changes
}

// GetTaskDelta 获取任务数量变化
func (g *PlanGenerator) GetTaskDelta() int {
	return g.taskCount - g.lastTaskCount
}

// GetReport 获取执行报告（增强版）
func (g *PlanGenerator) GetReport() map[string]interface{} {
	g.mu.Lock()
	defer g.mu.Unlock()

	var completed []TaskResult
	var failed []TaskResult
	pending := make([]string, 0)

	for id, task := range g.onceTasks {
		switch task.Status {
		case "completed":
			completed = append(completed, TaskResult{
				ID:      id,
				Result:  task.Result,
				Command: task.Command,
			})
		case "failed":
			failed = append(failed, TaskResult{
				ID:      id,
				Result:  task.Error,
				Command: task.Command,
			})
		case "pending":
			pending = append(pending, id)
		}
	}

	for id, task := range g.periodicTasks {
		switch task.Status {
		case "completed":
			completed = append(completed, TaskResult{
				ID:      id,
				Result:  task.Result,
				Command: task.Command,
			})
		case "failed":
			failed = append(failed, TaskResult{
				ID:      id,
				Result:  task.Error,
				Command: task.Command,
			})
		case "pending":
			pending = append(pending, id)
		}
	}

	return map[string]interface{}{
		"completed":      completed,
		"failed":         failed,
		"pending":        pending,
		"total_tasks":    g.taskCount,
		"periodic_count": len(g.periodicTasks),
		"once_count":     len(g.onceTasks),
		"timestamp":      time.Now().Format(time.RFC3339),
	}
}

// GetAllPlans 获取所有任务计划
func (g *PlanGenerator) GetAllPlans() []*TaskPlan {
	g.mu.Lock()
	defer g.mu.Unlock()

	var plans []*TaskPlan
	for _, t := range g.periodicTasks {
		plans = append(plans, t)
	}
	for _, t := range g.onceTasks {
		plans = append(plans, t)
	}
	return plans
}

// RemoveCompletedTask 移除已完成的任务
func (g *PlanGenerator) RemoveCompletedTask(id string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if task, exists := g.onceTasks[id]; exists && task.Status == "completed" {
		delete(g.onceTasks, id)
		g.taskCount--

		if g.memory != nil {
			g.memory.Write("task_removed", id, "Completed task removed", nil)
		}
		return true
	}
	return false
}

// WriteReportToFile 将报告写入文件
func WriteReportToFile(report map[string]interface{}, path string) error {
	content := fmt.Sprintf(`# Cerebellum Report
Generated: %s

## Summary
- Total Tasks: %d
- Periodic Tasks: %d
- Once Tasks: %d
- Completed: %d
- Failed: %d
- Pending: %d

## Completed (%d)
`,
		report["timestamp"],
		report["total_tasks"],
		report["periodic_count"],
		report["once_count"],
		len(report["completed"].([]TaskResult)),
		len(report["failed"].([]TaskResult)),
		len(report["pending"].([]string)),
		len(report["completed"].([]TaskResult)))

	for _, c := range report["completed"].([]TaskResult) {
		content += fmt.Sprintf("- **%s**: %s\n", c.ID, c.Result)
	}

	content += fmt.Sprintf("\n## Pending (%d)\n", len(report["pending"].([]string)))
	for _, p := range report["pending"].([]string) {
		content += fmt.Sprintf("- [ ] %s\n", p)
	}

	content += fmt.Sprintf("\n## Failed (%d)\n", len(report["failed"].([]TaskResult)))
	for _, f := range report["failed"].([]TaskResult) {
		content += fmt.Sprintf("- **%s**: %s\n", f.ID, f.Result)
	}

	return os.WriteFile(path, []byte(content), 0644)
}

// ParseBrainTasks 从JSON解析大脑任务
func ParseBrainTasks(data []byte) ([]BrainTask, error) {
	var tasks []BrainTask
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetMemoryStats 获取记忆统计
func (g *PlanGenerator) GetMemoryStats() (map[string]int, error) {
	if g.memory == nil {
		return nil, fmt.Errorf("memory not initialized")
	}
	return g.memory.GetStats()
}

// GetRecentMemory 获取最近的记忆
func (g *PlanGenerator) GetRecentMemory(n int) ([]memory.MemoryEntry, error) {
	if g.memory == nil {
		return nil, fmt.Errorf("memory not initialized")
	}
	return g.memory.ReadRecent(n)
}

// SaveTasks 保存任务状态到磁盘
func (g *PlanGenerator) SaveTasks() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.dataDir == "" {
		return nil
	}

	log.Printf("DEBUG SaveTasks: Saving %d periodic tasks", len(g.periodicTasks))

	// 调试：检查所有周期性任务的 Interval
	for id, task := range g.periodicTasks {
		if task.Interval == "" {
			log.Printf("WARNING: Task %s has empty Interval before save!", id)
		}
	}

	// 确保目录存在
	if err := os.MkdirAll(g.dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// 保存周期性任务
	periodicFile := filepath.Join(g.dataDir, "periodic_tasks.json")
	periodicData, err := json.MarshalIndent(g.periodicTasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal periodic tasks: %w", err)
	}
	log.Printf("DEBUG: Writing %d bytes to periodic_tasks.json", len(periodicData))
	if err := os.WriteFile(periodicFile, periodicData, 0644); err != nil {
		return fmt.Errorf("failed to write periodic tasks: %w", err)
	}

	// 保存一次性任务
	onceFile := filepath.Join(g.dataDir, "once_tasks.json")
	onceData, err := json.MarshalIndent(g.onceTasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal once tasks: %w", err)
	}
	if err := os.WriteFile(onceFile, onceData, 0644); err != nil {
		return fmt.Errorf("failed to write once tasks: %w", err)
	}

	return nil
}

// LoadTasks 从磁盘加载任务状态
func (g *PlanGenerator) LoadTasks() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.dataDir == "" {
		return nil
	}

	// 加载周期性任务（合并到现有map中）
	periodicFile := filepath.Join(g.dataDir, "periodic_tasks.json")
	if _, err := os.Stat(periodicFile); err == nil {
		data, err := os.ReadFile(periodicFile)
		if err != nil {
			return fmt.Errorf("failed to read periodic tasks: %w", err)
		}
		var loadedTasks map[string]*TaskPlan
		if err := json.Unmarshal(data, &loadedTasks); err != nil {
			return fmt.Errorf("failed to unmarshal periodic tasks: %w", err)
		}
		// 合并到现有map，保留新添加的任务
		for id, task := range loadedTasks {
			log.Printf("DEBUG LoadTasks: Loaded task %s, interval='%s'", id, task.Interval)
			// 修复：确保 interval 有默认值
			if task.Interval == "" {
				log.Printf("WARNING: Loaded task %s has empty interval, setting to 1m", id)
				task.Interval = "1m"
			}
			if _, exists := g.periodicTasks[id]; !exists {
				g.periodicTasks[id] = task
			}
		}
	}

	// 加载一次性任务
	onceFile := filepath.Join(g.dataDir, "once_tasks.json")
	if _, err := os.Stat(onceFile); err == nil {
		data, err := os.ReadFile(onceFile)
		if err != nil {
			return fmt.Errorf("failed to read once tasks: %w", err)
		}
		if err := json.Unmarshal(data, &g.onceTasks); err != nil {
			return fmt.Errorf("failed to unmarshal once tasks: %w", err)
		}
	}

	// 更新任务计数
	g.taskCount = len(g.periodicTasks) + len(g.onceTasks)
	g.lastTaskCount = g.taskCount

	return nil
}

// SetDataDir 设置数据目录
func (g *PlanGenerator) SetDataDir(dir string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.dataDir = dir
}

// GetPendingTasks 获取所有待处理任务
func (g *PlanGenerator) GetPendingTasks() []*TaskPlan {
	g.mu.Lock()
	defer g.mu.Unlock()

	var pending []*TaskPlan

	// 检查周期性任务
	for _, task := range g.periodicTasks {
		if task.Status == "pending" || task.Status == "failed" {
			pending = append(pending, task)
		}
	}

	// 检查一次性任务
	for _, task := range g.onceTasks {
		if task.Status == "pending" || task.Status == "failed" {
			pending = append(pending, task)
		}
	}

	return pending
}

// GetResumableTasks 获取可恢复执行的任务（用于重启后继续）
func (g *PlanGenerator) GetResumableTasks() []*TaskPlan {
	g.mu.Lock()
	defer g.mu.Unlock()

	var resumable []*TaskPlan
	now := time.Now()

	// 周期性任务：如果过了执行时间或状态为pending/failed
	for _, task := range g.periodicTasks {
		if task.Status == "failed" || task.Status == "pending" ||
			now.After(task.NextRun) || now.Equal(task.NextRun) {
			resumable = append(resumable, task)
		}
	}

	// 一次性任务：只恢复pending和failed的
	for _, task := range g.onceTasks {
		if task.Status == "pending" || task.Status == "failed" {
			resumable = append(resumable, task)
		}
	}

	return resumable
}

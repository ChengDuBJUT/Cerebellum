package task

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// TaskType 任务类型
type TaskType string

const (
	TaskTypePeriodic TaskType = "periodic" // 周期任务
	TaskTypeOnce     TaskType = "once"     // 一次性任务
)

// BrainTask 大脑分配的任务
type BrainTask struct {
	ID       string            `json:"id"`
	Type     TaskType          `json:"type"`
	Interval string            `json:"interval,omitempty"` // 如 "1h", "30m"
	Command  string            `json:"command"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// TaskPlan 小脑生成的任务计划
type TaskPlan struct {
	ID        string    `json:"id"`
	Type      TaskType  `json:"type"`
	Command   string    `json:"command"`
	CreatedAt time.Time `json:"created_at"`
	NextRun   time.Time `json:"next_run,omitempty"`
	LastRun   time.Time `json:"last_run,omitempty"`
	ExecCount int       `json:"exec_count"`
	Status    string    `json:"status"` // pending, running, completed, failed
	Result    string    `json:"result,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// TaskReport 执行报告
type TaskReport struct {
	Completed []TaskResult `json:"completed"`
	Pending   []string     `json:"pending"`
	Failed    []TaskResult `json:"failed"`
}

// TaskResult 完成的任务结果
type TaskResult struct {
	ID      string `json:"id"`
	Result  string `json:"result"`
	Command string `json:"command"`
}

// PlanGenerator 生成任务计划
type PlanGenerator struct {
	periodicTasks map[string]*TaskPlan
	onceTasks     map[string]*TaskPlan
	report        *TaskReport
}

// NewPlanGenerator 创建计划生成器
func NewPlanGenerator() *PlanGenerator {
	return &PlanGenerator{
		periodicTasks: make(map[string]*TaskPlan),
		onceTasks:     make(map[string]*TaskPlan),
		report: &TaskReport{
			Completed: []TaskResult{},
			Pending:   []string{},
			Failed:    []TaskResult{},
		},
	}
}

// GeneratePlan 从大脑任务生成计划
func (g *PlanGenerator) GeneratePlan(tasks []BrainTask) {
	now := time.Now()

	for _, task := range tasks {
		if task.Type == TaskTypePeriodic {
			if _, exists := g.periodicTasks[task.ID]; !exists {
				nextRun := g.calcNextRun(task.Interval, now)
				g.periodicTasks[task.ID] = &TaskPlan{
					ID:        task.ID,
					Type:      TaskTypePeriodic,
					Command:   task.Command,
					CreatedAt: now,
					NextRun:   nextRun,
					Status:    "pending",
				}
				g.report.Pending = append(g.report.Pending, task.ID)
			}
		} else if task.Type == TaskTypeOnce {
			if _, exists := g.onceTasks[task.ID]; !exists {
				g.onceTasks[task.ID] = &TaskPlan{
					ID:        task.ID,
					Type:      TaskTypeOnce,
					Command:   task.Command,
					CreatedAt: now,
					NextRun:   now,
					Status:    "pending",
				}
				g.report.Pending = append(g.report.Pending, task.ID)
			}
		}
	}
}

// calcNextRun 计算下次执行时间
func (g *PlanGenerator) calcNextRun(interval string, from time.Time) time.Time {
	d, err := time.ParseDuration(interval)
	if err != nil {
		return from.Add(time.Hour) // 默认1小时
	}
	return from.Add(d)
}

// ExecuteTasks 执行待执行的任务
func (g *PlanGenerator) ExecuteTasks(executor func(command string) (string, error)) {
	now := time.Now()

	// 执行一次性任务
	for id, task := range g.onceTasks {
		if task.Status == "pending" || task.Status == "failed" {
			task.Status = "running"
			task.LastRun = now

			result, err := executor(task.Command)
			if err != nil {
				task.Status = "failed"
				task.Error = err.Error()
				g.report.Failed = append(g.report.Failed, TaskResult{
					ID:      id,
					Result:  err.Error(),
					Command: task.Command,
				})
			} else {
				task.Status = "completed"
				task.Result = result
				task.ExecCount++
				g.report.Completed = append(g.report.Completed, TaskResult{
					ID:      id,
					Result:  result,
					Command: task.Command,
				})
			}
		}
	}

	// 执行周期任务
	for _, task := range g.periodicTasks {
		if now.Equal(task.NextRun) || now.After(task.NextRun) {
			task.Status = "running"
			task.LastRun = now
			task.ExecCount++

			result, err := executor(task.Command)
			if err != nil {
				task.Status = "failed"
				task.Error = err.Error()
			} else {
				task.Status = "completed"
				task.Result = result
			}

			// 计算下次执行时间
			task.NextRun = g.calcNextRun(task.Command, now)
		} else {
			task.Status = "pending"
		}
	}
}

// GetReport 获取执行报告
func (g *PlanGenerator) GetReport() TaskReport {
	return *g.report
}

// GetAllPlans 获取所有任务计划
func (g *PlanGenerator) GetAllPlans() []*TaskPlan {
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
	if task, exists := g.onceTasks[id]; exists && task.Status == "completed" {
		delete(g.onceTasks, id)
		return true
	}
	return false
}

// WriteReportToFile 将报告写入文件
func WriteReportToFile(report TaskReport, path string) error {
	content := fmt.Sprintf(`# Cerebellum Report
Generated: %s

## Completed (%d)
`, time.Now().Format(time.RFC3339), len(report.Completed))
	for _, c := range report.Completed {
		content += fmt.Sprintf(`- **%s**: %s
`, c.ID, c.Result)
	}

	content += fmt.Sprintf(`
## Pending (%d)
`, len(report.Pending))
	for _, p := range report.Pending {
		content += fmt.Sprintf(`- [ ] %s
`, p)
	}

	content += fmt.Sprintf(`
## Failed (%d)
`, len(report.Failed))
	for _, f := range report.Failed {
		content += fmt.Sprintf(`- **%s**: %s (error: %s)
`, f.ID, f.Command, f.Result)
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

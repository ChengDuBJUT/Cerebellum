package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"cerebellum/internal/config"
	"cerebellum/internal/llm"
	"cerebellum/internal/memory"
	"cerebellum/internal/store"
	"cerebellum/internal/task"
)

// Server represents the HTTP server
type Server struct {
	cfg            *config.Config
	store          *store.MarkdownStore
	llm            *llm.OllamaClient
	planner        *task.PlanGenerator
	systemIdentity string
	memory         *memory.JSONLMemory
	mu             sync.Mutex
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, store *store.MarkdownStore, llm *llm.OllamaClient) *Server {
	// Load system identity from skill-Cerebellum-EN.md
	systemIdentity := "You are Cerebellum, a helpful AI assistant."
	if content, err := os.ReadFile("skill-Cerebellum-EN.md"); err == nil {
		systemIdentity = string(content)
	}

	// Initialize memory system
	mem, err := memory.NewJSONLMemory("./data")
	if err != nil {
		log.Printf("Warning: Failed to initialize memory: %v", err)
		mem = nil
	}

	// Initialize planner with memory and data directory
	planner := task.NewPlanGenerator(mem)
	planner.SetDataDir("./data")

	// Load previous tasks from disk
	if err := planner.LoadTasks(); err != nil {
		log.Printf("Info: No previous tasks to load: %v", err)
	} else {
		log.Println("✓ Previous tasks loaded from disk")
		// Show resumable tasks
		resumable := planner.GetResumableTasks()
		if len(resumable) > 0 {
			log.Printf("✓ Found %d tasks to resume", len(resumable))
		}
	}

	return &Server{
		cfg:            cfg,
		store:          store,
		llm:            llm,
		planner:        planner,
		systemIdentity: systemIdentity,
		memory:         mem,
	}
}

// StartTaskExecutor 启动任务执行器（带持久化和智能报告）
func (s *Server) StartTaskExecutor() {
	// 立即执行一次，恢复之前的任务
	s.mu.Lock()
	resumableTasks := s.planner.GetResumableTasks()
	if len(resumableTasks) > 0 {
		log.Printf("Resuming %d tasks from previous session", len(resumableTasks))
		s.planner.ExecuteTasks(s.executeCommand)
		// 保存执行后的状态
		if err := s.planner.SaveTasks(); err != nil {
			log.Printf("Warning: Failed to save tasks: %v", err)
		}
	}
	s.mu.Unlock()

	// 启动定时器
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()

		// 执行任务
		s.planner.ExecuteTasks(s.executeCommand)

		// 保存任务状态到磁盘
		if err := s.planner.SaveTasks(); err != nil {
			log.Printf("Warning: Failed to save tasks: %v", err)
		}

		// 检查是否有显著变化（变化数 > 1），如果有则触发报告
		if s.planner.HasSignificantChanges() {
			changes := s.planner.GetAndClearChanges()
			log.Printf("Significant changes detected (%d), notifying brain...", len(changes))

			// 将变化记录到内存
			if s.memory != nil {
				for _, change := range changes {
					s.memory.Write("task_change", change.TaskID,
						fmt.Sprintf("Task %s: %s -> %s", change.Type, change.OldStatus, change.NewStatus),
						change)
				}
			}

			// 触发大脑报告（这里可以添加HTTP回调或Webhook）
			s.notifyBrain(changes)
		}

		s.mu.Unlock()
	}
}

// notifyBrain 通知大脑有重要变化
func (s *Server) notifyBrain(changes []task.TaskChange) {
	// 记录到日志（实际可以实现HTTP回调通知大脑）
	log.Printf("[BRAIN NOTIFICATION] %d task changes reported", len(changes))

	// 这里可以添加：发送到大脑的HTTP端点
	// report := map[string]interface{}{
	//     "timestamp": time.Now().Format(time.RFC3339),
	//     "type":      "task_changes",
	//     "count":     len(changes),
	//     "changes":   changes,
	//     "summary":   s.planner.GetReport(),
	// }
	// go s.sendToBrain(report)
}

// executeCommand 执行命令
func (s *Server) executeCommand(command string) (string, error) {
	prompt := fmt.Sprintf(`<system_instructions>
You are Cerebellum task executor. Execute the following command and provide the result.

Command: %s

Please execute this command and return the result.
</system_instructions>`, command)

	return s.llm.Generate(prompt)
}

// === Handler Functions ===

// HandleHealth 健康检查
func (s *Server) HandleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":      "ok",
		"llm_host":    s.llm.GetHost(),
		"llm_model":   s.llm.GetModel(),
		"tasks_count": len(s.store.GetTasks()),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleTasksRequest 大脑分配任务的请求
type HandleTasksRequest struct {
	Tasks []task.BrainTask `json:"tasks"`
}

// HandleTasksResponse 大脑分配任务的响应
type HandleTasksResponse struct {
	Status    string `json:"status"`
	TaskCount int    `json:"task_count"`
	Message   string `json:"message,omitempty"`
}

// HandleAPITasks POST /api/tasks - 大脑分配任务给小脑
func (s *Server) HandleAPITasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req HandleTasksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.planner.GeneratePlan(req.Tasks)
	planCount := len(s.planner.GetAllPlans())
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HandleTasksResponse{
		Status:    "accepted",
		TaskCount: planCount,
		Message:   fmt.Sprintf("Received %d tasks from brain", len(req.Tasks)),
	})
}

// HandleAPIReport GET /api/report - 获取执行报告
func (s *Server) HandleAPIReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	report := s.planner.GetReport()
	plans := s.planner.GetAllPlans()
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"report":          report,
		"total_plans":     len(plans),
		"pending_count":   len(report["pending"].([]string)),
		"completed_count": len(report["completed"].([]task.TaskResult)),
		"failed_count":    len(report["failed"].([]task.TaskResult)),
	})
}

// HandleAPIStatus GET /api/status - 获取小脑状态
func (s *Server) HandleAPIStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	report := s.planner.GetReport()
	plans := s.planner.GetAllPlans()
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":          "running",
		"llm_host":        s.llm.GetHost(),
		"llm_model":       s.llm.GetModel(),
		"total_tasks":     len(plans),
		"pending_tasks":   report["pending_count"],
		"completed_tasks": report["completed_count"],
		"failed_tasks":    report["failed_count"],
		"last_updated":    time.Now().Format(time.RFC3339),
	})
}

// HandleAPITaskDelete DELETE /api/task/{id} - 删除已完成任务
func (s *Server) HandleAPITaskDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/task/")
	if id == "" {
		http.Error(w, "Task ID required", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	removed := s.planner.RemoveCompletedTask(id)
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if removed {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "deleted",
			"id":     id,
		})
	} else {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "not_found",
			"id":      id,
			"message": "Task not found or not completed",
		})
	}
}

// === 原有端点 ===

// ChatRequest 聊天请求
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Response string `json:"response"`
}

// HandleChat POST /chat - 聊天
func (s *Server) HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	systemPrompt := s.buildSystemPrompt()
	prompt := fmt.Sprintf(`<system_instructions>
%s
</system_instructions>

<user_message>
%s
</user_message>

Please respond as a helpful assistant.`, systemPrompt, req.Message)

	response, err := s.llm.Generate(prompt)
	if err != nil {
		response = fmt.Sprintf("Error generating response: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{Response: response})
}

// buildSystemPrompt 构建系统提示
func (s *Server) buildSystemPrompt() string {
	Tasks := s.store.GetTasks()
	if len(Tasks) == 0 {
		return s.systemIdentity
	}

	var parts []string
	parts = append(parts, s.systemIdentity)
	parts = append(parts, "\n## Your Current Capabilities:\n")

	for _, t := range Tasks {
		parts = append(parts, fmt.Sprintf("- %s: %s", t.ID, t.Prompt))
	}

	return strings.Join(parts, "\n")
}

// ExecuteRequest 执行请求
type ExecuteRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// ExecuteResponse 执行响应
type ExecuteResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Error      string            `json:"error,omitempty"`
}

// HandleExecute POST /execute - HTTP 请求代理
func (s *Server) HandleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Method == "" {
		req.Method = "GET"
	}

	client := &http.Client{}
	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ExecuteResponse{
			Error: fmt.Sprintf("Failed to create request: %v", err),
		})
		return
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ExecuteResponse{
			Error: fmt.Sprintf("Failed to execute request: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	headers := make(map[string]string)
	for k := range resp.Header {
		headers[k] = resp.Header.Get(k)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ExecuteResponse{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       string(body),
	})
}

// TasksResponse 任务列表响应
type TasksResponse struct {
	Tasks []TaskInfo `json:"tasks"`
}

// TaskInfo 任务信息
type TaskInfo struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Prompt string `json:"prompt"`
}

// HandleTasks GET /tasks - 获取任务列表
func (s *Server) HandleTasks(w http.ResponseWriter, r *http.Request) {
	Tasks := s.store.GetTasks()
	var taskInfos []TaskInfo

	for _, t := range Tasks {
		taskInfos = append(taskInfos, TaskInfo{
			ID:     t.ID,
			Type:   t.Type,
			Prompt: t.Prompt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TasksResponse{Tasks: taskInfos})
}

// HandleReload POST /reload - 重新加载 brain.md
func (s *Server) HandleReload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := s.store.Reload(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to reload: %v", err), http.StatusInternalServerError)
		return
	}

	Tasks := s.store.GetTasks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "reloaded",
		"tasks_count": len(Tasks),
	})
}

// SaveTasks 保存任务状态到磁盘（用于优雅关闭）
func (s *Server) SaveTasks() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.planner.SaveTasks()
}

// HandleSetBeacon POST /api/beacon - 设置记忆信标
func (s *Server) HandleSetBeacon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name     string                 `json:"name"`
		Metadata map[string]interface{} `json:"metadata,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Beacon name required", http.StatusBadRequest)
		return
	}

	if s.memory == nil {
		http.Error(w, "Memory system not initialized", http.StatusInternalServerError)
		return
	}

	if err := s.memory.SetBeacon(req.Name, req.Metadata); err != nil {
		http.Error(w, fmt.Sprintf("Failed to set beacon: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "beacon_set",
		"name":   req.Name,
		"time":   time.Now().Format(time.RFC3339),
	})
}

// HandleReadMemory GET /api/memory?beacon=xxx&type=xxx - 读取记忆
func (s *Server) HandleReadMemory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.memory == nil {
		http.Error(w, "Memory system not initialized", http.StatusInternalServerError)
		return
	}

	beacon := r.URL.Query().Get("beacon")
	entryType := r.URL.Query().Get("type")

	if beacon == "" {
		// 如果没有指定信标，返回最近的100条
		entries, err := s.memory.ReadRecent(100)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read memory: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"entries": entries,
			"count":   len(entries),
		})
		return
	}

	// 从信标读取
	entries, err := s.memory.ReadSinceBeacon(beacon, entryType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read memory: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"beacon":  beacon,
		"entries": entries,
		"count":   len(entries),
	})
}

// HandleListBeacons GET /api/beacons - 列出所有信标
func (s *Server) HandleListBeacons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.memory == nil {
		http.Error(w, "Memory system not initialized", http.StatusInternalServerError)
		return
	}

	beacons, err := s.memory.ListBeacons()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list beacons: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"beacons": beacons,
		"count":   len(beacons),
	})
}

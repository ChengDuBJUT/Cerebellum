package store

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// BrainTask 从 brain.md 解析的任务
type BrainTask struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Prompt string `json:"prompt"`
}

// MarkdownStore handles reading and parsing of brain.md
type MarkdownStore struct {
	path      string
	lastMod   time.Time
	Tasks     []BrainTask
	lastError error
}

// NewMarkdownStore creates a new markdown store
func NewMarkdownStore(path string) (*MarkdownStore, error) {
	store := &MarkdownStore{
		path: path,
	}
	// 如果文件不存在，创建空任务
	if _, err := os.Stat(path); os.IsNotExist(err) {
		store.Tasks = []BrainTask{}
		return store, nil
	}
	if err := store.Reload(); err != nil {
		return nil, err
	}
	return store, nil
}

// GetPath returns the file path
func (s *MarkdownStore) GetPath() string {
	return s.path
}

// GetLastMod returns the last modification time
func (s *MarkdownStore) GetLastMod() time.Time {
	return s.lastMod
}

// GetTasks returns the loaded tasks
func (s *MarkdownStore) GetTasks() []BrainTask {
	return s.Tasks
}

// GetLastError returns the last error
func (s *MarkdownStore) GetLastError() error {
	return s.lastError
}

// Reload reloads the file and parses tasks
func (s *MarkdownStore) Reload() error {
	// 如果文件不存在，跳过
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		s.Tasks = []BrainTask{}
		s.lastError = nil
		return nil
	}

	file, err := os.Open(s.path)
	if err != nil {
		s.lastError = fmt.Errorf("failed to open file: %w", err)
		return s.lastError
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		s.lastError = fmt.Errorf("failed to stat file: %w", err)
		return s.lastError
	}

	s.lastMod = info.ModTime()
	Tasks, err := parseBrainFile(file)
	if err != nil {
		s.lastError = fmt.Errorf("failed to parse file: %w", err)
		return s.lastError
	}

	s.Tasks = Tasks
	s.lastError = nil
	return nil
}

// HasChanged checks if the file has been modified since last load
func (s *MarkdownStore) HasChanged() bool {
	info, err := os.Stat(s.path)
	if err != nil {
		return false
	}
	return !info.ModTime().Equal(s.lastMod)
}

// parseTaskLine parses a single task line like "- **task_id**: description"
func parseTaskLine(line string) *BrainTask {
	// Remove leading bullet point
	line = strings.TrimPrefix(line, "- ")
	line = strings.TrimPrefix(line, "* ")

	// Parse **id**: description format
	if !strings.HasPrefix(line, "**") {
		return nil
	}

	endIdx := strings.Index(line[2:], "**")
	if endIdx < 0 {
		return nil
	}

	id := line[2 : 2+endIdx]
	desc := strings.TrimSpace(line[2+endIdx+2:])
	if strings.HasPrefix(desc, ":") {
		desc = strings.TrimSpace(desc[1:])
	}

	return &BrainTask{
		ID:     id,
		Type:   "simple_qa",
		Prompt: desc,
	}
}

// parseBrainFile parses all tasks from the brain.md file
func parseBrainFile(file *os.File) ([]BrainTask, error) {
	scanner := bufio.NewScanner(file)
	var Tasks []BrainTask

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and headers
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse task lines
		if strings.HasPrefix(line, "- **") || strings.HasPrefix(line, "* **") {
			t := parseTaskLine(line)
			if t != nil {
				Tasks = append(Tasks, *t)
			}
		}
	}

	return Tasks, scanner.Err()
}

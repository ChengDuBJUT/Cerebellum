package memory

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// MemoryEntry 记忆条目
type MemoryEntry struct {
	Timestamp time.Time       `json:"timestamp"`
	Type      string          `json:"type"`
	TaskID    string          `json:"task_id,omitempty"`
	Content   string          `json:"content"`
	Data      json.RawMessage `json:"data,omitempty"`
}

// JSONLMemory JSONL 记忆管理器
type JSONLMemory struct {
	filePath string
	mu       sync.Mutex
	maxSize  int64
}

// NewJSONLMemory 创建新的 JSONL 记忆管理器
func NewJSONLMemory(dataDir string) (*JSONLMemory, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	m := &JSONLMemory{
		filePath: filepath.Join(dataDir, "cerebellum_memory.jsonl"),
		maxSize:  10 * 1024 * 1024,
	}

	if _, err := os.Stat(m.filePath); os.IsNotExist(err) {
		file, err := os.Create(m.filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create memory file: %w", err)
		}
		file.Close()
	}

	return m, nil
}

// Write 写入记忆
func (m *JSONLMemory) Write(entryType string, taskID string, content string, data interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.rotateIfNeeded(); err != nil {
		return err
	}

	entry := MemoryEntry{
		Timestamp: time.Now(),
		Type:      entryType,
		TaskID:    taskID,
		Content:   content,
	}

	if data != nil {
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}
		entry.Data = dataBytes
	}

	file, err := os.OpenFile(m.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open memory file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(entry); err != nil {
		return fmt.Errorf("failed to encode entry: %w", err)
	}

	return nil
}

// ReadAll 读取所有记忆
func (m *JSONLMemory) ReadAll() ([]MemoryEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	file, err := os.Open(m.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open memory file: %w", err)
	}
	defer file.Close()

	var entries []MemoryEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry MemoryEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}

// ReadRecent 读取最近的 n 条记忆
func (m *JSONLMemory) ReadRecent(n int) ([]MemoryEntry, error) {
	entries, err := m.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(entries) <= n {
		return entries, nil
	}

	return entries[len(entries)-n:], nil
}

// ReadByType 读取特定类型的记忆
func (m *JSONLMemory) ReadByType(entryType string, limit int) ([]MemoryEntry, error) {
	entries, err := m.ReadAll()
	if err != nil {
		return nil, err
	}

	var filtered []MemoryEntry
	for i := len(entries) - 1; i >= 0 && len(filtered) < limit; i-- {
		if entries[i].Type == entryType {
			filtered = append([]MemoryEntry{entries[i]}, filtered...)
		}
	}

	return filtered, nil
}

// GetStats 获取记忆统计
func (m *JSONLMemory) GetStats() (map[string]int, error) {
	entries, err := m.ReadAll()
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int)
	for _, entry := range entries {
		stats[entry.Type]++
	}

	return stats, nil
}

// Clear 清空记忆
func (m *JSONLMemory) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return os.Remove(m.filePath)
}

// rotateIfNeeded 如果需要则轮转文件
func (m *JSONLMemory) rotateIfNeeded() error {
	info, err := os.Stat(m.filePath)
	if err != nil {
		return nil
	}

	if info.Size() < m.maxSize {
		return nil
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := m.filePath + "." + timestamp + ".bak"

	if err := os.Rename(m.filePath, backupPath); err != nil {
		return fmt.Errorf("failed to rotate memory file: %w", err)
	}

	return nil
}

// SetMaxSize 设置最大文件大小
func (m *JSONLMemory) SetMaxSize(size int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.maxSize = size
}

// SetBeacon 设置记忆信标
func (m *JSONLMemory) SetBeacon(name string, metadata map[string]interface{}) error {
	entry := MemoryEntry{
		Timestamp: time.Now(),
		Type:      "beacon",
		TaskID:    name,
		Content:   fmt.Sprintf("Beacon set: %s", name),
	}

	if metadata != nil {
		dataBytes, _ := json.Marshal(metadata)
		entry.Data = dataBytes
	}

	return m.Write("beacon", name, entry.Content, metadata)
}

// ReadSinceBeacon 读取从信标以来的记忆
func (m *JSONLMemory) ReadSinceBeacon(beaconName string, entryType string) ([]MemoryEntry, error) {
	allEntries, err := m.ReadAll()
	if err != nil {
		return nil, err
	}

	var beaconTime time.Time
	for _, entry := range allEntries {
		if entry.Type == "beacon" && entry.TaskID == beaconName {
			beaconTime = entry.Timestamp
			break
		}
	}

	if beaconTime.IsZero() {
		return nil, fmt.Errorf("beacon '%s' not found", beaconName)
	}

	var result []MemoryEntry
	for _, entry := range allEntries {
		if entry.Timestamp.After(beaconTime) || entry.Timestamp.Equal(beaconTime) {
			if entryType == "" || entry.Type == entryType {
				result = append(result, entry)
			}
		}
	}

	return result, nil
}

// ReadBetweenBeacons 读取两个信标之间的记忆
func (m *JSONLMemory) ReadBetweenBeacons(startBeacon, endBeacon string, entryType string) ([]MemoryEntry, error) {
	allEntries, err := m.ReadAll()
	if err != nil {
		return nil, err
	}

	var startTime, endTime time.Time
	for _, entry := range allEntries {
		if entry.Type == "beacon" && entry.TaskID == startBeacon {
			startTime = entry.Timestamp
		}
		if entry.Type == "beacon" && entry.TaskID == endBeacon {
			endTime = entry.Timestamp
		}
	}

	if startTime.IsZero() {
		return nil, fmt.Errorf("start beacon '%s' not found", startBeacon)
	}
	if endTime.IsZero() {
		endTime = time.Now()
	}

	var result []MemoryEntry
	for _, entry := range allEntries {
		if (entry.Timestamp.After(startTime) || entry.Timestamp.Equal(startTime)) &&
			entry.Timestamp.Before(endTime) {
			if entryType == "" || entry.Type == entryType {
				result = append(result, entry)
			}
		}
	}

	return result, nil
}

// ListBeacons 列出所有信标
func (m *JSONLMemory) ListBeacons() ([]MemoryEntry, error) {
	return m.ReadByType("beacon", 100)
}

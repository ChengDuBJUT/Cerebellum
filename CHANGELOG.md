# Changelog

## [v1.0.1] - 2026-02-09

### Added

#### Beacon-Based Memory System (记忆信标系统)
- **New API Endpoints**:
  - `POST /api/beacon` - 设置记忆信标/检查点
  - `GET /api/beacons` - 列出所有信标
  - `GET /api/memory` - 读取记忆（支持时间窗口查询）
  - `GET /api/memory?beacon=NAME` - 读取指定信标之后的记忆
  - `GET /api/memory?beacon=NAME&type=TYPE` - 按类型过滤

- **Features**:
  - Named memory checkpoints for time-windowed queries
  - Persistent JSONL storage in `./data/cerebellum_memory.jsonl`
  - Metadata support for each beacon
  - Type filtering for efficient queries

#### Task Persistence (任务持久化)
- Automatic task state saving to `./data/periodic_tasks.json` and `./data/once_tasks.json`
- Task recovery on restart
- Task execution count tracking
- Graceful shutdown with task state save

#### Enhanced Task Execution (增强任务执行)
- Default interval fallback (30s) when not specified
- Interval field properly saved and restored
- Execution history tracking
- Smart brain notification on task changes

### Changed

#### Task Planner (`internal/task/planner.go`)
- Added `Interval` field to `TaskPlan` struct (non-optional)
- Added `TaskChange` struct for tracking task state changes
- Added `changes` slice with mutex for thread safety
- Added `memory` reference for logging to JSONL
- Added `SaveTasks()` and `LoadTasks()` for persistence
- Fixed `calcNextRun()` to use task.Interval correctly
- Enhanced `GeneratePlan()` with default interval handling

#### HTTP Server (`internal/server/http.go`)
- Added memory system initialization
- Added task persistence on startup and shutdown
- Enhanced `StartTaskExecutor()` with immediate task recovery
- Added task state saving every 30 seconds
- Added brain notification on significant changes
- New handlers: `HandleSetBeacon`, `HandleReadMemory`, `HandleListBeacons`

#### Documentation (`skill.md`, `skill-Cerebellum-EN.md`)
- Added Beacon API documentation
- Added usage examples for time-windowed queries
- Updated best practices section
- Added ETH price monitor demo documentation

### Fixed

- **Periodic task execution**: Fixed issue where tasks only executed once
  - Root cause: `Interval` field was being cleared during JSON serialization
  - Fix: Made `Interval` field non-optional (`json:"interval"` instead of `json:"interval,omitempty`)
  - Fix: Added default interval (30s) fallback

- **Task persistence**: Fixed issue where loaded tasks overwrote new tasks
  - Fix: Changed `json.Unmarshal` to merge instead of replace
  - Fix: Added interval restoration in `LoadTasks()`

### Demo Files Added

- `test-eth-monitor.bat` - Windows demo script
- `test-eth-monitor.sh` - Unix demo script
- `test-beacon-quick.bat` - Quick beacon test
- `DEMO-ETH-MONITOR.md` - Demo documentation

### Performance

- Task execution cycle: 30 seconds
- Memory file rotation: 10MB threshold
- All operations thread-safe with mutex

### Breaking Changes

None. This release is fully backward compatible.

---

## [v1.0.0] - 2026-02-08

Initial release with basic task execution capabilities.


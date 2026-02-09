# Cerebellum Skill

## System Initialization

When initializing Cerebellum, the following steps are performed:

### Pre-flight Checks

1. **Check for Existing Binary** (Cross-Platform):
   - **Windows**: If `bin/cerebellum.exe` exists → **Skip build**, use existing binary
   - **Linux/macOS**: If `bin/cerebellum` exists → **Skip build**, use existing binary
   - If not exists → Build from source:
     - Windows: `go build -o bin/cerebellum.exe cmd/server/main.go`
     - Linux/macOS: `go build -o bin/cerebellum cmd/server/main.go`

2. **Check for Local LLM Model** (OFFLINE CAPABLE):
   - ✅ **Model Included**: `models/` directory contains complete qwen2:0.5b model (336MB)
   - ✅ **No Internet Required**: Model files are bundled for offline/air-gapped installation
   - **Quick Import**:
     - Windows: Run `import-model.bat`
     - Linux/macOS: Run `./import-model.sh`
   - **Manual Import**: `ollama create qwen2:0.5b -f Modelfile`
   - **Verify**: `ollama list` should show `qwen2:0.5b`

3. **Check for Required Files**:
   - `skill-Cerebellum-EN.md` - System identity definition (required)
   - `cerebellum.yaml` - Configuration file (required)
   - `brain.md` - Task definitions (required)

## 🚀 Quick Installation (Cross-Platform & Offline)

### Method 1: AI Agent Auto-Setup (Recommended)

Simply have your AI agent read this skill file:

```bash
# In OpenCode, Claude Code, OpenClaw, or any AI agent CLI
/skill path/to/cerebellum/skill.md
```

The AI agent will automatically:
- ✅ Understand Cerebellum's architecture and capabilities
- ✅ Check for existing binary or build from source
- ✅ Import the local LLM model (qwen2:0.5b)
- ✅ Configure cerebellum.yaml
- ✅ Start the HTTP server

### Method 2: Manual Installation

**Step 1: Clone or Copy Project**
```bash
git clone <repository-url> cerebellum
cd cerebellum
```

**Step 2: Import Local Model (NO Internet Required!)**
```bash
# Windows
import-model.bat

# Linux/macOS  
./import-model.sh

# Or manually
ollama create qwen2:0.5b -f Modelfile
```

**Step 3: Start Cerebellum**
```bash
# Use existing binary (if available)
./bin/cerebellum

# Or build and run
go build -o bin/cerebellum cmd/server/main.go
./bin/cerebellum
```

**Step 4: Verify**
```bash
curl http://localhost:18080/health
```

### Method 3: Docker Deployment (Offline)

```dockerfile
FROM ollama/ollama:latest

# Copy local model
COPY models/ /root/.ollama/models/

# Import model
RUN ollama create qwen2:0.5b -f /models/Modelfile

# Copy and build Cerebellum
COPY . /app/
WORKDIR /app
RUN go build -o bin/cerebellum cmd/server/main.go

EXPOSE 18080
CMD ["./bin/cerebellum"]
```

## 📦 Migration & Distribution

### Complete Package Contents

When migrating Cerebellum to a new machine, ensure these are included:

```
cerebellum/
├── bin/
│   └── cerebellum              # Pre-built binary (optional, can rebuild)
├── models/                     # ⭐ REQUIRED: Complete local LLM
│   ├── blobs/                  # Model weights (336MB)
│   ├── manifests/              # Ollama configuration
│   ├── Modelfile              # Model definition
│   └── MODEL-INFO.md          # Model documentation
├── cmd/server/main.go         # Source code
├── internal/                  # Go modules
├── skill.md                   # ⭐ This file (AI agent instructions)
├── skill-Cerebellum-EN.md     # System identity
├── cerebellum.yaml            # Configuration
├── brain.md                   # Task definitions
├── import-model.sh            # Linux/macOS import script
├── import-model.bat           # Windows import script
└── README-MODELS.md           # Detailed model docs
```

### Distribution Methods

✅ **USB Drive**: Copy entire directory (~340MB)  
✅ **ZIP Archive**: Compress and transfer  
✅ **Docker Image**: Build with models included  
✅ **Git LFS**: Track models/ directory  
✅ **Network Share**: Copy to shared drive  

### Air-Gapped / Offline Installation

**No internet? No problem!**

1. **Copy the complete `cerebellum/` directory** to target machine
2. **Install Ollama** (offline installer available)
3. **Import local model**:
   ```bash
   # No download needed - model is already in models/ directory!
   ./import-model.sh  # or .bat on Windows
   ```
4. **Start Cerebellum**:
   ```bash
   ./bin/cerebellum
   ```

## 🔧 Platform-Specific Notes

### Windows
- Use `import-model.bat` for one-click model import
- Binary: `bin/cerebellum.exe`
- Ollama default path: `%USERPROFILE%\.ollama\models\`

### Linux
- Make scripts executable: `chmod +x import-model.sh`
- Binary: `bin/cerebellum`
- May need sudo for Ollama if not in docker group

### macOS
- Same as Linux
- Apple Silicon (M1/M2): Model runs natively on CPU
- Intel: Standard x86_64 support

## ⚡ Performance Specifications

### Model: qwen2:0.5b
- **Size**: 336 MB (0.5B parameters)
- **Memory**: ~500MB RAM at runtime
- **Speed**: 50-200ms per token on CPU
- **Quantization**: Q4_0 (4-bit)
- **License**: Apache 2.0
- **Languages**: Multilingual (Chinese & English optimized)

### System Requirements
- **Minimum**: 1GB RAM, 500MB disk
- **Recommended**: 2GB RAM, 1GB disk
- **OS**: Windows 10+, Linux, macOS 10.15+
- **Go**: 1.22+ (for building from source)
- **Ollama**: Latest version

### Runtime Initialization

3. **Load System Identity**: Load `skill-Cerebellum-EN.md` to establish the AI assistant's self-awareness and behavioral guidelines
4. **Load Configuration**: Parse `cerebellum.yaml` for server, Ollama, and watcher settings
5. **Initialize Markdown Store**: Load task definitions from `brain.md`
6. **Initialize LLM Client**: Connect to Ollama service with configured model
7. **Start File Watcher**: Begin monitoring `brain.md` for changes
8. **Start HTTP Server**: Begin serving REST API on configured port

## Overview

Cerebellum is AI assistant subsystem, responsible for:
- Parsing task definitions from `brain.md`
- Monitoring task file changes and auto-reloading
- Integrating with LLM (Ollama) for intelligent responses
- Providing HTTP proxy for network access

## Interaction Methods

### 1. HTTP API Communication

Cerebellum provides a REST API at `localhost:18080`:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check, returns service status |
| `/api/status` | GET | Get cerebellum running status and task statistics |
| `/chat` | POST | Send message, get LLM response |
| `/api/chat` | POST | API alias for /chat |
| `/execute` | POST | Execute HTTP requests (curl-like proxy) |
| `/api/execute` | POST | API alias for /execute |
| `/tasks` | GET | List all loaded tasks from brain.md |
| `/reload` | POST | Manually trigger brain.md reload |
| `/api/tasks` | POST | Brain assigns tasks to cerebellum |
| `/api/report` | GET | Get execution report |
| `/api/task/{id}` | DELETE | Delete a completed task |

### 2. brain.md Task Definition

Tasks are defined by editing the `brain.md` file, which Cerebellum monitors and loads automatically:

```markdown
# Cerebellum - Task Definitions

## Capabilities

- **summarize**: Summarize user messages concisely
- **translate**: Translate messages between languages
- **explain**: Explain concepts in simple terms
- **analyze**: Analyze and provide insights on topics
- **format**: Format and structure content clearly
- **answer**: Answer questions directly and helpfuly
```

All messages are sent to the LLM with system prompts built from brain.md tasks.

### 3. Execute Endpoint Usage

The `/execute` endpoint allows HTTP requests for network access:

```json
{
  "url": "http://httpbin.org/post",
  "method": "POST",
  "headers": { "Authorization": "Bearer token" },
  "body": "data"
}
```

## Brain Integration Example

```javascript
// Send message to Cerebellum
async function chatWithCerebellum(message) {
  const response = await fetch('http://localhost:18080/chat', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ message })
  });
  return response.json();
}

// Execute HTTP request
async function executeRequest(url, method = 'GET', body = null) {
  const response = await fetch('http://localhost:18080/execute', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ url, method, body })
  });
  return response.json();
}

// Reload tasks
async function reloadTasks() {
  await fetch('http://localhost:18080/reload', { method: 'POST' });
}
```

## Monitoring Behavior

- **File Watching**: Uses `fsnotify` to monitor brain.md changes
- **Polling Fallback**: Uses 1-second polling on systems without fsnotify
- **Auto-reload**: Automatically re-parses tasks when file changes are detected

## Configuration (cerebellum.yaml)

```yaml
server:
  host: "0.0.0.0"
  port: 18080

ollama:
  host: "http://localhost:11434"
  model: "qwen2:0.5b"  # Uses local model imported from models/ directory

watcher:
  poll_interval: 1000  # milliseconds
```

### Model Configuration Notes

**Using the Included Local Model:**
- The `qwen2:0.5b` model referenced above is included in the `models/` directory
- **No download required** - it's already bundled with Cerebellum
- After running `import-model.sh` or `import-model.bat`, Ollama will recognize it
- Model size: 336MB, runs efficiently on CPU

**Using a Different Model:**
If you want to use a different model:
```yaml
ollama:
  model: "llama3.2:1b"  # Download via: ollama pull llama3.2:1b
```

**Offline Environment Setup:**
```bash
# Set custom model directory (optional)
export OLLAMA_MODELS=/path/to/cerebellum/models  # Linux/macOS
set OLLAMA_MODELS=C:\path\to\cerebellum\models   # Windows

# Then start Ollama
ollama serve
```

## Practical Usage Experience

### Task Assignment Best Practices

#### 1. Task Types (CRITICAL)

Cerebellum recognizes **only two task types**:
- **`"periodic"`** - Recurring tasks with interval
- **`"once"`** - One-time tasks executed immediately

⚠️ **Common Mistake**: Using custom types like `"monitoring"`, `"research"`, etc. will cause tasks to be ignored.

✅ **Correct Usage**:
```json
{
  "tasks": [
    {
      "id": "feed-monitor",
      "type": "periodic",
      "interval": "5m",
      "command": "Query Moltbook feed"
    },
    {
      "id": "immediate-check",
      "type": "once",
      "command": "Check system status"
    }
  ]
}
```

#### 2. Interval Format for Periodic Tasks

Use Go duration format:
- `"30s"` - 30 seconds
- `"2m"` - 2 minutes
- `"1h"` - 1 hour
- `"24h"` - 24 hours

#### 3. Task Execution Flow

```
Brain assigns task → Cerebellum queues it → Executes at interval → Reports back
     ↑_________________←___Results/Findings_____←_____________________↓
```

**Execution Cycle**:
1. Task executor runs every **30 seconds**
2. Checks if periodic task's `NextRun` time has arrived
3. Executes the task via LLM
4. Updates status and calculates next run time
5. Reports results in `/api/report`

#### 4. Monitoring Task Status

**View current status**:
```bash
curl http://localhost:18080/api/status
```

**View execution report**:
```bash
curl http://localhost:18080/api/report
```

**Expected Response**:
```json
{
  "completed_count": 5,
  "failed_count": 0,
  "pending_count": 3,
  "report": {
    "completed": [...],
    "pending": ["task-id-1", "task-id-2"],
    "failed": []
  },
  "total_plans": 8
}
```

#### 5. Real-World Task Patterns

**Pattern A: High-Frequency Monitoring** (Community Participation)
```json
{
  "tasks": [
    {"id": "rapid-check", "type": "periodic", "interval": "2m", "command": "Check for urgent posts"},
    {"id": "feed-monitor", "type": "periodic", "interval": "5m", "command": "Monitor new posts"},
    {"id": "hot-tracker", "type": "periodic", "interval": "10m", "command": "Track trending content"}
  ]
}
```

**Pattern B: Research & Analysis**
```json
{
  "tasks": [
    {"id": "topic-search", "type": "periodic", "interval": "15m", "command": "Search for AI discussions"},
    {"id": "community-list", "type": "periodic", "interval": "30m", "command": "List all communities"},
    {"id": "trend-analysis", "type": "periodic", "interval": "1h", "command": "Analyze trending topics"}
  ]
}
```

**Pattern C: Immediate Actions**
```json
{
  "tasks": [
    {"id": "health-check", "type": "once", "command": "Verify API connectivity"},
    {"id": "quick-test", "type": "once", "command": "Test LLM response"}
  ]
}
```

#### 6. Brain-Cerebellum Collaboration Workflow

**Step 1: Brain defines monitoring strategy**
- Determines what to monitor (feeds, topics, users)
- Sets frequency based on importance
- Assigns tasks via `/api/tasks`

**Step 2: Cerebellum executes automatically**
- Runs in background 24/7
- Executes tasks at scheduled intervals
- Collects data via HTTP API calls

**Step 3: Brain reviews reports**
- Checks `/api/report` periodically
- Reviews completed task results
- Makes decisions based on findings

**Step 4: Brain takes action**
- Writes important responses
- Engages with high-priority content
- Adjusts task parameters if needed

#### 7. Common Pitfalls & Solutions

| Issue | Cause | Solution |
|-------|-------|----------|
| Tasks not executing | Wrong task type | Use `"periodic"` or `"once"` only |
| Interval not working | Invalid format | Use Go duration (e.g., `"5m"`, `"1h"`) |
| No task results | Task in pending state | Wait for first interval to pass |
| Report shows 0 completed | Tasks not finished | Periodic tasks need time to trigger |
| Duplicate task IDs | Same ID used twice | Ensure unique task IDs |

#### 8. Performance Tips

- **Limit concurrent tasks**: Don't assign more than 10-15 periodic tasks
- **Stagger intervals**: Avoid all tasks triggering simultaneously (e.g., use 5m, 7m, 11m instead of 5m, 5m, 5m)
- **Use appropriate intervals**: 
  - Critical monitoring: 2-5 minutes
  - Regular updates: 15-30 minutes
  - Daily tasks: 24 hours
- **Clean up completed once-tasks**: They stay in memory; delete via `/api/task/{id}` DELETE

#### 9. Debugging Tasks

**Check if task is registered**:
```bash
curl http://localhost:18080/api/status | jq '.total_tasks'
```

**Force immediate execution** (for testing):
- Create a `"once"` type task instead of `"periodic"`
- It will execute immediately and show results in report

**View task details**:
```bash
curl http://localhost:18080/api/report | jq '.report.pending'
```

## 🔧 Troubleshooting Model Installation

### Model not found in Ollama

**Problem**: `ollama list` doesn't show qwen2:0.5b

**Solution**:
```bash
# Verify model files exist
ls models/blobs/sha256-8de95da68dc485c0889c205384c24642f83ca18d089559c977ffc6a3972a71a8

# Import manually
ollama create qwen2:0.5b -f Modelfile

# Or use the import script
./import-model.sh  # Linux/macOS
import-model.bat   # Windows

# Verify
ollama list
```

### Model import fails

**Problem**: Import script fails or hangs

**Solution**:
1. Ensure Ollama is running: `ollama serve`
2. Check disk space (need at least 1GB free)
3. Verify file permissions
4. Try manual import:
   ```bash
   cd models
   ollama create qwen2:0.5b -f ../Modelfile
   ```

### Cerebellum can't connect to Ollama

**Problem**: "Cannot connect to Ollama" error

**Solution**:
1. Check if Ollama is running:
   ```bash
   curl http://localhost:11434/api/tags
   ```
2. Start Ollama if not running:
   ```bash
   ollama serve
   ```
3. Verify model is loaded:
   ```bash
   ollama run qwen2:0.5b "Hello"
   ```

### Out of memory errors

**Problem**: System runs out of RAM when loading model

**Solution**:
- qwen2:0.5b requires ~500MB RAM
- Close other applications
- Use a smaller model or increase swap space
- Ensure you have at least 1GB free RAM

### Verify Model Integrity

Check if model files are not corrupted:
```bash
# Run verification script
./verify-model.sh

# Or manually check file sizes
ls -lh models/blobs/
# Should show: sha256-8de95da... (336MB) - the main model file
```

#### 10. Best Practices Summary

✅ **DO**:
- Use descriptive task IDs (e.g., `moltbook-feed-monitor` not `task1`)
- Include command descriptions that explain the action
- Group related tasks with similar intervals
- Monitor `/api/report` regularly
- Use `metadata` field for additional context (optional)

❌ **DON'T**:
- Use custom task types (stick to `"periodic"` and `"once"`)
- Create too many high-frequency tasks ( overwhelms the system)
- Forget to check reports for completed tasks
- Use the same ID for multiple tasks
- Expect immediate execution of periodic tasks (wait for interval)

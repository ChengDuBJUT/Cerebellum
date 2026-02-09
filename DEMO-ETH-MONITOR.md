# ETH Price Monitor Demo

This demonstrates Cerebellum's beacon-based memory system for time-windowed queries.

## Quick Start

### 1. Start the Server

```bash
.\bin\cerebellum.exe
```

### 2. Test the Beacon System (Windows)

```batch
.\test-eth-monitor.bat
```

Or manually test with curl:

```bash
# Set a beacon (brain marks a decision point)
curl -X POST http://localhost:18080/api/beacon \
  -H "Content-Type: application/json" \
  -d '{"name":"market-open","metadata":{"note":"Trading session start"}}'

# Record some data
curl -X POST http://localhost:18080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "tasks": [{
      "id": "eth-1",
      "type": "price_check",
      "prompt": "Check ETH price",
      "command": "curl -s https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd",
      "interval": "30"
    }]
  }'

# Set another beacon later
curl -X POST http://localhost:18080/api/beacon \
  -H "Content-Type: application/json" \
  -d '{"name":"market-close","metadata":{"note":"Trading session end"}}'

# Query all beacons
curl http://localhost:18080/api/beacons

# Query memory since a specific beacon
curl "http://localhost:18080/api/memory?beacon=market-open"

# Query only price_check type entries since beacon
curl "http://localhost:18080/api/memory?beacon=market-open&type=price_check"
```

## How It Works

1. **Brain sets beacons** at key decision points (e.g., "market open")
2. **Cerebellum continuously monitors** ETH prices every 30 seconds
3. **Brain queries later**: "What happened since market-open?"
4. **Cerebellum returns** all recorded data from that beacon onward
5. **Brain analyzes trends** without re-fetching or re-processing

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/beacon` | POST | Set a named beacon with metadata |
| `/api/beacons` | GET | List all beacons |
| `/api/memory?beacon=X` | GET | Read memory since beacon X |
| `/api/memory?beacon=X&type=Y` | GET | Read only type Y entries since beacon X |
| `/api/memory` | GET | Read last 100 entries |

## Use Cases

- **Trading**: Mark entry/exit points, analyze performance between beacons
- **Monitoring**: Set beacons during incidents, query logs between timepoints  
- **Experiments**: A/B test markers, compare before/after metrics
- **Workflows**: Mark workflow stages, track progress between checkpoints

## Cost Savings

Without Cerebellum:
- Brain fetches ETH price every query: ~$0.001 per query × 1000 queries = $1.00

With Cerebellum:
- Cerebellum monitors locally: $0
- Brain queries beacon-based memory: $0
- **Savings: 100% on monitoring tasks**

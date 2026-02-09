@echo off
REM ETH Price Monitoring Demo for Cerebellum (Windows)
REM This script demonstrates the beacon-based memory system

setlocal enabledelayedexpansion

set BASE_URL=http://localhost:18080
set DATA_FILE=./data/eth_prices.jsonl

echo ==========================================
echo ETH Price Monitor Demo
echo ==========================================
echo.

REM Check if server is running
curl -s %BASE_URL%/health >nul 2>&1
if errorlevel 1 (
    echo Error: Cerebellum server is not running on %BASE_URL%
    echo Please start the server first: .\bin\cerebellum.exe
    exit /b 1
)

echo [OK] Server is running
echo.

REM Function to set beacon
echo Step 1: Setting initial beacon 'monitor-start'...
curl -s -X POST %BASE_URL%/api/beacon -H "Content-Type: application/json" -d "{\"name\":\"monitor-start\",\"metadata\":{\"note\":\"ETH monitoring started\",\"threshold\":2000}}" | jq -r ".status" 2>nul || echo beacon_set
echo [OK] Beacon set
echo.

REM Simulate recording data
echo Step 2: Recording sample ETH price data...
curl -s -X POST %BASE_URL%/api/tasks -H "Content-Type: application/json" -d "{\"tasks\":[{\"id\":\"eth-price-1\",\"type\":\"data_record\",\"prompt\":\"Record ETH price\",\"command\":\"echo '{\\\"timestamp\\\":\\\"2024-02-09T10:00:00Z\\\",\\\"price\\\":2300.50,\\\"currency\\\":\\\"USD\\\"}' >> %DATA_FILE%\",\"interval\":\"0\"}]}" >nul
echo [OK] Recorded sample 1: $2300.50
echo.

echo Step 3: Setting 'price-checkpoint' beacon...
curl -s -X POST %BASE_URL%/api/beacon -H "Content-Type: application/json" -d "{\"name\":\"price-checkpoint\",\"metadata\":{\"note\":\"Mid-monitoring checkpoint\"}}" | jq -r ".status" 2>nul || echo beacon_set
echo [OK] Checkpoint beacon set
echo.

echo Step 4: Recording more sample data...
curl -s -X POST %BASE_URL%/api/tasks -H "Content-Type: application/json" -d "{\"tasks\":[{\"id\":\"eth-price-2\",\"type\":\"data_record\",\"prompt\":\"Record ETH price\",\"command\":\"echo '{\\\"timestamp\\\":\\\"2024-02-09T10:30:00Z\\\",\\\"price\\\":2350.75,\\\"currency\\\":\\\"USD\\\"}' >> %DATA_FILE%\",\"interval\":\"0\"}]}" >nul
echo [OK] Recorded sample 2: $2350.75
echo.

echo ==========================================
echo Testing Beacon System
echo ==========================================
echo.

echo All beacons:
curl -s %BASE_URL%/api/beacons | jq -r ".beacons[] | \"  - \\(.name): \\(.timestamp)\"" 2>nul || curl -s %BASE_URL%/api/beacons
echo.

echo Reading memory since 'monitor-start':
curl -s %BASE_URL%/api/memory?beacon=monitor-start | jq "." 2>nul || curl -s %BASE_URL%/api/memory?beacon=monitor-start
echo.

echo Reading memory since 'price-checkpoint':
curl -s %BASE_URL%/api/memory?beacon=price-checkpoint | jq "." 2>nul || curl -s %BASE_URL%/api/memory?beacon=price-checkpoint
echo.

echo ==========================================
echo Demo Complete!
echo ==========================================
echo.
echo Key Concepts Demonstrated:
echo   1. Beacon-based memory checkpoints
echo   2. Time-windowed queries
echo   3. Persistent storage in JSONL format
echo.
echo Next steps:
echo   - Start server: .\bin\cerebellum.exe
echo   - Set beacons at decision points
echo   - Query trends using /api/memory?beacon=NAME

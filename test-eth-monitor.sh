#!/bin/bash
# ETH Price Monitoring Demo for Cerebellum
# This script demonstrates the beacon-based memory system

set -e

BASE_URL="http://localhost:18080"
DATA_FILE="./data/eth_prices.jsonl"

echo "=========================================="
echo "ETH Price Monitor Demo"
echo "=========================================="
echo ""

# Check if server is running
if ! curl -s "$BASE_URL/health" > /dev/null; then
    echo "Error: Cerebellum server is not running on $BASE_URL"
    echo "Please start the server first: ./bin/cerebellum.exe"
    exit 1
fi

echo "✓ Server is running"
echo ""

# Function to fetch ETH price
fetch_eth_price() {
    local url="https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd&include_24hr_change=true"
    
    curl -s -X POST "$BASE_URL/execute" \
        -H "Content-Type: application/json" \
        -d "{\"url\":\"$url\",\"method\":\"GET\",\"headers\":{\"Accept\":\"application/json\"}}" \
        | jq -r '.body' \
        | jq -r '.ethereum.usd'
}

# Function to record price to memory
record_price() {
    local price="$1"
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    curl -s -X POST "$BASE_URL/api/tasks" \
        -H "Content-Type: application/json" \
        -d "{\"tasks\":[{\"id\":\"eth-price-$timestamp\",\"type\":\"data_record\",\"prompt\":\"Record ETH price: $price at $timestamp\",\"command\":\"echo '{\\\"timestamp\\\":\\\"$timestamp\\\",\\\"price\\\":$price,\\\"currency\\\":\\\"USD\\\"}' >> $DATA_FILE\",\"interval\":\"0\",\"description\":\"ETH Price Monitor\"}]}" > /dev/null
    
    echo "✓ Recorded ETH price: \$$price at $timestamp"
}

# Function to set beacon
set_beacon() {
    local name="$1"
    local metadata="$2"
    
    curl -s -X POST "$BASE_URL/api/beacon" \
        -H "Content-Type: application/json" \
        -d "{\"name\":\"$name\",\"metadata\":$metadata}" \
        | jq -r '.status'
}

# Function to read memory since beacon
read_memory_since() {
    local beacon="$1"
    
    curl -s "$BASE_URL/api/memory?beacon=$beacon"
}

# Function to list beacons
list_beacons() {
    curl -s "$BASE_URL/api/beacons" | jq -r '.beacons[] | "  - \(.name): \(.timestamp)"'
}

echo "Step 1: Setting initial beacon 'monitor-start'..."
set_beacon "monitor-start" '{"note":"ETH monitoring started","threshold":2000}'
echo "✓ Beacon set"
echo ""

echo "Step 2: Fetching and recording ETH price (Sample 1)..."
PRICE1=$(fetch_eth_price)
record_price "$PRICE1"
echo ""

echo "Step 3: Waiting 5 seconds before next sample..."
sleep 5
echo ""

echo "Step 4: Setting 'price-checkpoint' beacon..."
set_beacon "price-checkpoint" '{"note":"Mid-monitoring checkpoint"}'
echo "✓ Checkpoint beacon set"
echo ""

echo "Step 5: Fetching and recording ETH price (Sample 2)..."
PRICE2=$(fetch_eth_price)
record_price "$PRICE2"
echo ""

echo "Step 6: Waiting 5 seconds..."
sleep 5
echo ""

echo "Step 7: Fetching and recording ETH price (Sample 3)..."
PRICE3=$(fetch_eth_price)
record_price "$PRICE3"
echo ""

echo "=========================================="
echo "Beacon System Demo"
echo "=========================================="
echo ""

echo "All beacons set:"
list_beacons
echo ""

echo "Reading memory since 'monitor-start':"
read_memory_since "monitor-start" | jq '.'
echo ""

echo "Reading memory since 'price-checkpoint':"
read_memory_since "price-checkpoint" | jq '.'
echo ""

echo "=========================================="
echo "Trend Analysis"
echo "=========================================="
echo ""

# Calculate trend
if command -v bc >/dev/null 2>&1; then
    CHANGE=$(echo "scale=2; $PRICE3 - $PRICE1" | bc)
    if (( $(echo "$CHANGE > 0" | bc -l) )); then
        echo "📈 Price trend: UP by \$$CHANGE"
    elif (( $(echo "$CHANGE < 0" | bc -l) )); then
        echo "📉 Price trend: DOWN by \$$CHANGE"
    else
        echo "➡️  Price trend: STABLE"
    fi
else
    echo "Price 1: \$$PRICE1"
    echo "Price 2: \$$PRICE2"
    echo "Price 3: \$$PRICE3"
fi

echo ""
echo "=========================================="
echo "Demo Complete!"
echo "=========================================="
echo ""
echo "Key Concepts Demonstrated:"
echo "  1. Periodic data collection (ETH prices)"
echo "  2. Beacon-based memory checkpoints"
echo "  3. Time-windowed queries"
echo "  4. Trend analysis from stored data"
echo ""
echo "Next steps:"
echo "  - Brain can set beacons at decision points"
echo "  - Query trends: 'What happened since beacon X?'"
echo "  - Historical analysis: 'Compare period A vs period B'"

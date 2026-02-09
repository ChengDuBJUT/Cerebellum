#!/bin/bash
# Import local qwen2:0.5b model into Ollama
# Usage: ./import-model.sh

echo "=========================================="
echo "Cerebellum Local Model Import Script"
echo "=========================================="
echo ""

# Check if Ollama is installed
if ! command -v ollama &> /dev/null; then
    echo "❌ Error: Ollama is not installed!"
    echo ""
    echo "Please install Ollama first:"
    echo "  macOS/Linux: curl -fsSL https://ollama.com/install.sh | sh"
    echo "  Windows: https://ollama.com/download/windows"
    exit 1
fi

echo "✓ Ollama found"

# Check if model files exist
MODEL_DIR="$(dirname "$0")/models"
if [ ! -d "$MODEL_DIR" ]; then
    echo "❌ Error: models directory not found!"
    echo "Expected location: $MODEL_DIR"
    exit 1
fi

echo "✓ Model directory found"

# Check for required files
REQUIRED_FILES=(
    "blobs/sha256-8de95da68dc485c0889c205384c24642f83ca18d089559c977ffc6a3972a71a8"
    "manifests/registry.ollama.ai/library/qwen2/0.5b"
)

for file in "${REQUIRED_FILES[@]}"; do
    if [ ! -f "$MODEL_DIR/$file" ]; then
        echo "❌ Error: Required file not found: $file"
        exit 1
    fi
done

echo "✓ All required files present"
echo ""

# Check if Ollama is running
if ! curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
    echo "⚠️  Ollama server is not running!"
    echo "Starting Ollama server..."
    ollama serve &
    sleep 5
    
    if ! curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
        echo "❌ Failed to start Ollama server"
        echo "Please start it manually: ollama serve"
        exit 1
    fi
    echo "✓ Ollama server started"
else
    echo "✓ Ollama server is running"
fi

echo ""
echo "Importing qwen2:0.5b model..."
echo "This may take a moment..."
echo ""

# Import the model
cd "$(dirname "$0")"
if ollama create qwen2:0.5b -f Modelfile; then
    echo ""
    echo "=========================================="
    echo "✅ Model imported successfully!"
    echo "=========================================="
    echo ""
    echo "Model name: qwen2:0.5b"
    echo "You can now use it with Cerebellum"
    echo ""
    echo "Verify installation:"
    echo "  ollama list"
    echo ""
    echo "Test the model:"
    echo '  curl http://localhost:11434/api/generate -d '\''{"model": "qwen2:0.5b", "prompt": "Hello"}'\'''
    echo ""
else
    echo ""
    echo "=========================================="
    echo "❌ Failed to import model"
    echo "=========================================="
    echo ""
    echo "Please check:"
    echo "1. Ollama is running correctly"
    echo "2. Model files are not corrupted"
    echo "3. You have sufficient disk space"
    exit 1
fi

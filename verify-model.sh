#!/bin/bash
# Verify model files integrity and size

echo "=========================================="
echo "Cerebellum Model Verification"
echo "=========================================="
echo ""

MODEL_DIR="$(dirname "$0")/models"

# Expected file sizes (in bytes)
declare -A EXPECTED_SIZES=(
    ["blobs/sha256-2184ab82477bc33a5e08fa209df88f0631a19e686320cce2cfe9e00695b2f0e6"]=488
    ["blobs/sha256-8de95da68dc485c0889c205384c24642f83ca18d089559c977ffc6a3972a71a8"]=352151968
    ["blobs/sha256-62fbfd9ed093d6e5ac83190c86eec5369317919f4b149598d2dbb38900e9faef"]=182
    ["blobs/sha256-c156170b718ec29139d3653d40ed1986fd92fb7e0959b5c71f3c48f62e6636f4"]=11344
    ["blobs/sha256-f02dd72bb2423204352eabc5637b44d79d17f109fdb510a7c51455892aa2d216"]=59
    ["manifests/registry.ollama.ai/library/qwen2/0.5b"]=856
)

ERRORS=0
TOTAL_SIZE=0

echo "Checking model files..."
echo ""

for file in "${!EXPECTED_SIZES[@]}"; do
    filepath="$MODEL_DIR/$file"
    expected_size=${EXPECTED_SIZES[$file]}
    
    if [ ! -f "$filepath" ]; then
        echo "❌ Missing: $file"
        ERRORS=$((ERRORS + 1))
        continue
    fi
    
    actual_size=$(stat -f%z "$filepath" 2>/dev/null || stat -c%s "$filepath" 2>/dev/null)
    TOTAL_SIZE=$((TOTAL_SIZE + actual_size))
    
    if [ "$actual_size" -eq "$expected_size" ]; then
        echo "✓ $file ($(numfmt --to=iec $actual_size))"
    else
        echo "⚠️  $file (expected: $expected_size, got: $actual_size)"
        ERRORS=$((ERRORS + 1))
    fi
done

echo ""
echo "=========================================="

if [ $ERRORS -eq 0 ]; then
    echo "✅ All files verified successfully!"
    echo ""
    echo "Total size: $(numfmt --to=iec $TOTAL_SIZE)"
    echo "Model: qwen2:0.5b (Qwen2 0.5B Instruct)"
    echo ""
    echo "You can now import the model:"
    echo "  ./import-model.sh"
else
    echo "❌ Found $ERRORS issue(s)"
    echo ""
    echo "Please ensure all model files are present and not corrupted."
    exit 1
fi

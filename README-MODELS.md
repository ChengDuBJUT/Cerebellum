# Cerebellum Local Models Setup

This guide explains how to set up Cerebellum to use the included local LLM model (qwen2:0.5b) instead of downloading from Ollama's registry.

## Included Model

**Model**: qwen2:0.5b (Qwen2 0.5B Instruct)
**Size**: 352 MB
**Location**: `models/` directory
**Purpose**: Local LLM inference for Cerebellum

## Quick Setup (Method 1 - Recommended)

### Option A: Use Ollama with Local Model

1. **Install Ollama** (if not already installed)
   ```bash
   # macOS/Linux
   curl -fsSL https://ollama.com/install.sh | sh
   
   # Windows
   # Download from https://ollama.com/download/windows
   ```

2. **Import the local model into Ollama**
   
   Create a Modelfile in the cerebellum directory:
   ```dockerfile
   FROM ./models/blobs/sha256-8de95da68dc485c0889c205384c24642f83ca18d089559c977ffc6a3972a71a8
   ```
   
   Then import:
   ```bash
   ollama create qwen2:0.5b-local -f Modelfile
   ```

3. **Update cerebellum.yaml**
   ```yaml
   ollama:
     host: "http://localhost:11434"
     model: "qwen2:0.5b-local"  # Use the local model name
   ```

### Option B: Direct Binary Usage (Advanced)

For advanced users who want to use the model directly without Ollama:

The model files are in GGUF format and can be used with:
- llama.cpp
- Other GGUF-compatible inference engines

However, Cerebellum is designed to work with Ollama for simplicity.

## Method 2: Offline/Air-Gapped Environment

If you're in an environment without internet access:

1. **Copy the models directory to your target machine**
   ```bash
   # On source machine
   tar -czf cerebellum-models.tar.gz models/
   
   # On target machine
   tar -xzf cerebellum-models.tar.gz
   ```

2. **Set OLLAMA_MODELS environment variable**
   ```bash
   # Linux/macOS
   export OLLAMA_MODELS=/path/to/cerebellum/models
   
   # Windows
   set OLLAMA_MODELS=C:\path\to\cerebellum\models
   ```

3. **Start Ollama**
   ```bash
   ollama serve
   ```

4. **Verify model is available**
   ```bash
   ollama list
   # Should show: qwen2:0.5b
   ```

## Configuration

### cerebellum.yaml

```yaml
server:
  host: "0.0.0.0"
  port: 18080

ollama:
  host: "http://localhost:11434"
  model: "qwen2:0.5b"  # This matches the manifest name

watcher:
  poll_interval: 1000
```

### Environment Variables

- `OLLAMA_HOST`: Ollama server URL (default: http://localhost:11434)
- `OLLAMA_MODELS`: Custom models directory path
- `CEREBELLUM_MODEL`: Override model name in config

## Verification

1. **Check Ollama is running**
   ```bash
   curl http://localhost:11434/api/tags
   ```

2. **Test the model**
   ```bash
   curl http://localhost:11434/api/generate -d '{
     "model": "qwen2:0.5b",
     "prompt": "Hello!"
   }'
   ```

3. **Start Cerebellum**
   ```bash
   ./bin/cerebellum
   ```

4. **Verify Cerebellum health**
   ```bash
   curl http://localhost:18080/health
   ```

## Model Details

### qwen2:0.5b Specifications

- **Architecture**: Transformer-based LLM
- **Parameters**: 0.5 billion (494M)
- **Context Length**: 32,768 tokens
- **Quantization**: Q4_0 (4-bit)
- **License**: Apache 2.0
- **Language**: Multilingual (optimized for Chinese and English)

### Performance Characteristics

- **Memory Usage**: ~500MB RAM
- **Inference Speed**: 50-200ms per token on CPU
- **Response Quality**: Good for task execution, monitoring, and filtering
- **Use Case**: Perfect for Cerebellum's local execution layer

## Troubleshooting

### Model not found
```bash
# Ensure Ollama can see the model
ollama list

# If not visible, you may need to import it
ollama create qwen2:0.5b -f Modelfile
```

### Permission denied
```bash
# Fix permissions
chmod -R 755 models/
```

### Out of memory
- The 0.5B model requires ~500MB RAM
- Ensure your system has at least 1GB free RAM
- Close other applications if needed

### Model loading slowly
- First load may take 10-30 seconds
- Subsequent loads are faster
- Model stays loaded in memory while Ollama is running

## Alternative Models

If you want to use a different model with Cerebellum:

1. **Download via Ollama**
   ```bash
   ollama pull llama3.2:1b  # 1B parameter model
   ollama pull phi3:mini    # Microsoft's Phi-3 mini
   ```

2. **Update cerebellum.yaml**
   ```yaml
   ollama:
     model: "llama3.2:1b"
   ```

3. **Restart Cerebellum**
   ```bash
   # Stop and start
   ./bin/cerebellum
   ```

## Migration Notes

When migrating Cerebellum to a new system:

1. ✅ Copy the entire project directory
2. ✅ Models are included in `models/` directory
3. ✅ No need to re-download from Ollama
4. ✅ Just start Ollama and Cerebellum on the new system

## File Structure

```
cerebellum/
├── bin/
│   └── cerebellum.exe          # Main executable
├── models/
│   ├── MODEL-INFO.md           # This file
│   ├── manifests/              # Ollama model manifests
│   │   └── registry.ollama.ai/
│   │       └── library/
│   │           └── qwen2/
│   │               └── 0.5b    # Model manifest
│   └── blobs/                  # Model binary files
│       ├── sha256-...          # Model weights (336MB)
│       └── sha256-...          # Config, template, etc.
├── cerebellum.yaml             # Configuration
└── ...
```

## Size Summary

- **Model weights**: 336 MB
- **Config & metadata**: ~15 KB
- **Total**: ~337 MB

This is small enough to:
- ✅ Include in Git LFS
- ✅ Bundle in Docker images
- ✅ Distribute via USB drive
- ✅ Run on modest hardware

## Support

For issues with the local model setup:
1. Check Ollama documentation: https://github.com/ollama/ollama
2. Verify model files are not corrupted
3. Ensure sufficient disk space and RAM
4. Check Cerebellum logs for specific errors

## License

The qwen2:0.5b model is licensed under Apache 2.0.
See models/blobs/sha256-c156170b718ec29139d3653d40ed1986fd92fb7e0959b5c71f3c48f62e6636f4 for full license text.

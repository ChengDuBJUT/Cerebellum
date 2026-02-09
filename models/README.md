# Models Directory

## Important Note

Due to GitHub's 100MB file size limit, model binary files are not included in this repository.

## Download Model

The qwen2:0.5b model (336MB) can be downloaded from:
- Ollama Hub: `ollama pull qwen2:0.5b`
- Direct download: Contact repository owner
- Alternative: Use the Modelfile to build from weights

## Manual Setup

1. Install Ollama: https://ollama.com
2. Download model:
   ```bash
   ollama pull qwen2:0.5b
   ```
3. Or use the import script after obtaining model files:
   ```bash
   ./import-model.sh
   ```

## Model Information

- **Model**: qwen2:0.5b (Qwen2 0.5B Instruct)
- **Size**: 336 MB
- **Format**: GGUF (Ollama compatible)
- **License**: Apache 2.0

## Alternative Models

You can use any Ollama-compatible model by updating `cerebellum.yaml`:
```yaml
ollama:
  model: "llama3.2:1b"  # or any other model
```

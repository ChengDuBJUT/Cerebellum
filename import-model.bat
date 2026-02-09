@echo off
chcp 65001 >nul
echo ==========================================
echo Cerebellum Local Model Import Script
echo ==========================================
echo.

REM Check if Ollama is installed
where ollama >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Error: Ollama is not installed!
    echo.
    echo Please install Ollama first:
    echo   https://ollama.com/download/windows
    pause
    exit /b 1
)

echo ✓ Ollama found

REM Check if model files exist
set "MODEL_DIR=%~dp0models"
if not exist "%MODEL_DIR%" (
    echo ❌ Error: models directory not found!
    echo Expected location: %MODEL_DIR%
    pause
    exit /b 1
)

echo ✓ Model directory found

REM Check for required files
if not exist "%MODEL_DIR%\blobs\sha256-8de95da68dc485c0889c205384c24642f83ca18d089559c977ffc6a3972a71a8" (
    echo ❌ Error: Model weights file not found!
    pause
    exit /b 1
)

if not exist "%MODEL_DIR%\manifests\registry.ollama.ai\library\qwen2\0.5b" (
    echo ❌ Error: Model manifest not found!
    pause
    exit /b 1
)

echo ✓ All required files present
echo.

REM Check if Ollama is running
curl -s http://localhost:11434/api/tags >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  Ollama server is not running!
    echo Starting Ollama server...
    start /b ollama serve
    timeout /t 5 /nobreak >nul
    
    curl -s http://localhost:11434/api/tags >nul 2>&1
    if %errorlevel% neq 0 (
        echo ❌ Failed to start Ollama server
        echo Please start it manually: ollama serve
        pause
        exit /b 1
    )
    echo ✓ Ollama server started
) else (
    echo ✓ Ollama server is running
)

echo.
echo Importing qwen2:0.5b model...
echo This may take a moment...
echo.

REM Import the model
cd /d "%~dp0"
ollama create qwen2:0.5b -f Modelfile
if %errorlevel% equ 0 (
    echo.
    echo ==========================================
    echo ✅ Model imported successfully!
    echo ==========================================
    echo.
    echo Model name: qwen2:0.5b
    echo You can now use it with Cerebellum
    echo.
    echo Verify installation:
    echo   ollama list
    echo.
    echo Test the model:
    echo   curl http://localhost:11434/api/generate -d "{\"model\": \"qwen2:0.5b\", \"prompt\": \"Hello\"}"
    echo.
) else (
    echo.
    echo ==========================================
    echo ❌ Failed to import model
    echo ==========================================
    echo.
    echo Please check:
    echo 1. Ollama is running correctly
    echo 2. Model files are not corrupted
    echo 3. You have sufficient disk space
)

pause

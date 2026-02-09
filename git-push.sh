#!/bin/bash
# Cerebellum Git 提交脚本
# 在 cerebellum 目录下运行此脚本

cd "C:\Users\Administrator\AppData\Roaming\npm\node_modules\openclaw\extensions\cerebellum"

# 1. 初始化 Git 仓库
echo "1. 初始化 Git 仓库..."
git init

# 2. 添加所有文件
echo "2. 添加文件到暂存区..."
git add .

# 3. 创建提交
echo "3. 创建提交..."
git commit -m "Initial commit: Cerebellum AI Agent Subsystem

Features:
- Two-tier architecture (Brain + Cerebellum)
- 80-95% API cost reduction
- Local LLM integration (Ollama qwen2:0.5b)
- Offline-capable with bundled model
- HTTP REST API for universal compatibility
- Support for OpenCode, Claude Code, OpenClaw

Includes:
- Complete Go source code
- Pre-built binaries
- Local LLM model (336MB)
- Import scripts (Windows/Linux/macOS)
- Comprehensive documentation"

# 4. 添加远程仓库
echo "4. 添加远程仓库..."
git remote add origin https://github.com/ChengDuBJUT/Cerebellum.git

# 5. 检查远程
echo "5. 检查远程仓库..."
git remote -v

echo ""
echo "=========================================="
echo "准备完成！"
echo "=========================================="
echo ""
echo "下一步："
echo ""
echo "如果你有仓库写入权限，运行："
echo "  git push -u origin main"
echo ""
echo "如果需要身份验证，会提示输入："
echo "  Username: 你的GitHub用户名"
echo "  Password: 你的GitHub Personal Access Token"
echo ""
echo "如果没有权限，可以先fork仓库到自己的账户，然后："
echo "  git remote set-url origin https://github.com/YOUR_USERNAME/Cerebellum.git"
echo "  git push -u origin main"
echo ""

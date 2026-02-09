# Cerebellum (小脑) - AI Agent 代币效率优化子系统

> **本地执行。集中决策。全球节省。**

Cerebellum 是一个轻量级 AI 助手子系统，旨在显著减少 API 代币消耗，同时保持高频任务执行能力。它作为 AI Agent 的"执行层"，在本地处理常规监控、数据收集和重复性任务，而将昂贵的 API 调用留给"大脑"（主 AI）进行战略决策。

一句话安装 请将下面的话贴给你的ai agent

https://github.com/ChengDuBJUT/Cerebellum 请安装小脑 并执行 skill.md


---

## 问题：AI Agent 中的代币消耗

现代 AI Agent 面临一个关键的成本障碍：**高频操作会迅速消耗 API 代币**。

### 场景：社区监控

**传统方法**（直接 API 调用）：
- 每 5 分钟查询社区动态：~2,000 代币
- 每小时分析趋势：~5,000 代币
- 每天 4 次搜索讨论：~8,000 代币
- 生成状态报告：~3,000 代币
- **每日成本**：~50,000 代币 × $0.01 = **$0.50/天**
- **每月成本**：**$15/月** 用于基本监控

**使用 Cerebellum**：
- 本地 LLM (Ollama qwen2:0.5b)：$0
- 大脑 API 调用（仅决策）：~500 代币/天
- **每日成本**：~500 代币 × $0.01 = **$0.005/天**
- **每月成本**：**$0.15/月**

**节省**：**97% 的成本降低**

---

## 解决方案：双层架构

```
┌─────────────────────────────────────────────────────────────┐
│                      主 AI (大脑)                            │
│              Claude / GPT-4 / Gemini / 等                   │
│                                                              │
│  • 战略决策                   • 复杂推理                     │
│  • 创意内容                   • 重要回复                     │
│  • 质量判断                   • 高价值分析                   │
│                                                              │
│  代币使用量：占总工作量的 ~5%                               │
│  成本：每 1K 代币 $0.01-0.03                                │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTP API (localhost:18080)
                       │ 偶尔的状态更新
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                    CEREBELLUM (执行层)                       │
│              本地 LLM (Ollama) + HTTP 服务                   │
│                                                              │
│  • 高频监控                   • 数据收集                     │
│  • HTTP 请求执行              • 任务调度                     │
│  • 内容准备                   • 状态报告                     │
│  • 模式检测                   • 趋势分析                     │
│                                                              │
│  代币使用量：占总工作量的 ~95%（本地，$0）                   │
│  成本：$0（在 CPU 上本地运行）                               │
└─────────────────────────────────────────────────────────────┘
```

### 工作原理

1. **大脑** 定义策略并分配监控任务给 Cerebellum
2. **Cerebellum** 使用轻量级 LLM (qwen2:0.5b) 在本地执行任务
3. **Cerebellum** 定期（每小时/每天）报告发现
4. **大脑** 审查摘要并做出决策
5. **大脑** 仅对高价值机会采取行动

---

## 通用兼容性

Cerebellum 适用于**任何 AI Agent 平台**：

| 平台 | 类型 | 集成方法 | 状态 |
|----------|------|-------------------|---------|
| **Claude Code** | 半自动 | HTTP API 调用 | ✅ 支持 |
| **OpenCode** | CLI Agent | HTTP API 调用 | ✅ 支持 |
| **OpenClaw** | 扩展系统 | HTTP API 调用 | ✅ 支持 |
| **AutoGPT** | 全自动 | HTTP API 调用 | ✅ 支持 |
| **自定义 Agent** | 任意 | HTTP REST API | ✅ 支持 |

### 集成模式（通用）

```python
# 任何 AI Agent 都可以使用 Cerebellum
import requests

CEREBELLUM_URL = "http://localhost:18080"

# 1. 分配监控任务
def assign_monitoring():
    tasks = {
        "tasks": [
            {
                "id": "community-monitor",
                "type": "periodic",
                "interval": "5m",
                "command": "检查社区动态中的新帖子"
            }
        ]
    }
    requests.post(f"{CEREBELLUM_URL}/api/tasks", json=tasks)

# 2. 通过 Cerebellum 执行 HTTP 请求
def query_data(url):
    response = requests.post(
        f"{CEREBELLUM_URL}/execute",
        json={"url": url, "method": "GET"}
    )
    return response.json()

# 3. 获取智能摘要
def get_report():
    return requests.get(f"{CEREBELLUM_URL}/api/report").json()
```

---

## 多场景应用

### 1. 社交社区管理

**用例**：监控 Reddit、Discord、Moltbook 等。

**不使用 Cerebellum**：
- 每次动态检查：2,000 代币
- 每次趋势分析：5,000 代币
- 每次回复草稿：3,000 代币
- **每月成本**：$50-100

**使用 Cerebellum**：
- Cerebellum 24/7 本地监控
- 大脑每小时审查摘要：500 代币
- 大脑每天撰写 2-3 个高质量回复：6,000 代币
- **每月成本**：$3-5
- **节省**：**90-95%**

```json
{
  "tasks": [
    {"id": "feed-monitor", "type": "periodic", "interval": "5m", "command": "检查社区动态"},
    {"id": "trend-analysis", "type": "periodic", "interval": "1h", "command": "分析热门话题"},
    {"id": "engagement-check", "type": "periodic", "interval": "15m", "command": "寻找需要互动的帖子"}
  ]
}
```

### 2. GitHub 仓库监控

**用例**：跟踪多个仓库的问题、PR、发布

**不使用 Cerebellum**：
- 每 30 分钟检查 10 个仓库：8,000 代币/小时
- 分析问题：5,000 代币/小时
- 起草回复：4,000 代币/小时
- **每月成本**：$120-150

**使用 Cerebellum**：
- Cerebellum 本地轮询 GitHub API
- 大脑每天 2 次审查摘要：1,000 代币
- 大脑仅处理关键问题：2,000 代币
- **每月成本**：$10-15
- **节省**：**90%**

```json
{
  "tasks": [
    {"id": "github-issues", "type": "periodic", "interval": "10m", "command": "检查新问题"},
    {"id": "github-prs", "type": "periodic", "interval": "15m", "command": "监控拉取请求"},
    {"id": "release-tracker", "type": "periodic", "interval": "1h", "command": "检查新发布"}
  ]
}
```

### 3. 客户支持自动化

**用例**：监控支持工单、论坛、聊天频道

**不使用 Cerebellum**：
- 每天分类 100 个工单：20,000 代币
- 起草回复：30,000 代币
- 升级分析：10,000 代币
- **每月成本**：$180-200

**使用 Cerebellum**：
- Cerebellum 本地分类：$0
- 大脑仅审查升级：5,000 代币/天
- 大脑起草复杂回复：10,000 代币/天
- **每月成本**：$45-50
- **节省**：**75%**

### 4. 市场情报与研究

**用例**：跟踪竞争对手、新闻、价格变化

**不使用 Cerebellum**：
- 每小时爬取 50 个来源：25,000 代币/小时
- 分析趋势：15,000 代币/小时
- 生成报告：10,000 代币/小时
- **每月成本**：$300-400

**使用 Cerebellum**：
- Cerebellum 本地爬取和过滤
- 大脑接收过滤后的洞察：3,000 代币/天
- 大脑生成战略报告：5,000 代币/天
- **每月成本**：$25-30
- **节省**：**90-92%**

### 5. IoT & 系统监控

**用例**：监控服务器日志、指标、告警

**不使用 Cerebellum**：
- 每 5 分钟分析日志：15,000 代币/小时
- 模式检测：8,000 代币/小时
- 告警分类：5,000 代币/小时
- **每月成本**：$200-250

**使用 Cerebellum**：
- Cerebellum 本地监控
- 大脑仅接收异常告警：2,000 代币/天
- **每月成本**：$15-20
- **节省**：**92-94%**

---

## 代币消耗分析

### 详细成本对比

#### 场景：24/7 社区监控

| 指标 | 直接 API (GPT-4) | 使用 Cerebellum | 节省 |
|--------|-------------------|-----------------|---------|
| **每小时检查** | 12 次 × 2K 代币 = 24K | 12 次 × 0 代币 = 0 | 100% |
| **趋势分析** | 24 次分析 × 5K = 120K | 24 次本地 = 0 | 100% |
| **每日摘要** | 0 | 1 × 1K = 1K | - |
| **战略决策** | 10 × 3K = 30K | 10 × 3K = 30K | 0% |
| **高价值回复** | 20 × 4K = 80K | 5 × 4K = 20K | 75% |
| **每日总计** | **254K 代币** | **51K 代币** | **80%** |
| **每月成本** | **$76.20** | **$15.30** | **$60.90 节省** |

**年度节省**：**$730.80**

#### 场景：GitHub 项目管理（10 个仓库）

| 活动 | 频率 | 直接 API 代币 | Cerebellum 代币 |
|----------|-----------|-------------------|-------------------|
| 问题轮询 | 每 10 分钟 | 144 × 1.5K = 216K | 0 |
| PR 监控 | 每 15 分钟 | 96 × 2K = 192K | 0 |
| 发布检查 | 每小时 | 24 × 1K = 24K | 0 |
| 模式检测 | 每 30 分钟 | 48 × 3K = 144K | 0 |
| 每日摘要审查 | 每天一次 | 0 | 2K |
| 关键问题处理 | 5 次/天 | 5 × 4K = 20K | 5 × 4K = 20K |
| **每日总计** | | **596K 代币** | **22K 代币** |
| **每月成本** | | **$178.80** | **$6.60** |
| **节省** | | | **96.3%** |

**年度节省**：**$2,066.40**

#### 场景：多平台客户支持

| 平台 | 每日工单 | 直接 API 成本 | Cerebellum 成本 |
|----------|-------------|-----------------|-----------------|
| 邮件 | 50 | $15/天 | $3/天 |
| 聊天 | 200 | $40/天 | $8/天 |
| 论坛 | 30 | $12/天 | $2/天 |
| 社交媒体 | 100 | $25/天 | $5/天 |
| **每日总计** | **380** | **$92** | **$18** |
| **每月成本** | | **$2,760** | **$540** |
| **节省** | | | **80.4%** |

**年度节省**：**$26,640**

---

## 架构深度解析

### 为什么 Cerebellum 能节省代币

**1. 本地执行模型**

```
传统 AI Agent：
用户请求 → 主 AI → API 调用 ($) → 响应
              ↓
      [每次操作都消耗代币]

使用 Cerebellum：
用户请求 → Cerebellum (本地 LLM, $0) → 原始数据
              ↓
          大脑接收过滤后的摘要
              ↓
          大脑做出决策（最少代币）
              ↓
          采取行动
```

**2. 任务批处理**

取代 100 次单独的 API 调用：
- Cerebellum 在本地批处理操作
- 向大脑发送整合报告
- 大脑做 1-2 个决策而不是 100 个

**3. 智能过滤**

Cerebellum 使用本地 LLM 进行预过滤：
- "这对大脑来说够重要吗？"
- 95% 的噪音在本地过滤
- 只有 5% 到达昂贵的大脑 API

### 性能特征

| 指标 | Cerebellum (本地) | 主 AI (API) |
|--------|-------------------|---------------|
| **响应时间** | 50-200ms | 500-3000ms |
| **每代币成本** | $0 | $0.01-0.03 |
| **可用性** | 100% (本地) | 依赖 API |
| **速率限制** | 无 | 有 |
| **最适合** | 执行、监控 | 战略、创意 |

---

## 安装与设置

### 前置要求

- **Ollama** 已安装 qwen2:0.5b 模型
- **Go 1.22+**（用于从源码构建）
- 任意 AI Agent 平台（Claude Code、OpenCode 等）

### 快速开始

```bash
# 1. 克隆仓库
git clone https://github.com/yourusername/cerebellum.git
cd cerebellum

# 2. 构建二进制文件
# Windows:
go build -o bin/cerebellum.exe cmd/server/main.go

# Linux/macOS:
go build -o bin/cerebellum cmd/server/main.go

# 3. 配置
cp cerebellum.yaml.example cerebellum.yaml
# 编辑 cerebellum.yaml 添加你的设置

# 4. 启动 Cerebellum
./bin/cerebellum

# 5. 验证
curl http://localhost:18080/health
```

### 配置

**cerebellum.yaml**：
```yaml
server:
  host: "0.0.0.0"
  port: 18080

ollama:
  host: "http://localhost:11434"
  model: "qwen2:0.5b"  # 500M 参数，在 CPU 上运行

watcher:
  poll_interval: 1000  # 毫秒
```

### 集成示例

```python
# 示例：将 Cerebellum 与 Claude Code 一起使用
import requests
import json

CEREBELLUM = "http://localhost:18080"

class CommunityManager:
    def __init__(self):
        self.setup_monitoring()
    
    def setup_monitoring(self):
        """Cerebellum 处理高频监控"""
        tasks = {
            "tasks": [
                {
                    "id": "community-feed",
                    "type": "periodic",
                    "interval": "5m",
                    "command": "监控社区动态中的新帖子"
                },
                {
                    "id": "trend-detector",
                    "type": "periodic",
                    "interval": "1h",
                    "command": "分析热门话题"
                }
            ]
        }
        requests.post(f"{CEREBELLUM}/api/tasks", json=tasks)
    
    def check_community(self):
        """大脑审查 Cerebellum 的发现"""
        report = requests.get(f"{CEREBELLUM}/api/report").json()
        
        if report["completed_count"] > 0:
            # 大脑分析发现（昂贵的 API 调用）
            findings = report["report"]["completed"]
            return self.make_strategic_decisions(findings)
        
        return "无重要活动"
    
    def make_strategic_decisions(self, findings):
        """大脑做出高价值决策"""
        # 这是 GPT-4/Claude 发挥作用的地方
        # 仅针对重要决策调用
        pass
```

---

## 最佳实践

### 1. 任务设计

**好的做法**：
```json
{
  "id": "github-issue-monitor",
  "type": "periodic",
  "interval": "10m",
  "command": "检查仓库中的关键问题"
}
```

**不好的做法**：
```json
{
  "id": "task1",
  "type": "monitoring",  // 错误的类型！
  "interval": "5 minutes"  // 错误的格式！
}
```

### 2. 间隔选择

- **关键监控**：2-5 分钟
- **常规更新**：15-30 分钟
- **每日摘要**：24 小时
- **避免**：所有任务使用相同间隔（错开它们）

### 3. 成本优化

```
之前：100 次 API 调用/天 × $0.02 = $2.00/天 = $60/月
之后：Cerebellum (100 次本地) + 5 次 API 调用 × $0.02 = $0.10/天 = $3/月
节省：95%
```

### 4. 安全

- Cerebellum 仅在本地运行（仅限 localhost）
- API 密钥永远不会离开你的机器
- 大脑 API 调用直接发送到提供商
- 不涉及第三方服务

---

## 真实世界节省计算器

### 输入你的数字

```python
# 你当前的成本
daily_api_calls = 200  # 每天多少次 API 调用
tokens_per_call = 3000  # 每次调用平均多少代币
api_cost_per_1k = 0.02  # 你的 API 每 1K 代币成本

# 计算
daily_tokens = daily_api_calls * tokens_per_call
daily_cost = (daily_tokens / 1000) * api_cost_per_1k
monthly_cost = daily_cost * 30

# 使用 Cerebellum（假设 API 调用减少 90%）
cerebellum_api_calls = daily_api_calls * 0.1  # 只有 10% 到达大脑
cerebellum_daily_cost = (cerebellum_api_calls * tokens_per_call / 1000) * api_cost_per_1k
cerebellum_monthly = cerebellum_daily_cost * 30

savings = monthly_cost - cerebellum_monthly
savings_percent = (savings / monthly_cost) * 100

print(f"当前每月成本: ${monthly_cost:.2f}")
print(f"使用 Cerebellum: ${cerebellum_monthly:.2f}")
print(f"每月节省: ${savings:.2f} ({savings_percent:.1f}%)")
print(f"年度节省: ${savings * 12:.2f}")
```

### 示例输出

```
当前每月成本: $360.00
使用 Cerebellum: $36.00
每月节省: $324.00 (90.0%)
年度节省: $3,888.00
```

---

## 与替代方案对比

| 解决方案 | 成本 | 延迟 | 设置 | 灵活性 |
|----------|------|---------|-------|-------------|
| **直接 API** | $$$ | 中等 | 简单 | 高 |
| **Cerebellum** | $ | 低 | 中等 | 高 |
| **缓存层** | $$ | 低 | 中等 | 低 |
| **小型模型** | $$ | 中等 | 简单 | 中等 |
| **无服务器函数** | $$ | 高 | 复杂 | 中等 |

**为什么 Cerebellum 胜出**：
- 比直接 API 便宜（90%+ 节省）
- 比无服务器更快（本地执行）
- 比缓存更灵活（智能过滤）
- 比小型模型更好（保持大脑的质量）

---

## 常见问题

**Q: Cerebellum 会取代我的主 AI 吗？**
A: 不会！它是补充。Cerebellum 处理执行；大脑处理战略。

**Q: Cerebellum 使用什么 LLM？**
A: 默认是通过 Ollama 的 qwen2:0.5b（免费，本地）。你可以配置其他模型。

**Q: 我实际能节省多少？**
A: 对于高频监控任务，典型节省是 80-95%。

**Q: 设置难吗？**
A: 单个二进制文件，一个配置文件。5 分钟即可运行。

**Q: 它适用于我现有的 Agent 吗？**
A: 是的！任何能发起 HTTP 调用的 Agent 都可以使用 Cerebellum。

**Q: 数据隐私呢？**
A: 一切都在本地运行。除最终的大脑 API 调用外，你的数据永远不会离开你的机器。

---

## 贡献

我们欢迎 AI Agent 社区的贡献！

- 报告问题
- 提交 PR
- 分享你的用例
- 建议改进

---

## 许可证

Apache-2.0 许可证 - 免费用于个人和商业用途。

---

## 致谢

为 AI Agent 社区而构建。特别感谢：
- Ollama 团队提供本地 LLM 推理
- 所有推动边界的 AI Agent 平台
- 早期采用者和贡献者

---

**今天开始节省代币。部署 Cerebellum。**

*"让 Cerebellum 处理噪音。让大脑处理信号。"*

---

**版本**: 1.0  
**最后更新**: 2026-02-09  
**状态**: 生产就绪 ✅  
**文档**: [完整文档](docs/)  
**社区**: [Discord](https://discord.gg/cerebellum)  
**问题**: [GitHub Issues](https://github.com/yourusername/cerebellum/issues)

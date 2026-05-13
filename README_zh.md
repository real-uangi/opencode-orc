中文 | [English](./README.md)

# opencode-orc

[opencode](https://github.com/anomalyco/opencode) 的轻量编排器。它以子进程方式运行 `opencode run --format json`，解析 JSONL 事件流，过滤并转换事件，最终输出简化的文本或 JSONL 流。

## 为什么需要？

Claude、Codex 等 AI 往往在官方提供的工具中才能发挥最大作用，但这些工具不支持外部的模型供应商。opencode 虽然支持外部供应商，却不支持为不同 Agent 指定不同模型（oh-my-opencode 解决了这个问题，但侵入性太大）。本项目的方案恰好避开了这些痛点——它保持轻量、零侵入，同时解锁了跨供应商的 Agent 委派能力。

opencode 支持绝大部分 AI 提供商，是最适合作为委派子 Agents 的桥接层。**任何 AI CLI 工具（包括 opencode 自身）都可以将 opencode-orc 作为子进程调用**，无需直接集成即可利用多提供商能力。这实现了任意深度的 Agent 嵌套——一个 AI Agent 可以通过 opencode-orc 派生另一个 AI Agent，形成强大的「套娃」编排模式。

opencode 输出冗长的 JSONL 事件，被其他 AI Agent 解析时会消耗大量 token。opencode-orc 将其压缩为简洁的文本摘要，大幅节省 token 开销：

```
[session] ses_abc123
[tools] 2 calls: read main.go, bash go build
我已审查代码，以下是发现...
[step] end_turn
[done] ok=true session=ses_abc123
```

## 安装

```bash
go install github.com/real-uangi/opencode-orc@latest
```

或从 [Releases](https://github.com/real-uangi/opencode-orc/releases) 下载。

### 让 AI 帮你配置

懒得看文档？把这句话发给你的 AI 助手：

> 「请阅读 https://github.com/real-uangi/opencode-orc/blob/main/AGENTS.md 并按照其中的要求完成配置。」

你的 AI 会自动完成安装、配置和 Skill 创建。

## 使用

```bash
opencode-orc "你的提示词"
```

### 参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-config` | `~/.config/opencode-orc/config.yaml` | 配置文件路径 |
| `-version` | | 显示版本 |
| `-models` | | 列出可用模型 |
| `-model` | | 指定模型（`provider/model` 格式） |

### 示例

```bash
# 简单查询
opencode-orc "main.go 做了什么？"

# 管道输出
opencode-orc "解释这段代码" > explanation.txt

# 使用自定义配置
opencode-orc -config ./my-config.yaml "审查这个 PR"

# 列出可用模型
opencode-orc -models

# 使用指定模型运行
opencode-orc -model deepseek/deepseek-chat "解释量子计算"
```

## 输出格式

### text（默认）

人类可读、省 token 的输出：

```
[session] ses_abc123
[tools] 3 calls: read main.go, bash go test ./..., grep TODO
所有测试通过，代码看起来很干净。
[step] end_turn
[done] ok=true session=ses_abc123
```

### jsonl

机器可解析的 JSON Lines：

```json
{"type":"session","sessionId":"ses_abc123"}
{"type":"tools","count":3,"summary":"read main.go, bash go test ./..., grep TODO"}
{"type":"text","text":"所有测试通过，代码看起来很干净。"}
{"type":"step","reason":"end_turn"}
{"type":"done","sessionId":"ses_abc123","ok":true}
```

## 配置

配置文件位置：`~/.config/opencode-orc/config.yaml`

```yaml
events:
  include:
    - step_start
    - tool_use
    - text
    - step_finish
    - error
  rules:
    step_start:
      keep:
        - sessionID
    tool_use:
      keep:
        - part.tool
        - part.state.status
        - part.state.input
        - part.state.error
        - part.state.metadata.exit
        - part.title
    text:
      keep:
        - part.text
    step_finish:
      keep:
        - part.reason
    error:
      keep:
        - error.name
        - error.data.message
output:
  format: text    # "text" 或 "jsonl"
  pretty: false   # JSON 缩进输出（仅 jsonl 模式）
```

## 事件类型

| 输出类型 | 说明 |
|----------|------|
| `session` | 会话信息（仅输出一次） |
| `text` | LLM 文本输出 |
| `tools` | 工具调用汇总 |
| `step` | 步骤结束（非工具调用） |
| `error` | 错误事件 |
| `done` | 最终状态 |

## 构建

```bash
go build -o opencode-orc .
```

## 许可证

MIT

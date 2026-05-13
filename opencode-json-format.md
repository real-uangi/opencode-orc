

# OpenCode JSON 输出格式说明

> 适用于：
>
> ```bash
> opencode run --format json
> ```
>
> 输出形式为 JSONL（JSON Lines）：
>
> - 每行一个 JSON 对象
> - 非 JSON 数组
> - 按事件流实时输出

---

# 顶层结构

所有事件统一结构：

```json
{
  "type": "事件类型",
  "timestamp": 1759406013703,
  "sessionID": "ses_xxx",
  "...": "事件数据"
}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| type | string | 事件类型 |
| timestamp | number | 毫秒时间戳 |
| sessionID | string | 会话 ID |

---

# 支持的事件类型

| type | 说明 |
|---|---|
| step_start | Agent step 开始 |
| step_finish | Agent step 完成 |
| text | 文本输出 |
| reasoning | 推理输出（需 `--thinking`） |
| tool_use | 工具调用完成/失败 |
| error | 会话错误 |

---

# 1. step_start

```json
{
  "type": "step_start",
  "timestamp": 1759406013703,
  "sessionID": "ses_xxx",
  "part": {
    "id": "prt_xxx",
    "sessionID": "ses_xxx",
    "messageID": "msg_xxx",
    "type": "step-start",
    "snapshot": "optional"
  }
}
```

## 字段说明

| 字段 | 类型 | 说明 |
|---|---|---|
| part.id | string | part ID |
| part.messageID | string | message ID |
| part.snapshot | string? | 快照 ID |

---

# 2. step_finish

```json
{
  "type": "step_finish",
  "timestamp": 1759406019999,
  "sessionID": "ses_xxx",
  "part": {
    "id": "prt_xxx",
    "sessionID": "ses_xxx",
    "messageID": "msg_xxx",
    "type": "step-finish",
    "reason": "stop",
    "snapshot": "optional",
    "cost": 0.0012,
    "tokens": {
      "total": 123,
      "input": 100,
      "output": 20,
      "reasoning": 3,
      "cache": {
        "read": 0,
        "write": 0
      }
    }
  }
}
```

## tokens

| 字段 | 说明 |
|---|---|
| total | 总 token |
| input | 输入 token |
| output | 输出 token |
| reasoning | reasoning token |
| cache.read | cache read token |
| cache.write | cache write token |

---

# 3. text

文本完成后输出。

不是 streaming delta。

```json
{
  "type": "text",
  "timestamp": 1759406015783,
  "sessionID": "ses_xxx",
  "part": {
    "id": "prt_xxx",
    "sessionID": "ses_xxx",
    "messageID": "msg_xxx",
    "type": "text",
    "text": "最终输出文本",
    "synthetic": false,
    "ignored": false,
    "time": {
      "start": 1759406015000,
      "end": 1759406015783
    },
    "metadata": {}
  }
}
```

## 字段说明

| 字段 | 类型 | 说明 |
|---|---|---|
| text | string | 最终文本 |
| synthetic | boolean | 是否 synthetic |
| ignored | boolean | 是否忽略 |
| metadata | object | 附加元数据 |

---

# 4. reasoning

需要：

```bash
--thinking
```

```json
{
  "type": "reasoning",
  "timestamp": 1759406017000,
  "sessionID": "ses_xxx",
  "part": {
    "id": "prt_xxx",
    "sessionID": "ses_xxx",
    "messageID": "msg_xxx",
    "type": "reasoning",
    "text": "推理内容",
    "metadata": {},
    "time": {
      "start": 1759406016000,
      "end": 1759406017000
    }
  }
}
```

---

# 5. tool_use

只在：

- completed
- error

时输出。

不会输出 pending/running。

---

## 5.1 completed

```json
{
  "type": "tool_use",
  "timestamp": 1759406018000,
  "sessionID": "ses_xxx",
  "part": {
    "id": "prt_xxx",
    "sessionID": "ses_xxx",
    "messageID": "msg_xxx",
    "type": "tool",
    "callID": "call_xxx",
    "tool": "read",
    "state": {
      "status": "completed",
      "input": {
        "file": "README.md"
      },
      "output": "文件内容",
      "title": "Read README.md",
      "metadata": {},
      "time": {
        "start": 1759406017000,
        "end": 1759406018000
      },
      "attachments": []
    },
    "metadata": {}
  }
}
```

---

## 5.2 error

```json
{
  "type": "tool_use",
  "timestamp": 1759406018000,
  "sessionID": "ses_xxx",
  "part": {
    "type": "tool",
    "tool": "bash",
    "state": {
      "status": "error",
      "input": {
        "command": "rm -rf /"
      },
      "error": "Permission denied",
      "metadata": {},
      "time": {
        "start": 1759406017000,
        "end": 1759406018000
      }
    }
  }
}
```

---

# 6. error

```json
{
  "type": "error",
  "timestamp": 1759406019000,
  "sessionID": "ses_xxx",
  "error": {
    "name": "APIError",
    "data": {
      "message": "Rate limited",
      "statusCode": 429,
      "isRetryable": true,
      "responseHeaders": {},
      "responseBody": "optional",
      "metadata": {}
    }
  }
}
```

---

# error.name 类型

可能值：

```text
ProviderAuthError
UnknownError
MessageOutputLengthError
MessageAbortedError
StructuredOutputError
ContextOverflowError
APIError
```

---

# 实际输出示例

```json
{"type":"step_start","timestamp":1759406013703,"sessionID":"ses_xxx","part":{"id":"prt_1","type":"step-start"}}
{"type":"text","timestamp":1759406015783,"sessionID":"ses_xxx","part":{"id":"prt_2","type":"text","text":"hello world"}}
{"type":"tool_use","timestamp":1759406018000,"sessionID":"ses_xxx","part":{"type":"tool","tool":"read","state":{"status":"completed"}}}
{"type":"step_finish","timestamp":1759406019999,"sessionID":"ses_xxx","part":{"id":"prt_3","type":"step-finish","reason":"stop"}}
```

---

# 注意事项

## 1. JSONL 不是 JSON

错误：

```json
[
  {...},
  {...}
]
```

正确：

```text
{"type":"text"...}
{"type":"tool_use"...}
{"type":"step_finish"...}
```

---

## 2. stdout 可能混入非 JSON 内容

例如：

- share URL
- warning
- permission 提示
- provider fallback

因此：

- 不要直接整体 parse
- 应逐行解析 JSON

推荐：

```ts
for await (const line of lines) {
  try {
    const event = JSON.parse(line)
  } catch {
    // ignore non-json line
  }
}
```

---

# 推荐事件处理顺序

推荐消费者按：

```text
step_start
  -> reasoning
  -> text
  -> tool_use
step_finish
```

组织状态机。

---

# 推荐 TypeScript 类型

```ts
type OpenCodeEvent =
  | StepStartEvent
  | StepFinishEvent
  | TextEvent
  | ReasoningEvent
  | ToolUseEvent
  | ErrorEvent
```

---
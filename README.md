[中文](./README_zh.md) | English

# opencode-orc

A lightweight orchestrator for [opencode](https://github.com/anomalyco/opencode). It runs `opencode run --format json` as a subprocess, parses the JSONL event stream, filters and transforms events, then outputs a simplified text or JSONL stream.

## Why?

opencode supports most AI providers, making it the ideal bridge layer for delegating sub-agents. Other AI CLI tools can call opencode-orc as a subprocess to leverage multi-provider capabilities without direct integration.

opencode outputs verbose JSONL events, consuming significant tokens when parsed by other AI agents. opencode-orc compresses them into concise text summaries, drastically reducing token usage:

```
[session] ses_abc123
[tools] 2 calls: read main.go, bash go build
I've reviewed the code and here are my findings...
[step] end_turn
[done] ok=true session=ses_abc123
```

## Install

```bash
go install github.com/real-uangi/opencode-orc@latest
```

Or download from [Releases](https://github.com/real-uangi/opencode-orc/releases).

## Usage

```bash
opencode-orc "your prompt here"
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `~/.config/opencode-orc/config.yaml` | Path to config file |
| `-version` | | Show version |

### Examples

```bash
# Simple query
opencode-orc "what does main.go do?"

# Pipe output
opencode-orc "explain this code" > explanation.txt

# Use custom config
opencode-orc -config ./my-config.yaml "review this PR"
```

## Output Formats

### text (default)

Human-readable, token-efficient output:

```
[session] ses_abc123
[tools] 3 calls: read main.go, bash go test ./..., grep TODO
All tests pass. The code looks clean.
[step] end_turn
[done] ok=true session=ses_abc123
```

### jsonl

Machine-parseable JSON Lines:

```json
{"type":"session","sessionId":"ses_abc123"}
{"type":"tools","count":3,"summary":"read main.go, bash go test ./..., grep TODO"}
{"type":"text","text":"All tests pass. The code looks clean."}
{"type":"step","reason":"end_turn"}
{"type":"done","sessionId":"ses_abc123","ok":true}
```

## Configuration

Config file location: `~/.config/opencode-orc/config.yaml`

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
  format: text    # "text" or "jsonl"
  pretty: false   # indent JSON output (jsonl mode only)
```

## Event Types

| Output Type | Description |
|-------------|-------------|
| `session` | Session info (emitted once) |
| `text` | LLM text output |
| `tools` | Aggregated tool call summary |
| `step` | Step finish (non-tool) |
| `error` | Error events |
| `done` | Final status |

## Build

```bash
go build -o opencode-orc .
```

## License

MIT

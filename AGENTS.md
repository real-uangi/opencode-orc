# AI Configuration Guide for opencode-orc

> This document is designed for AI agents. Human users: give this file or its URL to your AI assistant and say "please set up opencode-orc for me according to this guide."

## What is opencode-orc?

opencode-orc is a lightweight orchestrator for [opencode](https://github.com/anomalyco/opencode). It runs `opencode run --format json` as a subprocess, parses the JSONL event stream, filters and transforms events, then outputs a simplified text or JSONL stream.

**Key capability**: Any AI CLI tool (including opencode itself) can call opencode-orc as a subprocess to delegate tasks. This enables "nesting" or "meta-orchestration" — an AI agent can spawn another AI agent via opencode-orc.

## Installation

### Option 1: Go install

```bash
go install github.com/real-uangi/opencode-orc@latest
```

### Option 2: Download binary

Download from [Releases](https://github.com/real-uangi/opencode-orc/releases) and place in PATH.

### Option 3: Build from source

```bash
git clone https://github.com/real-uangi/opencode-orc.git
cd opencode-orc
go build -o opencode-orc .
```

## Verifying Installation

```bash
opencode-orc -version
```

## Configuration

### Config File Location

Default: `~/.config/opencode-orc/config.yaml`

You can specify a custom path with `-config` flag.

### Default Config

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

### Config Fields

- `events.include`: Which event types to include in output
- `events.rules.<type>.keep`: Which fields to keep for each event type (dot notation)
- `output.format`: Output format — `"text"` (human-readable) or `"jsonl"` (machine-parseable)
- `output.pretty`: Whether to indent JSON in jsonl mode

## Usage Examples

### Basic usage

```bash
opencode-orc "explain what this project does"
```

### With custom config

```bash
opencode-orc -config ./my-config.yaml "review this code"
```

### List available models

```bash
opencode-orc -models
```

### Run with a specific model

```bash
opencode-orc -model deepseek/deepseek-chat "explain quantum computing"
```

### Pipe output

```bash
opencode-orc "find bugs" > bugs.txt
```

### JSONL output (for programmatic use)

Set in config:
```yaml
output:
  format: jsonl
```

## Output Format Reference

### text (default)

```
[session] ses_abc123
[tools] 3 calls: read main.go, bash go test ./..., grep TODO
All tests pass. The code looks clean.
[step] end_turn
[done] ok=true session=ses_abc123
```

### jsonl

```json
{"type":"session","sessionId":"ses_abc123"}
{"type":"tools","count":3,"summary":"read main.go, bash go test ./..., grep TODO"}
{"type":"text","text":"All tests pass. The code looks clean."}
{"type":"step","reason":"end_turn"}
{"type":"done","sessionId":"ses_abc123","ok":true}
```

## Event Types

| Type | Description |
|------|-------------|
| `session` | Session info (emitted once) |
| `text` | LLM text output |
| `tools` | Aggregated tool call summary |
| `step` | Step finish reason |
| `error` | Error events |
| `done` | Final status with session ID |

## Creating a Skill for This Project

If your AI framework supports skills (e.g., OpenCode, Claude Code, etc.), you can create a skill that:

1. **Wraps opencode-orc as a tool/function**:
   - Name: `opencode_orc` or `delegate_to_opencode`
   - Description: "Delegate a task to opencode via opencode-orc orchestrator"
   - Parameter: `prompt` (string) — the task to delegate
   - Optional parameter: `config_path` (string) — path to custom config
   - Optional parameter: `model` (string) — model to use in `provider/model` format
   - Implementation: Execute `opencode-orc "<prompt>"` and return stdout

2. **Example skill definition** (pseudocode):
   ```yaml
   name: opencode-orc-delegate
   description: Delegate tasks to opencode AI via opencode-orc
   parameters:
     prompt:
       type: string
       description: The task description to send to opencode
     config_path:
       type: string
       description: Optional path to opencode-orc config file
       default: ""
     model:
       type: string
       description: Optional model to use (provider/model format)
       default: ""
   command: |
     args=""
     if [ -n "$config_path" ]; then
       args="$args -config $config_path"
     fi
     if [ -n "$model" ]; then
       args="$args -model $model"
     fi
     opencode-orc $args "$prompt"
   ```

3. **Nested/Meta usage**:
   Since opencode-orc can be called by any AI CLI (including opencode itself), you can create recursive delegation patterns:
   - Parent AI (any CLI) → calls opencode-orc → spawns opencode sub-agent
   - The sub-agent can itself call opencode-orc to spawn another agent
   - This creates arbitrary-depth agent nesting ("套娃")

## Prerequisites Check

Before using opencode-orc, verify:

1. opencode is installed and configured (`opencode --version`)
2. opencode-orc is in PATH (`opencode-orc -version`)
3. Config file exists at `~/.config/opencode-orc/config.yaml` (or specify with `-config`)

## Troubleshooting

- **"opencode-orc: command not found"**: Ensure the binary is in your PATH
- **Empty output**: Check that opencode is properly configured with API keys
- **Too verbose**: Switch `output.format` to `"text"` (default)
- **Need structured data**: Switch `output.format` to `"jsonl"`

## License

MIT

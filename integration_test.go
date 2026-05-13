package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"opencode-orc/buffer"
	"opencode-orc/config"
	"opencode-orc/filter"
	"opencode-orc/output"
	"opencode-orc/parser"
	"opencode-orc/types"
)

func TestIntegration_FilterEvents(t *testing.T) {
	// Sample JSONL input
	input := `{"type":"step_start","timestamp":1759406013703,"sessionID":"ses_xxx","part":{"id":"prt_1","type":"step-start"}}
{"type":"text","timestamp":1759406015783,"sessionID":"ses_xxx","part":{"id":"prt_2","type":"text","text":"hello world"}}
{"type":"tool_use","timestamp":1759406018000,"sessionID":"ses_xxx","part":{"type":"tool","tool":"read","state":{"status":"completed","input":{"file":"README.md"}}}}
{"type":"step_finish","timestamp":1759406019999,"sessionID":"ses_xxx","part":{"id":"prt_3","type":"step-finish","reason":"stop"}}
`

	// Setup
	cfg := config.DefaultConfig()
	p := parser.NewParser(strings.NewReader(input))
	f := filter.New(cfg)
	b := buffer.New()

	var outputBuf bytes.Buffer
	w := output.New(&outputBuf, false)

	// Process events
	for {
		event, err := p.ParseNext()
		if err != nil {
			break
		}

		filtered := f.Filter(event)
		if filtered == nil {
			continue
		}

		b.Update(filtered)
		w.WriteEvent(filtered)
	}

	// Write done event
	doneEvent := b.DoneEvent()
	w.WriteEvent(doneEvent)

	// Verify output
	lines := strings.Split(strings.TrimSpace(outputBuf.String()), "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 output lines, got %d", len(lines))
	}

	// Verify session event
	var sessionEvent types.SessionEvent
	json.Unmarshal([]byte(lines[0]), &sessionEvent)
	if sessionEvent.Type != "session" {
		t.Errorf("expected type session, got %s", sessionEvent.Type)
	}
	if sessionEvent.SessionID != "ses_xxx" {
		t.Errorf("expected sessionId ses_xxx, got %s", sessionEvent.SessionID)
	}

	// Verify text event
	var textEvent types.TextEvent
	json.Unmarshal([]byte(lines[1]), &textEvent)
	if textEvent.Type != "text" {
		t.Errorf("expected type text, got %s", textEvent.Type)
	}
	if textEvent.Text != "hello world" {
		t.Errorf("expected text 'hello world', got %s", textEvent.Text)
	}

	// Verify tool event
	var toolEvent types.ToolEvent
	json.Unmarshal([]byte(lines[2]), &toolEvent)
	if toolEvent.Type != "tool" {
		t.Errorf("expected type tool, got %s", toolEvent.Type)
	}
	if toolEvent.Tool != "read" {
		t.Errorf("expected tool read, got %s", toolEvent.Tool)
	}

	// Verify step event
	var stepEvent types.StepEvent
	json.Unmarshal([]byte(lines[3]), &stepEvent)
	if stepEvent.Type != "step" {
		t.Errorf("expected type step, got %s", stepEvent.Type)
	}
	if stepEvent.Reason != "stop" {
		t.Errorf("expected reason stop, got %s", stepEvent.Reason)
	}

	// Verify done event
	var doneEventResult types.DoneEvent
	json.Unmarshal([]byte(lines[4]), &doneEventResult)
	if doneEventResult.Type != "done" {
		t.Errorf("expected type done, got %s", doneEventResult.Type)
	}
	if !doneEventResult.Ok {
		t.Errorf("expected ok true, got false")
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/real-uangi/opencode-orc/buffer"
	"github.com/real-uangi/opencode-orc/config"
	"github.com/real-uangi/opencode-orc/filter"
	"github.com/real-uangi/opencode-orc/output"
	"github.com/real-uangi/opencode-orc/parser"
	"github.com/real-uangi/opencode-orc/types"
)

func TestIntegration_FilterEvents(t *testing.T) {
	// Sample JSONL input with tool calls
	input := `{"type":"step_start","timestamp":1759406013703,"sessionID":"ses_xxx","part":{"id":"prt_1","type":"step-start"}}
{"type":"tool_use","timestamp":1759406018000,"sessionID":"ses_xxx","part":{"type":"tool","tool":"read","state":{"status":"completed","input":{"file":"README.md"}}}}
{"type":"tool_use","timestamp":1759406018500,"sessionID":"ses_xxx","part":{"type":"tool","tool":"read","state":{"status":"completed","input":{"file":"go.mod"}}}}
{"type":"step_finish","timestamp":1759406019000,"sessionID":"ses_xxx","part":{"id":"prt_2","type":"step-finish","reason":"tool-calls"}}
{"type":"step_start","timestamp":1759406020000,"sessionID":"ses_xxx","part":{"id":"prt_3","type":"step-start"}}
{"type":"text","timestamp":1759406021000,"sessionID":"ses_xxx","part":{"id":"prt_4","type":"text","text":"hello world"}}
{"type":"step_finish","timestamp":1759406022000,"sessionID":"ses_xxx","part":{"id":"prt_5","type":"step-finish","reason":"stop"}}
`

	// Setup
	cfg := config.DefaultConfig()
	p := parser.NewParser(strings.NewReader(input))
	f := filter.New(cfg)
	b := buffer.New()

	var outputBuf bytes.Buffer
	w := output.New(&outputBuf, "jsonl", false)

	// Process events
	for {
		event, err := p.ParseNext()
		if err != nil {
			break
		}

		events := f.Filter(event)
		for _, filtered := range events {
			b.Update(filtered)
			w.WriteEvent(filtered)
		}
	}

	// Write done event
	doneEvent := b.DoneEvent()
	w.WriteEvent(doneEvent)

	// Verify output
	lines := strings.Split(strings.TrimSpace(outputBuf.String()), "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 output lines, got %d: %v", len(lines), lines)
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

	// Verify tools summary event
	var toolsEvent types.ToolsEvent
	json.Unmarshal([]byte(lines[1]), &toolsEvent)
	if toolsEvent.Type != "tools" {
		t.Errorf("expected type tools, got %s", toolsEvent.Type)
	}
	if toolsEvent.Count != 2 {
		t.Errorf("expected count 2, got %d", toolsEvent.Count)
	}
	if toolsEvent.Summary != "read README.md, read go.mod" {
		t.Errorf("expected summary 'read README.md, read go.mod', got %s", toolsEvent.Summary)
	}

	// Verify text event
	var textEvent types.TextEvent
	json.Unmarshal([]byte(lines[2]), &textEvent)
	if textEvent.Type != "text" {
		t.Errorf("expected type text, got %s", textEvent.Type)
	}
	if textEvent.Text != "hello world" {
		t.Errorf("expected text 'hello world', got %s", textEvent.Text)
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

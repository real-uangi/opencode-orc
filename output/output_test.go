package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/real-uangi/opencode-orc/types"
)

func TestWriteEvent_CompactJSON(t *testing.T) {
	var buf bytes.Buffer
	writer := New(&buf, "jsonl", false)

	event := types.SessionEvent{
		Type:      types.OutputTypeSession,
		SessionID: "ses_123",
	}

	err := writer.WriteEvent(event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := strings.TrimSpace(buf.String())

	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(output), &parsed)
	if err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if parsed["type"] != types.OutputTypeSession {
		t.Errorf("expected type '%s', got %v", types.OutputTypeSession, parsed["type"])
	}

	if parsed["sessionId"] != "ses_123" {
		t.Errorf("expected sessionId 'ses_123', got %v", parsed["sessionId"])
	}
}

func TestWriteEvent_PrettyJSON(t *testing.T) {
	var buf bytes.Buffer
	writer := New(&buf, "jsonl", true)

	event := types.DoneEvent{
		Type:      types.OutputTypeDone,
		SessionID: "ses_456",
		Ok:        true,
	}

	err := writer.WriteEvent(event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "\n") {
		t.Error("expected pretty output with newlines")
	}

	if !strings.Contains(output, "  ") {
		t.Error("expected indentation in pretty output")
	}

	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(output), &parsed)
	if err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if parsed["type"] != types.OutputTypeDone {
		t.Errorf("expected type '%s', got %v", types.OutputTypeDone, parsed["type"])
	}
}

func TestWriteEvent_MultipleEvents(t *testing.T) {
	var buf bytes.Buffer
	writer := New(&buf, "jsonl", false)

	events := []interface{}{
		types.TextEvent{Type: types.OutputTypeText, Text: "hello"},
		types.TextEvent{Type: types.OutputTypeText, Text: "world"},
	}

	for _, event := range events {
		err := writer.WriteEvent(event)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	for i, line := range lines {
		var parsed map[string]interface{}
		err := json.Unmarshal([]byte(line), &parsed)
		if err != nil {
			t.Fatalf("invalid JSON on line %d: %v", i, err)
		}
	}
}

func TestWriteEvent_ToolsEvent(t *testing.T) {
	var buf bytes.Buffer
	writer := New(&buf, "jsonl", false)

	event := types.ToolsEvent{
		Type:    types.OutputTypeTools,
		Count:   3,
		Summary: "read README.md, ls -la, grep pattern",
	}

	err := writer.WriteEvent(event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := strings.TrimSpace(buf.String())

	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(output), &parsed)
	if err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if parsed["type"] != types.OutputTypeTools {
		t.Errorf("expected type '%s', got %v", types.OutputTypeTools, parsed["type"])
	}

	if parsed["count"] != float64(3) {
		t.Errorf("expected count 3, got %v", parsed["count"])
	}

	if parsed["summary"] != "read README.md, ls -la, grep pattern" {
		t.Errorf("expected summary 'read README.md, ls -la, grep pattern', got %v", parsed["summary"])
	}
}

func TestWriteEvent_StepEvent(t *testing.T) {
	var buf bytes.Buffer
	writer := New(&buf, "jsonl", false)

	event := types.StepEvent{
		Type:   types.OutputTypeStep,
		Reason: "end_turn",
	}

	err := writer.WriteEvent(event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := strings.TrimSpace(buf.String())

	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(output), &parsed)
	if err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if parsed["reason"] != "end_turn" {
		t.Errorf("expected reason 'end_turn', got %v", parsed["reason"])
	}
}

func TestWriteEvent_NilEvent(t *testing.T) {
	var buf bytes.Buffer
	writer := New(&buf, "jsonl", false)

	err := writer.WriteEvent(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	if output != "null" {
		t.Errorf("expected 'null', got '%s'", output)
	}
}

func TestWriteEvent_EmptyStruct(t *testing.T) {
	var buf bytes.Buffer
	writer := New(&buf, "jsonl", false)

	type Empty struct{}
	err := writer.WriteEvent(Empty{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	if output != "{}" {
		t.Errorf("expected '{}', got '%s'", output)
	}
}

func TestWriteEvent_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	writer := New(&buf, "text", false)

	events := []interface{}{
		&types.SessionEvent{Type: "session", SessionID: "ses_abc"},
		&types.TextEvent{Type: "text", Text: "hello world"},
		&types.ToolsEvent{Type: "tools", Count: 2, Summary: "read main.go, bash go build"},
		&types.StepEvent{Type: "step", Reason: "end_turn"},
		&types.DoneEvent{Type: "done", SessionID: "ses_abc", Ok: true},
	}

	for _, event := range events {
		err := writer.WriteEvent(event)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines, got %d: %v", len(lines), lines)
	}

	expected := []string{
		"[session] ses_abc",
		"hello world",
		"[tools] 2 calls: read main.go, bash go build",
		"[step] end_turn",
		"[done] ok=true session=ses_abc",
	}

	for i, exp := range expected {
		if lines[i] != exp {
			t.Errorf("line %d: expected '%s', got '%s'", i, exp, lines[i])
		}
	}
}

func TestWriteEvent_TextFormat_Error(t *testing.T) {
	var buf bytes.Buffer
	writer := New(&buf, "text", false)

	events := []interface{}{
		&types.ToolEvent{Type: "error", Tool: "error", Status: "error", Action: "ProcessError", Error: "exit status 1"},
		&types.DoneEvent{Type: "done", SessionID: "ses_abc", Ok: false, Error: "exit status 1"},
	}

	for _, event := range events {
		err := writer.WriteEvent(event)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(lines), lines)
	}

	expected := []string{
		"[error] ProcessError: exit status 1",
		"[done] ok=false error=exit status 1 session=ses_abc",
	}

	for i, exp := range expected {
		if lines[i] != exp {
			t.Errorf("line %d: expected '%s', got '%s'", i, exp, lines[i])
		}
	}
}

package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"opencode-orc/types"
)

func TestWriteEvent_CompactJSON(t *testing.T) {
	var buf bytes.Buffer
	writer := New(&buf, false)

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
	writer := New(&buf, true)

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
	writer := New(&buf, false)

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

func TestWriteEvent_ToolEvent(t *testing.T) {
	var buf bytes.Buffer
	writer := New(&buf, false)

	exitCode := 0
	event := types.ToolEvent{
		Type:   types.OutputTypeTool,
		Tool:   "bash",
		Status: "completed",
		Action: "ls -la",
		Exit:   &exitCode,
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

	if parsed["tool"] != "bash" {
		t.Errorf("expected tool 'bash', got %v", parsed["tool"])
	}

	if parsed["action"] != "ls -la" {
		t.Errorf("expected action 'ls -la', got %v", parsed["action"])
	}
}

func TestWriteEvent_StepEvent(t *testing.T) {
	var buf bytes.Buffer
	writer := New(&buf, false)

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
	writer := New(&buf, false)

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
	writer := New(&buf, false)

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

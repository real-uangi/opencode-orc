package filter

import (
	"testing"

	"github.com/real-uangi/opencode-orc/types"
)

func TestFilterStepStart(t *testing.T) {
	cfg := &types.Config{
		Events: types.EventsConfig{
			Include: []string{types.EventTypeStepStart},
			Rules: map[string]types.EventRule{
				types.EventTypeStepStart: {
					Keep: []string{"sessionID"},
				},
			},
		},
	}

	filter := New(cfg)

	event := &types.RawEvent{
		Type:      types.EventTypeStepStart,
		SessionID: "ses_123",
		Part: map[string]interface{}{
			"id":        "prt_1",
			"messageID": "msg_1",
		},
	}

	events := filter.Filter(event)
	if len(events) == 0 {
		t.Fatal("expected events, got empty slice")
	}

	sessionEvent, ok := events[0].(*types.SessionEvent)
	if !ok {
		t.Fatalf("expected SessionEvent, got %T", events[0])
	}

	if sessionEvent.SessionID != "ses_123" {
		t.Errorf("expected sessionID ses_123, got %s", sessionEvent.SessionID)
	}
}

func TestFilterExcludedEventType(t *testing.T) {
	cfg := &types.Config{
		Events: types.EventsConfig{
			Include: []string{types.EventTypeStepStart},
		},
	}

	filter := New(cfg)

	event := &types.RawEvent{
		Type:      types.EventTypeText,
		SessionID: "ses_123",
		Part: map[string]interface{}{
			"text": "hello",
		},
	}

	events := filter.Filter(event)
	if len(events) != 0 {
		t.Errorf("expected empty slice for excluded event type, got %v", events)
	}
}

func TestFilterTextEvent(t *testing.T) {
	cfg := &types.Config{
		Events: types.EventsConfig{
			Include: []string{types.EventTypeText},
		},
	}

	filter := New(cfg)

	event := &types.RawEvent{
		Type:      types.EventTypeText,
		SessionID: "ses_123",
		Part: map[string]interface{}{
			"text": "hello world",
		},
	}

	events := filter.Filter(event)
	if len(events) == 0 {
		t.Fatal("expected events, got empty slice")
	}

	textEvent, ok := events[0].(*types.TextEvent)
	if !ok {
		t.Fatalf("expected TextEvent, got %T", events[0])
	}

	if textEvent.Text != "hello world" {
		t.Errorf("expected text 'hello world', got %s", textEvent.Text)
	}
}

func TestFilterStepFinishEvent(t *testing.T) {
	cfg := &types.Config{
		Events: types.EventsConfig{
			Include: []string{types.EventTypeStepFinish},
		},
	}

	filter := New(cfg)

	event := &types.RawEvent{
		Type:      types.EventTypeStepFinish,
		SessionID: "ses_123",
		Part: map[string]interface{}{
			"reason": "end_turn",
		},
	}

	events := filter.Filter(event)
	if len(events) == 0 {
		t.Fatal("expected events, got empty slice")
	}

	stepEvent, ok := events[0].(*types.StepEvent)
	if !ok {
		t.Fatalf("expected StepEvent, got %T", events[0])
	}

	if stepEvent.Reason != "end_turn" {
		t.Errorf("expected reason 'end_turn', got %s", stepEvent.Reason)
	}
}

func TestFilterToolUseBuffersAndSummarizes(t *testing.T) {
	cfg := &types.Config{
		Events: types.EventsConfig{
			Include: []string{types.EventTypeToolUse, types.EventTypeStepFinish},
		},
	}

	filter := New(cfg)

	// First tool call
	event1 := &types.RawEvent{
		Type:      types.EventTypeToolUse,
		SessionID: "ses_123",
		Part: map[string]interface{}{
			"part": map[string]interface{}{
				"tool": "bash",
				"state": map[string]interface{}{
					"status": "completed",
					"input": map[string]interface{}{
						"command": "ls -la",
					},
				},
			},
		},
	}

	events := filter.Filter(event1)
	if len(events) != 0 {
		t.Errorf("expected no output for tool_use, got %v", events)
	}

	// Second tool call
	event2 := &types.RawEvent{
		Type:      types.EventTypeToolUse,
		SessionID: "ses_123",
		Part: map[string]interface{}{
			"part": map[string]interface{}{
				"tool": "read",
				"state": map[string]interface{}{
					"status": "completed",
					"input": map[string]interface{}{
						"file": "README.md",
					},
				},
			},
		},
	}

	events = filter.Filter(event2)
	if len(events) != 0 {
		t.Errorf("expected no output for tool_use, got %v", events)
	}

	// Step finish should output tools summary
	event3 := &types.RawEvent{
		Type:      types.EventTypeStepFinish,
		SessionID: "ses_123",
		Part: map[string]interface{}{
			"reason": "tool-calls",
		},
	}

	events = filter.Filter(event3)
	if len(events) == 0 {
		t.Fatal("expected events, got empty slice")
	}

	toolsEvent, ok := events[0].(*types.ToolsEvent)
	if !ok {
		t.Fatalf("expected ToolsEvent, got %T", events[0])
	}

	if toolsEvent.Count != 2 {
		t.Errorf("expected count 2, got %d", toolsEvent.Count)
	}

	if toolsEvent.Summary != "ls -la, read README.md" {
		t.Errorf("expected summary 'ls -la, read README.md', got %s", toolsEvent.Summary)
	}
}

func TestFilterErrorEvent(t *testing.T) {
	cfg := &types.Config{
		Events: types.EventsConfig{
			Include: []string{types.EventTypeError},
		},
	}

	filter := New(cfg)

	event := &types.RawEvent{
		Type:      types.EventTypeError,
		SessionID: "ses_123",
		Error: map[string]interface{}{
			"error": map[string]interface{}{
				"name": "AuthError",
				"data": map[string]interface{}{
					"message": "unauthorized",
				},
			},
		},
	}

	events := filter.Filter(event)
	if len(events) == 0 {
		t.Fatal("expected events, got empty slice")
	}

	toolEvent, ok := events[0].(*types.ToolEvent)
	if !ok {
		t.Fatalf("expected ToolEvent, got %T", events[0])
	}

	if toolEvent.Tool != "error" {
		t.Errorf("expected tool 'error', got %s", toolEvent.Tool)
	}

	if toolEvent.Action != "AuthError" {
		t.Errorf("expected action 'AuthError', got %s", toolEvent.Action)
	}

	if toolEvent.Error != "unauthorized" {
		t.Errorf("expected error 'unauthorized', got %s", toolEvent.Error)
	}
}

package filter

import (
	"testing"

	"opencode-orc/types"
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
			"id":      "prt_1",
			"messageID": "msg_1",
		},
	}

	result := filter.Filter(event)
	if result == nil {
		t.Fatal("expected result, got nil")
	}

	sessionEvent, ok := result.(*types.SessionEvent)
	if !ok {
		t.Fatalf("expected SessionEvent, got %T", result)
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

	result := filter.Filter(event)
	if result != nil {
		t.Errorf("expected nil for excluded event type, got %v", result)
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

	result := filter.Filter(event)
	if result == nil {
		t.Fatal("expected result, got nil")
	}

	textEvent, ok := result.(*types.TextEvent)
	if !ok {
		t.Fatalf("expected TextEvent, got %T", result)
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

	result := filter.Filter(event)
	if result == nil {
		t.Fatal("expected result, got nil")
	}

	stepEvent, ok := result.(*types.StepEvent)
	if !ok {
		t.Fatalf("expected StepEvent, got %T", result)
	}

	if stepEvent.Reason != "end_turn" {
		t.Errorf("expected reason 'end_turn', got %s", stepEvent.Reason)
	}
}

func TestFilterToolUseBash(t *testing.T) {
	cfg := &types.Config{
		Events: types.EventsConfig{
			Include: []string{types.EventTypeToolUse},
		},
	}

	filter := New(cfg)

	event := &types.RawEvent{
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
					"metadata": map[string]interface{}{
						"exit": 0,
					},
				},
			},
		},
	}

	result := filter.Filter(event)
	if result == nil {
		t.Fatal("expected result, got nil")
	}

	toolEvent, ok := result.(*types.ToolEvent)
	if !ok {
		t.Fatalf("expected ToolEvent, got %T", result)
	}

	if toolEvent.Tool != "bash" {
		t.Errorf("expected tool 'bash', got %s", toolEvent.Tool)
	}

	if toolEvent.Action != "ls -la" {
		t.Errorf("expected action 'ls -la', got %s", toolEvent.Action)
	}

	if toolEvent.Exit == nil || *toolEvent.Exit != 0 {
		t.Errorf("expected exit code 0, got %v", toolEvent.Exit)
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

	result := filter.Filter(event)
	if result == nil {
		t.Fatal("expected result, got nil")
	}

	toolEvent, ok := result.(*types.ToolEvent)
	if !ok {
		t.Fatalf("expected ToolEvent, got %T", result)
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

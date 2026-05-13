package filter

import (
	"fmt"
	"strings"

	"opencode-orc/types"
)

// Filter processes events according to configuration
type Filter struct {
	config    *types.Config
	toolCalls []string // Buffer for tool call summaries in current step
	sessionID string   // Current session ID
	emitted   bool     // Whether session event has been emitted
}

// New creates a new event filter
func New(config *types.Config) *Filter {
	return &Filter{
		config:    config,
		toolCalls: make([]string, 0),
	}
}

// Filter processes a raw event and returns filtered output events
// Returns empty slice if the event should be skipped
func (f *Filter) Filter(event *types.RawEvent) []interface{} {
	// Check if event type is included
	if !f.isIncluded(event.Type) {
		return nil
	}

	switch event.Type {
	case types.EventTypeStepStart:
		return f.filterStepStart(event)
	case types.EventTypeToolUse:
		return f.filterToolUse(event)
	case types.EventTypeText:
		return f.filterText(event)
	case types.EventTypeStepFinish:
		return f.filterStepFinish(event)
	case types.EventTypeError:
		return f.filterError(event)
	default:
		return nil
	}
}

func (f *Filter) isIncluded(eventType string) bool {
	for _, t := range f.config.Events.Include {
		if t == eventType {
			return true
		}
	}
	return false
}

func (f *Filter) filterStepStart(event *types.RawEvent) []interface{} {
	f.sessionID = event.SessionID
	
	// Only emit session event once
	if !f.emitted {
		f.emitted = true
		return []interface{}{
			&types.SessionEvent{
				Type:      types.OutputTypeSession,
				SessionID: event.SessionID,
			},
		}
	}
	return nil
}

func (f *Filter) filterToolUse(event *types.RawEvent) []interface{} {
	part, ok := event.Part["part"].(map[string]interface{})
	if !ok {
		part = event.Part
	}

	tool := getStringFromMap(part, "tool")
	state, _ := part["state"].(map[string]interface{})

	// Build action summary and add to buffer
	action := f.buildToolAction(tool, state)
	f.toolCalls = append(f.toolCalls, action)

	// Don't output individual tool events
	return nil
}

func (f *Filter) buildToolAction(tool string, state map[string]interface{}) string {
	input, ok := state["input"].(map[string]interface{})
	if !ok {
		return tool
	}

	switch tool {
	case "read":
		if file, ok := input["file"].(string); ok {
			return fmt.Sprintf("read %s", file)
		}
	case "bash":
		if cmd, ok := input["command"].(string); ok {
			return cmd
		}
	case "write":
		if file, ok := input["file"].(string); ok {
			return fmt.Sprintf("write %s", file)
		}
	case "grep", "glob":
		if pattern, ok := input["pattern"].(string); ok {
			return fmt.Sprintf("%s %s", tool, pattern)
		}
	}

	return tool
}

func (f *Filter) filterText(event *types.RawEvent) []interface{} {
	part, ok := event.Part["part"].(map[string]interface{})
	if !ok {
		part = event.Part
	}

	text := getStringFromMap(part, "text")
	return []interface{}{
		&types.TextEvent{
			Type: types.OutputTypeText,
			Text: text,
		},
	}
}

func (f *Filter) filterStepFinish(event *types.RawEvent) []interface{} {
	part, ok := event.Part["part"].(map[string]interface{})
	if !ok {
		part = event.Part
	}

	reason := getStringFromMap(part, "reason")

	// If there are buffered tool calls, output a summary
	if len(f.toolCalls) > 0 {
		summary := strings.Join(f.toolCalls, ", ")
		toolsEvent := &types.ToolsEvent{
			Type:    types.OutputTypeTools,
			Count:   len(f.toolCalls),
			Summary: summary,
		}
		f.toolCalls = make([]string, 0) // Clear buffer
		return []interface{}{toolsEvent}
	}

	// Only output step event if it's not a tool-calls step
	if reason != "tool-calls" {
		return []interface{}{
			&types.StepEvent{
				Type:   types.OutputTypeStep,
				Reason: reason,
			},
		}
	}

	return nil
}

func (f *Filter) filterError(event *types.RawEvent) []interface{} {
	errData, ok := event.Error["error"].(map[string]interface{})
	if !ok {
		errData = event.Error
	}

	name := getStringFromMap(errData, "name")
	data, _ := errData["data"].(map[string]interface{})
	message := getStringFromMap(data, "message")

	return []interface{}{
		&types.ToolEvent{
			Type:   "error",
			Tool:   "error",
			Status: "error",
			Action: name,
			Error:  message,
		},
	}
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", val)
	}
	return ""
}

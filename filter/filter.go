package filter

import (
	"fmt"

	"opencode-orc/types"
)

// Filter processes events according to configuration
type Filter struct {
	config *types.Config
}

// New creates a new event filter
func New(config *types.Config) *Filter {
	return &Filter{
		config: config,
	}
}

// Filter processes a raw event and returns a filtered output event
// Returns nil if the event should be skipped
func (f *Filter) Filter(event *types.RawEvent) interface{} {
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

func (f *Filter) filterStepStart(event *types.RawEvent) *types.SessionEvent {
	return &types.SessionEvent{
		Type:      types.OutputTypeSession,
		SessionID: event.SessionID,
	}
}

func (f *Filter) filterToolUse(event *types.RawEvent) *types.ToolEvent {
	part, ok := event.Part["part"].(map[string]interface{})
	if !ok {
		part = event.Part
	}

	tool := getStringFromMap(part, "tool")
	state, _ := part["state"].(map[string]interface{})
	status := getStringFromMap(state, "status")

	// Build action summary
	action := f.buildToolAction(tool, state)

	// Get exit code for bash commands
	var exit *int
	if tool == "bash" {
		if metadata, ok := state["metadata"].(map[string]interface{}); ok {
			if exitVal, ok := metadata["exit"]; ok {
				if exitInt, ok := exitVal.(int); ok {
					exit = &exitInt
				} else if exitFloat, ok := exitVal.(float64); ok {
					exitInt := int(exitFloat)
					exit = &exitInt
				}
			}
		}
	}

	// Get error
	errMsg := ""
	if status == "error" {
		errMsg = getStringFromMap(state, "error")
	}

	return &types.ToolEvent{
		Type:   types.OutputTypeTool,
		Tool:   tool,
		Status: status,
		Action: action,
		Exit:   exit,
		Error:  errMsg,
	}
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

func (f *Filter) filterText(event *types.RawEvent) *types.TextEvent {
	part, ok := event.Part["part"].(map[string]interface{})
	if !ok {
		part = event.Part
	}

	text := getStringFromMap(part, "text")
	return &types.TextEvent{
		Type: types.OutputTypeText,
		Text: text,
	}
}

func (f *Filter) filterStepFinish(event *types.RawEvent) *types.StepEvent {
	part, ok := event.Part["part"].(map[string]interface{})
	if !ok {
		part = event.Part
	}

	reason := getStringFromMap(part, "reason")
	return &types.StepEvent{
		Type:   types.OutputTypeStep,
		Reason: reason,
	}
}

func (f *Filter) filterError(event *types.RawEvent) *types.ToolEvent {
	errData, ok := event.Error["error"].(map[string]interface{})
	if !ok {
		errData = event.Error
	}

	name := getStringFromMap(errData, "name")
	data, _ := errData["data"].(map[string]interface{})
	message := getStringFromMap(data, "message")

	return &types.ToolEvent{
		Type:   types.OutputTypeTool,
		Tool:   "error",
		Status: "error",
		Action: name,
		Error:  message,
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

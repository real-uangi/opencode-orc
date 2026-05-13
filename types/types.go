package types

import "fmt"

// Event types from opencode JSONL output
const (
	EventTypeStepStart  = "step_start"
	EventTypeStepFinish = "step_finish"
	EventTypeText       = "text"
	EventTypeToolUse    = "tool_use"
	EventTypeError      = "error"
	EventTypeReasoning  = "reasoning"
)

// Output event types
const (
	OutputTypeSession = "session"
	OutputTypeTools   = "tools"
	OutputTypeText    = "text"
	OutputTypeStep    = "step"
	OutputTypeDone    = "done"
)

// RawEvent represents a parsed JSONL event from opencode
type RawEvent struct {
	Type      string                 `json:"type"`
	Timestamp int64                  `json:"timestamp"`
	SessionID string                 `json:"sessionID"`
	Part      map[string]interface{} `json:"part,omitempty"`
	Error     map[string]interface{} `json:"error,omitempty"`
}

// Formatter can format itself as text
type Formatter interface {
	FormatText() string
}

// SessionEvent represents a session info event
type SessionEvent struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionId"`
}

func (e *SessionEvent) FormatText() string {
	return "[session] " + e.SessionID
}

// ToolEvent represents a tool use event
type ToolEvent struct {
	Type   string `json:"type"`
	Tool   string `json:"tool"`
	Status string `json:"status"`
	Action string `json:"action"`
	Exit   *int   `json:"exit,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (e *ToolEvent) FormatText() string {
	if e.Error != "" {
		return "[error] " + e.Action + ": " + e.Error
	}
	return "[" + e.Status + "] " + e.Action
}

// TextEvent represents a text output event
type TextEvent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (e *TextEvent) FormatText() string {
	return e.Text
}

// StepEvent represents a step finish event
type StepEvent struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

func (e *StepEvent) FormatText() string {
	return "[step] " + e.Reason
}

// ToolsEvent represents a summary of tool calls in a step
type ToolsEvent struct {
	Type    string `json:"type"`
	Count   int    `json:"count"`
	Summary string `json:"summary"`
}

func (e *ToolsEvent) FormatText() string {
	return fmt.Sprintf("[tools] %d calls: %s", e.Count, e.Summary)
}

// DoneEvent represents the final done event
type DoneEvent struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionId"`
	Ok        bool   `json:"ok"`
	Error     string `json:"error,omitempty"`
}

func (e *DoneEvent) FormatText() string {
	status := "ok=true"
	if !e.Ok {
		status = "ok=false"
	}
	if e.Error != "" {
		status += " error=" + e.Error
	}
	return fmt.Sprintf("[done] %s session=%s", status, e.SessionID)
}

// Config represents the application configuration
type Config struct {
	Events EventsConfig `yaml:"events"`
	Output OutputConfig `yaml:"output"`
}

// EventsConfig represents event filtering configuration
type EventsConfig struct {
	Include []string             `yaml:"include"`
	Rules   map[string]EventRule `yaml:"rules"`
}

// EventRule represents filtering rules for an event type
type EventRule struct {
	Keep    []string `yaml:"keep"`
	Discard []string `yaml:"discard"`
}

// OutputConfig represents output configuration
type OutputConfig struct {
	Format string `yaml:"format"`
	Pretty bool   `yaml:"pretty"`
}

// BufferState represents the state buffer for generating done event
type BufferState struct {
	SessionID string
	Ok        bool
	LastError string
}

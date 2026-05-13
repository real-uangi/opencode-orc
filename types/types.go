package types

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
	OutputTypeTool    = "tool"
	OutputTypeText    = "text"
	OutputTypeStep    = "step"
	OutputTypeDone    = "done"
)

// RawEvent represents a parsed JSONL event from opencode
type RawEvent struct {
	Type      string                 `json:"type"`
	Timestamp int64                  `json:"timestamp"`
	SessionID string                 `json:"sessionId"`
	Part      map[string]interface{} `json:"part,omitempty"`
	Error     map[string]interface{} `json:"error,omitempty"`
}

// SessionEvent represents a session info event
type SessionEvent struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionId"`
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

// TextEvent represents a text output event
type TextEvent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// StepEvent represents a step finish event
type StepEvent struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

// DoneEvent represents the final done event
type DoneEvent struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionId"`
	Ok        bool   `json:"ok"`
	Error     string `json:"error,omitempty"`
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

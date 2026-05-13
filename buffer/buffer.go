package buffer

import (
	"opencode-orc/types"
)

// Buffer manages state for generating the done event
type Buffer struct {
	state types.BufferState
}

// New creates a new buffer
func New() *Buffer {
	return &Buffer{
		state: types.BufferState{
			Ok: true,
		},
	}
}

// Update updates the buffer state based on an event
func (b *Buffer) Update(event interface{}) {
	switch e := event.(type) {
	case *types.SessionEvent:
		if b.state.SessionID == "" {
			b.state.SessionID = e.SessionID
		}
	case *types.ToolEvent:
		if e.Status == "error" && e.Tool != "error" {
			b.state.Ok = false
			b.state.LastError = e.Error
		}
	}
}

// UpdateFromError updates the buffer state from an error event
func (b *Buffer) UpdateFromError(name, message string) {
	b.state.Ok = false
	b.state.LastError = name + ": " + message
}

// DoneEvent generates the final done event
func (b *Buffer) DoneEvent() *types.DoneEvent {
	return &types.DoneEvent{
		Type:      types.OutputTypeDone,
		SessionID: b.state.SessionID,
		Ok:        b.state.Ok,
		Error:     b.state.LastError,
	}
}

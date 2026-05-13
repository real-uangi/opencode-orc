package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Writer handles output formatting and writing
type Writer struct {
	writer io.Writer
	pretty bool
}

// New creates a new output writer
func New(w io.Writer, pretty bool) *Writer {
	return &Writer{
		writer: w,
		pretty: pretty,
	}
}

// WriteEvent writes an event as JSONL
func (w *Writer) WriteEvent(event interface{}) error {
	var data []byte
	var err error

	if w.pretty {
		data, err = json.MarshalIndent(event, "", "  ")
	} else {
		data, err = json.Marshal(event)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	_, err = w.writer.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write event: %w", err)
	}

	_, err = w.writer.Write([]byte("\n"))
	if err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	return nil
}

// WriteError writes an error message to stderr
func WriteError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[opencode-orc] "+format+"\n", args...)
}

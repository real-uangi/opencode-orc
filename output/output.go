package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/real-uangi/opencode-orc/types"
)

// Writer handles output formatting and writing
type Writer struct {
	writer io.Writer
	format string
	pretty bool
}

// New creates a new output writer
// format: "jsonl" or "text"
func New(w io.Writer, format string, pretty bool) *Writer {
	return &Writer{
		writer: w,
		format: format,
		pretty: pretty,
	}
}

// WriteEvent writes an event in the configured format
func (w *Writer) WriteEvent(event interface{}) error {
	if w.format == "text" {
		return w.writeText(event)
	}
	return w.writeJSON(event)
}

func (w *Writer) writeText(event interface{}) error {
	var line string
	if f, ok := event.(types.Formatter); ok {
		line = f.FormatText()
	} else {
		line = fmt.Sprintf("%v", event)
	}

	_, err := w.writer.Write([]byte(line + "\n"))
	if err != nil {
		return fmt.Errorf("failed to write event: %w", err)
	}
	return nil
}

func (w *Writer) writeJSON(event interface{}) error {
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

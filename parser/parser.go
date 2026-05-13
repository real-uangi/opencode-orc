package parser

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"strings"

	"opencode-orc/types"
)

// Parser reads JSONL input and emits parsed events
type Parser struct {
	reader *bufio.Reader
}

// NewParser creates a new JSONL parser
func NewParser(r io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(r),
	}
}

// ParseNext reads and parses the next JSONL line
// Returns nil, nil when EOF is reached
// Returns nil, error for parse errors (non-JSON lines are skipped)
func (p *Parser) ParseNext() (*types.RawEvent, error) {
	for {
		line, err := p.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if strings.TrimSpace(line) == "" {
					return nil, io.EOF
				}
				// Process last line without newline
			} else {
				return nil, err
			}
		}

		line = strings.TrimSpace(line)
		if line == "" {
			if err == io.EOF {
				return nil, io.EOF
			}
			continue
		}

		// Try to parse as JSON
		var event types.RawEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			// Skip non-JSON lines
			log.Printf("[opencode-orc] skipping non-JSON line: %s", line[:min(50, len(line))])
			if err == io.EOF {
				return nil, io.EOF
			}
			continue
		}

		// Validate event type
		if event.Type == "" {
			log.Printf("[opencode-orc] skipping event with no type: %s", line[:min(50, len(line))])
			if err == io.EOF {
				return nil, io.EOF
			}
			continue
		}

		return &event, nil
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

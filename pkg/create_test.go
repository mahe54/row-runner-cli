package pkg

import (
	"fmt"
	"testing"

	"github.com/vbauerster/mpb/v8/decor"
)

func TestStatusUpdate(t *testing.T) {
	status := &Status{Current: Starting}
	text := "Processing 1/10"
	ws := uint(11)
	wcc := decor.WC{W: len(string(status.Current)) + len(text) + 2}

	decorator := statusUpdate(status, text, ws, wcc)

	// Create a mock Statistics
	stats := decor.Statistics{}

	// Call the decorator's Decor method
	result := decorator.Decor(stats)

	expected := fmt.Sprintf(" %s %s", status.Current, text)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

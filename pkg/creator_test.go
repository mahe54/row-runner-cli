package pkg

import (
	"strings"
	"testing"

	"github.com/vbauerster/mpb/v8"
)

type MockCreator struct{}

func (m *MockCreator) ProcessInput(input interface{}, bar *mpb.Bar, progress chan<- int, cancel <-chan struct{}, statusChan chan<- Status, status *Status) {
	// Mock implementation
}

func (m *MockCreator) ConvertToInput(data []string) interface{} {
	// Mock implementation
	return strings.Join(data, ",")
}

func TestReadDataFromReader(t *testing.T) {
	reader := strings.NewReader("test1,test2\ntest3,test4")
	creator := &MockCreator{}

	data, err := ReadDataFromReader(reader, creator)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(data) != 2 {
		t.Errorf("Expected 2 rows, got %v", len(data))
	}

	if data[0] != "test1,test2" || data[1] != "test3,test4" {
		t.Errorf("Unexpected data: %v", data)
	}
}

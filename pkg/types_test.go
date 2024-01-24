package pkg

import "testing"

func TestSet(t *testing.T) {
	s := &Status{}
	s.Set(Working)

	if s.Current != Working {
		t.Errorf("Expected status to be %v, but got %v", Working, s.Current)
	}
}

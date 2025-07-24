package xuuid

import (
	"testing"
)

func TestNew(t *testing.T) {
	id := New()

	if IsNilOrEmpty(id) {
		t.Errorf("expected uuid to be valid")
	}
}

func TestParseString(t *testing.T) {
	id, err := ParseString("786d2eb1-3260-429e-a65f-36d8f7889c04")
	if err != nil {
		t.Fatalf("expected err to be nil. got %v", err)
	}

	if id.IsNil() {
		t.Fatalf("expected id to be valid. got %v", id)
	}

	_, err = ParseString("6cc67735-ff0e-41fe-b721-any")
	if err == nil {
		t.Fatalf("expected invalid uuid string to be invalid. Got valid")
	}
}

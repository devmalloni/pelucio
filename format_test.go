package pelucio

import (
	"math/big"
	"testing"
)

func TestFromString(t *testing.T) {
	f, _ := FromString("1.23", 2)
	expected := big.NewInt(123)
	if f.Cmp(expected) != 0 {
		t.Errorf("FromString failed: got %v, want %v", f, expected)
	}
}

func TestToString(t *testing.T) {
	b := big.NewInt(123)
	s := ToString(b, 2)
	expected := "1.23"
	if s != expected {
		t.Errorf("ToString failed: got %v, want %v", s, expected)
	}
}

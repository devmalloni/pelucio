package xmap

import (
	"testing"
)

func TestToMap(t *testing.T) {
	s := []struct {
		A string
		B string
	}{}

	s = append(s, struct {
		A string
		B string
	}{"first", "second"}, struct {
		A string
		B string
	}{"third", "fourth"})
	res := ToMap(s, func(s struct {
		A string
		B string
	}) string {
		return s.A
	})

	if len(res) != 2 {
		t.Errorf("expected len to be 2, got %v", len(res))
	}

	if v, ok := res["first"]; !ok || v.B != "second" {
		t.Errorf("expected second string, got %v", v)
	}

	if v, ok := res["third"]; !ok || v.B != "fourth" {
		t.Errorf("expected second string, got %v", v)
	}
}

func TestValues(t *testing.T) {
	s := map[string]struct {
		A string
		B string
	}{
		"first":  {"a", "b"},
		"second": {"c", "d"},
	}

	v := Values(s)

	if len(v) != 2 {
		t.Fatalf("expected array result to have len 2, got %v", len(v))
	}
}

package xtime

import "time"

type Clock interface {
	Now() time.Time
	NilNow() *time.Time
}

type StdClock struct{}

func (p StdClock) Now() time.Time {
	return time.Now()
}

func (p StdClock) NilNow() *time.Time {
	n := time.Now()
	return &n
}

type StubClock struct {
	t time.Time
}

func NewStubClock(t time.Time) StubClock {
	return StubClock{t}
}

func (p StubClock) Now() time.Time {
	return p.t
}

func (p StubClock) NilNow() *time.Time {
	return &p.t
}

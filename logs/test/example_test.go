package test

import "testing"

type Assert struct {
	t *testing.T
}

func NewAssert(t *testing.T) *Assert {
	return &Assert{t: t}
}

func (a *Assert) Equal(got, want any) {
	a.t.Helper()
	if got != want {
		a.t.Errorf("got %v, want %v", got, want)
	}
}

func (a *Assert) NotEqual(got, want any) {
	a.t.Helper()
	if got == want {
		a.t.Errorf("got %v, want %v", got, want)
	}
}

func (a *Assert) True(got bool) {
	a.t.Helper()
	if !got {
		a.t.Errorf("got %v, want true", got)
	}
}

func (a *Assert) False(got bool) {
	a.t.Helper()
	if got {
		a.t.Errorf("got %v, want false", got)
	}
}

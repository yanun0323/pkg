package tester

import (
	"errors"
	"testing"
)

type require[T comparable] struct {
	t *testing.T
}

func NewRequire[T comparable](t *testing.T) *require[T] {
	return &require[T]{t: t}
}

func (a *require[T]) Equal(want, got T) {
	a.t.Helper()
	if got != want {
		a.t.Fatalf("got %v, want %v", got, want)
	}
}

func (a *require[T]) NotEqual(want, got T) {
	a.t.Helper()
	if got == want {
		a.t.Fatalf("got %v, want %v", got, want)
	}
}

func (a *require[T]) Nil(got T) {
	a.t.Helper()
	if !IsNil(got) {
		a.t.Fatalf("got %v, want nil", got)
	}
}

func (a *require[T]) NotNil(got T) {
	a.t.Helper()
	if IsNil(got) {
		a.t.Fatalf("got %v, want not nil", got)
	}
}

func (a *require[T]) Error(got error) {
	a.t.Helper()
	if got == nil {
		a.t.Fatalf("got %v, want error", got)
	}
}

func (a *require[T]) NoError(got error) {
	a.t.Helper()
	if got != nil {
		a.t.Fatalf("got %v, want no error", got)
	}
}

func (a *require[T]) True(got bool) {
	a.t.Helper()
	if !got {
		a.t.Fatalf("got %v, want true", got)
	}
}

func (a *require[T]) False(got bool) {
	a.t.Helper()
	if got {
		a.t.Fatalf("got %v, want false", got)
	}
}

func (a *require[T]) ErrorIs(want error, got error) {
	a.t.Helper()
	if !errors.Is(got, want) {
		a.t.Fatalf("got %v, want %v", got, want)
	}
}

func (a *require[T]) ErrorIsNot(want error, got error) {
	a.t.Helper()
	if errors.Is(got, want) {
		a.t.Fatalf("got %v, want %v", got, want)
	}
}

func (a *require[T]) ErrorAs(got error) {
	a.t.Helper()
	var want T
	if !errors.As(got, &want) {
		a.t.Fatalf("got %v, want %v", got, want)
	}
}

func (a *require[T]) ErrorAsNot(got error) {
	a.t.Helper()
	var want T
	if errors.As(got, &want) {
		a.t.Fatalf("got %v, want %v", got, want)
	}
}

func RequireEqual[T comparable](t *testing.T, want, got T) {
	t.Helper()
	require := NewRequire[T](t)
	require.Equal(want, got)
}

func RequireNotEqual[T comparable](t *testing.T, want, got T) {
	t.Helper()
	require := NewRequire[T](t)
	require.NotEqual(want, got)
}

func RequireNil[T comparable](t *testing.T, got T) {
	t.Helper()
	require := NewRequire[T](t)
	require.Nil(got)
}

func RequireNotNil[T comparable](t *testing.T, got T) {
	t.Helper()
	require := NewRequire[T](t)
	require.NotNil(got)
}

func RequireError(t *testing.T, got error) {
	t.Helper()
	require := NewRequire[error](t)
	require.Error(got)
}

func RequireNoError(t *testing.T, got error) {
	t.Helper()
	require := NewRequire[any](t)
	require.NoError(got)
}

func RequireTrue(t *testing.T, got bool) {
	t.Helper()
	require := NewRequire[bool](t)
	require.True(got)
}

func RequireFalse(t *testing.T, got bool) {
	t.Helper()
	require := NewRequire[bool](t)
	require.False(got)
}

func RequireErrorIs(t *testing.T, want error, got error) {
	t.Helper()
	require := NewRequire[error](t)
	require.ErrorIs(want, got)
}

func RequireErrorIsNot(t *testing.T, want error, got error) {
	t.Helper()
	require := NewRequire[error](t)
	require.ErrorIsNot(want, got)
}

func RequireErrorAs[T error](t *testing.T, got error) {
	t.Helper()
	require := NewRequire[error](t)
	require.ErrorAs(got)
}

func RequireErrorAsNot[T error](t *testing.T, got error) {
	t.Helper()
	require := NewRequire[error](t)
	require.ErrorAsNot(got)
}

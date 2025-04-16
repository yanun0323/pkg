package test

import (
	"errors"
	"testing"
)

type assert[T comparable] struct {
	t *testing.T
}

func NewAssert[T comparable](t *testing.T) *assert[T] {
	return &assert[T]{t: t}
}

func (a *assert[T]) Equal(got, want T) {
	a.t.Helper()
	if got != want {
		a.t.Errorf("got %v, want %v", got, want)
	}
}

func (a *assert[T]) NotEqual(got, want T) {
	a.t.Helper()
	if got == want {
		a.t.Errorf("got %v, want %v", got, want)
	}
}

func (a *assert[T]) Nil(got T) {
	a.t.Helper()
	if !IsNil(got) {
		a.t.Errorf("got %v, want nil", got)
	}
}

func (a *assert[T]) NotNil(got T) {
	a.t.Helper()
	if IsNil(got) {
		a.t.Errorf("got %v, want not nil", got)
	}
}

func (a *assert[T]) Error(got error) {
	a.t.Helper()
	if got == nil {
		a.t.Errorf("got %v, want error", got)
	}
}

func (a *assert[T]) NoError(got error) {
	a.t.Helper()
	if got != nil {
		a.t.Errorf("got %v, want no error", got)
	}
}

func (a *assert[T]) True(got bool) {
	a.t.Helper()
	if !got {
		a.t.Errorf("got %v, want true", got)
	}
}

func (a *assert[T]) False(got bool) {
	a.t.Helper()
	if got {
		a.t.Errorf("got %v, want false", got)
	}
}

func (a *assert[T]) ErrorIs(got error, want error) {
	a.t.Helper()
	if !errors.Is(got, want) {
		a.t.Errorf("got %v, want %v", got, want)
	}
}

func (a *assert[T]) ErrorIsNot(got error, want error) {
	a.t.Helper()
	if errors.Is(got, want) {
		a.t.Errorf("got %v, want %v", got, want)
	}
}

func (a *assert[T]) ErrorAs(got error) {
	a.t.Helper()
	var want T
	if !errors.As(got, &want) {
		a.t.Errorf("got %v, want %v", got, want)
	}
}

func (a *assert[T]) ErrorAsNot(got error) {
	a.t.Helper()
	var want T
	if errors.As(got, &want) {
		a.t.Errorf("got %v, want %v", got, want)
	}
}

func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	assert := NewAssert[T](t)
	assert.Equal(got, want)
}

func AssertNotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	assert := NewAssert[T](t)
	assert.NotEqual(got, want)
}

func AssertNil[T comparable](t *testing.T, got T) {
	t.Helper()
	assert := NewAssert[T](t)
	assert.Nil(got)
}

func AssertNotNil[T comparable](t *testing.T, got T) {
	t.Helper()
	assert := NewAssert[T](t)
	assert.NotNil(got)
}

func AssertError(t *testing.T, got error) {
	t.Helper()
	assert := NewAssert[error](t)
	assert.Error(got)
}

func AssertNoError(t *testing.T, got error) {
	t.Helper()
	assert := NewAssert[error](t)
	assert.NoError(got)
}

func AssertTrue(t *testing.T, got bool) {
	t.Helper()
	assert := NewAssert[bool](t)
	assert.True(got)
}

func AssertFalse(t *testing.T, got bool) {
	t.Helper()
	assert := NewAssert[bool](t)
	assert.False(got)
}

func AssertErrorIs(t *testing.T, got error, want error) {
	assert := NewAssert[error](t)
	assert.ErrorIs(got, want)
}

func AssertErrorIsNot(t *testing.T, got error, want error) {
	assert := NewAssert[error](t)
	assert.ErrorIsNot(got, want)
}

func AssertErrorAs[T error](t *testing.T, got error) {
	assert := NewAssert[error](t)
	assert.ErrorAs(got)
}

func AssertErrorAsNot[T error](t *testing.T, got error) {
	assert := NewAssert[error](t)
	assert.ErrorAsNot(got)
}

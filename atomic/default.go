package atomic

import "sync/atomic"

type Bool atomic.Bool

type Int32 atomic.Int32

type Int64 atomic.Int64

type Uint32 atomic.Uint32

type Uint64 atomic.Uint64

type Uintptr atomic.Uintptr

type Pointer[T any] atomic.Pointer[T]

// noCopy may be added to structs which must not be copied
// after the first use.
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
//
// Note that it must not be embedded, due to the Lock and Unlock methods.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

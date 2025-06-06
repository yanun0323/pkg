package sync

import "sync"

type Mutex sync.Mutex

type RWMutex sync.RWMutex

type Once sync.Once

type Cond sync.Cond

type WaitGroup sync.WaitGroup

type Locker sync.Locker

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

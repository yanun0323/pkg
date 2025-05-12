package sync

import (
	"log"
	"testing"
	"time"
)

func TestGoWithoutTimeout(t *testing.T) {
	ch := make(chan struct{})
	Go(func() {
		time.Sleep(time.Second)
		log.Println("done")
		ch <- struct{}{}
	})

	select {
	case <-ch:
	case <-time.After(time.Second * 2):
		t.Fatal("timeout")
	}

	time.Sleep(2 * time.Second)
}

func TestGoWithTimeout(t *testing.T) {
	ch := make(chan struct{})
	Go(func() {
		time.Sleep(3 * time.Second)
		log.Println("done")
		ch <- struct{}{}
	}, time.Second)

	select {
	case <-ch:
	case <-time.After(time.Second * 5):
		t.Fatal("timeout")
	}

	time.Sleep(2 * time.Second)
}

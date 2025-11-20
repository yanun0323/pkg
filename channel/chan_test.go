package channel

import (
	"testing"

	"github.com/yanun0323/pkg/tester"
)

func TestTryPushClosedChannel(t *testing.T) {
	ch := make(chan int, 10)
	close(ch)

	ok := TryPush(ch, 10)
	a := tester.NewAssert[int](t)
	a.False(ok)
}

func TestTryReceiveClosedChannel(t *testing.T) {
	ch := make(chan int, 10)
	close(ch)

	val, ok := TryReceive(ch)
	a := tester.NewAssert[int](t)
	a.False(ok)
	a.Equal(0, val)
}

func TestSafeClose(t *testing.T) {
	ch := make(chan int, 10)
	SafeClose(ch)
	SafeClose(ch)
	SafeClose(ch)
	SafeClose(ch)

	val, ok := <-ch
	a := tester.NewAssert[int](t)
	a.False(ok)
	a.Equal(0, val)
}

func TestIsClose(t *testing.T) {
	ch := make(chan int, 10)
	ch <- 99
	ch <- 999

	a := tester.NewAssert[int](t)
	a.False(IsClose(ch))

	val, ok := <-ch
	a.True(ok)
	a.Equal(99, val)

	close(ch)

	a.True(IsClose(ch))

	val, ok = <-ch
	a.True(ok)
	a.Equal(999, val)

	val, ok = <-ch
	a.False(ok)
	a.Equal(0, val)

}

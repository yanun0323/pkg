package channel

// SafeClose closes a channel if it is open.
//
// It is safe to call this function even if the channel is already closed.
func SafeClose[T any](ch chan T) {
	select {
	case _, ok := <-ch:
		if ok {
			close(ch)
		}
	default:
		close(ch)
	}
}

// NoneBlockingPush pushes a value to a channel. If the channel is full, it will skip the push without blocking.
//
// It is safe to call this function even if the channel is already closed.
func NoneBlockingPush[T any](ch chan<- T, data T) {
	select {
	case ch <- data:
	default:
	}
}

// NoneBlockingReceive receives a value from a channel. If the channel is empty, it will return a zero value and false without blocking.
//
// It is safe to call this function even if the channel is already closed.
func NoneBlockingReceive[T any](ch chan T) (T, bool) {
	select {
	case data := <-ch:
		return data, true
	default:
		return *new(T), false
	}
}

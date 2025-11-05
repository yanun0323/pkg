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

// TryPush pushes a value to a channel. If the channel is full, it will skip the push without blocking.
//
// It is safe to call this function even if the channel is already closed.
func TryPush[T any](ch chan<- T, data T) {
	select {
	case ch <- data:
	default:
	}
}

// TryReceive receives a value from a channel. If the channel is empty, it will return a zero value and false without blocking.
//
// It is safe to call this function even if the channel is already closed.
func TryReceive[T any](ch chan T) (T, bool) {
	select {
	case data, ok := <-ch:
		return data, ok
	default:
		return *new(T), false
	}
}

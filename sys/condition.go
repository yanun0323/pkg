package sys

// If returns trueV if cond is true, otherwise returns falseV
//
// # Caution:
//   - Both trueV and falseV are invoked before the condition is checked, ensure they are safe to execute to avoid runtime panics.
func If[T any](cond bool, trueV T, falseV T) T {
	if cond {
		return trueV
	}

	return falseV
}

package sys

func If[T any](cond bool, trueV T, falseV T) T {
	if cond {
		return trueV
	}

	return falseV
}

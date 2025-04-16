package sys

import (
	"github.com/yanun0323/pkg/test"
)

func IsNil(a any) bool {
	return test.IsNil(a)
}

func IsNotNil(a any) bool {
	return !IsNil(a)
}

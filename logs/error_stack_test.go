package logs

import (
	stderr "errors"
	"testing"

	"github.com/pkg/errors"
)

func TestErrStack(t *testing.T) {
	t.Log(GetStack(getError()))
	t.Log(GetStack(getNormalErr()))
}

func getError() error {
	return errors.WithStack(errors.New("test error"))
}

func getNormalErr() error {
	return stderr.New("test normal error")
}

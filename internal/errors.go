package internal

import (
	"fmt"
)

var (
	_ error = new(StopError)
)

type StopError struct {
	v   interface{}
	err error
}

func (e *StopError) Error() string {
	return fmt.Sprintf("occurred error on stopping %T. %s", e.v, e.err.Error())
}

func (e *StopError) Unwrap() error {
	return e.err
}

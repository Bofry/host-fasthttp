package response

type Error struct {
	msg string
	err error
}

func NewError(msg string) *Error {
	return &Error{
		msg: msg,
	}
}

func WrapError(err error, msg string) *Error {
	if err == nil {
		panic("argument err should not be nil")
	}

	if msg == "" {
		msg = err.Error()
	}

	return &Error{
		msg: msg,
		err: err,
	}
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Unwrap() error {
	return e.err
}

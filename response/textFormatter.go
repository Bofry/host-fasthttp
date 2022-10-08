package response

import (
	"fmt"

	http "github.com/valyala/fasthttp"
)

var _ ContentFormatter = TextFormatter("")

type TextFormatter string

func (f TextFormatter) Success(ctx *http.RequestCtx, body interface{}) error {
	buf, err := f.marshalBody(body)
	if err != nil {
		return err
	}

	Success(ctx, CONTENT_TYPE_TEXT, buf)
	return nil
}

func (f TextFormatter) Failure(ctx *http.RequestCtx, body interface{}, statusCode int) error {
	buf, err := f.marshalBody(body)
	if err != nil {
		return err
	}

	Failure(ctx, CONTENT_TYPE_TEXT, buf, statusCode)
	return nil
}

func (f TextFormatter) marshalBody(body interface{}) ([]byte, error) {
	var buf []byte

	switch v := body.(type) {
	case []byte:
		buf = v
	case string:
		buf = []byte(v)
	default:
		return nil, NewError(fmt.Sprintf("cannot cast %T to []byte", body))
	}
	return buf, nil
}

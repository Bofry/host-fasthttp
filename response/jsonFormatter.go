package response

import (
	"encoding/json"

	http "github.com/valyala/fasthttp"
)

var _ ContentFormatter = JsonFormatter("")

type JsonFormatter string

func (f JsonFormatter) Success(ctx *http.RequestCtx, body interface{}) error {
	buf, err := f.marshalBody(body)
	if err != nil {
		return err
	}
	Success(ctx, CONTENT_TYPE_JSON, buf)
	return nil
}

func (f JsonFormatter) Failure(ctx *http.RequestCtx, body interface{}, statusCode int) error {
	buf, err := f.marshalBody(body)
	if err != nil {
		return err
	}
	Failure(ctx, CONTENT_TYPE_JSON, buf, statusCode)
	return nil
}

func (f JsonFormatter) marshalBody(body interface{}) ([]byte, error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, WrapError(err, err.Error())
	}
	return buf, nil
}

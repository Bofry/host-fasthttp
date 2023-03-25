package response

import (
	_ "unsafe"

	"github.com/Bofry/host-fasthttp/internal/responseutil"
	http "github.com/valyala/fasthttp"
)

const (
	// response flag
	SUCCESS = responseutil.SUCCESS
	FAILURE = responseutil.FAILURE
	UNKNOWN = responseutil.UNKNOWN
)

const (
	CONTENT_TYPE_JSON string = "application/json; charset=utf-8"
	CONTENT_TYPE_TEXT string = "text/plain; charset=utf-8"
)

const (
	Json = JsonFormatter("")
	Text = TextFormatter("")
)

type (
	ResponseFlag  = responseutil.ResponseFlag
	ResponseState = responseutil.ResponseState
)

type ContentFormatter interface {
	Success(ctx *http.RequestCtx, body interface{}) error
	Failure(ctx *http.RequestCtx, body interface{}, statusCode int) error
}

//go:linkname ExtractResponseState github.com/Bofry/host-fasthttp/internal/responseutil.ExtractResponseState
func ExtractResponseState(ctx *http.RequestCtx) ResponseState

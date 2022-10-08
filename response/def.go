package response

import http "github.com/valyala/fasthttp"

const (
	// response flag
	SUCCESS ResponseFlag = iota
	FAILURE

	UNKNOWN ResponseFlag = -1
)

const (
	// response name in RequestCtx user store
	USER_STORE_RESPONSE_FLAG string = "github.com/Bofry/host-fasthttp/response::Response"

	CONTENT_TYPE_JSON string = "application/json; charset=utf-8"
	CONTENT_TYPE_TEXT string = "text/plain; charset=utf-8"
)

const (
	Json = JsonFormatter("")
	Text = TextFormatter("")
)

type ResponseFlag int

type Response interface {
	Flag() ResponseFlag
	StatusCode() int
}

type ContentFormatter interface {
	Success(ctx *http.RequestCtx, body interface{}) error
	Failure(ctx *http.RequestCtx, body interface{}, statusCode int) error
}

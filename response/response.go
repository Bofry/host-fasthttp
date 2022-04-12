package response

import http "github.com/valyala/fasthttp"

const (
	// response name in RequestCtx user store
	RESPONSE_INVARIANT_NAME string = "github.com/Bofry/host-fasthttp/response::Response"
)

func Success(ctx *http.RequestCtx, contentType string, body []byte) {
	ctx.Success(contentType, body)

	storeResponse(
		ctx,
		&responseImpl{
			flag:       SUCCESS,
			statusCode: ctx.Response.StatusCode(),
		},
	)
}

func Failure(ctx *http.RequestCtx, contentType string, message []byte, statusCode int) {
	ctx.SetStatusCode(statusCode)
	ctx.Success(contentType, message)

	storeResponse(
		ctx,
		&responseImpl{
			flag:       FAILURE,
			statusCode: statusCode,
		},
	)
}

func storeResponse(ctx *http.RequestCtx, resp Response) {
	ctx.SetUserValue(RESPONSE_INVARIANT_NAME, resp)
}

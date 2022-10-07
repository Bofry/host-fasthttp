package response

import http "github.com/valyala/fasthttp"

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

func GetResponseFlag(ctx *http.RequestCtx) Response {
	obj := ctx.UserValue(USER_STORE_RESPONSE_FLAG)
	v, ok := obj.(Response)
	if ok {
		return v
	}
	return nil
}

func storeResponse(ctx *http.RequestCtx, resp Response) {
	ctx.SetUserValue(USER_STORE_RESPONSE_FLAG, resp)
}

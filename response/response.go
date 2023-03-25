package response

import (
	"github.com/Bofry/host-fasthttp/internal/responseutil"
	http "github.com/valyala/fasthttp"
)

func Success(ctx *http.RequestCtx, contentType string, body []byte) {
	ctx.Success(contentType, body)

	responseutil.InjectResponseState(
		ctx,
		responseutil.CreateResponseState(
			SUCCESS,
			ctx.Response.StatusCode()),
	)
}

func Failure(ctx *http.RequestCtx, contentType string, message []byte, statusCode int) {
	ctx.SetStatusCode(statusCode)
	ctx.Success(contentType, message)

	responseutil.InjectResponseState(
		ctx,
		responseutil.CreateResponseState(
			FAILURE,
			statusCode,
		),
	)
}

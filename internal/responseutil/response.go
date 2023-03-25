package responseutil

import (
	http "github.com/valyala/fasthttp"

	"github.com/Bofry/host-fasthttp/internal/requestutil"
)

func CreateResponseState(flag ResponseFlag, statusCode int) ResponseState {
	return &ResponseStateImpl{
		flag:       flag,
		statusCode: statusCode,
	}
}

func ExtractResponseState(ctx *http.RequestCtx) ResponseState {
	obj := requestutil.ExtractResponseState(ctx)
	v, ok := obj.(ResponseState)
	if ok {
		return v
	}
	return nil
}

func InjectResponseState(ctx *http.RequestCtx, state ResponseState) {
	requestutil.InjectResponseState(ctx, state)
}

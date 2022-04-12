package middleware

import (
	"reflect"
	"unsafe"

	"github.com/Bofry/host-fasthttp/internal"
)

func isRequestHandler(rv reflect.Value) bool {
	if rv.IsValid() {
		return rv.Type().AssignableTo(typeOfRequestHandler)
	}
	return false
}

func asRequestHandler(rv reflect.Value) internal.RequestHandler {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfRequestHandler).Interface().(internal.RequestHandler); ok {
			return v
		}
	}
	return nil
}

func asFasthttpHost(rv reflect.Value) *internal.FasthttpHost {
	return reflect.NewAt(typeOfHost, unsafe.Pointer(rv.Pointer())).
		Interface().(*internal.FasthttpHost)
}

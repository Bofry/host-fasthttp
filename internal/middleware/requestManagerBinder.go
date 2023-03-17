package middleware

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Bofry/host"
	"github.com/Bofry/host-fasthttp/internal"
	"github.com/Bofry/structproto"
	"github.com/Bofry/structproto/tagresolver"
	"github.com/Bofry/structproto/util/reflectutil"
)

var _ structproto.StructBinder = new(RequestManagerBinder)

type RequestManagerBinder struct {
	registrar  *internal.FasthttpHostRegistrar
	appContext *host.AppContext
}

func (b *RequestManagerBinder) Init(context *structproto.StructProtoContext) error {
	return nil
}

func (b *RequestManagerBinder) Bind(field structproto.FieldInfo, rv reflect.Value) error {
	if !rv.IsValid() {
		return fmt.Errorf("specifiec argument 'rv' is invalid")
	}

	// assign zero if rv is nil
	rvRequestHandler := reflectutil.AssignZero(rv)
	binder := &RequestHandlerBinder{
		requestHandlerType: rvRequestHandler.Type().Name(),
		components: map[string]reflect.Value{
			host.APP_CONFIG_FIELD:           b.appContext.Config(),
			host.APP_SERVICE_PROVIDER_FIELD: b.appContext.ServiceProvider(),
		},
	}
	err := b.preformBindRequestHandler(rvRequestHandler, binder)
	if err != nil {
		return err
	}

	// register RequestHandlers
	return b.registerRoute(field.Name(), rvRequestHandler)
}

func (b *RequestManagerBinder) Deinit(context *structproto.StructProtoContext) error {
	return nil
}

func (b *RequestManagerBinder) preformBindRequestHandler(target reflect.Value, binder *RequestHandlerBinder) error {
	prototype, err := structproto.Prototypify(target,
		&structproto.StructProtoResolveOption{
			TagResolver: tagresolver.NoneTagResolver,
		})
	if err != nil {
		return err
	}

	return prototype.Bind(binder)
}

func (b *RequestManagerBinder) registerRoute(url string, rvRequestHandler reflect.Value) error {
	// register RequestHandlers
	count := rvRequestHandler.Type().NumMethod()
	for i := 0; i < count; i++ {
		method := rvRequestHandler.Type().Method(i)

		rvMethod := rvRequestHandler.Method(method.Index)
		if isRequestHandler(rvMethod) {
			handler := asRequestHandler(rvMethod)
			if handler != nil {
				// TODO: validate path make comply RFC3986
				b.registrar.AddRoute(strings.ToUpper(method.Name), url, handler)
			}
		}
	}
	return nil
}

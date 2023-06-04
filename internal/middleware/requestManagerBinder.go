package middleware

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Bofry/host"
	"github.com/Bofry/host-fasthttp/internal"
	"github.com/Bofry/structproto"
	"github.com/Bofry/structproto/reflecting"
	"github.com/Bofry/structproto/tagresolver"
)

var _ structproto.StructBinder = new(RequestManagerBinder)

type RequestManagerBinder struct {
	registrar *internal.FasthttpHostRegistrar
	app       *host.AppModule
}

func (b *RequestManagerBinder) Init(context *structproto.StructProtoContext) error {
	return nil
}

func (b *RequestManagerBinder) Bind(field structproto.FieldInfo, rv reflect.Value) error {
	if !rv.IsValid() {
		return fmt.Errorf("specifiec argument 'rv' is invalid")
	}

	// assign zero if rv is nil
	rvRequestHandler := reflecting.AssignZero(rv)
	binder := &RequestHandlerBinder{
		requestHandlerType: rvRequestHandler.Type().Name(),
		components: map[string]reflect.Value{
			host.APP_CONFIG_FIELD:           b.app.Config(),
			host.APP_SERVICE_PROVIDER_FIELD: b.app.ServiceProvider(),
		},
	}
	err := b.bindRequestHandler(rvRequestHandler, binder)
	if err != nil {
		return err
	}

	// register RequestHandlers
	return b.registerRoute(field.IDName(), field.Name(), rvRequestHandler)
}

func (b *RequestManagerBinder) Deinit(context *structproto.StructProtoContext) error {
	return nil
}

func (b *RequestManagerBinder) bindRequestHandler(target reflect.Value, binder *RequestHandlerBinder) error {
	prototype, err := structproto.Prototypify(target,
		&structproto.StructProtoResolveOption{
			TagResolver: tagresolver.NoneTagResolver,
		})
	if err != nil {
		return err
	}

	return prototype.Bind(binder)
}

func (b *RequestManagerBinder) registerRoute(moduleID, url string, rvRequestHandler reflect.Value) error {
	// register RequestHandlers
	count := rvRequestHandler.Type().NumMethod()
	for i := 0; i < count; i++ {
		method := rvRequestHandler.Type().Method(i)

		rvHandler := rvRequestHandler.Method(method.Index)
		if isRequestHandler(rvHandler) {
			handler := asRequestHandler(rvHandler)
			if handler != nil {
				// TODO: validate path make comply RFC3986
				b.registrar.AddRoute(strings.ToUpper(method.Name), url, handler, moduleID)
			}
		}
	}
	return nil
}

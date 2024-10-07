package middleware

import (
	"fmt"
	"os"
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
	var (
		moduleID = field.IDName()
		url      = field.Name()
	)

	{
		// NOTE: @ExpandEnv:"on"
		optExpandEnv := field.Tag().Get(TAG_OPT_EXPAND_ENV)
		if optExpandEnv != "off" || len(optExpandEnv) == 0 || optExpandEnv == "on" {
			url = os.ExpandEnv(url)
		}
	}

	{
		// NOTE: @BindMethod:"GET *FETCH"
		optBindMethod := field.Tag().Get(TAG_OPT_BIND_METHOD)
		if optBindMethod != "" {
			parts := strings.Split(optBindMethod, " ")
			if len(parts) < 2 {
				return fmt.Errorf("cannot resolve @BindMethod '%s' on %s", optBindMethod, moduleID)
			}
			var (
				httpMethod string = strings.ToUpper(strings.Trim(parts[0], " "))
				bindMethod string = strings.ToUpper(strings.Trim(parts[1], " "))
			)
			for i := 1; i < len(parts); i++ {
				part := parts[i]
				switch part[0] {
				case '*':
					bindMethod = strings.ToUpper(strings.Trim(part[1:], " "))
				default:
					return fmt.Errorf("cannot resolve @BindMethod '%s' on %s", optBindMethod, moduleID)
				}
			}

			err := b.registerRouteWithMethod(moduleID, url, httpMethod, bindMethod, rvRequestHandler)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return b.registerRoute(moduleID, url, rvRequestHandler)
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

func (b *RequestManagerBinder) registerRouteWithMethod(moduleID, url string, httpMethod string, bindMethod string, rvRequestHandler reflect.Value) error {
	// register RequestHandlers
	count := rvRequestHandler.Type().NumMethod()
	for i := 0; i < count; i++ {
		method := rvRequestHandler.Type().Method(i)

		rvHandler := rvRequestHandler.Method(method.Index)
		if isRequestHandler(rvHandler) {
			if strings.ToUpper(method.Name) == bindMethod {
				handler := asRequestHandler(rvHandler)
				if handler != nil {
					// TODO: validate path make comply RFC3986
					b.registrar.AddRoute(strings.ToUpper(httpMethod), url, handler, moduleID)
					return nil
				}
			}
		}
	}
	return fmt.Errorf("cannot find method '%s' on %s", bindMethod, moduleID)
}

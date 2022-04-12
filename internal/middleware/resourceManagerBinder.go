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

var _ structproto.StructBinder = new(ResourceManagerBinder)

type ResourceManagerBinder struct {
	router     internal.Router
	appContext *host.AppContext
}

func (b *ResourceManagerBinder) Init(context *structproto.StructProtoContext) error {
	return nil
}

func (b *ResourceManagerBinder) Bind(field structproto.FieldInfo, rv reflect.Value) error {
	if !rv.IsValid() {
		return fmt.Errorf("specifiec argument 'rv' is invalid")
	}

	// assign zero if rv is nil
	rvResource := reflectutil.AssignZero(rv)
	binder := &ResourceBinder{
		resourceType: rvResource.Type().Name(),
		components: map[string]reflect.Value{
			host.APP_CONFIG_FIELD:           b.appContext.Config(),
			host.APP_SERVICE_PROVIDER_FIELD: b.appContext.ServiceProvider(),
		},
	}
	err := b.preformBindResource(rvResource, binder)
	if err != nil {
		return err
	}

	// register RequestHandlers
	return b.registerRoute(field.Name(), rvResource)
}

func (b *ResourceManagerBinder) Deinit(context *structproto.StructProtoContext) error {
	return nil
}

func (b *ResourceManagerBinder) preformBindResource(target reflect.Value, binder *ResourceBinder) error {
	prototype, err := structproto.Prototypify(target,
		&structproto.StructProtoResolveOption{
			TagResolver: tagresolver.NoneTagResolver,
		})
	if err != nil {
		return err
	}

	return prototype.Bind(binder)
}

func (b *ResourceManagerBinder) registerRoute(url string, rvResource reflect.Value) error {
	// register RequestHandlers
	count := rvResource.Type().NumMethod()
	for i := 0; i < count; i++ {
		method := rvResource.Type().Method(i)

		rvMethod := rvResource.Method(method.Index)
		if isRequestHandler(rvMethod) {
			handler := asRequestHandler(rvMethod)
			if handler != nil {
				// TODO: validate path make comply RFC3986
				b.router.Add(strings.ToUpper(method.Name), url, handler)
			}
		}
	}
	return nil
}

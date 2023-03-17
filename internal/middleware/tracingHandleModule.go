package middleware

import (
	"context"
	"fmt"
	"reflect"

	. "github.com/Bofry/host-fasthttp/internal"
	"github.com/Bofry/trace"
)

const (
	tracerFieldName string = "Tracer"
)

var (
	typeOfTracer reflect.Type = reflect.TypeOf(new(trace.SeverityTracer))

	_ RequestHandleModule          = new(TracingHandleModule)
	_ RequestResourceProcessModule = new(TracingHandleModule)
)

type TracingHandleModule struct {
	tp     *trace.SeverityTracerProvider
	module RequestHandleModule
}

// ProcessRequestResource implements RequestResourceProcessModule
func (m *TracingHandleModule) ProcessRequestResource(rv reflect.Value) {
	rv = reflect.Indirect(rv)
	tracerValue := rv.FieldByName(tracerFieldName)

	if tracerValue.IsValid() && tracerValue.IsNil() && tracerValue.Type().AssignableTo(typeOfTracer) {
		name := rv.Type().Name()
		tracer := m.tp.Tracer(name)
		tracerValue.Set(reflect.ValueOf(tracer))
	}
}

// CanSetSuccessor implements RequestHandleModule
func (h *TracingHandleModule) CanSetSuccessor() bool {
	return true
}

// SetSuccessor implements RequestHandleModule
func (h *TracingHandleModule) SetSuccessor(successor RequestHandleModule) {
	h.module = successor
}

// ProcessRequest implements RequestHandleModule
func (h *TracingHandleModule) ProcessRequest(ctx *RequestCtx, recover *RecoverService) {
	if h.module != nil {
		fmt.Printf("Tracing Request: %s %s\n", string(ctx.Request.Header.Method()), string(ctx.Request.URI().Path()))

		recover.
			Defer(func(err interface{}) {
				if err != nil {
					defer func() {
						fmt.Printf("Tracing Error: %s %s [%v]\n", string(ctx.Request.Header.Method()), string(ctx.Request.URI().Path()), err)
					}()

					// NOTE: we should not handle error here, due to the underlying RequestHandler
					// will handle it.
				} else {
					fmt.Printf("Tracing Response: %s %s [%v]\n", string(ctx.Request.Header.Method()), string(ctx.Request.URI().Path()), ctx.Response.StatusCode())
				}
			}).
			Do(func() {
				h.module.ProcessRequest(ctx, recover)
			})
	}
}

// OnInitComplete implements RequestHandleModule
func (h *TracingHandleModule) OnInitComplete() {
	// ignored
}

// OnStart implements RequestHandleModule
func (h *TracingHandleModule) OnStart(ctx context.Context) error {
	// do nothing
	return nil
}

// OnStop implements RequestHandleModule
func (h *TracingHandleModule) OnStop(ctx context.Context) error {
	return h.tp.Shutdown(ctx)
}

package internal

import (
	"reflect"
	"strings"
	"sync"

	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

type TracerManager struct {
	TracerProvider    *trace.SeverityTracerProvider
	TextMapPropagator propagation.TextMapPropagator

	tracers map[reflect.Type]*trace.SeverityTracer

	mutex sync.Mutex
}

func NewTraceManager() *TracerManager {
	return &TracerManager{
		TracerProvider:    defaultTracerProvider,
		TextMapPropagator: defaultTextMapPropagator,
		tracers:           make(map[reflect.Type]*trace.SeverityTracer),
	}
}

func (m *TracerManager) GenerateManagedTracer(v interface{}) *trace.SeverityTracer {
	var rt reflect.Type
	if r, ok := v.(reflect.Type); ok {
		rt = r
	} else {
		rt = reflect.TypeOf(v)
	}

	for {
		if rt.Kind() != reflect.Ptr {
			break
		}
		rt = rt.Elem()
	}

	// find
	if tr, ok := m.tracers[rt]; ok {
		return tr
	}

	// create new
	return m.createManagedTracer(rt)
}

func (m *TracerManager) createManagedTracer(rt reflect.Type) *trace.SeverityTracer {
	for {
		if rt.Kind() != reflect.Ptr {
			break
		}
		rt = rt.Elem()
	}

	// create new
	m.mutex.Lock()
	defer m.mutex.Unlock()

	name := strings.Join([]string{
		rt.PkgPath(),
		rt.Name(),
	}, ".")

	tr := m.TracerProvider.Tracer(name)
	m.tracers[rt] = tr

	return tr
}

func (m *TracerManager) createTracer(name string) *trace.SeverityTracer {
	return m.TracerProvider.Tracer(name)
}

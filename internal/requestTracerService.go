package internal

import (
	"reflect"
	"sync"

	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

type RequestTracerService struct {
	TracerManager *TracerManager

	Enabled bool

	unhandledRequestTracer *trace.SeverityTracer

	tracers            map[string]*trace.SeverityTracer
	tracersInitializer sync.Once
}

func (s *RequestTracerService) Tracer(id string) *trace.SeverityTracer {
	if s.tracers != nil {
		if tr, ok := s.tracers[id]; ok {
			return tr
		}
	}
	return s.unhandledRequestTracer
}

func (s *RequestTracerService) TextMapPropagator() propagation.TextMapPropagator {
	return s.TracerManager.TextMapPropagator
}

func (s *RequestTracerService) init(requestManager interface{}) {
	if s.Enabled {
		trace.SetSpanExtractor(defaultSpanExtractor)

		s.makeTracerMap()
		s.buildTracer(requestManager)
	}
	s.makeUnhandledRequestTracer()
}

func (s *RequestTracerService) makeTracerMap() {
	s.tracersInitializer.Do(func() {
		s.tracers = make(map[string]*trace.SeverityTracer)
	})
}

func (s *RequestTracerService) buildTracer(requestManager interface{}) {
	var (
		rvManager reflect.Value = reflect.ValueOf(requestManager)
	)
	if rvManager.Kind() != reflect.Pointer || rvManager.IsNil() {
		return
	}

	rvManager = reflect.Indirect(rvManager)
	numOfHandles := rvManager.NumField()
	for i := 0; i < numOfHandles; i++ {
		rvRequest := rvManager.Field(i)
		if rvRequest.Kind() != reflect.Pointer || rvRequest.IsNil() {
			continue
		}

		rvRequest = reflect.Indirect(rvRequest)
		if rvRequest.Kind() == reflect.Struct {
			tracer := s.TracerManager.createManagedTracer(rvRequest.Type())

			info := rvManager.Type().Field(i)
			if _, ok := s.tracers[info.Name]; !ok {
				s.registerTracer(info.Name, tracer)
			}
		}
	}
}

func (s *RequestTracerService) registerTracer(id string, tracer *trace.SeverityTracer) {
	container := s.tracers

	if tracer != nil {
		if _, ok := container[id]; ok {
			FasthttpHostLogger.Fatalf("specified id '%s' already exists", id)
		}
		container[id] = tracer
	}
}

func (s *RequestTracerService) makeUnhandledRequestTracer() {
	s.unhandledRequestTracer = s.TracerManager.createTracer("")
}

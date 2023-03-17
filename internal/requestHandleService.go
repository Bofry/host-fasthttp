package internal

import (
	"context"
)

type RequestHandleService struct {
	modules []RequestHandleModule
}

func NewRequestHandleService() *RequestHandleService {
	return &RequestHandleService{}
}

func (s *RequestHandleService) Register(module RequestHandleModule) {
	size := len(s.modules)
	if size > 0 {
		last := s.modules[size-1]

		// ignore all new successor if the last RequestHandleModule cannot accept successor
		if !last.CanSetSuccessor() {
			return
		}

		last.SetSuccessor(module)
	}
	s.modules = append(s.modules, module)
}

func (s *RequestHandleService) ProcessRequest(ctx *RequestCtx, recover *RecoverService) {
	if handler := s.first(); handler != nil {
		handler.ProcessRequest(ctx, recover)
	}
}

func (s *RequestHandleService) first() RequestHandleModule {
	if len(s.modules) > 0 {
		return s.modules[0]
	}
	return nil
}

func (s *RequestHandleService) triggerStart(ctx context.Context) error {
	var err error

	for _, m := range s.modules {
		if err = m.OnStart(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *RequestHandleService) triggerStop(ctx context.Context) <-chan error {
	ch := make(chan error)

	go func() {
		defer close(ch)
		for _, m := range s.modules {
			FasthttpHostLogger.Printf("stopping %T", m)

			err := m.OnStop(ctx)
			if err != nil {
				ch <- &StopError{
					v:   m,
					err: err,
				}
			}
		}
	}()

	return ch
}

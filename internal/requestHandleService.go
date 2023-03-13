package internal

import (
	"context"
)

type RequestHandleService struct {
	handlers []RequestHandleModule
}

func (s *RequestHandleService) Register(successor RequestHandleModule) {
	size := len(s.handlers)
	if size > 0 {
		last := s.handlers[size-1]

		// ignore all new successor if the last RequestHandleModule cannot accept successor
		if !last.CanSetSuccessor() {
			return
		}

		last.SetSuccessor(successor)
	}
	s.handlers = append(s.handlers, successor)
}

func (s *RequestHandleService) ProcessRequest(ctx *RequestCtx, recover *RecoverService) {
	if handler := s.first(); handler != nil {
		handler.ProcessRequest(ctx, recover)
	}
}

func (s *RequestHandleService) first() RequestHandleModule {
	if len(s.handlers) > 0 {
		return s.handlers[0]
	}
	return nil
}

func (s *RequestHandleService) stop(ctx context.Context) <-chan error {
	ch := make(chan error)

	go func() {
		defer close(ch)
		for _, h := range s.handlers {
			err := h.OnStop(ctx)
			if err != nil {
				ch <- &StopError{
					v:   h,
					err: err,
				}
			}
		}
	}()

	return ch
}

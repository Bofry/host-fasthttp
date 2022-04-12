package internal

type RequestHandleService struct {
	handlers []RequestHandleModule
}

func (s *RequestHandleService) Register(successor RequestHandleModule) {
	size := len(s.handlers)
	if size > 0 {
		last := s.handlers[size-1]

		// ignore all new successor if the last RequestRouteResolveModule cannot accept successor
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

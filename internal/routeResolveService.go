package internal

var _ RouteResolver = new(RouteResolveService)

type RouteResolveService struct {
	resolvers []RouteResolveModule
}

func (s *RouteResolveService) Register(successor RouteResolveModule) {
	size := len(s.resolvers)
	if size > 0 {
		last := s.resolvers[size-1]

		// ignore all new successor if the last RequestRouteResolveModule cannot accept successor
		if !last.CanSetSuccessor() {
			return
		}

		last.SetSuccessor(successor)
	}
	s.resolvers = append(s.resolvers, successor)
}

func (s *RouteResolveService) ResolveHttpMethod(ctx *RequestCtx) string {
	if resolver := s.first(); resolver != nil {
		return resolver.ResolveHttpMethod(ctx)
	}
	return ""
}

func (s *RouteResolveService) ResolveHttpPath(ctx *RequestCtx) string {
	if resolver := s.first(); resolver != nil {
		return resolver.ResolveHttpPath(ctx)
	}
	return ""
}

func (s *RouteResolveService) first() RouteResolver {
	if len(s.resolvers) > 0 {
		return s.resolvers[0]
	}
	return nil
}

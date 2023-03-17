package internal

var _ RouteResolveModule = new(StdRouteResolveModule)

var (
	RouteResolveModuleInstance = new(StdRouteResolveModule)
)

type StdRouteResolveModule struct{}

func (r *StdRouteResolveModule) CanSetSuccessor() bool {
	return false
}

func (r *StdRouteResolveModule) SetSuccessor(successor RouteResolver) {
	panic("unsupported operation")
}

func (r *StdRouteResolveModule) ResolveHttpMethod(ctx *RequestCtx) string {
	return string(ctx.Method())
}

func (r *StdRouteResolveModule) ResolveHttpPath(ctx *RequestCtx) string {
	return string(ctx.Path())
}

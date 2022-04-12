package internal

var _ RouteResolveModule = new(RouteResolveModuleImpl)

var (
	RouteResolveModuleInstance = new(RouteResolveModuleImpl)
)

type RouteResolveModuleImpl struct{}

func (r *RouteResolveModuleImpl) CanSetSuccessor() bool {
	return false
}

func (r *RouteResolveModuleImpl) SetSuccessor(successor RouteResolver) {
	panic("unsupported operation")
}

func (r *RouteResolveModuleImpl) ResolveHttpMethod(ctx *RequestCtx) string {
	return string(ctx.Method())
}

func (r *RouteResolveModuleImpl) ResolveHttpPath(ctx *RequestCtx) string {
	return string(ctx.Path())
}

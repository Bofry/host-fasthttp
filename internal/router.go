package internal

type Router map[RoutePath]RouteComponent

func (r Router) Add(method string, path string, handler RequestHandler, handlerComponentID string) {
	key := RoutePath{
		Method: method,
		Path:   path,
	}
	r[key] = RouteComponent{
		RequestHandler:     handler,
		HandlerComponentID: handlerComponentID,
	}
}

func (r Router) Remove(method string, path string) {
	key := RoutePath{
		Method: method,
		Path:   path,
	}
	delete(r, key)
}

func (r Router) Get(method string, path string) RequestHandler {
	if r == nil {
		return nil
	}

	key := RoutePath{
		Method: method,
		Path:   path,
	}
	if v, ok := r[key]; ok {
		return v.RequestHandler
	}
	return nil
}

func (r Router) Has(path RoutePath) bool {
	if r == nil {
		return false
	}

	if _, ok := r[path]; ok {
		return true
	}
	return false
}

func (r Router) FindHandlerComponentID(method string, path string) string {
	if r == nil {
		return ""
	}

	key := RoutePath{
		Method: method,
		Path:   path,
	}
	if v, ok := r[key]; ok {
		return v.HandlerComponentID
	}
	return ""
}

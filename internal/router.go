package internal

type Router map[RoutePath]RequestHandler

func (r Router) Add(method string, path string, handler RequestHandler) {
	key := RoutePath{
		Method: method,
		Path:   path,
	}
	r[key] = handler
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
		return v
	}
	return nil
}

package responseutil

var _ ResponseState = new(ResponseStateImpl)

type ResponseStateImpl struct {
	flag       ResponseFlag
	statusCode int
}

func (r *ResponseStateImpl) Flag() ResponseFlag {
	return r.flag
}

func (r *ResponseStateImpl) StatusCode() int {
	return r.statusCode
}

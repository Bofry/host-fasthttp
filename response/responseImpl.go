package response

var _ Response = new(responseImpl)

type responseImpl struct {
	flag       ResponseFlag
	statusCode int
}

func (r *responseImpl) Flag() ResponseFlag {
	return r.flag
}

func (r *responseImpl) StatusCode() int {
	return r.statusCode
}

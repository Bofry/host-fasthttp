package internal

var _ RequestHandleModule = new(RequestHandleModuleImpl)

type RequestHandleModuleImpl struct {
	worker *FasthttpRequestWorker
}

func NewRequestHandleModuleImpl(worker *FasthttpRequestWorker) *RequestHandleModuleImpl {
	return &RequestHandleModuleImpl{
		worker: worker,
	}
}

func (r *RequestHandleModuleImpl) CanSetSuccessor() bool {
	return false
}

func (r *RequestHandleModuleImpl) SetSuccessor(successor RequestHandleModule) {
	panic("unsupported operation")
}

func (r *RequestHandleModuleImpl) ProcessRequest(ctx *RequestCtx, recover *RecoverService) {
	r.worker.internalProcessRequest(ctx, recover)
}

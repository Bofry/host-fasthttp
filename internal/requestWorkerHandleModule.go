package internal

import "context"

var _ RequestHandleModule = new(RequestWorkerHandleModule)

type RequestWorkerHandleModule struct {
	worker *RequestWorker
}

func NewRequestWorkerHandleModule(worker *RequestWorker) *RequestWorkerHandleModule {
	return &RequestWorkerHandleModule{
		worker: worker,
	}
}

// CanSetSuccessor implements RequestHandleModule
func (r *RequestWorkerHandleModule) CanSetSuccessor() bool {
	return false
}

// SetSuccessor implements RequestHandleModule
func (r *RequestWorkerHandleModule) SetSuccessor(successor RequestHandleModule) {
	panic("unsupported operation")
}

// ProcessRequest implements RequestHandleModule
func (r *RequestWorkerHandleModule) ProcessRequest(ctx *RequestCtx, recover *RecoverService) {
	r.worker.internalProcessRequest(ctx, recover)
}

// OnInitComplete implements RequestHandleModule
func (*RequestWorkerHandleModule) OnInitComplete() {
	// ignored
}

// OnStop implements RequestHandleModule
func (*RequestWorkerHandleModule) OnStop(ctx context.Context) error {
	// do nothing
	return nil
}

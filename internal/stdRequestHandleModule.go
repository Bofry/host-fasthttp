package internal

import "context"

var _ RequestHandleModule = new(StdRequestHandleModule)

type StdRequestHandleModule struct {
	worker *RequestWorker
}

func NewRequestWorkerHandleModule(worker *RequestWorker) *StdRequestHandleModule {
	return &StdRequestHandleModule{
		worker: worker,
	}
}

// CanSetSuccessor implements RequestHandleModule
func (r *StdRequestHandleModule) CanSetSuccessor() bool {
	return false
}

// SetSuccessor implements RequestHandleModule
func (r *StdRequestHandleModule) SetSuccessor(successor RequestHandleModule) {
	panic("unsupported operation")
}

// ProcessRequest implements RequestHandleModule
func (r *StdRequestHandleModule) ProcessRequest(ctx *RequestCtx, state RequestState, recover *RecoverService) {
	r.worker.internalProcessRequest(ctx, state, recover)
}

// OnInitComplete implements RequestHandleModule
func (*StdRequestHandleModule) OnInitComplete() {
	// ignored
}

// OnStart implements RequestHandleModule
func (*StdRequestHandleModule) OnStart(ctx context.Context) error {
	// do nothing
	return nil
}

// OnStop implements RequestHandleModule
func (*StdRequestHandleModule) OnStop(ctx context.Context) error {
	// do nothing
	return nil
}

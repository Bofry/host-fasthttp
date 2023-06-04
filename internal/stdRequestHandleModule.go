package internal

import "context"

var _ RequestHandleModule = new(StdRequestHandleModule)

type StdRequestHandleModule struct {
	worker *RequestWorker
}

func NewStdRequestHandleModule(worker *RequestWorker) *StdRequestHandleModule {
	return &StdRequestHandleModule{
		worker: worker,
	}
}

// CanSetSuccessor implements RequestHandleModule
func (*StdRequestHandleModule) CanSetSuccessor() bool {
	return false
}

// SetSuccessor implements RequestHandleModule
func (*StdRequestHandleModule) SetSuccessor(successor RequestHandleModule) {
	panic("unsupported operation")
}

// ProcessRequest implements RequestHandleModule
func (m *StdRequestHandleModule) ProcessRequest(ctx *RequestCtx, state RequestState, recover *Recover) {
	m.worker.internalProcessRequest(ctx, state, recover)
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

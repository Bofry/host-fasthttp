package internal

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/Bofry/host"
	"github.com/Bofry/trace"
	"github.com/valyala/fasthttp"
)

var _ host.Host = new(FasthttpHost)

type FasthttpHost struct {
	Server         *Server
	ListenAddress  string
	EnableCompress bool
	Version        string

	TracerProvider *trace.SeverityTracerProvider
	Logger         *log.Logger

	requestWorker  *RequestWorker
	requestManager interface{}

	requestResourceProcessService *RequestResourceProcessService

	wg          sync.WaitGroup
	mutex       sync.Mutex
	initialized bool
	running     bool
	disposed    bool
}

func (h *FasthttpHost) Start(ctx context.Context) {
	if h.disposed {
		FasthttpHostLogger.Panic("the FasthttpHost has been disposed")
	}
	if !h.initialized {
		FasthttpHostLogger.Panic("the FasthttpHost havn't be initialized yet")
	}
	if h.running {
		return
	}

	var err error
	h.mutex.Lock()
	defer func() {
		if err != nil {
			h.running = false
			h.disposed = true
		}
		h.mutex.Unlock()
	}()
	h.running = true

	s := h.Server

	h.requestWorker.start(ctx)

	go func() {
		FasthttpHostLogger.Printf("%% Notice: %s listening on address %s\n", h.Server.Name, h.ListenAddress)
		if err = s.ListenAndServe(h.ListenAddress); err != nil {
			FasthttpHostLogger.Fatalf("%% Error: error in ListenAndServe: %v\n", err)
		}
	}()
}

func (h *FasthttpHost) Stop(ctx context.Context) error {
	if h.disposed {
		return nil
	}
	if !h.running {
		return nil
	}

	var (
		server = h.Server
	)

	h.mutex.Lock()
	defer func() {
		h.running = false
		h.disposed = true
		h.mutex.Unlock()

		h.requestWorker.stop(ctx)
	}()

	err := server.Shutdown()
	h.wg.Wait()
	return err
}

func (h *FasthttpHost) preInit() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.requestWorker = NewRequestWorker()
	h.requestResourceProcessService = NewRequestHandlerReprocessService()
}

func (h *FasthttpHost) init() {
	if h.initialized {
		return
	}

	h.mutex.Lock()
	defer func() {
		h.initialized = true
		h.mutex.Unlock()
	}()

	if h.Server == nil {
		h.Server = &Server{}
	}

	h.requestWorker.init()
	h.processRequestResource()
	h.configRequestHandler()
	h.configCompress()
	h.configListenAddress()
}

func (h *FasthttpHost) processRequestResource() {
	if h.requestManager != nil {
		h.requestResourceProcessService.Process(h.requestManager)
	}
}

func (h *FasthttpHost) configRequestHandler() {
	s := h.Server
	var requestHandler RequestHandler

	if s.Handler != nil {
		requestHandler = s.Handler
	} else if h.requestWorker != nil {
		requestHandler = h.requestWorker.ProcessRequest
	}

	s.Handler = func(ctx *RequestCtx) {
		h.wg.Add(1)
		defer h.wg.Done()

		requestHandler(ctx)
	}
}

func (h *FasthttpHost) configCompress() {
	s := h.Server
	if h.EnableCompress {
		s.Handler = fasthttp.CompressHandler(s.Handler)
	}
}

func (h *FasthttpHost) configListenAddress() {
	host, port, err := splitHostPort(h.ListenAddress)
	if err != nil {
		panic(err)
	}

	if len(port) == 0 {
		port = DEFAULT_HTTP_PORT
	}
	h.ListenAddress = net.JoinHostPort(host, port)
}

func (h *FasthttpHost) onInitComplete() {

}

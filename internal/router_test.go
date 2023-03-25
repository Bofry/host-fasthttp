package internal

import "testing"

func TestRouter(t *testing.T) {
	router := make(Router)

	var breadcrumb bool = false

	postEchoRequestHandler := RequestHandler(func(ctx *RequestCtx) {
		defer func() {
			err := recover()
			if err != nil {
				panic(err)
			}
		}()

		breadcrumb = true
	})

	router.Add("POST", "/Echo", postEchoRequestHandler, "")
	requestHandler := router.Get("POST", "/Echo")

	if requestHandler == nil {
		t.Errorf("assert RequestHandler should not nil")
	}
	if requestHandler(nil); breadcrumb == false {
		t.Errorf("assert RequestHandler and postEchoRequestHandler are not same")
	}
}

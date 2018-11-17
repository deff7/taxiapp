package main

import (
	"log"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

type serviceHandler struct {
	manager *BidManager
}

func newServiceHandler(manager *BidManager) *serviceHandler {
	return &serviceHandler{manager: manager}
}

func (h *serviceHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/get":
		id := h.manager.GetRandomID()
		ctx.Write([]byte(id))
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func (h *serviceHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id := h.manager.GetRandomID()
	w.Write([]byte(id))
}

func main() {
	done := make(chan struct{})
	defer close(done)

	manager := NewBidManager(done)
	go func() {
		t := time.NewTicker(200 * time.Millisecond)
		for {
			select {
			case <-done:
				return
			case <-t.C:
				manager.Update()
			}
		}
	}()

	/*
		ln, err := reuseport.Listen("tcp4", "localhost:8000")
		if err != nil {
			log.Fatalf("error in reuseport listener: %s", err)
		}
		log.Fatal(fasthttp.Serve(ln, newServiceHandler(manager).HandleFastHTTP))
	*/

	http.Handle("/get", newServiceHandler(manager))
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

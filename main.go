package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type serviceHandler struct {
	manager *BidManager
}

func newServiceHandler(manager *BidManager) *serviceHandler {
	return &serviceHandler{manager: manager}
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

	http.Handle("/get", newServiceHandler(manager))
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

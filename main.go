package main

import (
	"log"
	"net/http"
	"time"
)

type serviceHandler struct {
	manager *BidManager
}

func newServiceHandler(manager *BidManager) *serviceHandler {
	return &serviceHandler{manager: manager}
}

func (h *serviceHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	bid := h.manager.GetRandom()
	w.Write([]byte(bid.ID))
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

	log.Fatal(http.ListenAndServe(":8000", nil))
}

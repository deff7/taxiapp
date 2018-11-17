package main

import (
	"sync"
	"testing"
)

func TestBidManager(t *testing.T) {
	var wg sync.WaitGroup
	done := make(chan struct{})
	defer close(done)
	m := NewBidManager(done)

	t.Run("concurrent updates causes no data races", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				m.Update()

				m.Lock()
				if len(m.bids) != bidsPoolSize {
					t.Fatalf("len(m.bids) must be %d, got %d", bidsPoolSize, len(m.bids))
				}
				m.Unlock()
			}()
		}
		wg.Wait()
	})
}

func BenchmarkGetRandom(b *testing.B) {
	done := make(chan struct{})
	defer close(done)

	manager := NewBidManager(done)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		manager.GetRandomID()
	}
}

func BenchmarkUpdate(b *testing.B) {
	done := make(chan struct{})
	defer close(done)

	manager := NewBidManager(done)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		manager.Update()
	}
}

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
				if len(m.bids) != bidsPoolSize {
					t.Fatalf("len(m.bids) must be %d, got %d", bidsPoolSize, len(m.bids))
				}
			}()
		}
		wg.Wait()
	})

	t.Run("concurrent fetches causes no data races", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				bid := m.GetRandom()
				if bid.TimesShowed == 0 {
					t.Fatal("showed counter not updated")
				}
			}()
		}
		wg.Wait()
	})
}

package main

import "testing"

func TestIDGenerator(t *testing.T) {

	t.Run("done signal allow to drain channel", func(t *testing.T) {
		done := make(chan struct{})
		g := NewIDGenerator(done)
		close(done)
		for range g {
		}
		_, ok := <-g
		if ok {
			t.Error("channel must be closed")
		}
	})

	t.Run("second ids array pass returns shuffled ids", func(t *testing.T) {
		done := make(chan struct{})
		g := NewIDGenerator(done)

		size := 26 * 26
		ids := make([][]BidID, 2)

		for i := 0; i < 2; i++ {
			ids[i] = make([]BidID, size)
			for j := 0; j < size; j++ {
				ids[i][j] = <-g
			}
		}

		diff := 0
		for i := 0; i < size; i++ {
			if ids[0][i] != ids[1][i] {
				diff++
			}
		}

		if diff == 0 {
			t.Error("arrays must be different")
		}
	})
}

package main

import "math/rand"

type IDGenerator struct {
}

func NewIDGenerator(done <-chan struct{}) <-chan BidID {
	var (
		ids        []BidID
		indicies   []int
		currentIdx int
		idChan     chan BidID

		r = rand.New(rand.NewSource(0))
	)

	size := int('z' - 'a')
	ids = make([]BidID, size*size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			ids[i*size+j] = BidID(string('a'+i) + string('a'+j))
		}
	}

	idChan = make(chan BidID, size*size)

	go func() {
		for {
			if len(indicies) == 0 || currentIdx >= len(ids) {
				indicies = r.Perm(len(ids))
				currentIdx = 0
			}
			id := ids[indicies[currentIdx]]

			select {
			case <-done:
				close(idChan)
				return
			case idChan <- id:
				currentIdx++
			}
		}
	}()
	return idChan
}

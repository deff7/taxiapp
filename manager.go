package main

import (
	"math/rand"
	"sync"
)

const bidsPoolSize = 50

type BidID string

type Bid struct {
	ID          BidID
	TimesShowed int
}

type BidManager struct {
	sync.Mutex
	bids   map[BidID]*Bid
	idChan <-chan BidID
}

func NewBidManager(done <-chan struct{}) *BidManager {
	m := &BidManager{}
	m.idChan = NewIDGenerator(done)
	m.bids = make(map[BidID]*Bid, bidsPoolSize)
	for i := 0; i < bidsPoolSize; i++ {
		m.AddBid()
	}
	return m
}

func (m *BidManager) AddBid() {
	bid := &Bid{ID: <-m.idChan}
	m.bids[bid.ID] = bid
}

func (m *BidManager) Update() {
	m.Lock()
	id := m.randomBidID()
	delete(m.bids, id)
	m.AddBid()
	m.Unlock()
}

func (m *BidManager) GetRandom() *Bid {
	var bid *Bid
	m.Lock()
	id := m.randomBidID()
	bid = m.bids[id]
	bid.TimesShowed++
	m.Unlock()
	return bid
}

func (m *BidManager) randomBidID() BidID {
	var id BidID
	i := rand.Intn(len(m.bids))
	for k := range m.bids {
		if i == 0 {
			id = k
			break
		}
		i--
	}
	return id
}

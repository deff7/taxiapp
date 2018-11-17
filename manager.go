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
	sync.RWMutex
	bids   map[*Bid]struct{}
	idChan <-chan BidID
}

func NewBidManager(done <-chan struct{}) *BidManager {
	m := &BidManager{}
	m.idChan = NewIDGenerator(done)
	m.bids = make(map[*Bid]struct{}, bidsPoolSize)
	for i := 0; i < bidsPoolSize; i++ {
		m.AddBid()
	}
	return m
}

func (m *BidManager) AddBid() {
	bid := &Bid{ID: <-m.idChan}
	m.bids[bid] = struct{}{}
}

func (m *BidManager) Update() {
	m.Lock()
	bid := m.randomBid()
	bid.ID = <-m.idChan
	bid.TimesShowed = 0
	m.Unlock()
}

func (m *BidManager) GetRandomID() BidID {
	var id BidID
	m.RLock()
	bid := m.randomBid()
	bid.TimesShowed++
	id = bid.ID
	m.RUnlock()
	return id
}

func (m *BidManager) randomBid() *Bid {
	var bid *Bid
	i := rand.Intn(len(m.bids))
	for v := range m.bids {
		if i == 0 {
			bid = v
			break
		}
		i--
	}
	return bid
}

package db

import "nft-indexer-api/indexer"

type EventStorage struct {
	Collections []indexer.CollectionCreatedEvent
	Mints []indexer.TokenMintedEvent
	quit chan struct{}
}

func NewEventStorage(idx *indexer.Indexer) (*EventStorage) {
	stor := EventStorage{
		Collections: []indexer.CollectionCreatedEvent{},
		Mints: []indexer.TokenMintedEvent{},
		quit: make(chan struct{}),
	}

	go func(){
		for {
			select {
			case events := <- idx.CollectionCreatedEvents:
				stor.Collections = append(stor.Collections, events...)
			case events := <- idx.TokenMintedEvents:
				stor.Mints = append(stor.Mints, events...)
			case <-stor.quit:
				break
			}
		}
	}()

	return &stor
}

func (stor *EventStorage) Close() {
	stor.quit <- struct{}{}
	close(stor.quit)
}
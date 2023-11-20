package indexer

import (
	"context"
	"log"
	"nft-indexer-api/rpc"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-co-op/gocron"
)

type Indexer struct {
	collectionMapper *CollectionCreatedEventMapper
	tokenMapper *TokenMintedEventMapper
	config WorkerConfig
	address string
	lastIndexedBlock uint64
	CollectionCreatedEvents chan []CollectionCreatedEvent
	TokenMintedEvents chan []TokenMintedEvent
	scheduler *gocron.Scheduler
	finalisationBlocksCount uint8
}

func NewIndexer(
	rpc *rpc.Rpc, 
	address string, 
	creationBlock uint64, 
	parralelRequests uint8, 
	blocksPerRequest uint8,
	finalisationBlocksCount uint8,
) (*Indexer, error) {
	abi, err := NewParsedAbi();
	if err != nil {
		return nil, err
	}

	collectionMapper := NewCollectionCreatedEventMapper(abi)
	tokenMapper := NewTokenMintedEventMapper(abi)

	config := WorkerConfig{
		MaxParralelRequests: parralelRequests / 2,
		BlocksPerRequest: blocksPerRequest,
		rpc: rpc,
	}

	scheduler := gocron.NewScheduler(time.UTC)

	return &Indexer{
		collectionMapper: collectionMapper,
		tokenMapper: tokenMapper,
		config: config,
		address: address,
		lastIndexedBlock: creationBlock,
		CollectionCreatedEvents: make(chan []CollectionCreatedEvent),
		TokenMintedEvents: make(chan []TokenMintedEvent),
		scheduler: scheduler,
		finalisationBlocksCount: finalisationBlocksCount,
	}, nil
}

func (idx *Indexer) Index(ctx context.Context) {
	log.Default().Print("Indexer job started")
	collectionsResult := make(chan []types.Log)
	tokensResult := make(chan []types.Log)

	defer close(collectionsResult)
	defer close(tokensResult)

	ctx, cancel := context.WithCancel(ctx)

	height, err := idx.config.rpc.GetCurrentHeight(ctx)

	if err != nil {
		log.Default().Print(err.Error())
		cancel()
		return
	}

	startBlock := idx.lastIndexedBlock;
	finishBlocks := height - uint64(idx.finalisationBlocksCount)

	if (finishBlocks <= startBlock) {
		return
	}

	go idx.config.PerformWork(ctx, WorkerTask{
		Topic: COLLECTION_CREATED_TOPIC,
		Address: idx.address,
		BlockStartNumber: startBlock,
		BlockFinishNumber: finishBlocks,
	}, collectionsResult)

	go idx.config.PerformWork(ctx, WorkerTask{
		Topic: TOKEN_MINTED_TOPIC,
		Address: idx.address,
		BlockStartNumber: startBlock,
		BlockFinishNumber: finishBlocks,
	}, tokensResult)

	timeout := time.After(time.Minute * 5)
	
	mappedCollections := []CollectionCreatedEvent{}
	mappedMints := []TokenMintedEvent{}

	for i := 0; i < 2; i++ {
		select {
		case logs := <-collectionsResult:
			log.Default().Printf("Collections job fetched %d logs", len(logs))
			
			for _, l := range logs {
				event, err := idx.collectionMapper.MapToCollectionCreatedEvent(l)
				if err != nil {
					log.Default().Print(err.Error())
					cancel()
					continue
				}
				mappedCollections = append(mappedCollections, *event)
			}
		case logs := <-tokensResult:
			log.Default().Printf("Mints job fetched %d logs", len(logs))
			for _, l := range logs {
				event, err := idx.tokenMapper.MapToTokenMintedEven(l)
				if err != nil {
					log.Default().Print(err.Error())
					cancel()
					continue
				}
				mappedMints = append(mappedMints, *event)
			}
		case <- timeout:
			log.Default().Print("Indexer timeout exceeded")
			return
		}
	}

	idx.CollectionCreatedEvents <- mappedCollections
	idx.TokenMintedEvents <- mappedMints

	idx.lastIndexedBlock = finishBlocks
}

func (idx *Indexer) Start() error {
	log.Default().Print("Indexer start requested")
	idx.scheduler.SingletonModeAll()
	_, err := idx.scheduler.Every(1).Minute().Do(func() {
		idx.Index(context.Background())
	})
	if err != nil {
		return err
	}
	idx.scheduler.StartAsync()
	log.Default().Print("Indexer started")
	return nil
}

func (idx *Indexer) Stop() {
	idx.scheduler.Stop()
	close(idx.CollectionCreatedEvents)
	close(idx.TokenMintedEvents)
	log.Default().Print("Indexer stopped")
}
package indexer

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	COLLECTION_CREATED_TOPIC = "0x3454b57f2dca4f5a54e8358d096ac9d1a0d2dab98991ddb89ff9ea1746260617"
)

type mappedCollectionCreatedEvent struct {
	Collection common.Address
	Name string
	Symbol string
}

type CollectionCreatedEvent struct {
	Collection string `json:"collection"`
	Name string `json:"name"`
	Symbol string `json:"symbol"`
}

type CollectionCreatedEventMapper struct {
	parsedABI *abi.ABI
}

func NewCollectionCreatedEventMapper(abi *abi.ABI) *CollectionCreatedEventMapper {

	return &CollectionCreatedEventMapper{
		parsedABI: abi,
	}
}

func (mapper *CollectionCreatedEventMapper) MapToCollectionCreatedEvent(log types.Log) (*CollectionCreatedEvent, error) {
	var event mappedCollectionCreatedEvent;
	err := mapper.parsedABI.UnpackIntoInterface(&event, "CollectionCreated", log.Data)

	if err != nil {
		return nil, err
	}

	return &CollectionCreatedEvent{
		Collection: event.Collection.Hex(),
		Name: event.Name,
		Symbol: event.Symbol,
	}, nil
}
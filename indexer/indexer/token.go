package indexer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	TOKEN_MINTED_TOPIC = "0xc9fee7cd4889f66f10ff8117316524260a5242e88e25e0656dfb3f4196a21917"
)

type mappedTokenMintedEvent struct {
	Collection common.Address
	Recipient common.Address
	TokenId *big.Int
	TokenURI string
}

type TokenMintedEvent struct {
	Collection string `json:"collections"`
	Recipient string `json:"recepient"`
	TokenId *big.Int `json:"tokenId"`
	TokenURI string `json:"tokenUri"`
}

type TokenMintedEventMapper struct {
	parsedABI *abi.ABI
}

func NewTokenMintedEventMapper(abi *abi.ABI) *TokenMintedEventMapper {

	return &TokenMintedEventMapper{
		parsedABI: abi,
	}
}

func (mapper *TokenMintedEventMapper) MapToTokenMintedEven(log types.Log) (*TokenMintedEvent, error) {
	var event mappedTokenMintedEvent;
	err := mapper.parsedABI.UnpackIntoInterface(&event, "TokenMinted", log.Data)

	if err != nil {
		return nil, err
	}

	return &TokenMintedEvent{
		Collection: event.Collection.Hex(),
		Recipient: event.Recipient.Hex(),
		TokenId: event.TokenId,
		TokenURI: event.TokenURI,
	}, nil
}
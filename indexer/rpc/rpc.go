package rpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Rpc struct {
	client *ethclient.Client
}

func NewEthRpc(ctx context.Context, url string) (*Rpc, error) {
	client, err := ethclient.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}

	rpc := Rpc{
		client: client,
	}

	return &rpc, nil
}

func (rpc *Rpc) GetEvents(ctx context.Context, from string, topic string, blockStart uint64, blockFinish uint64) ([]types.Log, error) {
	logs, err := rpc.client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(blockStart)),
		ToBlock: big.NewInt(int64(blockFinish)),

		Addresses: []common.Address{common.HexToAddress(from)},
		Topics: [][]common.Hash{{common.HexToHash(topic)},},
	})

	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (rpc *Rpc) GetCurrentHeight(ctx context.Context) (uint64, error) {
	height, err := rpc.client.BlockNumber(ctx)

	if err != nil {
		return 0, err
	}

	return height, nil
}

func (rpc *Rpc) Close() {
	rpc.client.Close()
}
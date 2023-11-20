package indexer

import (
	"bufio"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func newParsedAbi() (*abi.ABI, error) {
	abiFile, err := os.Open("./abi/NFTFactory.abi.json")

	if err != nil {
		return nil, err
	}

	parsedABI, err := abi.JSON(bufio.NewReader(abiFile))
	if err != nil {
		return nil, err
	}

	return &parsedABI, nil
}
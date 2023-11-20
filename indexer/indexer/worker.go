package indexer

import (
	"context"
	"log"
	"nft-indexer-api/rpc"

	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/sync/errgroup"
)

type worker struct {
	BlocksPerRequest uint8
	MaxParralelRequests uint8
	rpc *rpc.Rpc
}

type workerTask struct {
	Topic string
	Address string
	BlockStartNumber uint64
	BlockFinishNumber uint64
}

func (config *worker) performWork(ctx context.Context, task workerTask, workerResult chan<- []types.Log) error {
	blocksToFetch := task.BlockFinishNumber - task.BlockStartNumber
	tasksToPerform := blocksToFetch / uint64(config.BlocksPerRequest) + 1

	log.Default().Printf("Requested event fetch for range (%d-%d), tasks to perform %d", task.BlockStartNumber, task.BlockFinishNumber, tasksToPerform)

	throttle := make(chan struct{}, config.MaxParralelRequests);
	results := make(chan []types.Log, tasksToPerform)

	defer close(throttle)
	defer close(results)

	g, ctx := errgroup.WithContext(ctx);

	for i := uint64(0); i < tasksToPerform; i++ {
		taskBlockStart := min(
			task.BlockStartNumber + i * uint64(config.BlocksPerRequest), 
			task.BlockFinishNumber,
		)
		taskBlockFinish := min(
			taskBlockStart + uint64(config.BlocksPerRequest), 
			task.BlockFinishNumber,
		)

		if (taskBlockStart >= taskBlockFinish) {
			results <- []types.Log{}
			continue
		}

		g.Go(func () error {
			throttle <- struct{}{}
			defer func() {
				<-throttle
			}()

			logs, err := config.rpc.GetEvents(
				ctx,
				task.Address,
				task.Topic,
				taskBlockStart,
				taskBlockFinish,
			)

			if err != nil {
				return err
			}

			results <- logs

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	collectedResult := []types.Log{}

	for i := 0; i < int(tasksToPerform); i++ {
		collectedResult = append(collectedResult, <-results...)
	}

	workerResult <- collectedResult

	return nil
}
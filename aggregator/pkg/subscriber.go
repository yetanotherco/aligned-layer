package pkg

import "github.com/yetanotherco/aligned_layer/core/chainio"

func (agg *Aggregator) SubscribeToNewTasks() *chainio.ErrorPair {
	errorPairPtr := agg.subscribeToNewTasks()
	if errorPairPtr != nil {
		return errorPairPtr
	}

	for {
		select {
		case err := <-agg.taskSubscriber:
			agg.AggregatorConfig.BaseConfig.Logger.Info("Failed to subscribe to new tasks", "err", err)
			errorPairPtr = agg.subscribeToNewTasks()
			if errorPairPtr != nil {
				return errorPairPtr
			}
		case newBatch := <-agg.NewBatchChan:
			agg.AggregatorConfig.BaseConfig.Logger.Info("Adding new task")
			agg.AddNewTask(newBatch.BatchMerkleRoot, newBatch.SenderAddress, newBatch.TaskCreatedBlock)
		}
	}
}

func (agg *Aggregator) subscribeToNewTasks() *chainio.ErrorPair {
	errorPairPtr := agg.avsSubscriber.SubscribeToNewTasksV3(agg.NewBatchChan, agg.taskSubscriber)

	if errorPairPtr != nil {
		agg.AggregatorConfig.BaseConfig.Logger.Info("Failed to create task subscriber", "err", errorPairPtr)
	}

	return errorPairPtr
}

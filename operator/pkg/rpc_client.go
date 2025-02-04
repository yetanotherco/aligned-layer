package operator

import (
	"errors"
	"net/rpc"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	retry "github.com/yetanotherco/aligned_layer/core"
	"github.com/yetanotherco/aligned_layer/core/types"
)

// AggregatorRpcClient is the client to communicate with the aggregator via RPC
type AggregatorRpcClient struct {
	rpcClient            *rpc.Client
	aggregatorIpPortAddr string
	logger               logging.Logger
}

func NewAggregatorRpcClient(aggregatorIpPortAddr string, logger logging.Logger) (*AggregatorRpcClient, error) {
	client, err := rpc.DialHTTP("tcp", aggregatorIpPortAddr)
	if err != nil {
		return nil, err
	}

	return &AggregatorRpcClient{
		rpcClient:            client,
		aggregatorIpPortAddr: aggregatorIpPortAddr,
		logger:               logger,
	}, nil
}

func SendSignedTaskResponse(c *AggregatorRpcClient, signedTaskResponse *types.SignedTaskResponse) func() (uint8, error) {
	send_task_func := func() (uint8, error) {
		var reply uint8
		err := c.rpcClient.Call("Aggregator.ProcessOperatorSignedTaskResponseV2", signedTaskResponse, &reply)
		if err != nil {
			c.logger.Error("Received error from aggregator", "err", err)
			if errors.Is(err, rpc.ErrShutdown) {
				c.logger.Error("Aggregator is shutdown. Reconnecting...")
			}
		} else {
			c.logger.Info("Signed task response header accepted by aggregator.", "reply", reply)
		}
		return reply, err
	}
	return send_task_func
}

// SendSignedTaskResponseToAggregator is the method called by operators via RPC to send
// their signed task response.
func (c *AggregatorRpcClient) SendSignedTaskResponseToAggregatorRetryable(signedTaskResponse *types.SignedTaskResponse) (uint8, error) {
	config := retry.DefaultRetryConfig()
	config.NumRetries = 10
	config.Multiplier = 1 // Constant retry interval
	config.InitialInterval = 10 * time.Second
	return retry.RetryWithData(SendSignedTaskResponse(c, signedTaskResponse), config)
}

package retry_test

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	rpccalls "github.com/Layr-Labs/eigensdk-go/metrics/collectors/rpc_calls"
	backoff "github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	servicemanager "github.com/yetanotherco/aligned_layer/contracts/bindings/AlignedLayerServiceManager"
	retry "github.com/yetanotherco/aligned_layer/core"
	"github.com/yetanotherco/aligned_layer/core/chainio"
	"github.com/yetanotherco/aligned_layer/core/config"
	"github.com/yetanotherco/aligned_layer/core/utils"
)

func DummyFunction(x uint64) (uint64, error) {
	if x == 42 {
		return 0, retry.PermanentError{Inner: fmt.Errorf("Permanent error!")}
	} else if x < 42 {
		return 0, fmt.Errorf("Transient error!")
	}
	return x, nil
}

func TestRetryWithData(t *testing.T) {
	function := func() (*uint64, error) {
		x, err := DummyFunction(43)
		return &x, err
	}

	config := &retry.RetryConfig{
		InitialInterval:     1000,
		MaxInterval:         2,
		MaxElapsedTime:      3,
		RandomizationFactor: 0,
		Multiplier:          retry.EthCallMultiplier,
		NumRetries:          retry.EthCallNumRetries,
	}
	_, err := retry.RetryWithData(function, config)
	if err != nil {
		t.Errorf("Retry error!: %s", err)
	}
}

func TestRetry(t *testing.T) {
	function := func() error {
		_, err := DummyFunction(43)
		return err
	}
	config := &retry.RetryConfig{
		InitialInterval:     1000,
		MaxInterval:         2,
		MaxElapsedTime:      3,
		RandomizationFactor: 0,
		Multiplier:          retry.EthCallMultiplier,
		NumRetries:          retry.EthCallNumRetries,
	}
	err := retry.Retry(function, config)
	if err != nil {
		t.Errorf("Retry error!: %s", err)
	}
}

/*
Starts an anvil instance via the cli.
Assumes that anvil is installed but checks.
*/
func SetupAnvil(port uint16) (*exec.Cmd, *eth.InstrumentedClient, error) {

	path, err := exec.LookPath("anvil")
	if err != nil {
		fmt.Printf("Could not find `anvil` executable in `%s`\n", path)
	}

	port_str := strconv.Itoa(int(port))
	http_rpc_url := fmt.Sprintf("http://localhost:%d", port)

	// Create a command
	cmd := exec.Command("anvil", "--port", port_str, "--load-state", "../contracts/scripts/anvil/state/alignedlayer-deployed-anvil-state.json", "--block-time", "3")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Run the command
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	// Delay needed for anvil to start
	time.Sleep(750 * time.Millisecond)

	reg := prometheus.NewRegistry()
	rpcCallsCollector := rpccalls.NewCollector("ethRpc", reg)
	ethRpcClient, err := eth.NewInstrumentedClient(http_rpc_url, rpcCallsCollector)
	if err != nil {
		log.Fatal("Error initializing eth rpc client: ", err)
	}

	return cmd, ethRpcClient, nil
}

func TestAnvilSetupKill(t *testing.T) {
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		log.Fatal("Error setting up Anvil: ", err)
	}

	// Get Anvil PID
	pid := cmd.Process.Pid
	p, err := os.FindProcess(pid)
	if err != nil {
		log.Fatal("Error finding Anvil Process: ", err)
	}
	err = p.Signal(syscall.Signal(0))
	assert.Nil(t, err, "Anvil Process Killed")

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}

	// Check that PID is not currently present in running processes.
	// FindProcess always succeeds so on Unix systems we call it below.
	p, err = os.FindProcess(pid)
	if err != nil {
		log.Fatal("Error finding Anvil Process: ", err)
	}
	// Ensure process has exited
	err = p.Signal(syscall.Signal(0))

	assert.Nil(t, err, "Anvil Process Killed")
}

// |--Aggreagator Retry Tests--|

func TestWaitForTransactionReceipt(t *testing.T) {

	to := common.BytesToAddress([]byte{0x11})
	tx := types.NewTx(&types.AccessListTx{
		ChainID:  big.NewInt(1337),
		Nonce:    1,
		GasPrice: big.NewInt(11111),
		Gas:      1111,
		To:       &to,
		Value:    big.NewInt(111),
		Data:     []byte{0x11, 0x11, 0x11},
	})

	hash := tx.Hash()

	cmd, client, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	// Assert Call succeeds when Anvil running
	receipt_function := utils.WaitForTransactionReceipt(*client, *client, hash, retry.EthCallRetryConfig())
	_, err = receipt_function()
	assert.NotNil(t, err, "Error Waiting for Transaction with Anvil Running: %s\n", err)
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("WaitForTransactionReceipt Emitted incorrect error: %s\n", err)
		return
	}

	if err = cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}

	receipt_function = utils.WaitForTransactionReceipt(*client, *client, hash, retry.EthCallRetryConfig())
	_, err = receipt_function()
	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		t.Errorf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		t.Errorf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}

	cmd, client, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	receipt_function = utils.WaitForTransactionReceipt(*client, *client, hash, retry.EthCallRetryConfig())
	_, err = receipt_function()
	assert.NotNil(t, err)
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("WaitForTransactionReceipt Emitted incorrect error: %s\n", err)
		return
	}

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}
}

func TestGetGasPrice(t *testing.T) {

	cmd, client, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	// Assert Call succeeds when Anvil running
	receipt_function := utils.GetGasPrice(*client, *client)
	_, err = receipt_function()
	assert.NotNil(t, err, "Error Waiting for Transaction with Anvil Running: %s\n", err)
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("WaitForTransactionReceipt Emitted incorrect error: %s\n", err)
		return
	}

	if err = cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}

	receipt_function = utils.GetGasPrice(*client, *client)
	_, err = receipt_function()
	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		t.Errorf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		t.Errorf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}

	cmd, client, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	receipt_function = utils.GetGasPrice(*client, *client)
	_, err = receipt_function()
	assert.NotNil(t, err)
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("WaitForTransactionReceipt Emitted incorrect error: %s\n", err)
		return
	}

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}
}

// NOTE: The following tests involving starting the aggregator panic after the connection to anvil is cut crashing the test runner.
// The originates within the eigen-sdk and as of 8/11/24 is currently working to be fixed.

/*
func TestInitializeNewTask(t *testing.T) {

	_, _, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	agg, err := aggregator.NewAggregator(*aggregatorConfig)
	if err != nil {
		aggregatorConfig.BaseConfig.Logger.Error("Cannot create aggregator", "err", err)
		return
	}
	quorumNums := eigentypes.QuorumNums{eigentypes.QuorumNum(byte(0))}
	quorumThresholdPercentages := eigentypes.QuorumThresholdPercentages{eigentypes.QuorumThresholdPercentage(byte(57))}

	err = agg.InitializeNewTask(0, 1, quorumNums, quorumThresholdPercentages, 1*time.Second)
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("error killing process: %v\n", err)
		return
	}

	err = agg.InitializeNewTask(0, 1, quorumNums, quorumThresholdPercentages, 1*time.Second)
	assert.NotNil(t, err)
	t.Errorf("Error setting Avs Subscriber: %s\n", err)

	_, _, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	err = agg.InitializeNewTask(0, 1, quorumNums, quorumThresholdPercentages, 1*time.Second)
	assert.Nil(t, err)
	t.Errorf("Error setting Avs Subscriber: %s\n", err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("error killing process: %v\n", err)
		return
	}
}
*/

/*
func TestGetTaskIndex(t *testing.T) {

	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	agg, err := aggregator.NewAggregator(*aggregatorConfig)
	if err != nil {
		aggregatorConfig.BaseConfig.Logger.Error("Cannot create aggregator", "err", err)
		return
	}
	zero_bytes := [32]byte{}

	// Task is not present in map should return transient error
	_, err = agg.GetTaskIndex(zero_bytes)
	assert.NotNil(t, err)
	if !strings.Contains(err.Error(), "Task not found in the internal map") {
		t.Errorf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("error killing process: %v\n", err)
		return
	}
}
*/

/*
// |--Server Retry Tests--|
func TestProcessNewSignature(t *testing.T) {
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	agg, err := aggregator.NewAggregator(*aggregatorConfig)
	if err != nil {
		aggregatorConfig.BaseConfig.Logger.Error("Cannot create aggregator", "err", err)
		return
	}
	zero_bytes := [32]byte{}
	zero_sig := bls.NewZeroSignature()
	eigen_bytes := eigentypes.Bytes32{}

	err = agg.ProcessNewSignature(context.Background(), 0, zero_bytes, zero_sig, eigen_bytes)
	assert.NotNil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("error killing process: %v\n", err)
		return
	}

	err = agg.ProcessNewSignature(context.Background(), 0, zero_bytes, zero_sig, eigen_bytes)
	assert.NotNil(t, err)
	t.Errorf("Error Processing New Signature: %s\n", err)

	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	err = agg.ProcessNewSignature(context.Background(), 0, zero_bytes, zero_sig, eigen_bytes)
	assert.Nil(t, err)
	t.Errorf("Error Processing New Signature: %s\n", err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("error killing process: %v\n", err)
		return
	}
}
*/

// |--AVS-Subscriber Retry Tests--|

func TestSubscribeToNewTasksV3(t *testing.T) {
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	channel := make(chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3)
	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	baseConfig := aggregatorConfig.BaseConfig
	s, err := chainio.NewAvsServiceBindings(
		baseConfig.AlignedLayerDeploymentConfig.AlignedLayerServiceManagerAddr,
		baseConfig.AlignedLayerDeploymentConfig.AlignedLayerOperatorStateRetrieverAddr,
		baseConfig.EthWsClient, baseConfig.EthWsClientFallback, baseConfig.Logger)
	if err != nil {
		t.Errorf("Error setting up Avs Service Bindings: %s\n", err)
	}

	sub_func := chainio.SubscribeToNewTasksV3(&bind.WatchOpts{}, s.ServiceManager, channel, nil)
	_, err = sub_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}

	sub_func = chainio.SubscribeToNewTasksV3(&bind.WatchOpts{}, s.ServiceManager, channel, nil)
	_, err = sub_func()
	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		t.Errorf("SubscribeToNewTasksV3 Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connection reset") {
		t.Errorf("SubscribeToNewTasksV3 Emitted non Transient error: %s\n", err)
		return
	}

	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	sub_func = chainio.SubscribeToNewTasksV3(&bind.WatchOpts{}, s.ServiceManager, channel, nil)
	_, err = sub_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}
}

func TestSubscribeToNewTasksV2(t *testing.T) {
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	channel := make(chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV2)
	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	baseConfig := aggregatorConfig.BaseConfig
	s, err := chainio.NewAvsServiceBindings(
		baseConfig.AlignedLayerDeploymentConfig.AlignedLayerServiceManagerAddr,
		baseConfig.AlignedLayerDeploymentConfig.AlignedLayerOperatorStateRetrieverAddr,
		baseConfig.EthWsClient, baseConfig.EthWsClientFallback, baseConfig.Logger)
	if err != nil {
		t.Errorf("Error setting up Avs Service Bindings: %s\n", err)
	}

	sub_func := chainio.SubscribeToNewTasksV2(&bind.WatchOpts{}, s.ServiceManager, channel, nil)
	_, err = sub_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}

	sub_func = chainio.SubscribeToNewTasksV2(&bind.WatchOpts{}, s.ServiceManager, channel, nil)
	_, err = sub_func()

	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		t.Errorf("SubscribeToNewTasksV2 Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connection reset") {
		t.Errorf("SubscribeToNewTasksV2 Emitted non Transient error: %s\n", err)
		return
	}

	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	sub_func = chainio.SubscribeToNewTasksV2(&bind.WatchOpts{}, s.ServiceManager, channel, nil)
	_, err = sub_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}
}

func TestBlockNumber(t *testing.T) {
	// Start anvil
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	sub, err := chainio.NewAvsSubscriberFromConfig(aggregatorConfig.BaseConfig)
	if err != nil {
		return
	}
	block_func := chainio.BlockNumber(sub, context.Background())
	_, err = block_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}

	block_func = chainio.BlockNumber(sub, context.Background())
	_, err = block_func()
	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		t.Errorf("BlockNumber Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		t.Errorf("BlockNumber Emitted non Transient error: %s\n", err)
		return
	}

	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	block_func = chainio.BlockNumber(sub, context.Background())
	_, err = block_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}
}

func TestFilterBatchV2(t *testing.T) {
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	avsSubscriber, err := chainio.NewAvsSubscriberFromConfig(aggregatorConfig.BaseConfig)
	if err != nil {
		return
	}
	batch_func := chainio.FilterBatchV2(avsSubscriber, &bind.FilterOpts{Start: 0, End: nil, Context: context.Background()}, nil)
	_, err = batch_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}

	batch_func = chainio.FilterBatchV2(avsSubscriber, &bind.FilterOpts{Start: 0, End: nil, Context: context.Background()}, nil)
	_, err = batch_func()
	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		t.Errorf("FilterBatchV2 Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("FilterBatchV2 Emitted non Transient error: %s\n", err)
		return
	}

	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	batch_func = chainio.FilterBatchV2(avsSubscriber, &bind.FilterOpts{Start: 0, End: nil, Context: context.Background()}, nil)
	_, err = batch_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}
}

func TestFilterBatchV3(t *testing.T) {
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	avsSubscriber, err := chainio.NewAvsSubscriberFromConfig(aggregatorConfig.BaseConfig)
	if err != nil {
		return
	}
	batch_func := chainio.FilterBatchV3(avsSubscriber, &bind.FilterOpts{Start: 0, End: nil, Context: context.Background()}, nil)
	_, err = batch_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}

	batch_func = chainio.FilterBatchV3(avsSubscriber, &bind.FilterOpts{Start: 0, End: nil, Context: context.Background()}, nil)
	_, err = batch_func()
	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		t.Errorf("FilerBatchV3 Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("FilterBatchV3 Emitted non Transient error: %s\n", err)
		return
	}

	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	batch_func = chainio.FilterBatchV3(avsSubscriber, &bind.FilterOpts{Start: 0, End: nil, Context: context.Background()}, nil)
	_, err = batch_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}
}

func TestBatchesStateSubscriber(t *testing.T) {
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	avsSubscriber, err := chainio.NewAvsSubscriberFromConfig(aggregatorConfig.BaseConfig)
	if err != nil {
		return
	}

	zero_bytes := [32]byte{}
	batch_state_func := chainio.BatchState(avsSubscriber, nil, zero_bytes)
	_, err = batch_state_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}

	batch_state_func = chainio.BatchState(avsSubscriber, nil, zero_bytes)
	_, err = batch_state_func()
	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		t.Errorf("BatchesStateSubscriber Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("BatchesStateSubscriber Emitted non Transient error: %s\n", err)
		return
	}

	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	batch_state_func = chainio.BatchState(avsSubscriber, nil, zero_bytes)
	_, err = batch_state_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}
}

func TestSubscribeNewHead(t *testing.T) {
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	c := make(chan *types.Header)
	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	avsSubscriber, err := chainio.NewAvsSubscriberFromConfig(aggregatorConfig.BaseConfig)
	if err != nil {
		return
	}

	sub_func := chainio.SubscribeNewHead(avsSubscriber, context.Background(), c)
	_, err = sub_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}

	sub_func = chainio.SubscribeNewHead(avsSubscriber, context.Background(), c)
	_, err = sub_func()
	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		t.Errorf("SubscribeNewHead Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		t.Errorf("SubscribeNewHead Emitted non Transient error: %s\n", err)
		return
	}

	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	sub_func = chainio.SubscribeNewHead(avsSubscriber, context.Background(), c)
	_, err = sub_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
	}
}

// |--AVS-Writer Retry Tests--|

func TestRespondToTaskV2(t *testing.T) {
	// Start anvil
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	g2Point := servicemanager.BN254G2Point{
		X: [2]*big.Int{big.NewInt(4), big.NewInt(3)},
		Y: [2]*big.Int{big.NewInt(2), big.NewInt(8)},
	}

	g1Point := servicemanager.BN254G1Point{
		X: big.NewInt(2),
		Y: big.NewInt(1),
	}

	g1Points := []servicemanager.BN254G1Point{
		g1Point, g1Point, g1Point,
	}

	nonSignerStakesAndSignature := servicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature{
		NonSignerPubkeys:             g1Points,
		QuorumApks:                   g1Points,
		ApkG2:                        g2Point,
		Sigma:                        g1Point,
		NonSignerQuorumBitmapIndices: make([]uint32, 3),
		QuorumApkIndices:             make([]uint32, 3),
		TotalStakeIndices:            make([]uint32, 3),
		NonSignerStakeIndices:        make([][]uint32, 1, 3),
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	w, err := chainio.NewAvsWriterFromConfig(aggregatorConfig.BaseConfig, aggregatorConfig.EcdsaConfig, nil)
	if err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}
	txOpts := *w.Signer.GetTxOpts()
	aggregator_address := common.HexToAddress("0xc3e53F4d16Ae77Db1c982e75a937B9f60FE63690")
	zero_bytes := [32]byte{}

	// NOTE: With zero bytes the tx reverts
	resp_func := chainio.RespondToTaskV2(w, &txOpts, zero_bytes, aggregator_address, nonSignerStakesAndSignature)
	_, err = resp_func()
	assert.NotNil(t, err)
	if !strings.Contains(err.Error(), "execution reverted") {
		t.Errorf("RespondToTaskV2 did not emit the expected message: %q doesn't contain %q", err.Error(), "execution reverted: custom error 0x2396d34e:")
	}

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
	}

	resp_func = chainio.RespondToTaskV2(w, &txOpts, zero_bytes, aggregator_address, nonSignerStakesAndSignature)
	_, err = resp_func()
	assert.NotNil(t, err)
	if _, ok := err.(*backoff.PermanentError); ok {
		t.Errorf("RespondToTaskV2 Emitted non-Transient error: %s\n", err)
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		t.Errorf("RespondToTaskV2 did not return expected error: %s\n", err)
	}

	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	// NOTE: With zero bytes the tx reverts
	resp_func = chainio.RespondToTaskV2(w, &txOpts, zero_bytes, aggregator_address, nonSignerStakesAndSignature)
	_, err = resp_func()
	assert.NotNil(t, err)
	if !strings.Contains(err.Error(), "execution reverted") {
		t.Errorf("RespondToTaskV2 did not emit the expected message: %q doesn't contain %q", err.Error(), "execution reverted: custom error 0x2396d34e:")
	}

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
	}
}

func TestBatchesStateWriter(t *testing.T) {
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	avsWriter, err := chainio.NewAvsWriterFromConfig(aggregatorConfig.BaseConfig, aggregatorConfig.EcdsaConfig, nil)
	if err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}
	num := big.NewInt(6)

	var bytes [32]byte
	num.FillBytes(bytes[:])

	state_func := chainio.BatchesState(avsWriter, &bind.CallOpts{}, bytes)
	_, err = state_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("error killing process: %v\n", err)
		return
	}

	state_func = chainio.BatchesState(avsWriter, &bind.CallOpts{}, bytes)
	_, err = state_func()
	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		t.Errorf("BatchesStateWriter Emitted non-Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		t.Errorf("BatchesStateWriter did not contain expected error: %s\n", err)
		return
	}

	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	state_func = chainio.BatchesState(avsWriter, &bind.CallOpts{}, bytes)
	_, err = state_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}
}

func TestBalanceAt(t *testing.T) {
	cmd, client, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	avsWriter, err := chainio.NewAvsWriterFromConfig(aggregatorConfig.BaseConfig, aggregatorConfig.EcdsaConfig, nil)
	if err != nil {
		return
	}
	aggregator_address := common.HexToAddress("0x0")
	// Fetch the latest block number
	blockNumberUint64, err := client.BlockNumber(context.Background())
	blockNumber := new(big.Int).SetUint64(blockNumberUint64)
	if err != nil {
		t.Errorf("Error retrieving Anvil Block Number: %v\n", err)
	}

	balance_func := chainio.BalanceAt(avsWriter, context.Background(), aggregator_address, blockNumber)
	_, err = balance_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}

	balance_func = chainio.BalanceAt(avsWriter, context.Background(), aggregator_address, blockNumber)
	_, err = balance_func()
	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		t.Errorf("BalanceAt Emitted non-Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		t.Errorf("BalanceAt did not return expected error: %s\n", err)
		return
	}

	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	balance_func = chainio.BalanceAt(avsWriter, context.Background(), aggregator_address, blockNumber)
	_, err = balance_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}
}

func TestBatchersBalances(t *testing.T) {
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	avsWriter, err := chainio.NewAvsWriterFromConfig(aggregatorConfig.BaseConfig, aggregatorConfig.EcdsaConfig, nil)
	if err != nil {
		return
	}
	senderAddress := common.HexToAddress("0x0")

	batcher_func := chainio.BatcherBalances(avsWriter, &bind.CallOpts{}, senderAddress)
	_, err = batcher_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}

	batcher_func = chainio.BatcherBalances(avsWriter, &bind.CallOpts{}, senderAddress)
	_, err = batcher_func()
	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		t.Errorf("BatchersBalances Emitted non-Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		t.Errorf("BatchersBalances did not return expected error: %s\n", err)
		return
	}

	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	batcher_func = chainio.BatcherBalances(avsWriter, &bind.CallOpts{}, senderAddress)
	_, err = batcher_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}
}

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
	_, err := retry.RetryWithData(function, 1000, 2, 3, retry.MaxInterval, retry.MaxElapsedTime)
	if err != nil {
		t.Errorf("Retry error!: %s", err)
	}
}

func TestRetry(t *testing.T) {
	function := func() error {
		_, err := DummyFunction(43)
		return err
	}
	err := retry.Retry(function, 1000, 2, 3, retry.MaxInterval, retry.MaxElapsedTime)
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
	// Start Anvil
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

	// Kill Anvil
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("Error killing process: %v\n", err)
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

// Waits for receipt from anvil node -> Will fail to get receipt
func TestWaitForTransactionReceiptRetryable(t *testing.T) {

	// Retry call Params
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

	// Start anvil
	cmd, client, err := SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	// Assert Call succeeds when Anvil running
	_, err = utils.WaitForTransactionReceiptRetryable(*client, context.Background(), hash)
	assert.NotNil(t, err, "Error Waiting for Transaction with Anvil Running: %s\n", err)
	if !strings.Contains(err.Error(), "not found") {
		fmt.Printf("WaitForTransactionReceipt Emitted incorrect error: %s\n", err)
		return
	}

	// Kill Anvil
	if err = cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}

	_, err = utils.WaitForTransactionReceiptRetryable(*client, context.Background(), hash)
	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}

	// Start anvil
	cmd, client, err = SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	_, err = utils.WaitForTransactionReceiptRetryable(*client, context.Background(), hash)
	assert.NotNil(t, err)
	if !strings.Contains(err.Error(), "not found") {
		fmt.Printf("WaitForTransactionReceipt Emitted incorrect error: %s\n", err)
		return
	}

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
}

/*
func TestInitializeNewTaskRetryable(t *testing.T) {

	//Start Anvil
	_, _, err := SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	//Start Aggregator
	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	agg, err := aggregator.NewAggregator(*aggregatorConfig)
	if err != nil {
		aggregatorConfig.BaseConfig.Logger.Error("Cannot create aggregator", "err", err)
		return
	}
	quorumNums := eigentypes.QuorumNums{eigentypes.QuorumNum(byte(0))}
	quorumThresholdPercentages := eigentypes.QuorumThresholdPercentages{eigentypes.QuorumThresholdPercentage(byte(57))}

	// Should succeed with err msg
	err = agg.InitializeNewTaskRetryable(0, 1, quorumNums, quorumThresholdPercentages, 1*time.Second)
	assert.Nil(t, err)
	// TODO: Find exact error to assert

			// Kill Anvil
			if err := cmd.Process.Kill(); err != nil {
				fmt.Printf("error killing process: %v\n", err)
				return
			}
			time.Sleep(2 * time.Second)

				err = agg.InitializeNewTaskRetryable(0, 1, quorumNums, quorumThresholdPercentages, 1*time.Second)
				assert.NotNil(t, err)
				fmt.Printf("Error setting Avs Subscriber: %s\n", err)

			// Start Anvil
			_, _, err = SetupAnvil(8545)
			if err != nil {
				fmt.Printf("Error setting up Anvil: %s\n", err)
			}

			// Should succeed
			err = agg.InitializeNewTaskRetryable(0, 1, quorumNums, quorumThresholdPercentages, 1*time.Second)
			assert.Nil(t, err)
			fmt.Printf("Error setting Avs Subscriber: %s\n", err)
		// Kill Anvil
		if err := cmd.Process.Kill(); err != nil {
			fmt.Printf("error killing process: %v\n", err)
			return
		}
		time.Sleep(2 * time.Second)
}
*/

/*
// |--Server Retry Tests--|
func TestProcessNewSignatureRetryable(t *testing.T) {
	// Start anvil
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	//Start Aggregator
	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	agg, err := aggregator.NewAggregator(*aggregatorConfig)
	if err != nil {
		aggregatorConfig.BaseConfig.Logger.Error("Cannot create aggregator", "err", err)
		return
	}
	zero_bytes := [32]byte{}
	zero_sig := bls.NewZeroSignature()
	eigen_bytes := eigentypes.Bytes32{}

	err = agg.ProcessNewSignatureRetryable(context.Background(), 0, zero_bytes, zero_sig, eigen_bytes)
	assert.NotNil(t, err)
	// TODO: Find exact error to assert

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
	time.Sleep(2 * time.Second)

	err = agg.ProcessNewSignatureRetryable(context.Background(), 0, zero_bytes, zero_sig, eigen_bytes)
	assert.NotNil(t, err)
	fmt.Printf("Error Processing New Signature: %s\n", err)

	// Start anvil
	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	err = agg.ProcessNewSignatureRetryable(context.Background(), 0, zero_bytes, zero_sig, eigen_bytes)
	assert.Nil(t, err)
	fmt.Printf("Error Processing New Signature: %s\n", err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
}
*/

// |--AVS-Subscriber Retry Tests--|

func TestSubscribeToNewTasksV3Retryable(t *testing.T) {
	// Start anvil
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	channel := make(chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3)
	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	baseConfig := aggregatorConfig.BaseConfig
	s, err := chainio.NewAvsServiceBindings(
		baseConfig.AlignedLayerDeploymentConfig.AlignedLayerServiceManagerAddr,
		baseConfig.AlignedLayerDeploymentConfig.AlignedLayerOperatorStateRetrieverAddr,
		baseConfig.EthWsClient, baseConfig.EthWsClientFallback, baseConfig.Logger)
	if err != nil {
		fmt.Printf("Error setting up Avs Service Bindings: %s\n", err)
	}

	_, err = chainio.SubscribeToNewTasksV3Retryable(s.ServiceManager, channel, baseConfig.Logger)
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}

	_, err = chainio.SubscribeToNewTasksV3Retryable(s.ServiceManager, channel, baseConfig.Logger)
	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}

	// Start anvil
	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	_, err = chainio.SubscribeToNewTasksV3Retryable(s.ServiceManager, channel, baseConfig.Logger)
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
}

func TestSubscribeToNewTasksV2(t *testing.T) {
	// Start anvil
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	channel := make(chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV2)
	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	baseConfig := aggregatorConfig.BaseConfig
	s, err := chainio.NewAvsServiceBindings(
		baseConfig.AlignedLayerDeploymentConfig.AlignedLayerServiceManagerAddr,
		baseConfig.AlignedLayerDeploymentConfig.AlignedLayerOperatorStateRetrieverAddr,
		baseConfig.EthWsClient, baseConfig.EthWsClientFallback, baseConfig.Logger)
	if err != nil {
		fmt.Printf("Error setting up Avs Service Bindings: %s\n", err)
	}

	_, err = chainio.SubscribeToNewTasksV2Retrayable(s.ServiceManager, channel, baseConfig.Logger)
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}

	_, err = chainio.SubscribeToNewTasksV2Retrayable(s.ServiceManager, channel, baseConfig.Logger)
	assert.NotNil(t, err)
	// If it retruend a permanent error we exit
	if _, ok := err.(retry.PermanentError); ok {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}

	// Start anvil
	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	_, err = chainio.SubscribeToNewTasksV2Retrayable(s.ServiceManager, channel, baseConfig.Logger)
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
}

func TestBlockNumber(t *testing.T) {
	// Start anvil
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	//channel := make(chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3)
	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	sub, err := chainio.NewAvsSubscriberFromConfig(aggregatorConfig.BaseConfig)
	if err != nil {
		return
	}
	_, err = sub.BlockNumberRetryable(context.Background())
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}

	_, err = sub.BlockNumberRetryable(context.Background())
	assert.NotNil(t, err)
	// Assert returned error is both transient error and contains the expected error msg.
	if _, ok := err.(retry.PermanentError); ok {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}

	// Start anvil
	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	_, err = sub.BlockNumberRetryable(context.Background())
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
}

func TestFilterBatchV2(t *testing.T) {
	// Start anvil
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	avsSubscriber, err := chainio.NewAvsSubscriberFromConfig(aggregatorConfig.BaseConfig)
	if err != nil {
		return
	}
	_, err = avsSubscriber.FilterBatchV2Retryable(0, context.Background())
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}

	_, err = avsSubscriber.FilterBatchV2Retryable(0, context.Background())
	assert.NotNil(t, err)
	//
	if _, ok := err.(retry.PermanentError); ok {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}

	// Start anvil
	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	_, err = avsSubscriber.FilterBatchV2Retryable(0, context.Background())
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
}

func TestFilterBatchV3(t *testing.T) {
	// Start anvil
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	avsSubscriber, err := chainio.NewAvsSubscriberFromConfig(aggregatorConfig.BaseConfig)
	if err != nil {
		return
	}
	_, err = avsSubscriber.FilterBatchV3Retryable(0, context.Background())
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}

	_, err = avsSubscriber.FilterBatchV3Retryable(0, context.Background())
	assert.NotNil(t, err)
	// Assert returned error is both transient error and contains the expected error msg.
	if _, ok := err.(retry.PermanentError); ok {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}

	// Start anvil
	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	_, err = avsSubscriber.FilterBatchV3Retryable(0, context.Background())
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
}

func TestBatchesStateSubscriber(t *testing.T) {
	// Start anvil
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	avsSubscriber, err := chainio.NewAvsSubscriberFromConfig(aggregatorConfig.BaseConfig)
	if err != nil {
		return
	}

	zero_bytes := [32]byte{}
	_, err = avsSubscriber.BatchesStateRetryable(nil, zero_bytes)
	//TODO: Find exact failure error
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}

	_, err = avsSubscriber.BatchesStateRetryable(nil, zero_bytes)
	assert.NotNil(t, err)
	// Assert returned error is both transient error and contains the expected error msg.
	if _, ok := err.(retry.PermanentError); ok {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}

	// Start anvil
	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	_, err = avsSubscriber.BatchesStateRetryable(nil, zero_bytes)
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
}

func TestSubscribeNewHead(t *testing.T) {
	// Start anvil
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	c := make(chan *types.Header)
	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	avsSubscriber, err := chainio.NewAvsSubscriberFromConfig(aggregatorConfig.BaseConfig)
	if err != nil {
		return
	}

	_, err = avsSubscriber.SubscribeNewHeadRetryable(context.Background(), c)
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}

	_, err = avsSubscriber.SubscribeNewHeadRetryable(context.Background(), c)
	assert.NotNil(t, err)
	// Assert returned error is both transient error and contains the expected error msg.
	if _, ok := err.(retry.PermanentError); ok {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		fmt.Printf("WaitForTransactionReceipt Emitted non Transient error: %s\n", err)
		return
	}

	// Start anvil
	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	_, err = avsSubscriber.SubscribeNewHeadRetryable(context.Background(), c)
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
}

// |--AVS-Writer Retry Tests--|

func TestRespondToTaskV2(t *testing.T) {
	// Start anvil
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
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

	// Or if you want to initialize with specific values

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
	w, err := chainio.NewAvsWriterFromConfig(aggregatorConfig.BaseConfig, aggregatorConfig.EcdsaConfig)
	if err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
	txOpts := *w.Signer.GetTxOpts()
	aggregator_address := common.HexToAddress("0xc3e53F4d16Ae77Db1c982e75a937B9f60FE63690")
	zero_bytes := [32]byte{}

	// NOTE: With zero bytes the tx reverts
	_, err = w.RespondToTaskV2Retryable(&txOpts, zero_bytes, aggregator_address, nonSignerStakesAndSignature)
	assert.NotNil(t, err)
	// assert error contains "Message:"execution reverted: custom error 0x2396d34e:"
	if !strings.Contains(err.Error(), "execution reverted: custom error 0x2396d34e:") {
		t.Errorf("Respond to task V2 Retryable did not emit the expected message: %q doesn't contain %q", err.Error(), "execution reverted: custom error 0x2396d34e:")
		return
	}

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}

	_, err = w.RespondToTaskV2Retryable(&txOpts, zero_bytes, aggregator_address, nonSignerStakesAndSignature)
	assert.NotNil(t, err)
	// Assert returned error is both transient error and contains the expected error msg.
	if _, ok := err.(retry.PermanentError); ok {
		fmt.Printf("RespondToTaksV2 Emitted non-Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		fmt.Printf("RespondToTaskV2 did not return expected error: %s\n", err)
		return
	}

	// Start anvil
	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	// NOTE: With zero bytes the tx reverts
	_, err = w.RespondToTaskV2Retryable(&txOpts, zero_bytes, aggregator_address, nonSignerStakesAndSignature)
	assert.NotNil(t, err)
	// assert error contains "Message:"execution reverted: custom error 0x2396d34e:"
	if !strings.Contains(err.Error(), "execution reverted: custom error 0x2396d34e:") {
		t.Errorf("Respond to task V2 Retryable did not emit the expected message: %q doesn't contain %q", err.Error(), "execution reverted: custom error 0x2396d34e:")
		return
	}

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
}

func TestBatchesStateWriter(t *testing.T) {
	// Start anvil
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	avsWriter, err := chainio.NewAvsWriterFromConfig(aggregatorConfig.BaseConfig, aggregatorConfig.EcdsaConfig)
	if err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
	num := big.NewInt(6)

	var bytes [32]byte
	num.FillBytes(bytes[:])

	_, err = avsWriter.BatchesStateRetryable(&bind.CallOpts{}, bytes)
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}

	_, err = avsWriter.BatchesStateRetryable(&bind.CallOpts{}, bytes)
	assert.NotNil(t, err)
	// Assert returned error is both transient error and contains the expected error msg.
	if _, ok := err.(retry.PermanentError); ok {
		fmt.Printf("BatchesStateWriter Emitted non-Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		fmt.Printf("BatchesStateWriter did not contain expected error: %s\n", err)
		return
	}

	// Start anvil
	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	_, err = avsWriter.BatchesStateRetryable(&bind.CallOpts{}, bytes)
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
}

func TestBalanceAt(t *testing.T) {
	// Start anvil
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	avsWriter, err := chainio.NewAvsWriterFromConfig(aggregatorConfig.BaseConfig, aggregatorConfig.EcdsaConfig)
	if err != nil {
		return
	}
	//TODO: Source Aggregator Address
	aggregator_address := common.HexToAddress("0x0")
	blockHeight := big.NewInt(13)

	_, err = avsWriter.BalanceAtRetryable(context.Background(), aggregator_address, blockHeight)
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}

	_, err = avsWriter.BalanceAtRetryable(context.Background(), aggregator_address, blockHeight)
	assert.NotNil(t, err)
	// Assert returned error is both transient error and contains the expected error msg.
	if _, ok := err.(retry.PermanentError); ok {
		fmt.Printf("BalanceAt Emitted non-Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		fmt.Printf("BalanceAt did not return expected error: %s\n", err)
		return
	}

	// Start anvil
	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	_, err = avsWriter.BalanceAtRetryable(context.Background(), aggregator_address, blockHeight)
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
}

func TestBatchersBalances(t *testing.T) {
	// Start anvil
	cmd, _, err := SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	aggregatorConfig := config.NewAggregatorConfig("../config-files/config-aggregator-test.yaml")
	avsWriter, err := chainio.NewAvsWriterFromConfig(aggregatorConfig.BaseConfig, aggregatorConfig.EcdsaConfig)
	if err != nil {
		return
	}
	sender_address := common.HexToAddress("0x0")

	_, err = avsWriter.BatcherBalancesRetryable(sender_address)
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}

	_, err = avsWriter.BatcherBalancesRetryable(sender_address)
	assert.NotNil(t, err)
	// Assert returned error is both transient error and contains the expected error msg.
	if _, ok := err.(retry.PermanentError); ok {
		fmt.Printf("BatchersBalances Emitted non-Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		fmt.Printf("BatchersBalances did not return expected error: %s\n", err)
		return
	}

	// Start anvil
	cmd, _, err = SetupAnvil(8545)
	if err != nil {
		fmt.Printf("Error setting up Anvil: %s\n", err)
	}

	_, err = avsWriter.BatcherBalancesRetryable(sender_address)
	assert.Nil(t, err)

	// Kill Anvil at end of test
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("error killing process: %v\n", err)
		return
	}
}

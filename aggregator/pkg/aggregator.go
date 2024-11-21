package pkg

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/yetanotherco/aligned_layer/metrics"

	sdkclients "github.com/Layr-Labs/eigensdk-go/chainio/clients"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/services/avsregistry"
	blsagg "github.com/Layr-Labs/eigensdk-go/services/bls_aggregation"
	oppubkeysserv "github.com/Layr-Labs/eigensdk-go/services/operatorsinfo"
	eigentypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/crypto"
	servicemanager "github.com/yetanotherco/aligned_layer/contracts/bindings/AlignedLayerServiceManager"
	"github.com/yetanotherco/aligned_layer/core/chainio"
	"github.com/yetanotherco/aligned_layer/core/config"
	"github.com/yetanotherco/aligned_layer/core/types"
	"github.com/yetanotherco/aligned_layer/core/utils"
)

// FIXME(marian): Read this from Aligned contract directly
const QUORUM_NUMBER = byte(0)
const QUORUM_THRESHOLD = byte(67)

// Aggregator stores TaskResponse for a task here
type TaskResponses = []types.SignedTaskResponse

// BatchData stores the data of a batch, for use in map BatchIdentifierHash -> BatchData
type BatchData struct {
	BatchMerkleRoot [32]byte
	SenderAddress   [20]byte
	// Auxiliary data for id
	Index           uint32
	IdentifierHash  [32]byte
	CreatedBlock    uint64
}

type BatchDataCache struct {
	// Mutex to protect all fields in this structure
	// Contents are passed by copy, so only internal access needs protection
	mx *sync.Mutex

	// BLS Signature Service returns an Index
	// Since our ID is not an idx, we build this cache
	// Note: In case of a reboot, this doesn't need to be loaded,
	// and can start from zero
	identifierByIndex map[uint32][32]byte

	// Stores the TaskResponse for each batch by batchIdentifierHash
	byIdentifierHash  map[[32]byte]BatchData

	// This task index is to communicate with the local BLS Service
	// Note: In case of a reboot it can start from 0 again
	nextIndex         uint32
	// This is used to determine the range of tasks to delete on GC
	lastDeletedIndex  uint32

	logger logging.Logger
}

func NewBatchDataCache(logger logging.Logger) BatchDataCache {
	return BatchDataCache{
		identifierByIndex: make(map[uint32][32]byte),
		byIdentifierHash: make(map[[32]byte]BatchData),
		nextIndex: 0,
		lastDeletedIndex: 0,
		logger: logger,
	}
}

// Locks the mutex, logs about it and returns a function that does the opposite to defer or call.
func (cache *BatchDataCache) lock(routine string) func() {
	cache.logger.Info("- Waiting for lock: " + routine)
	cache.mx.Lock()
	cache.logger.Info("- Locked Resources: " + routine)

	return func() {
		cache.mx.Unlock()
		cache.logger.Info("- Unlocked Resources: " + routine)
	}
}

func (cache *BatchDataCache) GetTaskDataByIndex(idx uint32) (BatchData, bool) {
	unlock := cache.lock("Query by index")
	defer unlock()

	hash, ok := cache.identifierByIndex[idx]
	if !ok {
		return BatchData{}, false
	}

	data, ok := cache.byIdentifierHash[hash]
	if !ok {
		cache.logger.Warn("Dangling hash in batch data cache",
			"batchIndex", idx,
			"batchIdentifierHash", "0x"+hex.EncodeToString(hash[:]))
		delete(cache.identifierByIndex, idx)
		return BatchData{}, false
	}

	return data, ok
}

func (cache *BatchDataCache) GetTaskDataByIdentifierHash(idHash [32]byte) (BatchData, bool) {
	unlock := cache.lock("Query by identifier task")
	defer unlock()

	data, ok := cache.byIdentifierHash[idHash]
	if !ok {
		return data, false
	}
	if idHash != data.IdentifierHash {
		cache.logger.Error("Inconsistency detected in batch data",
			"batchIndex", data.Index,
			"expectedIdentifierHash", "0x"+hex.EncodeToString(idHash[:]),
			"foundIdentifierHash", "0x"+hex.EncodeToString(data.IdentifierHash[:]))
	}
	hash, ok := cache.identifierByIndex[data.Index]
	if !ok {
		cache.logger.Warn("Missing hash entry for batch data cache entry, inserting in map",
			"batchIndex", data.Index,
			"batchIdentifierHash", "0x"+hex.EncodeToString(idHash[:]))
		// NOTE: this case is more harmless because we're not breaking a different entry,
		// just recover by correcting the map.
		cache.identifierByIndex[data.Index] = idHash
		return data, true
	}
	if hash != data.IdentifierHash {
		cache.logger.Error("Inconsistency detected in identifier-index maps",
			"batchIndex", data.Index,
			"expectedIdentifierHash", "0x"+hex.EncodeToString(idHash[:]),
			"foundIdentifierHash", "0x"+hex.EncodeToString(hash[:]))
	}

	return data, ok
}

func (cache *BatchDataCache) PopulateBatchData(merkleRoot [32]byte, senderAddr [20]byte, idHash [32]byte, createdBlock uint32) (BatchData, error) {
	unlock := cache.lock("Add new task")
	defer unlock()

	idx := cache.nextIndex
	_, ok := cache.identifierByIndex[idx]
	for ok {
		cache.logger.Error("Attempted to reuse batch index, bumping", "index", idx)
		idx++
		_, ok = cache.identifierByIndex[idx]
	}
	// If we return early, this will be the first empty index
	cache.nextIndex = idx

	newData := BatchData{
		BatchMerkleRoot: merkleRoot,
		SenderAddress: senderAddr,
		IdentifierHash: idHash,
		// FIXME: why do we use two different types for this?
		CreatedBlock: uint64(createdBlock),
		Index: idx,
	}

	oldData, ok := cache.byIdentifierHash[idHash]
	if ok && (oldData.BatchMerkleRoot != merkleRoot ||
		oldData.SenderAddress != senderAddr ||
		oldData.IdentifierHash != idHash ||
		oldData.CreatedBlock != uint64(createdBlock)) {

		cache.logger.Error("Found different data on new hash addition, either a hash collision or a bug",
			"oldSenderAddress", oldData.SenderAddress,
			"oldMerkleRoot", oldData.BatchMerkleRoot,
			"oldIdentifierHash", oldData.IdentifierHash,
			"oldCreatedBlock", oldData.CreatedBlock,
			"newSenderAddress", newData.SenderAddress,
			"newMerkleRoot", newData.BatchMerkleRoot,
			"newIdentifierHash", newData.IdentifierHash,
			"newCreatedBlock", newData.CreatedBlock)

		return oldData, fmt.Errorf("Data mismatch for hash 0x%s", hex.EncodeToString(idHash[:]))
	}
	if ok {
		cache.logger.Warn("Trying to add duplicate entry to batch data cache", "identifierHash", "0x"+hex.EncodeToString(idHash[:]))
		return oldData, nil
	}

	cache.byIdentifierHash[idHash] = newData
	cache.identifierByIndex[idx] = idHash

	cache.nextIndex += 1
	cache.logger.Info("New task added", "batchIndex", idx, "batchIdentifierHash", "0x"+hex.EncodeToString(idHash[:]))

	return newData, nil
}

func (cache *BatchDataCache) DeleteOlderByIdentifierHash(id [32]byte) {
	unlock := cache.lock("Cleaning finalized tasks")
	defer unlock()

	data, ok := cache.byIdentifierHash[id]
	if !ok {
		return
	}
	var i uint32
	for i = cache.lastDeletedIndex; i <= data.Index; i++ {
		hash, ok := cache.identifierByIndex[i]
		if !ok {
			cache.logger.Warn("Task not found for cleanup", "taskIndex", i)
			continue
		}
		delete(cache.identifierByIndex, i)
		hashStr := "0x"+hex.EncodeToString(hash[:])
		_, ok = cache.byIdentifierHash[hash]
		if !ok {
			cache.logger.Warn("Task data not found for cleanup", "taskIndex", i, "identifierHash", hashStr)
			continue
		}
		delete(cache.byIdentifierHash, hash)
		cache.logger.Info("Task deleted", "taskIndex", i, "identifierHash", hashStr)
	}
	cache.lastDeletedIndex = i - 1
}

type Aggregator struct {
	AggregatorConfig      *config.AggregatorConfig
	NewBatchChan          chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3
	avsReader             *chainio.AvsReader
	avsSubscriber         *chainio.AvsSubscriber
	avsWriter             *chainio.AvsWriter
	taskSubscriber        chan error
	blsAggregationService blsagg.BlsAggregationService
	batchDataCache        BatchDataCache

	// Mutex to protect ethereum wallet
	walletMutex *sync.Mutex

	logger logging.Logger

	// Metrics
	metricsReg *prometheus.Registry
	metrics    *metrics.Metrics

	// Telemetry
	telemetry *Telemetry
}

func NewAggregator(aggregatorConfig config.AggregatorConfig) (*Aggregator, error) {
	newBatchChan := make(chan *servicemanager.ContractAlignedLayerServiceManagerNewBatchV3)

	logger := aggregatorConfig.BaseConfig.Logger

	// Metrics
	reg := prometheus.NewRegistry()
	aggregatorMetrics := metrics.NewMetrics(aggregatorConfig.Aggregator.MetricsIpPortAddress, reg, logger)

	// Telemetry
	aggregatorTelemetry := NewTelemetry(aggregatorConfig.Aggregator.TelemetryIpPortAddress, logger)

	avsReader, err := chainio.NewAvsReaderFromConfig(aggregatorConfig.BaseConfig, aggregatorConfig.EcdsaConfig)
	if err != nil {
		return nil, err
	}

	avsSubscriber, err := chainio.NewAvsSubscriberFromConfig(aggregatorConfig.BaseConfig)
	if err != nil {
		return nil, err
	}

	avsWriter, err := chainio.NewAvsWriterFromConfig(aggregatorConfig.BaseConfig, aggregatorConfig.EcdsaConfig, aggregatorMetrics)
	if err != nil {
		return nil, err
	}

	batchDataCache := NewBatchDataCache(aggregatorConfig.BaseConfig.Logger)

	chainioConfig := sdkclients.BuildAllConfig{
		EthHttpUrl:                 aggregatorConfig.BaseConfig.EthRpcUrl,
		EthWsUrl:                   aggregatorConfig.BaseConfig.EthWsUrl,
		RegistryCoordinatorAddr:    aggregatorConfig.BaseConfig.AlignedLayerDeploymentConfig.AlignedLayerRegistryCoordinatorAddr.Hex(),
		OperatorStateRetrieverAddr: aggregatorConfig.BaseConfig.AlignedLayerDeploymentConfig.AlignedLayerOperatorStateRetrieverAddr.Hex(),
		AvsName:                    "AlignedLayer",
		PromMetricsIpPortAddress:   ":9090",
	}

	aggregatorPrivateKey := aggregatorConfig.EcdsaConfig.PrivateKey

	clients, err := sdkclients.BuildAll(chainioConfig, aggregatorPrivateKey, logger)
	if err != nil {
		logger.Errorf("Cannot create sdk clients", "err", err)
		return nil, err
	}

	// This is a dummy "hash function" made to fulfill the BLS aggregator service API requirements.
	// When operators respond to a task, a call to `ProcessNewSignature` is made. In `v0.1.6` of the eigensdk,
	// this function required an argument `TaskResponseDigest`, which has changed to just `TaskResponse` in v0.1.9.
	// The digest we used in v0.1.6 was just the batch merkle root. To continue with the same idea, the hashing
	// function is set as the following one, which does nothing more than output the input it receives, which in
	// our case will be the batch merkle root. If wanted, we could define a real hash function here but there should
	// not be any need to re-hash the batch merkle root.
	hashFunction := func(taskResponse eigentypes.TaskResponse) (eigentypes.TaskResponseDigest, error) {
		taskResponseDigest, ok := taskResponse.([32]byte)
		if !ok {
			return eigentypes.TaskResponseDigest{}, fmt.Errorf("TaskResponse is not a 32-byte value")
		}
		return taskResponseDigest, nil
	}

	operatorPubkeysService := oppubkeysserv.NewOperatorsInfoServiceInMemory(context.Background(), clients.AvsRegistryChainSubscriber, clients.AvsRegistryChainReader, nil, oppubkeysserv.Opts{}, logger)
	avsRegistryService := avsregistry.NewAvsRegistryServiceChainCaller(avsReader.ChainReader, operatorPubkeysService, logger)
	blsAggregationService := blsagg.NewBlsAggregatorService(avsRegistryService, hashFunction, logger)

	aggregator := Aggregator{
		AggregatorConfig: &aggregatorConfig,
		avsReader:        avsReader,
		avsSubscriber:    avsSubscriber,
		avsWriter:        avsWriter,
		NewBatchChan:     newBatchChan,
		batchDataCache:   batchDataCache,
		walletMutex:      &sync.Mutex{},

		blsAggregationService: blsAggregationService,
		logger:                logger,
		metricsReg:            reg,
		metrics:               aggregatorMetrics,
		telemetry:             aggregatorTelemetry,
	}

	return &aggregator, nil
}

func (agg *Aggregator) Start(ctx context.Context) error {
	agg.logger.Infof("Starting aggregator...")

	go func() {
		err := agg.ServeOperators()
		if err != nil {
			agg.logger.Fatal("Error listening for tasks", "err", err)
		}
	}()

	var metricsErrChan <-chan error
	if agg.AggregatorConfig.Aggregator.EnableMetrics {
		metricsErrChan = agg.metrics.Start(ctx, agg.metricsReg)
	} else {
		metricsErrChan = make(chan error, 1)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-metricsErrChan:
			agg.logger.Fatal("Metrics server failed", "err", err)
		case blsAggServiceResp := <-agg.blsAggregationService.GetResponseChannel():
			agg.logger.Info("Received response from BLS aggregation service",
				"taskIndex", blsAggServiceResp.TaskIndex)

			go agg.handleBlsAggServiceResponse(blsAggServiceResp)
		}
	}
}

const MaxSentTxRetries = 5

func (agg *Aggregator) handleBlsAggServiceResponse(blsAggServiceResp blsagg.BlsAggregationServiceResponse) {
	defer func() {
		err := recover() //stops panics
		if err != nil {
			agg.logger.Error("handleBlsAggServiceResponse recovered from panic", "err", err)
		}
	}()

	taskData, ok := agg.GetTaskDataByIndex(blsAggServiceResp.TaskIndex)
	if !ok {
		agg.logger.Error("Missing task", "taskIndex", blsAggServiceResp.TaskIndex)
		return
	}
	batchIdentifierHash := taskData.IdentifierHash
	batchData := taskData // FIXME
	taskCreatedBlock := taskData.CreatedBlock

	// Finish task trace once the task is processed (either successfully or not)
	defer agg.telemetry.FinishTrace(batchData.BatchMerkleRoot)

	if blsAggServiceResp.Err != nil {
		agg.telemetry.LogTaskError(batchData.BatchMerkleRoot, blsAggServiceResp.Err)
		agg.logger.Error("BlsAggregationServiceResponse contains an error", "err", blsAggServiceResp.Err, "batchIdentifierHash", hex.EncodeToString(batchIdentifierHash[:]))
		return
	}
	nonSignerPubkeys := []servicemanager.BN254G1Point{}
	for _, nonSignerPubkey := range blsAggServiceResp.NonSignersPubkeysG1 {
		nonSignerPubkeys = append(nonSignerPubkeys, utils.ConvertToBN254G1Point(nonSignerPubkey))
	}
	quorumApks := []servicemanager.BN254G1Point{}
	for _, quorumApk := range blsAggServiceResp.QuorumApksG1 {
		quorumApks = append(quorumApks, utils.ConvertToBN254G1Point(quorumApk))
	}

	nonSignerStakesAndSignature := servicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature{
		NonSignerPubkeys:             nonSignerPubkeys,
		QuorumApks:                   quorumApks,
		ApkG2:                        utils.ConvertToBN254G2Point(blsAggServiceResp.SignersApkG2),
		Sigma:                        utils.ConvertToBN254G1Point(blsAggServiceResp.SignersAggSigG1.G1Point),
		NonSignerQuorumBitmapIndices: blsAggServiceResp.NonSignerQuorumBitmapIndices,
		QuorumApkIndices:             blsAggServiceResp.QuorumApkIndices,
		TotalStakeIndices:            blsAggServiceResp.TotalStakeIndices,
		NonSignerStakeIndices:        blsAggServiceResp.NonSignerStakeIndices,
	}

	agg.telemetry.LogQuorumReached(batchData.BatchMerkleRoot)

	agg.logger.Info("Threshold reached", "taskIndex", blsAggServiceResp.TaskIndex,
		"batchIdentifierHash", "0x"+hex.EncodeToString(batchIdentifierHash[:]))

	agg.logger.Info("Maybe waiting one block to send aggregated response onchain",
		"taskIndex", blsAggServiceResp.TaskIndex,
		"batchIdentifierHash", "0x"+hex.EncodeToString(batchIdentifierHash[:]),
		"taskCreatedBlock", taskCreatedBlock)

	err := agg.avsSubscriber.WaitForOneBlock(taskCreatedBlock)
	if err != nil {
		agg.logger.Error("Error waiting for one block, sending anyway", "err", err)
	}

	agg.logger.Info("Sending aggregated response onchain", "taskIndex", blsAggServiceResp.TaskIndex,
		"batchIdentifierHash", "0x"+hex.EncodeToString(batchIdentifierHash[:]), "merkleRoot", "0x"+hex.EncodeToString(batchData.BatchMerkleRoot[:]))
	receipt, err := agg.sendAggregatedResponse(batchIdentifierHash, batchData.BatchMerkleRoot, batchData.SenderAddress, nonSignerStakesAndSignature)
	if err == nil {
		// In some cases, we may fail to retrieve the receipt for the transaction.
		txHash := "Unknown"
		if receipt != nil {
			txHash = receipt.TxHash.String()
		}
		agg.telemetry.TaskSentToEthereum(batchData.BatchMerkleRoot, txHash)
		agg.logger.Info("Aggregator successfully responded to task",
			"taskIndex", blsAggServiceResp.TaskIndex,
			"batchIdentifierHash", "0x"+hex.EncodeToString(batchIdentifierHash[:]))

		return
	}

	agg.logger.Error("Aggregator failed to respond to task, this batch will be lost",
		"err", err,
		"taskIndex", blsAggServiceResp.TaskIndex,
		"merkleRoot", "0x"+hex.EncodeToString(batchData.BatchMerkleRoot[:]),
		"senderAddress", "0x"+hex.EncodeToString(batchData.SenderAddress[:]),
		"batchIdentifierHash", "0x"+hex.EncodeToString(batchIdentifierHash[:]))
	agg.telemetry.LogTaskError(batchData.BatchMerkleRoot, err)
}

// / Sends response to contract and waits for transaction receipt
// / Returns error if it fails to send tx or receipt is not found
func (agg *Aggregator) sendAggregatedResponse(batchIdentifierHash [32]byte, batchMerkleRoot [32]byte, senderAddress [20]byte, nonSignerStakesAndSignature servicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature) (*gethtypes.Receipt, error) {

	agg.walletMutex.Lock()
	agg.logger.Infof("- Locked Wallet Resources: Sending aggregated response for batch",
		"merkleRoot", hex.EncodeToString(batchMerkleRoot[:]),
		"senderAddress", hex.EncodeToString(senderAddress[:]),
		"batchIdentifierHash", hex.EncodeToString(batchIdentifierHash[:]))

	// This function is a callback that is called when the gas price is bumped on the avsWriter.SendAggregatedResponse
	onGasPriceBumped := func(bumpedGasPrice *big.Int) {
		agg.metrics.IncBumpedGasPriceForAggregatedResponse()
		agg.telemetry.BumpedTaskGasPrice(batchMerkleRoot, bumpedGasPrice.String())
	}
	receipt, err := agg.avsWriter.SendAggregatedResponse(
		batchIdentifierHash,
		batchMerkleRoot,
		senderAddress,
		nonSignerStakesAndSignature,
		agg.AggregatorConfig.Aggregator.GasBaseBumpPercentage,
		agg.AggregatorConfig.Aggregator.GasBumpIncrementalPercentage,
		agg.AggregatorConfig.Aggregator.TimeToWaitBeforeBump,
		onGasPriceBumped,
	)
	if err != nil {
		agg.walletMutex.Unlock()
		agg.logger.Infof("- Unlocked Wallet Resources: Error sending aggregated response for batch %s. Error: %s", hex.EncodeToString(batchIdentifierHash[:]), err)
		agg.telemetry.LogTaskError(batchMerkleRoot, err)
		return nil, err
	}

	agg.walletMutex.Unlock()
	agg.logger.Infof("- Unlocked Wallet Resources: Sending aggregated response for batch %s", hex.EncodeToString(batchIdentifierHash[:]))

	agg.metrics.IncAggregatedResponses()

	return receipt, nil
}

func (agg *Aggregator) AddNewTask(batchMerkleRoot [32]byte, senderAddress [20]byte, taskCreatedBlock uint32) {
	agg.telemetry.InitNewTrace(batchMerkleRoot)
	batchIdentifier := append(batchMerkleRoot[:], senderAddress[:]...)
	var batchIdentifierHash = *(*[32]byte)(crypto.Keccak256(batchIdentifier))

	agg.AggregatorConfig.BaseConfig.Logger.Info("Adding new task",
		"Batch merkle root", "0x"+hex.EncodeToString(batchMerkleRoot[:]),
		"Sender Address", "0x"+hex.EncodeToString(senderAddress[:]),
		"batchIdentifierHash", "0x"+hex.EncodeToString(batchIdentifierHash[:]))

	batchData, err := agg.batchDataCache.PopulateBatchData(batchMerkleRoot, senderAddress, batchIdentifierHash, taskCreatedBlock)
	if err != nil {
		agg.logger.Error("Failed to add task", "error", err)
		return
	}

	quorumNums := eigentypes.QuorumNums{eigentypes.QuorumNum(QUORUM_NUMBER)}
	quorumThresholdPercentages := eigentypes.QuorumThresholdPercentages{eigentypes.QuorumThresholdPercentage(QUORUM_THRESHOLD)}

	err = agg.blsAggregationService.InitializeNewTask(batchData.Index, taskCreatedBlock, quorumNums, quorumThresholdPercentages, agg.AggregatorConfig.Aggregator.BlsServiceTaskTimeout)
	if err != nil {
		agg.logger.Fatalf("BLS aggregation service error when initializing new task: %s", err)
	}

	agg.metrics.IncAggregatorReceivedTasks()
}

func (agg *Aggregator) GetTaskDataByIndex(idx uint32) (BatchData, bool) {
	return agg.batchDataCache.GetTaskDataByIndex(idx)
}

func (agg *Aggregator) GetTaskDataByIdentifierHash(idHash [32]byte) (BatchData, bool) {
	return agg.batchDataCache.GetTaskDataByIdentifierHash(idHash)
}

// Long-lived goroutine that periodically checks and removes old Tasks from stored Maps
// It runs every GarbageCollectorPeriod and removes all tasks older than GarbageCollectorTasksAge
// This was added because each task occupies memory in the maps, and we need to free it to avoid a memory leak
func (agg *Aggregator) ClearTasksFromMaps() {
	defer func() {
		err := recover() //stops panics
		if err != nil {
			agg.logger.Error("ClearTasksFromMaps Recovered from panic", "err", err)
		}
	}()

	agg.AggregatorConfig.BaseConfig.Logger.Info(fmt.Sprintf("- Removing finalized Task Infos from Maps every %v", agg.AggregatorConfig.Aggregator.GarbageCollectorPeriod))

	for {
		time.Sleep(agg.AggregatorConfig.Aggregator.GarbageCollectorPeriod)

		agg.AggregatorConfig.BaseConfig.Logger.Info("Cleaning finalized tasks from maps")
		oldTaskIdHash, err := agg.avsReader.GetOldTaskHash(agg.AggregatorConfig.Aggregator.GarbageCollectorTasksAge, agg.AggregatorConfig.Aggregator.GarbageCollectorTasksInterval)
		if err != nil {
			agg.logger.Error("Error getting old task hash, skipping this garbage collect", "err", err)
			continue // Retry in the next iteration
		}
		if oldTaskIdHash == nil {
			agg.logger.Warn("No old tasks found")
			continue // Retry in the next iteration
		}
		agg.batchDataCache.DeleteOlderByIdentifierHash(*oldTaskIdHash)
	}
}

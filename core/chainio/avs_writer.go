package chainio

import (
	"context"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/avsregistry"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/signer"
	"github.com/ethereum/go-ethereum/common"
	servicemanager "github.com/yetanotherco/aligned_layer/contracts/bindings/AlignedLayerServiceManager"
	"github.com/yetanotherco/aligned_layer/core/config"
	"github.com/yetanotherco/aligned_layer/core/utils"
)

type AvsWriter struct {
	*avsregistry.ChainWriter
	AvsContractBindings *AvsServiceBindings
	logger              logging.Logger
	Signer              signer.Signer
	Client              eth.Client
}

func NewAvsWriterFromConfig(baseConfig *config.BaseConfig, ecdsaConfig *config.EcdsaConfig) (*AvsWriter, error) {

	buildAllConfig := clients.BuildAllConfig{
		EthHttpUrl:                 baseConfig.EthRpcUrl,
		EthWsUrl:                   baseConfig.EthWsUrl,
		RegistryCoordinatorAddr:    baseConfig.AlignedLayerDeploymentConfig.AlignedLayerRegistryCoordinatorAddr.String(),
		OperatorStateRetrieverAddr: baseConfig.AlignedLayerDeploymentConfig.AlignedLayerOperatorStateRetrieverAddr.String(),
		AvsName:                    "AlignedLayer",
		PromMetricsIpPortAddress:   baseConfig.EigenMetricsIpPortAddress,
	}

	clients, err := clients.BuildAll(buildAllConfig, ecdsaConfig.PrivateKey, baseConfig.Logger)

	if err != nil {
		baseConfig.Logger.Error("Cannot build signer config", "err", err)
		return nil, err
	}

	avsServiceBindings, err := NewAvsServiceBindings(baseConfig.AlignedLayerDeploymentConfig.AlignedLayerServiceManagerAddr, baseConfig.AlignedLayerDeploymentConfig.AlignedLayerOperatorStateRetrieverAddr, baseConfig.EthRpcClient, baseConfig.EthRpcClientFallback, baseConfig.Logger)

	if err != nil {
		baseConfig.Logger.Error("Cannot create avs service bindings", "err", err)
		return nil, err
	}

	privateKeySigner, err := signer.NewPrivateKeySigner(ecdsaConfig.PrivateKey, baseConfig.ChainId)
	if err != nil {
		baseConfig.Logger.Error("Cannot create signer", "err", err)
		return nil, err
	}

	chainWriter := clients.AvsRegistryChainWriter

	return &AvsWriter{
		ChainWriter:         chainWriter,
		AvsContractBindings: avsServiceBindings,
		logger:              baseConfig.Logger,
		Signer:              privateKeySigner,
		Client:              baseConfig.EthRpcClient,
	}, nil
}

func (w *AvsWriter) SendTask(context context.Context, batchMerkleRoot [32]byte, batchDataPointer string) error {

	txOpts := w.Signer.GetTxOpts()

	tx, err := w.AvsContractBindings.ServiceManager.CreateNewTask(
		txOpts,
		batchMerkleRoot,
		batchDataPointer,
	)
	if err != nil {
		w.logger.Error("Error assembling CreateNewTask tx", "err", err)
		return err
	}

	_, err = utils.WaitForTransactionReceipt(w.Client, context, tx.Hash())
	if err != nil {
		return err
	}

	return nil
}

func (w *AvsWriter) SendAggregatedResponse(batchMerkleRoot [32]byte, senderAddress [20]byte, nonSignerStakesAndSignature servicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature) (*common.Hash, error) {
	txOpts := *w.Signer.GetTxOpts()
	txOpts.NoSend = true // simulate the transaction
	tx, err := w.AvsContractBindings.ServiceManager.RespondToTaskV2(&txOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
	if err != nil {
		// Retry with fallback
		tx, err = w.AvsContractBindings.ServiceManagerFallback.RespondToTaskV2(&txOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
		if err != nil {
			return nil, err
		}
	}

	// Send the transaction
	txOpts.NoSend = false
	txOpts.GasLimit = tx.Gas() * 110 / 100 // Add 10% to the gas limit
	tx, err = w.AvsContractBindings.ServiceManager.RespondToTaskV2(&txOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
	if err != nil {
		// Retry with fallback
		tx, err = w.AvsContractBindings.ServiceManagerFallback.RespondToTaskV2(&txOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
		if err != nil {
			return nil, err
		}
	}

	txHash := tx.Hash()

	return &txHash, nil
}

// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractAlignedLayerServiceManager

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BN254G1Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G1Point struct {
	X *big.Int
	Y *big.Int
}

// BN254G2Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

// IBLSSignatureCheckerNonSignerStakesAndSignature is an auto generated low-level Go binding around an user-defined struct.
type IBLSSignatureCheckerNonSignerStakesAndSignature struct {
	NonSignerQuorumBitmapIndices []uint32
	NonSignerPubkeys             []BN254G1Point
	QuorumApks                   []BN254G1Point
	ApkG2                        BN254G2Point
	Sigma                        BN254G1Point
	QuorumApkIndices             []uint32
	TotalStakeIndices            []uint32
	NonSignerStakeIndices        [][]uint32
}

// IBLSSignatureCheckerQuorumStakeTotals is an auto generated low-level Go binding around an user-defined struct.
type IBLSSignatureCheckerQuorumStakeTotals struct {
	SignedStakeForQuorum []*big.Int
	TotalStakeForQuorum  []*big.Int
}

// IRewardsCoordinatorRewardsSubmission is an auto generated low-level Go binding around an user-defined struct.
type IRewardsCoordinatorRewardsSubmission struct {
	StrategiesAndMultipliers []IRewardsCoordinatorStrategyAndMultiplier
	Token                    common.Address
	Amount                   *big.Int
	StartTimestamp           uint32
	Duration                 uint32
}

// IRewardsCoordinatorStrategyAndMultiplier is an auto generated low-level Go binding around an user-defined struct.
type IRewardsCoordinatorStrategyAndMultiplier struct {
	Strategy   common.Address
	Multiplier *big.Int
}

// ISignatureUtilsSignatureWithSaltAndExpiry is an auto generated low-level Go binding around an user-defined struct.
type ISignatureUtilsSignatureWithSaltAndExpiry struct {
	Signature []byte
	Salt      [32]byte
	Expiry    *big.Int
}

// ContractAlignedLayerServiceManagerMetaData contains all meta data concerning the ContractAlignedLayerServiceManager contract.
var ContractAlignedLayerServiceManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"__avsDirectory\",\"type\":\"address\",\"internalType\":\"contractIAVSDirectory\"},{\"name\":\"__rewardsCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRewardsCoordinator\"},{\"name\":\"__registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"__stakeRegistry\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"alignedAggregator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"avsDirectory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchersBalances\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchesState\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"taskCreatedBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"responded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"respondToTaskFeeLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blsApkRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBLSApkRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkPublicInput\",\"inputs\":[{\"name\":\"publicInput\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"checkSignatures\",\"inputs\":[{\"name\":\"msgHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.QuorumStakeTotals\",\"components\":[{\"name\":\"signedStakeForQuorum\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"},{\"name\":\"totalStakeForQuorum\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"}]},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createAVSRewardsSubmission\",\"inputs\":[{\"name\":\"rewardsSubmissions\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.RewardsSubmission[]\",\"components\":[{\"name\":\"strategiesAndMultipliers\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.StrategyAndMultiplier[]\",\"components\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"duration\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createNewTask\",\"inputs\":[{\"name\":\"batchMerkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"batchDataPointer\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"respondToTaskFeeLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"delegation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIDelegationManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"depositToBatcher\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"deregisterOperatorFromAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"disableVerifier\",\"inputs\":[{\"name\":\"verifierIdx\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"disabledVerifiers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"enableVerifier\",\"inputs\":[{\"name\":\"verifierIdx\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getOperatorRestakedStrategies\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRestakeableStrategies\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_rewardsInitiator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_alignedAggregator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_pauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"_initialPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initializeAggregator\",\"inputs\":[{\"name\":\"_alignedAggregator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initializePauser\",\"inputs\":[{\"name\":\"_pauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"_initialPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isVerifierDisabled\",\"inputs\":[{\"name\":\"verifierIdx\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseAll\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauserRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOperatorToAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"respondToTaskV2\",\"inputs\":[{\"name\":\"batchMerkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"senderAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rewardsInitiator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setAggregator\",\"inputs\":[{\"name\":\"_alignedAggregator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDisabledVerifiers\",\"inputs\":[{\"name\":\"bitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPauserRegistry\",\"inputs\":[{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRewardsInitiator\",\"inputs\":[{\"name\":\"newRewardsInitiator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStaleStakesForbidden\",\"inputs\":[{\"name\":\"value\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"staleStakesForbidden\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"trySignatureAndApkVerification\",\"inputs\":[{\"name\":\"msgHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"apk\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"pairingSuccessful\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"siganatureIsValid\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateAVSMetadataURI\",\"inputs\":[{\"name\":\"_metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyBatchInclusion\",\"inputs\":[{\"name\":\"proofCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"pubInputCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"provingSystemAuxDataCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proofGeneratorAddr\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"},{\"name\":\"batchMerkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"merkleProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"verificationDataBatchIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"senderAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBatchInclusion\",\"inputs\":[{\"name\":\"proofCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"pubInputCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"provingSystemAuxDataCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proofGeneratorAddr\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"},{\"name\":\"batchMerkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"merkleProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"verificationDataBatchIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BatchVerified\",\"inputs\":[{\"name\":\"batchMerkleRoot\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"senderAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BatcherBalanceUpdated\",\"inputs\":[{\"name\":\"batcher\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewBatchV2\",\"inputs\":[{\"name\":\"batchMerkleRoot\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"senderAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"taskCreatedBlock\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"batchDataPointer\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewBatchV3\",\"inputs\":[{\"name\":\"batchMerkleRoot\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"senderAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"taskCreatedBlock\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"batchDataPointer\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"respondToTaskFeeLimit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauserRegistrySet\",\"inputs\":[{\"name\":\"pauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsInitiatorUpdated\",\"inputs\":[{\"name\":\"prevRewardsInitiator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newRewardsInitiator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleStakesForbiddenUpdate\",\"inputs\":[{\"name\":\"value\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VerifierDisabled\",\"inputs\":[{\"name\":\"verifierIdx\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VerifierEnabled\",\"inputs\":[{\"name\":\"verifierIdx\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"BatchAlreadyResponded\",\"inputs\":[{\"name\":\"batchIdentifierHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BatchAlreadySubmitted\",\"inputs\":[{\"name\":\"batchIdentifierHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BatchDoesNotExist\",\"inputs\":[{\"name\":\"batchIdentifierHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InsufficientFunds\",\"inputs\":[{\"name\":\"batcher\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"required\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[{\"name\":\"param\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidDepositAmount\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidQuorumThreshold\",\"inputs\":[{\"name\":\"signedStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requiredStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SenderIsNotAggregator\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"alignedAggregator\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x6101806040523480156200001257600080fd5b50604051620063e9380380620063e9833981016040819052620000359162000419565b6001600160a01b0380851660805280841660a05280831660c052811660e05281848482846200006362000341565b50505050806001600160a01b0316610100816001600160a01b031681525050806001600160a01b031663683048356040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000c1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000e7919062000481565b6001600160a01b0316610120816001600160a01b031681525050806001600160a01b0316635df459466040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000140573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000166919062000481565b6001600160a01b0316610140816001600160a01b031681525050610120516001600160a01b031663df5cf7236040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001c2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001e8919062000481565b6001600160a01b0390811661016052851690506200023d57604051630b0f5aa160e11b815260206004820152600c60248201526b6176734469726563746f727960a01b60448201526064015b60405180910390fd5b6001600160a01b0383166200028b57604051630b0f5aa160e11b81526020600482015260126024820152713932bbb0b93239a1b7b7b93234b730ba37b960711b604482015260640162000234565b6001600160a01b038216620002e457604051630b0f5aa160e11b815260206004820152601360248201527f7265676973747279436f6f7264696e61746f7200000000000000000000000000604482015260640162000234565b6001600160a01b0381166200032d57604051630b0f5aa160e11b815260206004820152600d60248201526c7374616b65526567697374727960981b604482015260640162000234565b6200033762000341565b50505050620004a8565b600054610100900460ff1615620003ab5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840162000234565b60005460ff9081161015620003fe576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6001600160a01b03811681146200041657600080fd5b50565b600080600080608085870312156200043057600080fd5b84516200043d8162000400565b6020860151909450620004508162000400565b6040860151909350620004638162000400565b6060860151909250620004768162000400565b939692955090935050565b6000602082840312156200049457600080fd5b8151620004a18162000400565b9392505050565b60805160a05160c05160e05161010051610120516101405161016051615e33620005b6600039600081816108390152611d8c01526000818161056a0152611f9f01526000818161059e0152818161218c015261237c0152600081816106050152818161159201528181611a5201528181611bf90152611e400152600081816112ca0152818161141b015281816114b201528181613045015281816131be015261325d0152600081816110f10152818161118001528181611200015281816127dc015281816128a801528181612f80015261311901526000818161397f01528181613a3b0152613b1e0152600081816105cf015281816128300152818161290401526129830152615e336000f3fe60806040526004361061028c5760003560e01c8063800fb61f1161015a578063df5cf723116100c1578063f9120af61161007a578063f9120af6146108f3578063fa534dc014610913578063fabc1cbc14610933578063fc299dee14610953578063fce36c7d14610973578063fd4c3b7c1461099357600080fd5b8063df5cf72314610827578063e481af9d1461085b578063ea5ca34b14610870578063f2fde38b14610886578063f474b520146108a6578063f7013ef6146108d357600080fd5b8063a98fb35511610113578063a98fb35514610730578063ab21739a14610750578063b099627e14610770578063b753645e146107da578063b98d0908146107fa578063d66eaabd1461081457600080fd5b8063800fb61f14610672578063886f1195146106925780638da5cb5b146106b257806395c6d604146106d05780639926ee7d146106f0578063a364f4da1461071057600080fd5b80634223d551116101fe5780635df45946116101b75780635df4594614610558578063683048351461058c5780636b3aa72e146105c05780636d14a987146105f357806370a0823114610627578063715018a61461065d57600080fd5b80634223d5511461047b5780634a5bf6321461048e5780634ae07c37146104c6578063595c6a67146104f45780635ac86ab7146105095780635c975abb1461053957600080fd5b806318daeeaf1161025057806318daeeaf146103ae5780632585b25b146103ce5780632e1a7d4d146103ee57806333cfb7b71461040e5780633bc28c8c1461043b578063416c7e5e1461045b57600080fd5b806306045a91146102d357806310d67a2f14610308578063136439dd14610328578063137122b514610348578063171f1d5b1461037757600080fd5b366102ce5760fc546005906020908116036102c25760405162461bcd60e51b81526004016102b990614b81565b60405180910390fd5b6102cc33346109b3565b005b600080fd5b3480156102df57600080fd5b506102f36102ee366004614cf4565b610a43565b60405190151581526020015b60405180910390f35b34801561031457600080fd5b506102cc610323366004614d86565b610b65565b34801561033457600080fd5b506102cc610343366004614da3565b610c18565b34801561035457600080fd5b506102f3610363366004614dcb565b60cc54600160ff9092169190911b16151590565b34801561038357600080fd5b50610397610392366004614ea8565b610d57565b6040805192151583529015156020830152016102ff565b3480156103ba57600080fd5b506102cc6103c9366004614dcb565b610ee1565b3480156103da57600080fd5b506102cc6103e9366004614ef9565b610f29565b3480156103fa57600080fd5b506102cc610409366004614da3565b610fcb565b34801561041a57600080fd5b5061042e610429366004614d86565b6110cc565b6040516102ff9190614f25565b34801561044757600080fd5b506102cc610456366004614d86565b61157f565b34801561046757600080fd5b506102cc610476366004614f80565b611590565b6102cc610489366004614d86565b6116c7565b34801561049a57600080fd5b5060cb546104ae906001600160a01b031681565b6040516001600160a01b0390911681526020016102ff565b3480156104d257600080fd5b506104e66104e136600461525b565b6116fd565b6040516102ff9291906152f6565b34801561050057600080fd5b506102cc612631565b34801561051557600080fd5b506102f3610524366004614dcb565b60fc54600160ff9092169190911b9081161490565b34801561054557600080fd5b5060fc545b6040519081526020016102ff565b34801561056457600080fd5b506104ae7f000000000000000000000000000000000000000000000000000000000000000081565b34801561059857600080fd5b506104ae7f000000000000000000000000000000000000000000000000000000000000000081565b3480156105cc57600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006104ae565b3480156105ff57600080fd5b506104ae7f000000000000000000000000000000000000000000000000000000000000000081565b34801561063357600080fd5b5061054a610642366004614d86565b6001600160a01b0316600090815260ca602052604090205490565b34801561066957600080fd5b506102cc6126f8565b34801561067e57600080fd5b506102cc61068d366004614d86565b61270c565b34801561069e57600080fd5b5060fb546104ae906001600160a01b031681565b3480156106be57600080fd5b506033546001600160a01b03166104ae565b3480156106dc57600080fd5b506102f36106eb366004615387565b6127ac565b3480156106fc57600080fd5b506102cc61070b3660046153d2565b6127d1565b34801561071c57600080fd5b506102cc61072b366004614d86565b61289d565b34801561073c57600080fd5b506102cc61074b36600461547d565b612964565b34801561075c57600080fd5b506102cc61076b3660046154cd565b6129b8565b34801561077c57600080fd5b506107b861078b366004614da3565b60c9602052600090815260409020805460019091015463ffffffff821691640100000000900460ff169083565b6040805163ffffffff90941684529115156020840152908201526060016102ff565b3480156107e657600080fd5b506102cc6107f5366004614da3565b612d8a565b34801561080657600080fd5b506097546102f39060ff1681565b6102cc6108223660046154f4565b612d97565b34801561083357600080fd5b506104ae7f000000000000000000000000000000000000000000000000000000000000000081565b34801561086757600080fd5b5061042e612f7a565b34801561087c57600080fd5b5061054a60cc5481565b34801561089257600080fd5b506102cc6108a1366004614d86565b613326565b3480156108b257600080fd5b5061054a6108c1366004614d86565b60ca6020526000908152604090205481565b3480156108df57600080fd5b506102cc6108ee366004615546565b61339c565b3480156108ff57600080fd5b506102cc61090e366004614d86565b613575565b34801561091f57600080fd5b506102f361092e3660046155aa565b61359f565b34801561093f57600080fd5b506102cc61094e366004614da3565b61364a565b34801561095f57600080fd5b506065546104ae906001600160a01b031681565b34801561097f57600080fd5b506102cc61098e366004615627565b6137a6565b34801561099f57600080fd5b506102cc6109ae366004614dcb565b613b55565b806000036109d757604051632097692160e11b8152600481018290526024016102b9565b6001600160a01b038216600090815260ca6020526040812080548392906109ff9084906156b1565b90915550506001600160a01b038216600081815260ca6020908152604091829020549151918252600080516020615dbe833981519152910160405180910390a25050565b60fc54600090600290600490811603610a6e5760405162461bcd60e51b81526004016102b990614b81565b60006001600160a01b038416610a85575085610ab1565b8684604051602001610a989291906156c4565b6040516020818303038152906040528051906020012090505b600081815260c9602052604081205463ffffffff169003610ad6576000925050610b58565b600081815260c96020526040902054640100000000900460ff16610afe576000925050610b58565b60408051602081018d90529081018b9052606081018a90526001600160601b03198916608082015260009060940160408051601f1981840301815291905280516020820120909150610b52888a838a613b9c565b94505050505b5098975050505050505050565b60fb60009054906101000a90046001600160a01b03166001600160a01b031663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610bb8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bdc91906156df565b6001600160a01b0316336001600160a01b031614610c0c5760405162461bcd60e51b81526004016102b9906156fc565b610c1581613bb4565b50565b60fb5460405163237dfb4760e11b81523360048201526001600160a01b03909116906346fbf68e90602401602060405180830381865afa158015610c60573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c849190615746565b610ca05760405162461bcd60e51b81526004016102b990615763565b60fc5481811614610d195760405162461bcd60e51b815260206004820152603860248201527f5061757361626c652e70617573653a20696e76616c696420617474656d70742060448201527f746f20756e70617573652066756e6374696f6e616c697479000000000000000060648201526084016102b9565b60fc81905560405181815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d906020015b60405180910390a250565b60008060007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000187876000015188602001518860000151600060028110610d9f57610d9f6157ab565b60200201518951600160200201518a60200151600060028110610dc457610dc46157ab565b60200201518b60200151600160028110610de057610de06157ab565b602090810291909101518c518d830151604051610e3d9a99989796959401988952602089019790975260408801959095526060870193909352608086019190915260a085015260c084015260e08301526101008201526101200190565b6040516020818303038152906040528051906020012060001c610e6091906157c1565b9050610ed3610e79610e728884613cab565b8690613d3c565b610e81613dd1565b610ec9610eba85610eb4604080518082018252600080825260209182015281518083019092526001825260029082015290565b90613cab565b610ec38c613e91565b90613d3c565b886201d4c0613f20565b909890975095505050505050565b610ee961413a565b60cc8054600160ff841690811b199091169091556040517f5f52704e8e0190647930ccde0e43e14e89902d7d8c49c5f9e2544029f45ec12a90600090a250565b600054600390610100900460ff16158015610f4b575060005460ff8083169116105b610f675760405162461bcd60e51b81526004016102b9906157e3565b6000805461ffff191660ff831617610100179055610f858383614194565b6000805461ff001916905560405160ff821681527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a1505050565b60fc54600390600890811603610ff35760405162461bcd60e51b81526004016102b990614b81565b33600090815260ca60205260409020548211156110445733600081815260ca602052604090819020549051632e2a182f60e11b815260048101929092526024820184905260448201526064016102b9565b33600090815260ca602052604081208054849290611063908490615831565b909155505033600081815260ca6020908152604091829020549151918252600080516020615dbe833981519152910160405180910390a2604051339083156108fc029084906000818181858888f193505050501580156110c7573d6000803e3d6000fd5b505050565b6040516309aa152760e11b81526001600160a01b0382811660048301526060916000917f000000000000000000000000000000000000000000000000000000000000000016906313542a4e90602401602060405180830381865afa158015611138573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061115c9190615844565b60405163871ef04960e01b8152600481018290529091506000906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063871ef04990602401602060405180830381865afa1580156111c7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111eb919061585d565b90506001600160c01b038116158061128557507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561125c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112809190615886565b60ff16155b156112a55760408051600080825260208201909252905b50949350505050565b60006112b9826001600160c01b031661427a565b90506000805b8251811015611385577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316633ca5a5f5848381518110611309576113096157ab565b01602001516040516001600160e01b031960e084901b16815260f89190911c6004820152602401602060405180830381865afa15801561134d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113719190615844565b61137b90836156b1565b91506001016112bf565b506000816001600160401b038111156113a0576113a0614bd0565b6040519080825280602002602001820160405280156113c9578160200160208202803683370190505b5090506000805b84518110156115725760008582815181106113ed576113ed6157ab565b0160200151604051633ca5a5f560e01b815260f89190911c6004820181905291506000906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690633ca5a5f590602401602060405180830381865afa158015611462573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114869190615844565b905060005b81811015611567576040516356e4026d60e11b815260ff84166004820152602481018290527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063adc804da906044016040805180830381865afa158015611500573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061152491906158b8565b6000015186868151811061153a5761153a6157ab565b6001600160a01b03909216602092830291909101909101528461155c816158f9565b95505060010161148b565b5050506001016113d0565b5090979650505050505050565b61158761413a565b610c158161433c565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156115ee573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061161291906156df565b6001600160a01b0316336001600160a01b0316146116be5760405162461bcd60e51b815260206004820152605c60248201527f424c535369676e6174757265436865636b65722e6f6e6c79436f6f7264696e6160448201527f746f724f776e65723a2063616c6c6572206973206e6f7420746865206f776e6560648201527f72206f6620746865207265676973747279436f6f7264696e61746f7200000000608482015260a4016102b9565b610c15816143a5565b60fc546004906010908116036116ef5760405162461bcd60e51b81526004016102b990614b81565b6116f982346109b3565b5050565b6040805180820190915260608082526020820152600082604001515160405180604001604052806001815260200160008152505114801561175957508260a0015151604051806040016040528060018152602001600081525051145b801561178057508260c0015151604051806040016040528060018152602001600081525051145b80156117a757508260e0015151604051806040016040528060018152602001600081525051145b6118115760405162461bcd60e51b81526020600482015260416024820152600080516020615dde83398151915260448201527f7265733a20696e7075742071756f72756d206c656e677468206d69736d6174636064820152600d60fb1b608482015260a4016102b9565b825151602084015151146118895760405162461bcd60e51b815260206004820152604460248201819052600080516020615dde833981519152908201527f7265733a20696e707574206e6f6e7369676e6572206c656e677468206d69736d6064820152630c2e8c6d60e31b608482015260a4016102b9565b4363ffffffff168463ffffffff16106118f85760405162461bcd60e51b815260206004820152603c6024820152600080516020615dde83398151915260448201527f7265733a20696e76616c6964207265666572656e636520626c6f636b0000000060648201526084016102b9565b60408051808201825260008082526020808301829052835180850185526060808252818301528451808601865260018082529083019390935284518381528086019095529293919082810190803683370190505060208281019190915260408051808201825260018082526000919093015280518281528082019091529081602001602082028036833701905050815260408051808201909152606080825260208201528560200151516001600160401b038111156119b9576119b9614bd0565b6040519080825280602002602001820160405280156119e2578160200160208202803683370190505b5081526020860151516001600160401b03811115611a0257611a02614bd0565b604051908082528060200260200182016040528015611a2b578160200160208202803683370190505b5081602001819052506000611ad760405180604001604052806001815260200160008152507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611aae573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611ad29190615886565b6143ec565b905060005b876020015151811015611d6857611b2188602001518281518110611b0257611b026157ab565b6020026020010151805160009081526020918201519091526040902090565b83602001518281518110611b3757611b376157ab565b60209081029190910101528015611bf7576020830151611b58600183615831565b81518110611b6857611b686157ab565b602002602001015160001c83602001518281518110611b8957611b896157ab565b602002602001015160001c11611bf7576040805162461bcd60e51b8152602060048201526024810191909152600080516020615dde83398151915260448201527f7265733a206e6f6e5369676e65725075626b657973206e6f7420736f7274656460648201526084016102b9565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166304ec635184602001518381518110611c3c57611c3c6157ab565b60200260200101518b8b600001518581518110611c5b57611c5b6157ab565b60200260200101516040518463ffffffff1660e01b8152600401611c989392919092835263ffffffff918216602084015216604082015260600190565b602060405180830381865afa158015611cb5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611cd9919061585d565b6001600160c01b031683600001518281518110611cf857611cf86157ab565b602002602001018181525050611d5e610e72611d328486600001518581518110611d2457611d246157ab565b60200260200101511661447f565b8a602001518481518110611d4857611d486157ab565b60200260200101516144aa90919063ffffffff16565b9450600101611adc565b5050611d738361458d565b60975490935060ff16600081611d8a576000611e0c565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663c448feb86040518163ffffffff1660e01b8152600401602060405180830381865afa158015611de8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e0c9190615844565b905060005b604051806040016040528060018152602001600081525051811015612502578215611f9d578963ffffffff16827f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663249a0c4260405180604001604052806001815260200160008152508581518110611e9557611e956157ab565b01602001516040516001600160e01b031960e084901b16815260f89190911c6004820152602401602060405180830381865afa158015611ed9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611efd9190615844565b611f0791906156b1565b11611f9d5760405162461bcd60e51b81526020600482015260666024820152600080516020615dde83398151915260448201527f7265733a205374616b6552656769737472792075706461746573206d7573742060648201527f62652077697468696e207769746864726177616c44656c6179426c6f636b732060848201526577696e646f7760d01b60a482015260c4016102b9565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166368bccaac60405180604001604052806001815260200160008152508381518110611ff457611ff46157ab565b602001015160f81c60f81b60f81c8c8c60a001518581518110612019576120196157ab565b60209081029190910101516040516001600160e01b031960e086901b16815260ff909316600484015263ffffffff9182166024840152166044820152606401602060405180830381865afa158015612075573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906120999190615912565b6001600160401b0319166120bc8a604001518381518110611b0257611b026157ab565b67ffffffffffffffff1916146121585760405162461bcd60e51b81526020600482015260616024820152600080516020615dde83398151915260448201527f7265733a2071756f72756d41706b206861736820696e2073746f72616765206460648201527f6f6573206e6f74206d617463682070726f76696465642071756f72756d2061706084820152606b60f81b60a482015260c4016102b9565b61218889604001518281518110612171576121716157ab565b602002602001015187613d3c90919063ffffffff16565b95507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663c8294c56604051806040016040528060018152602001600081525083815181106121e1576121e16157ab565b602001015160f81c60f81b60f81c8c8c60c001518581518110612206576122066157ab565b60209081029190910101516040516001600160e01b031960e086901b16815260ff909316600484015263ffffffff9182166024840152166044820152606401602060405180830381865afa158015612262573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612286919061593d565b8560200151828151811061229c5761229c6157ab565b6001600160601b039092166020928302919091018201528501518051829081106122c8576122c86157ab565b6020026020010151856000015182815181106122e6576122e66157ab565b60200260200101906001600160601b031690816001600160601b0316815250506000805b8a60200151518110156124f85761237586600001518281518110612330576123306157ab565b602002602001015160405180604001604052806001815260200160008152508581518110612360576123606157ab565b016020015160f81c60ff161c60019081161490565b156124f0577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663f2be94ae604051806040016040528060018152602001600081525085815181106123d1576123d16157ab565b602001015160f81c60f81b60f81c8e896020015185815181106123f6576123f66157ab565b60200260200101518f60e001518881518110612414576124146157ab565b6020026020010151878151811061242d5761242d6157ab565b60209081029190910101516040516001600160e01b031960e087901b16815260ff909416600485015263ffffffff92831660248501526044840191909152166064820152608401602060405180830381865afa158015612491573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124b5919061593d565b87518051859081106124c9576124c96157ab565b602002602001018181516124dd919061595a565b6001600160601b03169052506001909101905b60010161230a565b5050600101611e11565b50505060008061251c8a868a606001518b60800151610d57565b915091508161258d5760405162461bcd60e51b81526020600482015260436024820152600080516020615dde83398151915260448201527f7265733a2070616972696e6720707265636f6d70696c652063616c6c206661696064820152621b195960ea1b608482015260a4016102b9565b806125ee5760405162461bcd60e51b81526020600482015260396024820152600080516020615dde83398151915260448201527f7265733a207369676e617475726520697320696e76616c69640000000000000060648201526084016102b9565b50506000878260200151604051602001612609929190615981565b60408051808303601f1901815291905280516020909101209299929850919650505050505050565b60fb5460405163237dfb4760e11b81523360048201526001600160a01b03909116906346fbf68e90602401602060405180830381865afa158015612679573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061269d9190615746565b6126b95760405162461bcd60e51b81526004016102b990615763565b60001960fc81905560405190815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d9060200160405180910390a2565b61270061413a565b61270a6000614628565b565b600054600290610100900460ff1615801561272e575060005460ff8083169116105b61274a5760405162461bcd60e51b81526004016102b9906157e3565b6000805461ffff191660ff83161761010017905561276782613575565b6000805461ff001916905560405160ff821681527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15050565b60008184846040516127bf9291906159c9565b60405180910390201490509392505050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146128195760405162461bcd60e51b81526004016102b9906159d9565b604051639926ee7d60e01b81526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690639926ee7d906128679085908590600401615a97565b600060405180830381600087803b15801561288157600080fd5b505af1158015612895573d6000803e3d6000fd5b505050505050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146128e55760405162461bcd60e51b81526004016102b9906159d9565b6040516351b27a6d60e11b81526001600160a01b0382811660048301527f0000000000000000000000000000000000000000000000000000000000000000169063a364f4da906024015b600060405180830381600087803b15801561294957600080fd5b505af115801561295d573d6000803e3d6000fd5b5050505050565b61296c61413a565b60405163a98fb35560e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063a98fb3559061292f908490600401615ae2565b60cb546001600160a01b031633146129f85760cb54604051632cbe419560e01b81523360048201526001600160a01b0390911660248201526044016102b9565b60fc54600190600290811603612a205760405162461bcd60e51b81526004016102b990614b81565b60005a905060008585604051602001612a3a9291906156c4565b60408051601f198184030181529181528151602092830120600081815260c990935290822080549193509163ffffffff9091169003612a8f576040516311cb69a760e11b8152600481018390526024016102b9565b8054640100000000900460ff1615612abd57604051634e78d7f960e11b8152600481018390526024016102b9565b805464ff00000000191664010000000017815560018101546001600160a01b038716600090815260ca60205260409020541015612b405760018101546001600160a01b038716600081815260ca602052604090819020549051632e2a182f60e11b81526004810192909252602482019290925260448101919091526064016102b9565b8054600090612b5790849063ffffffff16886116fd565b509050604360ff168160200151600081518110612b7657612b766157ab565b6020026020010151612b889190615af5565b6001600160601b031660648260000151600081518110612baa57612baa6157ab565b60200260200101516001600160601b0316612bc59190615b18565b1015612c585760648160000151600081518110612be457612be46157ab565b60200260200101516001600160601b0316612bff9190615b18565b604360ff168260200151600081518110612c1b57612c1b6157ab565b6020026020010151612c2d9190615af5565b60405163530f5c4560e11b815260048101929092526001600160601b031660248201526044016102b9565b6040516001600160a01b038816815288907f8511746b73275e06971968773119b9601fc501d7bdf3824d8754042d148940e29060200160405180910390a260003a5a612ca49087615831565b612cb190620111706156b1565b612cbb9190615b18565b9050600083600101548210612cd4578360010154612cd6565b815b6001600160a01b038a16600090815260ca6020526040812080549293508392909190612d03908490615831565b90915550506001600160a01b038916600081815260ca6020908152604091829020549151918252600080516020615dbe833981519152910160405180910390a260cb546040516001600160a01b039091169082156108fc029083906000818181858888f19350505050158015612d7d573d6000803e3d6000fd5b5050505050505050505050565b612d9261413a565b60cc55565b60fc54600090600190811603612dbf5760405162461bcd60e51b81526004016102b990614b81565b60008533604051602001612dd49291906156c4565b60408051601f198184030181529181528151602092830120600081815260c990935291205490915063ffffffff1615612e2357604051630c40bc4360e21b8152600481018290526024016102b9565b3415612e805733600090815260ca602052604081208054349290612e489084906156b1565b909155505033600081815260ca6020908152604091829020549151918252600080516020615dbe833981519152910160405180910390a25b33600090815260ca6020526040902054831115612ed15733600081815260ca602052604090819020549051632e2a182f60e11b815260048101929092526024820185905260448201526064016102b9565b604080516060810182526000602080830182815263ffffffff4381811686528587018a815288865260c99094529386902085518154935115156401000000000264ffffffffff1990941692169190911791909117815590516001909101559151909188917f8801fc966deb2c8f563a103c35c9e80740585c292cd97518587e6e7927e6af5591612f69913391908b908b908b90615b2f565b60405180910390a250505050505050565b606060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612fdc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130009190615886565b60ff1690508060000361302157505060408051600081526020810190915290565b6000805b828110156130cc57604051633ca5a5f560e01b815260ff821660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690633ca5a5f590602401602060405180830381865afa158015613094573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906130b89190615844565b6130c290836156b1565b9150600101613025565b506000816001600160401b038111156130e7576130e7614bd0565b604051908082528060200260200182016040528015613110578160200160208202803683370190505b5090506000805b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613175573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906131999190615886565b60ff1681101561331c57604051633ca5a5f560e01b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690633ca5a5f590602401602060405180830381865afa15801561320d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906132319190615844565b905060005b81811015613312576040516356e4026d60e11b815260ff84166004820152602481018290527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063adc804da906044016040805180830381865afa1580156132ab573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906132cf91906158b8565b600001518585815181106132e5576132e56157ab565b6001600160a01b039092166020928302919091019091015283613307816158f9565b945050600101613236565b5050600101613117565b5090949350505050565b61332e61413a565b6001600160a01b0381166133935760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016102b9565b610c1581614628565b600054610100900460ff16158080156133bc5750600054600160ff909116105b806133d65750303b1580156133d6575060005460ff166001145b6133f25760405162461bcd60e51b81526004016102b9906157e3565b6000805460ff191660011790558015613415576000805461ff0019166101001790555b6001600160a01b03861661345b57604051630b0f5aa160e11b815260206004820152600c60248201526b34b734ba34b0b627bbb732b960a11b60448201526064016102b9565b6001600160a01b0385166134a557604051630b0f5aa160e11b815260206004820152601060248201526f3932bbb0b93239a4b734ba34b0ba37b960811b60448201526064016102b9565b6001600160a01b0384166134f057604051630b0f5aa160e11b815260206004820152601160248201527030b634b3b732b220b3b3b932b3b0ba37b960791b60448201526064016102b9565b6134fa868661467a565b60cb80546001600160a01b0319166001600160a01b03861617905561351e86614628565b6135288383614194565b8015612895576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a1505050505050565b61357d61413a565b60cb80546001600160a01b0319166001600160a01b0392909216919091179055565b60fc546000906002906004908116036135ca5760405162461bcd60e51b81526004016102b990614b81565b6040516306045a9160e01b815230906306045a91906135fc908c908c908c908c908c908c908c90600090600401615b86565b602060405180830381865afa158015613619573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061363d9190615746565b9998505050505050505050565b60fb60009054906101000a90046001600160a01b03166001600160a01b031663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561369d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906136c191906156df565b6001600160a01b0316336001600160a01b0316146136f15760405162461bcd60e51b81526004016102b9906156fc565b60fc5419811960fc5419161461376f5760405162461bcd60e51b815260206004820152603860248201527f5061757361626c652e756e70617573653a20696e76616c696420617474656d7060448201527f7420746f2070617573652066756e6374696f6e616c697479000000000000000060648201526084016102b9565b60fc81905560405181815233907f3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c90602001610d4c565b6065546001600160a01b0316331461383b5760405162461bcd60e51b815260206004820152604c60248201527f536572766963654d616e61676572426173652e6f6e6c7952657761726473496e60448201527f69746961746f723a2063616c6c6572206973206e6f742074686520726577617260648201526b32399034b734ba34b0ba37b960a11b608482015260a4016102b9565b60005b81811015613b0657828282818110613858576138586157ab565b905060200281019061386a9190615be8565b61387b906040810190602001614d86565b6001600160a01b03166323b872dd333086868681811061389d5761389d6157ab565b90506020028101906138af9190615be8565b604080516001600160e01b031960e087901b1681526001600160a01b039485166004820152939092166024840152013560448201526064016020604051808303816000875af1158015613906573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061392a9190615746565b50600083838381811061393f5761393f6157ab565b90506020028101906139519190615be8565b613962906040810190602001614d86565b604051636eb1769f60e11b81523060048201526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000081166024830152919091169063dd62ed3e90604401602060405180830381865afa1580156139d0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139f49190615844565b9050838383818110613a0857613a086157ab565b9050602002810190613a1a9190615be8565b613a2b906040810190602001614d86565b6001600160a01b031663095ea7b37f000000000000000000000000000000000000000000000000000000000000000083878787818110613a6d57613a6d6157ab565b9050602002810190613a7f9190615be8565b60400135613a8d91906156b1565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af1158015613ad8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613afc9190615746565b505060010161383e565b5060405163fce36c7d60e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063fce36c7d906128679085908590600401615c6e565b613b5d61413a565b60cc8054600160ff841690811b9091179091556040517fec54a85c01b5fc7fb41be0f33eabc56f2981110da8317b9817bc7c718f6d7bfe90600090a250565b600083613baa8685856146f7565b1495945050505050565b6001600160a01b038116613c425760405162461bcd60e51b815260206004820152604960248201527f5061757361626c652e5f73657450617573657252656769737472793a206e657760448201527f50617573657252656769737472792063616e6e6f7420626520746865207a65726064820152686f206164647265737360b81b608482015260a4016102b9565b60fb54604080516001600160a01b03928316815291831660208301527f6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6910160405180910390a160fb80546001600160a01b0319166001600160a01b0392909216919091179055565b6040805180820190915260008082526020820152613cc7614aa7565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa90508080613cf657fe5b5080613d345760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b60448201526064016102b9565b505092915050565b6040805180820190915260008082526020820152613d58614ac5565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa90508080613d9357fe5b5080613d345760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b60448201526064016102b9565b613dd9614ae3565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b604080518082019091526000808252602082015260008080613ec1600080516020615d9e833981519152866157c1565b90505b613ecd816147f4565b9093509150600080516020615d9e8339815191528283098303613f06576040805180820190915290815260208101919091529392505050565b600080516020615d9e833981519152600182089050613ec4565b604080518082018252868152602080820186905282518084019093528683528201849052600091829190613f52614b08565b60005b600281101561410d576000613f6b826006615b18565b9050848260028110613f7f57613f7f6157ab565b60200201515183613f918360006156b1565b600c8110613fa157613fa16157ab565b6020020152848260028110613fb857613fb86157ab565b60200201516020015183826001613fcf91906156b1565b600c8110613fdf57613fdf6157ab565b6020020152838260028110613ff657613ff66157ab565b60200201515151836140098360026156b1565b600c8110614019576140196157ab565b6020020152838260028110614030576140306157ab565b60200201515160016020020151836140498360036156b1565b600c8110614059576140596157ab565b6020020152838260028110614070576140706157ab565b60200201516020015160006002811061408b5761408b6157ab565b60200201518361409c8360046156b1565b600c81106140ac576140ac6157ab565b60200201528382600281106140c3576140c36157ab565b6020020151602001516001600281106140de576140de6157ab565b6020020151836140ef8360056156b1565b600c81106140ff576140ff6157ab565b602002015250600101613f55565b50614116614b27565b60006020826101808560088cfa9151919c9115159b50909950505050505050505050565b6033546001600160a01b0316331461270a5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102b9565b60fb546001600160a01b03161580156141b557506001600160a01b03821615155b6142375760405162461bcd60e51b815260206004820152604760248201527f5061757361626c652e5f696e697469616c697a655061757365723a205f696e6960448201527f7469616c697a6550617573657228292063616e206f6e6c792062652063616c6c6064820152666564206f6e636560c81b608482015260a4016102b9565b60fc81905560405181815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d9060200160405180910390a26116f982613bb4565b60606000806142888461447f565b61ffff166001600160401b038111156142a3576142a3614bd0565b6040519080825280601f01601f1916602001820160405280156142cd576020820181803683370190505b5090506000805b8251821080156142e5575061010081105b1561331c576001811b93508584161561432c578060f81b83838151811061430e5761430e6157ab565b60200101906001600160f81b031916908160001a9053508160010191505b614335816158f9565b90506142d4565b606554604080516001600160a01b03928316815291831660208301527fe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3910160405180910390a1606580546001600160a01b0319166001600160a01b0392909216919091179055565b6097805460ff19168215159081179091556040519081527f40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc9060200160405180910390a150565b6000806143f884614876565b9050808360ff166001901b116144765760405162461bcd60e51b815260206004820152603f60248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206269746d61702065786365656473206d61782076616c75650060648201526084016102b9565b90505b92915050565b6000805b821561447957614494600184615831565b90921691806144a281615d7c565b915050614483565b60408051808201909152600080825260208201526102008261ffff16106145065760405162461bcd60e51b815260206004820152601060248201526f7363616c61722d746f6f2d6c6172676560801b60448201526064016102b9565b8161ffff16600103614519575081614479565b6040805180820190915260008082526020820181905284906001905b8161ffff168661ffff161061458257600161ffff871660ff83161c81169003614565576145628484613d3c565b93505b61456f8384613d3c565b92506201fffe600192831b169101614535565b509195945050505050565b604080518082019091526000808252602082015281511580156145b257506020820151155b156145d0575050604080518082019091526000808252602082015290565b604051806040016040528083600001518152602001600080516020615d9e833981519152846020015161460391906157c1565b61461b90600080516020615d9e833981519152615831565b905292915050565b919050565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166146e55760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b60648201526084016102b9565b6146ee82614628565b6116f98161433c565b60006020845161470791906157c1565b1561478e5760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a4016102b9565b8260205b8551811161129c576147a56002856157c1565b6000036147c9578160005280860151602052604060002091506002840493506147e2565b8086015160005281602052604060002091506002840493505b6147ed6020826156b1565b9050614792565b60008080600080516020615d9e8339815191526003600080516020615d9e83398151915286600080516020615d9e83398151915288890909089050600061486a827f0c19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f52600080516020615d9e8339815191526149fe565b91959194509092505050565b6000610100825111156148ff5760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a4016102b9565b815160000361491057506000919050565b60008083600081518110614926576149266157ab565b0160200151600160f89190911c81901b92505b84518110156149f557848181518110614954576149546157ab565b0160200151600160f89190911c1b91508282116149e95760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a4016102b9565b91811791600101614939565b50909392505050565b600080614a09614b27565b614a11614b45565b602080825281810181905260408201819052606082018890526080820187905260a082018690528260c08360056107d05a03fa92508280614a4e57fe5b5082614a9c5760405162461bcd60e51b815260206004820152601a60248201527f424e3235342e6578704d6f643a2063616c6c206661696c75726500000000000060448201526064016102b9565b505195945050505050565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b6040518060400160405280614af6614b63565b8152602001614b03614b63565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b60208082526019908201527f5061757361626c653a20696e6465782069732070617573656400000000000000604082015260600190565b80356001600160601b03198116811461462357600080fd5b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715614c0857614c08614bd0565b60405290565b60405161010081016001600160401b0381118282101715614c0857614c08614bd0565b604051601f8201601f191681016001600160401b0381118282101715614c5957614c59614bd0565b604052919050565b60006001600160401b03831115614c7a57614c7a614bd0565b614c8d601f8401601f1916602001614c31565b9050828152838383011115614ca157600080fd5b828260208301376000602084830101529392505050565b600082601f830112614cc957600080fd5b614cd883833560208501614c61565b9392505050565b6001600160a01b0381168114610c1557600080fd5b600080600080600080600080610100898b031215614d1157600080fd5b883597506020890135965060408901359550614d2f60608a01614bb8565b94506080890135935060a08901356001600160401b03811115614d5157600080fd5b614d5d8b828c01614cb8565b93505060c0890135915060e0890135614d7581614cdf565b809150509295985092959890939650565b600060208284031215614d9857600080fd5b813561447681614cdf565b600060208284031215614db557600080fd5b5035919050565b60ff81168114610c1557600080fd5b600060208284031215614ddd57600080fd5b813561447681614dbc565b600060408284031215614dfa57600080fd5b614e02614be6565b9050813581526020820135602082015292915050565b600082601f830112614e2957600080fd5b614e31614be6565b806040840185811115614e4357600080fd5b845b81811015614e5d578035845260209384019301614e45565b509095945050505050565b600060808284031215614e7a57600080fd5b614e82614be6565b9050614e8e8383614e18565b8152614e9d8360408401614e18565b602082015292915050565b6000806000806101208587031215614ebf57600080fd5b84359350614ed08660208701614de8565b9250614edf8660608701614e68565b9150614eee8660e08701614de8565b905092959194509250565b60008060408385031215614f0c57600080fd5b8235614f1781614cdf565b946020939093013593505050565b6020808252825182820181905260009190848201906040850190845b81811015614f665783516001600160a01b031683529284019291840191600101614f41565b50909695505050505050565b8015158114610c1557600080fd5b600060208284031215614f9257600080fd5b813561447681614f72565b803563ffffffff8116811461462357600080fd5b60006001600160401b03821115614fca57614fca614bd0565b5060051b60200190565b600082601f830112614fe557600080fd5b81356020614ffa614ff583614fb1565b614c31565b8083825260208201915060208460051b87010193508684111561501c57600080fd5b602086015b8481101561503f5761503281614f9d565b8352918301918301615021565b509695505050505050565b600082601f83011261505b57600080fd5b8135602061506b614ff583614fb1565b8083825260208201915060208460061b87010193508684111561508d57600080fd5b602086015b8481101561503f576150a48882614de8565b835291830191604001615092565b600082601f8301126150c357600080fd5b813560206150d3614ff583614fb1565b82815260059290921b840181019181810190868411156150f257600080fd5b8286015b8481101561503f5780356001600160401b038111156151155760008081fd5b6151238986838b0101614fd4565b8452509183019183016150f6565b6000610180828403121561514457600080fd5b61514c614c0e565b905081356001600160401b038082111561516557600080fd5b61517185838601614fd4565b8352602084013591508082111561518757600080fd5b6151938583860161504a565b602084015260408401359150808211156151ac57600080fd5b6151b88583860161504a565b60408401526151ca8560608601614e68565b60608401526151dc8560e08601614de8565b60808401526101208401359150808211156151f657600080fd5b61520285838601614fd4565b60a084015261014084013591508082111561521c57600080fd5b61522885838601614fd4565b60c084015261016084013591508082111561524257600080fd5b5061524f848285016150b2565b60e08301525092915050565b60008060006060848603121561527057600080fd5b8335925061528060208501614f9d565b915060408401356001600160401b0381111561529b57600080fd5b6152a786828701615131565b9150509250925092565b60008151808452602080850194506020840160005b838110156152eb5781516001600160601b0316875295820195908201906001016152c6565b509495945050505050565b604081526000835160408084015261531160808401826152b1565b90506020850151603f1984830301606085015261532e82826152b1565b925050508260208301529392505050565b60008083601f84011261535157600080fd5b5081356001600160401b0381111561536857600080fd5b60208301915083602082850101111561538057600080fd5b9250929050565b60008060006040848603121561539c57600080fd5b83356001600160401b038111156153b257600080fd5b6153be8682870161533f565b909790965060209590950135949350505050565b600080604083850312156153e557600080fd5b82356153f081614cdf565b915060208301356001600160401b038082111561540c57600080fd5b908401906060828703121561542057600080fd5b60405160608101818110838211171561543b5761543b614bd0565b60405282358281111561544d57600080fd5b61545988828601614cb8565b82525060208301356020820152604083013560408201528093505050509250929050565b60006020828403121561548f57600080fd5b81356001600160401b038111156154a557600080fd5b8201601f810184136154b657600080fd5b6154c584823560208401614c61565b949350505050565b6000806000606084860312156154e257600080fd5b83359250602084013561528081614cdf565b6000806000806060858703121561550a57600080fd5b8435935060208501356001600160401b0381111561552757600080fd5b6155338782880161533f565b9598909750949560400135949350505050565b600080600080600060a0868803121561555e57600080fd5b853561556981614cdf565b9450602086013561557981614cdf565b9350604086013561558981614cdf565b9250606086013561559981614cdf565b949793965091946080013592915050565b600080600080600080600060e0888a0312156155c557600080fd5b8735965060208801359550604088013594506155e360608901614bb8565b93506080880135925060a08801356001600160401b0381111561560557600080fd5b6156118a828b01614cb8565b92505060c0880135905092959891949750929550565b6000806020838503121561563a57600080fd5b82356001600160401b038082111561565157600080fd5b818501915085601f83011261566557600080fd5b81358181111561567457600080fd5b8660208260051b850101111561568957600080fd5b60209290920196919550909350505050565b634e487b7160e01b600052601160045260246000fd5b808201808211156144795761447961569b565b91825260601b6001600160601b031916602082015260340190565b6000602082840312156156f157600080fd5b815161447681614cdf565b6020808252602a908201527f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160408201526939903ab73830bab9b2b960b11b606082015260800190565b60006020828403121561575857600080fd5b815161447681614f72565b60208082526028908201527f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160408201526739903830bab9b2b960c11b606082015260800190565b634e487b7160e01b600052603260045260246000fd5b6000826157de57634e487b7160e01b600052601260045260246000fd5b500690565b6020808252602e908201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160408201526d191e481a5b9a5d1a585b1a5e995960921b606082015260800190565b818103818111156144795761447961569b565b60006020828403121561585657600080fd5b5051919050565b60006020828403121561586f57600080fd5b81516001600160c01b038116811461447657600080fd5b60006020828403121561589857600080fd5b815161447681614dbc565b6001600160601b0381168114610c1557600080fd5b6000604082840312156158ca57600080fd5b6158d2614be6565b82516158dd81614cdf565b815260208301516158ed816158a3565b60208201529392505050565b60006001820161590b5761590b61569b565b5060010190565b60006020828403121561592457600080fd5b815167ffffffffffffffff198116811461447657600080fd5b60006020828403121561594f57600080fd5b8151614476816158a3565b6001600160601b0382811682821603908082111561597a5761597a61569b565b5092915050565b63ffffffff60e01b8360e01b1681526000600482018351602080860160005b838110156159bc578151855293820193908201906001016159a0565b5092979650505050505050565b8183823760009101908152919050565b60208082526052908201527f536572766963654d616e61676572426173652e6f6e6c7952656769737472794360408201527f6f6f7264696e61746f723a2063616c6c6572206973206e6f742074686520726560608201527133b4b9ba393c9031b7b7b93234b730ba37b960711b608082015260a00190565b6000815180845260005b81811015615a7757602081850181015186830182015201615a5b565b506000602082860101526020601f19601f83011685010191505092915050565b60018060a01b0383168152604060208201526000825160606040840152615ac160a0840182615a51565b90506020840151606084015260408401516080840152809150509392505050565b602081526000614cd86020830184615a51565b6001600160601b03818116838216028082169190828114613d3457613d3461569b565b80820281158282048414176144795761447961569b565b6001600160a01b038616815263ffffffff851660208201526080604082018190528101839052828460a0830137600060a08483010152600060a0601f19601f86011683010190508260608301529695505050505050565b60006101008a83528960208401528860408401526001600160601b0319881660608401528660808401528060a0840152615bc281840187615a51565b60c084019590955250506001600160a01b039190911660e0909101529695505050505050565b60008235609e19833603018112615bfe57600080fd5b9190910192915050565b803561462381614cdf565b8183526000602080850194508260005b858110156152eb578135615c3681614cdf565b6001600160a01b0316875281830135615c4e816158a3565b6001600160601b0316878401526040968701969190910190600101615c23565b60208082528181018390526000906040808401600586901b8501820187855b88811015615d6e57878303603f190184528135368b9003609e19018112615cb357600080fd5b8a0160a0813536839003601e19018112615ccc57600080fd5b820188810190356001600160401b03811115615ce757600080fd5b8060061b3603821315615cf957600080fd5b828752615d098388018284615c13565b92505050615d18888301615c08565b6001600160a01b03168886015281870135878601526060615d3a818401614f9d565b63ffffffff16908601526080615d51838201614f9d565b63ffffffff16950194909452509285019290850190600101615c8d565b509098975050505050505050565b600061ffff808316818103615d9357615d9361569b565b600101939250505056fe30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd470ea46f246ccfc58f7a93aa09bc6245a6818e97b1a160d186afe78993a3b194a0424c535369676e6174757265436865636b65722e636865636b5369676e617475a26469706673582212205998fafeb5eab20d19cf0bd86f2eecb1e9eff99b69e6178a9513067dd2808c5564736f6c63430008180033",
}

// ContractAlignedLayerServiceManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractAlignedLayerServiceManagerMetaData.ABI instead.
var ContractAlignedLayerServiceManagerABI = ContractAlignedLayerServiceManagerMetaData.ABI

// ContractAlignedLayerServiceManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractAlignedLayerServiceManagerMetaData.Bin instead.
var ContractAlignedLayerServiceManagerBin = ContractAlignedLayerServiceManagerMetaData.Bin

// DeployContractAlignedLayerServiceManager deploys a new Ethereum contract, binding an instance of ContractAlignedLayerServiceManager to it.
func DeployContractAlignedLayerServiceManager(auth *bind.TransactOpts, backend bind.ContractBackend, __avsDirectory common.Address, __rewardsCoordinator common.Address, __registryCoordinator common.Address, __stakeRegistry common.Address) (common.Address, *types.Transaction, *ContractAlignedLayerServiceManager, error) {
	parsed, err := ContractAlignedLayerServiceManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractAlignedLayerServiceManagerBin), backend, __avsDirectory, __rewardsCoordinator, __registryCoordinator, __stakeRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractAlignedLayerServiceManager{ContractAlignedLayerServiceManagerCaller: ContractAlignedLayerServiceManagerCaller{contract: contract}, ContractAlignedLayerServiceManagerTransactor: ContractAlignedLayerServiceManagerTransactor{contract: contract}, ContractAlignedLayerServiceManagerFilterer: ContractAlignedLayerServiceManagerFilterer{contract: contract}}, nil
}

// ContractAlignedLayerServiceManager is an auto generated Go binding around an Ethereum contract.
type ContractAlignedLayerServiceManager struct {
	ContractAlignedLayerServiceManagerCaller     // Read-only binding to the contract
	ContractAlignedLayerServiceManagerTransactor // Write-only binding to the contract
	ContractAlignedLayerServiceManagerFilterer   // Log filterer for contract events
}

// ContractAlignedLayerServiceManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractAlignedLayerServiceManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractAlignedLayerServiceManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractAlignedLayerServiceManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractAlignedLayerServiceManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractAlignedLayerServiceManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractAlignedLayerServiceManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractAlignedLayerServiceManagerSession struct {
	Contract     *ContractAlignedLayerServiceManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                       // Call options to use throughout this session
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// ContractAlignedLayerServiceManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractAlignedLayerServiceManagerCallerSession struct {
	Contract *ContractAlignedLayerServiceManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                             // Call options to use throughout this session
}

// ContractAlignedLayerServiceManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractAlignedLayerServiceManagerTransactorSession struct {
	Contract     *ContractAlignedLayerServiceManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                             // Transaction auth options to use throughout this session
}

// ContractAlignedLayerServiceManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractAlignedLayerServiceManagerRaw struct {
	Contract *ContractAlignedLayerServiceManager // Generic contract binding to access the raw methods on
}

// ContractAlignedLayerServiceManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractAlignedLayerServiceManagerCallerRaw struct {
	Contract *ContractAlignedLayerServiceManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ContractAlignedLayerServiceManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractAlignedLayerServiceManagerTransactorRaw struct {
	Contract *ContractAlignedLayerServiceManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractAlignedLayerServiceManager creates a new instance of ContractAlignedLayerServiceManager, bound to a specific deployed contract.
func NewContractAlignedLayerServiceManager(address common.Address, backend bind.ContractBackend) (*ContractAlignedLayerServiceManager, error) {
	contract, err := bindContractAlignedLayerServiceManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManager{ContractAlignedLayerServiceManagerCaller: ContractAlignedLayerServiceManagerCaller{contract: contract}, ContractAlignedLayerServiceManagerTransactor: ContractAlignedLayerServiceManagerTransactor{contract: contract}, ContractAlignedLayerServiceManagerFilterer: ContractAlignedLayerServiceManagerFilterer{contract: contract}}, nil
}

// NewContractAlignedLayerServiceManagerCaller creates a new read-only instance of ContractAlignedLayerServiceManager, bound to a specific deployed contract.
func NewContractAlignedLayerServiceManagerCaller(address common.Address, caller bind.ContractCaller) (*ContractAlignedLayerServiceManagerCaller, error) {
	contract, err := bindContractAlignedLayerServiceManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerCaller{contract: contract}, nil
}

// NewContractAlignedLayerServiceManagerTransactor creates a new write-only instance of ContractAlignedLayerServiceManager, bound to a specific deployed contract.
func NewContractAlignedLayerServiceManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractAlignedLayerServiceManagerTransactor, error) {
	contract, err := bindContractAlignedLayerServiceManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerTransactor{contract: contract}, nil
}

// NewContractAlignedLayerServiceManagerFilterer creates a new log filterer instance of ContractAlignedLayerServiceManager, bound to a specific deployed contract.
func NewContractAlignedLayerServiceManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractAlignedLayerServiceManagerFilterer, error) {
	contract, err := bindContractAlignedLayerServiceManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerFilterer{contract: contract}, nil
}

// bindContractAlignedLayerServiceManager binds a generic wrapper to an already deployed contract.
func bindContractAlignedLayerServiceManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractAlignedLayerServiceManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractAlignedLayerServiceManager.Contract.ContractAlignedLayerServiceManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.ContractAlignedLayerServiceManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.ContractAlignedLayerServiceManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractAlignedLayerServiceManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.contract.Transact(opts, method, params...)
}

// AlignedAggregator is a free data retrieval call binding the contract method 0x4a5bf632.
//
// Solidity: function alignedAggregator() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) AlignedAggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "alignedAggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AlignedAggregator is a free data retrieval call binding the contract method 0x4a5bf632.
//
// Solidity: function alignedAggregator() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) AlignedAggregator() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.AlignedAggregator(&_ContractAlignedLayerServiceManager.CallOpts)
}

// AlignedAggregator is a free data retrieval call binding the contract method 0x4a5bf632.
//
// Solidity: function alignedAggregator() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) AlignedAggregator() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.AlignedAggregator(&_ContractAlignedLayerServiceManager.CallOpts)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) AvsDirectory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "avsDirectory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) AvsDirectory() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.AvsDirectory(&_ContractAlignedLayerServiceManager.CallOpts)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) AvsDirectory() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.AvsDirectory(&_ContractAlignedLayerServiceManager.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ContractAlignedLayerServiceManager.Contract.BalanceOf(&_ContractAlignedLayerServiceManager.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ContractAlignedLayerServiceManager.Contract.BalanceOf(&_ContractAlignedLayerServiceManager.CallOpts, account)
}

// BatchersBalances is a free data retrieval call binding the contract method 0xf474b520.
//
// Solidity: function batchersBalances(address ) view returns(uint256)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) BatchersBalances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "batchersBalances", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BatchersBalances is a free data retrieval call binding the contract method 0xf474b520.
//
// Solidity: function batchersBalances(address ) view returns(uint256)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) BatchersBalances(arg0 common.Address) (*big.Int, error) {
	return _ContractAlignedLayerServiceManager.Contract.BatchersBalances(&_ContractAlignedLayerServiceManager.CallOpts, arg0)
}

// BatchersBalances is a free data retrieval call binding the contract method 0xf474b520.
//
// Solidity: function batchersBalances(address ) view returns(uint256)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) BatchersBalances(arg0 common.Address) (*big.Int, error) {
	return _ContractAlignedLayerServiceManager.Contract.BatchersBalances(&_ContractAlignedLayerServiceManager.CallOpts, arg0)
}

// BatchesState is a free data retrieval call binding the contract method 0xb099627e.
//
// Solidity: function batchesState(bytes32 ) view returns(uint32 taskCreatedBlock, bool responded, uint256 respondToTaskFeeLimit)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) BatchesState(opts *bind.CallOpts, arg0 [32]byte) (struct {
	TaskCreatedBlock      uint32
	Responded             bool
	RespondToTaskFeeLimit *big.Int
}, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "batchesState", arg0)

	outstruct := new(struct {
		TaskCreatedBlock      uint32
		Responded             bool
		RespondToTaskFeeLimit *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TaskCreatedBlock = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.Responded = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.RespondToTaskFeeLimit = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// BatchesState is a free data retrieval call binding the contract method 0xb099627e.
//
// Solidity: function batchesState(bytes32 ) view returns(uint32 taskCreatedBlock, bool responded, uint256 respondToTaskFeeLimit)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) BatchesState(arg0 [32]byte) (struct {
	TaskCreatedBlock      uint32
	Responded             bool
	RespondToTaskFeeLimit *big.Int
}, error) {
	return _ContractAlignedLayerServiceManager.Contract.BatchesState(&_ContractAlignedLayerServiceManager.CallOpts, arg0)
}

// BatchesState is a free data retrieval call binding the contract method 0xb099627e.
//
// Solidity: function batchesState(bytes32 ) view returns(uint32 taskCreatedBlock, bool responded, uint256 respondToTaskFeeLimit)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) BatchesState(arg0 [32]byte) (struct {
	TaskCreatedBlock      uint32
	Responded             bool
	RespondToTaskFeeLimit *big.Int
}, error) {
	return _ContractAlignedLayerServiceManager.Contract.BatchesState(&_ContractAlignedLayerServiceManager.CallOpts, arg0)
}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) BlsApkRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "blsApkRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) BlsApkRegistry() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.BlsApkRegistry(&_ContractAlignedLayerServiceManager.CallOpts)
}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) BlsApkRegistry() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.BlsApkRegistry(&_ContractAlignedLayerServiceManager.CallOpts)
}

// CheckPublicInput is a free data retrieval call binding the contract method 0x95c6d604.
//
// Solidity: function checkPublicInput(bytes publicInput, bytes32 hash) pure returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) CheckPublicInput(opts *bind.CallOpts, publicInput []byte, hash [32]byte) (bool, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "checkPublicInput", publicInput, hash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckPublicInput is a free data retrieval call binding the contract method 0x95c6d604.
//
// Solidity: function checkPublicInput(bytes publicInput, bytes32 hash) pure returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) CheckPublicInput(publicInput []byte, hash [32]byte) (bool, error) {
	return _ContractAlignedLayerServiceManager.Contract.CheckPublicInput(&_ContractAlignedLayerServiceManager.CallOpts, publicInput, hash)
}

// CheckPublicInput is a free data retrieval call binding the contract method 0x95c6d604.
//
// Solidity: function checkPublicInput(bytes publicInput, bytes32 hash) pure returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) CheckPublicInput(publicInput []byte, hash [32]byte) (bool, error) {
	return _ContractAlignedLayerServiceManager.Contract.CheckPublicInput(&_ContractAlignedLayerServiceManager.CallOpts, publicInput, hash)
}

// CheckSignatures is a free data retrieval call binding the contract method 0x4ae07c37.
//
// Solidity: function checkSignatures(bytes32 msgHash, uint32 referenceBlockNumber, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) params) view returns((uint96[],uint96[]), bytes32)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) CheckSignatures(opts *bind.CallOpts, msgHash [32]byte, referenceBlockNumber uint32, params IBLSSignatureCheckerNonSignerStakesAndSignature) (IBLSSignatureCheckerQuorumStakeTotals, [32]byte, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "checkSignatures", msgHash, referenceBlockNumber, params)

	if err != nil {
		return *new(IBLSSignatureCheckerQuorumStakeTotals), *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(IBLSSignatureCheckerQuorumStakeTotals)).(*IBLSSignatureCheckerQuorumStakeTotals)
	out1 := *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return out0, out1, err

}

// CheckSignatures is a free data retrieval call binding the contract method 0x4ae07c37.
//
// Solidity: function checkSignatures(bytes32 msgHash, uint32 referenceBlockNumber, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) params) view returns((uint96[],uint96[]), bytes32)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) CheckSignatures(msgHash [32]byte, referenceBlockNumber uint32, params IBLSSignatureCheckerNonSignerStakesAndSignature) (IBLSSignatureCheckerQuorumStakeTotals, [32]byte, error) {
	return _ContractAlignedLayerServiceManager.Contract.CheckSignatures(&_ContractAlignedLayerServiceManager.CallOpts, msgHash, referenceBlockNumber, params)
}

// CheckSignatures is a free data retrieval call binding the contract method 0x4ae07c37.
//
// Solidity: function checkSignatures(bytes32 msgHash, uint32 referenceBlockNumber, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) params) view returns((uint96[],uint96[]), bytes32)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) CheckSignatures(msgHash [32]byte, referenceBlockNumber uint32, params IBLSSignatureCheckerNonSignerStakesAndSignature) (IBLSSignatureCheckerQuorumStakeTotals, [32]byte, error) {
	return _ContractAlignedLayerServiceManager.Contract.CheckSignatures(&_ContractAlignedLayerServiceManager.CallOpts, msgHash, referenceBlockNumber, params)
}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) Delegation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "delegation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) Delegation() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.Delegation(&_ContractAlignedLayerServiceManager.CallOpts)
}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) Delegation() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.Delegation(&_ContractAlignedLayerServiceManager.CallOpts)
}

// DisabledVerifiers is a free data retrieval call binding the contract method 0xea5ca34b.
//
// Solidity: function disabledVerifiers() view returns(uint256)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) DisabledVerifiers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "disabledVerifiers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DisabledVerifiers is a free data retrieval call binding the contract method 0xea5ca34b.
//
// Solidity: function disabledVerifiers() view returns(uint256)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) DisabledVerifiers() (*big.Int, error) {
	return _ContractAlignedLayerServiceManager.Contract.DisabledVerifiers(&_ContractAlignedLayerServiceManager.CallOpts)
}

// DisabledVerifiers is a free data retrieval call binding the contract method 0xea5ca34b.
//
// Solidity: function disabledVerifiers() view returns(uint256)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) DisabledVerifiers() (*big.Int, error) {
	return _ContractAlignedLayerServiceManager.Contract.DisabledVerifiers(&_ContractAlignedLayerServiceManager.CallOpts)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) GetOperatorRestakedStrategies(opts *bind.CallOpts, operator common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "getOperatorRestakedStrategies", operator)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.GetOperatorRestakedStrategies(&_ContractAlignedLayerServiceManager.CallOpts, operator)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.GetOperatorRestakedStrategies(&_ContractAlignedLayerServiceManager.CallOpts, operator)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) GetRestakeableStrategies(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "getRestakeableStrategies")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.GetRestakeableStrategies(&_ContractAlignedLayerServiceManager.CallOpts)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.GetRestakeableStrategies(&_ContractAlignedLayerServiceManager.CallOpts)
}

// IsVerifierDisabled is a free data retrieval call binding the contract method 0x137122b5.
//
// Solidity: function isVerifierDisabled(uint8 verifierIdx) view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) IsVerifierDisabled(opts *bind.CallOpts, verifierIdx uint8) (bool, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "isVerifierDisabled", verifierIdx)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsVerifierDisabled is a free data retrieval call binding the contract method 0x137122b5.
//
// Solidity: function isVerifierDisabled(uint8 verifierIdx) view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) IsVerifierDisabled(verifierIdx uint8) (bool, error) {
	return _ContractAlignedLayerServiceManager.Contract.IsVerifierDisabled(&_ContractAlignedLayerServiceManager.CallOpts, verifierIdx)
}

// IsVerifierDisabled is a free data retrieval call binding the contract method 0x137122b5.
//
// Solidity: function isVerifierDisabled(uint8 verifierIdx) view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) IsVerifierDisabled(verifierIdx uint8) (bool, error) {
	return _ContractAlignedLayerServiceManager.Contract.IsVerifierDisabled(&_ContractAlignedLayerServiceManager.CallOpts, verifierIdx)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) Owner() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.Owner(&_ContractAlignedLayerServiceManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) Owner() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.Owner(&_ContractAlignedLayerServiceManager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) Paused(opts *bind.CallOpts, index uint8) (bool, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "paused", index)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) Paused(index uint8) (bool, error) {
	return _ContractAlignedLayerServiceManager.Contract.Paused(&_ContractAlignedLayerServiceManager.CallOpts, index)
}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) Paused(index uint8) (bool, error) {
	return _ContractAlignedLayerServiceManager.Contract.Paused(&_ContractAlignedLayerServiceManager.CallOpts, index)
}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) Paused0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "paused0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) Paused0() (*big.Int, error) {
	return _ContractAlignedLayerServiceManager.Contract.Paused0(&_ContractAlignedLayerServiceManager.CallOpts)
}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) Paused0() (*big.Int, error) {
	return _ContractAlignedLayerServiceManager.Contract.Paused0(&_ContractAlignedLayerServiceManager.CallOpts)
}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) PauserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "pauserRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) PauserRegistry() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.PauserRegistry(&_ContractAlignedLayerServiceManager.CallOpts)
}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) PauserRegistry() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.PauserRegistry(&_ContractAlignedLayerServiceManager.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.RegistryCoordinator(&_ContractAlignedLayerServiceManager.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.RegistryCoordinator(&_ContractAlignedLayerServiceManager.CallOpts)
}

// RewardsInitiator is a free data retrieval call binding the contract method 0xfc299dee.
//
// Solidity: function rewardsInitiator() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) RewardsInitiator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "rewardsInitiator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardsInitiator is a free data retrieval call binding the contract method 0xfc299dee.
//
// Solidity: function rewardsInitiator() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) RewardsInitiator() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.RewardsInitiator(&_ContractAlignedLayerServiceManager.CallOpts)
}

// RewardsInitiator is a free data retrieval call binding the contract method 0xfc299dee.
//
// Solidity: function rewardsInitiator() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) RewardsInitiator() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.RewardsInitiator(&_ContractAlignedLayerServiceManager.CallOpts)
}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) StakeRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "stakeRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) StakeRegistry() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.StakeRegistry(&_ContractAlignedLayerServiceManager.CallOpts)
}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) StakeRegistry() (common.Address, error) {
	return _ContractAlignedLayerServiceManager.Contract.StakeRegistry(&_ContractAlignedLayerServiceManager.CallOpts)
}

// StaleStakesForbidden is a free data retrieval call binding the contract method 0xb98d0908.
//
// Solidity: function staleStakesForbidden() view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) StaleStakesForbidden(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "staleStakesForbidden")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// StaleStakesForbidden is a free data retrieval call binding the contract method 0xb98d0908.
//
// Solidity: function staleStakesForbidden() view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) StaleStakesForbidden() (bool, error) {
	return _ContractAlignedLayerServiceManager.Contract.StaleStakesForbidden(&_ContractAlignedLayerServiceManager.CallOpts)
}

// StaleStakesForbidden is a free data retrieval call binding the contract method 0xb98d0908.
//
// Solidity: function staleStakesForbidden() view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) StaleStakesForbidden() (bool, error) {
	return _ContractAlignedLayerServiceManager.Contract.StaleStakesForbidden(&_ContractAlignedLayerServiceManager.CallOpts)
}

// TrySignatureAndApkVerification is a free data retrieval call binding the contract method 0x171f1d5b.
//
// Solidity: function trySignatureAndApkVerification(bytes32 msgHash, (uint256,uint256) apk, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma) view returns(bool pairingSuccessful, bool siganatureIsValid)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) TrySignatureAndApkVerification(opts *bind.CallOpts, msgHash [32]byte, apk BN254G1Point, apkG2 BN254G2Point, sigma BN254G1Point) (struct {
	PairingSuccessful bool
	SiganatureIsValid bool
}, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "trySignatureAndApkVerification", msgHash, apk, apkG2, sigma)

	outstruct := new(struct {
		PairingSuccessful bool
		SiganatureIsValid bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.PairingSuccessful = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.SiganatureIsValid = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// TrySignatureAndApkVerification is a free data retrieval call binding the contract method 0x171f1d5b.
//
// Solidity: function trySignatureAndApkVerification(bytes32 msgHash, (uint256,uint256) apk, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma) view returns(bool pairingSuccessful, bool siganatureIsValid)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) TrySignatureAndApkVerification(msgHash [32]byte, apk BN254G1Point, apkG2 BN254G2Point, sigma BN254G1Point) (struct {
	PairingSuccessful bool
	SiganatureIsValid bool
}, error) {
	return _ContractAlignedLayerServiceManager.Contract.TrySignatureAndApkVerification(&_ContractAlignedLayerServiceManager.CallOpts, msgHash, apk, apkG2, sigma)
}

// TrySignatureAndApkVerification is a free data retrieval call binding the contract method 0x171f1d5b.
//
// Solidity: function trySignatureAndApkVerification(bytes32 msgHash, (uint256,uint256) apk, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma) view returns(bool pairingSuccessful, bool siganatureIsValid)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) TrySignatureAndApkVerification(msgHash [32]byte, apk BN254G1Point, apkG2 BN254G2Point, sigma BN254G1Point) (struct {
	PairingSuccessful bool
	SiganatureIsValid bool
}, error) {
	return _ContractAlignedLayerServiceManager.Contract.TrySignatureAndApkVerification(&_ContractAlignedLayerServiceManager.CallOpts, msgHash, apk, apkG2, sigma)
}

// VerifyBatchInclusion is a free data retrieval call binding the contract method 0x06045a91.
//
// Solidity: function verifyBatchInclusion(bytes32 proofCommitment, bytes32 pubInputCommitment, bytes32 provingSystemAuxDataCommitment, bytes20 proofGeneratorAddr, bytes32 batchMerkleRoot, bytes merkleProof, uint256 verificationDataBatchIndex, address senderAddress) view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) VerifyBatchInclusion(opts *bind.CallOpts, proofCommitment [32]byte, pubInputCommitment [32]byte, provingSystemAuxDataCommitment [32]byte, proofGeneratorAddr [20]byte, batchMerkleRoot [32]byte, merkleProof []byte, verificationDataBatchIndex *big.Int, senderAddress common.Address) (bool, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "verifyBatchInclusion", proofCommitment, pubInputCommitment, provingSystemAuxDataCommitment, proofGeneratorAddr, batchMerkleRoot, merkleProof, verificationDataBatchIndex, senderAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyBatchInclusion is a free data retrieval call binding the contract method 0x06045a91.
//
// Solidity: function verifyBatchInclusion(bytes32 proofCommitment, bytes32 pubInputCommitment, bytes32 provingSystemAuxDataCommitment, bytes20 proofGeneratorAddr, bytes32 batchMerkleRoot, bytes merkleProof, uint256 verificationDataBatchIndex, address senderAddress) view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) VerifyBatchInclusion(proofCommitment [32]byte, pubInputCommitment [32]byte, provingSystemAuxDataCommitment [32]byte, proofGeneratorAddr [20]byte, batchMerkleRoot [32]byte, merkleProof []byte, verificationDataBatchIndex *big.Int, senderAddress common.Address) (bool, error) {
	return _ContractAlignedLayerServiceManager.Contract.VerifyBatchInclusion(&_ContractAlignedLayerServiceManager.CallOpts, proofCommitment, pubInputCommitment, provingSystemAuxDataCommitment, proofGeneratorAddr, batchMerkleRoot, merkleProof, verificationDataBatchIndex, senderAddress)
}

// VerifyBatchInclusion is a free data retrieval call binding the contract method 0x06045a91.
//
// Solidity: function verifyBatchInclusion(bytes32 proofCommitment, bytes32 pubInputCommitment, bytes32 provingSystemAuxDataCommitment, bytes20 proofGeneratorAddr, bytes32 batchMerkleRoot, bytes merkleProof, uint256 verificationDataBatchIndex, address senderAddress) view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) VerifyBatchInclusion(proofCommitment [32]byte, pubInputCommitment [32]byte, provingSystemAuxDataCommitment [32]byte, proofGeneratorAddr [20]byte, batchMerkleRoot [32]byte, merkleProof []byte, verificationDataBatchIndex *big.Int, senderAddress common.Address) (bool, error) {
	return _ContractAlignedLayerServiceManager.Contract.VerifyBatchInclusion(&_ContractAlignedLayerServiceManager.CallOpts, proofCommitment, pubInputCommitment, provingSystemAuxDataCommitment, proofGeneratorAddr, batchMerkleRoot, merkleProof, verificationDataBatchIndex, senderAddress)
}

// VerifyBatchInclusion0 is a free data retrieval call binding the contract method 0xfa534dc0.
//
// Solidity: function verifyBatchInclusion(bytes32 proofCommitment, bytes32 pubInputCommitment, bytes32 provingSystemAuxDataCommitment, bytes20 proofGeneratorAddr, bytes32 batchMerkleRoot, bytes merkleProof, uint256 verificationDataBatchIndex) view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCaller) VerifyBatchInclusion0(opts *bind.CallOpts, proofCommitment [32]byte, pubInputCommitment [32]byte, provingSystemAuxDataCommitment [32]byte, proofGeneratorAddr [20]byte, batchMerkleRoot [32]byte, merkleProof []byte, verificationDataBatchIndex *big.Int) (bool, error) {
	var out []interface{}
	err := _ContractAlignedLayerServiceManager.contract.Call(opts, &out, "verifyBatchInclusion0", proofCommitment, pubInputCommitment, provingSystemAuxDataCommitment, proofGeneratorAddr, batchMerkleRoot, merkleProof, verificationDataBatchIndex)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyBatchInclusion0 is a free data retrieval call binding the contract method 0xfa534dc0.
//
// Solidity: function verifyBatchInclusion(bytes32 proofCommitment, bytes32 pubInputCommitment, bytes32 provingSystemAuxDataCommitment, bytes20 proofGeneratorAddr, bytes32 batchMerkleRoot, bytes merkleProof, uint256 verificationDataBatchIndex) view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) VerifyBatchInclusion0(proofCommitment [32]byte, pubInputCommitment [32]byte, provingSystemAuxDataCommitment [32]byte, proofGeneratorAddr [20]byte, batchMerkleRoot [32]byte, merkleProof []byte, verificationDataBatchIndex *big.Int) (bool, error) {
	return _ContractAlignedLayerServiceManager.Contract.VerifyBatchInclusion0(&_ContractAlignedLayerServiceManager.CallOpts, proofCommitment, pubInputCommitment, provingSystemAuxDataCommitment, proofGeneratorAddr, batchMerkleRoot, merkleProof, verificationDataBatchIndex)
}

// VerifyBatchInclusion0 is a free data retrieval call binding the contract method 0xfa534dc0.
//
// Solidity: function verifyBatchInclusion(bytes32 proofCommitment, bytes32 pubInputCommitment, bytes32 provingSystemAuxDataCommitment, bytes20 proofGeneratorAddr, bytes32 batchMerkleRoot, bytes merkleProof, uint256 verificationDataBatchIndex) view returns(bool)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerCallerSession) VerifyBatchInclusion0(proofCommitment [32]byte, pubInputCommitment [32]byte, provingSystemAuxDataCommitment [32]byte, proofGeneratorAddr [20]byte, batchMerkleRoot [32]byte, merkleProof []byte, verificationDataBatchIndex *big.Int) (bool, error) {
	return _ContractAlignedLayerServiceManager.Contract.VerifyBatchInclusion0(&_ContractAlignedLayerServiceManager.CallOpts, proofCommitment, pubInputCommitment, provingSystemAuxDataCommitment, proofGeneratorAddr, batchMerkleRoot, merkleProof, verificationDataBatchIndex)
}

// CreateAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xfce36c7d.
//
// Solidity: function createAVSRewardsSubmission(((address,uint96)[],address,uint256,uint32,uint32)[] rewardsSubmissions) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) CreateAVSRewardsSubmission(opts *bind.TransactOpts, rewardsSubmissions []IRewardsCoordinatorRewardsSubmission) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "createAVSRewardsSubmission", rewardsSubmissions)
}

// CreateAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xfce36c7d.
//
// Solidity: function createAVSRewardsSubmission(((address,uint96)[],address,uint256,uint32,uint32)[] rewardsSubmissions) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) CreateAVSRewardsSubmission(rewardsSubmissions []IRewardsCoordinatorRewardsSubmission) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.CreateAVSRewardsSubmission(&_ContractAlignedLayerServiceManager.TransactOpts, rewardsSubmissions)
}

// CreateAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xfce36c7d.
//
// Solidity: function createAVSRewardsSubmission(((address,uint96)[],address,uint256,uint32,uint32)[] rewardsSubmissions) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) CreateAVSRewardsSubmission(rewardsSubmissions []IRewardsCoordinatorRewardsSubmission) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.CreateAVSRewardsSubmission(&_ContractAlignedLayerServiceManager.TransactOpts, rewardsSubmissions)
}

// CreateNewTask is a paid mutator transaction binding the contract method 0xd66eaabd.
//
// Solidity: function createNewTask(bytes32 batchMerkleRoot, string batchDataPointer, uint256 respondToTaskFeeLimit) payable returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) CreateNewTask(opts *bind.TransactOpts, batchMerkleRoot [32]byte, batchDataPointer string, respondToTaskFeeLimit *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "createNewTask", batchMerkleRoot, batchDataPointer, respondToTaskFeeLimit)
}

// CreateNewTask is a paid mutator transaction binding the contract method 0xd66eaabd.
//
// Solidity: function createNewTask(bytes32 batchMerkleRoot, string batchDataPointer, uint256 respondToTaskFeeLimit) payable returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) CreateNewTask(batchMerkleRoot [32]byte, batchDataPointer string, respondToTaskFeeLimit *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.CreateNewTask(&_ContractAlignedLayerServiceManager.TransactOpts, batchMerkleRoot, batchDataPointer, respondToTaskFeeLimit)
}

// CreateNewTask is a paid mutator transaction binding the contract method 0xd66eaabd.
//
// Solidity: function createNewTask(bytes32 batchMerkleRoot, string batchDataPointer, uint256 respondToTaskFeeLimit) payable returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) CreateNewTask(batchMerkleRoot [32]byte, batchDataPointer string, respondToTaskFeeLimit *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.CreateNewTask(&_ContractAlignedLayerServiceManager.TransactOpts, batchMerkleRoot, batchDataPointer, respondToTaskFeeLimit)
}

// DepositToBatcher is a paid mutator transaction binding the contract method 0x4223d551.
//
// Solidity: function depositToBatcher(address account) payable returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) DepositToBatcher(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "depositToBatcher", account)
}

// DepositToBatcher is a paid mutator transaction binding the contract method 0x4223d551.
//
// Solidity: function depositToBatcher(address account) payable returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) DepositToBatcher(account common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.DepositToBatcher(&_ContractAlignedLayerServiceManager.TransactOpts, account)
}

// DepositToBatcher is a paid mutator transaction binding the contract method 0x4223d551.
//
// Solidity: function depositToBatcher(address account) payable returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) DepositToBatcher(account common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.DepositToBatcher(&_ContractAlignedLayerServiceManager.TransactOpts, account)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) DeregisterOperatorFromAVS(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "deregisterOperatorFromAVS", operator)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) DeregisterOperatorFromAVS(operator common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.DeregisterOperatorFromAVS(&_ContractAlignedLayerServiceManager.TransactOpts, operator)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) DeregisterOperatorFromAVS(operator common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.DeregisterOperatorFromAVS(&_ContractAlignedLayerServiceManager.TransactOpts, operator)
}

// DisableVerifier is a paid mutator transaction binding the contract method 0xfd4c3b7c.
//
// Solidity: function disableVerifier(uint8 verifierIdx) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) DisableVerifier(opts *bind.TransactOpts, verifierIdx uint8) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "disableVerifier", verifierIdx)
}

// DisableVerifier is a paid mutator transaction binding the contract method 0xfd4c3b7c.
//
// Solidity: function disableVerifier(uint8 verifierIdx) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) DisableVerifier(verifierIdx uint8) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.DisableVerifier(&_ContractAlignedLayerServiceManager.TransactOpts, verifierIdx)
}

// DisableVerifier is a paid mutator transaction binding the contract method 0xfd4c3b7c.
//
// Solidity: function disableVerifier(uint8 verifierIdx) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) DisableVerifier(verifierIdx uint8) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.DisableVerifier(&_ContractAlignedLayerServiceManager.TransactOpts, verifierIdx)
}

// EnableVerifier is a paid mutator transaction binding the contract method 0x18daeeaf.
//
// Solidity: function enableVerifier(uint8 verifierIdx) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) EnableVerifier(opts *bind.TransactOpts, verifierIdx uint8) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "enableVerifier", verifierIdx)
}

// EnableVerifier is a paid mutator transaction binding the contract method 0x18daeeaf.
//
// Solidity: function enableVerifier(uint8 verifierIdx) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) EnableVerifier(verifierIdx uint8) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.EnableVerifier(&_ContractAlignedLayerServiceManager.TransactOpts, verifierIdx)
}

// EnableVerifier is a paid mutator transaction binding the contract method 0x18daeeaf.
//
// Solidity: function enableVerifier(uint8 verifierIdx) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) EnableVerifier(verifierIdx uint8) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.EnableVerifier(&_ContractAlignedLayerServiceManager.TransactOpts, verifierIdx)
}

// Initialize is a paid mutator transaction binding the contract method 0xf7013ef6.
//
// Solidity: function initialize(address _initialOwner, address _rewardsInitiator, address _alignedAggregator, address _pauserRegistry, uint256 _initialPausedStatus) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address, _rewardsInitiator common.Address, _alignedAggregator common.Address, _pauserRegistry common.Address, _initialPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "initialize", _initialOwner, _rewardsInitiator, _alignedAggregator, _pauserRegistry, _initialPausedStatus)
}

// Initialize is a paid mutator transaction binding the contract method 0xf7013ef6.
//
// Solidity: function initialize(address _initialOwner, address _rewardsInitiator, address _alignedAggregator, address _pauserRegistry, uint256 _initialPausedStatus) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) Initialize(_initialOwner common.Address, _rewardsInitiator common.Address, _alignedAggregator common.Address, _pauserRegistry common.Address, _initialPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.Initialize(&_ContractAlignedLayerServiceManager.TransactOpts, _initialOwner, _rewardsInitiator, _alignedAggregator, _pauserRegistry, _initialPausedStatus)
}

// Initialize is a paid mutator transaction binding the contract method 0xf7013ef6.
//
// Solidity: function initialize(address _initialOwner, address _rewardsInitiator, address _alignedAggregator, address _pauserRegistry, uint256 _initialPausedStatus) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) Initialize(_initialOwner common.Address, _rewardsInitiator common.Address, _alignedAggregator common.Address, _pauserRegistry common.Address, _initialPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.Initialize(&_ContractAlignedLayerServiceManager.TransactOpts, _initialOwner, _rewardsInitiator, _alignedAggregator, _pauserRegistry, _initialPausedStatus)
}

// InitializeAggregator is a paid mutator transaction binding the contract method 0x800fb61f.
//
// Solidity: function initializeAggregator(address _alignedAggregator) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) InitializeAggregator(opts *bind.TransactOpts, _alignedAggregator common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "initializeAggregator", _alignedAggregator)
}

// InitializeAggregator is a paid mutator transaction binding the contract method 0x800fb61f.
//
// Solidity: function initializeAggregator(address _alignedAggregator) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) InitializeAggregator(_alignedAggregator common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.InitializeAggregator(&_ContractAlignedLayerServiceManager.TransactOpts, _alignedAggregator)
}

// InitializeAggregator is a paid mutator transaction binding the contract method 0x800fb61f.
//
// Solidity: function initializeAggregator(address _alignedAggregator) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) InitializeAggregator(_alignedAggregator common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.InitializeAggregator(&_ContractAlignedLayerServiceManager.TransactOpts, _alignedAggregator)
}

// InitializePauser is a paid mutator transaction binding the contract method 0x2585b25b.
//
// Solidity: function initializePauser(address _pauserRegistry, uint256 _initialPausedStatus) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) InitializePauser(opts *bind.TransactOpts, _pauserRegistry common.Address, _initialPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "initializePauser", _pauserRegistry, _initialPausedStatus)
}

// InitializePauser is a paid mutator transaction binding the contract method 0x2585b25b.
//
// Solidity: function initializePauser(address _pauserRegistry, uint256 _initialPausedStatus) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) InitializePauser(_pauserRegistry common.Address, _initialPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.InitializePauser(&_ContractAlignedLayerServiceManager.TransactOpts, _pauserRegistry, _initialPausedStatus)
}

// InitializePauser is a paid mutator transaction binding the contract method 0x2585b25b.
//
// Solidity: function initializePauser(address _pauserRegistry, uint256 _initialPausedStatus) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) InitializePauser(_pauserRegistry common.Address, _initialPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.InitializePauser(&_ContractAlignedLayerServiceManager.TransactOpts, _pauserRegistry, _initialPausedStatus)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) Pause(opts *bind.TransactOpts, newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "pause", newPausedStatus)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) Pause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.Pause(&_ContractAlignedLayerServiceManager.TransactOpts, newPausedStatus)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) Pause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.Pause(&_ContractAlignedLayerServiceManager.TransactOpts, newPausedStatus)
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) PauseAll(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "pauseAll")
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) PauseAll() (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.PauseAll(&_ContractAlignedLayerServiceManager.TransactOpts)
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) PauseAll() (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.PauseAll(&_ContractAlignedLayerServiceManager.TransactOpts)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) RegisterOperatorToAVS(opts *bind.TransactOpts, operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "registerOperatorToAVS", operator, operatorSignature)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) RegisterOperatorToAVS(operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.RegisterOperatorToAVS(&_ContractAlignedLayerServiceManager.TransactOpts, operator, operatorSignature)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) RegisterOperatorToAVS(operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.RegisterOperatorToAVS(&_ContractAlignedLayerServiceManager.TransactOpts, operator, operatorSignature)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.RenounceOwnership(&_ContractAlignedLayerServiceManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.RenounceOwnership(&_ContractAlignedLayerServiceManager.TransactOpts)
}

// RespondToTaskV2 is a paid mutator transaction binding the contract method 0xab21739a.
//
// Solidity: function respondToTaskV2(bytes32 batchMerkleRoot, address senderAddress, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) RespondToTaskV2(opts *bind.TransactOpts, batchMerkleRoot [32]byte, senderAddress common.Address, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "respondToTaskV2", batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
}

// RespondToTaskV2 is a paid mutator transaction binding the contract method 0xab21739a.
//
// Solidity: function respondToTaskV2(bytes32 batchMerkleRoot, address senderAddress, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) RespondToTaskV2(batchMerkleRoot [32]byte, senderAddress common.Address, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.RespondToTaskV2(&_ContractAlignedLayerServiceManager.TransactOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
}

// RespondToTaskV2 is a paid mutator transaction binding the contract method 0xab21739a.
//
// Solidity: function respondToTaskV2(bytes32 batchMerkleRoot, address senderAddress, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) RespondToTaskV2(batchMerkleRoot [32]byte, senderAddress common.Address, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.RespondToTaskV2(&_ContractAlignedLayerServiceManager.TransactOpts, batchMerkleRoot, senderAddress, nonSignerStakesAndSignature)
}

// SetAggregator is a paid mutator transaction binding the contract method 0xf9120af6.
//
// Solidity: function setAggregator(address _alignedAggregator) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) SetAggregator(opts *bind.TransactOpts, _alignedAggregator common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "setAggregator", _alignedAggregator)
}

// SetAggregator is a paid mutator transaction binding the contract method 0xf9120af6.
//
// Solidity: function setAggregator(address _alignedAggregator) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) SetAggregator(_alignedAggregator common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.SetAggregator(&_ContractAlignedLayerServiceManager.TransactOpts, _alignedAggregator)
}

// SetAggregator is a paid mutator transaction binding the contract method 0xf9120af6.
//
// Solidity: function setAggregator(address _alignedAggregator) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) SetAggregator(_alignedAggregator common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.SetAggregator(&_ContractAlignedLayerServiceManager.TransactOpts, _alignedAggregator)
}

// SetDisabledVerifiers is a paid mutator transaction binding the contract method 0xb753645e.
//
// Solidity: function setDisabledVerifiers(uint256 bitmap) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) SetDisabledVerifiers(opts *bind.TransactOpts, bitmap *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "setDisabledVerifiers", bitmap)
}

// SetDisabledVerifiers is a paid mutator transaction binding the contract method 0xb753645e.
//
// Solidity: function setDisabledVerifiers(uint256 bitmap) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) SetDisabledVerifiers(bitmap *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.SetDisabledVerifiers(&_ContractAlignedLayerServiceManager.TransactOpts, bitmap)
}

// SetDisabledVerifiers is a paid mutator transaction binding the contract method 0xb753645e.
//
// Solidity: function setDisabledVerifiers(uint256 bitmap) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) SetDisabledVerifiers(bitmap *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.SetDisabledVerifiers(&_ContractAlignedLayerServiceManager.TransactOpts, bitmap)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) SetPauserRegistry(opts *bind.TransactOpts, newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "setPauserRegistry", newPauserRegistry)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) SetPauserRegistry(newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.SetPauserRegistry(&_ContractAlignedLayerServiceManager.TransactOpts, newPauserRegistry)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) SetPauserRegistry(newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.SetPauserRegistry(&_ContractAlignedLayerServiceManager.TransactOpts, newPauserRegistry)
}

// SetRewardsInitiator is a paid mutator transaction binding the contract method 0x3bc28c8c.
//
// Solidity: function setRewardsInitiator(address newRewardsInitiator) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) SetRewardsInitiator(opts *bind.TransactOpts, newRewardsInitiator common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "setRewardsInitiator", newRewardsInitiator)
}

// SetRewardsInitiator is a paid mutator transaction binding the contract method 0x3bc28c8c.
//
// Solidity: function setRewardsInitiator(address newRewardsInitiator) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) SetRewardsInitiator(newRewardsInitiator common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.SetRewardsInitiator(&_ContractAlignedLayerServiceManager.TransactOpts, newRewardsInitiator)
}

// SetRewardsInitiator is a paid mutator transaction binding the contract method 0x3bc28c8c.
//
// Solidity: function setRewardsInitiator(address newRewardsInitiator) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) SetRewardsInitiator(newRewardsInitiator common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.SetRewardsInitiator(&_ContractAlignedLayerServiceManager.TransactOpts, newRewardsInitiator)
}

// SetStaleStakesForbidden is a paid mutator transaction binding the contract method 0x416c7e5e.
//
// Solidity: function setStaleStakesForbidden(bool value) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) SetStaleStakesForbidden(opts *bind.TransactOpts, value bool) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "setStaleStakesForbidden", value)
}

// SetStaleStakesForbidden is a paid mutator transaction binding the contract method 0x416c7e5e.
//
// Solidity: function setStaleStakesForbidden(bool value) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) SetStaleStakesForbidden(value bool) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.SetStaleStakesForbidden(&_ContractAlignedLayerServiceManager.TransactOpts, value)
}

// SetStaleStakesForbidden is a paid mutator transaction binding the contract method 0x416c7e5e.
//
// Solidity: function setStaleStakesForbidden(bool value) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) SetStaleStakesForbidden(value bool) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.SetStaleStakesForbidden(&_ContractAlignedLayerServiceManager.TransactOpts, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.TransferOwnership(&_ContractAlignedLayerServiceManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.TransferOwnership(&_ContractAlignedLayerServiceManager.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) Unpause(opts *bind.TransactOpts, newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "unpause", newPausedStatus)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) Unpause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.Unpause(&_ContractAlignedLayerServiceManager.TransactOpts, newPausedStatus)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) Unpause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.Unpause(&_ContractAlignedLayerServiceManager.TransactOpts, newPausedStatus)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) UpdateAVSMetadataURI(opts *bind.TransactOpts, _metadataURI string) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "updateAVSMetadataURI", _metadataURI)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) UpdateAVSMetadataURI(_metadataURI string) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.UpdateAVSMetadataURI(&_ContractAlignedLayerServiceManager.TransactOpts, _metadataURI)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) UpdateAVSMetadataURI(_metadataURI string) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.UpdateAVSMetadataURI(&_ContractAlignedLayerServiceManager.TransactOpts, _metadataURI)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.Withdraw(&_ContractAlignedLayerServiceManager.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.Withdraw(&_ContractAlignedLayerServiceManager.TransactOpts, amount)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerSession) Receive() (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.Receive(&_ContractAlignedLayerServiceManager.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerTransactorSession) Receive() (*types.Transaction, error) {
	return _ContractAlignedLayerServiceManager.Contract.Receive(&_ContractAlignedLayerServiceManager.TransactOpts)
}

// ContractAlignedLayerServiceManagerBatchVerifiedIterator is returned from FilterBatchVerified and is used to iterate over the raw logs and unpacked data for BatchVerified events raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerBatchVerifiedIterator struct {
	Event *ContractAlignedLayerServiceManagerBatchVerified // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAlignedLayerServiceManagerBatchVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAlignedLayerServiceManagerBatchVerified)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAlignedLayerServiceManagerBatchVerified)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAlignedLayerServiceManagerBatchVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAlignedLayerServiceManagerBatchVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAlignedLayerServiceManagerBatchVerified represents a BatchVerified event raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerBatchVerified struct {
	BatchMerkleRoot [32]byte
	SenderAddress   common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBatchVerified is a free log retrieval operation binding the contract event 0x8511746b73275e06971968773119b9601fc501d7bdf3824d8754042d148940e2.
//
// Solidity: event BatchVerified(bytes32 indexed batchMerkleRoot, address senderAddress)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) FilterBatchVerified(opts *bind.FilterOpts, batchMerkleRoot [][32]byte) (*ContractAlignedLayerServiceManagerBatchVerifiedIterator, error) {

	var batchMerkleRootRule []interface{}
	for _, batchMerkleRootItem := range batchMerkleRoot {
		batchMerkleRootRule = append(batchMerkleRootRule, batchMerkleRootItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.FilterLogs(opts, "BatchVerified", batchMerkleRootRule)
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerBatchVerifiedIterator{contract: _ContractAlignedLayerServiceManager.contract, event: "BatchVerified", logs: logs, sub: sub}, nil
}

// WatchBatchVerified is a free log subscription operation binding the contract event 0x8511746b73275e06971968773119b9601fc501d7bdf3824d8754042d148940e2.
//
// Solidity: event BatchVerified(bytes32 indexed batchMerkleRoot, address senderAddress)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) WatchBatchVerified(opts *bind.WatchOpts, sink chan<- *ContractAlignedLayerServiceManagerBatchVerified, batchMerkleRoot [][32]byte) (event.Subscription, error) {

	var batchMerkleRootRule []interface{}
	for _, batchMerkleRootItem := range batchMerkleRoot {
		batchMerkleRootRule = append(batchMerkleRootRule, batchMerkleRootItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.WatchLogs(opts, "BatchVerified", batchMerkleRootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAlignedLayerServiceManagerBatchVerified)
				if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "BatchVerified", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBatchVerified is a log parse operation binding the contract event 0x8511746b73275e06971968773119b9601fc501d7bdf3824d8754042d148940e2.
//
// Solidity: event BatchVerified(bytes32 indexed batchMerkleRoot, address senderAddress)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) ParseBatchVerified(log types.Log) (*ContractAlignedLayerServiceManagerBatchVerified, error) {
	event := new(ContractAlignedLayerServiceManagerBatchVerified)
	if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "BatchVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractAlignedLayerServiceManagerBatcherBalanceUpdatedIterator is returned from FilterBatcherBalanceUpdated and is used to iterate over the raw logs and unpacked data for BatcherBalanceUpdated events raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerBatcherBalanceUpdatedIterator struct {
	Event *ContractAlignedLayerServiceManagerBatcherBalanceUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAlignedLayerServiceManagerBatcherBalanceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAlignedLayerServiceManagerBatcherBalanceUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAlignedLayerServiceManagerBatcherBalanceUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAlignedLayerServiceManagerBatcherBalanceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAlignedLayerServiceManagerBatcherBalanceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAlignedLayerServiceManagerBatcherBalanceUpdated represents a BatcherBalanceUpdated event raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerBatcherBalanceUpdated struct {
	Batcher    common.Address
	NewBalance *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBatcherBalanceUpdated is a free log retrieval operation binding the contract event 0x0ea46f246ccfc58f7a93aa09bc6245a6818e97b1a160d186afe78993a3b194a0.
//
// Solidity: event BatcherBalanceUpdated(address indexed batcher, uint256 newBalance)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) FilterBatcherBalanceUpdated(opts *bind.FilterOpts, batcher []common.Address) (*ContractAlignedLayerServiceManagerBatcherBalanceUpdatedIterator, error) {

	var batcherRule []interface{}
	for _, batcherItem := range batcher {
		batcherRule = append(batcherRule, batcherItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.FilterLogs(opts, "BatcherBalanceUpdated", batcherRule)
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerBatcherBalanceUpdatedIterator{contract: _ContractAlignedLayerServiceManager.contract, event: "BatcherBalanceUpdated", logs: logs, sub: sub}, nil
}

// WatchBatcherBalanceUpdated is a free log subscription operation binding the contract event 0x0ea46f246ccfc58f7a93aa09bc6245a6818e97b1a160d186afe78993a3b194a0.
//
// Solidity: event BatcherBalanceUpdated(address indexed batcher, uint256 newBalance)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) WatchBatcherBalanceUpdated(opts *bind.WatchOpts, sink chan<- *ContractAlignedLayerServiceManagerBatcherBalanceUpdated, batcher []common.Address) (event.Subscription, error) {

	var batcherRule []interface{}
	for _, batcherItem := range batcher {
		batcherRule = append(batcherRule, batcherItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.WatchLogs(opts, "BatcherBalanceUpdated", batcherRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAlignedLayerServiceManagerBatcherBalanceUpdated)
				if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "BatcherBalanceUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBatcherBalanceUpdated is a log parse operation binding the contract event 0x0ea46f246ccfc58f7a93aa09bc6245a6818e97b1a160d186afe78993a3b194a0.
//
// Solidity: event BatcherBalanceUpdated(address indexed batcher, uint256 newBalance)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) ParseBatcherBalanceUpdated(log types.Log) (*ContractAlignedLayerServiceManagerBatcherBalanceUpdated, error) {
	event := new(ContractAlignedLayerServiceManagerBatcherBalanceUpdated)
	if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "BatcherBalanceUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractAlignedLayerServiceManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerInitializedIterator struct {
	Event *ContractAlignedLayerServiceManagerInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAlignedLayerServiceManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAlignedLayerServiceManagerInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAlignedLayerServiceManagerInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAlignedLayerServiceManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAlignedLayerServiceManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAlignedLayerServiceManagerInitialized represents a Initialized event raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractAlignedLayerServiceManagerInitializedIterator, error) {

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerInitializedIterator{contract: _ContractAlignedLayerServiceManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractAlignedLayerServiceManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAlignedLayerServiceManagerInitialized)
				if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) ParseInitialized(log types.Log) (*ContractAlignedLayerServiceManagerInitialized, error) {
	event := new(ContractAlignedLayerServiceManagerInitialized)
	if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractAlignedLayerServiceManagerNewBatchV2Iterator is returned from FilterNewBatchV2 and is used to iterate over the raw logs and unpacked data for NewBatchV2 events raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerNewBatchV2Iterator struct {
	Event *ContractAlignedLayerServiceManagerNewBatchV2 // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAlignedLayerServiceManagerNewBatchV2Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAlignedLayerServiceManagerNewBatchV2)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAlignedLayerServiceManagerNewBatchV2)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAlignedLayerServiceManagerNewBatchV2Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAlignedLayerServiceManagerNewBatchV2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAlignedLayerServiceManagerNewBatchV2 represents a NewBatchV2 event raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerNewBatchV2 struct {
	BatchMerkleRoot  [32]byte
	SenderAddress    common.Address
	TaskCreatedBlock uint32
	BatchDataPointer string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterNewBatchV2 is a free log retrieval operation binding the contract event 0x130d3e81af62e03ed6fff5e3bb343695ec513892cfad24d286486745dcc61437.
//
// Solidity: event NewBatchV2(bytes32 indexed batchMerkleRoot, address senderAddress, uint32 taskCreatedBlock, string batchDataPointer)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) FilterNewBatchV2(opts *bind.FilterOpts, batchMerkleRoot [][32]byte) (*ContractAlignedLayerServiceManagerNewBatchV2Iterator, error) {

	var batchMerkleRootRule []interface{}
	for _, batchMerkleRootItem := range batchMerkleRoot {
		batchMerkleRootRule = append(batchMerkleRootRule, batchMerkleRootItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.FilterLogs(opts, "NewBatchV2", batchMerkleRootRule)
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerNewBatchV2Iterator{contract: _ContractAlignedLayerServiceManager.contract, event: "NewBatchV2", logs: logs, sub: sub}, nil
}

// WatchNewBatchV2 is a free log subscription operation binding the contract event 0x130d3e81af62e03ed6fff5e3bb343695ec513892cfad24d286486745dcc61437.
//
// Solidity: event NewBatchV2(bytes32 indexed batchMerkleRoot, address senderAddress, uint32 taskCreatedBlock, string batchDataPointer)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) WatchNewBatchV2(opts *bind.WatchOpts, sink chan<- *ContractAlignedLayerServiceManagerNewBatchV2, batchMerkleRoot [][32]byte) (event.Subscription, error) {

	var batchMerkleRootRule []interface{}
	for _, batchMerkleRootItem := range batchMerkleRoot {
		batchMerkleRootRule = append(batchMerkleRootRule, batchMerkleRootItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.WatchLogs(opts, "NewBatchV2", batchMerkleRootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAlignedLayerServiceManagerNewBatchV2)
				if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "NewBatchV2", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewBatchV2 is a log parse operation binding the contract event 0x130d3e81af62e03ed6fff5e3bb343695ec513892cfad24d286486745dcc61437.
//
// Solidity: event NewBatchV2(bytes32 indexed batchMerkleRoot, address senderAddress, uint32 taskCreatedBlock, string batchDataPointer)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) ParseNewBatchV2(log types.Log) (*ContractAlignedLayerServiceManagerNewBatchV2, error) {
	event := new(ContractAlignedLayerServiceManagerNewBatchV2)
	if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "NewBatchV2", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractAlignedLayerServiceManagerNewBatchV3Iterator is returned from FilterNewBatchV3 and is used to iterate over the raw logs and unpacked data for NewBatchV3 events raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerNewBatchV3Iterator struct {
	Event *ContractAlignedLayerServiceManagerNewBatchV3 // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAlignedLayerServiceManagerNewBatchV3Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAlignedLayerServiceManagerNewBatchV3)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAlignedLayerServiceManagerNewBatchV3)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAlignedLayerServiceManagerNewBatchV3Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAlignedLayerServiceManagerNewBatchV3Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAlignedLayerServiceManagerNewBatchV3 represents a NewBatchV3 event raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerNewBatchV3 struct {
	BatchMerkleRoot       [32]byte
	SenderAddress         common.Address
	TaskCreatedBlock      uint32
	BatchDataPointer      string
	RespondToTaskFeeLimit *big.Int
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterNewBatchV3 is a free log retrieval operation binding the contract event 0x8801fc966deb2c8f563a103c35c9e80740585c292cd97518587e6e7927e6af55.
//
// Solidity: event NewBatchV3(bytes32 indexed batchMerkleRoot, address senderAddress, uint32 taskCreatedBlock, string batchDataPointer, uint256 respondToTaskFeeLimit)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) FilterNewBatchV3(opts *bind.FilterOpts, batchMerkleRoot [][32]byte) (*ContractAlignedLayerServiceManagerNewBatchV3Iterator, error) {

	var batchMerkleRootRule []interface{}
	for _, batchMerkleRootItem := range batchMerkleRoot {
		batchMerkleRootRule = append(batchMerkleRootRule, batchMerkleRootItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.FilterLogs(opts, "NewBatchV3", batchMerkleRootRule)
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerNewBatchV3Iterator{contract: _ContractAlignedLayerServiceManager.contract, event: "NewBatchV3", logs: logs, sub: sub}, nil
}

// WatchNewBatchV3 is a free log subscription operation binding the contract event 0x8801fc966deb2c8f563a103c35c9e80740585c292cd97518587e6e7927e6af55.
//
// Solidity: event NewBatchV3(bytes32 indexed batchMerkleRoot, address senderAddress, uint32 taskCreatedBlock, string batchDataPointer, uint256 respondToTaskFeeLimit)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) WatchNewBatchV3(opts *bind.WatchOpts, sink chan<- *ContractAlignedLayerServiceManagerNewBatchV3, batchMerkleRoot [][32]byte) (event.Subscription, error) {

	var batchMerkleRootRule []interface{}
	for _, batchMerkleRootItem := range batchMerkleRoot {
		batchMerkleRootRule = append(batchMerkleRootRule, batchMerkleRootItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.WatchLogs(opts, "NewBatchV3", batchMerkleRootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAlignedLayerServiceManagerNewBatchV3)
				if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "NewBatchV3", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewBatchV3 is a log parse operation binding the contract event 0x8801fc966deb2c8f563a103c35c9e80740585c292cd97518587e6e7927e6af55.
//
// Solidity: event NewBatchV3(bytes32 indexed batchMerkleRoot, address senderAddress, uint32 taskCreatedBlock, string batchDataPointer, uint256 respondToTaskFeeLimit)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) ParseNewBatchV3(log types.Log) (*ContractAlignedLayerServiceManagerNewBatchV3, error) {
	event := new(ContractAlignedLayerServiceManagerNewBatchV3)
	if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "NewBatchV3", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractAlignedLayerServiceManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerOwnershipTransferredIterator struct {
	Event *ContractAlignedLayerServiceManagerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAlignedLayerServiceManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAlignedLayerServiceManagerOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAlignedLayerServiceManagerOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAlignedLayerServiceManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAlignedLayerServiceManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAlignedLayerServiceManagerOwnershipTransferred represents a OwnershipTransferred event raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractAlignedLayerServiceManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerOwnershipTransferredIterator{contract: _ContractAlignedLayerServiceManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractAlignedLayerServiceManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAlignedLayerServiceManagerOwnershipTransferred)
				if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) ParseOwnershipTransferred(log types.Log) (*ContractAlignedLayerServiceManagerOwnershipTransferred, error) {
	event := new(ContractAlignedLayerServiceManagerOwnershipTransferred)
	if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractAlignedLayerServiceManagerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerPausedIterator struct {
	Event *ContractAlignedLayerServiceManagerPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAlignedLayerServiceManagerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAlignedLayerServiceManagerPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAlignedLayerServiceManagerPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAlignedLayerServiceManagerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAlignedLayerServiceManagerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAlignedLayerServiceManagerPaused represents a Paused event raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerPaused struct {
	Account         common.Address
	NewPausedStatus *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) FilterPaused(opts *bind.FilterOpts, account []common.Address) (*ContractAlignedLayerServiceManagerPausedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.FilterLogs(opts, "Paused", accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerPausedIterator{contract: _ContractAlignedLayerServiceManager.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *ContractAlignedLayerServiceManagerPaused, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.WatchLogs(opts, "Paused", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAlignedLayerServiceManagerPaused)
				if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) ParsePaused(log types.Log) (*ContractAlignedLayerServiceManagerPaused, error) {
	event := new(ContractAlignedLayerServiceManagerPaused)
	if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractAlignedLayerServiceManagerPauserRegistrySetIterator is returned from FilterPauserRegistrySet and is used to iterate over the raw logs and unpacked data for PauserRegistrySet events raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerPauserRegistrySetIterator struct {
	Event *ContractAlignedLayerServiceManagerPauserRegistrySet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAlignedLayerServiceManagerPauserRegistrySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAlignedLayerServiceManagerPauserRegistrySet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAlignedLayerServiceManagerPauserRegistrySet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAlignedLayerServiceManagerPauserRegistrySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAlignedLayerServiceManagerPauserRegistrySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAlignedLayerServiceManagerPauserRegistrySet represents a PauserRegistrySet event raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerPauserRegistrySet struct {
	PauserRegistry    common.Address
	NewPauserRegistry common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPauserRegistrySet is a free log retrieval operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) FilterPauserRegistrySet(opts *bind.FilterOpts) (*ContractAlignedLayerServiceManagerPauserRegistrySetIterator, error) {

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.FilterLogs(opts, "PauserRegistrySet")
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerPauserRegistrySetIterator{contract: _ContractAlignedLayerServiceManager.contract, event: "PauserRegistrySet", logs: logs, sub: sub}, nil
}

// WatchPauserRegistrySet is a free log subscription operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) WatchPauserRegistrySet(opts *bind.WatchOpts, sink chan<- *ContractAlignedLayerServiceManagerPauserRegistrySet) (event.Subscription, error) {

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.WatchLogs(opts, "PauserRegistrySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAlignedLayerServiceManagerPauserRegistrySet)
				if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "PauserRegistrySet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePauserRegistrySet is a log parse operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) ParsePauserRegistrySet(log types.Log) (*ContractAlignedLayerServiceManagerPauserRegistrySet, error) {
	event := new(ContractAlignedLayerServiceManagerPauserRegistrySet)
	if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "PauserRegistrySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractAlignedLayerServiceManagerRewardsInitiatorUpdatedIterator is returned from FilterRewardsInitiatorUpdated and is used to iterate over the raw logs and unpacked data for RewardsInitiatorUpdated events raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerRewardsInitiatorUpdatedIterator struct {
	Event *ContractAlignedLayerServiceManagerRewardsInitiatorUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAlignedLayerServiceManagerRewardsInitiatorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAlignedLayerServiceManagerRewardsInitiatorUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAlignedLayerServiceManagerRewardsInitiatorUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAlignedLayerServiceManagerRewardsInitiatorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAlignedLayerServiceManagerRewardsInitiatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAlignedLayerServiceManagerRewardsInitiatorUpdated represents a RewardsInitiatorUpdated event raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerRewardsInitiatorUpdated struct {
	PrevRewardsInitiator common.Address
	NewRewardsInitiator  common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterRewardsInitiatorUpdated is a free log retrieval operation binding the contract event 0xe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3.
//
// Solidity: event RewardsInitiatorUpdated(address prevRewardsInitiator, address newRewardsInitiator)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) FilterRewardsInitiatorUpdated(opts *bind.FilterOpts) (*ContractAlignedLayerServiceManagerRewardsInitiatorUpdatedIterator, error) {

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.FilterLogs(opts, "RewardsInitiatorUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerRewardsInitiatorUpdatedIterator{contract: _ContractAlignedLayerServiceManager.contract, event: "RewardsInitiatorUpdated", logs: logs, sub: sub}, nil
}

// WatchRewardsInitiatorUpdated is a free log subscription operation binding the contract event 0xe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3.
//
// Solidity: event RewardsInitiatorUpdated(address prevRewardsInitiator, address newRewardsInitiator)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) WatchRewardsInitiatorUpdated(opts *bind.WatchOpts, sink chan<- *ContractAlignedLayerServiceManagerRewardsInitiatorUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.WatchLogs(opts, "RewardsInitiatorUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAlignedLayerServiceManagerRewardsInitiatorUpdated)
				if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "RewardsInitiatorUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRewardsInitiatorUpdated is a log parse operation binding the contract event 0xe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3.
//
// Solidity: event RewardsInitiatorUpdated(address prevRewardsInitiator, address newRewardsInitiator)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) ParseRewardsInitiatorUpdated(log types.Log) (*ContractAlignedLayerServiceManagerRewardsInitiatorUpdated, error) {
	event := new(ContractAlignedLayerServiceManagerRewardsInitiatorUpdated)
	if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "RewardsInitiatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdateIterator is returned from FilterStaleStakesForbiddenUpdate and is used to iterate over the raw logs and unpacked data for StaleStakesForbiddenUpdate events raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdateIterator struct {
	Event *ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdate represents a StaleStakesForbiddenUpdate event raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdate struct {
	Value bool
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterStaleStakesForbiddenUpdate is a free log retrieval operation binding the contract event 0x40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc.
//
// Solidity: event StaleStakesForbiddenUpdate(bool value)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) FilterStaleStakesForbiddenUpdate(opts *bind.FilterOpts) (*ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdateIterator, error) {

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.FilterLogs(opts, "StaleStakesForbiddenUpdate")
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdateIterator{contract: _ContractAlignedLayerServiceManager.contract, event: "StaleStakesForbiddenUpdate", logs: logs, sub: sub}, nil
}

// WatchStaleStakesForbiddenUpdate is a free log subscription operation binding the contract event 0x40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc.
//
// Solidity: event StaleStakesForbiddenUpdate(bool value)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) WatchStaleStakesForbiddenUpdate(opts *bind.WatchOpts, sink chan<- *ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdate) (event.Subscription, error) {

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.WatchLogs(opts, "StaleStakesForbiddenUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdate)
				if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "StaleStakesForbiddenUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStaleStakesForbiddenUpdate is a log parse operation binding the contract event 0x40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc.
//
// Solidity: event StaleStakesForbiddenUpdate(bool value)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) ParseStaleStakesForbiddenUpdate(log types.Log) (*ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdate, error) {
	event := new(ContractAlignedLayerServiceManagerStaleStakesForbiddenUpdate)
	if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "StaleStakesForbiddenUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractAlignedLayerServiceManagerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerUnpausedIterator struct {
	Event *ContractAlignedLayerServiceManagerUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAlignedLayerServiceManagerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAlignedLayerServiceManagerUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAlignedLayerServiceManagerUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAlignedLayerServiceManagerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAlignedLayerServiceManagerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAlignedLayerServiceManagerUnpaused represents a Unpaused event raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerUnpaused struct {
	Account         common.Address
	NewPausedStatus *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) FilterUnpaused(opts *bind.FilterOpts, account []common.Address) (*ContractAlignedLayerServiceManagerUnpausedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.FilterLogs(opts, "Unpaused", accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerUnpausedIterator{contract: _ContractAlignedLayerServiceManager.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *ContractAlignedLayerServiceManagerUnpaused, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.WatchLogs(opts, "Unpaused", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAlignedLayerServiceManagerUnpaused)
				if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) ParseUnpaused(log types.Log) (*ContractAlignedLayerServiceManagerUnpaused, error) {
	event := new(ContractAlignedLayerServiceManagerUnpaused)
	if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractAlignedLayerServiceManagerVerifierDisabledIterator is returned from FilterVerifierDisabled and is used to iterate over the raw logs and unpacked data for VerifierDisabled events raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerVerifierDisabledIterator struct {
	Event *ContractAlignedLayerServiceManagerVerifierDisabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAlignedLayerServiceManagerVerifierDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAlignedLayerServiceManagerVerifierDisabled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAlignedLayerServiceManagerVerifierDisabled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAlignedLayerServiceManagerVerifierDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAlignedLayerServiceManagerVerifierDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAlignedLayerServiceManagerVerifierDisabled represents a VerifierDisabled event raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerVerifierDisabled struct {
	VerifierIdx uint8
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterVerifierDisabled is a free log retrieval operation binding the contract event 0xec54a85c01b5fc7fb41be0f33eabc56f2981110da8317b9817bc7c718f6d7bfe.
//
// Solidity: event VerifierDisabled(uint8 indexed verifierIdx)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) FilterVerifierDisabled(opts *bind.FilterOpts, verifierIdx []uint8) (*ContractAlignedLayerServiceManagerVerifierDisabledIterator, error) {

	var verifierIdxRule []interface{}
	for _, verifierIdxItem := range verifierIdx {
		verifierIdxRule = append(verifierIdxRule, verifierIdxItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.FilterLogs(opts, "VerifierDisabled", verifierIdxRule)
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerVerifierDisabledIterator{contract: _ContractAlignedLayerServiceManager.contract, event: "VerifierDisabled", logs: logs, sub: sub}, nil
}

// WatchVerifierDisabled is a free log subscription operation binding the contract event 0xec54a85c01b5fc7fb41be0f33eabc56f2981110da8317b9817bc7c718f6d7bfe.
//
// Solidity: event VerifierDisabled(uint8 indexed verifierIdx)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) WatchVerifierDisabled(opts *bind.WatchOpts, sink chan<- *ContractAlignedLayerServiceManagerVerifierDisabled, verifierIdx []uint8) (event.Subscription, error) {

	var verifierIdxRule []interface{}
	for _, verifierIdxItem := range verifierIdx {
		verifierIdxRule = append(verifierIdxRule, verifierIdxItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.WatchLogs(opts, "VerifierDisabled", verifierIdxRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAlignedLayerServiceManagerVerifierDisabled)
				if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "VerifierDisabled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVerifierDisabled is a log parse operation binding the contract event 0xec54a85c01b5fc7fb41be0f33eabc56f2981110da8317b9817bc7c718f6d7bfe.
//
// Solidity: event VerifierDisabled(uint8 indexed verifierIdx)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) ParseVerifierDisabled(log types.Log) (*ContractAlignedLayerServiceManagerVerifierDisabled, error) {
	event := new(ContractAlignedLayerServiceManagerVerifierDisabled)
	if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "VerifierDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractAlignedLayerServiceManagerVerifierEnabledIterator is returned from FilterVerifierEnabled and is used to iterate over the raw logs and unpacked data for VerifierEnabled events raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerVerifierEnabledIterator struct {
	Event *ContractAlignedLayerServiceManagerVerifierEnabled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractAlignedLayerServiceManagerVerifierEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAlignedLayerServiceManagerVerifierEnabled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractAlignedLayerServiceManagerVerifierEnabled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractAlignedLayerServiceManagerVerifierEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAlignedLayerServiceManagerVerifierEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAlignedLayerServiceManagerVerifierEnabled represents a VerifierEnabled event raised by the ContractAlignedLayerServiceManager contract.
type ContractAlignedLayerServiceManagerVerifierEnabled struct {
	VerifierIdx uint8
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterVerifierEnabled is a free log retrieval operation binding the contract event 0x5f52704e8e0190647930ccde0e43e14e89902d7d8c49c5f9e2544029f45ec12a.
//
// Solidity: event VerifierEnabled(uint8 indexed verifierIdx)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) FilterVerifierEnabled(opts *bind.FilterOpts, verifierIdx []uint8) (*ContractAlignedLayerServiceManagerVerifierEnabledIterator, error) {

	var verifierIdxRule []interface{}
	for _, verifierIdxItem := range verifierIdx {
		verifierIdxRule = append(verifierIdxRule, verifierIdxItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.FilterLogs(opts, "VerifierEnabled", verifierIdxRule)
	if err != nil {
		return nil, err
	}
	return &ContractAlignedLayerServiceManagerVerifierEnabledIterator{contract: _ContractAlignedLayerServiceManager.contract, event: "VerifierEnabled", logs: logs, sub: sub}, nil
}

// WatchVerifierEnabled is a free log subscription operation binding the contract event 0x5f52704e8e0190647930ccde0e43e14e89902d7d8c49c5f9e2544029f45ec12a.
//
// Solidity: event VerifierEnabled(uint8 indexed verifierIdx)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) WatchVerifierEnabled(opts *bind.WatchOpts, sink chan<- *ContractAlignedLayerServiceManagerVerifierEnabled, verifierIdx []uint8) (event.Subscription, error) {

	var verifierIdxRule []interface{}
	for _, verifierIdxItem := range verifierIdx {
		verifierIdxRule = append(verifierIdxRule, verifierIdxItem)
	}

	logs, sub, err := _ContractAlignedLayerServiceManager.contract.WatchLogs(opts, "VerifierEnabled", verifierIdxRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAlignedLayerServiceManagerVerifierEnabled)
				if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "VerifierEnabled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVerifierEnabled is a log parse operation binding the contract event 0x5f52704e8e0190647930ccde0e43e14e89902d7d8c49c5f9e2544029f45ec12a.
//
// Solidity: event VerifierEnabled(uint8 indexed verifierIdx)
func (_ContractAlignedLayerServiceManager *ContractAlignedLayerServiceManagerFilterer) ParseVerifierEnabled(log types.Log) (*ContractAlignedLayerServiceManagerVerifierEnabled, error) {
	event := new(ContractAlignedLayerServiceManagerVerifierEnabled)
	if err := _ContractAlignedLayerServiceManager.contract.UnpackLog(event, "VerifierEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

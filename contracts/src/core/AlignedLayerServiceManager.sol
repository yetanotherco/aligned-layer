// SPDX-License-Identifier: UNLICENSED
pragma solidity =0.8.12;

import {ServiceManagerBase, IAVSDirectory} from "eigenlayer-middleware/ServiceManagerBase.sol";
import {BLSSignatureChecker} from "eigenlayer-middleware/BLSSignatureChecker.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {Merkle} from "eigenlayer-core/contracts/libraries/Merkle.sol";
import {IRewardsCoordinator} from "eigenlayer-contracts/src/contracts/interfaces/IRewardsCoordinator.sol";
import {AlignedLayerServiceManagerStorage} from "./AlignedLayerServiceManagerStorage.sol";
import {IAlignedLayerServiceManager} from "./IAlignedLayerServiceManager.sol";

/**
 * @title Primary entrypoint for procuring services from Aligned.
 */
contract AlignedLayerServiceManager is
    IAlignedLayerServiceManager,
    ServiceManagerBase,
    BLSSignatureChecker,
    AlignedLayerServiceManagerStorage
{
    uint256 internal constant THRESHOLD_DENOMINATOR = 100;
    uint8 internal constant QUORUM_THRESHOLD_PERCENTAGE = 67;

    constructor(
        IAVSDirectory __avsDirectory,
        IRewardsCoordinator __rewardsCoordinator,
        IRegistryCoordinator __registryCoordinator,
        IStakeRegistry __stakeRegistry
    )
        BLSSignatureChecker(__registryCoordinator)
        ServiceManagerBase(
            __avsDirectory,
            __rewardsCoordinator,
            __registryCoordinator,
            __stakeRegistry
        )
    {
        _disableInitializers();
    }

    // @param _rewardsInitiator The address which is allowed to create AVS rewards submissions.
    function initialize(
        address _initialOwner,
        address _rewardsInitiator,
        address _alignedAggregator
    ) public initializer {
        __ServiceManagerBase_init(_initialOwner, _rewardsInitiator);
        alignedAggregator = _alignedAggregator; //can't do setAggregator(aggregator) since caller is not the owner
    }

    // This function is to be run only on upgrade
    // If a new contract is deployed, this function should be removed
    // Because this new value is also added in the initializer
    function initializeAggregator(
        address _alignedAggregator
    ) public reinitializer(2) {
        setAggregator(_alignedAggregator);
    }

    function createNewTask(
        bytes32 batchMerkleRoot,
        string calldata batchDataPointer
    ) external payable {
        bytes32 batchIdentifier = keccak256(
            abi.encodePacked(batchMerkleRoot, msg.sender)
        );

        if (batchesState[batchIdentifier].taskCreatedBlock != 0) {
            revert BatchAlreadySubmitted(batchIdentifier);
        }

        if (msg.value > 0) {
            batchersBalances[msg.sender] += msg.value;
            emit BatcherBalanceUpdated(
                msg.sender,
                batchersBalances[msg.sender]
            );
        }

        if (batchersBalances[msg.sender] == 0) {
            revert BatcherBalanceIsEmpty(msg.sender);
        }

        BatchState memory batchState;

        batchState.taskCreatedBlock = uint32(block.number);
        batchState.responded = false;

        batchesState[batchIdentifier] = batchState;

        emit NewBatchV2(
            batchMerkleRoot,
            msg.sender,
            uint32(block.number),
            batchDataPointer
        );
    }

    function respondToTaskV2(
        // (batchMerkleRoot,senderAddress) is signed as a way to verify the batch was right
        bytes32 batchMerkleRoot,
        address senderAddress,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external onlyAggregator {
        uint256 initialGasLeft = gasleft();

        bytes32 batchIdentifierHash = keccak256(
            abi.encodePacked(batchMerkleRoot, senderAddress)
        );

        /* CHECKING SIGNATURES & WHETHER THRESHOLD IS MET OR NOT */

        // Note: This is a hacky solidity way to see that the element exists
        // Value 0 would mean that the task is in block 0 so this can't happen.
        if (batchesState[batchIdentifierHash].taskCreatedBlock == 0) {
            revert BatchDoesNotExist(batchIdentifierHash);
        }

        // Check task hasn't been responsed yet
        if (batchesState[batchIdentifierHash].responded) {
            revert BatchAlreadyResponded(batchIdentifierHash);
        }

        if (batchersBalances[senderAddress] == 0) {
            revert BatcherHasNoBalance(senderAddress);
        }

        batchesState[batchIdentifierHash].responded = true;

        /* CHECKING SIGNATURES & WHETHER THRESHOLD IS MET OR NOT */
        // check that aggregated BLS signature is valid
        (QuorumStakeTotals memory quorumStakeTotals, ) = checkSignatures(
            batchIdentifierHash,
            batchesState[batchIdentifierHash].taskCreatedBlock,
            nonSignerStakesAndSignature
        );

        // check that signatories own at least a threshold percentage of each quourm
        if (
            quorumStakeTotals.signedStakeForQuorum[0] * THRESHOLD_DENOMINATOR <
            quorumStakeTotals.totalStakeForQuorum[0] *
                QUORUM_THRESHOLD_PERCENTAGE
        ) {
            revert InvalidQuorumThreshold(
                quorumStakeTotals.signedStakeForQuorum[0] *
                    THRESHOLD_DENOMINATOR,
                quorumStakeTotals.totalStakeForQuorum[0] *
                    QUORUM_THRESHOLD_PERCENTAGE
            );
        }

        emit BatchVerified(batchMerkleRoot, senderAddress);

        // Calculate estimation of gas used, check that batcher has sufficient funds
        // and send transaction cost to aggregator.
        uint256 finalGasLeft = gasleft();

        // 70k was measured by trial and error until the aggregator got paid a bit over what it needed
        uint256 txCost = (initialGasLeft - finalGasLeft + 70000) * tx.gasprice;

        if (batchersBalances[senderAddress] < txCost) {
            revert InsufficientFunds(
                senderAddress,
                txCost,
                batchersBalances[senderAddress]
            );
        }

        batchersBalances[senderAddress] -= txCost;
        emit BatcherBalanceUpdated(
            senderAddress,
            batchersBalances[senderAddress]
        );
        payable(alignedAggregator).transfer(txCost);
    }

    function verifyBatchInclusion(
        bytes32 proofCommitment,
        bytes32 pubInputCommitment,
        bytes32 provingSystemAuxDataCommitment,
        bytes20 proofGeneratorAddr,
        bytes32 batchMerkleRoot,
        bytes memory merkleProof,
        uint256 verificationDataBatchIndex,
        address senderAddress
    ) external view returns (bool) {
        bytes32 batchIdentifier;
        if (senderAddress == address(0)) {
            batchIdentifier = batchMerkleRoot;
        } else {
            batchIdentifier = keccak256(
                abi.encodePacked(batchMerkleRoot, senderAddress)
            );
        }

        if (batchesState[batchIdentifier].taskCreatedBlock == 0) {
            return false;
        }

        if (!batchesState[batchIdentifier].responded) {
            return false;
        }

        bytes memory leaf = abi.encodePacked(
            proofCommitment,
            pubInputCommitment,
            provingSystemAuxDataCommitment,
            proofGeneratorAddr
        );

        bytes32 hashedLeaf = keccak256(leaf);

        return
            Merkle.verifyInclusionKeccak(
                merkleProof,
                batchMerkleRoot,
                hashedLeaf,
                verificationDataBatchIndex
            );
    }

    // Old function signature for backwards compatibility
    function verifyBatchInclusion(
        bytes32 proofCommitment,
        bytes32 pubInputCommitment,
        bytes32 provingSystemAuxDataCommitment,
        bytes20 proofGeneratorAddr,
        bytes32 batchMerkleRoot,
        bytes memory merkleProof,
        uint256 verificationDataBatchIndex
    ) external view returns (bool) {
        return this.verifyBatchInclusion(
            proofCommitment,
            pubInputCommitment,
            provingSystemAuxDataCommitment,
            proofGeneratorAddr,
            batchMerkleRoot,
            merkleProof,
            verificationDataBatchIndex,
            address(0)
        );
    }

    function setAggregator(address _alignedAggregator) public onlyOwner {
        alignedAggregator = _alignedAggregator;
    }

    function withdraw(uint256 amount) external {
        if (batchersBalances[msg.sender] < amount) {
            revert InsufficientFunds(
                msg.sender,
                amount,
                batchersBalances[msg.sender]
            );
        }

        batchersBalances[msg.sender] -= amount;
        emit BatcherBalanceUpdated(msg.sender, batchersBalances[msg.sender]);

        payable(msg.sender).transfer(amount);
    }

    function balanceOf(address account) public view returns (uint256) {
        return batchersBalances[account];
    }

    function depositToBatcher(address account) external payable {
        _depositToBatcher(account, msg.value);
    }

    function _depositToBatcher(address account, uint256 amount) internal {
        if (amount == 0) {
            revert InvalidDepositAmount(amount);
        }
        batchersBalances[account] += amount;
        emit BatcherBalanceUpdated(account, batchersBalances[account]);
    }

    receive() external payable {
        _depositToBatcher(msg.sender, msg.value);
    }

    function checkPublicInput(
        bytes calldata publicInput,
        bytes32 hash
    ) public pure returns (bool) {
        return keccak256(publicInput) == hash;
    }

    modifier onlyAggregator() {
        if (msg.sender != alignedAggregator) {
            revert SenderIsNotAggregator(msg.sender, alignedAggregator);
        }
        _;
    }
}

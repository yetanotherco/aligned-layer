pragma solidity =0.8.12;

import {Initializable} from "@openzeppelin-upgrades/contracts/proxy/utils/Initializable.sol";
import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin-upgrades/contracts/security/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin-upgrades/contracts/proxy/utils/UUPSUpgradeable.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {IAlignedLayerServiceManager} from "./IAlignedLayerServiceManager.sol";

contract BatcherPaymentService is
    Initializable,
    OwnableUpgradeable,
    PausableUpgradeable,
    UUPSUpgradeable
{
    using ECDSA for bytes32;

    // CONSTANTS
    uint256 public constant UNLOCK_BLOCK_COUNT = 100;

    // EVENTS
    event PaymentReceived(address indexed sender, uint256 amount);
    event FundsWithdrawn(address indexed recipient, uint256 amount);

    // ERRORS
    error OnlyBatcherAllowed(address caller);
    error NoLeavesSubmitted();
    error NoProofSubmitterSignatures();
    error NotEnoughLeaves(uint256 leavesQty, uint256 signaturesQty);
    error LeavesNotPowerOfTwo(uint256 leavesQty);
    error NoGasForAggregator();
    error NoGasPerProof();
    error InsufficientGasForAggregator(uint256 required, uint256 available);
    error UserHasNoFundsToUnlock(address user);
    error UserHasNoFundsToLock(address user);
    error PayerInsufficientBalance(uint256 balance, uint256 amount);
    error FundsLocked(uint256 unlockBlock, uint256 currentBlock);
    error InvalidSignature();
    error InvalidNonce(uint256 expected, uint256 actual);
    error SignerInsufficientBalance(
        address signer,
        uint256 balance,
        uint256 required
    );
    error InvalidMerkleRoot(bytes32 expected, bytes32 actual);

    struct SignatureData {
        bytes signature;
        uint256 nonce;
    }

    struct UserInfo {
        uint256 balance;
        uint256 unlockBlock;
        uint256 nonce;
    }

    // STORAGE
    address public batcherWallet;

    IAlignedLayerServiceManager public alignedLayerServiceManager;

    // map to user data
    mapping(address => UserInfo) public userData;

    // storage gap for upgradeability
    uint256[24] private __GAP;

    // CONSTRUCTOR & INITIALIZER
    constructor() {
        _disableInitializers();
    }

    // MODIFIERS
    modifier onlyBatcher() {
        if (msg.sender != batcherWallet) {
            revert OnlyBatcherAllowed(msg.sender);
        }
        _;
    }

    function initialize(
        address _alignedLayerServiceManager,
        address _batcherPaymentServiceOwner,
        address _batcherWallet
    ) public initializer {
        __Ownable_init(); // default is msg.sender
        __UUPSUpgradeable_init();
        _transferOwnership(_batcherPaymentServiceOwner);

        alignedLayerServiceManager = IAlignedLayerServiceManager(
            _alignedLayerServiceManager
        );
        batcherWallet = _batcherWallet;
    }

    // PAYABLE FUNCTIONS
    receive() external payable {
        userData[msg.sender].balance += msg.value;
        emit PaymentReceived(msg.sender, msg.value);
    }

    // PUBLIC FUNCTIONS
    function createNewTask(
        bytes32 batchMerkleRoot,
        string calldata batchDataPointer,
        bytes32[] calldata leaves, // padded to the next power of 2
        SignatureData[] calldata signatures, // actual length (proof sumbitters == proofs submitted)
        uint256 gasForAggregator,
        uint256 gasPerProof
    ) external onlyBatcher whenNotPaused {
        uint256 leavesQty = leaves.length;
        uint256 signaturesQty = signatures.length;

        uint256 feeForAggregator = gasForAggregator * tx.gasprice;
        uint256 feePerProof = gasPerProof * tx.gasprice;

        if (leavesQty <= 0) {
            revert NoLeavesSubmitted();
        }

        if (signaturesQty <= 0) {
            revert NoProofSubmitterSignatures();
        }

        if (leavesQty < signaturesQty) {
            revert NotEnoughLeaves(leavesQty, signaturesQty);
        }

        if ((leavesQty & (leavesQty - 1)) != 0) {
            revert LeavesNotPowerOfTwo(leavesQty);
        }

        if (feeForAggregator <= 0) {
            revert NoGasForAggregator();
        }

        if (feePerProof <= 0) {
            revert NoGasPerProof();
        }

        if (feePerProof * signaturesQty <= feeForAggregator) {
            revert InsufficientGasForAggregator(
                feeForAggregator,
                feePerProof * signaturesQty
            );
        }

        _checkMerkleRootAndVerifySignatures(
            leaves,
            batchMerkleRoot,
            signatures,
            feePerProof
        );

        // call alignedLayerServiceManager
        // with value to fund the task's response
        alignedLayerServiceManager.createNewTask{value: feeForAggregator}(
            batchMerkleRoot,
            batchDataPointer
        );

        payable(batcherWallet).transfer(
            (feePerProof * signaturesQty) - feeForAggregator
        );
    }

    function unlock() external whenNotPaused {
        if (userData[msg.sender].balance <= 0) {
            revert UserHasNoFundsToUnlock(msg.sender);
        }

        userData[msg.sender].unlockBlock = block.number + UNLOCK_BLOCK_COUNT;
    }

    function lock() external whenNotPaused {
        if (userData[msg.sender].balance <= 0) {
            revert UserHasNoFundsToLock(msg.sender);
        }
        userData[msg.sender].unlockBlock = 0;
    }

    function withdraw(uint256 amount) external whenNotPaused {
        UserInfo storage user_data = userData[msg.sender];
        if (user_data.balance < amount) {
            revert PayerInsufficientBalance(user_data.balance, amount);
        }

        if (
            user_data.unlockBlock == 0 || user_data.unlockBlock > block.number
        ) {
            revert FundsLocked(user_data.unlockBlock, block.number);
        }

        user_data.balance -= amount;
        payable(msg.sender).transfer(amount);
        emit FundsWithdrawn(msg.sender, amount);
    }

    function pause() public onlyOwner {
        _pause();
    }

    function unpause() public onlyOwner {
        _unpause();
    }

    function _authorizeUpgrade(
        address newImplementation
    ) internal override onlyOwner {}

    function _checkMerkleRootAndVerifySignatures(
        bytes32[] calldata leaves,
        bytes32 batchMerkleRoot,
        SignatureData[] calldata signatures,
        uint256 feePerProof
    ) private {
        uint256 numNodesInLayer = leaves.length / 2;
        bytes32[] memory layer = new bytes32[](numNodesInLayer);

        uint32 i = 0;

        // Calculate the hash of the next layer of the Merkle tree
        // and verify the signatures up to numNodesInLayer
        for (i = 0; i < numNodesInLayer; i++) {
            layer[i] = keccak256(
                abi.encodePacked(leaves[2 * i], leaves[2 * i + 1])
            );

            _verifySignatureAndDecreaseBalance(
                leaves[i],
                signatures[i],
                feePerProof
            );
        }

        // Verify the rest of the signatures
        for (; i < signatures.length; i++) {
            _verifySignatureAndDecreaseBalance(
                leaves[i],
                signatures[i],
                feePerProof
            );
        }

        // The next layer above has half as many nodes
        numNodesInLayer /= 2;

        // Continue calculating Merkle root for remaining layers
        while (numNodesInLayer != 0) {
            // Overwrite the first numNodesInLayer nodes in layer with the pairwise hashes of their children
            for (i = 0; i < numNodesInLayer; i++) {
                layer[i] = keccak256(
                    abi.encodePacked(layer[2 * i], layer[2 * i + 1])
                );
            }

            // The next layer above has half as many nodes
            numNodesInLayer /= 2;
        }

        if (leaves.length == 1) {
            if (leaves[0] != batchMerkleRoot) {
                revert InvalidMerkleRoot(batchMerkleRoot, leaves[0]);
            }
        } else {
            if (layer[0] != batchMerkleRoot) {
                revert InvalidMerkleRoot(batchMerkleRoot, layer[0]);
            }
        }
    }

    function _verifySignatureAndDecreaseBalance(
        bytes32 hash,
        SignatureData calldata signatureData,
        uint256 feePerProof
    ) private {
        bytes32 noncedHash = keccak256(
            abi.encodePacked(hash, signatureData.nonce)
        );

        address signer = noncedHash.recover(signatureData.signature);

        if (signer == address(0)) {
            revert InvalidSignature();
        }

        UserInfo storage user_data = userData[signer];

        if (user_data.nonce != signatureData.nonce) {
            revert InvalidNonce(user_data.nonce, signatureData.nonce);
        }
        user_data.nonce++;

        if (user_data.balance < feePerProof) {
            revert SignerInsufficientBalance(
                signer,
                user_data.balance,
                feePerProof
            );
        }

        user_data.balance -= feePerProof;
    }

    function user_balances(address account) public view returns (uint256) {
        return userData[account].balance;
    }

    function user_nonces(address account) public view returns (uint256) {
        return userData[account].nonce;
    }

    function user_unlock_block(address account) public view returns (uint256) {
        return userData[account].unlockBlock;
    }
}

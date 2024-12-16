// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";

contract ClaimableAirdrop is
    ReentrancyGuard,
    Initializable,
    PausableUpgradeable,
    Ownable2StepUpgradeable
{
    /// @notice Address of the token contract to claim.
    address public tokenProxy;

    /// @notice Address of the wallet that has the tokens to distribute to the claimants.
    address public claimSupplier;

    /// @notice Timestamp until which the claimants can claim the tokens.
    uint256 public limitTimestampToClaim;

    /// @notice Merkle root of the claimants.
    bytes32 public claimMerkleRoot;

    /// @notice Mapping of the claimants that have claimed the tokens.
    /// @dev true if the claimant has claimed the tokens.
    mapping(address => bool) public hasClaimed;

    /// @notice Event emitted when a claimant claims the tokens.
    event TokenClaimed(address indexed to, uint256 indexed amount);

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @notice Initializes the contract.
    /// @dev This initializer should be called only once.
    /// @param _owner address of the owner of the token.
    /// @param _tokenProxy address of the token contract.
    /// @param _claimSupplier address of the wallet that has the tokens to distribute to the claimants.
    /// @param _limitTimestampToClaim timestamp until which the claimants can claim the tokens.
    /// @param _claimMerkleRoot Merkle root of the claimants.
    function initialize(
        address _owner,
        address _tokenProxy,
        address _claimSupplier,
        uint256 _limitTimestampToClaim,
        bytes32 _claimMerkleRoot
    ) public initializer {
        require(_owner != address(0), "Invalid owner address");
        require(
            _tokenProxy != address(0) && _tokenProxy != address(this),
            "Invalid token contract address"
        );
        require(
            _claimSupplier != address(0) && _claimSupplier != address(this),
            "Invalid token owner address"
        );
        require(_limitTimestampToClaim > block.timestamp, "Invalid timestamp");
        require(_claimMerkleRoot != 0, "Invalid Merkle root");

        __Ownable2Step_init(); // default is msg.sender
        _transferOwnership(_owner);
        __Pausable_init();

        tokenProxy = _tokenProxy;
        claimSupplier = _claimSupplier;
        limitTimestampToClaim = _limitTimestampToClaim;
        claimMerkleRoot = _claimMerkleRoot;
    }

    /// @notice Claim the tokens.
    /// @param amount amount of tokens to claim.
    /// @param merkleProof Merkle proof of the claim.
    function claim(
        uint256 amount,
        bytes32[] calldata merkleProof
    ) public nonReentrant whenNotPaused {
        require(
            !hasClaimed[msg.sender],
            "Account has already claimed the drop"
        );
        require(
            block.timestamp <= limitTimestampToClaim,
            "Drop is no longer claimable"
        );

        bytes32 leaf = keccak256(
            bytes.concat(keccak256(abi.encode(msg.sender, amount)))
        );
        bool verifies = MerkleProof.verify(merkleProof, claimMerkleRoot, leaf);

        require(verifies, "Invalid Merkle proof");

        require(
            IERC20(tokenProxy).allowance(claimSupplier, address(this)) >=
                amount,
            "Insufficient token allowance"
        );

        bool success = IERC20(tokenProxy).transferFrom(
            claimSupplier,
            msg.sender,
            amount
        );

        require(success, "Failed to transfer funds");

        hasClaimed[msg.sender] = true;

        emit TokenClaimed(msg.sender, amount);
    }

    /// @notice Update the Merkle root.
    /// @param newRoot new Merkle root.
    function updateMerkleRoot(bytes32 newRoot) external whenPaused onlyOwner {
        require(newRoot != 0 && newRoot != claimMerkleRoot, "Invalid root");
        claimMerkleRoot = newRoot;
    }

    /// @notice Extend the claim period.
    /// @param newTimestamp new timestamp until which the claimants can claim the tokens.
    function extendClaimPeriod(
        uint256 newTimestamp
    ) external whenPaused onlyOwner {
        require(
            newTimestamp > limitTimestampToClaim &&
                newTimestamp > block.timestamp,
            "Can only extend from current timestamp"
        );
        limitTimestampToClaim = newTimestamp;
    }

    /// @notice Pause the contract.
    function pause() public whenNotPaused onlyOwner {
        _pause();
    }

    /// @notice Unpause the contract.
    function unpause() public whenPaused onlyOwner {
        _unpause();
    }
}

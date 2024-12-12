// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title ClaimableAirdropV2Example contract
/// @notice This is an example of how a ClaimableAirdrop contract upgrade
/// should look like
/// @dev Minimum necessary changes:
/// - The initialize function is marked as reinitializer with version 2.
/// @custom:oz-upgrades-from ClaimableAirdropV1
contract ClaimableAirdropV2Example is
    ReentrancyGuard,
    Initializable,
    OwnableUpgradeable,
    UUPSUpgradeable
{
    address public tokenContractAddress;
    address public tokenOwnerAddress;
    uint256 public limitTimestampToClaim;
    bytes32 public claimMerkleRoot;

    mapping(address => bool) public hasClaimed;

    event TokenClaimed(address indexed to, uint256 indexed amount);

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(
        address _initialOwner,
        address _tokenContractAddress,
        address _tokenOwnerAddress,
        uint256 _limitTimestampToClaim,
        bytes32 _claimMerkleRoot
    ) public initializer nonReentrant {
        __Ownable_init(_initialOwner);
        __UUPSUpgradeable_init();

        tokenContractAddress = _tokenContractAddress;
        tokenOwnerAddress = _tokenOwnerAddress;
        limitTimestampToClaim = _limitTimestampToClaim;
        claimMerkleRoot = _claimMerkleRoot;
    }

    function claim(
        uint256 amount,
        bytes32[] calldata merkleProof
    ) public nonReentrant {
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

        hasClaimed[msg.sender] = true;

        bool success = IERC20(tokenContractAddress).transferFrom(
            tokenOwnerAddress,
            msg.sender,
            amount
        );

        require(success, "Failed to transfer funds");

        emit TokenClaimed(msg.sender, amount);
    }

    function _authorizeUpgrade(
        address newImplementation
    ) internal override onlyOwner {}
}

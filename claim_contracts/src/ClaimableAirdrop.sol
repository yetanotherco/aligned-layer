// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

contract ClaimableAirdrop is
    ReentrancyGuard,
    Initializable,
    PausableUpgradeable,
    OwnableUpgradeable
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
        address _multisig,
        address _tokenContractAddress,
        address _tokenOwnerAddress,
        uint256 _limitTimestampToClaim,
        bytes32 _claimMerkleRoot
    ) public initializer {
        __Ownable_init(_multisig);
        __Pausable_init();

        require(_multisig != address(0), "Invalid multisig address");
        require(
            _tokenContractAddress != address(0),
            "Invalid token contract address"
        );
        require(
            _tokenOwnerAddress != address(0),
            "Invalid token owner address"
        );
        require(_limitTimestampToClaim > block.timestamp, "Invalid timestamp");
        require(_claimMerkleRoot != 0, "Invalid Merkle root");

        tokenContractAddress = _tokenContractAddress;
        tokenOwnerAddress = _tokenOwnerAddress;
        limitTimestampToClaim = _limitTimestampToClaim;
        claimMerkleRoot = _claimMerkleRoot;
    }

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
            IERC20(tokenContractAddress).allowance(
                tokenOwnerAddress,
                address(this)
            ) >= amount,
            "Insufficient token allowance"
        );

        bool success = IERC20(tokenContractAddress).transferFrom(
            tokenOwnerAddress,
            msg.sender,
            amount
        );

        require(success, "Failed to transfer funds");

        hasClaimed[msg.sender] = true;

        emit TokenClaimed(msg.sender, amount);
    }

    function pause() public onlyOwner {
        _pause();
    }

    function unpause() public onlyOwner {
        _unpause();
    }
}

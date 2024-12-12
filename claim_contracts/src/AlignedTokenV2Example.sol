// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20PermitUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title AlignedTokenV2Example contract
/// @notice This is an example of how a AlignedToken contract upgrade
/// should look like.
/// @dev Minimum necessary changes:
/// - The initialize function is marked as reinitializer with version 2.
/// @custom:oz-upgrades-from AlignedTokenV1
contract AlignedTokenV2Example is
    Initializable,
    ERC20Upgradeable,
    OwnableUpgradeable,
    ERC20PermitUpgradeable,
    UUPSUpgradeable
{
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(
        address initialOwner,
        address beneficiary1,
        uint256 beneficiary1Part,
        address beneficiary2,
        uint256 beneficiary2Part,
        address beneficiary3,
        uint256 beneficiary3Part
    ) public initializer {
        __ERC20_init("AlignedTokenV2", "ALI");
        __Ownable_init(initialOwner);
        __ERC20Permit_init("AlignedTokenV2");
        __UUPSUpgradeable_init();

        _mint(beneficiary1, beneficiary1Part);
        _mint(beneficiary2, beneficiary2Part);
        _mint(beneficiary3, beneficiary3Part);
    }

    function _authorizeUpgrade(
        address newImplementation
    ) internal override onlyOwner {}
}

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20PermitUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

contract AlignedTokenV1 is
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
        address _initialOwner,
        address _beneficiary1,
        uint256 _beneficiary1Part,
        address _beneficiary2,
        uint256 _beneficiary2Part,
        address _beneficiary3,
        uint256 _beneficiary3Part
    ) public initializer {
        __ERC20_init("AlignedTokenV1", "ALI");
        __Ownable_init(_initialOwner);
        __ERC20Permit_init("AlignedTokenV1");
        __UUPSUpgradeable_init();

        _mint(_beneficiary1, _beneficiary1Part);
        _mint(_beneficiary2, _beneficiary2Part);
        _mint(_beneficiary3, _beneficiary3Part);
    }

    function _authorizeUpgrade(
        address newImplementation
    ) internal override onlyOwner {}
}

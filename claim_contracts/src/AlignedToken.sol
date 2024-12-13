// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";

contract AlignedToken is Initializable, ERC20Upgradeable {
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(
        address _beneficiary1,
        uint256 _beneficiary1Part,
        address _beneficiary2,
        uint256 _beneficiary2Part,
        address _beneficiary3,
        uint256 _beneficiary3Part
    ) public initializer {
        __ERC20_init("AlignedToken", "ALI");
        _mint(_beneficiary1, _beneficiary1Part);
        _mint(_beneficiary2, _beneficiary2Part);
        _mint(_beneficiary3, _beneficiary3Part);
    }
}

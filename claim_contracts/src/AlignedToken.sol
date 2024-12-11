// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20PermitUpgradeable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

contract AlignedToken is Initializable, ERC20Upgradeable, ReentrancyGuard {
    struct Supply {
        address beneficiary;
        uint256 amount;
    }

    constructor() {
        // Ensure that initialization methods are run only once.
        // This is a safeguard against accidental reinitialization.
        _disableInitializers();
    }

    function initialize(
        address beneficiary1,
        uint256 amount1,
        address beneficiary2,
        uint256 amount2,
        address beneficiary3,
        uint256 amount3
    ) public initializer nonReentrant {
        __ERC20_init("AlignedToken", "ALI");
        _mint(beneficiary1, amount1);
        _mint(beneficiary2, amount2);
        _mint(beneficiary3, amount3);
    }
}

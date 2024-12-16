// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20PermitUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";

contract AlignedToken is
    Initializable,
    ERC20Upgradeable,
    EIP712Upgradeable,
    ERC20PermitUpgradeable,
    Ownable2StepUpgradeable
{
    /// @notice Name of the token.
    string public constant NAME = "Aligned Token";

    /// @notice Symbol of the token.
    string public constant SYMBOL = "ALIGN";

    /// @notice Version of the token implementation.
    string public constant VERSION = "1";

    /// @notice Total supply of the token.
    uint256 public constant TOTAL_SUPPLY = 10_000_000_000e18; // 10 billion

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @notice Initializes the contract.
    /// @dev This initializer should be called only once.
    /// @param _owner address of the owner of the token.
    /// @param _foundation address of the foundation.
    /// @param _claim address of the claim.
    function initialize(
        address _owner,
        address _foundation,
        address _claim
    ) public initializer {
        require(
            _foundation != address(0) &&
                _claim != address(0) &&
                _owner != address(0),
            "Invalid _foundation or _claim or _owner"
        );
        __ERC20_init(NAME, SYMBOL);
        __EIP712_init(NAME, VERSION);
        __ERC20Permit_init(NAME);
        __Ownable2Step_init(); // default is msg.sender
        _transferOwnership(_owner);
        _mint(_foundation, 7_300_000_000e18); // 7.3 billion
        _mint(_claim, 2_700_000_000e18); // 2.7 billion
    }

    /// @notice Prevents the owner from renouncing ownership.
    function renounceOwnership() public view override onlyOwner {
        revert("Cannot renounce ownership");
    }
}

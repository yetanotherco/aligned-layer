// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20PermitUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20BurnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";

contract AlignedToken is
    Initializable,
    ERC20Upgradeable,
    ERC20PermitUpgradeable,
    ERC20BurnableUpgradeable,
    Ownable2StepUpgradeable
{
    /// @notice Name of the token.
    string public constant NAME = "Aligned Token";

    /// @notice Symbol of the token.
    string public constant SYMBOL = "ALIGN";

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @notice Initializes the contract.
    /// @dev This initializer should be called only once.
    /// @param _foundation address of the foundation.
    /// @param _claimSupplier address of the claim supplier. This is the address
    /// that will give the tokens to the users that claim them.
    function initialize(
        address _foundation,
        address _claimSupplier
    ) public initializer {
        require(
            _foundation != address(0) && _claimSupplier != address(0),
            "Invalid _foundation or _claimSupplier"
        );
        __ERC20_init(NAME, SYMBOL);
        __ERC20Permit_init(NAME);
        __ERC20Burnable_init();
        __Ownable2Step_init(); // default is msg.sender
        _transferOwnership(_foundation);
        _mint(_foundation, 7_400_000_000e18); // 7.4 billion
        _mint(_claimSupplier, 2_600_000_000e18); // 2.6 billion
    }

    /// @notice Mints `amount` of tokens.
    function mint(address to, uint256 amount) external onlyOwner {
        _mint(to, amount);
    }

    /// @notice Prevents the owner from renouncing ownership.
    function renounceOwnership() public view override onlyOwner {
        revert("Cannot renounce ownership");
    }
}

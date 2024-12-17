// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20PermitUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20BurnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";

/// @title Aligned Token
/// @notice This contract is the implementation of the Aligned Token
/// @dev This contract is upgradeable and should be used only through the proxy contract
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

    /// @notice Supply of the token for the foundation.
    uint256 public constant FOUNDATION_SUPPLY = 7_400_000_000e18; // 7.4 billion

    /// @notice Supply of the token for the token distributor.
    uint256 public constant TOKEN_DISTRIBUTOR_SUPPLY = 2_600_000_000e18; // 2.6 billion

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @notice Initializes the contract.
    /// @dev This initializer should be called only once.
    /// @param _foundation address of the foundation.
    /// @param _tokenDistributor address of the token distributor. This is the address
    /// that will give the tokens to the users that claim them.
    function initialize(
        address _foundation,
        address _tokenDistributor
    ) external initializer {
        require(
            _foundation != address(0) && _tokenDistributor != address(0),
            "Invalid _foundation or _tokenDistributor"
        );
        __ERC20_init(NAME, SYMBOL);
        __ERC20Permit_init(NAME);
        __ERC20Burnable_init();
        __Ownable_init(_foundation);
        _mint(_foundation, FOUNDATION_SUPPLY);
        _mint(_tokenDistributor, TOKEN_DISTRIBUTOR_SUPPLY);
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

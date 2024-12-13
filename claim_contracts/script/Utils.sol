// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Vm} from "forge-std/Vm.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "../src/AlignedTokenV1.sol";
import "../src/ClaimableAirdropV1.sol";

library Utils {
    // Cheatcodes address, 0x7109709ECfa91a80626fF3989D68f67F5b1DD12D.
    address internal constant VM_ADDRESS =
        address(uint160(uint256(keccak256("hevm cheat code"))));
    Vm internal constant vm = Vm(VM_ADDRESS);

    /// @notice Address of the deterministic create2 factory.
    /// @dev This address corresponds to a contracts that is set in the storage
    /// in the genesis file. The same contract with the same address is deployed
    /// in every testnet, so if this script is run in a testnet instead of in a
    /// local environment, it should work.
    address constant DETERMINISTIC_CREATE2_ADDRESS =
        0x4e59b44847b379578588920cA78FbF26c0B4956C;

    function deployWithCreate2(
        bytes memory bytecode,
        bytes32 salt,
        address create2Factory
    ) internal returns (address) {
        if (bytecode.length == 0) {
            revert("Bytecode is not set");
        }
        address contractAddress = vm.computeCreate2Address(
            salt,
            keccak256(bytecode),
            create2Factory
        );
        if (contractAddress.code.length != 0) {
            revert("Contract already deployed");
        }

        vm.broadcast();
        (bool success, bytes memory data) = create2Factory.call(
            abi.encodePacked(salt, bytecode)
        );
        contractAddress = bytesToAddress(data);

        if (!success) {
            revert(
                "Failed to deploy contract via create2: create2Factory call failed"
            );
        }

        if (contractAddress == address(0)) {
            revert(
                "Failed to deploy contract via create2: contract address is zero"
            );
        }

        if (contractAddress.code.length == 0) {
            revert(
                "Failed to deploy contract via create2: contract code is empty"
            );
        }

        return contractAddress;
    }

    function bytesToAddress(
        bytes memory addressOffset
    ) internal pure returns (address addr) {
        assembly {
            addr := mload(add(addressOffset, 20))
        }
    }

    // AlignedToken utils

    function deployAlignedTokenImplementation() internal returns (address) {
        vm.broadcast();
        AlignedTokenV1 _implementation = new AlignedTokenV1();
        return address(_implementation);
    }

    function alignedTokenProxyDeploymentData(
        address _proxyAdmin,
        address _implementation,
        address _beneficiary1,
        address _beneficiary2,
        address _beneficiary3,
        uint256 _mintAmount
    ) internal pure returns (bytes memory) {
        return
            abi.encodePacked(
                type(TransparentUpgradeableProxy).creationCode,
                abi.encode(
                    _implementation,
                    _proxyAdmin,
                    alignedTokenInitData(
                        _implementation,
                        _beneficiary1,
                        _beneficiary2,
                        _beneficiary3,
                        _mintAmount
                    )
                )
            );
    }

    function alignedTokenInitData(
        address _implementation,
        address _beneficiary1,
        address _beneficiary2,
        address _beneficiary3,
        uint256 _mintAmount
    ) internal pure returns (bytes memory) {
        return
            abi.encodeCall(
                AlignedTokenV1(_implementation).initialize,
                (
                    _beneficiary1,
                    _mintAmount,
                    _beneficiary2,
                    _mintAmount,
                    _beneficiary3,
                    _mintAmount
                )
            );
    }

    function alignedTokenUpgradeData(
        address _proxyAdmin,
        address _proxy,
        address _newImplementation,
        address _beneficiary1,
        address _beneficiary2,
        address _beneficiary3,
        uint256 _mintAmount
    ) internal pure returns (bytes memory) {
        return
            abi.encodeCall(
                ProxyAdmin(_proxyAdmin).upgradeAndCall,
                (
                    ITransparentUpgradeableProxy(_proxy),
                    _newImplementation,
                    Utils.alignedTokenInitData(
                        _newImplementation,
                        _beneficiary1,
                        _beneficiary2,
                        _beneficiary3,
                        _mintAmount
                    )
                )
            );
    }

    // ClaimableAirdrop utils

    function deployClaimableAirdropImplementation() internal returns (address) {
        vm.broadcast();
        ClaimableAirdropV1 _implementation = new ClaimableAirdropV1();
        return address(_implementation);
    }

    function claimableAirdropProxyDeploymentData(
        address _proxyAdmin,
        address _implementation,
        address _owner,
        address _tokenContractAddress,
        address _tokenOwnerAddress,
        uint256 _limitTimestampToClaim,
        bytes32 _claimMerkleRoot
    ) internal pure returns (bytes memory) {
        return
            abi.encodePacked(
                type(TransparentUpgradeableProxy).creationCode,
                abi.encode(
                    _implementation,
                    _proxyAdmin,
                    claimableAirdropInitData(
                        _implementation,
                        _owner,
                        _tokenContractAddress,
                        _tokenOwnerAddress,
                        _limitTimestampToClaim,
                        _claimMerkleRoot
                    )
                )
            );
    }

    function claimableAirdropInitData(
        address _implementation,
        address _owner,
        address _tokenContractAddress,
        address _tokenOwnerAddress,
        uint256 _limitTimestampToClaim,
        bytes32 _claimMerkleRoot
    ) internal pure returns (bytes memory) {
        return
            abi.encodeCall(
                ClaimableAirdropV1(_implementation).initialize,
                (
                    _owner,
                    _tokenContractAddress,
                    _tokenOwnerAddress,
                    _limitTimestampToClaim,
                    _claimMerkleRoot
                )
            );
    }

    function claimableAirdropUpgradeData(
        address _proxy,
        address _owner,
        address _newImplementation,
        address _tokenContractAddress,
        address _tokenOwnerAddress,
        uint256 _limitTimestampToClaim,
        bytes32 _claimMerkleRoot
    ) internal pure returns (bytes memory) {
        return
            abi.encodeCall(
                ProxyAdmin(_newImplementation).upgradeAndCall,
                (
                    ITransparentUpgradeableProxy(_proxy),
                    _newImplementation,
                    Utils.claimableAirdropInitData(
                        _owner,
                        _newImplementation,
                        _tokenContractAddress,
                        _tokenOwnerAddress,
                        _limitTimestampToClaim,
                        _claimMerkleRoot
                    )
                )
            );
    }

    // ProxyAdmin utils

    function deployProxyAdmin(address _safe) internal returns (address) {
        vm.broadcast();
        ProxyAdmin _proxyAdmin = new ProxyAdmin(_safe);
        return address(_proxyAdmin);
    }

    function proxyAdminDeploymentData(
        address _safe
    ) internal pure returns (bytes memory) {
        return
            abi.encodePacked(type(ProxyAdmin).creationCode, abi.encode(_safe));
    }
}

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Vm} from "forge-std/Vm.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import "../src/AlignedTokenV1.sol";
import "../src/AlignedTokenV2Example.sol";
import "../src/ClaimableAirdropV1.sol";
import "../src/ClaimableAirdropV2Example.sol";

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
        address create2Factory,
        uint256 signerPrivateKey
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

        vm.broadcast(signerPrivateKey);
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

    function deployAlignedTokenImplementation(
        uint256 _version
    ) internal returns (address) {
        address _implementation_address;
        if (_version == 1) {
            vm.broadcast();
            AlignedTokenV1 _implementation = new AlignedTokenV1();
            _implementation_address = address(_implementation);
        } else if (_version == 2) {
            vm.broadcast();
            AlignedTokenV2Example _implementation = new AlignedTokenV2Example();
            _implementation_address = address(_implementation);
        } else {
            revert("Unsupported version");
        }
        return _implementation_address;
    }

    function alignedTokenProxyDeploymentData(
        address _implementation,
        uint256 _version,
        address _safe,
        address _beneficiary1,
        address _beneficiary2,
        address _beneficiary3,
        uint256 _mintAmount
    ) internal pure returns (bytes memory) {
        return
            abi.encodePacked(
                type(ERC1967Proxy).creationCode,
                alignedTokenInitData(
                    _implementation,
                    _version,
                    _safe,
                    _beneficiary1,
                    _beneficiary2,
                    _beneficiary3,
                    _mintAmount
                )
            );
    }

    function alignedTokenInitData(
        address _implementation,
        uint256 _version,
        address _safe,
        address _beneficiary1,
        address _beneficiary2,
        address _beneficiary3,
        uint256 _mintAmount
    ) internal pure returns (bytes memory) {
        if (_version == 1) {
            return
                abi.encodeCall(
                    AlignedTokenV1(_implementation).initialize,
                    (
                        _safe,
                        _beneficiary1,
                        _mintAmount,
                        _beneficiary2,
                        _mintAmount,
                        _beneficiary3,
                        _mintAmount
                    )
                );
        } else if (_version == 2) {
            return
                abi.encodeCall(
                    AlignedTokenV2Example(_implementation).initialize,
                    (
                        _safe,
                        _beneficiary1,
                        _mintAmount,
                        _beneficiary2,
                        _mintAmount,
                        _beneficiary3,
                        _mintAmount
                    )
                );
        } else {
            revert("Unsupported version");
        }
    }

    function alignedTokenUpgradeData(
        address _newImplementation,
        uint256 _version,
        address _safe,
        address _beneficiary1,
        address _beneficiary2,
        address _beneficiary3,
        uint256 _mintAmount
    ) internal pure returns (bytes memory) {
        bytes memory _initData = Utils.alignedTokenInitData(
            _newImplementation,
            _version,
            _safe,
            _beneficiary1,
            _beneficiary2,
            _beneficiary3,
            _mintAmount
        );
        bytes memory _upgradeData;
        if (_version == 1) {
            _upgradeData = abi.encodeCall(
                AlignedTokenV1(_newImplementation).upgradeToAndCall,
                (_newImplementation, _initData)
            );
        } else if (_version == 2) {
            _upgradeData = abi.encodeCall(
                AlignedTokenV2Example(_newImplementation).upgradeToAndCall,
                (_newImplementation, _initData)
            );
        } else {
            revert("Unsupported version");
        }
        return _upgradeData;
    }

    function deployClaimableAirdropImplementation(
        uint256 _version
    ) internal returns (address) {
        address _implementation_address;
        if (_version == 1) {
            vm.broadcast();
            ClaimableAirdropV1 _implementation = new ClaimableAirdropV1();
            _implementation_address = address(_implementation);
        } else if (_version == 2) {
            vm.broadcast();
            ClaimableAirdropV2Example _implementation = new ClaimableAirdropV2Example();
            _implementation_address = address(_implementation);
        } else {
            revert("Unsupported version");
        }
        return _implementation_address;
    }

    function claimableAirdropProxyDeploymentData(
        address _implementation,
        uint256 _version,
        address _safe,
        address _tokenContractAddress,
        address _tokenOwnerAddress,
        uint256 _limitTimestampToClaim,
        bytes32 _claimMerkleRoot
    ) internal pure returns (bytes memory) {
        return
            abi.encodePacked(
                type(ERC1967Proxy).creationCode,
                claimableAirdropInitData(
                    _implementation,
                    _version,
                    _safe,
                    _tokenContractAddress,
                    _tokenOwnerAddress,
                    _limitTimestampToClaim,
                    _claimMerkleRoot
                )
            );
    }

    function claimableAirdropInitData(
        address _implementation,
        uint256 _version,
        address _safe,
        address _tokenContractAddress,
        address _tokenOwnerAddress,
        uint256 _limitTimestampToClaim,
        bytes32 _claimMerkleRoot
    ) internal pure returns (bytes memory) {
        if (_version == 1) {
            return
                abi.encodeCall(
                    ClaimableAirdropV1(_implementation).initialize,
                    (
                        _safe,
                        _tokenContractAddress,
                        _tokenOwnerAddress,
                        _limitTimestampToClaim,
                        _claimMerkleRoot
                    )
                );
        } else if (_version == 2) {
            return
                abi.encodeCall(
                    ClaimableAirdropV2Example(_implementation).initialize,
                    (
                        _safe,
                        _tokenContractAddress,
                        _tokenOwnerAddress,
                        _limitTimestampToClaim,
                        _claimMerkleRoot
                    )
                );
        } else {
            revert("Unsupported version");
        }
    }

    function claimableAirdropUpgradeData(
        address _newImplementation,
        uint256 _version,
        address _safe,
        address _tokenContractAddress,
        address _tokenOwnerAddress,
        uint256 _limitTimestampToClaim,
        bytes32 _claimMerkleRoot
    ) internal pure returns (bytes memory) {
        bytes memory _initData = Utils.claimableAirdropInitData(
            _newImplementation,
            _version,
            _safe,
            _tokenContractAddress,
            _tokenOwnerAddress,
            _limitTimestampToClaim,
            _claimMerkleRoot
        );
        bytes memory _upgradeData;
        if (_version == 1) {
            _upgradeData = abi.encodeCall(
                ClaimableAirdropV1(_newImplementation).upgradeToAndCall,
                (_newImplementation, _initData)
            );
        } else if (_version == 2) {
            _upgradeData = abi.encodeCall(
                ClaimableAirdropV2Example(_newImplementation).upgradeToAndCall,
                (_newImplementation, _initData)
            );
        } else {
            revert("Unsupported version");
        }
        return _upgradeData;
    }

    
}

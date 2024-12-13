// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../src/AlignedToken.sol";
import "../src/ClaimableAirdrop.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "forge-std/Script.sol";
import {Utils} from "./Utils.sol";

contract DeployAll is Script {
    function run() public {
        string memory root = vm.projectRoot();
        string memory path = string.concat(root, "/script-config/config.json");
        string memory config_json = vm.readFile(path);

        address _safe = stdJson.readAddress(config_json, ".safe");
        bytes32 _salt = stdJson.readBytes32(config_json, ".salt");
        address _deployer = stdJson.readAddress(config_json, ".deployer");
        address _beneficiary1 = stdJson.readAddress(
            config_json,
            ".beneficiary1"
        );
        address _beneficiary2 = stdJson.readAddress(
            config_json,
            ".beneficiary2"
        );
        address _beneficiary3 = stdJson.readAddress(
            config_json,
            ".beneficiary3"
        );
        uint256 _mintAmount = stdJson.readUint(config_json, ".mintAmount");
        uint256 _limitTimestampToClaim = stdJson.readUint(
            config_json,
            ".limitTimestampToClaim"
        );
        bytes32 _claimMerkleRoot = stdJson.readBytes32(
            config_json,
            ".claimMerkleRoot"
        );
        uint256 _holderPrivateKey = stdJson.readUint(
            config_json,
            ".holderPrivateKey"
        );

        ProxyAdmin _proxyAdmin = deployProxyAdmin(_safe, _salt, _deployer);

        TransparentUpgradeableProxy _tokenProxy = deployAlignedTokenProxy(
            address(_proxyAdmin),
            _salt,
            _deployer,
            _beneficiary1,
            _beneficiary2,
            _beneficiary3,
            _mintAmount
        );

        TransparentUpgradeableProxy _airdropProxy = deployClaimableAirdropProxy(
            address(_proxyAdmin),
            _safe,
            _beneficiary1,
            _salt,
            _deployer,
            address(_tokenProxy),
            _limitTimestampToClaim,
            _claimMerkleRoot
        );

        approve(
            address(_tokenProxy),
            address(_airdropProxy),
            _holderPrivateKey
        );
    }

    function deployProxyAdmin(
        address _safe,
        bytes32 _salt,
        address _deployer
    ) internal returns (ProxyAdmin) {
        bytes memory _proxyAdminDeploymentData = Utils.proxyAdminDeploymentData(
            _safe
        );
        address _proxyAdminCreate2Address = Utils.deployWithCreate2(
            _proxyAdminDeploymentData,
            _salt,
            _deployer
        );

        console.log(
            "Proxy Admin Address:",
            _proxyAdminCreate2Address,
            "Owner:",
            _safe
        );
        vm.serializeAddress("proxyAdmin", "address", _proxyAdminCreate2Address);
        vm.serializeBytes(
            "proxyAdmin",
            "deploymentData",
            _proxyAdminDeploymentData
        );

        return ProxyAdmin(_proxyAdminCreate2Address);
    }

    function deployAlignedTokenProxy(
        address _proxyAdmin,
        bytes32 _salt,
        address _deployer,
        address _beneficiary1,
        address _beneficiary2,
        address _beneficiary3,
        uint256 _mintAmount
    ) internal returns (TransparentUpgradeableProxy) {
        vm.broadcast();
        AlignedToken _token = new AlignedToken();

        bytes memory _alignedTokenDeploymentData = Utils
            .alignedTokenProxyDeploymentData(
                _proxyAdmin,
                address(_token),
                _beneficiary1,
                _beneficiary2,
                _beneficiary3,
                _mintAmount
            );
        address _alignedTokenProxy = Utils.deployWithCreate2(
            _alignedTokenDeploymentData,
            _salt,
            _deployer
        );

        console.log(
            "AlignedToken proxy deployed with address:",
            _alignedTokenProxy,
            "and admin:",
            _proxyAdmin
        );
        vm.serializeAddress("alignedToken", "address", _alignedTokenProxy);
        vm.serializeBytes(
            "alignedToken",
            "deploymentData",
            _alignedTokenDeploymentData
        );

        return TransparentUpgradeableProxy(payable(_alignedTokenProxy));
    }

    function deployClaimableAirdropProxy(
        address _proxyAdmin,
        address _owner,
        address _tokenOwner,
        bytes32 _salt,
        address _deployer,
        address _token,
        uint256 _limitTimestampToClaim,
        bytes32 _claimMerkleRoot
    ) internal returns (TransparentUpgradeableProxy) {
        vm.broadcast();
        ClaimableAirdrop _airdrop = new ClaimableAirdrop();

        bytes memory _airdropDeploymentData = Utils
            .claimableAirdropProxyDeploymentData(
                _proxyAdmin,
                address(_airdrop),
                _owner,
                _token,
                _tokenOwner,
                _limitTimestampToClaim,
                _claimMerkleRoot
            );
        address _airdropProxy = Utils.deployWithCreate2(
            _airdropDeploymentData,
            _salt,
            _deployer
        );

        console.log(
            "ClaimableAirdrop proxy deployed with address:",
            _airdropProxy,
            "and admin:",
            _proxyAdmin
        );
        vm.serializeAddress("claimableAirdrop", "address", _airdropProxy);
        vm.serializeBytes(
            "claimableAirdrop",
            "deploymentData",
            _airdropDeploymentData
        );

        return TransparentUpgradeableProxy(payable(_airdropProxy));
    }

    function approve(
        address _tokenContractProxy,
        address _airdropContractProxy,
        uint256 _holderPrivateKey
    ) public {
        vm.startBroadcast(_holderPrivateKey);
        (bool success, bytes memory data) = address(_tokenContractProxy).call(
            abi.encodeCall(
                IERC20.approve,
                (address(_airdropContractProxy), 1 ether)
            )
        );
        bool approved;
        assembly {
            approved := mload(add(data, 0x20))
        }

        if (!success || !approved) {
            revert("Failed to give approval to airdrop contract");
        }
        vm.stopBroadcast();

        console.log("Succesfully gave approval to airdrop contract");
    }
}

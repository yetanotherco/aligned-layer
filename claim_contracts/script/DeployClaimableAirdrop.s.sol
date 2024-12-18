// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../src/ClaimableAirdrop.sol";

import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import "forge-std/Script.sol";

import {Utils} from "./Utils.sol";

contract DeployClaimableAirdrop is Script {
    function run(string memory config) public {
        string memory root = vm.projectRoot();
        string memory path = string.concat(
            root,
            "/script-config/config.",
            config,
            ".json"
        );
        string memory config_json = vm.readFile(path);
        console.log(config);

        bytes32 _salt = stdJson.readBytes32(config_json, ".salt");
        address _deployer = stdJson.readAddress(config_json, ".deployer");
        address _foundation = stdJson.readAddress(config_json, ".foundation");
        address _safe = _foundation;
        address _tokenProxy = stdJson.readAddress(config_json, ".tokenProxy");
        uint256 _claimPrivateKey = vm.envUint("CLAIM_SUPPLIER_PRIVATE_KEY");
        uint256 _limitTimestampToClaim = stdJson.readUint(
            config_json,
            ".limitTimestampToClaim"
        );
        bytes32 _claimMerkleRoot = stdJson.readBytes32(
            config_json,
            ".claimMerkleRoot"
        );
        console.logBytes32(_claimMerkleRoot);

        TransparentUpgradeableProxy _airdropProxy = deployClaimableAirdropProxy(
            address(_safe),
            _safe,
            _foundation,
            _salt,
            _deployer,
            address(_tokenProxy),
            _limitTimestampToClaim,
            _claimMerkleRoot
        );

        approve(_tokenProxy, address(_airdropProxy), _claimPrivateKey);
    }

    function deployClaimableAirdropProxy(
        address _proxyOwner,
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
                _proxyOwner,
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

        address _proxyAdmin = Utils.getAdminAddress(_airdropProxy);

        console.log(
            "ClaimableAirdrop proxy deployed with address:",
            _airdropProxy,
            "and admin:",
            _proxyAdmin
        );
        vm.serializeAddress("claimableAirdrop", "address", _airdropProxy);
        vm.serializeAddress("claimableAirdropAdmin", "address", _proxyAdmin);
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
        uint256 _claimPrivateKey
    ) public {
        vm.startBroadcast(_claimPrivateKey);
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

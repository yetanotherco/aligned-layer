// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../src/AlignedToken.sol";
import "../src/ClaimableAirdrop.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "forge-std/Script.sol";
import {Utils} from "./Utils.sol";

contract DeployAlignedToken is Script {
    function run() public {
        string memory root = vm.projectRoot();
        string memory path = string.concat(root, "/script-config/config.json");
        string memory config_json = vm.readFile(path);

        address _safe = stdJson.readAddress(config_json, ".safe");
        bytes32 _salt = stdJson.readBytes32(config_json, ".salt");
        address _deployer = stdJson.readAddress(config_json, ".deployer");
        address _foundation = stdJson.readAddress(config_json, ".foundation");
        address _claimSupplier = stdJson.readAddress(
            config_json,
            ".claimSupplier"
        );

        ProxyAdmin _proxyAdmin = deployProxyAdmin(_safe, _salt, _deployer);

        console.log(
            "Proxy Admin deployed at address:",
            address(_proxyAdmin),
            "with owner:",
            _safe
        );

        TransparentUpgradeableProxy _tokenProxy = deployAlignedTokenProxy(
            address(_proxyAdmin),
            _salt,
            _deployer,
            _safe,
            _foundation,
            _claimSupplier
        );

        console.log(
            string.concat(
                "Aligned Token Proxy deployed at address:",
                vm.toString(address(_tokenProxy)),
                "with proxy admin:",
                vm.toString(address(_proxyAdmin)),
                "and owner:",
                vm.toString(_safe)
            )
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

        return ProxyAdmin(_proxyAdminCreate2Address);
    }

    function deployAlignedTokenProxy(
        address _proxyAdmin,
        bytes32 _salt,
        address _deployer,
        address _owner,
        address _foundation,
        address _claim
    ) internal returns (TransparentUpgradeableProxy) {
        vm.broadcast();
        AlignedToken _token = new AlignedToken();

        bytes memory _alignedTokenDeploymentData = Utils
            .alignedTokenProxyDeploymentData(
                _proxyAdmin,
                address(_token),
                _owner,
                _foundation,
                _claim
            );
        address _alignedTokenProxy = Utils.deployWithCreate2(
            _alignedTokenDeploymentData,
            _salt,
            _deployer
        );
        return TransparentUpgradeableProxy(payable(_alignedTokenProxy));
    }
}

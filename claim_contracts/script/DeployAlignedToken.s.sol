// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../src/AlignedToken.sol";
import "../src/ClaimableAirdrop.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "forge-std/Script.sol";
import {Vm} from "forge-std/Vm.sol";
import {Utils} from "./Utils.sol";

contract DeployAlignedToken is Script {
    function run(string memory config) public {
        string memory root = vm.projectRoot();
        string memory path = string.concat(root, "/script-config/config.", config, ".json");
        string memory config_json = vm.readFile(path);

        address _safe = stdJson.readAddress(config_json, ".safe");
        bytes32 _salt = stdJson.readBytes32(config_json, ".salt");
        address _deployer = stdJson.readAddress(config_json, ".deployer");
        address _foundation = stdJson.readAddress(config_json, ".foundation");
        address _claimSupplier = stdJson.readAddress(config_json, ".claimSupplier");

        TransparentUpgradeableProxy _tokenProxy =
            deployAlignedTokenProxy(_safe, _salt, _deployer, _foundation, _claimSupplier);

        console.log(
            string.concat(
                "Aligned Token Proxy deployed at address: ",
                vm.toString(address(_tokenProxy)),
                " with proxy admin: ",
                vm.toString(Utils.getAdminAddress(address(_tokenProxy))),
                " and owner: ",
                vm.toString(_safe)
            )
        );

        string memory deployedAddressesJson = "deployedAddressesJson";
        string memory finalJson = vm.serializeAddress(deployedAddressesJson, "tokenProxy", address(_tokenProxy));

        vm.writeJson(finalJson, _getOutputPath("deployed_token_addresses.json"));
    }

    function _getOutputPath(string memory fileName) internal returns (string memory) {
        string memory outputDir = "script-out/";

        // Create output directory if it doesn't exist
        if (!vm.exists(outputDir)) {
            vm.createDir(outputDir, true);
        }

        return string.concat(outputDir, fileName);
    }

    function deployAlignedTokenProxy(
        address _proxyAdmin,
        bytes32 _salt,
        address _deployer,
        address _foundation,
        address _claim
    ) internal returns (TransparentUpgradeableProxy) {
        vm.broadcast();
        AlignedToken _token = new AlignedToken();

        bytes memory _alignedTokenDeploymentData =
            Utils.alignedTokenProxyDeploymentData(_proxyAdmin, address(_token), _foundation, _claim);
        address _alignedTokenProxy = Utils.deployWithCreate2(_alignedTokenDeploymentData, _salt, _deployer);
        return TransparentUpgradeableProxy(payable(_alignedTokenProxy));
    }
}

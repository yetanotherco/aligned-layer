// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract ProxyAdminDeploymentData is Script {
    function run(address _safe) public {
        bytes memory _proxyDeploymentData = Utils.proxyAdminDeploymentData(
            _safe
        );

        console.logBytes(_proxyDeploymentData);

        string memory _out = vm.serializeBytes(
            "implementation",
            "deploymentData",
            _proxyDeploymentData
        );

        string memory _path = string.concat(
            vm.projectRoot(),
            "/script-out/claimable_airdrop_implementation.json"
        );

        vm.writeJson(_out, _path);
    }
}

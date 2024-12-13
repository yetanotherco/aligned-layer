// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract DeployProxyAdmin is Script {
    function run(address _safe) public {
        address _proxyAdmin = Utils.deployProxyAdmin(_safe);

        console.log("Claimable Airdrop Implementation Address:", _proxyAdmin);

        vm.serializeAddress("implementation", "address", _proxyAdmin);

        string memory _out = vm.serializeBytes(
            "implementation",
            "deploymentData",
            Utils.proxyAdminDeploymentData(_safe)
        );

        string memory _path = string.concat(
            vm.projectRoot(),
            "/script-out/claimable_airdrop_implementation.json"
        );

        vm.writeJson(_out, _path);
    }
}

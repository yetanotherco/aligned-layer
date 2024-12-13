// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract DeployAlignedTokenProxy is Script {
    function run(
        address _proxyAdmin,
        address _implementation,
        address _beneficiary1,
        address _beneficiary2,
        address _beneficiary3,
        uint256 _mintAmount
    ) public {
        bytes memory _deploymentData = Utils.alignedTokenInitData(
            _implementation,
            _beneficiary1,
            _beneficiary2,
            _beneficiary3,
            _mintAmount
        );

        vm.broadcast();
        TransparentUpgradeableProxy _proxy = new TransparentUpgradeableProxy(
            _implementation,
            _proxyAdmin,
            _deploymentData
        );

        console.log("Aligned Token Proxy Address:", address(_proxy));

        vm.serializeAddress("proxy", "address", address(_proxy));

        string memory _out = vm.serializeBytes(
            "proxy",
            "deploymentData",
            Utils.alignedTokenProxyDeploymentData(
                _proxyAdmin,
                _implementation,
                _beneficiary1,
                _beneficiary2,
                _beneficiary3,
                _mintAmount
            )
        );

        string memory _path = string.concat(
            vm.projectRoot(),
            "/script-out/aligned_token_proxy.json"
        );

        vm.writeJson(_out, _path);
    }
}

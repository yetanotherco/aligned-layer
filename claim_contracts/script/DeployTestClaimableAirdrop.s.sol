// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../test/TestClaimableAirdrop.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "forge-std/Script.sol";
import {Vm} from "forge-std/Vm.sol";
import {Utils} from "./Utils.sol";

contract DeployTestClaimableAirdrop is Script {
    function run(string memory config) public {
        string memory root = vm.projectRoot();
        string memory path = string.concat(
            root,
            "/script-config/config.",
            config,
            ".json"
        );
        string memory config_json = vm.readFile(path);

        address _foundation = stdJson.readAddress(config_json, ".foundation");
        address _tokenDistributor = stdJson.readAddress(
            config_json,
            ".tokenDistributor"
        );
        address _tokenProxy = stdJson.readAddress(config_json, ".tokenProxy");

        vm.broadcast();
        TestClaimableAirdrop _airdrop = new TestClaimableAirdrop();

        console.log(
            "Test Claimable Airdrop deployed at address:",
            address(_airdrop)
        );

        vm.broadcast();
        TransparentUpgradeableProxy _airdropProxy = new TransparentUpgradeableProxy(
                address(_airdrop),
                _foundation,
                Utils.claimableAirdropInitData(
                    address(_airdrop),
                    _foundation,
                    _tokenProxy,
                    _tokenDistributor
                )
            );

        bytes memory _alignedTokenProxyConstructorData = Utils
            .claimableAirdropProxyConstructorData(
                address(_airdrop),
                _foundation,
                _tokenProxy,
                _tokenDistributor
            );

        console.log(
            string.concat(
                "Test Claimable Airdrop Proxy deployed at address: ",
                vm.toString(address(_airdropProxy)),
                " with proxy admin: ",
                vm.toString(Utils.getAdminAddress(address(_airdropProxy))),
                " and owner: ",
                vm.toString(_foundation),
                " with constructor args: ",
                vm.toString(_alignedTokenProxyConstructorData)
            )
        );
    }
}

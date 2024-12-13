// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../../src/AlignedTokenV1.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract UpgradeAlignedTokenImplementation is Script {
    function run(
        address _proxyAdmin,
        address _proxy,
        address _newImplementation,
        address _beneficiary1,
        address _beneficiary2,
        address _beneficiary3,
        uint256 _mintAmount
    ) public {
        vm.broadcast();
        (bool success, ) = _proxyAdmin.call(
            Utils.alignedTokenUpgradeData(
                _proxyAdmin,
                _proxy,
                _newImplementation,
                _beneficiary1,
                _beneficiary2,
                _beneficiary3,
                _mintAmount
            )
        );

        if (!success) {
            revert("Failed to give approval to airdrop contract");
        }
        vm.stopBroadcast();

        console.log("Succesfully gave approval to airdrop contract");
    }
}

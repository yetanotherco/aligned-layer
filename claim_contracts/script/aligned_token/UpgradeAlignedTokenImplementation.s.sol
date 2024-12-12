// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../src/AlignedTokenV1.sol";
import "../src/AlignedTokenV2Example.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import "forge-std/Script.sol";
import {Utils} from "./Utils.sol";

contract UpgradeAlignedTokenImplementation is Script {
    function run(
        address _proxy,
        address _newImplementation,
        uint256 _version,
        address _safe,
        address _beneficiary1,
        address _beneficiary2,
        address _beneficiary3,
        uint256 _mintAmount
    ) public {
        bytes memory _upgradeData = Utils.alignedTokenUpgradeData(
            _newImplementation,
            _version,
            _safe,
            _beneficiary1,
            _beneficiary2,
            _beneficiary3,
            _mintAmount
        );

        vm.broadcast();
        _proxy.call(_upgradeData);
    }
}

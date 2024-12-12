// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract AlignedTokenUpgradeData is Script {
    function run(
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
        console.logBytes(_upgradeData);
        vm.writeFile(
            string.concat(
                vm.projectRoot(),
                "/script-out/aligned_token_upgrade_data.hex"
            ),
            vm.toString(_upgradeData)
        );
    }
}

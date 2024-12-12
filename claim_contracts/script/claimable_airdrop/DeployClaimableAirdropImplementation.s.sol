// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../../src/ClaimableAirdropV1.sol";
import "../../src/ClaimableAirdropV2Example.sol";
import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract DeployTokenImplementation is Script {
    function run(uint256 _version) public {
        address _implementation_address = Utils
            .deployClaimableAirdropImplementation(_version);

        console.log(
            string.concat(
                "Claimable Airdrop Implementation v",
                vm.toString(_version),
                "Address:"
            ),
            _implementation_address
        );

        vm.serializeAddress(
            "implementation",
            "address",
            _implementation_address
        );

        string memory out;
        if (_version == 1) {
            out = vm.serializeBytes(
                "implementation",
                "deploymentData",
                type(ClaimableAirdropV1).creationCode
            );
        } else if (_version == 2) {
            out = vm.serializeBytes(
                "implementation",
                "deploymentData",
                type(ClaimableAirdropV2Example).creationCode
            );
        } else {
            revert("Unsupported version");
        }

        string memory path = string.concat(
            vm.projectRoot(),
            "/script-out/claimable_airdrop_implementation_v",
            vm.toString(_version),
            ".json"
        );

        vm.writeJson(out, path);
    }
}

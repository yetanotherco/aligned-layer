// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../../src/AlignedTokenV1.sol";
import "../../src/AlignedTokenV2Example.sol";
import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract DeployTokenImplementation is Script {
    function run(uint256 _version) public {
        address _implementation_address = Utils
            .deployAlignedTokenImplementation(_version);

        console.log(
            string.concat(
                "Token Implementation v",
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
                type(AlignedTokenV1).creationCode
            );
        } else if (_version == 2) {
            out = vm.serializeBytes(
                "implementation",
                "deploymentData",
                type(AlignedTokenV2Example).creationCode
            );
        } else {
            revert("Unsupported version");
        }

        string memory path = string.concat(
            vm.projectRoot(),
            "/script-out/aligned_token_implementation_v",
            vm.toString(_version),
            ".json"
        );

        vm.writeJson(out, path);
    }
}

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../../src/ClaimableAirdrop.sol";
import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract DeployTokenImplementation is Script {
    function run() public {
        address _implementation_address = Utils
            .deployClaimableAirdropImplementation();

        console.log(
            "Claimable Airdrop Implementation Address:",
            _implementation_address
        );

        vm.serializeAddress(
            "implementation",
            "address",
            _implementation_address
        );

        string memory _out = vm.serializeBytes(
            "implementation",
            "deploymentData",
            type(ClaimableAirdrop).creationCode
        );

        string memory _path = string.concat(
            vm.projectRoot(),
            "/script-out/claimable_airdrop_implementation.json"
        );

        vm.writeJson(_out, _path);
    }
}

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../../src/AlignedTokenV1.sol";
import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract AlignedTokenCreate2 is Script {
    function run(bytes32 _salt, address _deployer) public {
        address _create2Address = vm.computeCreate2Address(
            _salt,
            keccak256(type(AlignedTokenV1).creationCode),
            _deployer
        );
        console.logAddress(_create2Address);
        vm.writeFile(
            string.concat(
                vm.projectRoot(),
                "/script-out/aligned_token_create2_address.hex"
            ),
            vm.toString(_create2Address)
        );
    }
}

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../../src/ClaimableAirdrop.sol";
import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract ClaimableAirdropCreate2 is Script {
    function run(bytes32 _salt, address _deployer) public {
        address _create2Address = vm.computeCreate2Address(
            _salt,
            keccak256(type(ClaimableAirdrop).creationCode),
            _deployer
        );

        console.logAddress(_create2Address);

        vm.writeFile(
            string.concat(
                vm.projectRoot(),
                "/script-out/claimable_airdrop_create2_address.hex"
            ),
            vm.toString(_create2Address)
        );
    }
}

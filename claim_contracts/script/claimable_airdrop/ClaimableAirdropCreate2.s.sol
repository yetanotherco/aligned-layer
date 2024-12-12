// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "../../src/ClaimableAirdropV1.sol";
import "../../src/ClaimableAirdropV2Example.sol";
import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract ClaimableAirdropCreate2 is Script {
    function run(uint256 _version, bytes32 _salt, address _deployer) public {
        address _create2Address;
        if (_version == 1) {
            _create2Address = vm.computeCreate2Address(
                _salt,
                keccak256(type(ClaimableAirdropV1).creationCode),
                _deployer
            );
        } else if (_version == 2) {
            _create2Address = vm.computeCreate2Address(
                _salt,
                keccak256(type(ClaimableAirdropV2Example).creationCode),
                _deployer
            );
        } else {
            revert("Unsupported version");
        }

        console.logAddress(_create2Address);

        vm.writeFile(
            string.concat(
                vm.projectRoot(),
                "/script-out/claimable_airdrop_v",
                vm.toString(_version),
                "_create2_address.hex"
            ),
            vm.toString(_create2Address)
        );
    }
}

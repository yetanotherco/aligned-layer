// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract AlignedTokenInitData is Script {
    function run(
        address _implementation,
        uint256 _version,
        address _safe,
        address _tokenContractAddress,
        address _tokenOwnerAddress,
        uint256 _limitTimestampToClaim,
        bytes32 _claimMerkleRoot
    ) public {
        bytes memory data = Utils.claimableAirdropInitData(
            _implementation,
            _version,
            _safe,
            _tokenContractAddress,
            _tokenOwnerAddress,
            _limitTimestampToClaim,
            _claimMerkleRoot
        );

        console.logBytes(data);

        vm.writeFile(
            string.concat(
                vm.projectRoot(),
                "/script-out/claimable_airdrop_init_data.hex"
            ),
            vm.toString(data)
        );
    }
}

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract AlignedTokenProxyDeploymentData is Script {
    function run(
        address _proxyAdmin,
        address _implementation,
        address _tokenContractAddress,
        address _tokenOwnerAddress,
        uint256 _limitTimestampToClaim,
        bytes32 _claimMerkleRoot
    ) public {
        bytes memory data = Utils.claimableAirdropProxyDeploymentData(
            _proxyAdmin,
            _implementation,
            _tokenContractAddress,
            _tokenOwnerAddress,
            _limitTimestampToClaim,
            _claimMerkleRoot
        );
        console.logBytes(data);
        vm.writeFile(
            string.concat(
                vm.projectRoot(),
                "/script-out/claimable_airdrop_proxy_deployment_data.hex"
            ),
            vm.toString(data)
        );
    }
}

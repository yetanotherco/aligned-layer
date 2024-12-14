// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract ClaimableAirdropnProxyCreate2 is Script {
    function run(
        address _multisig,
        address _proxyAdmin,
        bytes32 _salt,
        address _deployer,
        address _implementation,
        address _tokenContractAddress,
        address _tokenOwnerAddress,
        uint256 _limitTimestampToClaim,
        bytes32 _claimMerkleRoot
    ) public {
        address _create2Address = vm.computeCreate2Address(
            _salt,
            keccak256(
                Utils.claimableAirdropProxyDeploymentData(
                    _proxyAdmin,
                    _implementation,
                    _multisig,
                    _tokenContractAddress,
                    _tokenOwnerAddress,
                    _limitTimestampToClaim,
                    _claimMerkleRoot
                )
            ),
            _deployer
        );

        console.logAddress(_create2Address);

        vm.writeFile(
            string.concat(
                vm.projectRoot(),
                "/script-out/claimable_airdrop_proxy_create2_address.hex"
            ),
            vm.toString(_create2Address)
        );
    }
}

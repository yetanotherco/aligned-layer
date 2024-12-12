// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract AlignedTokenProxyCreate2 is Script {
    function run(
        bytes32 _salt,
        address _deployer,
        address _implementation,
        uint256 _version,
        address _safe,
        address _beneficiary1,
        address _beneficiary2,
        address _beneficiary3,
        uint256 _mintAmount
    ) public {
        address _create2Address = vm.computeCreate2Address(
            _salt,
            keccak256(
                Utils.alignedTokenProxyDeploymentData(
                    _implementation,
                    _version,
                    _safe,
                    _beneficiary1,
                    _beneficiary2,
                    _beneficiary3,
                    _mintAmount
                )
            ),
            _deployer
        );
        console.logAddress(_create2Address);
        vm.writeFile(
            string.concat(
                vm.projectRoot(),
                "/script-out/aligned_token_proxy_create2_address.hex"
            ),
            vm.toString(_create2Address)
        );
    }
}

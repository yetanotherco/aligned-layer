// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import "forge-std/Script.sol";
import {Utils} from "../Utils.sol";

contract DeployAlignedTokenProxy is Script {
    function run(
        address _implementation,
        uint256 _version,
        address _safe,
        address _tokenContractAddress,
        address _tokenOwnerAddress,
        uint256 _limitTimestampToClaim,
        bytes32 _claimMerkleRoot
    ) public {
        bytes memory _proxyDeploymentData = Utils.claimableAirdropInitData(
            _implementation,
            _version,
            _safe,
            _tokenContractAddress,
            _tokenOwnerAddress,
            _limitTimestampToClaim,
            _claimMerkleRoot
        );
        vm.broadcast();
        ERC1967Proxy _proxy = new ERC1967Proxy(
            _implementation,
            _proxyDeploymentData
        );
        console.log("Claimable Airdrop Proxy Address:", address(_proxy));

        vm.serializeAddress("proxy", "address", address(_proxy));
        string memory _out = vm.serializeBytes(
            "proxy",
            "deploymentData",
            Utils.claimableAirdropProxyDeploymentData(
                _implementation,
                _version,
                _safe,
                _tokenContractAddress,
                _tokenOwnerAddress,
                _limitTimestampToClaim,
                _claimMerkleRoot
            )
        );
        string memory _path = string.concat(
            vm.projectRoot(),
            "/script-out/claimable_airdrop_proxy.json"
        );
        vm.writeJson(_out, _path);
    }
}

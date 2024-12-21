// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// import "../src/TestToken.sol";
import "../src/ExampleAlignedTokenV2.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {ERC1967Utils} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Utils.sol";
import "forge-std/Script.sol";
import {Vm} from "forge-std/Vm.sol";
import {Utils} from "./Utils.sol";

/// @notice This script upgrades the ClaimableAirdrop contract to ClaimableAirdropV2.
/// @dev The `ProxyAdmin` owner must be the runner of this script since it is the
/// one that will call the upgradeAndCall function of the `ProxyAdmin`.
contract UpgradeToAlignedTokenV2 is Script {
    function run(string memory config) public {
        string memory root = vm.projectRoot();
        string memory path = string.concat(
            root,
            "/script-config/config.",
            config,
            ".json"
        );
        string memory config_json = vm.readFile(path);

        address _currentTokenProxy = stdJson.readAddress(
            config_json,
            ".tokenProxy"
        );

        vm.startBroadcast();
        ExampleAlignedTokenV2 _newToken = new ExampleAlignedTokenV2();

        address _adminAddress = getAdminAddress(_currentTokenProxy);

        ProxyAdmin(_adminAddress).upgradeAndCall(
            ITransparentUpgradeableProxy(_currentTokenProxy),
            address(_newToken),
            abi.encodeCall(ExampleAlignedTokenV2.reinitialize, ())
        );

        vm.stopBroadcast();
    }

    function getAdminAddress(address proxy) internal view returns (address) {
        address CHEATCODE_ADDRESS = 0x7109709ECfa91a80626fF3989D68f67F5b1DD12D;
        Vm vm = Vm(CHEATCODE_ADDRESS);
        
        bytes32 adminSlot = vm.load(proxy, ERC1967Utils.ADMIN_SLOT);
        return address(uint160(uint256(adminSlot)));
    }
}

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Script.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract ApproveERC20 is Script {
    function run(
        address tokenContractProxy,
        address airdropContractProxy,
        uint256 amounToClaim
    ) public {
        vm.startBroadcast();
        (bool success, bytes memory data) = address(tokenContractProxy).call(
            abi.encodeCall(
                IERC20.approve,
                (address(airdropContractProxy), amounToClaim)
            )
        );
        bool approved;
        assembly {
            approved := mload(add(data, 0x20))
        }

        if (!success || !approved) {
            revert("Failed to give approval to airdrop contract");
        }
        vm.stopBroadcast();

        console.log("Succesfully gave approval to airdrop contract");
    }
}

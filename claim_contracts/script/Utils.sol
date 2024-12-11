// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Vm} from "forge-std/Vm.sol";

library Utils {
    // Cheatcodes address, 0x7109709ECfa91a80626fF3989D68f67F5b1DD12D.
    address internal constant VM_ADDRESS =
        address(uint160(uint256(keccak256("hevm cheat code"))));
    Vm internal constant vm = Vm(VM_ADDRESS);

    /// @notice Address of the deterministic create2 factory.
    /// @dev This address corresponds to a contracts that is set in the storage
    /// in the genesis file. The same contract with the same address is deployed
    /// in every testnet, so if this script is run in a testnet instead of in a
    /// local environment, it should work.
    address constant DETERMINISTIC_CREATE2_ADDRESS =
        0x4e59b44847b379578588920cA78FbF26c0B4956C;

    function deployWithCreate2(
        bytes memory bytecode,
        bytes32 salt,
        address create2Factory,
        uint256 signerPrivateKey
    ) internal returns (address) {
        if (bytecode.length == 0) {
            revert("Bytecode is not set");
        }
        address contractAddress = vm.computeCreate2Address(
            salt,
            keccak256(bytecode),
            create2Factory
        );
        if (contractAddress.code.length != 0) {
            revert("Contract already deployed");
        }

        vm.broadcast(signerPrivateKey);
        (bool success, bytes memory data) = create2Factory.call(
            abi.encodePacked(salt, bytecode)
        );
        contractAddress = bytesToAddress(data);

        if (!success) {
            revert(
                "Failed to deploy contract via create2: create2Factory call failed"
            );
        }

        if (contractAddress == address(0)) {
            revert(
                "Failed to deploy contract via create2: contract address is zero"
            );
        }

        if (contractAddress.code.length == 0) {
            revert(
                "Failed to deploy contract via create2: contract code is empty"
            );
        }

        return contractAddress;
    }

    function bytesToAddress(
        bytes memory addressOffset
    ) internal pure returns (address addr) {
        assembly {
            addr := mload(add(addressOffset, 20))
        }
    }
}

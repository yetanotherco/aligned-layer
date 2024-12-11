// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console} from "forge-std/Script.sol";
import {AlignedToken} from "../src/AlignedToken.sol";
import {ClaimableAirdrop} from "../src/ClaimableAirdrop.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {Utils} from "./Utils.sol";

contract DeployScript is Script {
    struct Supply {
        address beneficiary;
        uint256 amount;
    }

    function setUp() public {}

    function run(
        address safe,
        address owner1,
        address owner2,
        address owner3,
        uint256 mintAmount,
        uint256 limitTimestamp,
        bytes32 merkleRoot
    ) public {
        console.log("Deploying contracts");
        uint256 deployer = vm.envUint("DEPLOYER_PRIVATE_KEY");

        ProxyAdmin contractsProxyAdmin = deployProxyAdmin(safe, deployer);

        TransparentUpgradeableProxy tokenContractProxy = deployTokenContractProxy(
                address(contractsProxyAdmin),
                owner1,
                owner2,
                owner3,
                mintAmount,
                deployer
            );

        TransparentUpgradeableProxy airdropContractProxy = deployAirdropContractProxy(
                address(contractsProxyAdmin),
                address(tokenContractProxy),
                owner3,
                limitTimestamp,
                merkleRoot,
                deployer
            );

        vm.startBroadcast();
        (deployer);
        (bool success, bytes memory data) = address(tokenContractProxy).call(
            abi.encodeCall(
                IERC20.approve,
                (address(airdropContractProxy), mintAmount)
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

    function deployProxyAdmin(
        address tokenContractProxyAdminOwner,
        uint256 signerPrivateKey
    ) internal returns (ProxyAdmin) {
        bytes memory bytecode = abi.encodePacked(
            type(ProxyAdmin).creationCode,
            abi.encode(tokenContractProxyAdminOwner)
        );
        bytes32 salt = bytes32(0);
        address proxyAdminAddress = Utils.deployWithCreate2(
            bytecode,
            salt,
            Utils.DETERMINISTIC_CREATE2_ADDRESS,
            signerPrivateKey
        );
        console.log(
            "Aligned Proxy Admin deployed at:",
            proxyAdminAddress,
            "with owner:",
            tokenContractProxyAdminOwner
        );
        return ProxyAdmin(proxyAdminAddress);
    }

    function deployTokenContractProxy(
        address tokenContractProxyAdmin,
        address owner1,
        address owner2,
        address owner3,
        uint256 mintAmount,
        uint256 signerPrivateKey
    ) internal returns (TransparentUpgradeableProxy) {
        vm.broadcast(signerPrivateKey);
        AlignedToken tokenContract = new AlignedToken();
        bytes memory bytecode = abi.encodePacked(
            type(TransparentUpgradeableProxy).creationCode,
            abi.encode(
                address(tokenContract),
                tokenContractProxyAdmin,
                abi.encodeCall(
                    tokenContract.initialize,
                    (owner1, mintAmount, owner2, mintAmount, owner3, mintAmount)
                )
            )
        );
        bytes32 salt = bytes32(0);
        address tokenContractProxyAddress = Utils.deployWithCreate2(
            bytecode,
            salt,
            Utils.DETERMINISTIC_CREATE2_ADDRESS,
            signerPrivateKey
        );
        console.log(
            "Aligned Token Proxy deployed at:",
            tokenContractProxyAddress,
            "with admin:",
            tokenContractProxyAdmin
        );
        return TransparentUpgradeableProxy(payable(tokenContractProxyAddress));
    }

    function deployAirdropContractProxy(
        address airdropContractProxyAdmin,
        address tokenContractProxyAddress,
        address tokenOwnerAddress,
        uint256 limitTimestamp,
        bytes32 merkleRoot,
        uint256 signerPrivateKey
    ) internal returns (TransparentUpgradeableProxy) {
        vm.broadcast(signerPrivateKey);
        ClaimableAirdrop airdropContract = new ClaimableAirdrop();
        bytes memory bytecode = abi.encodePacked(
            type(TransparentUpgradeableProxy).creationCode,
            abi.encode(
                address(airdropContract),
                airdropContractProxyAdmin,
                abi.encodeCall(
                    airdropContract.initialize,
                    (
                        tokenContractProxyAddress,
                        tokenOwnerAddress,
                        limitTimestamp,
                        merkleRoot
                    )
                )
            )
        );
        bytes32 salt = bytes32(0);
        address airdropContractProxy = Utils.deployWithCreate2(
            bytecode,
            salt,
            Utils.DETERMINISTIC_CREATE2_ADDRESS,
            signerPrivateKey
        );
        console.log(
            "Airdrop Proxy deployed at:",
            airdropContractProxy,
            "with admin:",
            airdropContractProxyAdmin
        );
        return TransparentUpgradeableProxy(payable(airdropContractProxy));
    }
}

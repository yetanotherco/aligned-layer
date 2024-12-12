.PHONY: help deploy-aligned-token-implementation deploy-aligned-token-proxy deploy-claimable-airdrop-implementation deploy-claimable-airdrop-proxy upgrade-aligned-token-implementation aligned-token-proxy-deploy-data aligned-token-init-data aligned-token-upgrade-data aligned-token-create2 aligned-token-proxy-create2

RPC_URL?=http://localhost:8545
PRIVATE_KEY?=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
VERSION?=1

help: ## 📚 Show help for each of the Makefile recipes
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Deployments

MINT?="1 ether"
OWNER?=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
BENEFICIARY1?=0x70997970C51812dc3A010C7d01b50e0d17dc79C8
BENEFICIARY2?=0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
BENEFICIARY3?=0x90F79bf6EB2c4f870365E785982E1f101E93b906

deploy-aligned-token-implementation: ## 🚀 Deploy the AlignedToken implementation contract
	cd script/aligned_token && \
	forge script DeployAlignedTokenImplementation.s.sol \
	--sig "run(uint256)" \
	$(VERSION) \
	--rpc-url $(RPC_URL) \
	--private-key $(PRIVATE_KEY) \
	--broadcast

deploy-aligned-token-proxy: ## 🚀 Deploy the AlignedToken proxy contract
	cd script/aligned_token && \
	forge script DeployAlignedTokenProxy.s.sol \
	--sig "run(address,uint256,address,address,address,address,uint256)" \
	$(IMPLEMENTATION) $(VERSION) $(OWNER) $(BENEFICIARY1) $(BENEFICIARY2) $(BENEFICIARY3) $(MINT) \
	--rpc-url $(RPC_URL) \
	--private-key $(PRIVATE_KEY) \
	--broadcast

deploy-claimable-airdrop-implementation: ## 🚀 Deploy the ClaimableAirdrop implementation contract
	cd script/claimable_airdrop && \
	forge script DeployClaimableAirdropImplementation.s.sol \
	--sig "run(uint256)" \
	$(VERSION) \
	--rpc-url $(RPC_URL) \
	--private-key $(PRIVATE_KEY) \
	--broadcast

deploy-claimable-airdrop-proxy: ## 🚀 Deploy the ClaimableAirdrop proxy contract
	cd script/claimable_airdrop && \
	forge script DeployClaimableAirdropProxy.s.sol \
	--sig "run(address,uint256,address,address,address,address)" \
	$(IMPLEMENTATION) $(VERSION) $(OWNER) $(BENEFICIARY1) $(BENEFICIARY2) $(BENEFICIARY3) $(MINT) \
	--rpc-url $(RPC_URL) \
	--private-key $(PRIVATE_KEY) \
	--broadcast

# Upgrades

upgrade-aligned-token-implementation: ## 🚀 Upgrade the AlignedToken implementation contract
	cd script/aligned_token && \
	forge script UpgradeAlignedTokenImplementation.s.sol \
	--sig "function run(address,address,uint256,address,address,address,address,uint256)" \
	$(PROXY) $(IMPLEMENTATION) $(VERSION) $(OWNER) $(BENEFICIARY1) $(BENEFICIARY2) $(BENEFICIARY3) $(MINT)\
	--rpc-url $(RPC_URL) \
	--private-key $(PRIVATE_KEY) \
	--broadcast

# Deployment Data

aligned-token-proxy-deploy-data: ## 🚀 Generate the deployment data for the AlignedToken proxy contract
	cd script/aligned_token && \
	forge script AlignedTokenProxyDeploymentData.s.sol \
	--sig "run(address,uint256, address,address,address,address,uint256)" \
	$(IMPLEMENTATION) $(VERSION) $(OWNER) $(BENEFICIARY1) $(BENEFICIARY2) $(BENEFICIARY3) $(MINT)

aligned-token-init-data: ## 🚀 Generate the init data for the AlignedToken proxy contract
	cd script/aligned_token && \
	forge script AlignedTokenInitData.s.sol \
	--sig "run(address,uint256, address,address,address,address,uint256)" \
	$(IMPLEMENTATION) $(VERSION) $(OWNER) $(BENEFICIARY1) $(BENEFICIARY2) $(BENEFICIARY3) $(MINT)

aligned-token-upgrade-data: ## 🚀 Generate the upgrade data for the AlignedToken proxy contract
	cd script/aligned_token && \
	forge script AlignedTokenUpgradeData.s.sol \
	--sig "run(address,uint256, address,address,address,address,uint256)" \
	$(IMPLEMENTATION) $(VERSION) $(OWNER) $(BENEFICIARY1) $(BENEFICIARY2) $(BENEFICIARY3) $(MINT)

SALT?=0x0000000000000000000000000000000000000000000000000000000000000000
# Sepolia Safe CreateCall contract.
DEPLOYER?=0x9b35Af71d77eaf8d7e40252370304687390A1A52

aligned-token-create2: ## 🚀 Generate the create2 data for the AlignedToken proxy contract
	cd script/aligned_token && \
	forge script AlignedTokenCreate2.s.sol \
	--sig "run(uint256,bytes32,address)" \
	$(VERSION) $(SALT) $(DEPLOYER)

aligned-token-proxy-create2: ## 🚀 Generate the create2 data for the AlignedToken proxy contract
	cd script/aligned_token && \
	forge script AlignedTokenCreate2.s.sol \
	--sig "run(address,uint256, address,address,address,address,uint256,bytes32,address)" \
	$(IMPLEMENTATION) $(VERSION) $(OWNER) $(BENEFICIARY1) $(BENEFICIARY2) $(BENEFICIARY3) $(MINT) $(SALT) $(DEPLOYER

# Test targets

test-token: ## 🚀 Test the AlignedToken contract
	cast call $(ADDRESS) "name()(string)" --rpc-url http://localhost:8545
	cast call $(ADDRESS) "totalSupply()(uint256)" --rpc-url http://localhost:8545

test-airdrop: ## 🚀 Test the ClaimableAirdrop contract
	cast call $(ADDRESS) "hasClaimed(address)(bool)" --rpc-url http://localhost:8545

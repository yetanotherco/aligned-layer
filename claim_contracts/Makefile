.PHONY: help deploy-aligned-token-implementation deploy-aligned-token-proxy deploy-claimable-airdrop-implementation deploy-claimable-airdrop-proxy upgrade-aligned-token-implementation aligned-token-proxy-deploy-data aligned-token-init-data aligned-token-upgrade-data aligned-token-create2 aligned-token-proxy-create2


help: ## 📚 Show help for each of the Makefile recipes
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Deployments

RPC_URL?=http://localhost:8545
PRIVATE_KEY?=0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d

deploy-all: ## 🚀 Deploy all contracts
	cd script && forge script DeployAll.s.sol --private-key $(PRIVATE_KEY) --rpc-url $(RPC_URL) --broadcast

CONFIG?=example
deploy-token: ## 🚀 Deploy the token contract
	cd script && \
	forge script DeployAlignedToken.s.sol \
	--sig "run(string)" \
	$(CONFIG) \
	--private-key $(PRIVATE_KEY) \
	--rpc-url $(RPC_URL) \
	--broadcast

MINT?="1 ether"
OWNER?=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
BENEFICIARY1?=0x70997970C51812dc3A010C7d01b50e0d17dc79C8
BENEFICIARY2?=0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
BENEFICIARY3?=0x90F79bf6EB2c4f870365E785982E1f101E93b906

CLAIM_TIME_LIMIT?=1632960000
MERKLE_ROOT?=0x90076b5fb9a6c81d9fce83dfd51760987b8c49e7c861ea25b328e6e63d2cd3df
TOKEN_OWNER?=$(OWNER)

deploy-proxy-admin: ## 🚀 Deploy the ProxyAdmin contract
	cd script/proxy_admin && \
	forge script DeployProxyAdmin.s.sol \
	--sig "run(address)" \
	$(OWNER) \
	--rpc-url $(RPC_URL) \
	--private-key $(PRIVATE_KEY) \
	--broadcast

deploy-aligned-token-implementation: ## 🚀 Deploy the AlignedToken implementation contract
	cd script/aligned_token && \
	forge script DeployAlignedTokenImplementation.s.sol \
	--rpc-url $(RPC_URL) \
	--private-key $(PRIVATE_KEY) \
	--broadcast

deploy-aligned-token-proxy: ## 🚀 Deploy the AlignedToken proxy contract
	cd script/aligned_token && \
	forge script DeployAlignedTokenProxy.s.sol \
	--sig "run(address,address,address,address,address,address,uint256)" \
	$(PROXY_ADMIN) $(IMPLEMENTATION) $(OWNER) $(BENEFICIARY1) $(BENEFICIARY2) $(BENEFICIARY3) $(MINT) \
	--rpc-url $(RPC_URL) \
	--private-key $(PRIVATE_KEY) \
	--broadcast

deploy-claimable-airdrop-implementation: ## 🚀 Deploy the ClaimableAirdrop implementation contract
	cd script/claimable_airdrop && \
	forge script DeployClaimableAirdropImplementation.s.sol \
	--rpc-url $(RPC_URL) \
	--private-key $(PRIVATE_KEY) \
	--broadcast

deploy-claimable-airdrop-proxy: ## 🚀 Deploy the ClaimableAirdrop proxy contract
	cd script/claimable_airdrop && \
	forge script DeployClaimableAirdropProxy.s.sol \
	--sig "run(address,address,address,address,address,uint256,bytes32)" \
	$(PROXY_ADMIN) $(IMPLEMENTATION) $(OWNER) $(TOKEN) $(BENEFICIARY1) $(CLAIM_TIME_LIMIT) $(MERKLE_ROOT) \
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

# Misc

approve:
	cd script && \
	forge script ApproveERC20.s.sol \
	--sig "run(address,address,uint256)" \
	$(TOKEN) $(AIRDROP) $(AMOUNT) \
	--rpc-url $(RPC_URL) \
	--private-key $(HOLDER_PRIVATE_KEY) \
	--broadcast

# Test targets

test-token:
	cast call $(ADDRESS) "name()(string)" --rpc-url $(RPC_URL)
	cast call $(ADDRESS) "totalSupply()(uint256)" --rpc-url $(RPC_URL)

# The following target needs the proof API running on localhost:4000
AMOUNT_TO_CLAIM=$(shell curl -S -H "Content-Type: application/json" http://localhost:4000/api/proof/\$(CLAIMER) | jq -r .amount)
MERKLE_PROOF_TO_CLAIM=$(shell curl -S -H "Content-Type: application/json" http://localhost:4000/api/proof/\$(CLAIMER) | jq .proof | tr -d '"\n ')
test-claim:
	cast send $(AIRDROP) --private-key $(CLAIMER_PRIVATE_KEY) "claim(uint256,bytes32[])" $(AMOUNT_TO_CLAIM) "$(MERKLE_PROOF_TO_CLAIM)" --rpc-url $(RPC_URL)

test-claimed:
	cast call $(AIRDROP) "hasClaimed(address)(bool)" $(CLAIMER) --rpc-url $(RPC_URL)
	cast balance --erc20 $(TOKEN) $(CLAIMER) --rpc-url $(RPC_URL)

OWNER_PRIVATE_KEY?=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

test-pause:
	cast send $(AIRDROP) --private-key $(OWNER_PRIVATE_KEY) "pause()" --rpc-url $(RPC_URL)

test-unpause:
	cast send $(AIRDROP) --private-key $(OWNER_PRIVATE_KEY) "unpause()" --rpc-url $(RPC_URL)

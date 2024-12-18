.PHONY: help deploy-aligned-token-implementation deploy-aligned-token-proxy deploy-claimable-airdrop-implementation deploy-claimable-airdrop-proxy upgrade-aligned-token-implementation aligned-token-proxy-deploy-data aligned-token-init-data aligned-token-upgrade-data aligned-token-create2 aligned-token-proxy-create2


help: ## 📚 Show help for each of the Makefile recipes
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Deployments

RPC_URL?=http://localhost:8545
PRIVATE_KEY?=0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a

deploy-all: ## 🚀 Deploy all contracts
	cd script && \
	forge script DeployAll.s.sol \
	--sig "run(string)" \
	$(CONFIG) \
	--private-key $(PRIVATE_KEY) \
	--rpc-url $(RPC_URL) \
	--broadcast \
	-vvv

deploy-token: ## 🚀 Deploy the token contract
	cd script && \
	forge script DeployAlignedToken.s.sol \
	--sig "run(string)" \
	$(CONFIG) \
	--private-key $(PRIVATE_KEY) \
	--rpc-url $(RPC_URL) \
	--broadcast \
	--verbosity 3

deploy-token-prod: ## 🚀 Deploy the token contract
	cd script && \
	forge script DeployAlignedToken.s.sol \
	--sig "run(string)" \
	$(PROD_CONFIG) \
	--keystore $(KEYSTORE_PATH) \
	--rpc-url $(PROD_RPC_URL) \
	--broadcast \
	--verbosity 3

CONFIG?=anvil
deploy-claimable: ## 🚀 Deploy the token contract
	cd script && \
	forge script DeployClaimableAirdrop.s.sol \
	--sig "run(string)" \
	$(CONFIG) \
	--private-key $(PRIVATE_KEY) \
	--rpc-url $(RPC_URL) \
	--broadcast \
	--verbosity 3


update_token_proxy:
	@NEW_TOKEN_PROXY=$$(jq -r '.tokenProxy' "script-out/deployed_token_addresses.json") && \
	jq --arg new_proxy "$$NEW_TOKEN_PROXY" '.tokenProxy = $$new_proxy' $(CONFIG) > $(CONFIG).tmp \
	&& mv $(CONFIG).tmp $(CONFIG)

upgrade-token: ## 🚀 Upgrade the token contract
	cd script && \
	forge script UpgradeToken.s.sol \
	--sig "run(string)" \
	$(CONFIG) \
	--private-key $(PRIVATE_KEY) \
	--rpc-url $(RPC_URL) \
	--broadcast \
	--verbosity 3

# Upgrades

upgrade-aligned-token-implementation: ## 🚀 Upgrade the AlignedToken implementation contract
	cd script/aligned_token && \
	forge script UpgradeAlignedTokenImplementation.s.sol \
	--sig "function run(address,address,uint256,address,address,address,address,uint256)" \
	$(PROXY) $(IMPLEMENTATION) $(VERSION) $(OWNER) $(BENEFICIARY1) $(BENEFICIARY2) $(BENEFICIARY3) $(MINT)\
	--rpc-url $(RPC_URL) \
	--private-key $(PRIVATE_KEY) \
	--broadcast

# Test targets

test-token:
	cast call $(ADDRESS) "name()(string)" --rpc-url $(RPC_URL)
	cast call $(ADDRESS) "symbol()(string)" --rpc-url $(RPC_URL)
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

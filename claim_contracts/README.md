# Claimable Airdrop deployment steps

1. Create a `config.json` based on `config.example.json`. Make sure the following fields are correct:
    - salt
    - deployer
    - foundation
    - tokenProxy
    - limitTimestampToClaim
    - claimMerkleRoot
1. Set the token supplier private key in an env variable called `CLAIM_SUPPLIER_PRIVATE_KEY`
1. Run `make deploy-claimable`

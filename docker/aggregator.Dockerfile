FROM aligned_base AS builder

WORKDIR /aligned_layer

RUN go build -o ./aligned-layer-aggregator aggregator/cmd/main.go

FROM debian:bookworm-slim

WORKDIR /aggregator

COPY --from=builder /aligned_layer/aligned-layer-aggregator /usr/local/bin/aligned-layer-aggregator
COPY --from=builder /aligned_layer/config-files/config-aggregator-docker.yaml ./config-files/config-aggregator-docker.yaml
COPY --from=builder /aligned_layer/contracts/script/output/devnet/alignedlayer_deployment_output.json ./contracts/script/output/devnet/alignedlayer_deployment_output.json
COPY --from=builder /aligned_layer/contracts/script/output/devnet/eigenlayer_deployment_output.json ./contracts/script/output/devnet/eigenlayer_deployment_output.json
COPY --from=builder /aligned_layer/config-files/anvil.aggregator.ecdsa.key.json ./config-files/anvil.aggregator.ecdsa.key.json
COPY --from=builder /aligned_layer/config-files/anvil.aggregator.bls.key.json ./config-files/anvil.aggregator.bls.key.json

CMD ["aligned-layer-aggregator", "--config", "config-files/config-aggregator-docker.yaml"]

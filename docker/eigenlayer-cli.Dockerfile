FROM golang:1.22.2-bookworm

COPY config-files/ ./config-files

RUN go install github.com/Layr-Labs/eigenlayer-cli/cmd/eigenlayer@v0.12.0-beta

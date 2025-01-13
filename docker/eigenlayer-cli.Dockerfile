FROM golang:1.22.2-bookworm

COPY config-files/ ./config-files

# Install eigenlayer cli using downloading the binary
RUN curl -sSfL https://raw.githubusercontent.com/layr-labs/eigenlayer-cli/master/scripts/install.sh | sh -s -- v0.11.0

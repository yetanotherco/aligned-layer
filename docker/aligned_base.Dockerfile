FROM debian:bookworm-slim AS base

ARG BUILDARCH
ENV GO_VERSION=1.22.2

RUN apt update -y && apt upgrade -y
RUN apt install -y wget \
                   tar \
                   curl \
                   git \
                   make \
                   clang \
                   pkg-config \
                   openssl \
                   libssl-dev \
                   yq \
                   jq

RUN wget https://golang.org/dl/go$GO_VERSION.linux-${BUILDARCH}.tar.gz
RUN tar -C /usr/local -xzf go$GO_VERSION.linux-${BUILDARCH}.tar.gz
RUN rm go$GO_VERSION.linux-${BUILDARCH}.tar.gz
RUN apt clean -y
RUN rm -rf /var/lib/apt/lists/*
ENV PATH="/usr/local/go/bin:${PATH}"

# Install go deps
RUN go install github.com/maoueh/zap-pretty@latest
RUN go install github.com/ethereum/go-ethereum/cmd/abigen@latest
RUN go install github.com/Layr-Labs/eigenlayer-cli/cmd/eigenlayer@latest

# Install rust
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

# Install Mold Linker
RUN git clone --branch stable https://github.com/rui314/mold.git
WORKDIR mold
RUN ./install-build-deps.sh
RUN cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_CXX_COMPILER=c++ -B build
RUN cmake --build build -j24
RUN sudo cmake --build build --target install


# Build Aligned Layer
ENV RELEASE_FLAG=--release
ENV TARGET_REL_PATH=release
ENV CARGO_NET_GIT_FETCH_WITH_CLI=true
ENV RUSTFLAGS="-C link-arg=-fuse-ld=mold"
ENV PATH=$PATH:/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin

WORKDIR /aligned_layer
COPY operator/ /aligned_layer/operator/
COPY batcher/ /aligned_layer/batcher/

# build_sp1_linux
WORKDIR /aligned_layer/operator/sp1/lib
RUN cargo build ${RELEASE_FLAG}
RUN cp /aligned_layer/operator/sp1/lib/target/${TARGET_REL_PATH}/libsp1_verifier_ffi.so /aligned_layer/operator/sp1/lib/libsp1_verifier_ffi.so

# build_risc_zero_linux
WORKDIR /aligned_layer/operator/risc_zero/lib
RUN cargo build ${RELEASE_FLAG}
RUN cp /aligned_layer/operator/risc_zero/lib/target/${TARGET_REL_PATH}/librisc_zero_verifier_ffi.so /aligned_layer/operator/risc_zero/lib/librisc_zero_verifier_ffi.so

# build_sp1_linux_old
WORKDIR /aligned_layer/operator/sp1_old/lib
RUN cargo build ${RELEASE_FLAG}
RUN cp /aligned_layer/operator/sp1_old/lib/target/${TARGET_REL_PATH}/libsp1_verifier_old_ffi.so /aligned_layer/operator/sp1_old/lib/libsp1_verifier_old_ffi.so

# build_risc_zero_linux_old
WORKDIR /aligned_layer/operator/risc_zero_old/lib
RUN cargo build ${RELEASE_FLAG}
RUN cp /aligned_layer/operator/risc_zero_old/lib/target/${TARGET_REL_PATH}/librisc_zero_verifier_old_ffi.so /aligned_layer/operator/risc_zero_old/lib/librisc_zero_verifier_old_ffi.so

# build_merkle_tree_linux
WORKDIR /aligned_layer/operator/merkle_tree/lib
RUN cargo build ${RELEASE_FLAG}
RUN cp /aligned_layer/operator/merkle_tree/lib/target/${TARGET_REL_PATH}/libmerkle_tree.so /aligned_layer/operator/merkle_tree/lib/libmerkle_tree.so
RUN cp /aligned_layer/operator/merkle_tree/lib/target/${TARGET_REL_PATH}/libmerkle_tree.a /aligned_layer/operator/merkle_tree/lib/libmerkle_tree.a

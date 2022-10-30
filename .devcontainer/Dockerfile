# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.245.2/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version (use -bullseye variants on local arm64/Apple Silicon): 1, 1.19, 1.18, 1-bullseye, 1.19-bullseye, 1.18-bullseye, 1-buster, 1.19-buster, 1.18-buster
ARG VARIANT="1.19-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Choice] Node.js version: none, lts/*, 18, 16, 14
ARG NODE_VERSION="none"
RUN if [ "${NODE_VERSION}" != "none" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

# packages
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends wget cmake ninja-build g++ curl xz-utils nano hugo sudo zsh python3-pip

# WABT SOURCE
WORKDIR /workspaces/parigot
RUN mkdir deps
WORKDIR /workspaces/parigot/deps
RUN git clone -q  --recursive https://github.com/WebAssembly/wabt
WORKDIR /workspaces/parigot/deps/wabt
RUN git submodule update -q --init

# build all the code in wabt
RUN mkdir -p /workspaces/parigot/deps/wabt/build
WORKDIR /workspaces/parigot/deps/wabt/build
RUN cmake ..
RUN make
WORKDIR /workspaces/parigot/deps/wabt
RUN rm -rf build

#binary downloads, wasmtime 1.0.1
WORKDIR /workspaces/parigot/deps
RUN wget -q https://github.com/bytecodealliance/wasmtime/releases/download/v1.0.1/wasmtime-v1.0.1-x86_64-linux.tar.xz
RUN unxz wasmtime-v1.0.1-x86_64-linux.tar.xz
RUN tar xf wasmtime*.tar
RUN rm wasmtime*tar

# go installs
USER vscode
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest
RUN go install github.com/bufbuild/buf/cmd/buf@latest


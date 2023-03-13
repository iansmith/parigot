FROM arm64v8/ubuntu:jammy
RUN sed -i -e 's/archive.ubuntu.com\|security.ubuntu.com/old-releases.ubuntu.com/g' /etc/apt/sources.list
RUN apt-get update --fix-missing
RUN apt-get install -y wget git cmake ninja-build g++ python3 curl xz-utils nano \
    hugo sudo openjdk-18-jdk-headless zsh python3-pip ca-certificates gnupg file

RUN useradd -m -d /home/parigot -s /bin/bash -G sudo parigot
RUN usermod -aG sudo parigot

## we allow sudo for convenience if you need to modify something inside
## the container without going through this Dockerfile
RUN mkdir /home/parigot/deps
RUN mkdir /home/parigot/tmp
WORKDIR /home/parigot/deps
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# install go1.19
ENV GOLANG_VERSION="1.19.6"
ENV PROC_ARCH=arm64
RUN wget -q https://go.dev/dl/go${GOLANG_VERSION}.linux-${PROC_ARCH}.tar.gz
RUN tar xzf go*linux*tar.gz
RUN rm go*linux*.tar.gz 
RUN mv go go${GOLANG_VERSION}
ENV GOROOT=/home/parigot/tools/go${GOLANG_VERSION}

# get all the wabt project code
# WORKDIR /home/parigot/deps
# RUN git clone -q  --recursive https://github.com/WebAssembly/wabt
# WORKDIR /home/parigot/deps/wabt
# RUN git submodule update -q --init

# build all the code in wabt
# RUN mkdir /home/parigot/deps/wabt/build
# WORKDIR /home/parigot/deps/wabt/build
# RUN cmake ..
# RUN make
# WORKDIR /home/parigot/deps/wabt
# RUN rm -rf build

# wasmtime
ENV WASMTIME_VERSION=4.0.0
WORKDIR /home/parigot/deps
RUN wget -q https://github.com/bytecodealliance/wasmtime/releases/download/v${WASMTIME_VERSION}/wasmtime-v${WASMTIME_VERSION}-x86_64-linux.tar.xz
RUN unxz wasmtime*.tar.xz
RUN tar xf wasmtime*.tar
RUN rm wasmtime*tar

# get ready for symlinks
RUN mkdir /home/parigot/deps/bin
WORKDIR /home/parigot/deps/bin

# protocol buffers 
ENV PB_VERSION="21.12"
ENV PB_REL="https://github.com/protocolbuffers/protobuf/releases"
RUN curl -s -LO $PB_REL/download/v${PB_VERSION}/protoc-${PB_VERSION}-linux-x86_64.zip
RUN mkdir /home/parigot/tools
WORKDIR /home/parigot/tools
ENV PB_REL=""

# build then link the wabt tools in place
#RUN ln -s /home/parigot/deps/wabt /home/parigot/tools

# symlink all the single binaries in tools/bin
RUN mkdir /home/parigot/tools/bin
WORKDIR /home/parigot/tools/bin
RUN ln -s /usr/bin/python3 /home/parigot/tools/bin/python
RUN ln -s /home/parigot/deps/wasmtime-v${WASMTIME_VERSION}-x86_64-linux/wasmtime /home/parigot/tools/bin
RUN ln -s /home/parigot/deps/go${GOLANG_VERSION} /home/parigot/tools

# antlr 4.9 is a platform-neutral lib (java)
# ENV ANTLR_VERSION=4.11.1
# RUN mkdir -p /home/parigot/tools/lib
# WORKDIR /home/parigot/tools/lib
# RUN wget -q https://www.antlr.org/download/antlr-${ANTLR_VERSION}-complete.jar

ADD zshrc /home/parigot/.zshrc
ADD p10k.zsh /home/parigot/.p10k.zsh
ADD ohmyzsh /home/parigot/.oh-my-zsh
RUN chown -R parigot /home/parigot/.zshrc  /home/parigot/.p10k.zsh /home/parigot/.oh-my-zsh
RUN chgrp -R parigot /home/parigot/.zshrc /home/parigot/.p10k.zsh /home/parigot/.oh-my-zsh
#RUN chmod -R g+w /home/parigot/.zshrc /home/parigot/.p10k.zsh /home/parigot/.oh-my-zsh

RUN mkdir /var/run/parigot
RUN chown parigot /var/run/parigot
RUN chgrp parigot /var/run/parigot
#RUN chmod g+w /var/run/parigot

# safety
WORKDIR /home/parigot
RUN chown -R parigot .
RUN chgrp -R parigot .

#zsh
RUN chsh -s /usr/bin/zsh parigot

#temporary
USER parigot
#go progs
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/home/parigot/tools/go${GOLANG_VERSION}/bin:/home/parigot/.local/bin

ENV PROTOC_GO_VERSION=1.28.1
ENV PROTOC_GRPC_VERSION=1.51.0
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GO_VERSION}
#RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GRPC_VERSION}
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest
RUN go install github.com/bufbuild/buf/cmd/buf@latest
RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install honnef.co/go/tools/cmd/staticcheck@latest
RUN go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
#RUN go install github.com/onsi/ginkgo/v2/ginkgo

RUN ln -s /home/parigot/go/bin/* /home/parigot/tools/bin
RUN ln -s /workspaces/parigot/build/protoc-gen-parigot /home/parigot/tools/bin/protoc-gen-parigot

# add error correction

# set some vars
WORKDIR /home/parigot/src
ENV TOOLS=/home/parigot/tools
ENV GOTOOLS=/home/parigot/tools/go${GOLANG_VERSION}
ENV WABTTOOLS=/home/parigot/tools/wabt
ENV PARIGOT_TOOLS=/home/parigot/tools
ENV PARIGOT_IMPORT_PATH=/home/parigot/go/pkg/mod/github.com/gogo/protobuf@v1.3.2/protobuf
ENV PATH=/usr/local/sbin:/usr/sbin:/usr/bin:/sbin:/bin:$TOOLS/bin:$GOTOOLS/bin:$WABTTOOLS/bin:/home/parigot/.local/bin

ENTRYPOINT ["bash"]


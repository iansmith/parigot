FROM ubuntu:jammy-20221003
RUN apt-get update
RUN apt-get install -y wget git cmake ninja-build g++ python3 curl xz-utils nano hugo

RUN useradd -m -d /home/parigot -s /bin/bash -G sudo parigot
RUN usermod -aG sudo parigot

## we allow sudo for convenience if you need to modify something inside
## the container without going through this Dockerfile
RUN mkdir /home/parigot/deps
WORKDIR /home/parigot/deps
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

RUN wget -q https://go.dev/dl/go1.19.2.linux-amd64.tar.gz
RUN tar xzf go*linux*tar.gz
RUN rm go1.19.2.linux-amd64.tar.gz 
RUN mv go go1.19.2
ENV GOROOT=/home/parigot/tools/go1.19.2

RUN wget -q https://github.com/tinygo-org/tinygo/releases/download/v0.26.0/tinygo_0.26.0_amd64.deb
RUN dpkg -i tinygo_0.26.0_amd64.deb
RUN rm tinygo_0.26.0_amd64.deb

# for the tinygo config file related to undefined wasm symbols
RUN chgrp parigot /usr/local/lib/tinygo/targets
RUN chgrp parigot /usr/local/lib/tinygo/targets/wasm-undefined.txt
RUN chmod g+w /usr/local/lib/tinygo/targets
RUN chmod g+w /usr/local/lib/tinygo/targets/wasm-undefined.txt
RUN cp /usr/local/lib/tinygo/targets/wasm-undefined.txt /usr/local/lib/tinygo/targets/wasm-undefined.txt.orig
RUN chmod g+w /usr/local/lib/tinygo/targets/wasm-undefined.txt.orig

WORKDIR /home/parigot/deps
RUN git clone -q  --recursive https://github.com/WebAssembly/wabt
WORKDIR /home/parigot/deps/wabt
RUN git submodule update -q --init

RUN mkdir /home/parigot/deps/wabt/build
WORKDIR /home/parigot/deps/wabt/build
RUN cmake ..
RUN make
WORKDIR /home/parigot/deps/wabt
RUN rm -rf build

#binary downloads
WORKDIR /home/parigot/deps
RUN wget -q https://github.com/bytecodealliance/wasmtime/releases/download/v1.0.1/wasmtime-v1.0.1-x86_64-linux.tar.xz
RUN unxz wasmtime-v1.0.1-x86_64-linux.tar.xz
RUN tar xf wasmtime*.tar
RUN rm wasmtime*tar

RUN mkdir /home/parigot/deps/bin
WORKDIR /home/parigot/deps/bin

ENV PB_REL="https://github.com/protocolbuffers/protobuf/releases"
RUN curl -s -LO $PB_REL/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip
RUN mkdir /home/parigot/tools
WORKDIR /home/parigot/tools
RUN ln -s /home/parigot/deps/go1.19.2 /home/parigot/tools
RUN ln -s /home/parigot/deps/wabt /home/parigot/tools
RUN ln -s /usr/local/lib/tinygo /home/parigot/tools/tinygo0.26
ENV TINYGOROOT=/home/parigot/tools/tinygo0.26

ENV PB_REL=""

RUN mkdir /home/parigot/tools/bin
WORKDIR /home/parigot/tools/bin
RUN ln -s /usr/local/bin/tinygo /home/parigot/tools/bin
# hugo puts itself in /usr/bin, ugh!
#RUN ln -s /usr/local/bin/hugo /home/parigot/tools/bin
RUN ln -s /home/parigot/deps/wasmtime-v1.0.1-x86_64-linux/wasmtime /home/parigot/tools/bin

WORKDIR /home/parigot
RUN chown -R parigot .
RUN chgrp -R parigot .

#temporary
USER parigot
#go progs
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/home/parigot/deps/go1.19.2/bin

RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest
RUN go install github.com/bufbuild/buf/cmd/buf@latest
RUN ln -s /home/parigot/go/bin/* /home/parigot/tools/bin
RUN ln -s /home/parigot/src/build/protoc-gen-parigot /home/parigot/tools/bin/protoc-gen-parigot

# git config
echo '[safe]' > /home/parigot/gitconfig

# final touches
ADD enable-parigot /home/parigot/
ADD .gitconfig /home/parigot/

RUN mkdir /home/parigot/src
WORKDIR /home/parigot/src
ENV TOOLS=/home/parigot/tools
ENV GOTOOLS=/home/parigot/tools/go1.19.2
ENV WABTTOOLS=/home/parigot/tools/wabt

ENV PATH=/usr/local/sbin:/usr/sbin:/usr/bin:/sbin:/bin:$TOOLS/bin:$GOTOOLS/bin:$WABTTOOLS/bin

ENTRYPOINT ["bash"]


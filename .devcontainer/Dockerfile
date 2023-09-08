
FROM ubuntu:23.04
RUN apt-get update --fix-missing && apt-get upgrade -y
RUN apt-get clean
RUN apt-get install -y curl wget git xz-utils nano make
RUN apt-get install -y sudo ca-certificates file dpkg 
RUN apt-get install -y zip unzip zsh clang-16

RUN useradd -m -d /home/parigot -s /bin/bash -G sudo parigot
RUN usermod -aG sudo parigot

# we allow sudo
RUN mkdir /home/parigot/deps
RUN mkdir /home/parigot/tmp
WORKDIR /home/parigot/deps
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

## golang version
ENV GOLANG_VERSION_HOST="1.21.0"
ENV GOLANG_VERSION_HOST_SHORT="1.21"
## wasmtime version
ENV WASMTIME_VERSION=10.0.1


# hugo version for website
ENV HUGO_VERSION=0.114.1
# protoc because some plugins don't work with buf
ENV PROTOC_VERSION=23.3

#figure out ARCH and download correct tarballs
RUN set -e; \
	unamearch="$(uname -m)" ; \
	unameos="$(uname -s)"; \
	if [ "$unameos" != "Linux" ] ; \
	then \
	echo "build failed, operating system inside container not linux -- $unameos"; return 1;\
	fi ;\
	case $unamearch in \
		aarch64) wget -O go.tar.gz -q https://go.dev/dl/go${GOLANG_VERSION_HOST}.linux-arm64.tar.gz ; \
		wget -O wasmtime.tar.xz -q https://github.com/bytecodealliance/wasmtime/releases/download/v${WASMTIME_VERSION}/wasmtime-v${WASMTIME_VERSION}-aarch64-linux.tar.xz ; \
		wget -O hugo_${HUGO_VERSION}.deb -q https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_linux-arm64.deb; \
		wget -O protoc-${PROTOC_VERSION}.zip https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-aarch_64.zip;; \
		x86_64) wget -O go.tar.gz -q https://go.dev/dl/go${GOLANG_VERSION_HOST}.linux-amd64.tar.gz ; \
		wget -O wasmtime.tar.xz -q https://github.com/bytecodealliance/wasmtime/releases/download/v${WASMTIME_VERSION}/wasmtime-v${WASMTIME_VERSION}-x86_64-linux.tar.xz ; \
		wget -O hugo_${HUGO_VERSION}.deb -q https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_linux-amd64.deb; \
		wget -O protoc-${PROTOC_VERSION}.zip https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip;; \
        *) echo "unknown architecture type, or not running in a linux container: aborting..."; return 1 ;; \
	esac;

WORKDIR /home/parigot/deps

ENV GOROOT=
ENV GOROOT_BOOTSTRAP=

# install Hugo
WORKDIR /home/parigot/deps
RUN dpkg -i hugo_${HUGO_VERSION}.deb
RUN rm -f hugo_${HUGO_VERSION}.deb

# install protoc
RUN mkdir /home/parigot/deps/protoc-${PROTOC_VERSION}
WORKDIR /home/parigot/deps/protoc-${PROTOC_VERSION}
RUN mv ../protoc-${PROTOC_VERSION}.zip .
RUN unzip protoc-${PROTOC_VERSION}.zip

# install golang
WORKDIR /home/parigot/deps
RUN tar xvf go.tar.gz
RUN mv go go${GOLANG_VERSION_HOST}
RUN rm go.tar.gz

## install wasmtime
WORKDIR /home/parigot/deps
RUN unxz wasmtime.tar.xz
RUN tar xf wasmtime.tar
RUN rm wasmtime.tar

# protocol buffers 
ENV PB_VERSION="21.12"
ENV PB_REL="https://github.com/protocolbuffers/protobuf/releases"
#RUN curl -s -LO $PB_REL/download/v${PB_VERSION}/protoc-${PB_VERSION}-linux-x86_64.zip
RUN mkdir /home/parigot/tools
WORKDIR /home/parigot/tools
ENV PB_REL=""

## tinygo does not support enough of reflect package to run "standard" go protobufs
## using the google tools ... there are other implementation of protobufs that
## do work 
## https://github.com/knqyf263/go-plugin
## https://github.com/planetscale/vtprotobuf

# symlink all the single binaries in tools/bin
RUN mkdir /home/parigot/tools/bin
WORKDIR /home/parigot/tools/bin

#RUN ln -s /home/parigot/deps/go${GOLANG_VERSION_GUEST} /home/parigot/tools
#RUN ln -s /home/parigot/deps/go${GOLANG_VERSION_GUEST}/bin/go /home/parigot/tools/bin/go${GOLANG_VERSION_GUEST}

RUN ln -s /usr/bin/python3 /home/parigot/tools/bin/python
RUN ln -s /home/parigot/deps/go${GOLANG_VERSION_HOST} /home/parigot/tools
RUN ln -s /home/parigot/deps/go${GOLANG_VERSION_HOST}/bin/go /home/parigot/tools/bin/go${GOLANG_VERSION_HOST}
## if something called "go" is not in your path, the dev container will not run
RUN ln -s /home/parigot/deps/go${GOLANG_VERSION_HOST}/bin/go /home/parigot/tools/bin/go

RUN ln -s /usr/lib/llvm-16/bin/clang /home/parigot/tools/bin/clang 
RUN ln -s /usr/lib/llvm-16/bin/clang /home/parigot/tools/bin/cc 

RUN ln -s /usr/local/bin/hugo /home/parigot/tools/bin

RUN ln -s /home/parigot/deps/protoc-${PROTOC_VERSION}/bin/protoc /home/parigot/tools/bin

# RUN ln -s /usr/local/tinygo/build/tinygo /home/parigot/tools/bin
# RUN ln -s /usr/lib/llvm-15/bin/clang /home/parigot/tools/bin
# RUN ln -s /usr/lib/llvm-15/bin/clang /usr/bin  # compiled in somewhere


# docker socket dir?
RUN mkdir /var/run/parigot
RUN chown parigot /var/run/parigot /home/parigot
RUN chgrp parigot /var/run/parigot /home/parigot
#RUN chmod g+w /var/run/parigot

# safety
#WORKDIR /home/parigot
#RUN chown -R parigot .
#RUN chgrp -R parigot .
RUN mkdir -p /home/parigot/tools/bin
RUN chown -R parigot /home/parigot/tools
RUN chgrp -R parigot /home/parigot/tools
RUN mkdir -p /home/parigot/deps
RUN chown -R parigot /home/parigot/deps
RUN chgrp -R parigot /home/parigot/deps

#zsh    
RUN chsh -s /usr/bin/zsh parigot

#temporary PATH hack
USER parigot
ENV PATH=/workspaces/parigot/build:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/home/parigot/tools/bin


ENV PROTOC_GO_VERSION=1.28.1
ENV PROTOC_GRPC_VERSION=1.51.0
RUN  ls -l /home/parigot/deps/
RUN echo go${GOLANG_VERSION_HOST} install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go${GOLANG_VERSION_HOST} install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go${GOLANG_VERSION_HOST} install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GO_VERSION}
RUN go${GOLANG_VERSION_HOST} install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest
RUN go${GOLANG_VERSION_HOST} install github.com/bufbuild/buf/cmd/buf@v1.26
RUN go${GOLANG_VERSION_HOST} install golang.org/x/tools/gopls@latest
RUN go${GOLANG_VERSION_HOST} install github.com/go-delve/delve/cmd/dlv@latest
RUN go${GOLANG_VERSION_HOST} install honnef.co/go/tools/cmd/staticcheck@latest
RUN go${GOLANG_VERSION_HOST} install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.20.0
RUN go${GOLANG_VERSION_HOST} install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
RUN go${GOLANG_VERSION_HOST} install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go${GOLANG_VERSION_HOST} install github.com/tetratelabs/wazero/cmd/wazero@v1.5.0
RUN go${GOLANG_VERSION_HOST} install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest
#RUN go${GOLANG_VERSION_HOST} install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GRPC_VERSION}
#RUN go${GOLANG_VERSION_HOST} install github.com/onsi/ginkgo/v2/ginkgo

RUN ln -s /home/parigot/go/bin/* /home/parigot/tools/bin
RUN ln -s /workspaces/parigot/build/protoc-gen-parigot /home/parigot/tools/bin/protoc-gen-parigot

RUN mkdir ~/.ssh
RUN chmod 700 ~/.ssh
RUN ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

# set some vars
WORKDIR /workspaces/parigot
#ENV TOOLS=/home/parigot/tools
ENV WABTTOOLS=/home/parigot/tools/wabt
ENV TOOLS=/home/parigot/tools
#ENV PARIGOT_IMPORT_PATH=/home/parigot/go/pkg/mod/github.com/gogo/protobuf@v1.3.2/protobuf
ENV PATH=/workspaces/parigot/build:/usr/local/sbin:/usr/sbin:/usr/bin:/sbin:/bin:$TOOLS/bin:$WABTTOOLS/bin
ENV CC=clang
ENV GOFLAGS="-buildvcs=false"

ENTRYPOINT ["bash"]


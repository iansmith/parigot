#
# sample use of this file
# docker buildx build  --target parigot-base --platform=linux/amd64,linux/arm64 --push -t iansmith/parigot-base:0.4 .
#
FROM ubuntu:mantic as build

## without this the macros that define pread64 => pread are not used to enable
## the C code to build on alpine
ENV CGO_ENABLED=1
ENV CGO_CPPFLAGS="-D_LARGEFILE64_SOURCE"

RUN apt-get update
RUN apt-get install -y golang build-essential sqlite3 ca-certificates

ENV GOLANG_VERSION=1.21
RUN mkdir -p ~/go/bin
RUN ln -s /usr/bin/go /root/go/bin/go${GOLANG_VERSION}

ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/root/go/bin:/workspaces/parigot/build

RUN which go
RUN ls -l /root/go/bin

ENV PROTOC_VERSION_GO=1.31
RUN go${GOLANG_VERSION} install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go${GOLANG_VERSION} install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_VERSION_GO}
RUN go${GOLANG_VERSION} install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest
RUN go${GOLANG_VERSION} install github.com/bufbuild/buf/cmd/buf@v1.26
RUN go${GOLANG_VERSION} install golang.org/x/tools/gopls@latest
RUN go${GOLANG_VERSION} install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
RUN go${GOLANG_VERSION} install github.com/tetratelabs/wazero/cmd/wazero@v1.5.0
RUN go${GOLANG_VERSION} install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/root/go/bin:/workspaces/parigot/build

FROM ubuntu:mantic as parigot-base
ENV GOFLAGS="-buildvcs=false"
ENV CGO_CPPFLAGS="-D_LARGEFILE64_SOURCE"
ENV BUILDKIT_PROGRESS=plain
ENV CGO_ENABLED=1
ENV CC=/usr/bin/gcc
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/root/go/bin:/workspaces/parigot/build

#RUN echo 'https://dl-cdn.alpinelinux.org/alpine/edge/testing' >> /etc/apk/repositories 
#RUN apk update
#RUN apk add go make sqlite protobuf git openssh 
RUN apt-get update
RUN apt-get install -y golang build-essential sqlite3 ca-certificates git  libprotobuf-dev
COPY --from=build /root/go/bin /root/go/bin
RUN echo 'export CGO_ENABLED=1' > /root/.bashrc
RUN echo 'export CGO_CPPFLAGS=-D_LARGEFILE64_SOURCE' >> /root/.bashrc
RUN echo 'expath PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/root/go/bin:/workspaces/parigot/build' >> /root/.bashrc

RUN openssl version

FROM alpine:edge as parigot-openslice
RUN echo 'https://dl-cdn.alpinelinux.org/alpine/edge/testing' >> /etc/apk/repositories 
RUN apk update
RUN apk add go make etcd sqlite sqlite-dev rabbitmq-server protobuf git openssh 
COPY --from=parigot-base /root/go/bin /root/go/bin
ENV GOFLAGS="-buildvcs=false"
ENV BUILDKIT_PROGRESS=plain
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/root/go/bin:/workspaces/parigot/build


FROM ubuntu:jammy-20221003
RUN apt-get update
RUN apt-get install --fix-missing -y wget git cmake ninja-build g++ python3 curl xz-utils nano \
    hugo sudo openjdk-18-jdk-headless zsh python3-pip ca-certificates gnupg

RUN useradd -m -d /home/parigot -s /bin/bash -G sudo parigot
RUN usermod -aG sudo parigot

## we allow sudo for convenience if you need to modify something inside
## the container without going through this Dockerfile
RUN mkdir /home/parigot/deps
RUN mkdir /home/parigot/tmp
WORKDIR /home/parigot/deps
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# install go1.19
RUN wget -q https://go.dev/dl/go1.19.4.linux-amd64.tar.gz
RUN tar xzf go*linux*tar.gz
RUN rm go1.19.4.linux-amd64.tar.gz 
RUN mv go go1.19.4
ENV GOROOT=/home/parigot/tools/go1.19.4

# get all the wabt project code
WORKDIR /home/parigot/deps
RUN git clone -q  --recursive https://github.com/WebAssembly/wabt
WORKDIR /home/parigot/deps/wabt
RUN git submodule update -q --init

# build all the code in wabt
RUN mkdir /home/parigot/deps/wabt/build
WORKDIR /home/parigot/deps/wabt/build
RUN cmake ..
RUN make
WORKDIR /home/parigot/deps/wabt
RUN rm -rf build

#binary downloads, wasmtime 1.0.1
WORKDIR /home/parigot/deps
RUN wget -q https://github.com/bytecodealliance/wasmtime/releases/download/v3.0.1/wasmtime-v3.0.1-x86_64-linux.tar.xz
RUN unxz wasmtime-*.xz
RUN tar xf wasmtime*.tar
RUN rm wasmtime*tar

# get ready for symlinks
RUN mkdir /home/parigot/deps/bin
WORKDIR /home/parigot/deps/bin

# protocol buffers 3.15
ENV PB_REL="https://github.com/protocolbuffers/protobuf/releases"
RUN curl -s -LO $PB_REL/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip
RUN mkdir /home/parigot/tools
WORKDIR /home/parigot/tools

# link the 3 big directories of stuff go,tinygo, and wabt
RUN ln -s /home/parigot/deps/go1.19.4 /home/parigot/tools
RUN ln -s /home/parigot/deps/wabt /home/parigot/tools
ENV PB_REL=""

# symlink all the single binaries in tools/bin
RUN mkdir /home/parigot/tools/bin
WORKDIR /home/parigot/tools/bin
RUN ln -s /usr/local/bin/tinygo /home/parigot/tools/bin
# hugo puts itself in /usr/bin, ugh!
#RUN ln -s /usr/local/bin/hugo /home/parigot/tools/bin
RUN ln -s /home/parigot/deps/wasmtime-v3.0.1-x86_64-linux/wasmtime /home/parigot/tools/bin

# antlr 4.9 is a platform-neutral lib (java)
RUN mkdir -p /home/parigot/tools/lib
WORKDIR /home/parigot/tools/lib
RUN wget -q https://www.antlr.org/download/antlr-4.9.3-complete.jar

# dotfiles for convenience
RUN mkdir /home/parigot/dotfiles
RUN chown parigot /home/parigot/dotfiles
RUN chgrp parigot /home/parigot/dotfiles
RUN mkdir /home/parigot/dotfiles/ssh
RUN chgrp parigot /home/parigot/dotfiles/ssh
RUN chown parigot /home/parigot/dotfiles/ssh
RUN chmod 700 /home/parigot/dotfiles/ssh
ADD recv.sh /home/parigot/dotfiles
RUN chmod +x /home/parigot/dotfiles/recv.sh

#RUN ln -s /home/parigot/ssh /home/parigot/.ssh
#RUN ln -s /home/parigot/dotfiles/gitconfig /home/parigot/.gitconfig
#RUN ln -s /home/parigot/dotfiles/zshrc /home/parigot/.zshrc
#RUN ln -s /home/parigot/dotfiles/oh-my-zsh /home/parigot/.oh-my-zsh

# tmp dir
RUN chgrp parigot /home/parigot/tmp
RUN chmod g+w /home/parigot/tmp

# socket dir
RUN mkdir /var/run/parigot
RUN chgrp parigot /var/run/parigot
RUN chmod g+w /var/run/parigot


# source directory
RUN mkdir /home/parigot/src

# safety
WORKDIR /home/parigot
RUN chown -R parigot .
RUN chgrp -R parigot .
RUN chsh -s /usr/bin/zsh parigot

# sample enable script
ADD enable-parigot /home/parigot/
RUN chown parigot /home/parigot/enable-parigot
RUN chgrp parigot /home/parigot/enable-parigot

#temporary
USER parigot
#go progs
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/home/parigot/deps/go1.19.4/bin:/home/parigot/.local/bin

RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest
RUN go install github.com/bufbuild/buf/cmd/buf@latest
RUN ln -s /home/parigot/go/bin/* /home/parigot/tools/bin
RUN ln -s /home/parigot/src/build/protoc-gen-parigot /home/parigot/tools/bin/protoc-gen-parigot

# add error correction
RUN pip3 install thefuck

# set some vars
WORKDIR /home/parigot/src
ENV TOOLS=/home/parigot/tools
ENV GOTOOLS=/home/parigot/tools/go1.19.4
ENV WABTTOOLS=/home/parigot/tools/wabt

ENV PATH=/usr/local/sbin:/usr/sbin:/usr/bin:/sbin:/bin:$TOOLS/bin:$GOTOOLS/bin:$WABTTOOLS/bin:/home/parigot/.local/bin

ENTRYPOINT ["bash"]


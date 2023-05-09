FROM tinygo/tinygo

## gets mounted there
WORKDIR /home/tinygo
RUN mkdir parigot

RUN sudo apt-get update
RUN sudo apt-get install -y clang-13  libclang-common-11-dev libclang-cpp11 libclang1-11 libllvm11 libomp-11-dev libomp5-11 llvm-11 llvm-11-dev llvm-11-runtime llvm-11-tools sqlite3

RUN sudo ln -s /usr/bin/clang-13 /usr/bin/clang

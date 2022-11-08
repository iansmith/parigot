FROM ubuntu:jammy-20221101
RUN apt-get update
ADD runner /runner
ADD data.wasm /
EXPOSE 13339/udp
ENTRYPOINT ["/runner","/data.wasm"]

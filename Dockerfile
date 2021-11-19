# compile fx-core
FROM golang:1.17-alpine as builder

# default mainnet
ARG NETWORK=mainnet

COPY . /app

RUN apk add --no-cache git build-base linux-headers && \
    cd /app && FX_BUILD_OPTIONS=${NETWORK} make go-build

# build fx-core
FROM alpine:latest

WORKDIR root

COPY --from=builder /app/build/bin/fxcored /usr/bin/fxcored

EXPOSE 26656/tcp 26657/tcp 26660/tcp 9090/tcp 1317/tcp

VOLUME ["/root"]

ENTRYPOINT ["fxcored"]
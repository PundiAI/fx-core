# compile fx-core
FROM golang:1.18.2-alpine3.16 as builder

# default mainnet
ARG NETWORK=mainnet

RUN apk add --no-cache git build-base linux-headers

WORKDIR /app

# download and cache go mod
COPY ./go.* ./
RUN go env -w GO111MODULE=on && go mod download

COPY . .

RUN FX_BUILD_OPTIONS=${NETWORK} make build

# build fx-core
FROM alpine:3.16

WORKDIR root

COPY --from=builder /app/build/bin/fxcored /usr/bin/fxcored

EXPOSE 26656/tcp 26657/tcp 26660/tcp 9090/tcp 1317/tcp 8545/tcp 8546/tcp

VOLUME ["/root"]

ENTRYPOINT ["fxcored"]
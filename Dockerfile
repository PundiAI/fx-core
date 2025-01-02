FROM golang:1.23.1-alpine3.19 as builder

RUN apk add --no-cache git build-base linux-headers binutils-gold

WORKDIR /app

COPY . .

RUN make build

FROM alpine:3.19

WORKDIR root

COPY --from=builder /app/build/bin/fxcored /usr/bin/fxcored

EXPOSE 26656/tcp 26657/tcp 26660/tcp 9090/tcp 1317/tcp 8545/tcp 8546/tcp

VOLUME ["/root"]

ENTRYPOINT ["fxcored"]

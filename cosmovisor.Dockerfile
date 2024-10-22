FROM golang:1.23.0-alpine3.19 as builder

RUN apk add --no-cache git build-base linux-headers

RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.6.0

WORKDIR /app

COPY . .

RUN make build

FROM ghcr.io/functionx/fxcorevisor:7.5.0 as fxv7_5

FROM alpine:3.19

WORKDIR /root

ENV DAEMON_HOME=/root/.fxcore
ENV DAEMON_NAME=fxcored
ENV DAEMON_ALLOW_DOWNLOAD_BINARIES=false
ENV DAEMON_DOWNLOAD_MUST_HAVE_CHECKSUM=false
ENV DAEMON_RESTART_AFTER_UPGRADE=true
ENV DAEMON_RESTART_DELAY=1s
ENV DAEMON_POLL_INTERVAL=3s
ENV UNSAFE_SKIP_BACKUP=true
ENV DAEMON_PREUPGRADE_MAX_RETRIES=3
ENV COSMOVISOR_DISABLE_LOGS=false
ENV COSMOVISOR_COLOR_LOGS=true

ENV PATH="/root/.fxcore/cosmovisor/current/bin:${PATH}"

COPY --from=fxv7_5 /root/.fxcore/cosmovisor /root/.fxcore/cosmovisor

COPY --from=builder /go/bin/cosmovisor /usr/bin/cosmovisor
COPY --from=builder /app/build/bin/fxcored /root/.fxcore/cosmovisor/upgrades/v8.0.x/bin/fxcored

RUN cosmovisor init /root/.fxcore/cosmovisor/genesis/bin/fxcored

EXPOSE 26656/tcp 26657/tcp 26660/tcp 9090/tcp 1317/tcp 8545/tcp 8546/tcp

VOLUME ["/root"]

ENTRYPOINT ["cosmovisor", "run"]


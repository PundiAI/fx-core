FROM golang:1.21.4-alpine3.18 as builder

RUN apk add --no-cache git build-base linux-headers

RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

WORKDIR /app

COPY . .

RUN make build

FROM ghcr.io/functionx/fxcorevisor-lite:7.4.0-rc6 as fxv7_4

FROM alpine:3.18

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

COPY --from=fxv7_4 /root/.fxcore/cosmovisor/upgrades /root/.fxcore/cosmovisor/upgrades

COPY --from=builder /go/bin/cosmovisor /usr/bin/cosmovisor
COPY --from=builder /app/build/bin/fxcored /usr/bin/fxcored
COPY --from=builder /app/build/bin/fxcored /root/.fxcore/cosmovisor/upgrades/v7.5.x/bin/fxcored

RUN cosmovisor init /root/.fxcore/cosmovisor/upgrades/v6.0.x/bin/fxcored

EXPOSE 26656/tcp 26657/tcp 26660/tcp 9090/tcp 1317/tcp 8545/tcp 8546/tcp

VOLUME ["/root"]

ENTRYPOINT ["cosmovisor", "run"]

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/config"

	fxconfig "github.com/functionx/fx-core/server/config"
)

func Test_configTomlConfig_output(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.BaseConfig.Moniker = "anonymous"
	c := configTomlConfig{config: cfg}
	err := c.output(func(out []byte) error {
		assert.Equal(t, tmConfigJson, string(out))
		return nil
	})
	assert.NoError(t, err)
}

func Test_appTomlConfig_output(t *testing.T) {
	c := appTomlConfig{config: fxconfig.DefaultConfig()}
	err := c.output(func(out []byte) error {
		assert.Equal(t, appConfigJson, string(out))
		return nil
	})
	assert.NoError(t, err)
}

const tmConfigJson = `{
  "abci": "socket",
  "consensus": {
    "create_empty_blocks": true,
    "create_empty_blocks_interval": 0,
    "double_sign_check_height": 0,
    "home": "",
    "peer_gossip_sleep_duration": 100000000,
    "peer_query_maj23_sleep_duration": 2000000000,
    "skip_timeout_commit": false,
    "timeout_commit": 1000000000,
    "timeout_precommit": 1000000000,
    "timeout_precommit_delta": 500000000,
    "timeout_prevote": 1000000000,
    "timeout_prevote_delta": 500000000,
    "timeout_propose": 3000000000,
    "timeout_propose_delta": 500000000,
    "wal_file": "data/cs.wal/wal"
  },
  "db_backend": "goleveldb",
  "db_dir": "data",
  "fast_sync": true,
  "fastsync": {
    "version": "v0"
  },
  "filter_peers": false,
  "genesis_file": "config/genesis.json",
  "home": "",
  "instrumentation": {
    "max_open_connections": 3,
    "namespace": "tendermint",
    "prometheus": false,
    "prometheus_listen_addr": ":26660"
  },
  "log_format": "plain",
  "log_level": "info",
  "mempool": {
    "broadcast": true,
    "cache_size": 10000,
    "home": "",
    "keep-invalid-txs-in-cache": false,
    "max_batch_bytes": 0,
    "max_tx_bytes": 1048576,
    "max_txs_bytes": 1073741824,
    "recheck": true,
    "size": 5000,
    "wal_dir": ""
  },
  "moniker": "anonymous",
  "node_key_file": "config/node_key.json",
  "p2p": {
    "addr_book_file": "config/addrbook.json",
    "addr_book_strict": true,
    "allow_duplicate_ip": false,
    "dial_timeout": 3000000000,
    "external_address": "",
    "flush_throttle_timeout": 100000000,
    "handshake_timeout": 20000000000,
    "home": "",
    "laddr": "tcp://0.0.0.0:26656",
    "max_num_inbound_peers": 40,
    "max_num_outbound_peers": 10,
    "max_packet_msg_payload_size": 1024,
    "persistent_peers": "",
    "persistent_peers_max_dial_period": 0,
    "pex": true,
    "private_peer_ids": "",
    "recv_rate": 5120000,
    "seed_mode": false,
    "seeds": "",
    "send_rate": 5120000,
    "test_dial_fail": false,
    "test_fuzz": false,
    "test_fuzz_config": {
      "Mode": 0,
      "MaxDelay": 3000000000,
      "ProbDropRW": 0.2,
      "ProbDropConn": 0,
      "ProbSleep": 0
    },
    "unconditional_peer_ids": "",
    "upnp": false
  },
  "priv_validator_key_file": "config/priv_validator_key.json",
  "priv_validator_laddr": "",
  "priv_validator_state_file": "data/priv_validator_state.json",
  "proxy_app": "tcp://127.0.0.1:26658",
  "rpc": {
    "cors_allowed_headers": [
      "Origin",
      "Accept",
      "Content-Type",
      "X-Requested-With",
      "X-Server-Time"
    ],
    "cors_allowed_methods": [
      "HEAD",
      "GET",
      "POST"
    ],
    "cors_allowed_origins": [],
    "experimental_close_on_slow_client": false,
    "experimental_subscription_buffer_size": 200,
    "experimental_websocket_write_buffer_size": 200,
    "grpc_laddr": "",
    "grpc_max_open_connections": 900,
    "home": "",
    "laddr": "tcp://127.0.0.1:26657",
    "max_body_bytes": 1000000,
    "max_header_bytes": 1048576,
    "max_open_connections": 900,
    "max_subscription_clients": 100,
    "max_subscriptions_per_client": 5,
    "pprof_laddr": "",
    "timeout_broadcast_tx_commit": 10000000000,
    "tls_cert_file": "",
    "tls_key_file": "",
    "unsafe": false
  },
  "statesync": {
    "chunk_fetchers": 4,
    "chunk_request_timeout": 10000000000,
    "discovery_time": 15000000000,
    "enable": false,
    "rpc_servers": null,
    "temp_dir": "",
    "trust_hash": "",
    "trust_height": 0,
    "trust_period": 604800000000000
  },
  "tx_index": {
    "indexer": "kv",
    "psql-conn": ""
  }
}`

const appConfigJson = `{
  "api": {
    "address": "tcp://0.0.0.0:1317",
    "enable": false,
    "enabled-unsafe-cors": false,
    "max-open-connections": 1000,
    "rpc-max-body-bytes": 1000000,
    "rpc-read-timeout": 10,
    "rpc-write-timeout": 0,
    "swagger": false
  },
  "evm": {
    "max-tx-gas-wanted": 500000,
    "tracer": ""
  },
  "grpc": {
    "address": "0.0.0.0:9090",
    "enable": true
  },
  "halt-height": 0,
  "halt-time": 0,
  "iavl-cache-size": 781250,
  "index-events": [],
  "inter-block-cache": true,
  "json-rpc": {
    "address": "0.0.0.0:8545",
    "api": [
      "eth",
      "net",
      "web3"
    ],
    "block-range-cap": 10000,
    "enable": true,
    "evm-timeout": 5000000000,
    "feehistory-cap": 100,
    "filter-cap": 200,
    "gas-cap": 25000000,
    "http-idle-timeout": 120000000000,
    "http-timeout": 30000000000,
    "logs-cap": 10000,
    "txfee-cap": 1,
    "ws-address": "0.0.0.0:8546"
  },
  "min-retain-blocks": 0,
  "minimum-gas-prices": "",
  "pruning": "default",
  "pruning-interval": "0",
  "pruning-keep-every": "0",
  "pruning-keep-recent": "0",
  "state-sync": {
    "snapshot-interval": 0,
    "snapshot-keep-recent": 2
  },
  "telemetry": {
    "enable-hostname": false,
    "enable-hostname-label": false,
    "enable-service-label": false,
    "enabled": false,
    "global-labels": [],
    "prometheus-retention-time": 0,
    "service-name": ""
  },
  "tls": {
    "certificate-path": "",
    "key-path": ""
  }
}`

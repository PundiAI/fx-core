package config

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/config"

	fxtypes "github.com/functionx/fx-core/v7/types"
)

func Test_output(t *testing.T) {
	type args struct {
		ctx     client.Context
		content interface{}
	}
	clientCtx := func() client.Context {
		return client.Context{
			Output:       new(bytes.Buffer),
			OutputFormat: "json",
		}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "app.toml output grpc.enable",
			args: args{
				ctx:     clientCtx(),
				content: true,
			},
		},
		{
			name: "app.toml output bypass-min-fee.msg-types empty",
			args: args{
				ctx:     clientCtx(),
				content: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, output(tt.args.ctx, tt.args.content))
			assert.Equal(t, tt.args.ctx.Output.(*bytes.Buffer).String(), fmt.Sprintf("%v\n", tt.args.content))
		})
	}
}

func Test_configTomlConfig_output(t *testing.T) {
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
    "namespace": "cometbft",
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
    "ttl-duration": 0,
    "ttl-num-blocks": 0,
    "version": "v0",
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
}
`

	cfg := config.DefaultConfig()
	cfg.BaseConfig.Moniker = "anonymous"
	c := configTomlConfig{config: cfg}
	buf := new(bytes.Buffer)
	clientCtx := client.Context{
		Output:       buf,
		OutputFormat: "json",
	}
	assert.NoError(t, c.output(clientCtx))
	assert.Equal(t, tmConfigJson, buf.String())
}

func Test_appTomlConfig_output(t *testing.T) {
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
  "app-db-backend": "",
  "bypass-min-fee": {
    "msg-max-gas-usage": 300000,
    "msg-types": []
  },
  "evm": {
    "max-tx-gas-wanted": 0,
    "tracer": ""
  },
  "grpc": {
    "address": "0.0.0.0:9090",
    "enable": true,
    "max-recv-msg-size": 10485760,
    "max-send-msg-size": 2147483647
  },
  "grpc-web": {
    "address": "0.0.0.0:9091",
    "enable": true,
    "enable-unsafe-cors": false
  },
  "halt-height": 0,
  "halt-time": 0,
  "iavl-cache-size": 781250,
  "iavl-disable-fastnode": false,
  "iavl-lazy-loading": false,
  "index-events": [],
  "inter-block-cache": true,
  "json-rpc": {
    "address": "127.0.0.1:8545",
    "allow-unprotected-txs": false,
    "api": [
      "eth",
      "net",
      "web3"
    ],
    "block-range-cap": 10000,
    "enable": true,
    "enable-indexer": false,
    "evm-timeout": 5000000000,
    "feehistory-cap": 100,
    "filter-cap": 200,
    "fix-revert-gas-refund-height": 0,
    "gas-cap": 30000000,
    "http-idle-timeout": 120000000000,
    "http-timeout": 30000000000,
    "logs-cap": 10000,
    "max-open-connections": 0,
    "metrics-address": "127.0.0.1:6065",
    "txfee-cap": 1,
    "ws-address": "127.0.0.1:8546"
  },
  "min-retain-blocks": 0,
  "minimum-gas-prices": "4000000000000FX",
  "pruning": "default",
  "pruning-interval": "0",
  "pruning-keep-recent": "0",
  "rosetta": {
    "address": ":8080",
    "blockchain": "app",
    "denom-to-suggest": "FX",
    "enable": false,
    "enable-fee-suggestion": false,
    "gas-to-suggest": 200000,
    "network": "network",
    "offline": false,
    "retries": 3
  },
  "state-sync": {
    "snapshot-interval": 0,
    "snapshot-keep-recent": 2
  },
  "store": {
    "streamers": []
  },
  "streamers": {
    "file": {
      "fsync": false,
      "keys": [
        "*"
      ],
      "output-metadata": true,
      "prefix": "",
      "stop-node-on-error": true,
      "write_dir": ""
    }
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
}
`

	_, v := AppConfig(fxtypes.GetDefGasPrice())
	cfg := v.(Config)
	c := appTomlConfig{config: &cfg}
	buf := new(bytes.Buffer)
	clientCtx := client.Context{
		Output:       buf,
		OutputFormat: "json",
	}
	assert.NoError(t, c.output(clientCtx))
	assert.Equal(t, appConfigJson, buf.String())
}

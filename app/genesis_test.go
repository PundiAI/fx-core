package app_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v7/app"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

func TestNewDefaultGenesisByDenom(t *testing.T) {
	const genesisData = `{
  "arbitrum": {
	"params": {
	  "gravity_id": "fx-arbitrum-bridge",
	  "average_block_time": "7000",
	  "external_batch_timeout": "43200000",
	  "average_external_block_time": "500",
	  "signed_window": "30000",
	  "slash_fraction": "0.800000000000000000",
	  "oracle_set_update_power_change_percent": "0.100000000000000000",
	  "ibc_transfer_timeout_height": "20000",
	  "oracles": [],
	  "delegate_threshold": {
		"denom": "FX",
		"amount": "10000000000000000000000"
	  },
	  "delegate_multiple": "10",
	  "bridge_call_timeout": "604800000",
	  "bridge_call_max_gas_limit": "30000000"
	},
	"last_observed_event_nonce": "0",
	"last_observed_block_height": {
	  "external_block_height": "0",
	  "block_height": "0"
	},
	"oracles": [],
	"oracle_sets": [],
	"bridge_tokens": [],
	"unbatched_transfers": [],
	"batches": [],
	"oracle_set_confirms": [],
	"batch_confirms": [],
	"attestations": [],
	"proposal_oracle": {
	  "oracles": []
	},
	"last_observed_oracle_set": {
	  "nonce": "0",
	  "members": [],
	  "height": "0"
	},
	"last_slashed_batch_block": "0",
	"last_slashed_oracle_set_nonce": "0"
  },
  "auth": {
	"params": {
	  "max_memo_characters": "256",
	  "tx_sig_limit": "7",
	  "tx_size_cost_per_byte": "10",
	  "sig_verify_cost_ed25519": "590",
	  "sig_verify_cost_secp256k1": "1000"
	},
	"accounts": []
  },
  "authz": {
	"authorization": []
  },
  "avalanche": {
	"params": {
	  "gravity_id": "fx-avalanche-bridge",
	  "average_block_time": "7000",
	  "external_batch_timeout": "43200000",
	  "average_external_block_time": "2000",
	  "signed_window": "30000",
	  "slash_fraction": "0.800000000000000000",
	  "oracle_set_update_power_change_percent": "0.100000000000000000",
	  "ibc_transfer_timeout_height": "20000",
	  "oracles": [],
	  "delegate_threshold": {
		"denom": "FX",
		"amount": "10000000000000000000000"
	  },
	  "delegate_multiple": "10",
	  "bridge_call_timeout": "604800000",
	  "bridge_call_max_gas_limit": "30000000"
	},
	"last_observed_event_nonce": "0",
	"last_observed_block_height": {
	  "external_block_height": "0",
	  "block_height": "0"
	},
	"oracles": [],
	"oracle_sets": [],
	"bridge_tokens": [],
	"unbatched_transfers": [],
	"batches": [],
	"oracle_set_confirms": [],
	"batch_confirms": [],
	"attestations": [],
	"proposal_oracle": {
	  "oracles": []
	},
	"last_observed_oracle_set": {
	  "nonce": "0",
	  "members": [],
	  "height": "0"
	},
	"last_slashed_batch_block": "0",
	"last_slashed_oracle_set_nonce": "0"
  },
  "bank": {
	"params": {
	  "send_enabled": [],
	  "default_send_enabled": true
	},
	"balances": [
	  {
		"address": "cosmos1c602zv38ht8xu8u2qcmymyl55mcyvvjrzq9ur3",
		"coins": [
		  {
			"denom": "FX",
			"amount": "378600525462891000000000000"
		  }
		]
	  }
	],
	"supply": [
	  {
		"denom": "FX",
		"amount": "378604525462891000000000000"
	  }
	],
	"denom_metadata": [
	  {
		"description": "The native staking token of the Function X",
		"denom_units": [
		  {
			"denom": "FX",
			"exponent": 0,
			"aliases": []
		  }
		],
		"base": "FX",
		"display": "FX",
		"name": "Function X",
		"symbol": "FX",
		"uri": "",
		"uri_hash": ""
	  }
	]
  },
  "bsc": {
	"params": {
	  "gravity_id": "fx-bsc-bridge",
	  "average_block_time": "7000",
	  "external_batch_timeout": "43200000",
	  "average_external_block_time": "3000",
	  "signed_window": "30000",
	  "slash_fraction": "0.800000000000000000",
	  "oracle_set_update_power_change_percent": "0.100000000000000000",
	  "ibc_transfer_timeout_height": "20000",
	  "oracles": [],
	  "delegate_threshold": {
		"denom": "FX",
		"amount": "10000000000000000000000"
	  },
	  "delegate_multiple": "10",
	  "bridge_call_timeout": "604800000",
	  "bridge_call_max_gas_limit": "30000000"
	},
	"last_observed_event_nonce": "0",
	"last_observed_block_height": {
	  "external_block_height": "0",
	  "block_height": "0"
	},
	"oracles": [],
	"oracle_sets": [],
	"bridge_tokens": [],
	"unbatched_transfers": [],
	"batches": [],
	"oracle_set_confirms": [],
	"batch_confirms": [],
	"attestations": [],
	"proposal_oracle": {
	  "oracles": []
	},
	"last_observed_oracle_set": {
	  "nonce": "0",
	  "members": [],
	  "height": "0"
	},
	"last_slashed_batch_block": "0",
	"last_slashed_oracle_set_nonce": "0"
  },
  "capability": {
	"index": "1",
	"owners": []
  },
  "crisis": {
	"constant_fee": {
	  "denom": "FX",
	  "amount": "13333000000000000000000"
	}
  },
  "distribution": {
	"params": {
	  "community_tax": "0.400000000000000000",
	  "base_proposer_reward": "0.010000000000000000",
	  "bonus_proposer_reward": "0.040000000000000000",
	  "withdraw_addr_enabled": true
	},
	"fee_pool": {
	  "community_pool": []
	},
	"delegator_withdraw_infos": [],
	"previous_proposer": "",
	"outstanding_rewards": [],
	"validator_accumulated_commissions": [],
	"validator_historical_rewards": [],
	"validator_current_rewards": [],
	"delegator_starting_infos": [],
	"validator_slash_events": []
  },
  "erc20": {
	"params": {
	  "enable_erc20": true,
	  "enable_evm_hook": true,
	  "ibc_timeout": "43200s"
	},
	"token_pairs": []
  },
  "eth": {
	"params": {
	  "gravity_id": "fx-bridge-eth",
	  "average_block_time": "7000",
	  "external_batch_timeout": "43200000",
	  "average_external_block_time": "15000",
	  "signed_window": "30000",
	  "slash_fraction": "0.800000000000000000",
	  "oracle_set_update_power_change_percent": "0.100000000000000000",
	  "ibc_transfer_timeout_height": "20000",
	  "oracles": [],
	  "delegate_threshold": {
		"denom": "FX",
		"amount": "10000000000000000000000"
	  },
	  "delegate_multiple": "10",
	  "bridge_call_timeout": "604800000",
	  "bridge_call_max_gas_limit": "30000000"
	},
	"last_observed_event_nonce": "0",
	"last_observed_block_height": {
	  "external_block_height": "0",
	  "block_height": "0"
	},
	"oracles": [],
	"oracle_sets": [],
	"bridge_tokens": [],
	"unbatched_transfers": [],
	"batches": [],
	"oracle_set_confirms": [],
	"batch_confirms": [],
	"attestations": [],
	"proposal_oracle": {
	  "oracles": []
	},
	"last_observed_oracle_set": {
	  "nonce": "0",
	  "members": [],
	  "height": "0"
	},
	"last_slashed_batch_block": "0",
	"last_slashed_oracle_set_nonce": "0"
  },
  "evidence": {
	"evidence": []
  },
  "evm": {
	"accounts": [],
	"params": {
	  "evm_denom": "FX",
	  "enable_create": true,
	  "enable_call": true,
	  "extra_eips": [],
	  "chain_config": {
		"homestead_block": "0",
		"dao_fork_block": "0",
		"dao_fork_support": true,
		"eip150_block": "0",
		"eip150_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
		"eip155_block": "0",
		"eip158_block": "0",
		"byzantium_block": "0",
		"constantinople_block": "0",
		"petersburg_block": "0",
		"istanbul_block": "0",
		"muir_glacier_block": "0",
		"berlin_block": "0",
		"london_block": "0",
		"arrow_glacier_block": "0",
		"gray_glacier_block": "0",
		"merge_netsplit_block": "0",
		"shanghai_block": "0",
		"cancun_block": "0"
	  },
	  "allow_unprotected_txs": false
	}
  },
  "feegrant": {
	"allowances": []
  },
  "feemarket": {
	"params": {
	  "no_base_fee": false,
	  "base_fee_change_denominator": 8,
	  "elasticity_multiplier": 2,
	  "enable_height": "0",
	  "base_fee": "500000000000",
	  "min_gas_price": "500000000000.000000000000000000",
	  "min_gas_multiplier": "0.000000000000000000"
	},
	"block_gas": "0"
  },
  "fxtransfer": {},
  "genutil": {
	"gen_txs": []
  },
  "gov": {
	"starting_proposal_id": "1",
	"deposits": [],
	"votes": [],
	"proposals": [],
	"deposit_params": {
	  "min_deposit": [
		{
		  "denom": "FX",
		  "amount": "10000000000000000000000"
		}
	  ],
	  "max_deposit_period": "1209600s"
	},
	"voting_params": {
	  "voting_period": "1209600s"
	},
	"tally_params": {
	  "quorum": "0.400000000000000000",
	  "threshold": "0.500000000000000000",
	  "veto_threshold": "0.334000000000000000"
	}
  },
  "ibc": {
	"client_genesis": {
	  "clients": [],
	  "clients_consensus": [],
	  "clients_metadata": [],
	  "params": {
		"allowed_clients": [
		  "07-tendermint"
		]
	  },
	  "create_localhost": false,
	  "next_client_sequence": "0"
	},
	"connection_genesis": {
	  "connections": [],
	  "client_connection_paths": [],
	  "next_connection_sequence": "0",
	  "params": {
		"max_expected_time_per_block": "30000000000"
	  }
	},
	"channel_genesis": {
	  "channels": [],
	  "acknowledgements": [],
	  "commitments": [],
	  "receipts": [],
	  "send_sequences": [],
	  "recv_sequences": [],
	  "ack_sequences": [],
	  "next_channel_sequence": "0"
	}
  },
  "layer2": {
	"params": {
	  "gravity_id": "fx-layer2-bridge",
	  "average_block_time": "7000",
	  "external_batch_timeout": "43200000",
	  "average_external_block_time": "2000",
	  "signed_window": "30000",
	  "slash_fraction": "0.800000000000000000",
	  "oracle_set_update_power_change_percent": "0.100000000000000000",
	  "ibc_transfer_timeout_height": "20000",
	  "oracles": [],
	  "delegate_threshold": {
		"denom": "FX",
		"amount": "10000000000000000000000"
	  },
	  "delegate_multiple": "10",
	  "bridge_call_timeout": "604800000",
	  "bridge_call_max_gas_limit": "30000000"
	},
	"last_observed_event_nonce": "0",
	"last_observed_block_height": {
	  "external_block_height": "0",
	  "block_height": "0"
	},
	"oracles": [],
	"oracle_sets": [],
	"bridge_tokens": [],
	"unbatched_transfers": [],
	"batches": [],
	"oracle_set_confirms": [],
	"batch_confirms": [],
	"attestations": [],
	"proposal_oracle": {
	  "oracles": []
	},
	"last_observed_oracle_set": {
	  "nonce": "0",
	  "members": [],
	  "height": "0"
	},
	"last_slashed_batch_block": "0",
	"last_slashed_oracle_set_nonce": "0"
  },
  "migrate": {},
  "mint": {
	"minter": {
	  "inflation": "0.350000000000000000",
	  "annual_provisions": "0.000000000000000000"
	},
	"params": {
	  "mint_denom": "FX",
	  "inflation_rate_change": "0.300000000000000000",
	  "inflation_max": "0.416762000000000000",
	  "inflation_min": "0.170000000000000000",
	  "goal_bonded": "0.510000000000000000",
	  "blocks_per_year": "6311520"
	}
  },
  "optimism": {
	"params": {
	  "gravity_id": "fx-optimism-bridge",
	  "average_block_time": "7000",
	  "external_batch_timeout": "43200000",
	  "average_external_block_time": "500",
	  "signed_window": "30000",
	  "slash_fraction": "0.800000000000000000",
	  "oracle_set_update_power_change_percent": "0.100000000000000000",
	  "ibc_transfer_timeout_height": "20000",
	  "oracles": [],
	  "delegate_threshold": {
		"denom": "FX",
		"amount": "10000000000000000000000"
	  },
	  "delegate_multiple": "10",
	  "bridge_call_timeout": "604800000",
	  "bridge_call_max_gas_limit": "30000000"
	},
	"last_observed_event_nonce": "0",
	"last_observed_block_height": {
	  "external_block_height": "0",
	  "block_height": "0"
	},
	"oracles": [],
	"oracle_sets": [],
	"bridge_tokens": [],
	"unbatched_transfers": [],
	"batches": [],
	"oracle_set_confirms": [],
	"batch_confirms": [],
	"attestations": [],
	"proposal_oracle": {
	  "oracles": []
	},
	"last_observed_oracle_set": {
	  "nonce": "0",
	  "members": [],
	  "height": "0"
	},
	"last_slashed_batch_block": "0",
	"last_slashed_oracle_set_nonce": "0"
  },
  "params": {},
  "polygon": {
	"params": {
	  "gravity_id": "fx-polygon-bridge",
	  "average_block_time": "7000",
	  "external_batch_timeout": "43200000",
	  "average_external_block_time": "2000",
	  "signed_window": "30000",
	  "slash_fraction": "0.800000000000000000",
	  "oracle_set_update_power_change_percent": "0.100000000000000000",
	  "ibc_transfer_timeout_height": "20000",
	  "oracles": [],
	  "delegate_threshold": {
		"denom": "FX",
		"amount": "10000000000000000000000"
	  },
	  "delegate_multiple": "10",
	  "bridge_call_timeout": "604800000",
	  "bridge_call_max_gas_limit": "30000000"
	},
	"last_observed_event_nonce": "0",
	"last_observed_block_height": {
	  "external_block_height": "0",
	  "block_height": "0"
	},
	"oracles": [],
	"oracle_sets": [],
	"bridge_tokens": [],
	"unbatched_transfers": [],
	"batches": [],
	"oracle_set_confirms": [],
	"batch_confirms": [],
	"attestations": [],
	"proposal_oracle": {
	  "oracles": []
	},
	"last_observed_oracle_set": {
	  "nonce": "0",
	  "members": [],
	  "height": "0"
	},
	"last_slashed_batch_block": "0",
	"last_slashed_oracle_set_nonce": "0"
  },
  "slashing": {
	"params": {
	  "signed_blocks_window": "20000",
	  "min_signed_per_window": "0.050000000000000000",
	  "downtime_jail_duration": "600s",
	  "slash_fraction_double_sign": "0.050000000000000000",
	  "slash_fraction_downtime": "0.001000000000000000"
	},
	"signing_infos": [],
	"missed_blocks": []
  },
  "staking": {
	"params": {
	  "unbonding_time": "1814400s",
	  "max_validators": 20,
	  "max_entries": 7,
	  "historical_entries": 20000,
	  "bond_denom": "FX",
	  "min_commission_rate": "0.000000000000000000"
	},
	"last_total_power": "0",
	"last_validator_powers": [],
	"validators": [],
	"delegations": [],
	"unbonding_delegations": [],
	"redelegations": [],
	"exported": false,
	"allowances": []
  },
  "transfer": {
	"port_id": "transfer",
	"denom_traces": [],
	"params": {
	  "send_enabled": true,
	  "receive_enabled": true
	}
  },
  "tron": {
	"params": {
	  "gravity_id": "fx-tron-bridge",
	  "average_block_time": "7000",
	  "external_batch_timeout": "43200000",
	  "average_external_block_time": "3000",
	  "signed_window": "30000",
	  "slash_fraction": "0.800000000000000000",
	  "oracle_set_update_power_change_percent": "0.100000000000000000",
	  "ibc_transfer_timeout_height": "20000",
	  "oracles": [],
	  "delegate_threshold": {
		"denom": "FX",
		"amount": "10000000000000000000000"
	  },
	  "delegate_multiple": "10",
	  "bridge_call_timeout": "604800000",
	  "bridge_call_max_gas_limit": "30000000"
	},
	"last_observed_event_nonce": "0",
	"last_observed_block_height": {
	  "external_block_height": "0",
	  "block_height": "0"
	},
	"oracles": [],
	"oracle_sets": [],
	"bridge_tokens": [],
	"unbatched_transfers": [],
	"batches": [],
	"oracle_set_confirms": [],
	"batch_confirms": [],
	"attestations": [],
	"proposal_oracle": {
	  "oracles": []
	},
	"last_observed_oracle_set": {
	  "nonce": "0",
	  "members": [],
	  "height": "0"
	},
	"last_slashed_batch_block": "0",
	"last_slashed_oracle_set_nonce": "0"
  },
  "upgrade": {},
  "vesting": {}
}`

	encodingConfig := app.MakeEncodingConfig()
	genAppState := app.NewDefAppGenesisByDenom(fxtypes.DefaultDenom, encodingConfig.Codec)

	genAppStateStr, err := json.MarshalIndent(genAppState, "", "  ")
	assert.NoError(t, err)

	assert.JSONEqf(t, genesisData, string(genAppStateStr), "")
}

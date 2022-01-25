package cmd

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

func QueryBlockResultsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block-results <height>",
		Short: "Query for a transaction by hash in a committed block",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			var height int64
			if len(args) > 0 {
				height, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					blockHeight, err := hexutil.DecodeUint64(args[0])
					if err != nil {
						return err
					}
					height = int64(blockHeight)
				}
			} else {
				status, err := clientCtx.Client.Status(context.Background())
				if err != nil {
					return err
				}
				height = status.SyncInfo.LatestBlockHeight
			}

			blockResults, err := clientCtx.Client.BlockResults(context.Background(), &height)
			if err != nil {
				return err
			}
			output, err := ParseBlockResults(clientCtx.Codec, blockResults)
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput([]byte(output))
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func ParseBlockResults(cdc codec.JSONCodec, blockResults *coretypes.ResultBlockResults) (string, error) {
	consensusParamUpdates, err := cdc.MarshalJSON(blockResults.ConsensusParamUpdates)
	if err != nil {
		return "", err
	}
	var beginBlockEvents []map[string]interface{}
	for _, event := range blockResults.BeginBlockEvents {
		var attributes []map[string]interface{}
		for _, attribute := range event.Attributes {
			attributes = append(attributes, map[string]interface{}{
				"index": attribute.Index,
				"key":   string(attribute.Key),
				"value": string(attribute.Value),
			})
		}
		beginBlockEvents = append(beginBlockEvents, map[string]interface{}{
			"type":       event.Type,
			"attributes": attributes,
		})
	}
	var endBlockEvents []map[string]interface{}
	for _, event := range blockResults.EndBlockEvents {
		var attributes []map[string]interface{}
		for _, attribute := range event.Attributes {
			attributes = append(attributes, map[string]interface{}{
				"index": attribute.Index,
				"key":   string(attribute.Key),
				"value": string(attribute.Value),
			})
		}
		endBlockEvents = append(endBlockEvents, map[string]interface{}{
			"type":       event.Type,
			"attributes": attributes,
		})
	}
	var txsResults []map[string]interface{}
	for _, txResult := range blockResults.TxsResults {
		var txResultEvents []map[string]interface{}
		for _, event := range txResult.Events {
			var attributes []map[string]interface{}
			for _, attribute := range event.Attributes {
				attributes = append(attributes, map[string]interface{}{
					"index": attribute.Index,
					"key":   string(attribute.Key),
					"value": string(attribute.Value),
				})
			}
			txResultEvents = append(txResultEvents, map[string]interface{}{
				"type":       event.Type,
				"attributes": attributes,
			})
		}
		txsResults = append(txsResults, map[string]interface{}{
			"code":       txResult.Code,
			"data":       string(txResult.Data),
			"log":        txResult.Log,
			"info":       txResult.Info,
			"gas_wanted": txResult.GasWanted,
			"gas_used":   txResult.GasUsed,
			"events":     txResultEvents,
			"codespace":  txResult.Codespace,
		})
	}
	var validatorUpdates []json.RawMessage
	for _, valUp := range blockResults.ValidatorUpdates {
		valUpData, err := cdc.MarshalJSON(&valUp)
		if err != nil {
			return "", err
		}
		validatorUpdates = append(validatorUpdates, valUpData)
	}
	output, err := json.Marshal(map[string]interface{}{
		"height":                  blockResults.Height,
		"txs_results":             txsResults,
		"begin_block_events":      beginBlockEvents,
		"end_block_events":        endBlockEvents,
		"validator_updates":       validatorUpdates,
		"consensus_param_updates": json.RawMessage(consensusParamUpdates),
	})
	if err != nil {
		return "", err
	}
	return string(output), nil
}

package cli

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/tendermint/tendermint/abci/types"

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
			return PrintOutput(clientCtx, output)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func ParseBlockResults(cdc codec.JSONCodec, blockResults *coretypes.ResultBlockResults) (interface{}, error) {
	consensusParamUpdates, err := cdc.MarshalJSON(blockResults.ConsensusParamUpdates)
	if err != nil {
		return nil, err
	}
	var beginBlockEvents []map[string]interface{}
	for _, event := range blockResults.BeginBlockEvents {
		beginBlockEvents = append(beginBlockEvents, map[string]interface{}{
			"type":       event.Type,
			"attributes": AttributesToMap(event.Attributes),
		})
	}
	var endBlockEvents []map[string]interface{}
	for _, event := range blockResults.EndBlockEvents {
		endBlockEvents = append(endBlockEvents, map[string]interface{}{
			"type":       event.Type,
			"attributes": AttributesToMap(event.Attributes),
		})
	}
	var txsResults []map[string]interface{}
	for _, txResult := range blockResults.TxsResults {
		txsResults = append(txsResults, TxResultToMap(txResult))
	}
	var validatorUpdates []json.RawMessage
	for _, valUp := range blockResults.ValidatorUpdates {
		valUpData, err := cdc.MarshalJSON(&valUp)
		if err != nil {
			return nil, err
		}
		validatorUpdates = append(validatorUpdates, valUpData)
	}
	return map[string]interface{}{
		"height":                  blockResults.Height,
		"txs_results":             txsResults,
		"begin_block_events":      beginBlockEvents,
		"end_block_events":        endBlockEvents,
		"validator_updates":       validatorUpdates,
		"consensus_param_updates": json.RawMessage(consensusParamUpdates),
	}, nil
}

func TxResponseToMap(cdc codec.JSONCodec, txResponse *sdk.TxResponse) map[string]interface{} {
	if txResponse == nil {
		return map[string]interface{}{}
	}
	var txResultEvents []map[string]interface{}
	for _, event := range txResponse.Events {
		txResultEvents = append(txResultEvents, map[string]interface{}{
			"type":       event.Type,
			"attributes": AttributesToMap(event.Attributes),
		})
	}
	tx, err := cdc.MarshalJSON(txResponse.Tx)
	if err != nil {
		return nil
	}
	txData, err := hex.DecodeString(txResponse.Data)
	if err != nil {
		return nil
	}
	var txMsgData = sdk.TxMsgData{
		Data: make([]*sdk.MsgData, 0),
	}
	if err := proto.Unmarshal(txData, &txMsgData); err != nil {
		return nil
	}
	return map[string]interface{}{
		"height":     txResponse.Height,
		"txhash":     txResponse.TxHash,
		"codespace":  txResponse.Codespace,
		"code":       txResponse.Code,
		"data":       txMsgData,
		"raw_log":    txResponse.RawLog,
		"logs":       txResponse.Logs,
		"info":       txResponse.Info,
		"gas_wanted": txResponse.GasWanted,
		"gas_used":   txResponse.GasUsed,
		"tx":         json.RawMessage(tx),
		"timestamp":  txResponse.Timestamp,
		"events":     txResultEvents,
	}
}

func TxResultToMap(txResult *types.ResponseDeliverTx) map[string]interface{} {
	if txResult == nil {
		return map[string]interface{}{}
	}
	var txResultEvents []map[string]interface{}
	for _, event := range txResult.Events {
		txResultEvents = append(txResultEvents, map[string]interface{}{
			"type":       event.Type,
			"attributes": AttributesToMap(event.Attributes),
		})
	}
	var txMsgData = sdk.TxMsgData{
		Data: make([]*sdk.MsgData, 0),
	}
	if err := proto.Unmarshal(txResult.Data, &txMsgData); err != nil {
		return nil
	}
	return map[string]interface{}{
		"code":       txResult.Code,
		"data":       txMsgData,
		"log":        txResult.Log,
		"info":       txResult.Info,
		"gas_wanted": txResult.GasWanted,
		"gas_used":   txResult.GasUsed,
		"events":     txResultEvents,
		"codespace":  txResult.Codespace,
	}
}

func AttributesToMap(attrs []types.EventAttribute) []map[string]interface{} {
	var attributes []map[string]interface{}
	for _, attribute := range attrs {
		attributes = append(attributes, map[string]interface{}{
			"index": attribute.Index,
			"key":   string(attribute.Key),
			"value": string(attribute.Value),
		})
	}
	return attributes
}

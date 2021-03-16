package cmd

import (
	"context"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
)

func QueryBlockResultsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block-results <height>",
		Short: "Query for a transaction by hash in a committed block",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			height, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			blockResults, err := clientCtx.Client.BlockResults(context.Background(), &height)
			if err != nil {
				return err
			}
			consensusParamUpdates, err := clientCtx.JSONMarshaler.MarshalJSON(blockResults.ConsensusParamUpdates)
			if err != nil {
				return err
			}
			var beginBlockEvents []map[string]interface{}
			for _, event := range blockResults.BeginBlockEvents {
				var attributes []map[string]interface{}
				for _, attribute := range event.Attributes {
					attributes = append(attributes, map[string]interface{}{
						"Index": attribute.Index,
						"Key":   string(attribute.Key),
						"Value": string(attribute.Value),
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
						"Index": attribute.Index,
						"Key":   string(attribute.Key),
						"Value": string(attribute.Value),
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
							"Index": attribute.Index,
							"Key":   string(attribute.Key),
							"Value": string(attribute.Value),
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
				valUpData, err := clientCtx.JSONMarshaler.MarshalJSON(&valUp)
				if err != nil {
					return err
				}
				validatorUpdates = append(validatorUpdates, valUpData)
			}
			output, err := json.Marshal(map[string]interface{}{
				"Height":                blockResults.Height,
				"ConsensusParamUpdates": json.RawMessage(consensusParamUpdates),
				"BeginBlockEvents":      beginBlockEvents,
				"EndBlockEvents":        endBlockEvents,
				"TxsResults":            txsResults,
				"ValidatorUpdates":      validatorUpdates,
			})
			if err != nil {
				return err
			}
			return PrintOutput(clientCtx, output)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func PrintOutput(ctx client.Context, out []byte) error {
	if ctx.OutputFormat == "text" {
		// handle text format by decoding and re-encoding JSON as YAML
		var j interface{}

		err := json.Unmarshal(out, &j)
		if err != nil {
			return err
		}

		out, err = yaml.Marshal(j)
		if err != nil {
			return err
		}
	}

	writer := ctx.Output
	if writer == nil {
		writer = os.Stdout
	}

	_, err := writer.Write(out)
	if err != nil {
		return err
	}

	if ctx.OutputFormat != "text" {
		// append new-line for formats besides YAML
		_, err = writer.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/version"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/spf13/cobra"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	flagEvents = "events"
	flagType   = "type"

	typeHash   = "hash"
	typeAccSeq = "acc_seq"
	typeSig    = "signature"

	eventFormat = "{eventType}.{eventAttribute}={value}"
)

// QueryTxsByEventsCmd returns a command to search through transactions by events.
func QueryTxsByEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "txs",
		Short: "Query for paginated transactions that match a set of events",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Search for transactions that match the exact given events where results are paginated.
Each event takes the form of '%s'. Please refer
to each module's documentation for the full set of events to query for. Each module
documents its respective events under 'xx_events.md'.

Example:
$ %s query txs --%s 'message.sender=cosmos1...&message.action=withdraw_delegator_reward' --page 1 --limit 30
$ %s query txs --%s 'message.sender=cosmos1...&message.module=distribution' --page 1 --limit 30
`, eventFormat, version.AppName, flagEvents, version.AppName, flagEvents),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			eventsRaw, _ := cmd.Flags().GetString(flagEvents)
			eventsStr := strings.Trim(eventsRaw, "'")

			var events []string
			if strings.Contains(eventsStr, "&") {
				events = strings.Split(eventsStr, "&")
			} else {
				events = append(events, eventsStr)
			}

			var tmEvents []string

			for _, event := range events {
				if !strings.Contains(event, "=") {
					return fmt.Errorf("invalid event; event %s should be of the format: %s", event, eventFormat)
				} else if strings.Count(event, "=") > 1 {
					return fmt.Errorf("invalid event; event %s should be of the format: %s", event, eventFormat)
				}

				tokens := strings.Split(event, "=")
				if tokens[0] == tmtypes.TxHeightKey {
					event = fmt.Sprintf("%s=%s", tokens[0], tokens[1])
				} else {
					event = fmt.Sprintf("%s='%s'", tokens[0], tokens[1])
				}

				tmEvents = append(tmEvents, event)
			}

			page, _ := cmd.Flags().GetInt(flags.FlagPage)
			limit, _ := cmd.Flags().GetInt(flags.FlagLimit)

			txsResult, err := authtx.QueryTxsByEvents(clientCtx, tmEvents, page, limit, "")
			if err != nil {
				return err
			}
			var txsArray = make([]interface{}, 0)
			for i := 0; i < len(txsResult.Txs); i++ {
				txsArray = append(txsArray, TxResponseToMap(clientCtx.Codec, txsResult.Txs[i]))
			}
			return PrintOutput(clientCtx, map[string]interface{}{
				"total_count": txsResult.TotalCount,
				"count":       txsResult.Count,
				"page_number": txsResult.PageNumber,
				"page_total":  txsResult.PageTotal,
				"limit":       txsResult.Limit,
				"txs":         txsArray,
			})
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().Int(flags.FlagPage, rest.DefaultPage, "Query a specific page of paginated results")
	cmd.Flags().Int(flags.FlagLimit, rest.DefaultLimit, "Query number of transactions results per page returned")
	cmd.Flags().String(flagEvents, "", fmt.Sprintf("list of transaction events in the form of %s", eventFormat))
	_ = cmd.MarkFlagRequired(flagEvents)

	return cmd
}

// QueryTxCmd implements the default command for a tx query.
func QueryTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx --type=[hash|acc_seq|signature] [hash|acc_seq|signature]",
		Short: "Query for a transaction by hash, \"<addr>/<seq>\" combination or comma-separated signatures in a committed block",
		Long: strings.TrimSpace(fmt.Sprintf(`
Example:
$ %s query tx <hash>
$ %s query tx --%s=%s <addr>/<sequence>
$ %s query tx --%s=%s <sig1_base64>,<sig2_base64...>
`,
			version.AppName,
			version.AppName, flagType, typeAccSeq,
			version.AppName, flagType, typeSig)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			typ, _ := cmd.Flags().GetString(flagType)

			switch typ {
			case typeHash:
				{
					txHash := args[0]
					if txHash == "" {
						return fmt.Errorf("argument should be a tx hash")
					}
					txHash = strings.TrimPrefix(txHash, "0x")

					// If hash is given, then query the tx by hash.
					output, err := authtx.QueryTx(clientCtx, txHash)
					if err != nil {
						return err
					}

					if output.Empty() {
						return fmt.Errorf("no transaction found with hash %s", txHash)
					}
					return PrintOutput(clientCtx, TxResponseToMap(clientCtx.Codec, output))
				}
			case typeSig:
				{
					sigParts, err := parseSigArgs(args)
					if err != nil {
						return err
					}
					tmEvents := make([]string, len(sigParts))
					for i, sig := range sigParts {
						tmEvents[i] = fmt.Sprintf("%s.%s='%s'", sdk.EventTypeTx, sdk.AttributeKeySignature, sig)
					}

					txs, err := authtx.QueryTxsByEvents(clientCtx, tmEvents, rest.DefaultPage, query.DefaultLimit, "")
					if err != nil {
						return err
					}
					if len(txs.Txs) == 0 {
						return fmt.Errorf("found no txs matching given signatures")
					}
					if len(txs.Txs) > 1 {
						// This case means there's a bug somewhere else in the code. Should not happen.
						return errors.ErrLogic.Wrapf("found %d txs matching given signatures", len(txs.Txs))
					}
					return PrintOutput(clientCtx, TxResponseToMap(clientCtx.Codec, txs.Txs[0]))
				}
			case typeAccSeq:
				{
					if args[0] == "" {
						return fmt.Errorf("`acc_seq` type takes an argument '<addr>/<seq>'")
					}

					tmEvents := []string{
						fmt.Sprintf("%s.%s='%s'", sdk.EventTypeTx, sdk.AttributeKeyAccountSequence, args[0]),
					}
					txs, err := authtx.QueryTxsByEvents(clientCtx, tmEvents, rest.DefaultPage, query.DefaultLimit, "")
					if err != nil {
						return err
					}
					if len(txs.Txs) == 0 {
						return fmt.Errorf("found no txs matching given address and sequence combination")
					}
					if len(txs.Txs) > 1 {
						// This case means there's a bug somewhere else in the code. Should not happen.
						return fmt.Errorf("found %d txs matching given address and sequence combination", len(txs.Txs))
					}

					return PrintOutput(clientCtx, TxResponseToMap(clientCtx.Codec, txs.Txs[0]))
				}
			default:
				return fmt.Errorf("unknown --%s value %s", flagType, typ)
			}
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().String(flagType, typeHash, fmt.Sprintf("The type to be used when querying tx, can be one of \"%s\", \"%s\", \"%s\"", typeHash, typeAccSeq, typeSig))

	return cmd
}

// parseSigArgs parses comma-separated signatures from the CLI arguments.
func parseSigArgs(args []string) ([]string, error) {
	if len(args) != 1 || args[0] == "" {
		return nil, fmt.Errorf("argument should be comma-separated signatures")
	}

	return strings.Split(args[0], ","), nil
}

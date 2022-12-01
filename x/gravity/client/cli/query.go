package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	abcitype "github.com/tendermint/tendermint/abci/types"

	"github.com/functionx/fx-core/v3/x/gravity/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the eth module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
		Hidden:                     true,
	}

	cmd.AddCommand([]*cobra.Command{
		// query module params
		CmdGetParams(),

		// query delegate relation
		CmdGetDelegateKeyByValidator(),
		CmdGetDelegateKeyByOrchestrator(),
		CmdGetDelegateKeyByEth(),

		// query Validator set
		CmdGetCurrentValset(),
		CmdGetValsetRequest(),

		// need validator consensus sign
		// 1. validator set change confirm
		CmdGetLastValSetRequests(),
		CmdGetPendingValsetRequest(),
		CmdGetValsetConfirm(),
		CmdGetValsetConfirms(),

		// 2. request batch confirm
		CmdGetPendingOutgoingTXBatchRequest(),
		CmdBatchConfirm(),
		CmdBatchConfirms(),

		// send to eth
		CmdBatchRequestByNonce(),
		CmdGetPendingSendToEth(),
		CmdOutgoingTxBatches(),
		CmdGetBatchFees(),

		CmdGetLastObservedBlockHeight(),
		CmdProjectedBatchTimeoutHeight(),

		// denom <-> bep20
		CmdGetDenomToERC20Token(),
		CmdGetERC20TokenToDenom(),
		CmdGetBridgeTokens(),

		// validator event nonce
		// 1. fxcore validator event nonce
		CmdGetValidatorEventNonce(),
		// 2. eth event nonce block height
		CmdGetValidatorEventBlockHeight(),
	}...)

	for _, command := range cmd.Commands() {
		flags.AddQueryFlagsToCmd(command)
	}
	return cmd
}

func CmdGetParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Query the current gravity parameters information",
		Example: "fxcored q gravity params",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// nolint
			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetDelegateKeyByValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegate-key-by-validator [validator]",
		Short:   "Query delegate eth and fx key for a given validator",
		Example: "fxcored q gravity delegate-key-by-validator fxvaloper1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8wpk9j2",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			validator, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			req := &types.QueryDelegateKeyByValidatorRequest{
				ValidatorAddress: validator.String(),
			}

			// nolint
			res, err := queryClient.GetDelegateKeyByValidator(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetDelegateKeyByOrchestrator() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegate-key-by-orchestrator [orchestrator]",
		Short:   "Query delegate eth and validator key for a given orchestrator",
		Example: "fxcored q gravity delegate-key-by-orchestrator fx1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8xst6v2",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			validator, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// nolint
			res, err := queryClient.GetDelegateKeyByOrchestrator(cmd.Context(), &types.QueryDelegateKeyByOrchestratorRequest{
				OrchestratorAddress: validator.String(),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetDelegateKeyByEth() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegate-key-by-eth",
		Short:   "Query delegate key by eth address",
		Example: "fxcored q gravity delegate-key-by-eth 0xb86d4DC8e2C57190c1cfb834fE69A7a65E2756C2",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			if !gethCommon.IsHexAddress(args[0]) {
				return fmt.Errorf("target address is invalid!address:[%s]", args[0])
			}
			if !gethCommon.IsHexAddress(args[0]) {
				return fmt.Errorf("contract address is invalid!address:[%s]", args[0])
			}
			// nolint
			res, err := queryClient.GetDelegateKeyByEth(cmd.Context(), &types.QueryDelegateKeyByEthRequest{
				EthAddress: args[0],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetCurrentValset() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "current-valset",
		Short:   "Query current valset",
		Example: "fxcored q gravity current-valset",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// nolint
			res, err := queryClient.CurrentValset(cmd.Context(), &types.QueryCurrentValsetRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetValsetRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "valset-request [nonce]",
		Short:   "Query requested valset with a particular nonce",
		Example: "fxcored q gravity valset-request 1",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			var nonce uint64
			if len(args) <= 0 {
				queryAbciResp, err := clientCtx.QueryABCI(abcitype.RequestQuery{
					Path: fmt.Sprintf("store/%s/key", types.ModuleName),
					Data: types.LatestValsetNonce,
				})
				if err != nil {
					return err
				}
				nonce = sdk.BigEndianToUint64(queryAbciResp.Value)
			} else {
				if len(args) < 1 {
					return fmt.Errorf("require particular nonce")
				}
				nonce, err = strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					return err
				}
			}
			// nolint
			res, err := queryClient.ValsetRequest(cmd.Context(), &types.QueryValsetRequestRequest{
				Nonce: nonce,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetLastValSetRequests() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-valset-requests",
		Short: "Query last valset requests",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// nolint
			res, err := queryClient.LastValsetRequests(cmd.Context(), &types.QueryLastValsetRequestsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

func CmdGetPendingValsetRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pending-valset-request [orchestrator]",
		Short:   "Query the latest valset request which has not been signed by a particular validator orchestrator",
		Example: "fxcored q gravity pending-valset-request fx1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8xst6v2",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			orchestrator, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			req := &types.QueryLastPendingValsetRequestByAddrRequest{
				Address: orchestrator.String(),
			}

			// nolint
			res, err := queryClient.LastPendingValsetRequestByAddr(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetValsetConfirm() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "valset-confirm [nonce] [orchestrator]",
		Short:   "Query valset confirmation with a particular nonce from a particular validator orchestrator",
		Example: "fxcored q gravity valset-confirm 1 fx1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8xst6v2",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			// nolint
			res, err := queryClient.ValsetConfirm(cmd.Context(), &types.QueryValsetConfirmRequest{
				Nonce:   nonce,
				Address: args[1],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetValsetConfirms() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "valset-confirms [nonce]",
		Short:   "Query valset confirmations with a particular nonce",
		Example: "fxcored q gravity valset-confirms 1",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			// nolint
			res, err := queryClient.ValsetConfirmsByNonce(cmd.Context(), &types.QueryValsetConfirmsByNonceRequest{
				Nonce: nonce,
			})
			if err != nil {
				return err
			}
			if err = clientCtx.PrintString(fmt.Sprintf("confirm size:[%d]\n", len(res.Confirms))); err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetPendingOutgoingTXBatchRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pending-batch-request [orchestrator]",
		Short:   "Query the latest outgoing TX batch request which has not been signed by a particular validator orchestrator",
		Example: "fxcored q gravity pending-batch-request fx1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8xst6v2",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// nolint
			res, err := queryClient.LastPendingBatchRequestByAddr(cmd.Context(), &types.QueryLastPendingBatchRequestByAddrRequest{
				Address: args[0],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdBatchConfirm() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "batch-confirm  [token-contract] [nonce] [orchestrator]",
		Short:   "Query outgoing tx batches confirm by validator orchestrator",
		Example: "fxcored q gravity batch-confirm 0x30dA8589BFa1E509A319489E014d384b87815D89 1 fx1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8xst6v2",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			if !gethCommon.IsHexAddress(args[0]) {
				return fmt.Errorf("contract address is invalid!address:[%s]", args[0])
			}
			nonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			orchestrator, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}
			// nolint
			res, err := queryClient.BatchConfirm(cmd.Context(), &types.QueryBatchConfirmRequest{
				TokenContract: args[0],
				Nonce:         nonce,
				Address:       orchestrator.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdBatchConfirms() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "batch-confirms [contract] [nonce]",
		Short:   "Query outgoing tx batches confirms",
		Example: "fxcored q gravity batch-confirms 0x30dA8589BFa1E509A319489E014d384b87815D89 1",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			if !gethCommon.IsHexAddress(args[0]) {
				return fmt.Errorf("contract address is invalid!address:[%s]", args[0])
			}
			// nolint
			res, err := queryClient.BatchConfirms(cmd.Context(), &types.QueryBatchConfirmsRequest{
				TokenContract: args[0],
				Nonce:         uint64(nonce),
			})
			if err != nil {
				return err
			}
			if err = clientCtx.PrintString(fmt.Sprintf("confirm size:[%d]\n", len(res.Confirms))); err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdBatchRequestByNonce() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-request [token-contract] [nonce]",
		Short: "Query outgoing tx batches",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			// nolint
			res, err := queryClient.BatchRequestByNonce(cmd.Context(), &types.QueryBatchRequestByNonceRequest{
				TokenContract: args[1],
				Nonce:         nonce,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetPendingSendToEth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-send-to-eth [address]",
		Short: "query pending send to external txs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			if _, err := sdk.AccAddressFromBech32(args[0]); err != nil {
				return nil
			}
			// nolint
			res, err := queryClient.GetPendingSendToEth(cmd.Context(), &types.QueryPendingSendToEthRequest{
				SenderAddress: args[0],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdOutgoingTxBatches() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "outgoing-tx-batches",
		Short:   "Query outgoing tx batches",
		Example: "fxcored q gravity outgoing-tx-batches",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// nolint
			res, err := queryClient.OutgoingTxBatches(cmd.Context(), &types.QueryOutgoingTxBatchesRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetBatchFees() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "batch-fees",
		Short:   "Query a list of withdrawal transaction fees to be processed",
		Example: "fxcored q gravity batch-fees",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// nolint
			res, err := queryClient.BatchFees(cmd.Context(), &types.QueryBatchFeeRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetLastObservedBlockHeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-observed-block-height",
		Short: "Query last observed block height",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// nolint
			res, err := queryClient.LastObservedBlockHeight(cmd.Context(), &types.QueryLastObservedBlockHeightRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdProjectedBatchTimeoutHeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "projected-batch-timeout-height",
		Short: "Query projected batch timeout height",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// nolint
			res, err := queryClient.ProjectedBatchTimeoutHeight(cmd.Context(), &types.QueryProjectedBatchTimeoutHeightRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetDenomToERC20Token() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "erc20 [denom]",
		Short:   "Query contract address from denom",
		Example: "fxcored q gravity erc20 eth0x2170ed0880ac9a755fd29b2688956bd959f933f8",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// nolint
			res, err := queryClient.DenomToERC20(cmd.Context(), &types.QueryDenomToERC20Request{
				Denom: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetERC20TokenToDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "denom [token-contract]",
		Short:   "Query denom from contract address",
		Example: "fxcored q gravity denom 0x2170ed0880ac9a755fd29b2688956bd959f933f8",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			contractAddress := args[0]
			if !gethCommon.IsHexAddress(contractAddress) {
				return fmt.Errorf("invalid contract address:%s", contractAddress)
			}

			// nolint
			res, err := queryClient.ERC20ToDenom(cmd.Context(), &types.QueryERC20ToDenomRequest{
				Erc20: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetBridgeTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridge-tokens [chain-name]",
		Short: "Query bridge token list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// nolint
			res, err := queryClient.BridgeTokens(cmd.Context(), &types.QueryBridgeTokensRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

func CmdGetValidatorEventNonce() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "event-nonce [orchestrator]",
		Short:   "Query last event nonce by validator orchestratorAddress",
		Example: "fxcored q gravity event-nonce fx1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8xst6v2",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			orchestratorAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			// nolint
			res, err := queryClient.LastEventNonceByAddr(cmd.Context(), &types.QueryLastEventNonceByAddrRequest{
				Address: orchestratorAddress.String(),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetValidatorEventBlockHeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "event-block-height [orchestrator]",
		Short:   "Query last event block height by validator orchestrator",
		Example: "fxcored q gravity event-block-height fx1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8xst6v2",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			orchestratorAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			// nolint
			res, err := queryClient.LastEventBlockHeightByAddr(cmd.Context(), &types.QueryLastEventBlockHeightByAddrRequest{
				Address: orchestratorAddress.String(),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

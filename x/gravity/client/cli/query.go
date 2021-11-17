package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	abcitype "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/x/gravity/types"
)

const (
	flagValsetLatest = "latest"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the eth module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
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
		CmdGetPendingValsetRequest(),
		CmdGetValsetConfirm(),
		CmdGetValsetConfirms(),

		// 2. request batch confirm
		CmdGetPendingOutgoingTXBatchRequest(),
		CmdBatchConfirm(),
		CmdBatchConfirms(),

		// send to eth
		CmdOutgoingTxBatches(),
		CmdGetBatchFees(),

		// denom <-> bep20
		CmdGetDenomToERC20Token(),
		CmdGetERC20TokenToDenom(),

		// validator event nonce
		// 1. fxcore validator event nonce
		CmdGetValidatorEventNonce(),
		// 2. eth event nonce block height
		CmdGetValidatorEventBlockHeight(),

		// eth -> fxcore -> ibc
		// 1. query eth -> fxcore -> ibc transfer sequence block height
		CmdIbcSequenceHeight(),
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
		Short:   "Get delegate eth and fx key for a given validator",
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
		Short:   "Get delegate eth and validator key for a given orchestrator",
		Example: "fxcored q gravity delegate-key-by-orchestrator fx1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8xst6v2",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			validator, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

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
		Short:   "query delegate key by eth address",
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

			req := &types.QueryCurrentValsetRequest{}

			res, err := queryClient.CurrentValset(cmd.Context(), req)
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
		Short:   "Get requested valset with a particular nonce",
		Example: fmt.Sprintf("fxcored q gravity valset-request 1"),
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
				nonce = types.UInt64FromBytes(queryAbciResp.Value)
			} else {
				if len(args) < 1 {
					return fmt.Errorf("require particular nonce")
				}
				nonce, err = strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					return err
				}
			}
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

func CmdGetPendingValsetRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pending-valset-request [orchestrator]",
		Short:   "Get the latest valset request which has not been signed by a particular validator orchestrator",
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
		Short:   "Get valset confirmation with a particular nonce from a particular validator orchestrator",
		Example: "fxcored q gravity valset-confirm 1 fx1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8xst6v2",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
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
		Short:   "Get valset confirmations with a particular nonce",
		Example: "fxcored q gravity valset-confirms 1",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
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
		Short:   "Get the latest outgoing TX batch request which has not been signed by a particular validator orchestrator",
		Example: "fxcored q gravity pending-batch-request fx1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8xst6v2",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

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
		Use:     "batch-confirm  [contract] [nonce] [orchestrator]",
		Short:   "query outgoing tx batches confirm by validator orchestrator",
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
			res, err := queryClient.BatchConfirm(cmd.Context(), &types.QueryBatchConfirmRequest{
				ContractAddress: args[0],
				Nonce:           nonce,
				Address:         orchestrator.String(),
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
		Short:   "query outgoing tx batches confirms",
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
			res, err := queryClient.BatchConfirms(cmd.Context(), &types.QueryBatchConfirmsRequest{
				ContractAddress: args[0],
				Nonce:           uint64(nonce),
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

func CmdOutgoingTxBatches() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "outgoing-tx-batches",
		Short:   "query outgoing tx batches",
		Example: "fxcored q gravity outgoing-tx-batches",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

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
		Short:   "Gets a list of withdrawal transaction fees to be processed",
		Example: "fxcored q gravity batch-fees",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.BatchFees(cmd.Context(), &types.QueryBatchFeeRequest{})
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
		Use:     "denom [contractAddress]",
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

func CmdGetValidatorEventNonce() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "event-nonce [orchestrator]",
		Short:   "query last event nonce by validator orchestratorAddress",
		Example: "fxcored q gravity event-nonce fx1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8xst6v2",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			orchestratorAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
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
		Short:   "query last event block height by validator orchestrator",
		Example: "fxcored q gravity event-block-height fx1zgpzdf2uqla7hkx85wnn4p2r3duwqzd8xst6v2",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			orchestratorAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
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

func CmdIbcSequenceHeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ibc-sequence-height [sourcePort] [sourceChannel] [sequence]",
		Short:   "query eth -> ibc sequence block height",
		Example: "fxcored q gravity ibc-sequence-height transfer channel-0 1",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			sequence, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			res, err := queryClient.GetIbcSequenceHeightByChannel(cmd.Context(), &types.QueryIbcSequenceHeightRequest{
				SourcePort:    args[0],
				SourceChannel: args[1],
				Sequence:      sequence,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

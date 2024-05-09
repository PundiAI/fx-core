package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	tranfsertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	abcitype "github.com/tendermint/tendermint/abci/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func GetQueryCmd(subCmd ...*cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the crosschain module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(subCmd...)
	return cmd
}

func GetQuerySubCmds(chainName string) []*cobra.Command {
	cmds := []*cobra.Command{
		// query module params
		CmdGetParams(chainName),

		// query Oracle
		CmdGetOracle(chainName),
		CmdGetOracles(chainName),
		CmdGetOracleReward(chainName),
		CmdGetOracleDelegateAddr(chainName),
		CmdGetProposalOracles(chainName),

		// query oracle set
		CmdGetCurrentOracleSet(chainName),
		CmdGetOracleSetRequest(chainName),

		// need oracle consensus sign
		// 1. oracle set change confirm
		CmdGetLastOracleSetRequests(chainName),
		CmdGetPendingOracleSetRequest(chainName),
		CmdGetOracleSetConfirm(chainName),
		CmdGetOracleSetConfirms(chainName),

		// 2. request batch confirm
		CmdGetPendingOutgoingTXBatchRequest(chainName),
		CmdBatchConfirm(chainName),
		CmdBatchConfirms(chainName),

		// send to external
		CmdBatchRequestByNonce(chainName),
		CmdGetPendingSendToExternal(chainName),
		CmdOutgoingTxBatches(chainName),
		CmdGetBatchFees(chainName),

		CmdGetLastObservedBlockHeight(chainName),
		CmdProjectedBatchTimeoutHeight(chainName),

		// denom <-> external token
		CmdGetDenomToExternalToken(chainName),
		CmdGetExternalTokenToDenom(chainName),
		CmdGetBridgeTokens(chainName),
		CmdGetBridgeCoinByDenom(chainName),

		// 1. oracle event nonce
		CmdGetOracleEventNonce(chainName),
		// 2. event nonce block height
		CmdGetOracleEventBlockHeight(chainName),

		// help cmd.
		CmdCovertBridgeToken(chainName),

		// bridge call
		CmdGetBridgeCalls(chainName),
		CmdGetBridgeCall(chainName),
		CmdBridgeCallByAddr(chainName),
		CmdBridgeCallConfirm(chainName),
		CmdLastPendingBridgeCall(chainName),
	}

	for _, command := range cmds {
		flags.AddQueryFlagsToCmd(command)
	}
	return cmds
}

func CmdGetParams(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current parameters information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(&res.Params)
		},
	}
	return cmd
}

func CmdGetOracles(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracles",
		Short: "Query current oracles",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Oracles(cmd.Context(), &types.QueryOraclesRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetProposalOracles(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-oracles",
		Short: "Query proposal oracles address",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			abciResp, err := clientCtx.QueryABCI(abcitype.RequestQuery{
				Data: types.ProposalOracleKey,
				Path: fmt.Sprintf("/store/%s/key", chainName),
			})
			if err != nil {
				return err
			}
			var proposalOracle types.ProposalOracle
			if err := clientCtx.LegacyAmino.Unmarshal(abciResp.Value, &proposalOracle); err != nil {
				return err
			}
			return clientCtx.PrintProto(&proposalOracle)
		},
	}
	return cmd
}

func CmdGetOracleReward(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward [oracle-address]",
		Short: "Query oracle reward",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := distrtypes.NewQueryClient(clientCtx)

			oracleAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			data := append(oracleAddr.Bytes(), []byte(chainName)...)
			rewards, err := queryClient.DelegationTotalRewards(cmd.Context(), &distrtypes.QueryDelegationTotalRewardsRequest{
				DelegatorAddress: sdk.AccAddress(crypto.Keccak256(data)[12:]).String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(rewards)
		},
	}
	return cmd
}

func CmdGetOracleDelegateAddr(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "oracle-delegate [oracle-address]",
		Short:  "Query oracle delegate address",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			oracleAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			data := append(oracleAddr.Bytes(), []byte(chainName)...)
			cmd.Println(sdk.AccAddress(crypto.Keccak256(data)[12:]).String())
			return nil
		},
	}
	return cmd
}

func CmdGetOracle(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracle [oracle-address|bridge-address]",
		Short: "Query oracle for a given oracle address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			if externalAddress, err := getContractAddr(args[0]); err == nil {
				res, err := queryClient.GetOracleByExternalAddr(cmd.Context(), &types.QueryOracleByExternalAddrRequest{
					ExternalAddress: externalAddress,
					ChainName:       chainName,
				})
				if err != nil {
					return err
				}
				return clientCtx.PrintProto(res.Oracle)
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.GetOracleByAddr(cmd.Context(), &types.QueryOracleByAddrRequest{
				OracleAddress: address.String(),
				ChainName:     chainName,
			})
			if err != nil {
				res, err = queryClient.GetOracleByBridgerAddr(cmd.Context(), &types.QueryOracleByBridgerAddrRequest{
					BridgerAddress: address.String(),
					ChainName:      chainName,
				})
				if err != nil {
					return err
				}
			}
			return clientCtx.PrintProto(res.Oracle)
		},
	}
	return cmd
}

func CmdGetCurrentOracleSet(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-oracle-set",
		Short: "Query current oracle-set",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CurrentOracleSet(cmd.Context(), &types.QueryCurrentOracleSetRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.OracleSet)
		},
	}
	return cmd
}

func CmdGetOracleSetRequest(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracle-set-request [nonce]",
		Short: "Query requested oracle-set with a particular nonce",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			var nonce uint64
			if len(args) == 0 {
				queryAbciResp, err := clientCtx.QueryABCI(abcitype.RequestQuery{
					Path: fmt.Sprintf("store/%s/key", chainName),
					Data: types.LatestOracleSetNonce,
				})
				if err != nil {
					return err
				}
				nonce = sdk.BigEndianToUint64(queryAbciResp.Value)
			} else {
				var err error
				nonce, err = strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					return err
				}
			}
			res, err := queryClient.OracleSetRequest(cmd.Context(), &types.QueryOracleSetRequestRequest{
				ChainName: chainName,
				Nonce:     nonce,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.OracleSet)
		},
	}
	return cmd
}

func CmdGetLastOracleSetRequests(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-oracle-set-requests",
		Short: "Query last oracle set requests",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.LastOracleSetRequests(cmd.Context(), &types.QueryLastOracleSetRequestsRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetPendingOracleSetRequest(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-oracle-set-request [bridger]",
		Short: "Query the latest oracle-set request which has not been signed by a particular oracle bridger",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			bridgerAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.LastPendingOracleSetRequestByAddr(cmd.Context(), &types.QueryLastPendingOracleSetRequestByAddrRequest{
				BridgerAddress: bridgerAddr.String(),
				ChainName:      chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetOracleSetConfirm(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracle-set-confirm [nonce] [bridger-address]",
		Short: "Query oracle-set confirmation with a particular nonce from a particular oracle bridger",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			bridgerAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			res, err := queryClient.OracleSetConfirm(cmd.Context(), &types.QueryOracleSetConfirmRequest{
				Nonce:          nonce,
				BridgerAddress: bridgerAddr.String(),
				ChainName:      chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.Confirm)
		},
	}
	return cmd
}

func CmdGetOracleSetConfirms(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracle-set-confirms [nonce]",
		Short: "Query oracle-set confirmations with a particular nonce",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			res, err := queryClient.OracleSetConfirmsByNonce(cmd.Context(), &types.QueryOracleSetConfirmsByNonceRequest{
				Nonce:     nonce,
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetPendingOutgoingTXBatchRequest(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-batch-request [bridger-address]",
		Short: "Query the latest outgoing TX batch request which has not been signed by a particular oracle bridger address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			bridgerAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.LastPendingBatchRequestByAddr(cmd.Context(), &types.QueryLastPendingBatchRequestByAddrRequest{
				BridgerAddress: bridgerAddr.String(),
				ChainName:      chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.Batch)
		},
	}
	return cmd
}

func CmdBatchConfirm(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-confirm [token-contract] [nonce] [bridger-address]",
		Short: "Query outgoing tx batches confirm by oracle bridger address",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			tokenContract, err := getContractAddr(args[0])
			if err != nil {
				return err
			}
			nonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			bridgerAddr, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}
			res, err := queryClient.BatchConfirm(cmd.Context(), &types.QueryBatchConfirmRequest{
				ChainName:      chainName,
				TokenContract:  tokenContract,
				Nonce:          nonce,
				BridgerAddress: bridgerAddr.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.Confirm)
		},
	}
	return cmd
}

func CmdBatchConfirms(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-confirms [token-contract] [nonce]",
		Short: "Query outgoing tx batches confirms",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			tokenContract, err := getContractAddr(args[0])
			if err != nil {
				return err
			}
			nonce, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			res, err := queryClient.BatchConfirms(cmd.Context(), &types.QueryBatchConfirmsRequest{
				TokenContract: tokenContract,
				Nonce:         uint64(nonce),
				ChainName:     chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdBatchRequestByNonce(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-request [token-contract] [nonce]",
		Short: "Query outgoing tx batches",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			tokenContract, err := getContractAddr(args[0])
			if err != nil {
				return err
			}
			nonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			res, err := queryClient.BatchRequestByNonce(cmd.Context(), &types.QueryBatchRequestByNonceRequest{
				ChainName:     chainName,
				TokenContract: tokenContract,
				Nonce:         nonce,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.Batch)
		},
	}
	return cmd
}

func CmdGetPendingSendToExternal(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-send-to-external [address]",
		Short: "Query pending send to external txs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.GetPendingSendToExternal(cmd.Context(), &types.QueryPendingSendToExternalRequest{
				ChainName:     chainName,
				SenderAddress: addr.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdOutgoingTxBatches(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outgoing-tx-batches",
		Short: "Query outgoing tx batches",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.OutgoingTxBatches(cmd.Context(), &types.QueryOutgoingTxBatchesRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetBatchFees(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-fees",
		Short: "Query a list of send to external transaction fees to be processed",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.BatchFees(cmd.Context(), &types.QueryBatchFeeRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetLastObservedBlockHeight(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-observed-block-height",
		Short: "Query last observed block height",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.LastObservedBlockHeight(cmd.Context(), &types.QueryLastObservedBlockHeightRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdProjectedBatchTimeoutHeight(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "projected-batch-timeout-height",
		Short: "Query projected batch timeout height",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ProjectedBatchTimeoutHeight(cmd.Context(), &types.QueryProjectedBatchTimeoutHeightRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetDenomToExternalToken(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token [denom]",
		Short: "Query contract address from denom",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			denom := args[0]
			if err := sdk.ValidateDenom(denom); err != nil {
				return err
			}
			res, err := queryClient.DenomToToken(cmd.Context(), &types.QueryDenomToTokenRequest{
				ChainName: chainName,
				Denom:     denom,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetExternalTokenToDenom(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "denom [token-contract]",
		Short: "Query denom from contract address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			tokenAddress, err := getContractAddr(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.TokenToDenom(cmd.Context(), &types.QueryTokenToDenomRequest{
				ChainName: chainName,
				Token:     tokenAddress,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetBridgeTokens(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridge-tokens",
		Short: "Query bridge token list",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.BridgeTokens(cmd.Context(), &types.QueryBridgeTokensRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetOracleEventNonce(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event-nonce [bridger-address]",
		Short: "Query last event nonce by bridger address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			bridgerAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.LastEventNonceByAddr(cmd.Context(), &types.QueryLastEventNonceByAddrRequest{
				ChainName:      chainName,
				BridgerAddress: bridgerAddr.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetOracleEventBlockHeight(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event-block-height [bridger-address]",
		Short: "Query last event block height by bridger address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			bridgerAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.LastEventBlockHeightByAddr(cmd.Context(), &types.QueryLastEventBlockHeightByAddrRequest{
				BridgerAddress: bridgerAddr.String(),
				ChainName:      chainName,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdCovertBridgeToken(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "covert-bridge-token [token-contract] [channel ibc]",
		Short: "Covert bridge claim token name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			tokenContract, err := getContractAddr(args[0])
			if err != nil {
				return err
			}
			channelIbc := args[1]
			coinName := types.NewBridgeDenom(chainName, tokenContract)
			if len(channelIbc) > 0 {
				coinName = tranfsertypes.DenomTrace{
					Path:      channelIbc,
					BaseDenom: coinName,
				}.IBCDenom()
			}
			indent, err := json.MarshalIndent(map[string]interface{}{
				"chain_name":     chainName,
				"token_contract": tokenContract,
				"channel_ibc":    channelIbc,
				"coin_name":      coinName,
			}, "", "  ")
			if err != nil {
				return err
			}
			return clientCtx.PrintBytes(indent)
		},
	}
	return cmd
}

func CmdGetBridgeCoinByDenom(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridge-coin [denom]",
		Short: "Query bridge coin from contract address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			denom := args[0]
			res, err := queryClient.BridgeCoinByDenom(cmd.Context(), &types.QueryBridgeCoinByDenomRequest{
				ChainName: chainName,
				Denom:     denom,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdGetBridgeCalls(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridge-calls",
		Short: "Query bridge calls",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.BridgeCalls(cmd.Context(), &types.QueryBridgeCallsRequest{
				ChainName:  chainName,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddPaginationFlagsToCmd(cmd, "bridgecalls")
	return cmd
}

func CmdGetBridgeCall(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridge-call [nonce]",
		Short: "Query bridge call by event nonce",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.BridgeCallByNonce(cmd.Context(), &types.QueryBridgeCallByNonceRequest{
				ChainName:  chainName,
				EventNonce: nonce,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdBridgeCallByAddr(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridge-call-by-sender [address]",
		Short: "Query bridge call by sender",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			senderAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.BridgeCallBySender(cmd.Context(), &types.QueryBridgeCallBySenderRequest{
				ChainName:     chainName,
				SenderAddress: senderAddr.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdBridgeCallConfirm(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridge-call-confirm [nonce]",
		Short: "Query bridge call confirm by event nonce",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.BridgeCallConfirmByNonce(cmd.Context(), &types.QueryBridgeCallConfirmByNonceRequest{
				ChainName:  chainName,
				EventNonce: nonce,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

func CmdLastPendingBridgeCall(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-pending-bridge-call [bridger-address]",
		Short: "Query last pending bridge call for bridger address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			bridgerAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.LastPendingBridgeCallByAddr(cmd.Context(), &types.QueryLastPendingBridgeCallByAddrRequest{
				ChainName:      chainName,
				BridgerAddress: bridgerAddr.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	return cmd
}

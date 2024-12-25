package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	abcitype "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	tranfsertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func GetQueryCmd(moduleName string, subNames ...string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        moduleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", moduleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	for _, chainName := range subNames {
		cmd.AddCommand(GetQueryCmd(chainName))
	}
	if len(subNames) == 0 {
		cmd.AddCommand(getQuerySubCmds(moduleName)...)
	}
	return cmd
}

func getQuerySubCmds(chainName string) []*cobra.Command {
	cmds := []*cobra.Command{
		// query Oracle
		CmdGetOracle(chainName),
		CmdGetOracleReward(chainName),
		CmdGetOracleDelegateAddr(chainName),
		CmdGetProposalOracles(chainName),

		// query oracle set
		CmdGetOracleSetRequest(chainName),

		// denom <-> external token
		CmdCovertBridgeToken(chainName),
		//
		// event nonce
		CmdGetLastObservedEventNonce(chainName),
	}

	for _, command := range cmds {
		flags.AddQueryFlagsToCmd(command)
	}
	return cmds
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

			if err := types.ValidateExternalAddr(chainName, args[0]); err == nil {
				res, err := queryClient.GetOracleByExternalAddr(cmd.Context(), &types.QueryOracleByExternalAddrRequest{
					ExternalAddress: args[0],
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
				queryABCIResp, err := clientCtx.QueryABCI(abcitype.RequestQuery{
					Path: fmt.Sprintf("store/%s/key", chainName),
					Data: types.LatestOracleSetNonce,
				})
				if err != nil {
					return err
				}
				nonce = sdk.BigEndianToUint64(queryABCIResp.Value)
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

func CmdGetLastObservedEventNonce(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-observed-nonce",
		Short: "Query last observed event nonce",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryABCIResp, err := clientCtx.QueryABCI(abcitype.RequestQuery{
				Path: fmt.Sprintf("store/%s/key", chainName),
				Data: types.LastObservedEventNonceKey,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintString(fmt.Sprintf("%d\n", sdk.BigEndianToUint64(queryABCIResp.Value)))
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

			tokenContract := args[0]
			if err := types.ValidateExternalAddr(chainName, tokenContract); err != nil {
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

package cli

import (
	"context"
	"encoding/hex"
	"fmt"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/types"
	erc20types "github.com/functionx/fx-core/x/erc20/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/x/migrate/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "migrate transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(GetMigrateAccountCmd())
	return cmd
}

func GetMigrateAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [to-address]",
		Short: "migrate account to new address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddress := cliCtx.GetFromAddress()
			toAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			ctx := context.Background()

			//check migrate
			queryClient := types.NewQueryClient(cliCtx)
			if _, err := queryClient.MigrateCheckAccount(ctx, &types.QueryMigrateCheckAccountRequest{From: fromAddress.String(), To: toAddress.String()}); err != nil {
				return err
			}

			//convert coin
			msgs, err := getConvertCoinMsg(cliCtx, ctx, fromAddress, toAddress)
			if err != nil {
				return err
			}

			//migrate account
			msg, err := getMigrateAccountMsg(cliCtx, fromAddress, toAddress)
			msgs = append(msgs, msg)
			//sign and broadcast tx
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msgs...)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getConvertCoinMsg(cliCtx client.Context, ctx context.Context, from, to sdk.AccAddress) ([]sdk.Msg, error) {
	//query balances
	bankClient := banktypes.NewQueryClient(cliCtx)
	respBalances, err := bankClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{Address: from.String()})
	if err != nil {
		return nil, err
	}
	if len(respBalances.Balances) == 0 {
		return nil, nil
	}
	//query pairs
	erc20Client := erc20types.NewQueryClient(cliCtx)
	respPairs, err := erc20Client.TokenPairs(ctx, &erc20types.QueryTokenPairsRequest{})
	if err != nil {
		return nil, err
	}
	supportDenom := make(map[string]bool, len(respPairs.TokenPairs))
	for _, p := range respPairs.TokenPairs {
		supportDenom[p.Denom] = true
	}

	msgs := make([]sdk.Msg, 0, len(respBalances.Balances))
	for _, b := range respBalances.Balances {
		if b.Denom == fxtypes.DefaultDenom || !supportDenom[b.Denom] {
			continue
		}
		msg := erc20types.NewMsgConvertCoin(b, common.BytesToAddress(to.Bytes()), from)
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func getMigrateAccountMsg(cliCtx client.Context, from, to sdk.AccAddress) (sdk.Msg, error) {
	toInfo, _ := cliCtx.Keyring.KeyByAddress(to)
	sign, _, err := cliCtx.Keyring.Sign(toInfo.GetName(), types.MigrateAccountSignatureHash(from, to))
	if err != nil {
		return nil, fmt.Errorf("sign migrate signature error %v", err)
	}
	return types.NewMsgMigrateAccount(from, to, hex.EncodeToString(sign)), nil
}

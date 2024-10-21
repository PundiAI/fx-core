package cli

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/v8/contract"
	fxtypes "github.com/functionx/fx-core/v8/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
	"github.com/functionx/fx-core/v8/x/migrate/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Migrate transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(GetMigrateAccountCmd())
	return cmd
}

func GetMigrateAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [hex-address]",
		Short: "migrate account to new address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddress := cliCtx.GetFromAddress()

			if err := contract.ValidateEthereumAddress(args[0]); err != nil {
				return err
			}
			hexAddress := common.HexToAddress(args[0])
			toAddress := sdk.AccAddress(hexAddress.Bytes())
			if _, err := cliCtx.Keyring.KeyByAddress(toAddress); err != nil {
				return err
			}

			ctx := context.Background()

			// check migrate
			queryClient := types.NewQueryClient(cliCtx)
			if _, err := queryClient.MigrateCheckAccount(ctx, &types.QueryMigrateCheckAccountRequest{From: fromAddress.String(), To: hexAddress.String()}); err != nil {
				return err
			}

			// convert coin
			msgs, err := getConvertCoinMsg(ctx, cliCtx, fromAddress, toAddress)
			if err != nil {
				return err
			}

			// migrate account
			msg, err := getMigrateAccountMsg(cliCtx, fromAddress, hexAddress)
			if err != nil {
				return err
			}
			msgs = append(msgs, msg)
			// sign and broadcast tx
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msgs...)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getConvertCoinMsg(ctx context.Context, cliCtx client.Context, from, to sdk.AccAddress) ([]sdk.Msg, error) {
	bankClient := banktypes.NewQueryClient(cliCtx)
	respBalances, err := bankClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{Address: from.String()})
	if err != nil {
		return nil, err
	}
	if len(respBalances.Balances) == 0 {
		return nil, nil
	}

	erc20Client := erc20types.NewQueryClient(cliCtx)
	respPairs, err := erc20Client.TokenPairs(ctx, &erc20types.QueryTokenPairsRequest{})
	if err != nil {
		return nil, err
	}
	supportDenom := make(map[string]bool, len(respPairs.Erc20Tokens))
	for _, p := range respPairs.Erc20Tokens {
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

func getMigrateAccountMsg(cliCtx client.Context, from sdk.AccAddress, to common.Address) (sdk.Msg, error) {
	toInfo, _ := cliCtx.Keyring.KeyByAddress(sdk.AccAddress(to.Bytes()))
	sign, _, err := cliCtx.Keyring.Sign(toInfo.Name, types.MigrateAccountSignatureHash(from, to.Bytes()), signing.SignMode_SIGN_MODE_DIRECT)
	if err != nil {
		return nil, fmt.Errorf("sign migrate signature error %w", err)
	}
	msg := types.NewMsgMigrateAccount(from, to, hex.EncodeToString(sign))
	if err := msg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("validate basic error %w", err)
	}
	return msg, nil
}

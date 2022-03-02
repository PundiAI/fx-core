package cli

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/functionx/fx-core/x/migrate/types"
	"github.com/spf13/cobra"
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
		Use:   "migrate-account [to-address]",
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
			toInfo, err := cliCtx.Keyring.KeyByAddress(toAddress)
			sign, _, err := cliCtx.Keyring.Sign(toInfo.GetName(), types.MigrateAccountSignatureHash(fromAddress, toAddress))
			if err != nil {
				return fmt.Errorf("sign migrate signature error %v", err)
			}
			msg := types.NewMsgMigrateAccount(fromAddress, toAddress, hex.EncodeToString(sign))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

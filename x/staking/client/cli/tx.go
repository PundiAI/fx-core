package cli

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/v6/x/staking/types"
)

// NewTxCmd returns a root CLI command handler for all x/staking transaction commands.
func NewTxCmd() *cobra.Command {
	stakingTxCmd := &cobra.Command{
		Use:                        stakingtypes.ModuleName,
		Short:                      "Staking transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingTxCmd.AddCommand(
		cli.NewCreateValidatorCmd(),
		cli.NewEditValidatorCmd(),
		cli.NewDelegateCmd(),
		cli.NewRedelegateCmd(),
		cli.NewUnbondCmd(),
		cli.NewCancelUnbondingDelegation(),
		NewGrantPrivilegeCmd(),
		NewEditConsensusPubKeyCmd(),
	)

	return stakingTxCmd
}

func NewGrantPrivilegeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant-privilege [validator-address] [to-address]",
		Short: "Grant validator privilege to an account",
		Long: strings.TrimSpace(
			fmt.Sprintf(`create a new grant authorization to an address to execute a transaction on validator:

Examples:
 $ %s tx %s grant-privilege fxvaloper1.. fx1.. --from=fx1..
	`, version.AppName, stakingtypes.ModuleName)),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			fromAddr := clientCtx.GetFromAddress()
			toAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			toInfo, err := clientCtx.Keyring.KeyByAddress(toAddr)
			if err != nil {
				return err
			}
			// signature
			sign, _, err := clientCtx.Keyring.Sign(toInfo.Name, types.GrantPrivilegeSignatureData(valAddr, fromAddr, toAddr))
			if err != nil {
				return fmt.Errorf("sign grant privilege signature error %w", err)
			}
			msg := &types.MsgGrantPrivilege{
				ValidatorAddress: valAddr.String(),
				FromAddress:      fromAddr.String(),
				ToPubkey:         toInfo.PubKey,
				Signature:        hex.EncodeToString(sign),
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewEditConsensusPubKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-consensus-pubkey [validator-address] [pubkey]",
		Short: "Edit an existing validator consensus public key",
		Long: strings.TrimSpace(
			fmt.Sprintf(`edit an existing validator consensus public key:

Examples:
 $ %s tx %s edit-consensus-pubkey fxvaloper1.. '{"@type":"/cosmos.crypto.ed25519.PubKey","key":"...."}' --from=fx1..
	`, version.AppName, stakingtypes.ModuleName)),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			fromAddr := clientCtx.GetFromAddress()
			var pk cryptotypes.PubKey
			if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(args[1]), &pk); err != nil {
				return err
			}

			msg, err := types.NewMsgEditConsensusPubKey(valAddr, fromAddr, pk)
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

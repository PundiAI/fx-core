package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ethermint "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/intrarelayer/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewTxCmd returns a root CLI command handler for certain modules/intrarelayer transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "intrarelayer subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewConvertCoinCmd(),
		NewConvertFIP20Cmd(),
	)
	return txCmd
}

// NewConvertCoinCmd returns a CLI command handler for converting cosmos coins
func NewConvertCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert-coin [coin] [receiver_hex]",
		Short: "Convert a Cosmos coin to FIP20",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			var receiver string
			sender := cliCtx.GetFromAddress()

			if len(args) == 2 {
				receiver = args[1]
				if err := ethermint.ValidateAddress(receiver); err != nil {
					return fmt.Errorf("invalid receiver hex address %w", err)
				}
			} else {
				uid := viper.GetString(flags.FlagFrom)
				key, err := cliCtx.Keyring.Key(uid)
				if err != nil {
					return fmt.Errorf("get account %s pubKey error %v", uid, err)
				}
				eip55Address, err := types.PubKeyToEIP55Address(key.GetPubKey())
				if err != nil {
					return err
				}
				receiver = eip55Address.Hex()
			}

			msg := &types.MsgConvertCoin{
				Coin:     coin,
				Receiver: receiver,
				Sender:   sender.String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewConvertFIP20Cmd returns a CLI command handler for converting FIP20s
func NewConvertFIP20Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert-fip20 [contract-address] [amount] [receiver]",
		Short: "Convert an FIP20 token to Cosmos coin",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			contract := args[0]
			if err := ethermint.ValidateAddress(contract); err != nil {
				return fmt.Errorf("invalid FIP20 contract address %w", err)
			}

			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid amount %s", args[1])
			}

			uid := viper.GetString(flags.FlagFrom)
			if len(uid) == 0 {
				return fmt.Errorf("empty from account")
			}
			key, err := cliCtx.Keyring.Key(uid)
			if err != nil {
				return fmt.Errorf("get account %s pubKey error %v", uid, err)
			}

			receiver := cliCtx.GetFromAddress()
			if len(args) == 3 {
				receiver, err = sdk.AccAddressFromBech32(args[2])
				if err != nil {
					return err
				}
			}

			pubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, key.GetPubKey())
			if err != nil {
				return err
			}

			msg := &types.MsgConvertFIP20{
				ContractAddress: contract,
				Amount:          amount,
				Receiver:        receiver.String(),
				Sender:          key.GetAddress().String(),
				PubKey:          pubKey,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewRegisterCoinProposalCmd implements the command to submit a community-pool-spend proposal
func NewRegisterCoinProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-coin [metadata]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a register coin proposal",
		Long: `Submit a proposal to register a Cosmos coin to the intrarelayer along with an initial deposit.
Upon passing, the 
The proposal details must be supplied via a JSON file.`,
		Example: fmt.Sprintf(`$ %s tx gov submit-proposal register-coin <path/to/metadata.json> --from=<key_or_address>

Where metadata.json contains (example):
	
{
  "description": "PundiX PURSE",
  "denom_units": [
    {
      "denom": "ibc/0000000000000000000000000000000000000000000000000000000000000000",
      "exponent": 0,
      "aliases": []
    },
    {
      "denom": "PURSE",
      "exponent": 18,
      "aliases": []
    }
  ],
  "base": "ibc/0000000000000000000000000000000000000000000000000000000000000000",
  "display": "PURSE"
}`, version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			metadata, err := ParseMetadata(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			content := types.NewRegisterCoinProposal(title, description, metadata)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "1FX", "deposit of proposal")
	if err := cmd.MarkFlagRequired(cli.FlagTitle); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired(cli.FlagDescription); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired(cli.FlagDeposit); err != nil {
		panic(err)
	}
	return cmd
}

// NewRegisterFIP20ProposalCmd implements the command to submit a community-pool-spend proposal
func NewRegisterFIP20ProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "register-fip20 [fip20-address]",
		Args:    cobra.ExactArgs(1),
		Short:   "Submit a proposal to register an FIP20 token",
		Long:    "Submit a proposal to register an FIP20 token to the intrarelayer along with an initial deposit.",
		Example: fmt.Sprintf("$ %s tx gov submit-proposal register-fip20 <fip20-address> --from=<key_or_address>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			fip20Addr := args[0]
			from := clientCtx.GetFromAddress()
			content := types.NewRegisterFIP20Proposal(title, description, fip20Addr)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "1FX", "deposit of proposal")
	if err := cmd.MarkFlagRequired(cli.FlagTitle); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired(cli.FlagDescription); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired(cli.FlagDeposit); err != nil {
		panic(err)
	}
	return cmd
}

// NewToggleTokenRelayProposalCmd implements the command to submit a community-pool-spend proposal
func NewToggleTokenRelayProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "toggle-token-relay [token]",
		Args:    cobra.ExactArgs(1),
		Short:   "Submit a toggle token relay proposal",
		Long:    "Submit a proposal to toggle the relaying of a token pair along with an initial deposit.",
		Example: fmt.Sprintf("$ %s tx gov submit-proposal toggle-token-relay <denom_or_contract> --from=<key_or_address>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			token := args[0]
			content := types.NewToggleTokenRelayProposal(title, description, token)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "1FX", "deposit of proposal")
	if err := cmd.MarkFlagRequired(cli.FlagTitle); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired(cli.FlagDescription); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired(cli.FlagDeposit); err != nil {
		panic(err)
	}
	return cmd
}

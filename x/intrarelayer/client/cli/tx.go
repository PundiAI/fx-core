package cli

import (
	"fmt"
	"github.com/spf13/viper"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	ethermint "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/intrarelayer/types"
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
		NewConvertERC20Cmd(),
		NewProposalCmd(),
	)
	return txCmd
}

// NewConvertCoinCmd returns a CLI command handler for converting cosmos coins
func NewConvertCoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert-coin [coin] [receiver_hex]",
		Short: "Convert a Cosmos coin to ERC20",
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

// NewConvertERC20Cmd returns a CLI command handler for converting ERC20s
func NewConvertERC20Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert-erc20 [contract-address] [amount] [receiver]",
		Short: "Convert an ERC20 token to Cosmos coin",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			contract := args[0]
			if err := ethermint.ValidateAddress(contract); err != nil {
				return fmt.Errorf("invalid ERC20 contract address %w", err)
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

			msg := &types.MsgConvertERC20{
				ContractAddress: contract,
				Amount:          amount,
				Receiver:        receiver.String(),
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

// NewProposalCmd returns a CLI command handler for proposal
func NewProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "proposal",
		Short:                      "intrarelayer proposals",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewInitIntrarelayerProposalCmd(),
		NewRegisterCoinProposalCmd(),
		NewRegisterERC20ProposalCmd(),
		NewToggleTokenRelayProposalCmd(),
		//NewUpdateTokenPairERC20ProposalCmd(),
	)
	cmd.PersistentFlags().String(cli.FlagTitle, "", "title of proposal")
	cmd.PersistentFlags().String(cli.FlagDescription, "", "description of proposal")
	cmd.PersistentFlags().String(cli.FlagDeposit, "1FX", "deposit of proposal")
	return cmd
}

// NewInitIntrarelayerProposalCmd init intrarelayer proposal
func NewInitIntrarelayerProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init-params",
		Short:   "Submit a init params proposal",
		Example: "fxcored tx intrarelayer proposal init-params --deposit=10000000000000000000000FX --title=\"Init intrarelayer module params\" --description=\"about init intrarelayer params description\"",
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

			enableIntrarelayer := viper.GetBool(flagInitIntrarelayerEnableIntrarelayer)
			enableEvmHook := viper.GetBool(flagInitIntrarelayerEnableEvmHook)
			tokenPairVotingPeriod := viper.GetDuration(flagInitIntrarelayerTokenPairVotingPeriod)

			params := types.NewParams(enableIntrarelayer, tokenPairVotingPeriod, enableEvmHook)
			if err := params.Validate(); err != nil {
				return err
			}

			proposal := types.NewInitIntrarelayerProposal(title, description, &params)
			if err := proposal.ValidateBasic(); err != nil {
				return err
			}

			fromAddress := clientCtx.GetFromAddress()
			msg, err := govtypes.NewMsgSubmitProposal(proposal, deposit, fromAddress)
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
	cmd.Flags().Bool(flagInitIntrarelayerEnableIntrarelayer, true, "enable intrarelayer")
	cmd.Flags().Bool(flagInitIntrarelayerEnableEvmHook, true, "enable emv hook")
	cmd.Flags().Duration(flagInitIntrarelayerTokenPairVotingPeriod, time.Hour*24*2, "token pair voting period")
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
  "description": "staking, gas and governance token of the Evmos testnets"
  "denom_units": [
		{
			"denom": "aphoton",
			"exponent": 0,
			"aliases": ["atto photon"]
		},
		{
			"denom": "photon",
			"exponent": 18
		}
	],
	"base": "aphoton",
	"display: "photon",
	"name": "Photon",
	"symbol": "PHOTON"
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

			metadata, err := ParseMetadata(clientCtx.JSONMarshaler, args[0])
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
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewRegisterERC20ProposalCmd implements the command to submit a community-pool-spend proposal
func NewRegisterERC20ProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "register-erc20 [erc20-address]",
		Args:    cobra.ExactArgs(1),
		Short:   "Submit a proposal to register an ERC20 token",
		Long:    "Submit a proposal to register an ERC20 token to the intrarelayer along with an initial deposit.",
		Example: fmt.Sprintf("$ %s tx gov submit-proposal register-erc20 <path/to/proposal.json> --from=<key_or_address>", version.AppName),
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

			erc20Addr := args[0]
			from := clientCtx.GetFromAddress()
			content := types.NewRegisterERC20Proposal(title, description, erc20Addr)

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
	flags.AddTxFlagsToCmd(cmd)
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
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewUpdateTokenPairERC20ProposalCmd implements the command to submit a community-pool-spend proposal
// Deprecated: unused proposal
func NewUpdateTokenPairERC20ProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update-token-pair-erc20 [erc20_address] [new_erc20_address]",
		Args:    cobra.ExactArgs(2),
		Short:   "Submit a update token pair ERC20 proposal",
		Long:    `Submit a proposal to update the ERC20 address of a token pair along with an initial deposit.`,
		Example: fmt.Sprintf("$ %s tx gov submit-proposal update-token-pair-erc20 <path/to/proposal.json> --from=<key_or_address>", version.AppName),
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

			erc20Addr := args[0]
			newERC20Addr := args[1]

			from := clientCtx.GetFromAddress()
			content := types.NewUpdateTokenPairERC20Proposal(title, description, erc20Addr, newERC20Addr)
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
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

const (
	flagInitIntrarelayerEnableIntrarelayer    = "enable-intrarelayer"
	flagInitIntrarelayerEnableEvmHook         = "enable-evm-hook"
	flagInitIntrarelayerTokenPairVotingPeriod = "token-pair-voting-period"
)

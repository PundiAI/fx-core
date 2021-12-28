package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	feemarkettypes "github.com/functionx/fx-core/x/feemarket/types"
	"github.com/spf13/viper"
	"math"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	rpctypes "github.com/functionx/fx-core/rpc/ethereum/types"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        evmtypes.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", evmtypes.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(NewRawTxCmd())
	return cmd
}

// NewRawTxCmd command build cosmos transaction from raw ethereum transaction
func NewRawTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raw [tx-hex]",
		Short: "Build cosmos transaction from raw ethereum transaction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := hexutil.Decode(args[0])
			if err != nil {
				return errors.Wrap(err, "failed to decode ethereum tx hex bytes")
			}

			msg := &evmtypes.MsgEthereumTx{}
			if err := msg.UnmarshalBinary(data); err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			rsp, err := rpctypes.NewQueryClient(clientCtx).Params(cmd.Context(), &evmtypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			tx, err := msg.BuildTx(clientCtx.TxConfig.NewTxBuilder(), rsp.Params.EvmDenom)
			if err != nil {
				return err
			}

			if clientCtx.GenerateOnly {
				json, err := clientCtx.TxConfig.TxJSONEncoder()(tx)
				if err != nil {
					return err
				}

				return clientCtx.PrintString(fmt.Sprintf("%s\n", json))
			}

			if !clientCtx.SkipConfirm {
				out, err := clientCtx.TxConfig.TxJSONEncoder()(tx)
				if err != nil {
					return err
				}

				_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", out)

				buf := bufio.NewReader(os.Stdin)
				ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf, os.Stderr)

				if err != nil || !ok {
					_, _ = fmt.Fprintf(os.Stderr, "%s\n", "canceled transaction")
					return err
				}
			}

			txBytes, err := clientCtx.TxConfig.TxEncoder()(tx)
			if err != nil {
				return err
			}

			// broadcast to a Tendermint node
			res, err := clientCtx.BroadcastTx(txBytes)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func InitEvmParamsProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init-evm-params",
		Short:   "Submit a init evm params proposal",
		Example: fmt.Sprintf(`$ %s tx gov submit-proposal init-evm-params --evm-denom=<denom> --from=<key_or_address>`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			initProposalAmount, err := sdk.ParseCoinsNormalized(viper.GetString(cli.FlagDeposit))
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

			evmParams, err := getEvmParamsByFlags(cmd)
			if err != nil {
				return err
			}
			feeMarketParams, err := getFeeMarkerParamsByFlags(cmd)
			proposal := &evmtypes.InitEvmParamsProposal{
				Title:           title,
				Description:     description,
				EvmParams:       evmParams,
				FeemarketParams: feeMarketParams,
			}
			fromAddress := cliCtx.GetFromAddress()
			msg, err := govtypes.NewMsgSubmitProposal(proposal, initProposalAmount, fromAddress)
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "1FX", "deposit of proposal")
	cmd.Flags().String(flagEvmParamsEvmDenom, "FX", "evm denom represents the token denomination used to run the EVM state transitions.")
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

func getFeeMarkerParamsByFlags(cmd *cobra.Command) (*feemarkettypes.Params, error) {
	NoBaseFee := true
	var BaseFeeChangeDenominator uint32 = 8
	var ElasticityMultiplier uint32 = 2
	var InitialBaseFee int64 = 1000000000
	var EnableHeight int64 = math.MaxInt64
	return &feemarkettypes.Params{
		NoBaseFee:                NoBaseFee,
		BaseFeeChangeDenominator: BaseFeeChangeDenominator,
		ElasticityMultiplier:     ElasticityMultiplier,
		InitialBaseFee:           InitialBaseFee,
		EnableHeight:             EnableHeight,
	}, nil
}

func getEvmParamsByFlags(cmd *cobra.Command) (*evmtypes.Params, error) {
	evmParamsEvmDenom, err := cmd.Flags().GetString(flagEvmParamsEvmDenom)
	if err != nil {
		return nil, err
	}
	homesteadBlock := sdk.ZeroInt()
	daoForkBlock := sdk.ZeroInt()
	eip150Block := sdk.ZeroInt()
	eip155Block := sdk.ZeroInt()
	eip158Block := sdk.ZeroInt()
	byzantiumBlock := sdk.ZeroInt()
	constantinopleBlock := sdk.ZeroInt()
	petersburgBlock := sdk.ZeroInt()
	istanbulBlock := sdk.ZeroInt()
	muirGlacierBlock := sdk.ZeroInt()
	berlinBlock := sdk.ZeroInt()
	londonBlock := sdk.ZeroInt()

	return &evmtypes.Params{
		EvmDenom:     evmParamsEvmDenom,
		EnableCreate: true,
		EnableCall:   true,
		ExtraEIPs:    nil,
		ChainConfig: evmtypes.ChainConfig{
			HomesteadBlock:      &homesteadBlock,
			DAOForkBlock:        &daoForkBlock,
			DAOForkSupport:      true,
			EIP150Block:         &eip150Block,
			EIP150Hash:          common.Hash{}.String(),
			EIP155Block:         &eip155Block,
			EIP158Block:         &eip158Block,
			ByzantiumBlock:      &byzantiumBlock,
			ConstantinopleBlock: &constantinopleBlock,
			PetersburgBlock:     &petersburgBlock,
			IstanbulBlock:       &istanbulBlock,
			MuirGlacierBlock:    &muirGlacierBlock,
			BerlinBlock:         &berlinBlock,
			LondonBlock:         &londonBlock,
		},
	}, nil
}

const (
	flagEvmParamsEvmDenom = "evm-denom"
)

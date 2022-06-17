package cli

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	gethcommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/x/crosschain/types"
)

const (
	flagProposalTitle       = "title"
	flagProposalDescription = "desc"
	flagInitParamsOracles   = "oracles"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Cross chain transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand([]*cobra.Command{
		CmdUpdateChainOraclesProposal(),

		// set bridger address
		CmdCreateOracleBridger(),
		// add oracle stake
		CmdAddOracleDelegate(),
		// send to external chain
		CmdSendToExternal(),
		CmdCancelSendToExternal(),
		CmdRequestBatch(),

		// oracle consensus confirm
		CmdOracleSetConfirm(),
		CmdRequestBatchConfirm(),
	}...)

	return cmd
}

func CmdUpdateChainOraclesProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update-crosschain-oracles [chain-name] [initial proposal stake]",
		Short:   "update cross chain oracles",
		Example: "fxcored tx crosschain update-crosschain-oracles bsc 100000000000000000000FX --title=\"Update Bsc chain oracles\", --desc=\"oracles description\" --oracles <oracles>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			chainName := args[0]
			initProposalAmount, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}
			title, err := cmd.Flags().GetString(flagProposalTitle)
			if err != nil {
				return err
			}
			description, err := cmd.Flags().GetString(flagProposalDescription)
			if err != nil {
				return err
			}

			oracles, err := cmd.Flags().GetStringSlice(flagInitParamsOracles)
			if err != nil {
				return err
			}
			for i, oracle := range oracles {
				oracleAddr, err := sdk.AccAddressFromBech32(oracle)
				if err != nil {
					return err
				}
				oracles[i] = oracleAddr.String()
			}
			proposal := &types.UpdateChainOraclesProposal{
				Title:       title,
				Description: description,
				Oracles:     oracles,
				ChainName:   chainName,
			}
			fromAddress := cliCtx.GetFromAddress()
			msg, err := govtypes.NewMsgSubmitProposal(proposal, initProposalAmount, fromAddress)
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(flagProposalTitle, "", "proposal title")
	cmd.Flags().String(flagProposalDescription, "", "proposal desc")
	cmd.Flags().StringSlice(flagInitParamsOracles, nil, "list of Oracles that have permission to participate in consensus, using comma split")
	return cmd
}

func CmdCreateOracleBridger() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-oracle-bridger [chain-name] [bridger-address] [external-address] [delegate-amount]",
		Short: "Allows oracle to delegate their voting responsibilities to a given key.",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			bridgerAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			externalAddress := args[2]
			amount, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}
			msg := types.MsgBondedOracle{
				OracleAddress:   cliCtx.GetFromAddress().String(),
				BridgerAddress:  bridgerAddr.String(),
				ExternalAddress: externalAddress,
				DelegateAmount:  amount,
				ChainName:       args[0],
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdAddOracleDelegate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-oracle-delegate [chain-name] [delegate-amount]",
		Short: "Allows oracle add delegate.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}
			msg := types.MsgAddDelegate{
				OracleAddress: cliCtx.GetFromAddress().String(),
				Amount:        amount,
				ChainName:     args[0],
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdSendToExternal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-to-external [chain-name] [external-dest] [amount] [bridge-fee]",
		Short: "Adds a new entry to the transaction pool to withdraw an amount from the Ethereum bridge contract",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			externalDestAddr := args[1]
			if strings.HasPrefix(externalDestAddr, "0x") {
				if !gethcommon.IsHexAddress(externalDestAddr) {
					return fmt.Errorf("target address is invalid!address: [%s]", externalDestAddr)
				}
				externalDestAddr = gethcommon.HexToAddress(externalDestAddr).Hex()
			}

			amount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return sdkerrors.Wrap(err, "amount")
			}
			bridgeFee, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return sdkerrors.Wrap(err, "bridge fee")
			}

			msg := types.MsgSendToExternal{
				Sender:    cliCtx.GetFromAddress().String(),
				Dest:      externalDestAddr,
				Amount:    amount,
				BridgeFee: bridgeFee,
				ChainName: args[0],
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCancelSendToExternal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-send-to-external [chain-name] [tx-ID]",
		Short: "Cancel transaction send to external",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgCancelSendToExternal{
				TransactionId: txId,
				Sender:        cliCtx.GetFromAddress().String(),
				ChainName:     args[0],
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRequestBatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build-batch [chain-name] [token-denom] [minimum-fee] [external-fee-receive]",
		Short: "Build a new batch on the fx side for pooled withdrawal transactions",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			minimumFee, ok := sdk.NewIntFromString(args[2])
			if !ok || minimumFee.IsNegative() {
				return fmt.Errorf("miniumu fee is valid, %v\n", args[2])
			}
			feeReceive := args[3]
			if strings.HasPrefix(feeReceive, "0x") {
				if !gethcommon.IsHexAddress(feeReceive) {
					return fmt.Errorf("invalid feeReceive address: %v", feeReceive)
				}
				feeReceive = gethcommon.HexToAddress(feeReceive).Hex()
			}
			baseFee := sdk.ZeroInt()
			baseFeeStr, err := cmd.Flags().GetString("base-fee")
			if err == nil {
				baseFeeStr = strings.TrimSpace(baseFeeStr)
				if len(baseFeeStr) > 0 {
					baseFee, ok = sdk.NewIntFromString(baseFeeStr)
					if !ok {
						return fmt.Errorf("invalid baseFee:%v", baseFeeStr)
					}
				}
			}
			denom := args[1]
			queryClient := types.NewQueryClient(clientCtx)
			token, err := queryClient.DenomToToken(cmd.Context(), &types.QueryDenomToTokenRequest{
				Denom:     denom,
				ChainName: args[0],
			})
			if err != nil {
				return err
			}
			err = clientCtx.PrintString(fmt.Sprintf("build-batch:\n\tdenom:%s\n\ttoken:%s\n\tchannelIBC:%s\n\n", denom, token.Token, token.ChannelIbc))
			if err != nil {
				fmt.Printf("print denom data err:%v\n", err)
			}
			msg := &types.MsgRequestBatch{
				Sender:     clientCtx.GetFromAddress().String(),
				Denom:      denom,
				MinimumFee: minimumFee,
				FeeReceive: feeReceive,
				ChainName:  args[0],
				BaseFee:    &baseFee,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String("base-fee", "", "requestBatch baseFee, is empty is sdk.ZeroInt")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRequestBatchConfirm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-batch-confirm [chain-name] [contract-address] [nonce] [private-key]",
		Short: "Send valset confirm msg",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddress := clientCtx.GetFromAddress()
			tokenContract := args[1]
			if strings.HasPrefix(tokenContract, "0x") {
				if !gethcommon.IsHexAddress(tokenContract) {
					return fmt.Errorf("invalid contract address:%v", tokenContract)
				}
				tokenContract = gethcommon.HexToAddress(tokenContract).Hex()
			}

			nonce, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			privateKey, err := recoveryPrivateKeyByKeystore(args[3])
			if err != nil {
				return err
			}
			externalAddress := ethCrypto.PubkeyToAddress(privateKey.PublicKey)

			queryClient := types.NewQueryClient(clientCtx)
			batchRequestByNonceResp, err := queryClient.BatchRequestByNonce(cmd.Context(), &types.QueryBatchRequestByNonceRequest{
				Nonce:         nonce,
				TokenContract: tokenContract,
				ChainName:     args[0],
			})
			if err != nil {
				return err
			}
			if batchRequestByNonceResp.Batch == nil {
				return fmt.Errorf("not found batch request by nonce!!!tokenContract:[%v], nonce:[%v]", tokenContract, nonce)
			}
			// Determine whether it has been confirmed
			batchConfirmResp, err := queryClient.BatchConfirm(cmd.Context(), &types.QueryBatchConfirmRequest{
				Nonce:          nonce,
				TokenContract:  tokenContract,
				BridgerAddress: fromAddress.String(),
				ChainName:      args[0],
			})
			if err != nil {
				return err
			}
			if batchConfirmResp.GetConfirm() != nil {
				confirm := batchConfirmResp.GetConfirm()
				return clientCtx.PrintString(fmt.Sprintf(`already confirm requestBatch:
	nonce:[%v]
	tokenContract:[%v]
	bridgerAddress:[%v]
	externalAddress:[%v]
	signature:[%v]
`,
					confirm.Nonce, confirm.TokenContract, confirm.BridgerAddress, confirm.ExternalAddress, confirm.Signature))
			}
			paramsResp, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{ChainName: args[0]})
			if err != nil {
				return err
			}
			checkpoint, err := batchRequestByNonceResp.GetBatch().GetCheckpoint(paramsResp.Params.GetGravityId())
			if err != nil {
				return err
			}
			signature, err := types.NewEthereumSignature(checkpoint, privateKey)
			if err != nil {
				return err
			}
			msg := &types.MsgConfirmBatch{
				Nonce:           nonce,
				TokenContract:   tokenContract,
				ExternalAddress: externalAddress.String(),
				BridgerAddress:  fromAddress.String(),
				Signature:       hex.EncodeToString(signature),
				ChainName:       args[0],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdOracleSetConfirm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracle-set-confirm [chain-name] [nonce] [private-key]",
		Short: "Send oracle-set confirm msg",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddress := clientCtx.GetFromAddress()

			nonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			privateKey, err := recoveryPrivateKeyByKeystore(args[2])
			if err != nil {
				return err
			}
			externalAddress := ethCrypto.PubkeyToAddress(privateKey.PublicKey)

			queryClient := types.NewQueryClient(clientCtx)
			oracleSetRequestResp, err := queryClient.OracleSetRequest(cmd.Context(), &types.QueryOracleSetRequestRequest{Nonce: nonce, ChainName: args[0]})
			if err != nil {
				return err
			}
			// Determine whether it has been confirmed
			oracleSetConfirmResp, err := queryClient.OracleSetConfirm(cmd.Context(), &types.QueryOracleSetConfirmRequest{
				Nonce:          nonce,
				BridgerAddress: fromAddress.String(),
				ChainName:      args[0],
			})
			if err != nil {
				return err
			}
			if oracleSetConfirmResp.GetConfirm() != nil {
				confirm := oracleSetConfirmResp.GetConfirm()
				return fmt.Errorf(`already confirm oracleSet:
	nonce:[%v]
	bridgerAddress:[%v]
	externalAddress:[%v]
	signature:[%v]
`, confirm.Nonce, confirm.BridgerAddress, confirm.ExternalAddress, confirm.Signature)
			}
			paramsResp, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{
				ChainName: args[0],
			})
			if err != nil {
				return err
			}
			checkpoint, err := oracleSetRequestResp.GetOracleSet().GetCheckpoint(paramsResp.Params.GetGravityId())
			if err != nil {
				return err
			}
			signature, err := types.NewEthereumSignature(checkpoint, privateKey)
			if err != nil {
				return err
			}
			msg := &types.MsgOracleSetConfirm{
				Nonce:           nonce,
				BridgerAddress:  fromAddress.String(),
				ExternalAddress: externalAddress.String(),
				Signature:       hex.EncodeToString(signature),
				ChainName:       args[0],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func recoveryPrivateKeyByKeystore(privateKey string) (*ecdsa.PrivateKey, error) {
	var ethPrivateKey *ecdsa.PrivateKey
	if _, err := os.Stat(privateKey); err == nil {
		file, err := ioutil.ReadFile(privateKey)
		if err != nil {
			return nil, err
		}
		stdinReader, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return nil, err
		}
		password := strings.TrimSpace(stdinReader)
		key, err := keystore.DecryptKey(file, password)
		if err != nil {
			return nil, err
		}
		ethPrivateKey = key.PrivateKey
	} else {
		key, err := ethCrypto.HexToECDSA(privateKey)
		if err != nil {
			return nil, fmt.Errorf("invalid eth private key: %s", err.Error())
		}
		ethPrivateKey = key
	}
	return ethPrivateKey, nil
}

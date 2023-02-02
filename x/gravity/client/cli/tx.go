// nolint:staticcheck
package cli

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	gethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	"github.com/functionx/fx-core/v3/x/gravity/types"
)

const (
	flagEthKeyType  = "eth-key-type"
	flagEthKeystore = "eth-keystore"
	flagEthPassword = "eth-password"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Gravity transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
		Hidden:                     true,
	}

	cmd.AddCommand([]*cobra.Command{
		// set delegate address
		CmdSetOrchestratorAddress(),
		// send to eth
		CmdSendToEth(),
		CmdCancelSendToEth(),
		CmdRequestBatch(),

		// validator consensus confirm
		CmdValidatorSetConfirm(),
		CmdRequestBatchConfirm(),
	}...)

	return cmd
}

func CmdSetOrchestratorAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-orchestrator-address [validator-address] [orchestrator-address] [eth-address]",
		Short: "Allows validators to delegate their voting responsibilities to a given key.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			validatorAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			orchestratorAddress, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			ethAddress := args[2]
			if !gethcommon.IsHexAddress(ethAddress) {
				return fmt.Errorf("invalid eth address: %v", ethAddress)
			}
			msg := types.MsgSetOrchestratorAddress{
				Validator:    validatorAddress.String(),
				Orchestrator: orchestratorAddress.String(),
				EthAddress:   ethAddress,
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdSendToEth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-to-eth [eth-dest] [amount] [bridge-fee]",
		Short: "Adds a new entry to the transaction pool to withdraw an amount from the Ethereum bridge contract",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			sender := cliCtx.GetFromAddress()

			if !gethcommon.IsHexAddress(args[0]) {
				return fmt.Errorf("invalid eth-dest address: %v", args[0])
			}
			amount, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return sdkerrors.Wrap(err, "amount")
			}
			bridgeFee, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return sdkerrors.Wrap(err, "bridge fee")
			}

			if len(amount) > 1 || len(bridgeFee) > 1 {
				return fmt.Errorf("coin amounts too long, expecting just 1 coin amount for both amount and bridgeFee")
			}

			// Make the message
			msg := types.MsgSendToEth{
				Sender:    sender.String(),
				EthDest:   args[0],
				Amount:    amount[0],
				BridgeFee: bridgeFee[0],
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCancelSendToEth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-send-to-eth [txID]",
		Short: "Cancel transaction send to eth",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			senderAddr := cliCtx.GetFromAddress()
			txId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			msg := &types.MsgCancelSendToEth{
				TransactionId: txId,
				Sender:        senderAddr.String(),
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRequestBatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build-batch [token_contract] [minimum_fee] [eth_fee_receive]",
		Short: "Build a new batch on the fxcore side for pooled withdrawal transactions",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddress := cliCtx.GetFromAddress()

			minimumFee, ok := sdk.NewIntFromString(args[1])
			if !ok || minimumFee.IsNegative() {
				return fmt.Errorf("miniumu fee is valid, fee: %v", args[1])
			}
			ethFeeReceive := args[2]
			if !gethcommon.IsHexAddress(ethFeeReceive) {
				return fmt.Errorf("invalid ethFeeReceive address: %v", args[2])
			}
			baseFee := sdk.ZeroInt()
			baseFeeStr, err := cmd.Flags().GetString("base-fee")
			if err == nil {
				baseFeeStr = strings.TrimSpace(baseFeeStr)
				if len(baseFeeStr) > 0 {
					baseFee, ok = sdk.NewIntFromString(baseFeeStr)
					if !ok {
						return fmt.Errorf("invalid baseFee: %v", baseFeeStr)
					}
				}
			}
			msg := &types.MsgRequestBatch{
				Sender:     fromAddress.String(),
				Denom:      args[0],
				MinimumFee: minimumFee,
				FeeReceive: ethFeeReceive,
				BaseFee:    baseFee,
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String("base-fee", "", "requestBatch baseFee, is empty is sdk.ZeroInt")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdRequestBatchConfirm will be deprecated
//
//gocyclo:ignore
func CmdRequestBatchConfirm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-batch-confirm [contractAddress] [nonce] [...hexEthPrivate]",
		Short: "Send valset confirm msg",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddress := clientCtx.GetFromAddress()
			contractAddress := args[0]
			if !gethcommon.IsHexAddress(contractAddress) {
				return fmt.Errorf("invalid contract address: %v", contractAddress)
			}

			nonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			ethKeyType, err := cmd.Flags().GetString(flagEthKeyType)
			if err != nil {
				return err
			}

			var ethPrivateKey *ecdsa.PrivateKey
			switch ethKeyType {
			case "hex":
				if len(args) < 3 {
					return fmt.Errorf("eth-key-type=hex must input hexEthPrivateKey")
				}
				ethPrivateKey, err = ethcrypto.HexToECDSA(args[2])
				if err != nil {
					return err
				}
			case "keystore":
				keystoreFile, err := cmd.Flags().GetString(flagEthKeystore)
				if err != nil {
					return err
				}
				passwordFile, err := cmd.Flags().GetString(flagEthPassword)
				if err != nil {
					return err
				}
				ethPrivateKey, err = recoveryPrivateKeyByKeystore(keystoreFile, passwordFile)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown eth-key-type flag: %v, support:(keystore|hex)", ethKeyType)
			}
			ethAddress := ethcrypto.PubkeyToAddress(ethPrivateKey.PublicKey)

			queryClient := types.NewQueryClient(clientCtx)
			batchRequestByNonceResp, err := queryClient.BatchRequestByNonce(cmd.Context(), &types.QueryBatchRequestByNonceRequest{
				Nonce:         nonce,
				TokenContract: contractAddress,
			})
			if err != nil {
				return err
			}
			if batchRequestByNonceResp.Batch == nil {
				return fmt.Errorf("not found batch request by nonce! contractAddress: %v, nonce: %v", contractAddress, nonce)
			}
			// Determine whether it has been confirmed
			batchConfirmResp, err := queryClient.BatchConfirm(cmd.Context(), &types.QueryBatchConfirmRequest{
				Nonce:         nonce,
				TokenContract: contractAddress,
				Address:       fromAddress.String(),
			})
			if err != nil {
				return err
			}
			if batchConfirmResp.GetConfirm() != nil {
				confirm := batchConfirmResp.GetConfirm()
				return clientCtx.PrintProto(confirm)
			}
			paramsResp, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}
			outgoingTxBatch := &crosschaintypes.OutgoingTxBatch{
				BatchNonce:    batchRequestByNonceResp.Batch.BatchNonce,
				BatchTimeout:  batchRequestByNonceResp.Batch.BatchTimeout,
				Transactions:  make([]*crosschaintypes.OutgoingTransferTx, len(batchRequestByNonceResp.Batch.Transactions)),
				TokenContract: batchRequestByNonceResp.Batch.TokenContract,
				Block:         batchRequestByNonceResp.Batch.Block,
				FeeReceive:    batchRequestByNonceResp.Batch.FeeReceive,
			}
			for i := 0; i < len(batchRequestByNonceResp.Batch.Transactions); i++ {
				outgoingTxBatch.Transactions[i] = &crosschaintypes.OutgoingTransferTx{
					Id:          batchRequestByNonceResp.Batch.Transactions[i].Id,
					Sender:      batchRequestByNonceResp.Batch.Transactions[i].Sender,
					DestAddress: batchRequestByNonceResp.Batch.Transactions[i].DestAddress,
					Token: crosschaintypes.ERC20Token{
						Contract: batchRequestByNonceResp.Batch.Transactions[i].Erc20Token.Contract,
						Amount:   batchRequestByNonceResp.Batch.Transactions[i].Erc20Token.Amount,
					},
					Fee: crosschaintypes.ERC20Token{
						Contract: batchRequestByNonceResp.Batch.Transactions[i].Erc20Fee.Contract,
						Amount:   batchRequestByNonceResp.Batch.Transactions[i].Erc20Fee.Amount,
					},
				}
			}
			checkpoint, err := outgoingTxBatch.GetCheckpoint(paramsResp.Params.GetGravityId())
			if err != nil {
				return err
			}
			signature, err := crosschaintypes.NewEthereumSignature(checkpoint, ethPrivateKey)
			if err != nil {
				return err
			}
			msg := &types.MsgConfirmBatch{
				Nonce:         nonce,
				TokenContract: contractAddress,
				EthSigner:     ethAddress.String(),
				Orchestrator:  fromAddress.String(),
				Signature:     hex.EncodeToString(signature),
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagEthKeyType, "keystore", "eth private key type(keystore|hex), default:keystore")
	cmd.Flags().String(flagEthKeystore, "", "eth keystore file")
	cmd.Flags().String(flagEthPassword, "", "eth keystore password file")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdValidatorSetConfirm will be deprecated
//
//gocyclo:ignore
func CmdValidatorSetConfirm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "valset-confirm [nonce] [...hexEthPrivate]",
		Short: "Send valset confirm msg",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddress := clientCtx.GetFromAddress()

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			ethKeyType, err := cmd.Flags().GetString(flagEthKeyType)
			if err != nil {
				return err
			}

			var ethPrivateKey *ecdsa.PrivateKey
			switch ethKeyType {
			case "hex":
				if len(args) < 2 {
					return fmt.Errorf("eth-key-type=hex must input hexEthPrivateKey")
				}
				ethPrivateKey, err = ethcrypto.HexToECDSA(args[1])
				if err != nil {
					return errors.WithStack(err)
				}
			case "keystore":
				keystoreFile, err := cmd.Flags().GetString(flagEthKeystore)
				if err != nil {
					return err
				}
				passwordFile, err := cmd.Flags().GetString(flagEthPassword)
				if err != nil {
					return err
				}
				ethPrivateKey, err = recoveryPrivateKeyByKeystore(keystoreFile, passwordFile)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown eth-key-type flag: %v, support:(keystore|hex)", ethKeyType)
			}
			ethAddress := ethcrypto.PubkeyToAddress(ethPrivateKey.PublicKey)

			queryClient := types.NewQueryClient(clientCtx)
			valsetRequestResp, err := queryClient.ValsetRequest(cmd.Context(), &types.QueryValsetRequestRequest{Nonce: nonce})
			if err != nil {
				return err
			}
			// Determine whether it has been confirmed
			valsetConfirmResp, err := queryClient.ValsetConfirm(cmd.Context(), &types.QueryValsetConfirmRequest{
				Nonce:   nonce,
				Address: fromAddress.String(),
			})
			if err != nil {
				return err
			}
			if valsetConfirmResp.GetConfirm() != nil {
				confirm := valsetConfirmResp.GetConfirm()
				return clientCtx.PrintProto(confirm)
			}
			paramsResp, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}
			oracleSet := crosschaintypes.OracleSet{
				Nonce:   valsetRequestResp.Valset.Nonce,
				Members: make([]crosschaintypes.BridgeValidator, len(valsetRequestResp.Valset.Members)),
				Height:  valsetRequestResp.Valset.Height,
			}
			for i := 0; i < len(valsetRequestResp.Valset.Members); i++ {
				oracleSet.Members[i] = crosschaintypes.BridgeValidator{
					Power:           valsetRequestResp.Valset.Members[i].Power,
					ExternalAddress: valsetRequestResp.Valset.Members[i].EthAddress,
				}
			}
			checkpoint, err := oracleSet.GetCheckpoint(paramsResp.Params.GetGravityId())
			if err != nil {
				return err
			}
			signature, err := crosschaintypes.NewEthereumSignature(checkpoint, ethPrivateKey)
			if err != nil {
				return err
			}
			msg := &types.MsgValsetConfirm{
				Nonce:        nonce,
				Orchestrator: fromAddress.String(),
				EthAddress:   ethAddress.String(),
				Signature:    hex.EncodeToString(signature),
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagEthKeyType, "keystore", "eth private key type(keystore|hex), default:keystore")
	cmd.Flags().String(flagEthKeystore, "", "eth keystore file")
	cmd.Flags().String(flagEthPassword, "", "eth keystore password file")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func recoveryPrivateKeyByKeystore(keystoreFile, passwordFile string) (*ecdsa.PrivateKey, error) {
	keystoreData, err := os.ReadFile(keystoreFile)
	if err != nil {
		return nil, errors.WithMessagef(err, "keystoreFile: %s", keystoreFile)
	}
	passwordData, err := os.ReadFile(passwordFile)
	if err != nil {
		return nil, errors.WithMessagef(err, "passwordFile: %s", keystoreFile)
	}
	decryptKey, err := keystore.DecryptKey(keystoreData, string(passwordData))
	if err != nil {
		return nil, errors.WithMessagef(err, "decryptKey err")
	}
	return decryptKey.PrivateKey, nil
}

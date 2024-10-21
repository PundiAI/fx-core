package cli

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func GetTxCmd(moduleName string, subNames ...string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        moduleName,
		Short:                      fmt.Sprintf("%s%s transaction subcommands", strings.ToUpper(moduleName[:1]), moduleName[1:]),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	for _, chainName := range subNames {
		cmd.AddCommand(GetTxCmd(chainName))
	}
	if len(subNames) == 0 {
		cmd.AddCommand(getTxSubCmds(moduleName)...)
	}
	return cmd
}

func getTxSubCmds(chainName string) []*cobra.Command {
	cmds := []*cobra.Command{
		// oracle consensus confirm
		CmdOracleSetConfirm(chainName),
		CmdRequestBatchConfirm(chainName),
	}
	for _, command := range cmds {
		flags.AddTxFlagsToCmd(command)
	}
	return cmds
}

func CmdRequestBatchConfirm(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-batch-confirm [contract-address] [nonce] [private-key]",
		Short: "Send batch confirm msg",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fromAddress := clientCtx.GetFromAddress()

			tokenContract := args[0]
			nonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			privateKey, err := recoveryPrivateKeyByKeystore(args[2])
			if err != nil {
				return err
			}
			externalAddress := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

			queryClient := types.NewQueryClient(clientCtx)
			outgoingTxBatchResp, err := queryClient.OutgoingTxBatch(cmd.Context(), &types.QueryOutgoingTxBatchRequest{
				Nonce:         nonce,
				TokenContract: tokenContract,
				ChainName:     chainName,
			})
			if err != nil {
				return err
			}
			if outgoingTxBatchResp.Batch == nil {
				return fmt.Errorf("not found batch request by nonce, tokenContract: %v, nonce: %v", tokenContract, nonce)
			}
			// Determine whether it has been confirmed
			batchConfirmResp, err := queryClient.BatchConfirm(cmd.Context(), &types.QueryBatchConfirmRequest{
				Nonce:          nonce,
				TokenContract:  tokenContract,
				BridgerAddress: fromAddress.String(),
				ChainName:      chainName,
			})
			if err != nil {
				return err
			}
			if batchConfirmResp.GetConfirm() != nil {
				confirm := batchConfirmResp.GetConfirm()
				return clientCtx.PrintProto(confirm)
			}
			paramsResp, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{
				ChainName: chainName,
			})
			if err != nil {
				return err
			}
			checkpoint, err := outgoingTxBatchResp.GetBatch().GetCheckpoint(paramsResp.Params.GetGravityId())
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
				ChainName:       chainName,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

func CmdOracleSetConfirm(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracle-set-confirm [nonce] [private-key]",
		Short: "Send oracle-set confirm msg",
		Args:  cobra.ExactArgs(2),
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
			privateKey, err := recoveryPrivateKeyByKeystore(args[1])
			if err != nil {
				return err
			}
			externalAddress := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

			queryClient := types.NewQueryClient(clientCtx)
			oracleSetRequestResp, err := queryClient.OracleSetRequest(cmd.Context(), &types.QueryOracleSetRequestRequest{
				Nonce: nonce, ChainName: chainName,
			})
			if err != nil {
				return err
			}
			// Determine whether it has been confirmed
			oracleSetConfirmResp, err := queryClient.OracleSetConfirm(cmd.Context(), &types.QueryOracleSetConfirmRequest{
				Nonce:          nonce,
				BridgerAddress: fromAddress.String(),
				ChainName:      chainName,
			})
			if err != nil {
				return err
			}
			if oracleSetConfirmResp.GetConfirm() != nil {
				confirm := oracleSetConfirmResp.GetConfirm()
				return clientCtx.PrintProto(confirm)
			}
			paramsResp, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{
				ChainName: chainName,
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
				ChainName:       chainName,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

func recoveryPrivateKeyByKeystore(privateKey string) (*ecdsa.PrivateKey, error) {
	var ethPrivateKey *ecdsa.PrivateKey
	if _, err := os.Stat(privateKey); err == nil {
		file, err := os.ReadFile(privateKey)
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
		key, err := ethcrypto.HexToECDSA(privateKey)
		if err != nil {
			return nil, fmt.Errorf("invalid eth private key: %s", err.Error())
		}
		ethPrivateKey = key
	}
	return ethPrivateKey, nil
}

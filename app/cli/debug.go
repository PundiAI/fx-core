package cli

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/cosmos-sdk/types/bech32/legacybech32"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/tendermint/tendermint/privval"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/client"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func Debug() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Tool for helping with debugging your application",
		RunE:  client.ValidateCmd,
	}
	cmd.AddCommand(
		HexToString(),
		Base64ToString(),
		ModuleAddressCmd(),
		CovertTxDataToHash(),
		ParseTx(),
		ChecksumEthAddress(),
		PubkeyCmd(),
		VerifyTx(),
	)
	return cmd
}

func HexToString() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hex-to-str [hex]",
		Short: "Hex to string tools",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hexStr := strings.TrimPrefix(args[0], "0x")
			decodeString, err := hex.DecodeString(hexStr)
			if err != nil {
				return err
			}
			cmd.Println(string(decodeString))
			return nil
		},
	}
	return cmd
}

func Base64ToString() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "base64-to-str [hex]",
		Short: "Base64 to string tools",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			decodeString, err := base64.StdEncoding.DecodeString(args[0])
			if err != nil {
				return err
			}
			cmd.Println(string(decodeString))
			return nil
		},
	}
	return cmd
}

func ModuleAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module-addr <module name>",
		Short: "Get module address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println(types.NewModuleAddress(args[0]).String())
			return nil
		},
	}
	return cmd
}

func ParseTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parse-tx [base64TxData]",
		Short:   "parse tx  base64 tx data and print",
		Example: "fxcored debug parse-tx CucHC===...",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			txBytes, err := base64.StdEncoding.DecodeString(args[0])
			if err != nil {
				return err
			}
			tx, err := clientCtx.TxConfig.TxDecoder()(txBytes)
			if err != nil {
				return err
			}
			jsonMarshal, err := clientCtx.TxConfig.TxJSONEncoder()(tx)
			if err != nil {
				return err
			}
			return clientCtx.PrintBytes(jsonMarshal)
		},
	}
	return cmd
}

func VerifyTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "verify-tx [base64TxData]",
		Short:   "verify tx",
		Example: "fxcored debug verify-tx CucHC===...",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			txBytes, err := base64.StdEncoding.DecodeString(args[0])
			if err != nil {
				return err
			}
			tx, err := clientCtx.TxConfig.TxDecoder()(txBytes)
			if err != nil {
				return err
			}

			builder, err := clientCtx.TxConfig.WrapTxBuilder(tx)
			if err != nil {
				return err
			}
			stdTx := builder.GetTx()

			sigTx, ok := tx.(authsigning.SigVerifiableTx)
			if !ok {
				return errors.New("invalid transaction type")
			}
			// stdSigs contains the sequence number, account number, and signatures.
			// When simulating, this would just be a 0-length slice.
			sigs, err := sigTx.GetSignaturesV2()
			if err != nil {
				return fmt.Errorf("get signature error %s", err.Error())
			}
			signerAddrs := sigTx.GetSigners()

			// check that signer length and signature length are the same
			if len(sigs) != len(signerAddrs) {
				return fmt.Errorf("invalid number of signer;  expected: %d, got %d", len(signerAddrs), len(sigs))
			}
			status, err := clientCtx.Client.Status(cmd.Context())
			if err != nil {
				return err
			}
			chainId := status.NodeInfo.Network
			queryClient := authtypes.NewQueryClient(clientCtx)
			for i, sig := range sigs {
				accountResponse, err := queryClient.Account(cmd.Context(), &authtypes.QueryAccountRequest{Address: sdk.AccAddress(signerAddrs[i]).String()})
				if err != nil {
					return err
				}
				var acc authtypes.AccountI
				err = clientCtx.InterfaceRegistry.UnpackAny(accountResponse.GetAccount(), &acc)
				if err != nil {
					return err
				}
				// retrieve pubkey
				pubKey := acc.GetPubKey()
				sequence := sig.Sequence
				signerData := authsigning.SignerData{
					ChainID:       chainId,
					AccountNumber: acc.GetAccountNumber(),
					Sequence:      sequence,
				}

				bz := legacytx.StdSignBytes(
					chainId, acc.GetAccountNumber(), sequence, stdTx.GetTimeoutHeight(),
					legacytx.StdFee{Amount: stdTx.GetFee(), Gas: stdTx.GetGas()},
					tx.GetMsgs(), stdTx.GetMemo(),
				)
				err = clientCtx.PrintBytes(bz)
				if err != nil {
					return err
				}

				err = authsigning.VerifySignature(pubKey, signerData, sig.Data, clientCtx.TxConfig.SignModeHandler(), tx)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	return cmd
}

func CovertTxDataToHash() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tx-hash [base64TxData]",
		Short:   "covert base64 tx data to txHash",
		Example: "fxcored debug tx-hash CucHC===...",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBytes, err := base64.StdEncoding.DecodeString(args[0])
			if err != nil {
				return err
			}
			hashBytes := sha256.Sum256(txBytes)
			cmd.Println(fmt.Sprintf("%X\n", hashBytes))
			return nil
		},
	}
	return cmd
}

func ChecksumEthAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checksum [eth address]",
		Short: "checksum eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			if !gethCommon.IsHexAddress(args[0]) {
				return fmt.Errorf("not hex address:%s", args[0])
			}
			return clientCtx.PrintString(fmt.Sprintf("%s\n", gethCommon.HexToAddress(args[0]).Hex()))
		},
	}
	return cmd
}

func PubkeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pubkey [pubkey]",
		Short: "Decode a pubkey from proto JSON",
		Example: fmt.Sprintf(`
$ %s debug pubkey '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AurroA7jvfPd1AadmmOvWM2rJSwipXfRf8yD6pLbA2DJ"}'
$ %s debug pubkey '{"@type":"/cosmos.crypto.ed25519.PubKey","key":"eKlxn6Xoe9LNmD53omoNQrVrws5KT73hfmqeCSqL87A="}'
`, version.AppName, version.AppName),
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			var pubkey cryptotypes.PubKey
			if len(args) <= 0 {
				serverCtx := server.GetServerContextFromCmd(cmd)
				serverCfg := serverCtx.Config
				privValidator := privval.LoadFilePV(serverCfg.PrivValidatorKeyFile(), serverCfg.PrivValidatorStateFile())
				valPubKey, err := privValidator.GetPubKey()
				if err != nil {
					return err
				}
				pubkey, err = cryptocodec.FromTmPubKeyInterface(valPubKey)
				if err != nil {
					return err
				}
			}
			if err = clientCtx.Codec.UnmarshalInterfaceJSON([]byte(args[0]), &pubkey); err != nil {
				if pubkey, err = legacybech32.UnmarshalPubKey(legacybech32.ConsPK, args[0]); err == nil {
				} else if pubkey, err = legacybech32.UnmarshalPubKey(legacybech32.AccPK, args[0]); err == nil {
				} else if pubkey, err = legacybech32.UnmarshalPubKey(legacybech32.ValPK, args[0]); err == nil {
				} else {
					return fmt.Errorf("pubkey '%s' invalid", args[0])
				}
			}
			var data []byte
			switch pubkey.Type() {
			case "ed25519":
				data, err = json.MarshalIndent(map[string]interface{}{
					"address":        strings.ToUpper(hex.EncodeToString(pubkey.Address().Bytes())),
					"val_cons_pub":   pubkey,
					"pub_key_hex":    hex.EncodeToString(pubkey.Bytes()),
					"pub_key_base64": base64.StdEncoding.EncodeToString(pubkey.Bytes()),
				}, "", "  ")
			case "secp256k1":
				data, err = json.MarshalIndent(map[string]interface{}{
					"eip55_address":  common.BytesToAddress(pubkey.Address()).String(),
					"acc_address":    sdk.AccAddress(pubkey.Address().Bytes()).String(),
					"val_address":    sdk.ValAddress(pubkey.Address().Bytes()).String(),
					"pub_key_hex":    hex.EncodeToString(pubkey.Bytes()),
					"pub_key_base64": base64.StdEncoding.EncodeToString(pubkey.Bytes()),
				}, "", "  ")
			default:
				return fmt.Errorf("invalied public key type")
			}
			if err != nil {
				return err
			}
			return clientCtx.PrintString(string(data))
		},
	}
}

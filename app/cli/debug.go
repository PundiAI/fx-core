package cli

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/debug"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32/legacybech32" //nolint
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/privval"
)

func Debug() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Tool for helping with debugging your application",
		RunE:  client.ValidateCmd,
	}
	cmd.AddCommand(
		ToStringCmd(),
		ToBytes32Cmd(),
		ModuleAddressCmd(),
		ChecksumEthAddressCmd(),
		CovertTxDataToHashCmd(),
		DecodeSimulateTxCmd(),
		VerifyTxCmd(),
		PubkeyCmd(),
		AddrCmd(),
		debug.RawBytesCmd(),
	)
	return cmd
}

func ToStringCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "2str [hex/base64/base58] [data]",
		Short: "Decode to string tools",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var decodeString []byte
			switch args[0] {
			case "hex":
				decodeString, err = hexutil.Decode(args[1])
				if err != nil {
					return err
				}
			case "base64":
				decodeString, err = base64.StdEncoding.DecodeString(args[1])
				if err != nil {
					return err
				}
			case "base58":
				decodeString = base58.Decode(args[1])
			default:
				return fmt.Errorf("invalid encode type: %s", args[0])
			}
			cmd.Println(string(decodeString))
			return nil
		},
	}
}

func ToBytes32Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "2bytes32",
		Short: "String to bytes32 hex",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args[0]) > 32 {
				return fmt.Errorf("input data leng greater than 32")
			}
			var byte32 [32]byte
			copy(byte32[:], args[0])
			cmd.Println(hex.EncodeToString(byte32[:]))
			return nil
		},
	}
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

func VerifyTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "verify-tx [base64TxData]",
		Short:   "Verify tx",
		Example: fmt.Sprintf("%s debug verify-tx 'CucHC...==='", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			txBytes, err := base64.StdEncoding.DecodeString(args[0])
			if err != nil {
				return err
			}
			sdkTx, err := clientCtx.TxConfig.TxDecoder()(txBytes)
			if err != nil {
				return err
			}

			builder, err := clientCtx.TxConfig.WrapTxBuilder(sdkTx)
			if err != nil {
				return err
			}
			stdTx := builder.GetTx()

			sigTx, ok := sdkTx.(authsigning.SigVerifiableTx)
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
				accountResponse, err := queryClient.Account(cmd.Context(), &authtypes.QueryAccountRequest{Address: signerAddrs[i].String()})
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
					sdkTx.GetMsgs(), stdTx.GetMemo(),
				)
				if err = clientCtx.PrintString(string(bz) + "\n"); err != nil {
					return err
				}

				if err = authsigning.VerifySignature(pubKey, signerData, sig.Data, clientCtx.TxConfig.SignModeHandler(), sdkTx); err != nil {
					return err
				}
			}
			return nil
		},
	}
	return cmd
}

func CovertTxDataToHashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tx-hash [base64TxData]",
		Short:   "Covert base64 tx data to txHash",
		Example: fmt.Sprintf("%s debug tx-hash 'CucHC...==='", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBytes, err := base64.StdEncoding.DecodeString(args[0])
			if err != nil {
				return err
			}
			hashBytes := sha256.Sum256(txBytes)
			cmd.Println(fmt.Sprintf("%X", hashBytes))
			return nil
		},
	}
	return cmd
}

func DecodeSimulateTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simulate-tx [base64/hex]",
		Short: "decode base64 or hex tx data to json",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			var txBytes []byte
			if useHex, _ := cmd.Flags().GetBool("hex"); useHex {
				txBytes, err = hexutil.Decode(args[0])
			} else {
				txBytes, err = base64.StdEncoding.DecodeString(args[0])
			}
			if err != nil {
				return err
			}

			var simulateReq = new(tx.SimulateRequest)
			if err := proto.Unmarshal(txBytes, simulateReq); err != nil {
				return err
			}
			return clientCtx.PrintProto(simulateReq)
		},
	}
	cmd.Flags().BoolP("hex", "x", false, "Treat input as hexadecimal instead of base64")
	return cmd
}

func ChecksumEthAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checksum [eth address]",
		Short: "Checksum eth address",
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
			} else {
				if err = clientCtx.Codec.UnmarshalInterfaceJSON([]byte(args[0]), &pubkey); err != nil {
					// nolint
					if pubkey, err = legacybech32.UnmarshalPubKey(legacybech32.ConsPK, args[0]); err == nil {
					} else if pubkey, err = legacybech32.UnmarshalPubKey(legacybech32.AccPK, args[0]); err == nil {
					} else if pubkey, err = legacybech32.UnmarshalPubKey(legacybech32.ValPK, args[0]); err == nil {
					} else {
						return fmt.Errorf("pubkey '%s' invalid", args[0])
					}
				}
			}
			pubkeyJson, err := clientCtx.Codec.MarshalInterfaceJSON(pubkey)
			if err != nil {
				return err
			}
			var data []byte
			switch pubkey.Type() {
			case "ed25519":
				data, err = json.MarshalIndent(map[string]interface{}{
					"address":        strings.ToUpper(hex.EncodeToString(pubkey.Address().Bytes())),
					"val_cons_pub":   json.RawMessage(pubkeyJson),
					"pub_key_hex":    hex.EncodeToString(pubkey.Bytes()),
					"pub_key_base64": base64.StdEncoding.EncodeToString(pubkey.Bytes()),
				}, "", "  ")
			case "secp256k1":
				data, err = json.MarshalIndent(map[string]interface{}{
					"acc_address":    sdk.AccAddress(pubkey.Address().Bytes()).String(),
					"val_address":    sdk.ValAddress(pubkey.Address().Bytes()).String(),
					"pub_key_hex":    hex.EncodeToString(pubkey.Bytes()),
					"pub_key_base64": base64.StdEncoding.EncodeToString(pubkey.Bytes()),
				}, "", "  ")
			case "eth_secp256k1":
				data, err = json.MarshalIndent(map[string]interface{}{
					"eip55_address":  common.BytesToAddress(pubkey.Address()).String(),
					"acc_address":    sdk.AccAddress(pubkey.Address().Bytes()).String(),
					"val_address":    sdk.ValAddress(pubkey.Address().Bytes()).String(),
					"pub_key_hex":    hex.EncodeToString(pubkey.Bytes()),
					"pub_key_base64": base64.StdEncoding.EncodeToString(pubkey.Bytes()),
				}, "", "  ")
			default:
				return fmt.Errorf("invalied public key type: %s", pubkey.Type())
			}
			if err != nil {
				return err
			}
			return clientCtx.PrintString(string(data))
		},
	}
}

func AddrCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "addr [address]",
		Short: "Convert an address between hex and bech32",
		Long: fmt.Sprintf(`Convert an address between hex encoding and bech32.

Example:
$ %s debug addr fx1e0jnq2sun3dzjh8p2xq95kk0expwmd7sd7r5ye
			`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			addrString := args[0]
			var addr []byte

			// try hex, then bech32
			addr, err = hexutil.Decode(addrString)
			if err != nil {
				var err2 error
				addr, err2 = sdk.AccAddressFromBech32(addrString)
				if err2 != nil {
					var err3 error
					addr, err3 = sdk.ValAddressFromBech32(addrString)

					if err3 != nil {
						return fmt.Errorf("expected hex or bech32. Got errors: hex: %v, bech32 acc: %v, bech32 val: %v", err, err2, err3)

					}
				}
			}

			return PrintOutput(clientCtx, map[string]interface{}{
				"base64_address": addr,
				"hex_address":    hex.EncodeToString(addr),
				"eip55_address":  common.BytesToAddress(addr).String(),
				"acc_address":    sdk.AccAddress(addr),
				"val_address":    sdk.ValAddress(addr),
			})
		},
	}
}

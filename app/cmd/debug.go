package cmd

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/tendermint/tendermint/privval"
	"strconv"
	"strings"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/x/auth/types"

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
		HexExternalAddress(),
		PubkeyCmd(),
		AddrCmd(),
		RawBytesCmd(),
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

func HexExternalAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "checksum [externalAddress]",
		Short:   "HexToAddress and hex address",
		Example: "fxcored q crosschain checksum 0x3f6795b8abe0775a88973469909ade1405f7ac09",
		Args:    cobra.ExactArgs(1),
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
		Long:  "Decode a pubkey from proto JSON and display it's address",
		Example: fmt.Sprintf(
			`"$ %s debug pubkey '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AurroA7jvfPd1AadmmOvWM2rJSwipXfRf8yD6pLbA2DJ"}'`,
			version.AppName,
		),
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			var pk cryptotypes.PubKey
			if len(args) <= 0 {
				serverCtx := server.GetServerContextFromCmd(cmd)
				serverCfg := serverCtx.Config
				privValidator := privval.LoadFilePV(serverCfg.PrivValidatorKeyFile(), serverCfg.PrivValidatorStateFile())
				valPubKey, err := privValidator.GetPubKey()
				if err != nil {
					return err
				}
				pk, err = cryptocodec.FromTmPubKeyInterface(valPubKey)
				if err != nil {
					return err
				}
			} else {
				err := clientCtx.JSONMarshaler.UnmarshalInterfaceJSON([]byte(args[0]), &pk)
				if err != nil {
					return err
				}
			}
			addr := pk.Address()
			consPubAddress, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, pk)
			if err != nil {
				return err
			}
			data, err := json.Marshal(map[string]interface{}{
				"BytesAddress": fmt.Sprintf("%v", addr.Bytes()),
				"HexAddress":   addr.String(),
				"EIP55Address": common.BytesToAddress(addr),
				"AccAddress":   sdk.AccAddress(addr).String(),
				"ConsAddress":  sdk.ConsAddress(addr).String(),
				"PubKeyHex":    hex.EncodeToString(pk.Bytes()),
				"ConsPub":      consPubAddress,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintString(string(data))
		},
	}
}

func AddrCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addr [address]",
		Short: "Convert an address between hex and bech32",
		Long:  "Convert an address between hex encoding and bech32.",
		Example: fmt.Sprintf(
			`$ %s debug addr fx1walm8hnaylycy5s45uc9n3cc50mwhhhcep9ksl
$ %s debug addr 0xA588C66983a81e800Db4dF74564F09f91c026351`, version.AppName, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			addrString := args[0]
			cfg := sdk.GetConfig()

			var addr []byte
			switch {
			case common.IsHexAddress(addrString):
				addr = common.HexToAddress(addrString).Bytes()
			case strings.HasPrefix(addrString, cfg.GetBech32ValidatorAddrPrefix()):
				addr, _ = sdk.ValAddressFromBech32(addrString)
			case strings.HasPrefix(addrString, cfg.GetBech32AccountAddrPrefix()):
				addr, _ = sdk.AccAddressFromBech32(addrString)
			default:
				return fmt.Errorf("expected a valid hex or bech32 address (acc prefix %s), got '%s'", cfg.GetBech32AccountAddrPrefix(), addrString)
			}

			addressPrefix, _ := cmd.Flags().GetString("addr-prefix")
			UpdateAddressPrefix(addressPrefix)

			data, err := json.Marshal(map[string]interface{}{
				"BytesAddress": fmt.Sprintf("%v", addr),
				"HexAddress":   bytes.HexBytes(addr).String(),
				"EIP55Address": common.BytesToAddress(addr),
				"AccAddress":   sdk.AccAddress(addr),
				"ValAddress":   sdk.ValAddress(addr),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintString(string(data))
		},
	}
	cmd.Flags().String("addr-prefix", "fx", "custom address prefix")
	return cmd
}

func UpdateAddressPrefix(prefix string) {
	config := sdk.GetConfig()
	*config = *sdk.NewConfig()

	config.SetBech32PrefixForAccount(prefix, prefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(prefix+sdk.PrefixValidator+sdk.PrefixOperator, prefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(prefix+sdk.PrefixValidator+sdk.PrefixConsensus, prefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
	config.Seal()
}

func RawBytesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "raw-bytes [raw-bytes]",
		Short:   "Convert raw bytes output (eg. [10 21 13 255]) to hex",
		Example: fmt.Sprintf(`$ %s debug raw-bytes [72 101 108 108 111 44 32 112 108 97 121 103 114 111 117 110 100]`, version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			stringBytes := args[0]
			stringBytes = strings.Trim(stringBytes, "[")
			stringBytes = strings.Trim(stringBytes, "]")
			spl := strings.Split(stringBytes, " ")

			var byteArray []byte
			for _, s := range spl {
				b, err := strconv.ParseInt(s, 10, 8)
				if err != nil {
					return err
				}
				byteArray = append(byteArray, byte(b))
			}
			fmt.Printf("%X\n", byteArray)
			return nil
		},
	}
}

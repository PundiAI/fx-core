package cmd

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/bytes"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/bech32"
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
		ReEncodeAddrCommand(),
		HexToFxAddrCommand(),
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

func ReEncodeAddrCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "re-encode-addr [address] [prefix]",
		Short:   "Address bech32 decode",
		Example: "fxcored debug re-encode-addr fx1hajqu45kq3d0ewt7wtevhzlxgjfweja5gn7ppl px",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, bytes, err := bech32.DecodeAndConvert(args[0])
			if err != nil {
				return err
			}
			address, err := bech32.ConvertAndEncode(args[1], bytes)
			if err != nil {
				return err
			}
			cmd.Println(address)
			return nil
		},
	}
	return cmd
}

func HexToFxAddrCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hex-to-addr [hex] [prefix]",
		Short: "Hex to fx address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			hexStr := strings.TrimPrefix(args[0], "0x")
			decodeString, err := hex.DecodeString(hexStr)
			if err != nil {
				return err
			}
			if len(decodeString) != 32 {
				return fmt.Errorf("invalid hex string")
			}
			address, err := bech32.ConvertAndEncode(args[1], decodeString[len(decodeString)-20:])
			if err != nil {
				return err
			}
			cmd.Println(address)
			return nil
		},
	}
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

// getPubKeyFromString decodes SDK PubKey using JSON marshaler.
func getPubKeyFromString(ctx client.Context, pkstr string) (cryptotypes.PubKey, error) {
	var pk cryptotypes.PubKey
	err := ctx.JSONMarshaler.UnmarshalInterfaceJSON([]byte(pkstr), &pk)
	return pk, err
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
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			pk, err := getPubKeyFromString(clientCtx, args[0])
			if err != nil {
				return err
			}

			addr := pk.Address()
			cmd.Printf("Address (EIP-55): %s\n", common.BytesToAddress(addr))
			cmd.Printf("Bech32 Acc: %s\n", sdk.AccAddress(addr))
			cmd.Println("PubKey Hex:", hex.EncodeToString(pk.Bytes()))
			return nil
		},
	}
}

func AddrCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "addr [address]",
		Short: "Convert an address between hex and bech32",
		Long:  "Convert an address between hex encoding and bech32.",
		Example: fmt.Sprintf(
			`$ %s debug addr ethm10jmp6sgh4cc6zt3e8gw05wavvejgr5pw2unfju
$ %s debug addr 0xA588C66983a81e800Db4dF74564F09f91c026351`, version.AppName, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			cmd.Println("Address bytes:", addr)
			cmd.Printf("Address (hex): %s\n", bytes.HexBytes(addr))
			cmd.Printf("Address (EIP-55): %s\n", common.BytesToAddress(addr))
			cmd.Printf("Bech32 Acc: %s\n", sdk.AccAddress(addr))
			cmd.Printf("Bech32 Val: %s\n", sdk.ValAddress(addr))
			return nil
		},
	}
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

			byteArray := []byte{}
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

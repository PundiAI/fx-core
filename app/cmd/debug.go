package cmd

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

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
			hexStr := args[0]
			if strings.HasPrefix(hexStr, "0x") {
				hexStr = hexStr[2:]
			}
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
			hexStr := args[0]
			if strings.HasPrefix(hexStr, "0x") {
				hexStr = hexStr[2:]
			}
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
			clientCtx, err := client.GetClientQueryContext(cmd)
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

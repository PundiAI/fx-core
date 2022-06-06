package keys

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/app/cli"
)

const prefixFlag = "prefix"

func ParseAddressCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse [address or name]",
		Short: "Parse address from hex to bech32 and vice versa",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			addrStr := args[0]

			var addr []byte
			keyInfo, err := clientCtx.Keyring.Key(addrStr)
			if err != nil {
				// try hex, then bech32
				addr, err = hex.DecodeString(addrStr)
				if err != nil {
					_, addr, err = bech32.DecodeAndConvert(addrStr)
					if err != nil {
						return err
					}
				}
				keyInfo, _ = clientCtx.Keyring.KeyByAddress(sdk.AccAddress(addr))
			} else {
				addr = keyInfo.GetAddress().Bytes()
			}
			prefix, err := cmd.Flags().GetString(prefixFlag)
			if err != nil {
				return err
			}
			accAddress, err := bech32.ConvertAndEncode(prefix, addr)
			if err != nil {
				return err
			}
			valAddress, err := bech32.ConvertAndEncode(prefix+sdk.PrefixValidator+sdk.PrefixOperator, addr)
			if err != nil {
				return err
			}

			outputMap := map[string]interface{}{
				"base64_address": addr,
				"hex_address":    hex.EncodeToString(addr),
				"eip55_address":  common.BytesToAddress(addr).String(),
				"acc_address":    accAddress,
				"val_address":    valAddress,
			}
			if keyInfo != nil {
				outputMap["name"] = keyInfo.GetName()
				outputMap["algo"] = keyInfo.GetAlgo()
				outputMap["pubkey"] = keyInfo.GetPubKey()
				outputMap["type"] = keyInfo.GetType()
				path, err := keyInfo.GetPath()
				if err == nil {
					outputMap["path"] = path
				}
			}

			return cli.PrintOutput(clientCtx, outputMap)
		},
	}
	cmd.Flags().String(prefixFlag, "fx", "custom address prefix")
	return cmd
}

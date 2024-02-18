package keys

import (
	"encoding/hex"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/v7/client/cli"
)

const prefixFlag = "prefix"

func ParseAddressCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decode [address or name]",
		Short: "Decode address from hex to bech32 and vice versa",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			addrStr := args[0]

			outputMap := make(map[string]interface{})

			var addr []byte
			k, err := clientCtx.Keyring.Key(addrStr)
			if err != nil {
				// try hex, then bech32
				addr, err = hexutil.Decode(addrStr)
				if err != nil {
					_, addr, err = bech32.DecodeAndConvert(addrStr)
					if err != nil {
						return err
					}
				} else {
					outputMap["eip55_address"] = common.BytesToAddress(addr).String()
				}
				k, _ = clientCtx.Keyring.KeyByAddress(sdk.AccAddress(addr))
			} else {
				addr, err = k.GetAddress()
				if err != nil {
					return err
				}
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

			outputMap["base64"] = addr
			outputMap["hex"] = hex.EncodeToString(addr)
			outputMap["acc_address"] = accAddress
			outputMap["val_address"] = valAddress
			if k != nil {
				pubKey, err := k.GetPubKey()
				if err != nil {
					return err
				}
				outputMap["name"] = k.Name
				outputMap["pubkey"] = pubKey
				outputMap["type"] = k.GetType()
				ledger := k.GetLedger()
				if ledger != nil {
					outputMap["path"] = *ledger.Path
				}
				if strings.EqualFold(pubKey.Type(), ethsecp256k1.KeyType) {
					outputMap["eip55_address"] = common.BytesToAddress(addr).String()
				}
			}

			return cli.PrintOutput(clientCtx, outputMap)
		},
	}
	cmd.Flags().String(prefixFlag, "fx", "custom address prefix")
	return cmd
}

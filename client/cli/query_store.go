package cli

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
)

func QueryStoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store <store name> <hex key>",
		Short: "Query for a blockchain store",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			hexKey, err := hex.DecodeString(args[1])
			if err != nil {
				return err
			}
			data, height, err := clientCtx.QueryStore(hexKey, args[0])
			if err != nil {
				return err
			}
			return PrintOutput(clientCtx, map[string]interface{}{
				"block_height": height,
				"data":         hex.EncodeToString(data),
			})
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func QueryValidatorByConsAddr() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator [validator-consAddr]",
		Short: "Query details about an individual validator cons address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			consAddr, err := sdk.ConsAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			opAddr, _, err := clientCtx.QueryStore(types.GetValidatorByConsAddrKey(consAddr), types.StoreKey)
			if err != nil {
				return err
			}
			if opAddr == nil {
				return fmt.Errorf("not found validator by consAddress:%s", consAddr.String())
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Validator(context.Background(), &types.QueryValidatorRequest{
				ValidatorAddr: sdk.ValAddress(opAddr).String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(&res.Validator)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

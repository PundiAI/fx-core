package server

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/rosetta"
	"github.com/spf13/cobra"

	fxtypes "github.com/functionx/fx-core/v7/types"
)

// RosettaCommand builds the rosetta root command given
// a protocol buffers serializer/deserializer
func RosettaCommand(ir codectypes.InterfaceRegistry, cdc codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rosetta",
		Short: "spin up a rosetta server",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := rosetta.FromFlags(cmd.Flags())
			if err != nil {
				return err
			}
			grpcFlag := cmd.Flag(rosetta.FlagGRPCEndpoint)
			if grpcFlag.DefValue != grpcFlag.Value.String() && strings.HasPrefix(grpcFlag.Value.String(), "map") {
				conf.GRPCEndpoint = grpcFlag.DefValue
			}

			protoCodec, ok := cdc.(*codec.ProtoCodec)
			if !ok {
				return fmt.Errorf("exoected *codec.ProtoCodec, got: %T", cdc)
			}
			conf.WithCodec(ir, protoCodec)

			rosettaSrv, err := rosetta.ServerFromConfig(conf)
			if err != nil {
				return err
			}
			return rosettaSrv.Start()
		},
	}
	cmd.Flags().String(rosetta.FlagBlockchain, rosetta.DefaultBlockchain, "the blockchain type")
	cmd.Flags().String(rosetta.FlagNetwork, rosetta.DefaultNetwork, "the network name")
	cmd.Flags().String(rosetta.FlagTendermintEndpoint, rosetta.DefaultTendermintEndpoint, "the tendermint rpc endpoint, without tcp://")
	cmd.Flags().String(rosetta.FlagGRPCEndpoint, rosetta.DefaultGRPCEndpoint, "the app gRPC endpoint")
	cmd.Flags().String(rosetta.FlagAddr, rosetta.DefaultAddr, "the address rosetta will bind to")
	cmd.Flags().Int(rosetta.FlagRetries, rosetta.DefaultRetries, "the number of retries that will be done before quitting")
	cmd.Flags().Bool(rosetta.FlagOffline, rosetta.DefaultOffline, "run rosetta only with construction API")
	cmd.Flags().Bool(rosetta.FlagEnableFeeSuggestion, rosetta.DefaultEnableFeeSuggestion, "enable default fee suggestion")
	cmd.Flags().Int(rosetta.FlagGasToSuggest, flags.DefaultGasLimit, "default gas for fee suggestion")
	cmd.Flags().String(rosetta.FlagDenomToSuggest, fxtypes.GetDefaultNodeHome(), "default denom for fee suggestion")
	cmd.Flags().String(rosetta.FlagPricesToSuggest, fxtypes.GetDefGasPrice().String(), "default prices for fee suggestion")

	return cmd
}

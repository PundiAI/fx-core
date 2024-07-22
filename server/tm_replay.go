package server

import (
	tmcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/consensus"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
)

// ReplayCmd allows replaying of messages from the WAL.
func ReplayCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "replay",
		Short: "Replay messages from WAL",
		Run: func(cmd *cobra.Command, args []string) {
			serverCtx := server.GetServerContextFromCmd(cmd)
			tmcfg.EnsureRoot(serverCtx.Config.RootDir)

			consensus.RunReplayFile(serverCtx.Config.BaseConfig, serverCtx.Config.Consensus, false)
		},
	}
}

// ReplayConsoleCmd allows replaying of messages from the WAL in a
// console.
func ReplayConsoleCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "replay-console",
		Aliases: []string{"replay_console"},
		Short:   "Replay messages from WAL in a console",
		Run: func(cmd *cobra.Command, args []string) {
			serverCtx := server.GetServerContextFromCmd(cmd)
			tmcfg.EnsureRoot(serverCtx.Config.RootDir)
			consensus.RunReplayFile(serverCtx.Config.BaseConfig, serverCtx.Config.Consensus, true)
		},
	}
}

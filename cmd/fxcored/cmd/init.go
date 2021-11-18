package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"os"
	"path/filepath"

	appCmd "github.com/functionx/fx-core/app/cmd"
	"github.com/functionx/fx-core/app/fxcore"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/tendermint/tendermint/types"
	bip39 "github.com/tyler-smith/go-bip39"
)

const (
	// FlagOverwrite defines a flag to overwrite an existing genesis JSON file.
	FlagOverwrite = "overwrite"

	// FlagRecover defines a flag to initialize the private validator key from a specific seed.
	FlagRecover = "recover"

	// FlagDenom defines a flag to set the default coin denomination
	FlagDenom = "denom"
)

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.JSONMarshaler

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", tmrand.Str(6))
			}

			// Get bip39 mnemonic
			var mnemonic string
			flagRecover, err := cmd.Flags().GetBool(FlagRecover)
			if err != nil {
				return err
			}
			if flagRecover {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				mnemonic, err := input.GetString("Enter your bip39 mnemonic", inBuf)
				if err != nil {
					return err
				}

				if !bip39.IsMnemonicValid(mnemonic) {
					return errors.New("invalid mnemonic")
				}
			}

			nodeID, _, err := genutil.InitializeNodeValidatorFilesFromMnemonic(config, mnemonic)
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			genFile := config.GenesisFile()
			overwrite, _ := cmd.Flags().GetBool(FlagOverwrite)

			if !overwrite && tmos.FileExists(genFile) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}
			flagDenom, err := cmd.Flags().GetString(FlagDenom)
			if err != nil || flagDenom == "" {
				return fmt.Errorf("invalid staking denom: %v", err)
			}
			appState, err := json.MarshalIndent(fxcore.NewDefAppGenesisByDenom(flagDenom, cdc), "", " ")
			if err != nil {
				return fmt.Errorf("failed to marshall default genesis state: %s", err.Error())
			}

			genDoc := &types.GenesisDoc{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
				genDoc.ConsensusParams = CustomConsensusParams()
			} else {
				genDoc, err = types.GenesisDocFromFile(genFile)
				if err != nil {
					return fmt.Errorf("failed to read genesis doc from file: %s", err.Error())
				}
			}

			genDoc.ChainID = chainID
			genDoc.Validators = nil
			genDoc.AppState = appState
			if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
				return fmt.Errorf("failed to export gensis file: %s", err.Error())
			}

			toPrint := appCmd.NewPrintInfo(config.Moniker, chainID, nodeID, "", appState)

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

			out, err := json.MarshalIndent(toPrint, "", " ")
			if err != nil {
				return err
			}
			return clientCtx.PrintBytes(sdk.MustSortJSON(out))
		},
	}

	cmd.Flags().String(cli.HomeFlag, fxcore.DefaultNodeHome, "node's home directory")
	cmd.Flags().Bool(FlagOverwrite, false, "overwrite the genesis.json file")
	cmd.Flags().Bool(FlagRecover, false, "provide seed phrase to recover existing key instead of creating")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(FlagDenom, fxcore.MintDenom, "set the default coin denomination")
	cmd.Flags().StringP(cli.OutputFlag, "o", "json", "Output format (text|json)")
	return cmd
}

func CustomConsensusParams() *tmproto.ConsensusParams {
	result := types.DefaultConsensusParams()
	result.Block.MaxBytes = 1048576 //1M
	result.Block.MaxGas = -1
	result.Block.TimeIotaMs = 1000
	result.Evidence.MaxAgeNumBlocks = 1000000
	result.Evidence.MaxBytes = 100000
	result.Evidence.MaxAgeDuration = 172800000000000
	return result
}

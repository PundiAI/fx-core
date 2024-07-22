package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	tmcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/libs/cli"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	clientcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
	"github.com/tyler-smith/go-bip39"
)

type PrintInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

func NewPrintInfo(moniker, chainID, nodeID, genTxsDir string, appMessage json.RawMessage) PrintInfo {
	return PrintInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
//
//gocyclo:ignore
func InitCmd(nodeHome string, genesisState map[string]json.RawMessage, consensusParams *tmtypes.ConsensusParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, application and client configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			chainID, err := cmd.Flags().GetString(flags.FlagChainID)
			if err != nil {
				return err
			}

			// Get bip39 mnemonic
			var mnemonic string
			flagRecover, err := cmd.Flags().GetBool(clientcli.FlagRecover)
			if err != nil {
				return err
			}
			if flagRecover {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				mnemonic, err = input.GetString("Enter your bip39 mnemonic", inBuf)
				if err != nil {
					return err
				}

				if !bip39.IsMnemonicValid(mnemonic) {
					return errors.New("invalid mnemonic")
				}
			}

			// Get initial height
			initHeight, _ := cmd.Flags().GetInt64(flags.FlagInitHeight)
			if initHeight < 1 {
				initHeight = 1
			}

			nodeID, _, err := genutil.InitializeNodeValidatorFilesFromMnemonic(config, mnemonic)
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			genFile := config.GenesisFile()
			overwrite, _ := cmd.Flags().GetBool(clientcli.FlagOverwrite)

			if !overwrite && tmos.FileExists(genFile) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}
			appState, err := json.MarshalIndent(genesisState, "", " ")
			if err != nil {
				return fmt.Errorf("failed to marshall default genesis state: %s", err.Error())
			}

			genDoc := &tmtypes.GenesisDoc{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
				genDoc.ConsensusParams = consensusParams
			} else {
				genDoc, err = tmtypes.GenesisDocFromFile(genFile)
				if err != nil {
					return fmt.Errorf("failed to read genesis doc from file: %s", err.Error())
				}
			}

			genDoc.ChainID = chainID
			genDoc.Validators = nil
			genDoc.AppState = appState
			genDoc.InitialHeight = initHeight
			if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
				return fmt.Errorf("failed to export gensis file: %s", err.Error())
			}

			tmcfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

			toPrint := NewPrintInfo(config.Moniker, chainID, nodeID, "", appState)
			out, err := json.MarshalIndent(toPrint, "", " ")
			if err != nil {
				return err
			}
			return clientCtx.PrintString(string(sdk.MustSortJSON(out)) + "\n")
		},
	}

	cmd.Flags().String(cli.HomeFlag, nodeHome, "node's home directory")
	cmd.Flags().Bool(clientcli.FlagOverwrite, false, "overwrite the genesis.json file")
	cmd.Flags().Bool(clientcli.FlagRecover, false, "provide seed phrase to recover existing key instead of creating")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().StringP(cli.OutputFlag, "o", "json", "Output format (text|json)")
	cmd.Flags().Int64(flags.FlagInitHeight, 1, "specify the initial block height at genesis")
	return cmd
}

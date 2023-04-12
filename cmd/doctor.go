package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmcfg "github.com/tendermint/tendermint/config"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"

	fxcfg "github.com/functionx/fx-core/v3/server/config"
	fxtypes "github.com/functionx/fx-core/v3/types"
)

func doctorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Query node gas prices",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)

			printInfo()
			checkOS()
			checkFxCoreVersion()
			if err := checkGenesis(serverCtx.Config.GenesisFile()); err != nil {
				return err
			}
			if err := checkAppConfig(serverCtx.Viper); err != nil {
				return err
			}
			if err := checkTmConfig(serverCtx.Config); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func printInfo() {
	fmt.Printf(`
Please note that these warnings are just used to help the fxCore maintainers
If everything you use "fxcored" for is working fine: please don't worry; 
just ignore this. Thanks!
`)
	fmt.Println()
}

func checkOS() {
	fmt.Println("Computer Info:")
	fmt.Printf("\tOS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func checkFxCoreVersion() {
	fmt.Println("fxcored Info:")
	info := version.NewInfo()
	fmt.Println("\tVersion: ", info.Version)
	fmt.Println("\tGit Commit: ", info.GitCommit)
	fmt.Println("\tBuild Tags: ", info.BuildTags)
	fmt.Println("\tGo Version: ", info.GoVersion)
	fmt.Println("\tCosmos SDK Version: ", info.CosmosSdkVersion)
}

func checkGenesis(genesisFile string) error {
	fmt.Print("Genesis: ")
	genesisSha256, err := getGenesisSha256(genesisFile)
	if err != nil {
		return err
	}
	switch genesisSha256 {
	case fxtypes.MainnetGenesisHash:
		fmt.Println("Mainnet")
	case fxtypes.TestnetGenesisHash:
		fmt.Println("Testnet")
	default:
		fmt.Println("Unknown!")
	}
	return nil
}

func getGenesisSha256(genesisFile string) (string, error) {
	genesisFileData, err := os.ReadFile(genesisFile)
	if err != nil {
		return "", err
	}
	genesisDoc, err := tmtypes.GenesisDocFromJSON(genesisFileData)
	if err != nil {
		return "", err
	}
	genesisBytes, err := tmjson.Marshal(genesisDoc)
	if err != nil {
		return "", err
	}
	return fxtypes.Sha256Hex(genesisBytes), nil
}

func checkAppConfig(viper *viper.Viper) error {
	config.SetConfigTemplate(fxcfg.DefaultConfigTemplate())
	appConfig := fxcfg.DefaultConfig()
	if err := viper.Unmarshal(appConfig); err != nil {
		return err
	}
	return nil
}

func checkTmConfig(config *tmcfg.Config) error {
	return nil
}

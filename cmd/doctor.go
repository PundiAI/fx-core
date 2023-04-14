package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/version"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmcfg "github.com/tendermint/tendermint/config"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/functionx/fx-core/v3/client/grpc"
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			printPrompt()
			printOSInfo()
			printSelfInfo()
			network, err := checkGenesis(serverCtx.Config.GenesisFile())
			if err != nil {
				return err
			}
			upgradeInfo, err := checkUpgradeInfo(serverCtx.Config.RootDir)
			if err != nil {
				return err
			}
			needUpgrade, err := checkBlockchainData(clientCtx, network, upgradeInfo)
			if err != nil {
				return err
			}
			if err := checkAppConfig(serverCtx.Viper); err != nil {
				return err
			}
			if err := checkTmConfig(serverCtx.Config, needUpgrade); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(flags.FlagNode, "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	cmd.Flags().String(flags.FlagGRPC, "", "the gRPC endpoint to use for this chain")
	cmd.Flags().Bool(flags.FlagGRPCInsecure, false, "allow gRPC over insecure channels, if not TLS the server must use TLS")
	return cmd
}

func printPrompt() {
	fmt.Printf(`
Please note that these warnings are just used to help the fxCore maintainers
If everything you use "fxcored" for is working fine: please don't worry; 
just ignore this. Thanks!
`)
	fmt.Println()
}

func printOSInfo() {
	fmt.Println("Computer Info:")
	fmt.Printf("\tOS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func printSelfInfo() {
	fmt.Println("fxcored Info:")
	info := version.NewInfo()
	fmt.Println("\tVersion: ", info.Version)
	fmt.Println("\tGit Commit: ", info.GitCommit)
	fmt.Println("\tBuild Tags: ", info.BuildTags)
	fmt.Println("\tGo Version: ", info.GoVersion)
	fmt.Println("\tCosmos SDK Version: ", info.CosmosSdkVersion)
}

func checkGenesis(genesisFile string) (string, error) {
	fmt.Println("Genesis:")
	fmt.Println("\tFile: ", genesisFile)
	genesisSha256, err := getGenesisSha256(genesisFile)
	if err != nil {
		return "", err
	}
	switch genesisSha256 {
	case fxtypes.MainnetGenesisHash:
		fmt.Println("\tNetwork: Mainnet")
		return fxtypes.MainnetGenesisHash, nil
	case fxtypes.TestnetGenesisHash:
		fmt.Println("\tNetwork: Testnet")
		return fxtypes.TestnetGenesisHash, nil
	default:
		fmt.Println("\tNetwork: Unknown!")
		return "", nil
	}
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

func checkUpgradeInfo(homeDir string) (*upgradetypes.Plan, error) {
	fmt.Println("Upgrade Info:")
	file := filepath.Join(homeDir, "data", upgradetypes.UpgradeInfoFilename)
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("\tNot found")
			return nil, nil
		}
		return nil, fmt.Errorf("read upgrade info error: %s", err.Error())
	}
	upgradeInfo := new(upgradetypes.Plan)
	if err := json.Unmarshal(data, upgradeInfo); err != nil {
		return nil, err
	}
	fmt.Println("\tName: ", upgradeInfo.Name)
	fmt.Println("\tHeight: ", upgradeInfo.Height)
	return upgradeInfo, nil
}

func checkBlockchainData(cliCtx client.Context, network string, upgradeInfo *upgradetypes.Plan) (bool, error) {
	fmt.Println("Blockchain Data:")
	fmt.Printf("\tRemote Node: %v or %v\n", cliCtx.Viper.Get(flags.FlagGRPC), cliCtx.Viper.Get(flags.FlagNode))
	cli := grpc.NewClient(cliCtx)
	chainId, err := cli.GetChainId()
	if err != nil {
		return false, err
	}
	fmt.Println("\tChain ID: ", chainId)
	blockHeight, err := cli.GetBlockHeight()
	if err != nil {
		return false, err
	}
	fmt.Printf("\tBlock Height: %d\n", blockHeight)
	response, err := cli.TMServiceClient().GetSyncing(context.Background(), &tmservice.GetSyncingRequest{})
	if err != nil {
		return false, nil
	}
	fmt.Println("\tSyncing: ", response.Syncing)
	infoResponse, err := cli.TMServiceClient().GetNodeInfo(context.Background(), &tmservice.GetNodeInfoRequest{})
	if err != nil {
		return false, nil
	}
	fmt.Println("\tNode Info: ", response.Syncing)
	fmt.Println("\t\tVersion: ", infoResponse.ApplicationVersion.Version)
	fmt.Println("\t\tGit Commit: ", infoResponse.ApplicationVersion.GitCommit)
	fmt.Println("\t\tBuild Tags: ", infoResponse.ApplicationVersion.BuildTags)
	fmt.Println("\t\tGo Version: ", infoResponse.ApplicationVersion.GoVersion)
	fmt.Println("\t\tCosmos SDK Version: ", infoResponse.ApplicationVersion.CosmosSdkVersion)
	resp, err := cli.UpgradeQuery().CurrentPlan(context.Background(), &upgradetypes.QueryCurrentPlanRequest{})
	if err != nil {
		return false, err
	}
	if resp.Plan != nil && !resp.Plan.Equal(upgradeInfo) {
		fmt.Println("\tUpgrade Plan:")
		fmt.Println("\t\tName: ", resp.Plan.Name)
		fmt.Println("\t\tHeight: ", resp.Plan.Height)
	}
	if chainId != network {
		fmt.Printf("\tWarn: The remote node chainId(%s) does not match the local genesis chainId(%s)\n", chainId, network)
		return false, nil
	}
	if network == fxtypes.MainnetChainId {
		if blockHeight < fxtypes.MainnetBlockHeightV2 {
			fmt.Println("Version: V1")
		} else if blockHeight < fxtypes.MainnetBlockHeightV3 {
			fmt.Println("Version: V2")
		}
		fmt.Println("Version: V3")
	}
	if network == fxtypes.TestnetChainId {
		switch blockHeight {
		case fxtypes.TestnetBlockHeightV2:
			fmt.Println("Version: V1")
		case fxtypes.TestnetBlockHeightV3:
			fmt.Println("Version: V2")
		}
		fmt.Println("Version: V3")
	}
	return resp.Plan != nil && response.Syncing, nil
}

func checkAppConfig(viper *viper.Viper) error {
	fmt.Println("App Config:")
	fmt.Println("\tFile: ", viper.ConfigFileUsed())
	config.SetConfigTemplate(fxcfg.DefaultConfigTemplate())
	appConfig := fxcfg.DefaultConfig()
	if err := viper.Unmarshal(appConfig); err != nil {
		return err
	}
	if err := appConfig.ValidateBasic(); err != nil {
		fmt.Println("\tWarn: ", err.Error())
	}
	return nil
}

func checkTmConfig(config *tmcfg.Config, needUpgrade bool) error {
	fmt.Println("Tendermint Config: ")
	fmt.Printf("\tFile: %s/config/config.toml\n", config.RootDir)
	if err := config.ValidateBasic(); err != nil {
		fmt.Println("\tWarn: ", err.Error())
	}
	if needUpgrade && config.Consensus.DoubleSignCheckHeight > 0 {
		fmt.Println("Warn: double_sign_check_height is greater than 0")
	}
	return nil
}

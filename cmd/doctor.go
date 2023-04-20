package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	cosmosserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmcfg "github.com/tendermint/tendermint/config"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/privval"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/functionx/fx-core/v4/client/grpc"
	"github.com/functionx/fx-core/v4/server"
	fxcfg "github.com/functionx/fx-core/v4/server/config"
	fxtypes "github.com/functionx/fx-core/v4/types"
)

func doctorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check your system for potential problems",
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := cosmosserver.GetServerContextFromCmd(cmd)
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
			bc, err := getBlockchain(clientCtx, serverCtx)
			if err != nil {
				return err
			}
			needUpgrade, err := checkBlockchainData(bc, network, serverCtx.Config.PrivValidatorKeyFile(), upgradeInfo)
			if err != nil {
				return err
			}
			if err := checkAppConfig(serverCtx.Viper); err != nil {
				return err
			}
			if err := checkTmConfig(serverCtx.Config, needUpgrade); err != nil {
				return err
			}
			if err := checkCosmovisor(serverCtx.Config.RootDir, bc); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(flags.FlagHome, fxtypes.GetDefaultNodeHome(), "The application home directory")
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
	fmt.Println("\tCPU: ", runtime.NumCPU())
	memory, err := mem.VirtualMemory()
	if err != nil {
		return
	}
	fmt.Printf("\tMemory Total: %v MB, Available: %v MB, UsedPercent: %f%%\n",
		memory.Total/1024/1024, memory.Available/1024/1024, memory.UsedPercent)
}

func printSelfInfo() {
	fmt.Println("fxcored Info:")
	info := version.NewInfo()
	fmt.Println("\tVersion: ", info.Version)
	fmt.Println("\tGit Commit: ", info.GitCommit)
	fmt.Println("\tBuild Tags: ", info.BuildTags)
	fmt.Println("\tGo Version: ", runtime.Version())
	fmt.Println("\tCosmos SDK Version: ", info.CosmosSdkVersion)
}

func checkGenesis(genesisFile string) (string, error) {
	fmt.Println("Genesis:")
	fmt.Println("\tFile: ", genesisFile)

	genesisSha256, err := getGenesisSha256(genesisFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("\tWarning: Not found Genesis file!")
			return "", nil
		}
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
		fmt.Println("\tWarning: Unknown Network!")
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
			fmt.Println("\tNot found file: ", file)
			return nil, nil
		}
		return nil, fmt.Errorf("read upgrade info error: %s", err.Error())
	}
	upgradeInfo := new(upgradetypes.Plan)
	if err := json.Unmarshal(data, upgradeInfo); err != nil {
		return nil, err
	}
	fmt.Println("\tFile: ", file)
	fmt.Println("\tName: ", upgradeInfo.Name)
	fmt.Println("\tHeight: ", upgradeInfo.Height)
	return upgradeInfo, nil
}

type blockchain interface {
	GetChainId() (string, error)
	GetBlockHeight() (int64, error)
	GetSyncing() (bool, error)
	GetNodeInfo() (*tmservice.VersionInfo, error)
	CurrentPlan() (*upgradetypes.Plan, error)
	GetValidators() ([]stakingtypes.Validator, error)
}

func getBlockchain(cliCtx client.Context, serverCtx *cosmosserver.Context) (blockchain, error) {
	fmt.Println("Blockchain Data:")
	newClient := grpc.NewClient(cliCtx)
	_, err := newClient.GetBlockHeight()
	if err == nil {
		fmt.Printf("\tRemote Node: %v or %v\n", serverCtx.Viper.Get(flags.FlagGRPC), serverCtx.Viper.Get(flags.FlagNode))
		return newClient, nil
	} else {
		if len(serverCtx.Config.RootDir) <= 0 {
			fmt.Println("\tNot found root dir")
			return nil, nil
		}
		database, err := server.NewDatabase(serverCtx.Config.RootDir, serverCtx.Config.DBBackend, cliCtx.Codec)
		if err != nil || database == nil {
			return nil, nil
		}
		fmt.Printf("\tData Dir: %s/data\n", serverCtx.Config.RootDir)
		return database, nil
	}
}

//gocyclo:ignore
func checkBlockchainData(n blockchain, network, privValidatorKeyFile string, upgradeInfo *upgradetypes.Plan) (bool, error) {
	if n == nil {
		return false, nil
	}
	chainId, err := n.GetChainId()
	if err != nil {
		return false, err
	}
	fmt.Println("\tChain ID: ", chainId)
	blockHeight, err := n.GetBlockHeight()
	if err != nil {
		return false, err
	}
	fmt.Printf("\tBlock Height: %d\n", blockHeight)
	syncing, err := n.GetSyncing()
	if err != nil {
		return false, nil
	}
	fmt.Println("\tSyncing: ", syncing)
	pvKey := privval.FilePVKey{}
	keyJSONBytes, err := os.ReadFile(privValidatorKeyFile)
	if err != nil {
		return false, err
	}
	err = tmjson.Unmarshal(keyJSONBytes, &pvKey)
	if err != nil {
		return false, err
	}
	validators, err := n.GetValidators()
	if err != nil {
		return false, err
	}
	for _, validator := range validators {
		if strings.EqualFold(
			sdk.ValAddress(pvKey.Address.Bytes()).String(),
			validator.GetOperator().String(),
		) {
			fmt.Println("\tNode Type: This node is a validator")
		}
	}
	app, err := n.GetNodeInfo()
	if err != nil {
		return false, nil
	}
	fmt.Println("\tNode Info: ")
	if len(app.Version) > 0 {
		fmt.Println("\t\tVersion: ", app.Version)
	}
	if len(app.GitCommit) > 0 {
		fmt.Println("\t\tGit Commit: ", app.GitCommit)
	}
	if len(app.BuildTags) > 0 {
		fmt.Println("\t\tBuild Tags: ", app.BuildTags)
	}
	if len(app.GoVersion) > 0 {
		fmt.Println("\t\tGo Version: ", app.GoVersion)
	}
	if len(app.CosmosSdkVersion) > 0 {
		fmt.Println("\t\tCosmos SDK Version: ", app.CosmosSdkVersion)
	}
	plan, err := n.CurrentPlan()
	if err != nil {
		return false, err
	}
	if plan != nil && !plan.Equal(upgradeInfo) {
		fmt.Println("\tUpgrade Plan:")
		fmt.Println("\t\tName: ", plan.Name)
		fmt.Println("\t\tHeight: ", plan.Height)
	}
	if chainId != network {
		fmt.Printf("\tWarning: The remote node chainId(%s) does not match the local genesis chainId(%s)\n", chainId, network)
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
	return plan != nil && syncing, nil
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
		fmt.Println("\tWarning: ", err.Error())
	}
	return nil
}

func checkTmConfig(config *tmcfg.Config, needUpgrade bool) error {
	fmt.Println("Tendermint Config: ")
	fmt.Printf("\tFile: %s/config/config.toml\n", config.RootDir)
	if err := config.ValidateBasic(); err != nil {
		fmt.Println("\tWarning: ", err.Error())
	}
	if needUpgrade && config.Consensus.DoubleSignCheckHeight > 0 {
		fmt.Println("Warning: double_sign_check_height is greater than 0")
		fmt.Println("Warning: double_sign_check_height is greater than 0")
	}
	if config.P2P.Seeds == "" {
		fmt.Println("\tWarning: seeds is empty")
	}
	return nil
}

func checkCosmovisor(rootPath string, n blockchain) error {
	cosmovisorPath := filepath.Join(rootPath, "cosmovisor")
	exist, isDir, err := checkDirFile(cosmovisorPath)
	if err != nil {
		return err
	}
	if !exist {
		fmt.Println("Cosmovisor: Not installed")
		return nil
	}
	if !isDir {
		fmt.Println("Cosmovisor: Not directory")
		return nil
	}
	fmt.Println("Cosmovisor:")
	fmt.Println("\tPath:", cosmovisorPath)

	defer func() {
		fmt.Println("\tList Path:")
		if err := printDirectory(cosmovisorPath, 0, []bool{false}, "\t\t"); err != nil {
			fmt.Println("\t", err.Error())
		}
	}()

	current := filepath.Join(cosmovisorPath, "current")
	checkCosmovisorUpgrade(current, "")

	genesis := filepath.Join(cosmovisorPath, "genesis")
	checkCosmovisorUpgrade(genesis, "")

	upgradePath := filepath.Join(cosmovisorPath, "upgrades")
	exist, isDir, err = checkDirFile(upgradePath)
	upgradePlans := make([]*upgradetypes.Plan, 0)
	if err == nil && exist && isDir {
		fmt.Println("\tUpgrades:", upgradePath)
		entries, err := os.ReadDir(upgradePath)
		if err == nil {
			if len(entries) > 0 {
				for _, entry := range entries {
					if entry.IsDir() {
						upgrade := filepath.Join(upgradePath, entry.Name())
						upgradeInfo := checkCosmovisorUpgrade(upgrade, "\t")
						if upgradeInfo != nil {
							upgradePlans = append(upgradePlans, upgradeInfo)
						}
					} else {
						fmt.Printf("\tUpgrades %s: is not directory\n", entry.Name())
					}
				}
			} else {
				fmt.Println("\t\tWarning: Not exist upgrade version")
			}
		} else {
			fmt.Println("\tUpgrades: ", err.Error())
		}
	} else {
		errDir(exist, isDir, err, "Upgrades", "")
	}

	return checkCosmovisorCurrentVersion(upgradePlans, n)
}

// checkCosmovisorUpgrade check cosmovisor upgrade plan
//
//gocyclo:ignore
func checkCosmovisorUpgrade(path, t string) *upgradetypes.Plan {
	basePath := filepath.Base(path)
	title := "Version " + basePath
	if basePath == "genesis" {
		title = "Genesis"
	} else if basePath == "current" {
		title = "Current"
	}
	exist, isDir, err := checkDirFile(path)
	if err == nil && exist && isDir {
		fmt.Printf("%s\t%s: %s\n", t, title, path)
		bin := filepath.Join(path, "bin")
		exist, isDir, err = checkDirFile(bin)
		if err == nil && exist && isDir {
			binFxcored := filepath.Join(bin, "fxcored")
			exist, isDir, err = checkDirFile(binFxcored)
			if err == nil && exist && !isDir {
				output, err := exec.Command(binFxcored, "version").Output()
				if err != nil {
					fmt.Printf("%s\t\tfxcored exec: %s\n", t, err.Error())
				} else {
					fmt.Printf("%s\t\tfxcored version: %s\n", t, bytes.Trim(output, "\n"))
				}
			} else {
				errFile(exist, isDir, err, "fxcored", t+"\t")
			}
		} else {
			errDir(exist, isDir, err, "bin", t+"\t")
		}

		if basePath != "genesis" && basePath != "current" {
			upgradeFile := filepath.Join(path, upgradetypes.UpgradeInfoFilename)
			exist, isDir, err = checkDirFile(upgradeFile)
			if err == nil && exist && !isDir {
				fmt.Printf("%s\t\t%s: %s\n", t, upgradetypes.UpgradeInfoFilename, upgradeFile)
				upgradeInfo, err := os.ReadFile(upgradeFile)
				if err != nil {
					fmt.Printf("%s\t\t%s read: %s\n", t, upgradetypes.UpgradeInfoFilename, err.Error())
					return nil
				}
				fmt.Printf("%s\t\t%s content: %s\n", t, upgradetypes.UpgradeInfoFilename, bytes.Trim(upgradeInfo, "\n"))
				var plan upgradetypes.Plan
				if err := json.Unmarshal(upgradeInfo, &plan); err != nil {
					fmt.Printf("%s\t\t%s error: %s\n", t, upgradetypes.UpgradeInfoFilename, err.Error())
					return nil
				}
				if plan.Name != basePath {
					fmt.Printf("%s\t\tupgrade plan not match: %s\n", t, plan.Name)
					return nil
				}
				return &plan
			}
			errDir(exist, isDir, err, upgradetypes.UpgradeInfoFilename, t+"\t")
		}
		return nil
	}
	errDir(exist, isDir, err, title, t)
	return nil
}

func checkCosmovisorCurrentVersion(upgradePlans []*upgradetypes.Plan, n blockchain) error {
	sort.SliceStable(upgradePlans, func(i, j int) bool {
		return upgradePlans[i].Height < upgradePlans[j].Height
	})
	if len(upgradePlans) == 0 {
		return nil
	}
	height, err := n.GetBlockHeight()
	if err != nil {
		return err
	}
	var currentPlan *upgradetypes.Plan
	for _, info := range upgradePlans {
		if height <= info.Height {
			break
		}
		currentPlan = info
	}

	if currentPlan == nil {
		// genesis
		fmt.Println("\tCurrent plan: genesis")
	} else {
		fmt.Println("\tCurrent plan:", currentPlan.Name)
	}

	return nil
}

func checkDirFile(path string) (exist, dir bool, err error) {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}
	return true, stat.IsDir(), nil
}

func errDir(exist, isDir bool, err error, title, t string) {
	if err != nil {
		fmt.Printf("%s\t%s: %s\n", title, t, err.Error())
	} else if !exist {
		fmt.Printf("%s\t%s: Not exist!\n", t, title)
	} else if !isDir {
		fmt.Printf("%s\t%s: Not directory!\n", t, title)
	}
}

func errFile(exist, isDir bool, err error, title, t string) {
	if err != nil {
		fmt.Printf("%s\t%s: %s\n", t, title, err.Error())
	} else if !exist {
		fmt.Printf("%s\t%s: Not exist!\n", t, title)
	} else if isDir {
		fmt.Printf("%s\t%s: Is directory!\n", t, title)
	}
}

func printDirectory(path string, depth int, last []bool, t string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	printPath := path
	if depth > 0 {
		printPath = filepath.Base(path)
	}
	printListing(printPath, depth, last, t)
	for idx, entry := range entries {
		currentLast := idx == len(entries)-1
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if (info.Mode() & os.ModeSymlink) == os.ModeSymlink {
			fullPath, err := os.Readlink(filepath.Join(path, entry.Name()))
			if err != nil {
				return err
			} else {
				printListing(entry.Name()+" -> "+fullPath, depth+1, append(last, currentLast), t)
			}
		} else if entry.IsDir() {
			if err = printDirectory(filepath.Join(path, entry.Name()), depth+1, append(last, currentLast), t); err != nil {
				return err
			}
		} else {
			printListing(entry.Name(), depth+1, append(last, currentLast), t)
		}
	}
	return nil
}

func printListing(entry string, depth int, last []bool, t string) {
	if depth == 0 {
		fmt.Printf("%s%s\n", t, entry)
	} else {
		indent := ""
		newLast := last[1:]
		for i := 0; i < len(newLast)-1; i++ {
			if newLast[i] {
				indent = fmt.Sprintf("%s    ", indent)
			} else {
				indent = fmt.Sprintf("%s│   ", indent)
			}
		}
		sepStr := "├── "
		if last[len(last)-1] {
			sepStr = "└── "
		}
		fmt.Printf("%s%s%s%s\n", t, indent, sepStr, entry)
	}
}

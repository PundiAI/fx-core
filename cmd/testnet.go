package cmd

import (
	"fmt"
	"html/template"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"

	"github.com/functionx/fx-core/v6/app"
	"github.com/functionx/fx-core/v6/testutil"
	"github.com/functionx/fx-core/v6/testutil/network"
	fxtypes "github.com/functionx/fx-core/v6/types"
)

const (
	flagOutputDir   = "output-dir"
	flagDockerImage = "docker-image"
)

// testnetCmd get cmd to initialize all files for tendermint testnet and application
func testnetCmd(encCfg app.EncodingConfig) *cobra.Command {
	networkConfig := testutil.DefaultNetworkConfig(encCfg)
	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a fxcore local testnet",
		Long: `testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:
	$ fxcored testnet -v 4 -output-dir ./testnet
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			outputDir, _ := cmd.Flags().GetString(flagOutputDir)
			stakingTokens, _ := cmd.Flags().GetInt64("staking-tokens")
			networkConfig.StakingTokens = sdkmath.NewIntWithDecimal(stakingTokens, fxtypes.DenomUnit)
			bondedTokens, _ := cmd.Flags().GetInt64("bonded-tokens")
			networkConfig.BondedTokens = sdkmath.NewIntWithDecimal(bondedTokens, fxtypes.DenomUnit)
			validators, err := network.GenerateGenesisAndValidators(outputDir, &networkConfig)
			if err != nil {
				return err
			}
			cmd.Println(fmt.Sprintf("Successfully initialized %d node directories", networkConfig.NumValidators))

			dockerImage, _ := cmd.Flags().GetString(flagDockerImage)
			if len(dockerImage) == 0 {
				return nil
			}
			if err = generateDockerComposeYml(validators, outputDir, dockerImage); err != nil {
				return err
			}
			cmd.Println("Please run: docker-compose up -d")
			return nil
		},
	}

	cmd.Flags().String(flagOutputDir, "./testnet", "Directory to store initialization data for the testnet")
	cmd.Flags().String(flagDockerImage, "", "set docker run image")
	cmd.Flags().IntVarP(&networkConfig.NumValidators, "validator-number", "v", 4, "Number of validators to initialize the testnet with")
	cmd.Flags().StringSliceVar(&networkConfig.Mnemonics, "mnemonics", []string{}, "Mnemonics for the validator accounts")
	cmd.Flags().StringVar(&networkConfig.BondDenom, "bond-denom", networkConfig.BondDenom, "Coin denomination used for bonds")
	cmd.Flags().StringVar(&networkConfig.RPCAddress, "rpc-address", "", "RPC listen address (including port)")
	cmd.Flags().StringVar(&networkConfig.APIAddress, "api-address", "", "REST API listen address (including port)")
	cmd.Flags().StringVar(&networkConfig.GRPCAddress, "grpc-address", "", "GRPC server listen address (including port)")
	cmd.Flags().StringVar(&networkConfig.JSONRPCAddress, "jsonrpc-address", "", "JSON-RPC listen address (including port)")
	cmd.Flags().StringVar(&networkConfig.ChainID, "chain-id", networkConfig.ChainID, "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().StringVar(&networkConfig.MinGasPrices, "min-gas-prices", networkConfig.MinGasPrices, "Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum")
	cmd.Flags().DurationVar(&networkConfig.TimeoutCommit, "timeout-commit", 5*time.Second, "The amount of time to wait before a commit is returned as a timeout")
	cmd.Flags().Int64("staking-tokens", 5000, "the amount of tokens each validator has available to stake")
	cmd.Flags().Int64("bonded-tokens", 100, "the amount of tokens each validator stakes")
	return cmd
}

func generateDockerComposeYml(validators []*network.Validator, outputDir, dockerImage string) error {
	chainId := validators[0].ClientCtx.ChainID
	IPAddress := "172.20.0.2"
	data := map[string]interface{}{
		"Subnet": fmt.Sprintf("%s/16", getSubnet(IPAddress)),
	}
	persistentPeers := make([]string, 0, len(validators))
	services := make([]map[string]interface{}, 0)
	for i, validator := range validators {
		ip, err := getNextIP(i, IPAddress)
		if err != nil {
			return err
		}
		nodeName := validator.Ctx.Config.Moniker
		nodeDir := filepath.Join(outputDir, nodeName, strings.ToLower(chainId))

		validator.AppConfig.API.Address = fmt.Sprintf("tcp://%s:1317", ip)
		validator.AppConfig.GRPC.Address = fmt.Sprintf("%s:9090", ip)
		validator.AppConfig.JSONRPC.Address = fmt.Sprintf("%s:8545", ip)
		validator.AppConfig.JSONRPC.WsAddress = fmt.Sprintf("%s:8546", ip)
		config.WriteConfigFile(filepath.Join(nodeDir, "config/app.toml"), validator.AppConfig)

		validator.Ctx.Config.DBBackend = "goleveldb"
		validator.Ctx.Config.P2P.PersistentPeers = strings.Join(persistentPeers, ",")
		validator.Ctx.Config.P2P.ListenAddress = fmt.Sprintf("%s:26656", ip)
		validator.Ctx.Config.P2P.ExternalAddress = fmt.Sprintf("%s:26656", ip)
		validator.Ctx.Config.RPC.ListenAddress = fmt.Sprintf("tcp://%s:26657", ip)
		tmcfg.WriteConfigFile(filepath.Join(nodeDir, "config/config.toml"), validator.Ctx.Config)

		persistentPeers = append(persistentPeers, fmt.Sprintf("%s@%s", validator.NodeID, validator.Ctx.Config.P2P.ExternalAddress))

		ports := []string{fmt.Sprintf("%d:26657", 26657+(i*2))}
		if i == 0 {
			ports = append(ports, "1317:1317")
			ports = append(ports, "9090:9090")
			ports = append(ports, "8545:8545")
			ports = append(ports, "8546:8546")
		}
		services = append(services, map[string]interface{}{
			"Image":         dockerImage,
			"ContainerName": nodeName,
			"IPv4Address":   ip,
			"Volumes":       fmt.Sprintf("%s/%s/%s:/root/.%s", outputDir, nodeName, chainId, chainId),
			"Ports":         ports,
		})
	}
	data["Services"] = services

	const fxChainDockerComposeYmlTemplate = `version: '3'

services:
  {{range .Services}}
  {{.ContainerName}}:
    container_name: {{.ContainerName}}
    image: {{.Image}}
    command: start
    ports:{{range .Ports}}
      - "{{.}}"{{end}}
    volumes:
      - {{.Volumes}}
    networks:
      chain-net:
        ipv4_address: {{.IPv4Address}}
  {{end}}
networks:
  chain-net:
    driver: bridge
    ipam:
      driver: default
      config:
      - subnet: {{.Subnet}}
`
	tmpl, err := template.New("").Parse(fxChainDockerComposeYmlTemplate)
	if err != nil {
		return err
	}
	f, err := os.OpenFile("docker-compose.yml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, data)
}

func getNextIP(i int, ip string) (string, error) {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		return "", fmt.Errorf("%v: non ipv4 address", ip)
	}

	for j := 0; j < i; j++ {
		ipv4[3]++
	}

	return ipv4.String(), nil
}

func getSubnet(ip string) string {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		panic(fmt.Errorf("%v: non ipv4 address", ip))
	}
	ipv4[3] = 0
	return ipv4.String()
}

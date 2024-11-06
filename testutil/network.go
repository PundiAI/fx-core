package testutil

import (
	"time"

	sdkmath "cosmossdk.io/math"
	pruningtypes "cosmossdk.io/store/pruning/types"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibccoretypes "github.com/cosmos/ibc-go/v8/modules/core/types"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	hd2 "github.com/evmos/ethermint/crypto/hd"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v8/app"
	fxcfg "github.com/functionx/fx-core/v8/server/config"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/testutil/network"
	fxtypes "github.com/functionx/fx-core/v8/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
)

// DefaultNetworkConfig returns a sane default configuration suitable for nearly all
// testing requirements.
func DefaultNetworkConfig(opts ...func(config *network.Config)) network.Config {
	newApp := helpers.NewApp()
	chainID := fxtypes.MainnetChainId
	cfg := network.Config{
		Codec:             newApp.AppCodec(),
		InterfaceRegistry: newApp.InterfaceRegistry(),
		TxConfig:          newApp.GetTxConfig(),
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor: func(appConfig *fxcfg.Config, ctx *server.Context) servertypes.Application {
			return app.New(
				ctx.Logger,
				dbm.NewMemDB(),
				nil,
				true,
				make(map[int64]bool),
				ctx.Config.RootDir,
				ctx.Viper,
				baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(appConfig.Pruning)),
				baseapp.SetMinGasPrices(appConfig.MinGasPrices),
				baseapp.SetChainID(chainID),
			)
		},
		GenesisState:    NoSupplyGenesisState(newApp.AppCodec(), newApp.ModuleBasics),
		TimeoutCommit:   500 * time.Millisecond,
		StakingTokens:   sdk.TokensFromConsensusPower(5000, sdk.DefaultPowerReduction), // 500_000
		BondedTokens:    sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction),  // 10_000
		NumValidators:   4,
		ChainID:         chainID,
		BondDenom:       fxtypes.DefaultDenom,
		MinGasPrices:    fxtypes.GetDefGasPrice().String(),
		PruningStrategy: pruningtypes.PruningOptionNothing,
		// RPCAddress:      "tcp://localhost:26657",
		// JSONRPCAddress:  "localhost:8545",
		// APIAddress:      "localhost:1317",
		// GRPCAddress:     "localhost:9090",
		EnableJSONRPC:   false,
		EnableAPI:       false,
		EnableTMLogging: false,
		CleanupDir:      true,
		SigningAlgo:     string(hd.Secp256k1Type),
		KeyringOptions: []keyring.Option{
			hd2.EthSecp256k1Option(),
		},
		BypassMinFeeMsgTypes: []string{
			sdk.MsgTypeURL(&distributiontypes.MsgSetWithdrawAddress{}),
		},
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func NoSupplyGenesisState(cdc codec.JSONCodec, moduleBasics module.BasicManager) app.GenesisState {
	genesisState := app.NewDefAppGenesisByDenom(cdc, moduleBasics)

	// reset supply
	bankState := banktypes.DefaultGenesisState()
	bankState.DenomMetadata = []banktypes.Metadata{fxtypes.NewFXMetaData()}
	genesisState[banktypes.ModuleName] = cdc.MustMarshalJSON(bankState)

	var govGenState govv1.GenesisState
	cdc.MustUnmarshalJSON(genesisState[govtypes.ModuleName], &govGenState)
	votingPeriod := time.Millisecond
	govGenState.Params.VotingPeriod = &votingPeriod

	genesisState[govtypes.ModuleName] = cdc.MustMarshalJSON(&govGenState)

	var evmGenState evmtypes.GenesisState
	cdc.MustUnmarshalJSON(genesisState[evmtypes.ModuleName], &evmGenState)
	evmGenState.Params.EvmDenom = fxtypes.DefaultDenom
	genesisState[evmtypes.ModuleName] = cdc.MustMarshalJSON(&evmGenState)

	return genesisState
}

func IbcGenesisState(cdc codec.Codec, genesisState app.GenesisState) app.GenesisState {
	clientState := clienttypes.DefaultGenesisState()

	// src chain 07-tendermint-0 connection-0 channel-0
	// dst chain 07-tendermint-0 connection-0 channel-0
	clientId := "07-tendermint-0"
	connectionId := connectiontypes.FormatConnectionIdentifier(0)
	channelId := "channel-0"

	// 1. ibc client state
	clientState.Clients = []clienttypes.IdentifiedClientState{
		clienttypes.NewIdentifiedClientState(clientId, &ibctm.ClientState{
			ChainId:      tmrand.Str(10),
			LatestHeight: clienttypes.NewHeight(0, 1),
			// if ibc timeout  ctx.BlockTime() > TrustingPeriod + clientState.ClientsConsensus.Timestamp
			TrustingPeriod: 8 * time.Minute,
		}),
	}

	clientState.ClientsConsensus = []clienttypes.ClientConsensusStates{
		{
			ClientId: clientId,
			ConsensusStates: []clienttypes.ConsensusStateWithHeight{
				clienttypes.NewConsensusStateWithHeight(clienttypes.NewHeight(0, 1),
					ibctm.NewConsensusState(time.Now(), types.NewMerkleRoot(tmrand.Bytes(32)), nil),
				),
			},
		},
	}

	counterparty := connectiontypes.NewCounterparty(clientId, connectionId, types.NewMerklePrefix(tmrand.Bytes(32)))
	// 2. ibc connection state
	connState := connectiontypes.DefaultGenesisState()
	connState.Connections = []connectiontypes.IdentifiedConnection{
		connectiontypes.NewIdentifiedConnection(connectionId, connectiontypes.NewConnectionEnd(
			connectiontypes.OPEN, clientId, counterparty,
			[]*connectiontypes.Version{connectiontypes.DefaultIBCVersion}, 0)),
	}

	// 3. ibc channel state
	channelState := channeltypes.DefaultGenesisState()

	channelState.Channels = []channeltypes.IdentifiedChannel{
		channeltypes.NewIdentifiedChannel(transfertypes.PortID, channelId,
			channeltypes.NewChannel(channeltypes.OPEN, channeltypes.UNORDERED,
				channeltypes.NewCounterparty(transfertypes.PortID, channelId),
				[]string{connectionId}, connectiontypes.DefaultIBCVersionIdentifier)),
	}
	channelState.SendSequences = []channeltypes.PacketSequence{channeltypes.NewPacketSequence(transfertypes.PortID, channelId, 1)}

	// for ibc test
	capabilityState := buildCapabilityGenesisState()
	genesisState[capabilitytypes.ModuleName] = cdc.MustMarshalJSON(&capabilityState)
	genesisState[exported.ModuleName] = cdc.MustMarshalJSON(&ibccoretypes.GenesisState{
		ClientGenesis:     clientState,
		ConnectionGenesis: connState,
		ChannelGenesis:    channelState,
	})
	return genesisState
}

func BankGenesisState(cdc codec.Codec, genesisState app.GenesisState) app.GenesisState {
	bankState := banktypes.DefaultGenesisState()
	coins := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1e8).Mul(sdkmath.NewInt(1e18))))
	bankState.Balances = append(bankState.Balances, banktypes.Balance{Address: authtypes.NewModuleAddress(ethtypes.ModuleName).String(), Coins: coins})
	genesisState[banktypes.ModuleName] = cdc.MustMarshalJSON(bankState)
	return genesisState
}

func GovGenesisState(cdc codec.Codec, genesisState app.GenesisState, votingPeriod time.Duration) app.GenesisState {
	var govGenState govv1.GenesisState
	cdc.MustUnmarshalJSON(genesisState[govtypes.ModuleName], &govGenState)
	govGenState.Params.VotingPeriod = &votingPeriod

	genesisState[govtypes.ModuleName] = cdc.MustMarshalJSON(&govGenState)
	return genesisState
}

func SlashingGenesisState(cdc codec.Codec, genesisState app.GenesisState, signedBlocksWindow int64, minSignedPerWindow sdkmath.LegacyDec, downtimeJailDuration time.Duration) app.GenesisState {
	var slashingState slashingtypes.GenesisState
	cdc.MustUnmarshalJSON(genesisState[slashingtypes.ModuleName], &slashingState)
	slashingState.Params.SignedBlocksWindow = signedBlocksWindow
	slashingState.Params.MinSignedPerWindow = minSignedPerWindow
	slashingState.Params.DowntimeJailDuration = downtimeJailDuration

	genesisState[slashingtypes.ModuleName] = cdc.MustMarshalJSON(&slashingState)
	return genesisState
}

func buildCapabilityGenesisState() capabilitytypes.GenesisState {
	capabilityState := capabilitytypes.GenesisState{}
	capabilityState.Index = 3
	capabilityState.Owners = []capabilitytypes.GenesisOwners{
		{
			Index: 1,
			IndexOwners: capabilitytypes.CapabilityOwners{
				Owners: []capabilitytypes.Owner{
					capabilitytypes.NewOwner("ibc", "ports/transfer"),
					capabilitytypes.NewOwner("transfer", "ports/transfer"),
				},
			},
		},
		{
			Index: 2,
			IndexOwners: capabilitytypes.CapabilityOwners{
				Owners: []capabilitytypes.Owner{
					capabilitytypes.NewOwner("ibc", "capabilities/ports/transfer/channels/channel-0"),
					capabilitytypes.NewOwner("transfer", "capabilities/ports/transfer/channels/channel-0"),
				},
			},
		},
	}
	return capabilityState
}

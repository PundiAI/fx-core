package testutil

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v3/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v3/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibccoretypes "github.com/cosmos/ibc-go/v3/modules/core/types"
	ibctmtypes "github.com/cosmos/ibc-go/v3/modules/light-clients/07-tendermint/types"
	hd2 "github.com/evmos/ethermint/crypto/hd"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/tendermint/tendermint/libs/rand"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/testutil/network"
	fxtypes "github.com/functionx/fx-core/v3/types"
)

// DefaultNetworkConfig returns a sane default configuration suitable for nearly all
// testing requirements.
func DefaultNetworkConfig(encCfg app.EncodingConfig, opts ...func(config *network.Config)) network.Config {
	fxtypes.SetConfig(true)
	cfg := network.Config{
		Codec:             encCfg.Codec,
		TxConfig:          encCfg.TxConfig,
		LegacyAmino:       encCfg.Amino,
		InterfaceRegistry: encCfg.InterfaceRegistry,
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor: func(val network.Validator) servertypes.Application {
			return app.New(val.Ctx.Logger, dbm.NewMemDB(),
				nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
				encCfg,
				app.EmptyAppOptions{},
				baseapp.SetPruning(storetypes.NewPruningOptionsFromString(val.AppConfig.Pruning)),
				baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
			)
		},
		GenesisState:    NoSupplyGenesisState(encCfg.Codec),
		TimeoutCommit:   500 * time.Millisecond,
		ChainID:         fxtypes.MainnetChainId,
		NumValidators:   4,
		BondDenom:       fxtypes.DefaultDenom,
		MinGasPrices:    fmt.Sprintf("4000000000000%s", fxtypes.DefaultDenom),
		AccountTokens:   sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction),
		StakingTokens:   sdk.TokensFromConsensusPower(5000, sdk.DefaultPowerReduction),
		BondedTokens:    sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction),
		PruningStrategy: storetypes.PruningOptionNothing,
		CleanupDir:      true,
		SigningAlgo:     string(hd.Secp256k1Type),
		KeyringOptions: []keyring.Option{
			hd2.EthSecp256k1Option(),
		},
		PrintMnemonic: false,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func NoSupplyGenesisState(cdc codec.Codec) app.GenesisState {
	genesisState := app.NewDefAppGenesisByDenom(fxtypes.DefaultDenom, cdc)

	// reset supply
	bankState := banktypes.DefaultGenesisState()
	bankState.DenomMetadata = []banktypes.Metadata{fxtypes.GetFXMetaData(fxtypes.DefaultDenom)}
	genesisState[banktypes.ModuleName] = cdc.MustMarshalJSON(bankState)

	var govGenState govtypes.GenesisState
	cdc.MustUnmarshalJSON(genesisState[govtypes.ModuleName], &govGenState)
	govGenState.VotingParams.VotingPeriod = time.Millisecond

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
		clienttypes.NewIdentifiedClientState(clientId, &ibctmtypes.ClientState{
			ChainId:      rand.Str(10),
			LatestHeight: clienttypes.NewHeight(0, 1),
			// if ibc timeout  ctx.BlockTime() > TrustingPeriod + consensusStateAny.Timestamp
			TrustingPeriod: time.Hour,
		}),
	}

	clientState.ClientsConsensus = []clienttypes.ClientConsensusStates{
		{
			ClientId: clientId,
			ConsensusStates: []clienttypes.ConsensusStateWithHeight{
				clienttypes.NewConsensusStateWithHeight(clienttypes.NewHeight(0, 1),
					ibctmtypes.NewConsensusState(time.Now(), types.NewMerkleRoot(rand.Bytes(32)), nil),
				),
			},
		},
	}

	counterparty := connectiontypes.NewCounterparty(clientId, connectionId, types.NewMerklePrefix(rand.Bytes(32)))
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
	genesisState[host.ModuleName] = cdc.MustMarshalJSON(&ibccoretypes.GenesisState{
		ClientGenesis:     clientState,
		ConnectionGenesis: connState,
		ChannelGenesis:    channelState,
	})
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

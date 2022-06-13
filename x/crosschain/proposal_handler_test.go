package crosschain_test

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	"github.com/functionx/fx-core/app/helpers"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/x/crosschain"
	"github.com/functionx/fx-core/x/crosschain/types"
)

func TestUpdateOracleProposal(t *testing.T) {
	minStakeAmount = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil))
	// get test env
	myApp, ctx, oracleAddressList, orchestratorAddressList, ethKeys, h, gh := createProposalTestEnv(t)
	var err error

	normalMsg := &types.MsgSetOrchestratorAddress{
		Oracle:          oracleAddressList[0].String(),
		Orchestrator:    orchestratorAddressList[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		StakeAmount:     sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: minStakeAmount},
		ChainName:       chainName,
	}
	_, err = h(ctx, normalMsg)
	require.NoError(t, err)

	keeper := myApp.BscKeeper

	addTwoOracleProposal := &types.UpdateChainOraclesProposal{
		Title:       "zzz",
		Description: "zzz",
		Oracles:     []string{oracleAddressList[0].String(), oracleAddressList[1].String(), oracleAddressList[2].String()},
		ChainName:   chainName,
	}
	err = gh(ctx, addTwoOracleProposal)
	require.NoError(t, err)
	require.True(t, keeper.IsProposalOracle(ctx, oracleAddressList[0].String()))
	require.True(t, keeper.IsProposalOracle(ctx, oracleAddressList[1].String()))
	require.True(t, keeper.IsProposalOracle(ctx, oracleAddressList[2].String()))
	require.False(t, keeper.IsProposalOracle(ctx, oracleAddressList[3].String()))

	deleteOneOracleProposal := &types.UpdateChainOraclesProposal{
		Title:       "zzz",
		Description: "zzz",
		Oracles:     []string{oracleAddressList[1].String(), oracleAddressList[2].String()},
		ChainName:   chainName,
	}
	err = gh(ctx, deleteOneOracleProposal)
	require.ErrorIs(t, err, types.ErrInvalid)
	require.EqualValues(t, fmt.Sprintf("max change power!maxChangePower:%v,deletePower:%v: %v", "0", "1", types.ErrInvalid), err.Error())

	var oracles []string
	for i := 0; i < types.MaxOracleSize+1; i++ {
		oracles = append(oracles, sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String())
	}
	addMaxOracleSizePlusOneProposal := &types.UpdateChainOraclesProposal{
		Title:       "zzz",
		Description: "zzz",
		Oracles:     oracles,
		ChainName:   chainName,
	}
	err = gh(ctx, addMaxOracleSizePlusOneProposal)
	require.ErrorIs(t, err, types.ErrInvalid)
	require.EqualValues(t, fmt.Sprintf("oracle length must be less than or equal : %d: %v", types.MaxOracleSize, types.ErrInvalid), err.Error())
}

func createProposalTestEnv(t *testing.T) (myApp *app.App, ctx sdk.Context, oracleAddressList, orchestratorAddressList []sdk.AccAddress, ethKeys []*ecdsa.PrivateKey, handler sdk.Handler, govHandler govtypes.Handler) {
	initBalances := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := helpers.GenerateGenesisValidator(t, 2,
		sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, initBalances)))
	myApp = helpers.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx = myApp.BaseApp.NewContext(false, tmproto.Header{})
	ctx = ctx.WithBlockHeight(2000000)
	oracleAddressList = helpers.AddTestAddrsIncremental(myApp, ctx, generateAccountNum, minStakeAmount.Mul(sdk.NewInt(1000)))
	orchestratorAddressList = helpers.AddTestAddrs(myApp, ctx, generateAccountNum, sdk.ZeroInt())
	ethKeys = genEthKey(generateAccountNum)
	// chain module oracle list
	var oracles []string
	for _, account := range oracleAddressList {
		oracles = append(oracles, account.String())
	}

	var err error
	// init bsc params by proposal
	proposalHandler := crosschain.NewCrossChainProposalHandler(myApp.CrosschainKeeper)
	err = proposalHandler(ctx, &types.InitCrossChainParamsProposal{
		Title:       "init bsc chain params",
		Description: "init fx chain <-> bsc chain params",
		Params:      defaultModuleParams(oracles),
		ChainName:   chainName,
	})
	require.NoError(t, err)

	crosschianHandler := crosschain.NewHandler(myApp.CrosschainKeeper)
	// To add a proxy handler, execute msg validateBasic
	proxyHandler := func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		require.NoError(t, msg.ValidateBasic(), fmt.Sprintf("msg %s validate basic error", sdk.MsgTypeURL(msg)))
		return crosschianHandler(ctx, msg)
	}
	return myApp, ctx, oracleAddressList, orchestratorAddressList, ethKeys, proxyHandler, proposalHandler
}

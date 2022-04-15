package crosschain_test

import (
	"crypto/ecdsa"
	"fmt"
	fxtypes "github.com/functionx/fx-core/types"
	"math/big"
	"testing"

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
	minDepositAmount = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil))
	// get test env
	fxcore, ctx, oracleAddressList, orchestratorAddressList, ethKeys, h, gh := createProposalTestEnv(t)
	var err error

	normalMsg := &types.MsgSetOrchestratorAddress{
		Oracle:          oracleAddressList[0].String(),
		Orchestrator:    orchestratorAddressList[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(ethKeys[0].PublicKey).Hex(),
		Deposit:         sdk.Coin{Denom: depositToken, Amount: minDepositAmount},
		ChainName:       chainName,
	}
	_, err = h(ctx, normalMsg)
	require.NoError(t, err)

	keeper := fxcore.BscKeeper

	addTwoOracleProposal := &types.UpdateChainOraclesProposal{
		Title:       "zzz",
		Description: "zzz",
		Oracles:     []string{oracleAddressList[0].String(), oracleAddressList[1].String(), oracleAddressList[2].String()},
		ChainName:   chainName,
	}
	err = gh(ctx, addTwoOracleProposal)
	require.NoError(t, err)
	require.True(t, keeper.IsOracle(ctx, oracleAddressList[0].String()))
	require.True(t, keeper.IsOracle(ctx, oracleAddressList[1].String()))
	require.True(t, keeper.IsOracle(ctx, oracleAddressList[2].String()))
	require.False(t, keeper.IsOracle(ctx, oracleAddressList[3].String()))

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
	accounts := app.CreateIncrementalAccounts(types.MaxOracleSize + 1)
	for i := 0; i < len(accounts); i++ {
		oracles = append(oracles, accounts[i].String())
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

func createProposalTestEnv(t *testing.T) (fxcore *app.App, ctx sdk.Context, oracleAddressList, orchestratorAddressList []sdk.AccAddress, ethKeys []*ecdsa.PrivateKey, handler sdk.Handler, govHandler govtypes.Handler) {
	initBalances := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := app.GenerateGenesisValidator(2,
		sdk.NewCoins(sdk.NewCoin(fxtypes.MintDenom, initBalances)))
	fxcore = app.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx = fxcore.BaseApp.NewContext(false, tmproto.Header{})
	ctx = ctx.WithBlockHeight(2000000)
	oracleAddressList = app.AddTestAddrsIncremental(fxcore, ctx, GenerateAccountNum, minDepositAmount.Mul(sdk.NewInt(1000)))
	orchestratorAddressList = app.AddTestAddrs(fxcore, ctx, GenerateAccountNum, sdk.ZeroInt())
	ethKeys = genEthKey(GenerateAccountNum)
	// chain module oracle list
	var oracles []string
	for _, account := range oracleAddressList {
		oracles = append(oracles, account.String())
	}

	var err error
	// init bsc params by proposal
	proposalHandler := crosschain.NewCrossChainProposalHandler(fxcore.CrosschainKeeper)
	err = proposalHandler(ctx, &types.InitCrossChainParamsProposal{
		Title:       "init bsc chain params",
		Description: "init fx chain <-> bsc chain params",
		Params:      defaultModuleParams(oracles),
		ChainName:   chainName,
	})
	require.NoError(t, err)

	crosschianHandler := crosschain.NewHandler(fxcore.CrosschainKeeper)
	// To add a proxy handler, execute msg validateBasic
	proxyHandler := func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		require.NoError(t, msg.ValidateBasic(), fmt.Sprintf("msg %s/%s validate basic error", msg.Route(), msg.Type()))
		return crosschianHandler(ctx, msg)
	}
	return fxcore, ctx, oracleAddressList, orchestratorAddressList, ethKeys, proxyHandler, proposalHandler
}

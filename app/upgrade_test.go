package app_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v5/app"
	v5 "github.com/functionx/fx-core/v5/app/upgrades/v5"
	"github.com/functionx/fx-core/v5/client/jsonrpc"
	"github.com/functionx/fx-core/v5/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v5/types"
)

func Test_TestnetUpgrade(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test: ", t.Name())

	fxtypes.SetConfig(true)
	fxtypes.SetChainId(fxtypes.TestnetChainId) // only for testnet

	testCases := []struct {
		name                  string
		fromVersion           int
		toVersion             int
		LocalStoreBlockHeight uint64
		plan                  upgradetypes.Plan
	}{
		{
			name:        "upgrade v5.0.x",
			fromVersion: 4,
			toVersion:   5,
			plan: upgradetypes.Plan{
				Name: v5.Upgrade.UpgradeName,
				Info: "local test upgrade v5.0.x",
			},
		},
	}

	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, filepath.Join(fxtypes.GetDefaultNodeHome(), "data"))
	require.NoError(t, err)

	makeEncodingConfig := app.MakeEncodingConfig()
	myApp := app.New(log.NewFilter(log.NewTMLogger(os.Stdout), log.AllowAll()),
		db, nil, false, map[int64]bool{}, fxtypes.GetDefaultNodeHome(), 0,
		makeEncodingConfig, app.EmptyAppOptions{})
	// todo default DefaultStoreLoader  New module verification failed
	myApp.SetStoreLoader(upgradetypes.UpgradeStoreLoader(myApp.LastBlockHeight()+1, v5.Upgrade.StoreUpgrades()))
	err = myApp.LoadLatestVersion()
	require.NoError(t, err)

	ctx := newContext(t, myApp)
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.plan.Height = ctx.BlockHeight()

			myApp.UpgradeKeeper.ApplyUpgrade(ctx, testCase.plan)
		})
	}

	checkSlashPeriod(t, ctx, myApp)

	myApp.EthKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
	myApp.BscKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
	myApp.TronKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
	myApp.PolygonKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
	myApp.AvalancheKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + 1))
}

func newContext(t *testing.T, myApp *app.App) sdk.Context {
	chainId := fxtypes.MainnetChainId
	if os.Getenv("CHAIN_ID") == fxtypes.TestnetChainId {
		chainId = fxtypes.TestnetChainId
	}
	ctx := myApp.NewUncachedContext(false, tmproto.Header{
		ChainID: chainId, Height: myApp.LastBlockHeight(),
	})
	// set the first validator to proposer
	validators := myApp.StakingKeeper.GetAllValidators(ctx)
	assert.True(t, len(validators) > 0)
	var pubKey cryptotypes.PubKey
	assert.NoError(t, myApp.AppCodec().UnpackAny(validators[0].ConsensusPubkey, &pubKey))
	ctx = ctx.WithProposer(pubKey.Address().Bytes())
	return ctx
}

func checkSlashPeriod(t *testing.T, ctx sdk.Context, myApp *app.App) {
	for val := range v5.ValidatorSlashHeightTestnetFXV4 {
		valAddr, err := sdk.ValAddressFromBech32(val)
		require.NoError(t, err)
		delegations := myApp.StakingKeeper.GetValidatorDelegations(ctx, valAddr)

		for _, del := range delegations {
			cacheCtx, _ := ctx.CacheContext()
			_, err := myApp.DistrKeeper.DelegationRewards(sdk.WrapSDKContext(cacheCtx), &distributiontypes.QueryDelegationRewardsRequest{
				DelegatorAddress: del.DelegatorAddress,
				ValidatorAddress: del.ValidatorAddress,
			})
			assert.NoError(t, err)
		}

		// withdraw
		for _, del := range delegations {
			cacheCtx, _ := ctx.CacheContext()
			_, err := myApp.DistrKeeper.WithdrawDelegationRewards(cacheCtx, del.GetDelegatorAddr(), del.GetValidatorAddr())
			assert.NoError(t, err)
		}

		// undelegate
		for _, del := range delegations {
			cacheCtx, _ := ctx.CacheContext()
			_, err := myApp.StakingKeeper.Undelegate(cacheCtx, del.GetDelegatorAddr(), del.GetValidatorAddr(), del.GetShares())
			assert.NoError(t, err)
		}

		// delegate
		count := 0
		for _, del := range delegations {
			cacheCtx, _ := ctx.CacheContext()
			fxBalance := myApp.BankKeeper.GetBalance(ctx, del.GetDelegatorAddr(), fxtypes.DefaultDenom)
			if fxBalance.Amount.IsZero() {
				continue
			}
			validator, found := myApp.StakingKeeper.GetValidator(ctx, del.GetValidatorAddr())
			require.True(t, found)
			_, err := myApp.StakingKeeper.Delegate(cacheCtx, del.GetDelegatorAddr(), fxBalance.Amount, stakingtypes.Unbonded, validator, true)
			assert.NoError(t, err)
			count++
		}
		assert.True(t, count > 0)
	}
}

func TestSlashPeriodTestnetFXV4(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	fxtypes.SetConfig(true)
	rpc := jsonrpc.NewNodeRPC(jsonrpc.NewClient("https://testnet-fx-json.functionx.io:26657"))
	query := fmt.Sprintf("block.height > %d AND slash.reason = 'missing_signature'", fxtypes.TestnetBlockHeightV4)

	blockSearch, err := rpc.BlockSearch(query, 1, 100, "")
	require.NoError(t, err)

	slashedFXV4 := make(map[string][]int64, len(v5.ValidatorSlashHeightTestnetFXV4))
	slashedValFXV4 := make([]string, 0, len(v5.ValidatorSlashHeightTestnetFXV4))
	for _, block := range blockSearch.Blocks {
		results, err := rpc.BlockResults(block.Block.Height)
		require.NoError(t, err)
		for _, result := range results.BeginBlockEvents {
			if result.Type != slashingtypes.EventTypeLiveness {
				continue
			}
			for _, attr := range result.Attributes {
				if string(attr.Key) != slashingtypes.AttributeKeyAddress {
					continue
				}
				valAddr, err := rpc.GetValAddressByCons(string(attr.Value))
				assert.NoError(t, err)
				if _, ok := slashedFXV4[valAddr.String()]; ok {
					slashedFXV4[valAddr.String()] = append(slashedFXV4[valAddr.String()], block.Block.Height)
				} else {
					slashedFXV4[valAddr.String()] = []int64{block.Block.Height}
					slashedValFXV4 = append(slashedValFXV4, valAddr.String())
				}
			}
		}
	}
	eq := len(v5.ValidatorSlashHeightTestnetFXV4) == len(slashedFXV4)
	assert.True(t, eq)

	if eq {
		for val, h1 := range slashedFXV4 {
			h2, ok := v5.ValidatorSlashHeightTestnetFXV4[val]
			assert.True(t, ok, "val: %s", val)
			sort.SliceStable(h1, func(i, j int) bool {
				return h1[i] < h1[j]
			})
			eq = assert.ObjectsAreEqual(h1, h2)
			assert.True(t, eq, "val: %s", val)
		}
	}

	// print
	if !eq {
		sort.SliceStable(slashedValFXV4, func(i, j int) bool {
			return slashedValFXV4[i] < slashedValFXV4[j]
		})
		for _, addr := range slashedValFXV4 {
			heights := slashedFXV4[addr]
			sort.SliceStable(heights, func(i, j int) bool {
				return heights[i] < heights[j]
			})
			t.Log("val:", addr, "heights:", heights)
		}
	}
}

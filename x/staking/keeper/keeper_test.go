package keeper_test

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *app.App
	queryClient types.QueryClient
	signer      *helpers.Signer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// Test helpers
func (suite *KeeperTestSuite) SetupTest() {
	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(suite.T(), err)
	suite.signer = helpers.NewSigner(priv)

	set, accs, balances := helpers.GenerateGenesisValidator(tmrand.Intn(10)+5, nil)
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), set, accs, balances...)

	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		Height:          suite.app.LastBlockHeight() + 1,
		ChainID:         fxtypes.ChainId(),
		ProposerAddress: set.Proposer.Address,
		Time:            time.Now().UTC(),
	})
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdk.OneInt())))
	suite.ctx = suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(1e18))

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, stakingkeeper.Querier{Keeper: suite.app.StakingKeeper.Keeper})
	suite.queryClient = types.NewQueryClient(queryHelper)

	for _, validator := range set.Validators {
		signingInfo := slashingtypes.NewValidatorSigningInfo(
			validator.Address.Bytes(),
			suite.ctx.BlockHeight(),
			0,
			time.Unix(0, 0),
			false,
			0,
		)
		suite.app.SlashingKeeper.SetValidatorSigningInfo(suite.ctx, validator.Address.Bytes(), signingInfo)
	}

	helpers.AddTestAddr(suite.app, suite.ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1000).Mul(sdk.NewInt(1e18)))))
}

func (suite *KeeperTestSuite) Commit() {
	height := suite.ctx.BlockHeight()
	suite.app.EndBlock(abci.RequestEndBlock{Height: height})
	suite.app.Commit()
	header := suite.ctx.BlockHeader()
	header.Height = height + 1
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: header,
		LastCommitInfo: abci.LastCommitInfo{
			Votes: []abci.VoteInfo{
				{
					Validator: abci.Validator{
						Address: suite.ctx.BlockHeader().ProposerAddress,
						Power:   math.MaxInt64,
					},
					SignedLastBlock: false,
				},
			},
		},
	})
	suite.ctx = suite.app.NewContext(false, header)
}

func (suite *KeeperTestSuite) RandSigner() *helpers.Signer {
	privKey := helpers.NewEthPrivKey()
	helpers.AddTestAddr(suite.app, suite.ctx, privKey.PubKey().Address().Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1e6).Mul(sdk.NewInt(1e18)))))
	return helpers.NewSigner(privKey)
}

func (suite *KeeperTestSuite) RandDelegates(setVal ...types.Validator) (signer *helpers.Signer, val sdk.ValAddress, lpToken common.Address, bondAmt sdk.Int, share sdk.Dec) {
	validators := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
	randVal := validators[tmrand.Intn(len(validators))]
	if len(setVal) > 0 {
		randVal = setVal[0]
	}

	lpToken, found := suite.app.StakingKeeper.GetValidatorLPToken(suite.ctx, randVal.GetOperator())
	suite.Require().True(found)

	signer = suite.RandSigner()
	bondAmt = sdk.NewInt(int64(tmrand.Int() + 100)).Mul(sdk.NewInt(1e18))
	helpers.AddTestAddr(suite.app, suite.ctx, signer.AccAddress().Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, bondAmt)))

	var err error
	share, err = suite.app.StakingKeeper.Delegate(suite.ctx, signer.AccAddress(), bondAmt, types.Unbonded, randVal, true)
	suite.Require().NoError(err)

	var res struct{ Value *big.Int }
	err = suite.app.EvmKeeper.CallContract(suite.ctx, signer.Address(), lpToken, fxtypes.GetLPToken().ABI, "balanceOf", &res, signer.Address())
	suite.Require().NoError(err)

	suite.Require().Equal(share.BigInt(), res.Value)

	return signer, randVal.GetOperator(), lpToken, bondAmt, share
}

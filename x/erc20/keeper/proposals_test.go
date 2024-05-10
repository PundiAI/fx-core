package keeper_test

import (
	"fmt"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/exp/slices"

	"github.com/functionx/fx-core/v7/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	"github.com/functionx/fx-core/v7/x/erc20/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
)

func (suite *KeeperTestSuite) setupRegisterERC20Pair() common.Address {
	contractAddr, err := suite.DeployContract(suite.signer.Address())
	suite.NoError(err)
	_, err = suite.app.Erc20Keeper.RegisterNativeERC20(suite.ctx, contractAddr)
	suite.NoError(err)
	return contractAddr
}

func (suite *KeeperTestSuite) setupRegisterERC20PairAddAliases() common.Address {
	contractAddr, err := suite.DeployContract(suite.signer.Address())
	suite.NoError(err)
	_, err = suite.app.Erc20Keeper.RegisterNativeERC20(suite.ctx, contractAddr, "eth0xdAC17F958D2ee523a2206206994597C13D831ec7")
	suite.NoError(err)
	return contractAddr
}

func (suite *KeeperTestSuite) setupRegisterCoin() (banktypes.Metadata, *types.TokenPair) {
	metadata := newMetadata()
	pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, metadata)
	suite.NoError(err)
	return metadata, pair
}

func (suite *KeeperTestSuite) TestRegisterCoinWithAlias() {
	metadata := newMetadata()

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
		errMsg   string
	}{
		{
			"ok",
			func() {},
			true,
			"",
		},
		{
			"intrarelaying is disabled globally",
			func() {
				params := types.DefaultParams()
				params.EnableErc20 = false
				err := suite.app.Erc20Keeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
			},
			false,
			"registration is currently disabled by governance: erc20 module is disabled",
		},
		{
			"denom already registered",
			func() {
				regPair := types.NewTokenPair(helpers.GenerateAddress(), metadata.Base, true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, regPair)
			},
			false,
			"coin denomination already registered: usdt: token pair already exists",
		},
		{
			"alias already registered denom",
			func() {
				regPair := types.NewTokenPair(helpers.GenerateAddress(), metadata.DenomUnits[0].Aliases[0], true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, regPair)
			},
			false,
			fmt.Sprintf("denom %s already registered: invalid metadata", metadata.DenomUnits[0].Aliases[0]),
		},
		{
			"denom register as alias",
			func() {
				suite.app.Erc20Keeper.SetAliasesDenom(suite.ctx, "test", metadata.Base)
			},
			false,
			"alias usdt already registered: invalid metadata",
		},
		{
			"alias already registered",
			func() {
				suite.app.Erc20Keeper.SetAliasesDenom(suite.ctx, metadata.Base, metadata.DenomUnits[0].Aliases[0])
			},
			false,
			fmt.Sprintf("alias %s already registered: invalid metadata", metadata.DenomUnits[0].Aliases[0]),
		},
		{
			"alias equal base",
			func() {
				metadata.DenomUnits[0].Aliases = []string{"usdt"}
			},
			false,
			"alias can not equal base, display or symbol: invalid metadata",
		},
		{
			"alias equal display",
			func() {
				metadata.DenomUnits[0].Aliases = []string{"display usdt"}
			},
			false,
			"alias can not equal base, display or symbol: invalid metadata",
		},
		{
			"alias equal symbol",
			func() {
				metadata.DenomUnits[0].Aliases = []string{"USDT"}
			},
			false,
			"alias can not equal base, display or symbol: invalid metadata",
		},
		{
			"alias empty",
			func() {
				metadata.DenomUnits[0].Aliases = []string{}
			},
			true,
			"",
		},
		{
			"fx",
			func() {
				metadata = fxtypes.GetFXMetaData()
			},
			false,
			"coin denomination already registered: FX: token pair already exists",
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			pair, tcErr := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, metadata)

			erc20ModuleAddr := common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes())
			contractAddr := crypto.CreateAddress(erc20ModuleAddr, suite.app.EvmKeeper.GetNonce(suite.ctx, erc20ModuleAddr)-1)
			expPair := &types.TokenPair{
				Erc20Address:  contractAddr.String(),
				Denom:         metadata.Base,
				Enabled:       true,
				ContractOwner: 1,
			}

			if tc.expPass {
				suite.NoError(tcErr, tc.name)
				suite.Equal(pair, expPair)
				for _, alias := range metadata.DenomUnits[0].Aliases {
					suite.True(suite.app.Erc20Keeper.IsAliasDenomRegistered(suite.ctx, alias))
				}
			} else {
				suite.Error(tcErr, tc.name)
				suite.EqualError(tcErr, tc.errMsg, tc.name)
			}
		})
	}
}

//gocyclo:ignore
func (suite *KeeperTestSuite) TestUpdateDenomAlias() {
	denom := fmt.Sprintf("test%s", helpers.GenerateAddress().Hex())
	metadata := newMetadata()

	testCases := []struct {
		name     string
		malleate func() (*types.TokenPair, error)
		expPass  bool
		alias    []string
	}{
		{
			name: "success - add alias",
			malleate: func() (*types.TokenPair, error) {
				pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, metadata)
				suite.NoError(err)

				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAliases(suite.ctx, "usdt", denom)
				if err != nil {
					return nil, err
				}
				suite.True(addAlias)

				return pair, nil
			},
			expPass: true,
			alias:   append(metadata.DenomUnits[0].Aliases, denom),
		},
		{
			name: "success - delete alias",
			malleate: func() (*types.TokenPair, error) {
				pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, metadata)
				suite.NoError(err)

				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAliases(suite.ctx, "usdt", metadata.DenomUnits[0].Aliases[0])
				if err != nil {
					return nil, err
				}
				suite.False(addAlias)
				return pair, nil
			},
			expPass: true,
			alias:   metadata.DenomUnits[0].Aliases[1:],
		},
		{
			name: "failed - denom not equal",
			malleate: func() (*types.TokenPair, error) {
				pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, metadata)
				suite.NoError(err)

				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAliases(suite.ctx, "abc", denom)
				if err != nil {
					return nil, err
				}
				suite.True(addAlias)
				return pair, nil
			},
			expPass: false,
			alias:   []string{},
		},
		{
			name: "failed - alias registered",
			malleate: func() (*types.TokenPair, error) {
				pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, metadata)
				suite.NoError(err)
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAliases(suite.ctx, "abc", metadata.DenomUnits[0].Aliases[0])
				if err != nil {
					return nil, err
				}
				suite.True(addAlias)
				return pair, nil
			},
			expPass: false,
			alias:   []string{},
		},
		{
			name: "success - empty alias",
			malleate: func() (*types.TokenPair, error) {
				pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, banktypes.Metadata{
					Description: "The cross chain token of Function X",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    "abc",
							Exponent: 0,
							Aliases:  nil,
						},
					},
					Base:    "abc",
					Display: "abc",
					Name:    "Token ABC",
					Symbol:  "ABC",
				})
				if err != nil {
					return nil, err
				}
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAliases(suite.ctx, "abc", denom)
				if err != nil {
					return nil, err
				}
				suite.True(addAlias)
				return pair, nil
			},
			expPass: true,
			alias:   []string{denom},
		},
		{
			name: "success - remove alias empty",
			malleate: func() (*types.TokenPair, error) {
				pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, banktypes.Metadata{
					Description: "The cross chain token of Function X",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    "abc",
							Exponent: 0,
							Aliases:  []string{denom},
						},
					},
					Base:    "abc",
					Display: "abc",
					Name:    "Token ABC",
					Symbol:  "ABC",
				})
				if err != nil {
					return nil, err
				}
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAliases(suite.ctx, "abc", denom)
				if err != nil {
					return nil, err
				}
				suite.False(addAlias)
				return pair, nil
			},
			expPass: true,
			alias:   []string{},
		},
		{
			name: "failed - aliases can not empty",
			malleate: func() (*types.TokenPair, error) {
				pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, metadata)
				suite.NoError(err)

				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAliases(suite.ctx, "abc", denom)
				if err != nil {
					return nil, err
				}
				suite.True(addAlias)
				return pair, nil
			},
			expPass: false,
			alias:   []string{},
		},
		{
			name: "failed - alias denom not equal with update denom",
			malleate: func() (*types.TokenPair, error) {
				pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, metadata)
				suite.NoError(err)

				suite.app.Erc20Keeper.SetAliasesDenom(suite.ctx, "abc", denom)

				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAliases(suite.ctx, "usdt", denom)
				if err != nil {
					return nil, err
				}
				suite.True(addAlias)
				return pair, nil
			},
			expPass: false,
			alias:   []string{},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			pair, tcErr := tc.malleate()

			if tc.expPass {
				suite.NoError(tcErr, tc.name)
				md, found := suite.app.Erc20Keeper.GetValidMetadata(suite.ctx, pair.Denom)
				suite.True(found)
				if len(tc.alias) == 0 && len(md.DenomUnits[0].Aliases) == 0 {
					return // remove all alias
				}
				suite.Equal(md.DenomUnits[0].Aliases, tc.alias)
				for _, alias := range tc.alias {
					aliasRegistered := suite.app.Erc20Keeper.IsAliasDenomRegistered(suite.ctx, alias)
					suite.True(aliasRegistered)
				}
				for _, alias := range metadata.DenomUnits[0].Aliases {
					if !slices.Contains(tc.alias, alias) {
						aliasRegistered := suite.app.Erc20Keeper.IsAliasDenomRegistered(suite.ctx, alias)
						suite.False(aliasRegistered, alias)
					}
				}
			} else {
				suite.Error(tcErr, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestRegisterERC20() {
	testCases := []struct {
		name     string
		malleate func(contractAddr common.Address)
		expPass  bool
	}{
		{
			"intrarelaying is disabled globally",
			func(contractAddr common.Address) {
				params := types.DefaultParams()
				params.EnableErc20 = false
				err := suite.app.Erc20Keeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
			},
			false,
		},
		{
			"token ERC20 already registered",
			func(contractAddr common.Address) {
				pair := types.NewTokenPair(contractAddr, "test", true, types.OWNER_EXTERNAL)
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, pair)
			},
			false,
		},
		{
			"denom already registered",
			func(contractAddr common.Address) {
				pair := types.NewTokenPair(contractAddr, "test", true, types.OWNER_EXTERNAL)
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, pair)
			},
			false,
		},
		{
			"alias already registered",
			func(contractAddr common.Address) {
				pair := types.NewTokenPair(contractAddr, "test", true, types.OWNER_EXTERNAL)
				suite.app.Erc20Keeper.SetAliasesDenom(suite.ctx, "test", pair.Denom)
			},
			false,
		},
		{
			"meta data already stored",
			func(contractAddr common.Address) {
				_, err := suite.app.Erc20Keeper.RegisterNativeERC20(suite.ctx, contractAddr, "eth0xdAC17F958D2ee523a2206206994597C13D831ec7")
				suite.NoError(err)
			},
			false,
		},
		{
			"ok",
			func(contractAddr common.Address) {},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			contractAddr, err := suite.DeployContract(suite.signer.Address())
			suite.NoError(err)

			tc.malleate(contractAddr)

			coinName := "test"
			_, err = suite.app.Erc20Keeper.RegisterNativeERC20(suite.ctx, contractAddr, "eth0xdAC17F958D2ee523a2206206994597C13D831ec7")
			metadata, _ := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, coinName)
			if tc.expPass {
				suite.NoError(err, tc.name)
				// Metadata variables
				suite.Equal(coinName, metadata.Base)
				suite.Equal(coinName, metadata.Display)
				// Denom units
				suite.Equal(len(metadata.DenomUnits), 2)
				suite.Equal(coinName, metadata.DenomUnits[0].Denom)
				suite.Equal(uint32(18), metadata.DenomUnits[1].Exponent)
				suite.Equal(strings.ToUpper(coinName), metadata.DenomUnits[1].Denom)
			} else {
				suite.Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestToggleRelay() {
	testCases := []struct {
		name         string
		malleate     func() common.Address
		expPass      bool
		relayEnabled bool
	}{
		{
			"token not registered",
			func() common.Address {
				contractAddr, err := suite.DeployContract(suite.signer.Address())
				suite.NoError(err)
				return contractAddr
			},
			false,
			false,
		},
		{
			"token not registered - pair not found",
			func() common.Address {
				contractAddr, err := suite.DeployContract(suite.signer.Address())
				suite.NoError(err)
				return contractAddr
			},
			false,
			false,
		},
		{
			"disable relay",
			func() common.Address {
				contractAddr := suite.setupRegisterERC20Pair()
				_, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, contractAddr.String())
				suite.True(found)
				return contractAddr
			},
			true,
			false,
		},
		{
			"disable and enable relay",
			func() common.Address {
				contractAddr := suite.setupRegisterERC20Pair()
				pair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, contractAddr.String())
				suite.True(found)

				pair1, err := suite.app.Erc20Keeper.ToggleTokenConvert(suite.ctx, contractAddr.String())
				suite.NoError(err)
				pair.Enabled = !pair.Enabled
				suite.Equal(pair, pair1)
				return contractAddr
			},
			true,
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			contractAddr := tc.malleate()
			pair, err := suite.app.Erc20Keeper.ToggleTokenConvert(suite.ctx, contractAddr.String())
			if tc.expPass {
				suite.NoError(err, tc.name)
				suite.Equal(pair.Enabled, tc.relayEnabled)
			} else {
				suite.Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestRegisterCoinConversionInvariant() {
	metadata := newMetadata()
	pair, err := suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, metadata)
	suite.Require().NoError(err)
	suite.Require().True(pair.Enabled)
	initCoin := sdk.NewCoin(pair.Denom, sdkmath.NewInt(1000).MulRaw(1e18))
	helpers.AddTestAddr(suite.app, suite.ctx, suite.signer.AccAddress(), sdk.NewCoins(initCoin))
	beforeBalanceCoin := suite.app.BankKeeper.GetBalance(suite.ctx, suite.signer.AccAddress(), pair.Denom)
	beforeBalanceToken := suite.BalanceOf(pair.GetERC20Contract(), suite.signer.Address())
	// coin --> erc20  // Coin: initCoin, Receiver: suite.signer.Address().String(), Sender: suite.signer.AccAddress().String()
	_, err = suite.app.Erc20Keeper.ConvertCoin(suite.ctx, types.NewMsgConvertCoin(initCoin, suite.signer.Address(), suite.signer.AccAddress()))
	suite.Require().NoError(err)
	lockBalance := suite.app.BankKeeper.GetBalance(suite.ctx, authtypes.NewModuleAddress(types.ModuleName), initCoin.Denom)
	suite.Require().EqualValues(lockBalance.String(), initCoin.String())
	suite.Require().EqualValues(beforeBalanceCoin.Sub(initCoin).String(), sdk.NewCoin(pair.Denom, sdkmath.NewInt(0)).String())
	afterBalanceToken := suite.BalanceOf(pair.GetERC20Contract(), suite.signer.Address())
	suite.Require().EqualValues(afterBalanceToken.String(), new(big.Int).Add(beforeBalanceToken, initCoin.Amount.BigInt()).String())
	// erc20 --> coin
	_, err = suite.app.Erc20Keeper.ConvertERC20(suite.ctx, types.NewMsgConvertERC20(sdkmath.NewIntFromBigInt(afterBalanceToken), suite.signer.AccAddress(), common.HexToAddress(pair.Erc20Address), suite.signer.Address()))
	suite.Require().NoError(err)
	lockBalance = suite.app.BankKeeper.GetBalance(suite.ctx, authtypes.NewModuleAddress(types.ModuleName), initCoin.Denom)
	suite.Require().EqualValues(lockBalance.String(), sdk.NewCoin(initCoin.Denom, sdkmath.NewInt(0)).String())
	balanceCoin := suite.app.BankKeeper.GetBalance(suite.ctx, suite.signer.AccAddress(), pair.Denom)
	suite.Require().EqualValues(balanceCoin.String(), beforeBalanceCoin.String())
	balanceToken := suite.BalanceOf(pair.GetERC20Contract(), suite.signer.Address())
	suite.Require().EqualValues(balanceToken.String(), sdkmath.NewInt(0).String())
}

func (suite *KeeperTestSuite) TestRegisterERC20ConversionInvariant() {
	contact, err := suite.app.Erc20Keeper.DeployUpgradableToken(suite.ctx, suite.signer.Address(), "Test token", "TEST", 18)
	suite.Require().NoError(err)
	tokenPair, err := suite.app.Erc20Keeper.RegisterNativeERC20(suite.ctx, contact, crosschaintypes.NewBridgeDenom(ethtypes.ModuleName, helpers.GenerateAddress().String()))
	suite.Require().NoError(err)
	suite.Require().True(tokenPair.Enabled)
	suite.Require().EqualValues(tokenPair.Erc20Address, contact.String())
	beforeMintBalance := suite.BalanceOf(contact, suite.signer.Address())
	suite.Require().EqualValues(beforeMintBalance.String(), big.NewInt(0).String())
	initBalance := sdkmath.NewInt(1000).MulRaw(1e18)
	_, err = suite.app.EvmKeeper.ApplyContract(suite.ctx, suite.signer.Address(), contact, nil, contract.GetFIP20().ABI, "mint", suite.signer.Address(), initBalance.BigInt())
	suite.Require().NoError(err)
	afterMintBalance := suite.BalanceOf(contact, suite.signer.Address())
	suite.Require().EqualValues(afterMintBalance.String(), initBalance.String())
	beforeCoin := suite.app.BankKeeper.GetBalance(suite.ctx, suite.signer.AccAddress(), tokenPair.Denom)
	suite.Require().EqualValues(beforeCoin.String(), sdk.NewCoin(tokenPair.Denom, sdkmath.NewInt(0)).String())
	// ERC20 token -> coin
	_, err = suite.app.Erc20Keeper.ConvertERC20(suite.ctx, types.NewMsgConvertERC20(initBalance, suite.signer.AccAddress(), contact, suite.signer.Address()))
	suite.Require().NoError(err)
	afterBalance := suite.BalanceOf(contact, suite.signer.Address())
	suite.Require().EqualValues(afterBalance.String(), sdkmath.NewInt(0).String())
	lockBalance := suite.BalanceOf(contact, suite.app.Erc20Keeper.ModuleAddress())
	suite.Require().EqualValues(lockBalance.String(), initBalance.String())
	afterCoin := suite.app.BankKeeper.GetBalance(suite.ctx, suite.signer.AccAddress(), tokenPair.Denom)
	suite.Require().EqualValues(afterCoin.String(), sdk.NewCoin(tokenPair.Denom, initBalance).String())
	// coin -> erc20
	_, err = suite.app.Erc20Keeper.ConvertCoin(suite.ctx, types.NewMsgConvertCoin(afterCoin, suite.signer.Address(), suite.signer.AccAddress()))
	suite.Require().NoError(err)
	afterBalance = suite.BalanceOf(contact, suite.signer.Address())
	suite.Require().EqualValues(afterBalance.String(), afterCoin.Amount.String())
	lockBalance = suite.BalanceOf(contact, suite.app.Erc20Keeper.ModuleAddress())
	suite.Require().EqualValues(lockBalance.String(), sdkmath.NewInt(0).String())
	afterCoin = suite.app.BankKeeper.GetBalance(suite.ctx, suite.signer.AccAddress(), tokenPair.Denom)
	suite.Require().EqualValues(afterCoin.String(), sdk.NewCoin(tokenPair.Denom, sdkmath.NewInt(0)).String())
}

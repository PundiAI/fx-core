package keeper_test

import (
	"fmt"
	"strings"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/exp/slices"

	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (suite *KeeperTestSuite) setupRegisterERC20Pair() common.Address {
	contractAddr, err := suite.DeployContract(suite.signer.Address())
	suite.NoError(err)
	_, err = suite.app.Erc20Keeper.RegisterERC20(suite.ctx, contractAddr)
	suite.NoError(err)
	return contractAddr
}

func (suite *KeeperTestSuite) setupRegisterCoin() (banktypes.Metadata, *types.TokenPair) {
	metadata := newMetadata()
	pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, metadata)
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
				metadata = fxtypes.GetFXMetaData(fxtypes.DefaultDenom)
			},
			false,
			"coin denomination already registered: FX: token pair already exists",
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			pair, tcErr := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, metadata)

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

func (suite *KeeperTestSuite) TestUpdateDenomAlias() {
	denom := fmt.Sprintf("test%s", helpers.GenerateAddress().Hex())
	metadata := newMetadata()

	testCases := []struct {
		name     string
		malleate func() error
		expPass  bool
		alias    []string
	}{
		{
			name: "success - add alias",
			malleate: func() error {
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "usdt", denom)
				if err != nil {
					return err
				}
				suite.True(addAlias)

				return nil
			},
			expPass: true,
			alias:   append(metadata.DenomUnits[0].Aliases, denom),
		},
		{
			name: "success - delete alias",
			malleate: func() error {
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "usdt", metadata.DenomUnits[0].Aliases[0])
				if err != nil {
					return err
				}
				suite.False(addAlias)
				return nil
			},
			expPass: true,
			alias:   metadata.DenomUnits[0].Aliases[1:],
		},
		{
			name: "failed - denom not equal",
			malleate: func() error {
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "abc", denom)
				if err != nil {
					return err
				}
				suite.True(addAlias)
				return nil
			},
			expPass: false,
			alias:   []string{},
		},
		{
			name: "failed - alias registered",
			malleate: func() error {
				_, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, banktypes.Metadata{
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
					return err
				}
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "abc", denom)
				if err != nil {
					return err
				}
				suite.True(addAlias)
				return nil
			},
			expPass: false,
			alias:   []string{},
		},
		{
			name: "failed - metadata not support many to one",
			malleate: func() error {
				_, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, banktypes.Metadata{
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
					return err
				}
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "abc", denom)
				if err != nil {
					return err
				}
				suite.True(addAlias)
				return nil
			},
			expPass: false,
			alias:   []string{},
		},
		{
			name: "failed - aliases can not empty",
			malleate: func() error {
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "abc", denom)
				if err != nil {
					return err
				}
				suite.True(addAlias)
				return nil
			},
			expPass: false,
			alias:   []string{},
		},
		{
			name: "failed - alias denom not equal with update denom",
			malleate: func() error {
				suite.app.Erc20Keeper.SetAliasesDenom(suite.ctx, "abc", denom)

				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "usdt", denom)
				if err != nil {
					return err
				}
				suite.True(addAlias)
				return nil
			},
			expPass: false,
			alias:   []string{},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, metadata)
			suite.NoError(err)

			tcErr := tc.malleate()

			if tc.expPass {
				suite.NoError(tcErr, tc.name)
				md, found := suite.app.Erc20Keeper.HasDenomAlias(suite.ctx, pair.Denom)
				suite.True(found)
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
				_, err := suite.app.Erc20Keeper.RegisterERC20(suite.ctx, contractAddr)
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
			_, err = suite.app.Erc20Keeper.RegisterERC20(suite.ctx, contractAddr)
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

				pair1, err := suite.app.Erc20Keeper.ToggleRelay(suite.ctx, contractAddr.String())
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
			pair, err := suite.app.Erc20Keeper.ToggleRelay(suite.ctx, contractAddr.String())
			if tc.expPass {
				suite.NoError(err, tc.name)
				suite.Equal(pair.Enabled, tc.relayEnabled)
			} else {
				suite.Error(err, tc.name)
			}
		})
	}
}

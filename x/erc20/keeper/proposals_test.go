package keeper_test

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

const (
	contractMinterBurner = iota + 1
)

const (
	erc20Name       = "Coin Token"
	erc20Symbol     = "CTKN"
	erc20Decimals   = uint8(18)
	cosmosTokenBase = "acoin"
	defaultExponent = uint32(18)
)

func (suite *KeeperTestSuite) setupRegisterERC20Pair(contractType int) common.Address {
	suite.SetupTest()

	var contractAddr common.Address
	// Deploy contract
	switch contractType {
	default:
		contractAddr, _ = suite.DeployContract(suite.address, erc20Name, erc20Symbol, erc20Decimals)
	}
	//suite.Commit()

	_, err := suite.app.Erc20Keeper.RegisterERC20(suite.ctx, contractAddr)
	suite.Require().NoError(err)
	return contractAddr
}

func (suite *KeeperTestSuite) setupRegisterCoin() (banktypes.Metadata, *types.TokenPair) {
	validMetadata := banktypes.Metadata{
		Description: "description of the token",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    cosmosTokenBase,
				Exponent: uint32(0),
			}, {
				Denom:    erc20Symbol,
				Exponent: uint32(18),
			},
		},
		Base:    cosmosTokenBase,
		Display: cosmosTokenBase,
		Name:    erc20Name,
		Symbol:  erc20Symbol,
	}

	// pair := types.NewTokenPair(contractAddr, cosmosTokenBase, true, types.OWNER_MODULE)
	pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, validMetadata)
	suite.Require().NoError(err)
	//suite.Commit()
	return validMetadata, pair
}

func (suite *KeeperTestSuite) setupRegisterCoinUSDT(alias ...string) (banktypes.Metadata, *types.TokenPair) {
	if len(alias) == 0 {
		alias = []string{ethDenom, polygonDenom}
	}
	validMetadata := banktypes.Metadata{
		Description: "description of the token",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "usdt",
				Exponent: uint32(0),
				Aliases:  alias,
			}, {
				Denom:    "USDT",
				Exponent: uint32(18),
			},
		},
		Base:    "usdt",
		Display: "usdt",
		Name:    "Tether USD",
		Symbol:  "USDT",
	}

	pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, validMetadata)
	suite.Require().NoError(err)

	return validMetadata, pair
}

func (suite *KeeperTestSuite) setupRegisterCoinUSDTWithOutAlias() (banktypes.Metadata, *types.TokenPair) {
	validMetadata := banktypes.Metadata{
		Description: "description of the token",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "usdt",
				Exponent: uint32(0),
			}, {
				Denom:    "USDT",
				Exponent: uint32(18),
			},
		},
		Base:    "usdt",
		Display: "usdt",
		Name:    "Tether USD",
		Symbol:  "USDT",
	}

	pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, validMetadata)
	suite.Require().NoError(err)

	return validMetadata, pair
}

func (suite *KeeperTestSuite) TestRegisterCoin() {
	metadata := banktypes.Metadata{
		Description: "description",
		// NOTE: Denom units MUST be increasing
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    cosmosTokenBase,
				Exponent: uint32(0),
			},
			{
				Denom:    erc20Symbol,
				Exponent: uint32(defaultExponent),
			},
		},
		Base:    cosmosTokenBase,
		Display: cosmosTokenBase,
		Name:    erc20Name,
		Symbol:  erc20Symbol,
	}

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"intrarelaying is disabled globally",
			func() {
				params := types.DefaultParams()
				params.EnableErc20 = false
				suite.app.Erc20Keeper.SetParams(suite.ctx, params)
			},
			false,
		},
		{
			"denom already registered",
			func() {
				regPair := types.NewTokenPair(helpers.GenerateAddress(), metadata.Base, true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.SetDenomMap(suite.ctx, regPair.Denom, regPair.GetID())
				//suite.Commit()
			},
			false,
		},
		{
			"metadata different that stored",
			func() {
				metadata.Base = cosmosTokenBase
				validMetadata := banktypes.Metadata{
					Description: "description",
					// NOTE: Denom units MUST be increasing
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    cosmosTokenBase,
							Exponent: uint32(0),
						},
						{
							Denom:    "coin2",
							Exponent: defaultExponent,
						},
					},
					Base:    cosmosTokenBase,
					Display: cosmosTokenBase,
					Name:    erc20Name,
					Symbol:  erc20Symbol,
				}
				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, validMetadata)
			},
			false,
		},
		{
			"denom already registered alias",
			func() {
				suite.app.Erc20Keeper.SetAliasesDenom(suite.ctx, "usdt", metadata.Base)
			},
			false,
		},
		{
			"ok",
			func() {
				metadata.Base = cosmosTokenBase
			},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, metadata)
			//suite.Commit()

			contractAddr := crypto.CreateAddress(types.ModuleAddress, suite.app.EvmKeeper.GetNonce(suite.ctx, types.ModuleAddress)-1)
			expPair := &types.TokenPair{
				Erc20Address:  contractAddr.String(),
				Denom:         cosmosTokenBase,
				Enabled:       true,
				ContractOwner: 1,
			}

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().Equal(pair, expPair)
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestRegisterCoinWithManyToOne() {

	metadata := banktypes.Metadata{
		Description: "The cross chain token of the Function X",
		// NOTE: Denom units MUST be increasing
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "usdt",
				Exponent: uint32(0),
				Aliases:  []string{tronDenom, polygonDenom},
			},
			{
				Denom:    "usdtd",
				Exponent: 0,
			},
			{
				Denom:    "USDT",
				Exponent: 18,
			},
		},
		Base:    "usdt",
		Display: "usdtd",
		Name:    "Tether USD",
		Symbol:  "USDT",
	}

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
				suite.app.Erc20Keeper.SetParams(suite.ctx, params)
			},
			false,
			"registration is currently disabled by governance: erc20 module is disabled",
		},
		{
			"denom already registered",
			func() {
				regPair := types.NewTokenPair(helpers.GenerateAddress(), metadata.Base, true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.SetDenomMap(suite.ctx, regPair.Denom, regPair.GetID())
			},
			false,
			"coin denomination already registered: usdt: token pair already exists",
		},
		{
			"alias already registered denom",
			func() {
				regPair := types.NewTokenPair(helpers.GenerateAddress(), tronDenom, true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.SetDenomMap(suite.ctx, regPair.Denom, regPair.GetID())
			},
			false,
			"denom tronTR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t already registered: invalid metadata",
		},
		{
			"alias already registered",
			func() {
				suite.app.Erc20Keeper.SetAliasesDenom(suite.ctx, metadata.Base, tronDenom)
			},
			false,
			"alias tronTR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t already registered: invalid metadata",
		},
		{
			"alias equal base",
			func() {
				metadata.DenomUnits[0].Aliases = []string{tronDenom, polygonDenom, "usdt"}
			},
			false,
			"alias can not equal base, display or symbol: invalid metadata",
		},
		{
			"alias equal display",
			func() {
				metadata.DenomUnits[0].Aliases = []string{tronDenom, polygonDenom, "usdtd"}
			},
			false,
			"alias can not equal base, display or symbol: invalid metadata",
		},
		{
			"alias equal symbol",
			func() {
				metadata.DenomUnits[0].Aliases = []string{tronDenom, polygonDenom, "USDT"}
			},
			false,
			"alias can not equal base, display or symbol: invalid metadata",
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			pair, tcErr := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, metadata)

			contractAddr := crypto.CreateAddress(types.ModuleAddress, suite.app.EvmKeeper.GetNonce(suite.ctx, types.ModuleAddress)-1)
			expPair := &types.TokenPair{
				Erc20Address:  contractAddr.String(),
				Denom:         "usdt",
				Enabled:       true,
				ContractOwner: 1,
			}

			if tc.expPass {
				suite.Require().NoError(tcErr, tc.name)
				suite.Require().Equal(pair, expPair)
				suite.Require().True(suite.app.Erc20Keeper.IsAliasDenomRegistered(suite.ctx, tronDenom))
				suite.Require().True(suite.app.Erc20Keeper.IsAliasDenomRegistered(suite.ctx, polygonDenom))
			} else {
				suite.Require().Error(tcErr, tc.name)
				suite.Require().EqualError(tcErr, tc.errMsg, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestUpdateDenomAlias() {

	metadata := banktypes.Metadata{
		Description: "The cross chain token of the Function X",
		// NOTE: Denom units MUST be increasing
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "usdt",
				Exponent: uint32(0),
				Aliases:  []string{tronDenom, polygonDenom},
			},
			{
				Denom:    "usdtd",
				Exponent: 0,
			},
			{
				Denom:    "USDT",
				Exponent: 18,
			},
		},
		Base:    "usdt",
		Display: "usdtd",
		Name:    "Tether USD",
		Symbol:  "USDT",
	}

	testCases := []struct {
		name     string
		malleate func() error
		expPass  bool
		alias    []string
		errMsg   string
	}{
		{
			name: "success - add alias",
			malleate: func() error {
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "usdt", bscDenom)
				if err != nil {
					return err
				}
				suite.Require().True(addAlias)

				return nil
			},
			expPass: true,
			alias:   []string{tronDenom, polygonDenom, bscDenom},
			errMsg:  "",
		},
		{
			name: "success - delete alias",
			malleate: func() error {
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "usdt", polygonDenom)
				if err != nil {
					return err
				}
				suite.Require().False(addAlias)
				return nil
			},
			expPass: true,
			alias:   []string{tronDenom},
			errMsg:  "",
		},
		{
			name: "failed - denom not equal",
			malleate: func() error {
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "abc", polygonDenom)
				if err != nil {
					return err
				}
				suite.Require().True(addAlias)
				return nil
			},
			expPass: false,
			alias:   []string{},
			errMsg:  "",
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
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "abc", polygonDenom)
				if err != nil {
					return err
				}
				suite.Require().True(addAlias)
				return nil
			},
			expPass: false,
			alias:   []string{},
			errMsg:  "",
		},
		{
			name: "failed - metadata not found",
			malleate: func() error {
				suite.app.Erc20Keeper.SetDenomMap(suite.ctx, "abc", []byte{})

				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "abc", bscDenom)
				if err != nil {
					return err
				}
				suite.Require().True(addAlias)
				return nil
			},
			expPass: false,
			alias:   []string{},
			errMsg:  "",
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
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "abc", bscDenom)
				if err != nil {
					return err
				}
				suite.Require().True(addAlias)
				return nil
			},
			expPass: false,
			alias:   []string{},
			errMsg:  "",
		},
		{
			name: "failed - aliases can not empty",
			malleate: func() error {
				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "abc", tronDenom)
				if err != nil {
					return err
				}
				suite.Require().True(addAlias)
				return nil
			},
			expPass: false,
			alias:   []string{},
			errMsg:  "",
		},
		{
			name: "failed - alias denom not equal with update denom",
			malleate: func() error {
				suite.app.Erc20Keeper.SetAliasesDenom(suite.ctx, "abc", bscDenom)

				addAlias, err := suite.app.Erc20Keeper.UpdateDenomAlias(suite.ctx, "usdt", bscDenom)
				if err != nil {
					return err
				}
				suite.Require().True(addAlias)
				return nil
			},
			expPass: false,
			alias:   []string{},
			errMsg:  "",
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, metadata)
			suite.Require().NoError(err)

			tcErr := tc.malleate()

			if tc.expPass {
				suite.Require().NoError(tcErr, tc.name)
				md, found := suite.app.Erc20Keeper.HasDenomAlias(suite.ctx, pair.Denom)
				suite.Require().True(found)
				suite.Require().Equal(md.DenomUnits[0].Aliases, tc.alias)
				for _, alias := range tc.alias {
					aliasRegistered := suite.app.Erc20Keeper.IsAliasDenomRegistered(suite.ctx, alias)
					suite.Require().True(aliasRegistered)
				}
				if len(metadata.DenomUnits[0].Aliases) > len(tc.alias) {
					for _, alias := range metadata.DenomUnits[0].Aliases[len(tc.alias):] {
						aliasRegistered := suite.app.Erc20Keeper.IsAliasDenomRegistered(suite.ctx, alias)
						suite.Require().False(aliasRegistered)
					}
				}
			} else {
				suite.Require().Error(tcErr, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestRegisterERC20() {
	var (
		contractAddr common.Address
		pair         types.TokenPair
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"intrarelaying is disabled globally",
			func() {
				params := types.DefaultParams()
				params.EnableErc20 = false
				suite.app.Erc20Keeper.SetParams(suite.ctx, params)
			},
			false,
		},
		{
			"token ERC20 already registered",
			func() {
				suite.app.Erc20Keeper.SetERC20Map(suite.ctx, pair.GetERC20Contract(), pair.GetID())
			},
			false,
		},
		{
			"denom already registered",
			func() {
				suite.app.Erc20Keeper.SetDenomMap(suite.ctx, pair.Denom, pair.GetID())
			},
			false,
		},
		{
			"meta data already stored",
			func() {
				_, _, _, err := suite.app.Erc20Keeper.CreateCoinMetadata(suite.ctx, contractAddr)
				suite.Require().NoError(err)
			},
			false,
		},
		{
			"ok",
			func() {},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			var err error
			contractAddr, err = suite.DeployContract(suite.address, erc20Name, erc20Symbol, erc20Decimals)
			suite.Require().NoError(err)
			//suite.Commit()
			coinName := strings.ToLower(erc20Symbol)
			pair = types.NewTokenPair(contractAddr, coinName, true, types.OWNER_EXTERNAL)

			tc.malleate()

			_, err = suite.app.Erc20Keeper.RegisterERC20(suite.ctx, contractAddr)
			metadata, _ := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, coinName)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				// Metadata variables
				suite.Require().Equal(coinName, metadata.Base)
				suite.Require().Equal(coinName, metadata.Display)
				// Denom units
				suite.Require().Equal(len(metadata.DenomUnits), 2)
				suite.Require().Equal(coinName, metadata.DenomUnits[0].Denom)
				suite.Require().Equal(uint32(erc20Decimals), metadata.DenomUnits[1].Exponent)
				suite.Require().Equal(erc20Symbol, metadata.DenomUnits[1].Denom)
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestToggleRelay() {
	var (
		contractAddr common.Address
		id           []byte
		pair         types.TokenPair
	)

	testCases := []struct {
		name         string
		malleate     func()
		expPass      bool
		relayEnabled bool
	}{
		{
			"token not registered",
			func() {
				contractAddr, err := suite.DeployContract(suite.address, erc20Name, erc20Symbol, erc20Decimals)
				suite.Require().NoError(err)
				//suite.Commit()
				pair = types.NewTokenPair(contractAddr, cosmosTokenBase, true, types.OWNER_MODULE)
			},
			false,
			false,
		},
		{
			"token not registered - pair not found",
			func() {
				contractAddr, err := suite.DeployContract(suite.address, erc20Name, erc20Symbol, erc20Decimals)
				suite.Require().NoError(err)
				//suite.Commit()
				pair = types.NewTokenPair(contractAddr, cosmosTokenBase, true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.SetERC20Map(suite.ctx, common.HexToAddress(pair.Erc20Address), pair.GetID())
			},
			false,
			false,
		},
		{
			"disable relay",
			func() {
				contractAddr = suite.setupRegisterERC20Pair(contractMinterBurner)
				id = suite.app.Erc20Keeper.GetTokenPairID(suite.ctx, contractAddr.String())
				pair, _ = suite.app.Erc20Keeper.GetTokenPair(suite.ctx, id)
			},
			true,
			false,
		},
		{
			"disable and enable relay",
			func() {
				contractAddr = suite.setupRegisterERC20Pair(contractMinterBurner)
				id = suite.app.Erc20Keeper.GetTokenPairID(suite.ctx, contractAddr.String())
				pair, _ = suite.app.Erc20Keeper.GetTokenPair(suite.ctx, id)
				pair, _ = suite.app.Erc20Keeper.ToggleRelay(suite.ctx, contractAddr.String())
			},
			true,
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			var err error
			pair, err = suite.app.Erc20Keeper.ToggleRelay(suite.ctx, contractAddr.String())
			// Request the pair using the GetPairToken func to make sure that is updated on the db
			pair, _ = suite.app.Erc20Keeper.GetTokenPair(suite.ctx, id)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				if tc.relayEnabled {
					suite.Require().True(pair.Enabled)
				} else {
					suite.Require().False(pair.Enabled)
				}
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

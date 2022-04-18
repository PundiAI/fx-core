package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/suite"
)

type ParamsTestSuite struct {
	suite.Suite
}

func TestParamsTestSuite(t *testing.T) {
	suite.Run(t, new(ParamsTestSuite))
}

func (suite *ParamsTestSuite) TestParamKeyTable() {
	suite.Require().IsType(paramtypes.KeyTable{}, ParamKeyTable())
}

func (suite *ParamsTestSuite) TestParamsValidate() {
	testCases := []struct {
		name     string
		params   Params
		expError bool
	}{
		{"default", DefaultParams(), false},
		{
			"valid",
			NewParams(7, 3, 2000000000, MinBaseFee, MaxBaseFee, 50000000),
			false,
		},
		{
			"empty",
			Params{},
			true,
		},
		{
			"base fee change denominator is 0 ",
			NewParams(0, 3, 2000000000, MinBaseFee, MaxBaseFee, 0),
			true,
		},
		{
			"if max base fee gt 0, max base fee(1) must be gte min base fee(100)",
			NewParams(0, 3, 2000000000, sdk.NewInt(100), sdk.NewInt(1), 0),
			true,
		},
		{
			"valid",
			NewParams(7, 3, 2000000000, MinBaseFee, sdk.ZeroInt(), 50000000),
			false,
		},
		{
			"valid",
			NewParams(7, 3, 2000000000, sdk.NewInt(100), sdk.ZeroInt(), 50000000),
			false,
		},
	}

	for _, tc := range testCases {
		err := tc.params.Validate()

		if tc.expError {
			suite.Require().Error(err, tc.name)
		} else {
			suite.Require().NoError(err, tc.name)
		}
	}
}

func (suite *ParamsTestSuite) TestParamsValidatePriv() {
	suite.Require().Error(validateBaseFeeChangeDenominator(0))
	suite.Require().Error(validateBaseFeeChangeDenominator(uint32(0)))
	suite.Require().NoError(validateBaseFeeChangeDenominator(uint32(7)))
	suite.Require().Error(validateElasticityMultiplier(""))
	suite.Require().NoError(validateElasticityMultiplier(uint32(2)))
	suite.Require().Error(validateBaseFee(""))
	suite.Require().Error(validateBaseFee(int64(2000000000)))
	suite.Require().Error(validateBaseFee(sdk.NewInt(-2000000000)))
	suite.Require().NoError(validateBaseFee(sdk.NewInt(2000000000)))
	suite.Require().Error(validateMaxGas(""))
	suite.Require().Error(validateMaxGas(int64(2000000000)))
	suite.Require().Error(validateMaxGas(sdk.NewInt(-2000000000)))
	suite.Require().NoError(validateMaxGas(sdk.NewInt(2000000000)))
}

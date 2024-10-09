package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
	evmtypes "github.com/functionx/fx-core/v8/x/evm/types"
	govtypes "github.com/functionx/fx-core/v8/x/gov/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryCustomParams() {
	testCases := []struct {
		name     string
		msgUrl   string
		malleate func(expect govtypes.CustomParams) govtypes.CustomParams
	}{
		{
			name:   "distribution MsgCommunityPoolSpend",
			msgUrl: sdk.MsgTypeURL(&distributiontypes.MsgCommunityPoolSpend{}),
			malleate: func(expect govtypes.CustomParams) govtypes.CustomParams {
				expect.DepositRatio = govtypes.EGFCustomParamDepositRatio.String()
				expect.VotingPeriod = &govtypes.DefaultEGFCustomParamVotingPeriod
				return expect
			},
		},
		{
			name:   "evm MsgCallContract",
			msgUrl: sdk.MsgTypeURL(&evmtypes.MsgCallContract{}),
			malleate: func(expect govtypes.CustomParams) govtypes.CustomParams {
				expect.Quorum = govtypes.DefaultCustomParamQuorum25.String()
				return expect
			},
		},
		{
			name:   "erc20 MsgRegisterCoin",
			msgUrl: sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{}),
			malleate: func(expect govtypes.CustomParams) govtypes.CustomParams {
				expect.Quorum = govtypes.DefaultCustomParamQuorum25.String()
				return expect
			},
		},
		{
			name:   "erc20 MsgRegisterERC20",
			msgUrl: sdk.MsgTypeURL(&erc20types.MsgRegisterERC20{}),
			malleate: func(expect govtypes.CustomParams) govtypes.CustomParams {
				expect.Quorum = govtypes.DefaultCustomParamQuorum25.String()
				return expect
			},
		},
		{
			name:   "erc20 MsgToggleTokenConversion",
			msgUrl: sdk.MsgTypeURL(&erc20types.MsgToggleTokenConversion{}),
			malleate: func(expect govtypes.CustomParams) govtypes.CustomParams {
				expect.Quorum = govtypes.DefaultCustomParamQuorum25.String()
				return expect
			},
		},
		{
			name:   "erc20 MsgUpdateDenomAlias",
			msgUrl: sdk.MsgTypeURL(&erc20types.MsgUpdateDenomAlias{}),
			malleate: func(expect govtypes.CustomParams) govtypes.CustomParams {
				expect.Quorum = govtypes.DefaultCustomParamQuorum25.String()
				return expect
			},
		},
		{
			name:   "gov MsgUpdateSwitchParams",
			msgUrl: sdk.MsgTypeURL(&govtypes.MsgUpdateSwitchParams{}),
			malleate: func(expect govtypes.CustomParams) govtypes.CustomParams {
				appNewGenesisVotingPeriod := time.Hour * 24 * 14
				expect.VotingPeriod = &appNewGenesisVotingPeriod
				return expect
			},
		},
	}

	for _, tc := range testCases {
		expectParams := govtypes.CustomParams{
			DepositRatio: govtypes.DefaultCustomParamDepositRatio.String(),
			VotingPeriod: &govtypes.DefaultCustomParamVotingPeriod,
			Quorum:       govtypes.DefaultCustomParamQuorum40.String(),
		}
		suite.Run(tc.name, func() {
			expect := tc.malleate(expectParams)

			actual, err := suite.queryClient.CustomParams(suite.Ctx, &govtypes.QueryCustomParamsRequest{MsgUrl: tc.msgUrl})
			suite.Require().NoError(err)
			suite.Require().Equal(expect, actual.GetParams())
		})
	}
}

package helpers

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (s *BaseSuite) SetOracle(chainName string, online bool) crosschaintypes.Oracle {
	oracle := crosschaintypes.Oracle{
		OracleAddress:     GenAccAddress().String(),
		BridgerAddress:    GenAccAddress().String(),
		ExternalAddress:   GenExternalAddr(chainName),
		DelegateAmount:    sdkmath.NewInt(1e18).MulRaw(1000),
		StartHeight:       1,
		Online:            online,
		DelegateValidator: sdk.ValAddress(GenAccAddress()).String(),
		SlashTimes:        0,
	}
	keeper := s.App.CrosschainKeepers.GetKeeper(chainName)
	keeper.SetOracle(s.Ctx, oracle)
	keeper.SetOracleAddrByExternalAddr(s.Ctx, oracle.ExternalAddress, oracle.GetOracle())
	keeper.SetOracleAddrByBridgerAddr(s.Ctx, oracle.GetBridger(), oracle.GetOracle())
	return oracle
}

func (s *BaseSuite) GetERC20Token(baseDenom string) *erc20types.ERC20Token {
	erc20token, err := s.App.Erc20Keeper.GetERC20Token(s.Ctx, baseDenom)
	s.Require().NoError(err)
	return &erc20token
}

func (s *BaseSuite) GetBridgeToken(chainName, baseDenom string) erc20types.BridgeToken {
	bridgeToken, err := s.App.Erc20Keeper.GetBridgeToken(s.Ctx, chainName, baseDenom)
	s.Require().NoError(err)
	return bridgeToken
}

func (s *BaseSuite) AddBridgeToken(chainName, symbolOrAddr string, isNativeCoin bool, isIBC ...bool) erc20types.BridgeToken {
	keeper := s.App.Erc20Keeper
	var baseDenom string
	isNative := false
	if symbolOrAddr == fxtypes.LegacyFXDenom {
		baseDenom = fxtypes.FXDenom
	} else if isNativeCoin || symbolOrAddr == fxtypes.DefaultSymbol {
		erc20Token, err := keeper.RegisterNativeCoin(s.Ctx, symbolOrAddr, symbolOrAddr, 18)
		s.Require().NoError(err)
		baseDenom = erc20Token.Denom
	} else {
		isNative = true
		erc20Token, err := keeper.RegisterNativeERC20(s.Ctx, common.HexToAddress(symbolOrAddr))
		s.Require().NoError(err)
		baseDenom = erc20Token.Denom
	}
	if len(isIBC) > 0 && isIBC[0] {
		isNative = true
	}
	err := keeper.AddBridgeToken(s.Ctx, baseDenom, chainName, GenExternalAddr(chainName), isNative)
	s.Require().NoError(err)
	return s.GetBridgeToken(chainName, baseDenom)
}

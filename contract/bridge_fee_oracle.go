package contract

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/types"
)

type BridgeFeeOracleKeeper struct {
	Caller
	abi      abi.ABI
	from     common.Address
	contract common.Address
}

func NewBridgeFeeOracleKeeper(caller Caller, contract string) BridgeFeeOracleKeeper {
	return BridgeFeeOracleKeeper{
		Caller:   caller,
		abi:      GetBridgeFeeOracle().ABI,
		from:     common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes()),
		contract: common.HexToAddress(contract),
	}
}

func (k BridgeFeeOracleKeeper) Initialize(ctx context.Context) (*types.MsgEthereumTxResponse, error) {
	return k.Caller.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "initialize", common.HexToAddress(CrosschainAddress))
}

func (k BridgeFeeOracleKeeper) GetOwnerRole(ctx context.Context) (common.Hash, error) {
	var res struct{ Role common.Hash }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, k.contract, k.abi, "OWNER_ROLE", &res); err != nil {
		return common.Hash{}, err
	}
	return res.Role, nil
}

func (k BridgeFeeOracleKeeper) GetUpgradeRole(ctx context.Context) (common.Hash, error) {
	var res struct{ Role common.Hash }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, k.contract, k.abi, "UPGRADE_ROLE", &res); err != nil {
		return common.Hash{}, err
	}
	return res.Role, nil
}

func (k BridgeFeeOracleKeeper) GetQuoteRole(ctx context.Context) (common.Hash, error) {
	var res struct{ Role common.Hash }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, k.contract, k.abi, "QUOTE_ROLE", &res); err != nil {
		return common.Hash{}, err
	}
	return res.Role, nil
}

func (k BridgeFeeOracleKeeper) GrantRole(ctx context.Context, role common.Hash, account common.Address) (*types.MsgEthereumTxResponse, error) {
	return k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "grantRole", role, account)
}

func (k BridgeFeeOracleKeeper) SetDefaultOracle(ctx context.Context, oracle common.Address) (*types.MsgEthereumTxResponse, error) {
	return k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "setDefaultOracle", oracle)
}

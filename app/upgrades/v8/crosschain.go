package v8

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/app/keepers"
	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschainkeeper "github.com/pundiai/fx-core/v8/x/crosschain/keeper"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20 "github.com/pundiai/fx-core/v8/x/erc20/types"
	layer2types "github.com/pundiai/fx-core/v8/x/layer2/types"
)

func migrateCrosschainParams(ctx sdk.Context, keepers keepers.CrosschainKeepers) error {
	for _, k := range keepers.ToSlice() {
		params := k.GetParams(ctx)
		params.DelegateThreshold.Denom = fxtypes.DefaultDenom
		params.DelegateThreshold.Amount = fxtypes.SwapAmount(params.DelegateThreshold.Amount)
		if !params.DelegateThreshold.IsPositive() {
			return sdkerrors.ErrInvalidCoins.Wrapf("module %s invalid delegate threshold: %s",
				k.ModuleName(), params.DelegateThreshold.String())
		}
		if err := k.SetParams(ctx, &params); err != nil {
			return err
		}
	}
	return nil
}

func migrateCrosschainModuleAccount(ctx sdk.Context, ak authkeeper.AccountKeeper) error {
	addr, perms := ak.GetModuleAddressAndPermissions(crosschaintypes.ModuleName)
	if addr == nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("crosschain module empty permissions")
	}
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("crosschain account not exist")
	}
	baseAcc, ok := acc.(*authtypes.BaseAccount)
	if !ok {
		return sdkerrors.ErrInvalidAddress.Wrapf("crosschain account not base account")
	}
	macc := authtypes.NewModuleAccount(baseAcc, crosschaintypes.ModuleName, perms...)
	ak.SetModuleAccount(ctx, macc)
	return nil
}

func migrateOracleDelegateAmount(ctx sdk.Context, keepers keepers.CrosschainKeepers) {
	for _, k := range keepers.ToSlice() {
		k.IterateOracle(ctx, func(oracle crosschaintypes.Oracle) bool {
			oracle.DelegateAmount = fxtypes.SwapAmount(oracle.DelegateAmount)
			k.SetOracle(ctx, oracle)
			return false
		})
	}
}

func initBridgeAccount(ctx sdk.Context, ak authkeeper.AccountKeeper) {
	bridgeFeeCollector := authtypes.NewModuleAddress(crosschaintypes.BridgeFeeCollectorName)
	if account := ak.GetAccount(ctx, bridgeFeeCollector); account == nil {
		ak.SetAccount(ctx, ak.NewAccountWithAddress(ctx, bridgeFeeCollector))
	}
	bridgeCallFrom := authtypes.NewModuleAddress(crosschaintypes.BridgeCallSender)
	if account := ak.GetAccount(ctx, bridgeCallFrom); account == nil {
		ak.SetAccount(ctx, ak.NewAccountWithAddress(ctx, bridgeCallFrom))
	}
}

func cancelOutgoingTxPool(ctx sdk.Context, keepers keepers.CrosschainKeepers, erc20Keeper crosschaintypes.Erc20Keeper, evmKeeper crosschaintypes.EVMKeeper, bankKeeper crosschaintypes.BankKeeper, accountKeeper crosschaintypes.AccountKeeper) (err error) {
	for _, k := range keepers.ToSlice() {
		batchHandleTxs := make([]crosschaintypes.OutgoingTransferTx, 0)
		k.IterateOutgoingTxPool(ctx, func(tx crosschaintypes.OutgoingTransferTx) bool {
			if k.ModuleName() == layer2types.ModuleName && strings.EqualFold(tx.Token.Contract, "0x0fee9ab6B068385B99992fa9C349b40Bd8E392D1") {
				batchHandleTxs = append(batchHandleTxs, tx)
				return false
			}
			if err = handleCancel(ctx, k, erc20Keeper, evmKeeper, tx); err != nil {
				return true
			}
			return false
		})
		if err != nil {
			return err
		}
		if len(batchHandleTxs) == 0 {
			continue
		}
		cacheCtx, commit := ctx.CacheContext()
		if handleErr := handleLayer2Txs(cacheCtx, k, erc20Keeper, evmKeeper, bankKeeper, accountKeeper, batchHandleTxs); handleErr != nil {
			ctx.Logger().Error("handle layer2 txs error", "err", handleErr, "module", k.ModuleName(), "txs", len(batchHandleTxs))
			continue
		}
		commit()
	}
	return err
}

func handleLayer2Txs(ctx sdk.Context, k crosschainkeeper.Keeper, erc20Keeper crosschaintypes.Erc20Keeper, evmKeeper crosschaintypes.EVMKeeper, bankKeeper crosschaintypes.BankKeeper, accountKeeper crosschaintypes.AccountKeeper, batchHandleTxs []crosschaintypes.OutgoingTransferTx) error {
	// calculate total refund amount
	totalRefundAmount := sdkmath.NewInt(0)
	for _, tx := range batchHandleTxs {
		totalRefundAmount = totalRefundAmount.Add(tx.Token.Amount).Add(tx.Fee.Amount)
	}

	// get bridge token
	bridgeToken, err := k.GetBridgeToken(ctx, batchHandleTxs[0].Token.Contract)
	if err != nil {
		return err
	}
	// get erc20 token
	erc20Token, err := erc20Keeper.GetERC20Token(ctx, bridgeToken.Denom)
	if err != nil {
		return err
	}

	erc20Contract := contract.NewERC20TokenKeeper(evmKeeper)
	erc20ModuleAccount := accountKeeper.GetModuleAddress(erc20.ModuleName)
	balanceOf, err := erc20Contract.BalanceOf(ctx, common.HexToAddress(erc20Token.Erc20Address), common.BytesToAddress(erc20ModuleAccount.Bytes()))
	if err != nil {
		return err
	}
	// check erc20 balance is enough
	erc20HasEnough := balanceOf.Cmp(totalRefundAmount.BigInt()) >= 0
	if !erc20HasEnough {
		ctx.Logger().Info("not enough erc20 balance", "balance", balanceOf.String(), "totalRefundAmount", totalRefundAmount.String(), "diff", totalRefundAmount.BigInt().Sub(totalRefundAmount.BigInt(), balanceOf).String())
		return nil
	}

	if err = mintCoinToModule(ctx, k.ModuleName(), accountKeeper, bridgeToken, bankKeeper, totalRefundAmount); err != nil {
		return err
	}

	// handle cancel txs
	for _, tx := range batchHandleTxs {
		if err = handleCancel(ctx, k, erc20Keeper, evmKeeper, tx); err != nil {
			return err
		}
	}
	return nil
}

func mintCoinToModule(ctx sdk.Context, moduleName string, accountKeeper crosschaintypes.AccountKeeper, bridgeToken erc20.BridgeToken, bankKeeper crosschaintypes.BankKeeper, totalRefundAmount sdkmath.Int) error {
	// check crosschain module and k.moduleName has enough balance
	keeperModuleAddr := accountKeeper.GetModuleAddress(moduleName)
	crosschainModuleAddr := accountKeeper.GetModuleAddress(crosschaintypes.ModuleName)
	bridgeDenom := crosschaintypes.NewBridgeDenom(moduleName, bridgeToken.Contract)
	if crosschainModuleAddr == nil {
		return sdkerrors.ErrInvalidAddress.Wrap("crosschain module account not exist")
	}
	if keeperModuleAddr == nil {
		return sdkerrors.ErrInvalidAddress.Wrap(moduleName + " module account not exist")
	}
	bridgeDenomBalance := bankKeeper.GetBalance(ctx, keeperModuleAddr, bridgeDenom)
	baseDenomBalance := bankKeeper.GetBalance(ctx, crosschainModuleAddr, bridgeToken.Denom)
	if bridgeDenomBalance.Amount.LT(totalRefundAmount) {
		diff := totalRefundAmount.Sub(bridgeDenomBalance.Amount)
		if err := bankKeeper.MintCoins(ctx, moduleName, sdk.NewCoins(sdk.NewCoin(bridgeDenom, diff))); err != nil {
			return err
		}
	}
	if baseDenomBalance.Amount.LT(totalRefundAmount) {
		diff := totalRefundAmount.Sub(bridgeDenomBalance.Amount)
		if err := bankKeeper.MintCoins(ctx, crosschaintypes.ModuleName, sdk.NewCoins(sdk.NewCoin(bridgeToken.Denom, diff))); err != nil {
			return err
		}
	}
	return nil
}

func handleCancel(ctx sdk.Context, k crosschainkeeper.Keeper, erc20Keeper crosschaintypes.Erc20Keeper, evmKeeper crosschaintypes.EVMKeeper, transferTx crosschaintypes.OutgoingTransferTx) error {
	k.DeleteOutgoingTxPool(ctx, transferTx.Fee, transferTx.Id)

	sender, err := sdk.AccAddressFromBech32(transferTx.Sender)
	if err != nil {
		return err
	}
	baseCoin, err := k.DepositBridgeTokenToBaseCoin(ctx, sender, transferTx.Token.Amount.Add(transferTx.Fee.Amount), transferTx.Token.Contract)
	if err != nil {
		return err
	}
	if baseCoin.Denom == fxtypes.DefaultDenom {
		return nil
	}
	_, err = erc20Keeper.BaseCoinToEvm(ctx, evmKeeper, common.BytesToAddress(sender.Bytes()), baseCoin)
	return err
}

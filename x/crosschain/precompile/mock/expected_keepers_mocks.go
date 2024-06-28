// Code generated by MockGen. DO NOT EDIT.
// Source: x/crosschain/precompile/expected_keepers.go
//
// Generated by this command:
//
//	mockgen -source=x/crosschain/precompile/expected_keepers.go -package mock -destination x/crosschain/precompile/mock/expected_keepers_mocks.go
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	time "time"

	types "github.com/cosmos/cosmos-sdk/types"
	types0 "github.com/cosmos/cosmos-sdk/x/bank/types"
	types1 "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	common "github.com/ethereum/go-ethereum/common"
	types2 "github.com/functionx/fx-core/v7/types"
	types3 "github.com/functionx/fx-core/v7/x/erc20/types"
	gomock "go.uber.org/mock/gomock"
)

// MockErc20Keeper is a mock of Erc20Keeper interface.
type MockErc20Keeper struct {
	ctrl     *gomock.Controller
	recorder *MockErc20KeeperMockRecorder
}

// MockErc20KeeperMockRecorder is the mock recorder for MockErc20Keeper.
type MockErc20KeeperMockRecorder struct {
	mock *MockErc20Keeper
}

// NewMockErc20Keeper creates a new mock instance.
func NewMockErc20Keeper(ctrl *gomock.Controller) *MockErc20Keeper {
	mock := &MockErc20Keeper{ctrl: ctrl}
	mock.recorder = &MockErc20KeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockErc20Keeper) EXPECT() *MockErc20KeeperMockRecorder {
	return m.recorder
}

// ConvertDenomToTarget mocks base method.
func (m *MockErc20Keeper) ConvertDenomToTarget(ctx types.Context, from types.AccAddress, coin types.Coin, fxTarget types2.FxTarget) (types.Coin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConvertDenomToTarget", ctx, from, coin, fxTarget)
	ret0, _ := ret[0].(types.Coin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConvertDenomToTarget indicates an expected call of ConvertDenomToTarget.
func (mr *MockErc20KeeperMockRecorder) ConvertDenomToTarget(ctx, from, coin, fxTarget any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConvertDenomToTarget", reflect.TypeOf((*MockErc20Keeper)(nil).ConvertDenomToTarget), ctx, from, coin, fxTarget)
}

// GetIbcTimeout mocks base method.
func (m *MockErc20Keeper) GetIbcTimeout(ctx types.Context) time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIbcTimeout", ctx)
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// GetIbcTimeout indicates an expected call of GetIbcTimeout.
func (mr *MockErc20KeeperMockRecorder) GetIbcTimeout(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIbcTimeout", reflect.TypeOf((*MockErc20Keeper)(nil).GetIbcTimeout), ctx)
}

// GetTokenPair mocks base method.
func (m *MockErc20Keeper) GetTokenPair(ctx types.Context, tokenOrDenom string) (types3.TokenPair, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenPair", ctx, tokenOrDenom)
	ret0, _ := ret[0].(types3.TokenPair)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetTokenPair indicates an expected call of GetTokenPair.
func (mr *MockErc20KeeperMockRecorder) GetTokenPair(ctx, tokenOrDenom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenPair", reflect.TypeOf((*MockErc20Keeper)(nil).GetTokenPair), ctx, tokenOrDenom)
}

// GetTokenPairByAddress mocks base method.
func (m *MockErc20Keeper) GetTokenPairByAddress(ctx types.Context, address common.Address) (types3.TokenPair, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenPairByAddress", ctx, address)
	ret0, _ := ret[0].(types3.TokenPair)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetTokenPairByAddress indicates an expected call of GetTokenPairByAddress.
func (mr *MockErc20KeeperMockRecorder) GetTokenPairByAddress(ctx, address any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenPairByAddress", reflect.TypeOf((*MockErc20Keeper)(nil).GetTokenPairByAddress), ctx, address)
}

// HasDenomAlias mocks base method.
func (m *MockErc20Keeper) HasDenomAlias(ctx types.Context, denom string) (types0.Metadata, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasDenomAlias", ctx, denom)
	ret0, _ := ret[0].(types0.Metadata)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// HasDenomAlias indicates an expected call of HasDenomAlias.
func (mr *MockErc20KeeperMockRecorder) HasDenomAlias(ctx, denom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasDenomAlias", reflect.TypeOf((*MockErc20Keeper)(nil).HasDenomAlias), ctx, denom)
}

// HasOutgoingTransferRelation mocks base method.
func (m *MockErc20Keeper) HasOutgoingTransferRelation(ctx types.Context, moduleName string, txID uint64) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasOutgoingTransferRelation", ctx, moduleName, txID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasOutgoingTransferRelation indicates an expected call of HasOutgoingTransferRelation.
func (mr *MockErc20KeeperMockRecorder) HasOutgoingTransferRelation(ctx, moduleName, txID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasOutgoingTransferRelation", reflect.TypeOf((*MockErc20Keeper)(nil).HasOutgoingTransferRelation), ctx, moduleName, txID)
}

// IsOriginDenom mocks base method.
func (m *MockErc20Keeper) IsOriginDenom(ctx types.Context, denom string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsOriginDenom", ctx, denom)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsOriginDenom indicates an expected call of IsOriginDenom.
func (mr *MockErc20KeeperMockRecorder) IsOriginDenom(ctx, denom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsOriginDenom", reflect.TypeOf((*MockErc20Keeper)(nil).IsOriginDenom), ctx, denom)
}

// ModuleAddress mocks base method.
func (m *MockErc20Keeper) ModuleAddress() common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModuleAddress")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// ModuleAddress indicates an expected call of ModuleAddress.
func (mr *MockErc20KeeperMockRecorder) ModuleAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModuleAddress", reflect.TypeOf((*MockErc20Keeper)(nil).ModuleAddress))
}

// SetIBCTransferRelation mocks base method.
func (m *MockErc20Keeper) SetIBCTransferRelation(ctx types.Context, channel string, sequence uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetIBCTransferRelation", ctx, channel, sequence)
}

// SetIBCTransferRelation indicates an expected call of SetIBCTransferRelation.
func (mr *MockErc20KeeperMockRecorder) SetIBCTransferRelation(ctx, channel, sequence any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetIBCTransferRelation", reflect.TypeOf((*MockErc20Keeper)(nil).SetIBCTransferRelation), ctx, channel, sequence)
}

// ToTargetDenom mocks base method.
func (m *MockErc20Keeper) ToTargetDenom(ctx types.Context, denom, base string, aliases []string, fxTarget types2.FxTarget) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToTargetDenom", ctx, denom, base, aliases, fxTarget)
	ret0, _ := ret[0].(string)
	return ret0
}

// ToTargetDenom indicates an expected call of ToTargetDenom.
func (mr *MockErc20KeeperMockRecorder) ToTargetDenom(ctx, denom, base, aliases, fxTarget any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToTargetDenom", reflect.TypeOf((*MockErc20Keeper)(nil).ToTargetDenom), ctx, denom, base, aliases, fxTarget)
}

// MockBankKeeper is a mock of BankKeeper interface.
type MockBankKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockBankKeeperMockRecorder
}

// MockBankKeeperMockRecorder is the mock recorder for MockBankKeeper.
type MockBankKeeperMockRecorder struct {
	mock *MockBankKeeper
}

// NewMockBankKeeper creates a new mock instance.
func NewMockBankKeeper(ctrl *gomock.Controller) *MockBankKeeper {
	mock := &MockBankKeeper{ctrl: ctrl}
	mock.recorder = &MockBankKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBankKeeper) EXPECT() *MockBankKeeperMockRecorder {
	return m.recorder
}

// BurnCoins mocks base method.
func (m *MockBankKeeper) BurnCoins(ctx types.Context, moduleName string, amounts types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BurnCoins", ctx, moduleName, amounts)
	ret0, _ := ret[0].(error)
	return ret0
}

// BurnCoins indicates an expected call of BurnCoins.
func (mr *MockBankKeeperMockRecorder) BurnCoins(ctx, moduleName, amounts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BurnCoins", reflect.TypeOf((*MockBankKeeper)(nil).BurnCoins), ctx, moduleName, amounts)
}

// GetBalance mocks base method.
func (m *MockBankKeeper) GetBalance(ctx types.Context, addr types.AccAddress, denom string) types.Coin {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", ctx, addr, denom)
	ret0, _ := ret[0].(types.Coin)
	return ret0
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockBankKeeperMockRecorder) GetBalance(ctx, addr, denom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockBankKeeper)(nil).GetBalance), ctx, addr, denom)
}

// GetDenomMetaData mocks base method.
func (m *MockBankKeeper) GetDenomMetaData(ctx types.Context, denom string) (types0.Metadata, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDenomMetaData", ctx, denom)
	ret0, _ := ret[0].(types0.Metadata)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetDenomMetaData indicates an expected call of GetDenomMetaData.
func (mr *MockBankKeeperMockRecorder) GetDenomMetaData(ctx, denom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDenomMetaData", reflect.TypeOf((*MockBankKeeper)(nil).GetDenomMetaData), ctx, denom)
}

// GetSupply mocks base method.
func (m *MockBankKeeper) GetSupply(ctx types.Context, denom string) types.Coin {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSupply", ctx, denom)
	ret0, _ := ret[0].(types.Coin)
	return ret0
}

// GetSupply indicates an expected call of GetSupply.
func (mr *MockBankKeeperMockRecorder) GetSupply(ctx, denom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSupply", reflect.TypeOf((*MockBankKeeper)(nil).GetSupply), ctx, denom)
}

// MintCoins mocks base method.
func (m *MockBankKeeper) MintCoins(ctx types.Context, moduleName string, amounts types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MintCoins", ctx, moduleName, amounts)
	ret0, _ := ret[0].(error)
	return ret0
}

// MintCoins indicates an expected call of MintCoins.
func (mr *MockBankKeeperMockRecorder) MintCoins(ctx, moduleName, amounts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MintCoins", reflect.TypeOf((*MockBankKeeper)(nil).MintCoins), ctx, moduleName, amounts)
}

// SendCoins mocks base method.
func (m *MockBankKeeper) SendCoins(ctx types.Context, fromAddr, toAddr types.AccAddress, amt types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoins", ctx, fromAddr, toAddr, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoins indicates an expected call of SendCoins.
func (mr *MockBankKeeperMockRecorder) SendCoins(ctx, fromAddr, toAddr, amt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoins", reflect.TypeOf((*MockBankKeeper)(nil).SendCoins), ctx, fromAddr, toAddr, amt)
}

// SendCoinsFromAccountToModule mocks base method.
func (m *MockBankKeeper) SendCoinsFromAccountToModule(ctx types.Context, senderAddr types.AccAddress, recipientModule string, amt types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoinsFromAccountToModule", ctx, senderAddr, recipientModule, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoinsFromAccountToModule indicates an expected call of SendCoinsFromAccountToModule.
func (mr *MockBankKeeperMockRecorder) SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoinsFromAccountToModule", reflect.TypeOf((*MockBankKeeper)(nil).SendCoinsFromAccountToModule), ctx, senderAddr, recipientModule, amt)
}

// SendCoinsFromModuleToAccount mocks base method.
func (m *MockBankKeeper) SendCoinsFromModuleToAccount(ctx types.Context, senderModule string, recipientAddr types.AccAddress, amt types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoinsFromModuleToAccount", ctx, senderModule, recipientAddr, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoinsFromModuleToAccount indicates an expected call of SendCoinsFromModuleToAccount.
func (mr *MockBankKeeperMockRecorder) SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoinsFromModuleToAccount", reflect.TypeOf((*MockBankKeeper)(nil).SendCoinsFromModuleToAccount), ctx, senderModule, recipientAddr, amt)
}

// MockIBCTransferKeeper is a mock of IBCTransferKeeper interface.
type MockIBCTransferKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockIBCTransferKeeperMockRecorder
}

// MockIBCTransferKeeperMockRecorder is the mock recorder for MockIBCTransferKeeper.
type MockIBCTransferKeeperMockRecorder struct {
	mock *MockIBCTransferKeeper
}

// NewMockIBCTransferKeeper creates a new mock instance.
func NewMockIBCTransferKeeper(ctrl *gomock.Controller) *MockIBCTransferKeeper {
	mock := &MockIBCTransferKeeper{ctrl: ctrl}
	mock.recorder = &MockIBCTransferKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBCTransferKeeper) EXPECT() *MockIBCTransferKeeperMockRecorder {
	return m.recorder
}

// Transfer mocks base method.
func (m *MockIBCTransferKeeper) Transfer(goCtx context.Context, msg *types1.MsgTransfer) (*types1.MsgTransferResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transfer", goCtx, msg)
	ret0, _ := ret[0].(*types1.MsgTransferResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Transfer indicates an expected call of Transfer.
func (mr *MockIBCTransferKeeperMockRecorder) Transfer(goCtx, msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transfer", reflect.TypeOf((*MockIBCTransferKeeper)(nil).Transfer), goCtx, msg)
}

// MockAccountKeeper is a mock of AccountKeeper interface.
type MockAccountKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockAccountKeeperMockRecorder
}

// MockAccountKeeperMockRecorder is the mock recorder for MockAccountKeeper.
type MockAccountKeeperMockRecorder struct {
	mock *MockAccountKeeper
}

// NewMockAccountKeeper creates a new mock instance.
func NewMockAccountKeeper(ctrl *gomock.Controller) *MockAccountKeeper {
	mock := &MockAccountKeeper{ctrl: ctrl}
	mock.recorder = &MockAccountKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountKeeper) EXPECT() *MockAccountKeeperMockRecorder {
	return m.recorder
}

// GetModuleAddress mocks base method.
func (m *MockAccountKeeper) GetModuleAddress(moduleName string) types.AccAddress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModuleAddress", moduleName)
	ret0, _ := ret[0].(types.AccAddress)
	return ret0
}

// GetModuleAddress indicates an expected call of GetModuleAddress.
func (mr *MockAccountKeeperMockRecorder) GetModuleAddress(moduleName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModuleAddress", reflect.TypeOf((*MockAccountKeeper)(nil).GetModuleAddress), moduleName)
}

// MockCrosschainKeeper is a mock of CrosschainKeeper interface.
type MockCrosschainKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockCrosschainKeeperMockRecorder
}

// MockCrosschainKeeperMockRecorder is the mock recorder for MockCrosschainKeeper.
type MockCrosschainKeeperMockRecorder struct {
	mock *MockCrosschainKeeper
}

// NewMockCrosschainKeeper creates a new mock instance.
func NewMockCrosschainKeeper(ctrl *gomock.Controller) *MockCrosschainKeeper {
	mock := &MockCrosschainKeeper{ctrl: ctrl}
	mock.recorder = &MockCrosschainKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCrosschainKeeper) EXPECT() *MockCrosschainKeeperMockRecorder {
	return m.recorder
}

// PrecompileAddPendingPoolRewards mocks base method.
func (m *MockCrosschainKeeper) PrecompileAddPendingPoolRewards(ctx types.Context, txID uint64, sender types.AccAddress, reward types.Coin) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrecompileAddPendingPoolRewards", ctx, txID, sender, reward)
	ret0, _ := ret[0].(error)
	return ret0
}

// PrecompileAddPendingPoolRewards indicates an expected call of PrecompileAddPendingPoolRewards.
func (mr *MockCrosschainKeeperMockRecorder) PrecompileAddPendingPoolRewards(ctx, txID, sender, reward any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrecompileAddPendingPoolRewards", reflect.TypeOf((*MockCrosschainKeeper)(nil).PrecompileAddPendingPoolRewards), ctx, txID, sender, reward)
}

// PrecompileBridgeCall mocks base method.
func (m *MockCrosschainKeeper) PrecompileBridgeCall(ctx types.Context, sender, refund common.Address, coins types.Coins, to common.Address, data, memo []byte) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrecompileBridgeCall", ctx, sender, refund, coins, to, data, memo)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrecompileBridgeCall indicates an expected call of PrecompileBridgeCall.
func (mr *MockCrosschainKeeperMockRecorder) PrecompileBridgeCall(ctx, sender, refund, coins, to, data, memo any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrecompileBridgeCall", reflect.TypeOf((*MockCrosschainKeeper)(nil).PrecompileBridgeCall), ctx, sender, refund, coins, to, data, memo)
}

// PrecompileCancelPendingBridgeCall mocks base method.
func (m *MockCrosschainKeeper) PrecompileCancelPendingBridgeCall(ctx types.Context, nonce uint64, sender types.AccAddress) (types.Coins, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrecompileCancelPendingBridgeCall", ctx, nonce, sender)
	ret0, _ := ret[0].(types.Coins)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrecompileCancelPendingBridgeCall indicates an expected call of PrecompileCancelPendingBridgeCall.
func (mr *MockCrosschainKeeperMockRecorder) PrecompileCancelPendingBridgeCall(ctx, nonce, sender any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrecompileCancelPendingBridgeCall", reflect.TypeOf((*MockCrosschainKeeper)(nil).PrecompileCancelPendingBridgeCall), ctx, nonce, sender)
}

// PrecompileCancelSendToExternal mocks base method.
func (m *MockCrosschainKeeper) PrecompileCancelSendToExternal(ctx types.Context, txID uint64, sender types.AccAddress) (types.Coin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrecompileCancelSendToExternal", ctx, txID, sender)
	ret0, _ := ret[0].(types.Coin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrecompileCancelSendToExternal indicates an expected call of PrecompileCancelSendToExternal.
func (mr *MockCrosschainKeeperMockRecorder) PrecompileCancelSendToExternal(ctx, txID, sender any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrecompileCancelSendToExternal", reflect.TypeOf((*MockCrosschainKeeper)(nil).PrecompileCancelSendToExternal), ctx, txID, sender)
}

// PrecompileIncreaseBridgeFee mocks base method.
func (m *MockCrosschainKeeper) PrecompileIncreaseBridgeFee(ctx types.Context, txID uint64, sender types.AccAddress, addBridgeFee types.Coin) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrecompileIncreaseBridgeFee", ctx, txID, sender, addBridgeFee)
	ret0, _ := ret[0].(error)
	return ret0
}

// PrecompileIncreaseBridgeFee indicates an expected call of PrecompileIncreaseBridgeFee.
func (mr *MockCrosschainKeeperMockRecorder) PrecompileIncreaseBridgeFee(ctx, txID, sender, addBridgeFee any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrecompileIncreaseBridgeFee", reflect.TypeOf((*MockCrosschainKeeper)(nil).PrecompileIncreaseBridgeFee), ctx, txID, sender, addBridgeFee)
}

// TransferAfter mocks base method.
func (m *MockCrosschainKeeper) TransferAfter(ctx types.Context, sender types.AccAddress, receive string, coins, fee types.Coin, originToken, insufficientLiquidity bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferAfter", ctx, sender, receive, coins, fee, originToken, insufficientLiquidity)
	ret0, _ := ret[0].(error)
	return ret0
}

// TransferAfter indicates an expected call of TransferAfter.
func (mr *MockCrosschainKeeperMockRecorder) TransferAfter(ctx, sender, receive, coins, fee, originToken, insufficientLiquidity any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferAfter", reflect.TypeOf((*MockCrosschainKeeper)(nil).TransferAfter), ctx, sender, receive, coins, fee, originToken, insufficientLiquidity)
}

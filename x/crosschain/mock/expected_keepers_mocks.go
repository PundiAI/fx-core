// Code generated by MockGen. DO NOT EDIT.
// Source: x/crosschain/types/expected_keepers.go
//
// Generated by this command:
//
//	mockgen -source=x/crosschain/types/expected_keepers.go -package mock -destination x/crosschain/mock/expected_keepers_mocks.go
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	big "math/big"
	reflect "reflect"

	types "github.com/cosmos/cosmos-sdk/types"
	types0 "github.com/cosmos/cosmos-sdk/x/auth/types"
	types1 "github.com/cosmos/cosmos-sdk/x/bank/types"
	types2 "github.com/cosmos/cosmos-sdk/x/distribution/types"
	types3 "github.com/cosmos/cosmos-sdk/x/params/types"
	types4 "github.com/cosmos/cosmos-sdk/x/staking/types"
	types5 "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	common "github.com/ethereum/go-ethereum/common"
	types6 "github.com/evmos/ethermint/x/evm/types"
	types7 "github.com/functionx/fx-core/v7/types"
	types8 "github.com/functionx/fx-core/v7/x/crosschain/types"
	types9 "github.com/functionx/fx-core/v7/x/erc20/types"
	gomock "go.uber.org/mock/gomock"
)

// MockStakingKeeper is a mock of StakingKeeper interface.
type MockStakingKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockStakingKeeperMockRecorder
}

// MockStakingKeeperMockRecorder is the mock recorder for MockStakingKeeper.
type MockStakingKeeperMockRecorder struct {
	mock *MockStakingKeeper
}

// NewMockStakingKeeper creates a new mock instance.
func NewMockStakingKeeper(ctrl *gomock.Controller) *MockStakingKeeper {
	mock := &MockStakingKeeper{ctrl: ctrl}
	mock.recorder = &MockStakingKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStakingKeeper) EXPECT() *MockStakingKeeperMockRecorder {
	return m.recorder
}

// AfterDelegationModified mocks base method.
func (m *MockStakingKeeper) AfterDelegationModified(ctx types.Context, delAddr types.AccAddress, valAddr types.ValAddress) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AfterDelegationModified", ctx, delAddr, valAddr)
	ret0, _ := ret[0].(error)
	return ret0
}

// AfterDelegationModified indicates an expected call of AfterDelegationModified.
func (mr *MockStakingKeeperMockRecorder) AfterDelegationModified(ctx, delAddr, valAddr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AfterDelegationModified", reflect.TypeOf((*MockStakingKeeper)(nil).AfterDelegationModified), ctx, delAddr, valAddr)
}

// BeforeDelegationCreated mocks base method.
func (m *MockStakingKeeper) BeforeDelegationCreated(ctx types.Context, delAddr types.AccAddress, valAddr types.ValAddress) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeforeDelegationCreated", ctx, delAddr, valAddr)
	ret0, _ := ret[0].(error)
	return ret0
}

// BeforeDelegationCreated indicates an expected call of BeforeDelegationCreated.
func (mr *MockStakingKeeperMockRecorder) BeforeDelegationCreated(ctx, delAddr, valAddr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeforeDelegationCreated", reflect.TypeOf((*MockStakingKeeper)(nil).BeforeDelegationCreated), ctx, delAddr, valAddr)
}

// GetBondedValidatorsByPower mocks base method.
func (m *MockStakingKeeper) GetBondedValidatorsByPower(ctx types.Context) []types4.Validator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBondedValidatorsByPower", ctx)
	ret0, _ := ret[0].([]types4.Validator)
	return ret0
}

// GetBondedValidatorsByPower indicates an expected call of GetBondedValidatorsByPower.
func (mr *MockStakingKeeperMockRecorder) GetBondedValidatorsByPower(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBondedValidatorsByPower", reflect.TypeOf((*MockStakingKeeper)(nil).GetBondedValidatorsByPower), ctx)
}

// GetDelegation mocks base method.
func (m *MockStakingKeeper) GetDelegation(ctx types.Context, delAddr types.AccAddress, valAddr types.ValAddress) (types4.Delegation, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDelegation", ctx, delAddr, valAddr)
	ret0, _ := ret[0].(types4.Delegation)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetDelegation indicates an expected call of GetDelegation.
func (mr *MockStakingKeeperMockRecorder) GetDelegation(ctx, delAddr, valAddr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDelegation", reflect.TypeOf((*MockStakingKeeper)(nil).GetDelegation), ctx, delAddr, valAddr)
}

// GetUnbondingDelegation mocks base method.
func (m *MockStakingKeeper) GetUnbondingDelegation(ctx types.Context, delAddr types.AccAddress, valAddr types.ValAddress) (types4.UnbondingDelegation, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnbondingDelegation", ctx, delAddr, valAddr)
	ret0, _ := ret[0].(types4.UnbondingDelegation)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetUnbondingDelegation indicates an expected call of GetUnbondingDelegation.
func (mr *MockStakingKeeperMockRecorder) GetUnbondingDelegation(ctx, delAddr, valAddr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnbondingDelegation", reflect.TypeOf((*MockStakingKeeper)(nil).GetUnbondingDelegation), ctx, delAddr, valAddr)
}

// GetValidator mocks base method.
func (m *MockStakingKeeper) GetValidator(ctx types.Context, addr types.ValAddress) (types4.Validator, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidator", ctx, addr)
	ret0, _ := ret[0].(types4.Validator)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetValidator indicates an expected call of GetValidator.
func (mr *MockStakingKeeperMockRecorder) GetValidator(ctx, addr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidator", reflect.TypeOf((*MockStakingKeeper)(nil).GetValidator), ctx, addr)
}

// RemoveDelegation mocks base method.
func (m *MockStakingKeeper) RemoveDelegation(ctx types.Context, delegation types4.Delegation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveDelegation", ctx, delegation)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveDelegation indicates an expected call of RemoveDelegation.
func (mr *MockStakingKeeperMockRecorder) RemoveDelegation(ctx, delegation any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveDelegation", reflect.TypeOf((*MockStakingKeeper)(nil).RemoveDelegation), ctx, delegation)
}

// MockStakingMsgServer is a mock of StakingMsgServer interface.
type MockStakingMsgServer struct {
	ctrl     *gomock.Controller
	recorder *MockStakingMsgServerMockRecorder
}

// MockStakingMsgServerMockRecorder is the mock recorder for MockStakingMsgServer.
type MockStakingMsgServerMockRecorder struct {
	mock *MockStakingMsgServer
}

// NewMockStakingMsgServer creates a new mock instance.
func NewMockStakingMsgServer(ctrl *gomock.Controller) *MockStakingMsgServer {
	mock := &MockStakingMsgServer{ctrl: ctrl}
	mock.recorder = &MockStakingMsgServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStakingMsgServer) EXPECT() *MockStakingMsgServerMockRecorder {
	return m.recorder
}

// BeginRedelegate mocks base method.
func (m *MockStakingMsgServer) BeginRedelegate(goCtx context.Context, msg *types4.MsgBeginRedelegate) (*types4.MsgBeginRedelegateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginRedelegate", goCtx, msg)
	ret0, _ := ret[0].(*types4.MsgBeginRedelegateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginRedelegate indicates an expected call of BeginRedelegate.
func (mr *MockStakingMsgServerMockRecorder) BeginRedelegate(goCtx, msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginRedelegate", reflect.TypeOf((*MockStakingMsgServer)(nil).BeginRedelegate), goCtx, msg)
}

// Delegate mocks base method.
func (m *MockStakingMsgServer) Delegate(goCtx context.Context, msg *types4.MsgDelegate) (*types4.MsgDelegateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delegate", goCtx, msg)
	ret0, _ := ret[0].(*types4.MsgDelegateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delegate indicates an expected call of Delegate.
func (mr *MockStakingMsgServerMockRecorder) Delegate(goCtx, msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delegate", reflect.TypeOf((*MockStakingMsgServer)(nil).Delegate), goCtx, msg)
}

// Undelegate mocks base method.
func (m *MockStakingMsgServer) Undelegate(goCtx context.Context, msg *types4.MsgUndelegate) (*types4.MsgUndelegateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Undelegate", goCtx, msg)
	ret0, _ := ret[0].(*types4.MsgUndelegateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Undelegate indicates an expected call of Undelegate.
func (mr *MockStakingMsgServerMockRecorder) Undelegate(goCtx, msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Undelegate", reflect.TypeOf((*MockStakingMsgServer)(nil).Undelegate), goCtx, msg)
}

// MockDistributionKeeper is a mock of DistributionKeeper interface.
type MockDistributionKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockDistributionKeeperMockRecorder
}

// MockDistributionKeeperMockRecorder is the mock recorder for MockDistributionKeeper.
type MockDistributionKeeperMockRecorder struct {
	mock *MockDistributionKeeper
}

// NewMockDistributionKeeper creates a new mock instance.
func NewMockDistributionKeeper(ctrl *gomock.Controller) *MockDistributionKeeper {
	mock := &MockDistributionKeeper{ctrl: ctrl}
	mock.recorder = &MockDistributionKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDistributionKeeper) EXPECT() *MockDistributionKeeperMockRecorder {
	return m.recorder
}

// GetDelegatorStartingInfo mocks base method.
func (m *MockDistributionKeeper) GetDelegatorStartingInfo(ctx types.Context, val types.ValAddress, del types.AccAddress) types2.DelegatorStartingInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDelegatorStartingInfo", ctx, val, del)
	ret0, _ := ret[0].(types2.DelegatorStartingInfo)
	return ret0
}

// GetDelegatorStartingInfo indicates an expected call of GetDelegatorStartingInfo.
func (mr *MockDistributionKeeperMockRecorder) GetDelegatorStartingInfo(ctx, val, del any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDelegatorStartingInfo", reflect.TypeOf((*MockDistributionKeeper)(nil).GetDelegatorStartingInfo), ctx, val, del)
}

// MockDistributionMsgServer is a mock of DistributionMsgServer interface.
type MockDistributionMsgServer struct {
	ctrl     *gomock.Controller
	recorder *MockDistributionMsgServerMockRecorder
}

// MockDistributionMsgServerMockRecorder is the mock recorder for MockDistributionMsgServer.
type MockDistributionMsgServerMockRecorder struct {
	mock *MockDistributionMsgServer
}

// NewMockDistributionMsgServer creates a new mock instance.
func NewMockDistributionMsgServer(ctrl *gomock.Controller) *MockDistributionMsgServer {
	mock := &MockDistributionMsgServer{ctrl: ctrl}
	mock.recorder = &MockDistributionMsgServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDistributionMsgServer) EXPECT() *MockDistributionMsgServerMockRecorder {
	return m.recorder
}

// WithdrawDelegatorReward mocks base method.
func (m *MockDistributionMsgServer) WithdrawDelegatorReward(goCtx context.Context, msg *types2.MsgWithdrawDelegatorReward) (*types2.MsgWithdrawDelegatorRewardResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithdrawDelegatorReward", goCtx, msg)
	ret0, _ := ret[0].(*types2.MsgWithdrawDelegatorRewardResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WithdrawDelegatorReward indicates an expected call of WithdrawDelegatorReward.
func (mr *MockDistributionMsgServerMockRecorder) WithdrawDelegatorReward(goCtx, msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithdrawDelegatorReward", reflect.TypeOf((*MockDistributionMsgServer)(nil).WithdrawDelegatorReward), goCtx, msg)
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
func (m *MockBankKeeper) BurnCoins(ctx types.Context, name string, amt types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BurnCoins", ctx, name, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// BurnCoins indicates an expected call of BurnCoins.
func (mr *MockBankKeeperMockRecorder) BurnCoins(ctx, name, amt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BurnCoins", reflect.TypeOf((*MockBankKeeper)(nil).BurnCoins), ctx, name, amt)
}

// GetAllBalances mocks base method.
func (m *MockBankKeeper) GetAllBalances(ctx types.Context, addr types.AccAddress) types.Coins {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllBalances", ctx, addr)
	ret0, _ := ret[0].(types.Coins)
	return ret0
}

// GetAllBalances indicates an expected call of GetAllBalances.
func (mr *MockBankKeeperMockRecorder) GetAllBalances(ctx, addr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllBalances", reflect.TypeOf((*MockBankKeeper)(nil).GetAllBalances), ctx, addr)
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

// HasBalance mocks base method.
func (m *MockBankKeeper) HasBalance(ctx types.Context, addr types.AccAddress, amt types.Coin) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasBalance", ctx, addr, amt)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasBalance indicates an expected call of HasBalance.
func (mr *MockBankKeeperMockRecorder) HasBalance(ctx, addr, amt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasBalance", reflect.TypeOf((*MockBankKeeper)(nil).HasBalance), ctx, addr, amt)
}

// HasDenomMetaData mocks base method.
func (m *MockBankKeeper) HasDenomMetaData(ctx types.Context, denom string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasDenomMetaData", ctx, denom)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasDenomMetaData indicates an expected call of HasDenomMetaData.
func (mr *MockBankKeeperMockRecorder) HasDenomMetaData(ctx, denom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasDenomMetaData", reflect.TypeOf((*MockBankKeeper)(nil).HasDenomMetaData), ctx, denom)
}

// IterateAllDenomMetaData mocks base method.
func (m *MockBankKeeper) IterateAllDenomMetaData(ctx types.Context, cb func(types1.Metadata) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IterateAllDenomMetaData", ctx, cb)
}

// IterateAllDenomMetaData indicates an expected call of IterateAllDenomMetaData.
func (mr *MockBankKeeperMockRecorder) IterateAllDenomMetaData(ctx, cb any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IterateAllDenomMetaData", reflect.TypeOf((*MockBankKeeper)(nil).IterateAllDenomMetaData), ctx, cb)
}

// MintCoins mocks base method.
func (m *MockBankKeeper) MintCoins(ctx types.Context, name string, amt types.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MintCoins", ctx, name, amt)
	ret0, _ := ret[0].(error)
	return ret0
}

// MintCoins indicates an expected call of MintCoins.
func (mr *MockBankKeeperMockRecorder) MintCoins(ctx, name, amt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MintCoins", reflect.TypeOf((*MockBankKeeper)(nil).MintCoins), ctx, name, amt)
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

// ConvertCoin mocks base method.
func (m *MockErc20Keeper) ConvertCoin(goCtx context.Context, msg *types9.MsgConvertCoin) (*types9.MsgConvertCoinResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConvertCoin", goCtx, msg)
	ret0, _ := ret[0].(*types9.MsgConvertCoinResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConvertCoin indicates an expected call of ConvertCoin.
func (mr *MockErc20KeeperMockRecorder) ConvertCoin(goCtx, msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConvertCoin", reflect.TypeOf((*MockErc20Keeper)(nil).ConvertCoin), goCtx, msg)
}

// ConvertDenomToTarget mocks base method.
func (m *MockErc20Keeper) ConvertDenomToTarget(ctx types.Context, from types.AccAddress, coin types.Coin, fxTarget types7.FxTarget) (types.Coin, error) {
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

// DeleteOutgoingTransferRelation mocks base method.
func (m *MockErc20Keeper) DeleteOutgoingTransferRelation(ctx types.Context, moduleName string, txID uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteOutgoingTransferRelation", ctx, moduleName, txID)
}

// DeleteOutgoingTransferRelation indicates an expected call of DeleteOutgoingTransferRelation.
func (mr *MockErc20KeeperMockRecorder) DeleteOutgoingTransferRelation(ctx, moduleName, txID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOutgoingTransferRelation", reflect.TypeOf((*MockErc20Keeper)(nil).DeleteOutgoingTransferRelation), ctx, moduleName, txID)
}

// GetTokenPair mocks base method.
func (m *MockErc20Keeper) GetTokenPair(ctx types.Context, tokenOrDenom string) (types9.TokenPair, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenPair", ctx, tokenOrDenom)
	ret0, _ := ret[0].(types9.TokenPair)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetTokenPair indicates an expected call of GetTokenPair.
func (mr *MockErc20KeeperMockRecorder) GetTokenPair(ctx, tokenOrDenom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenPair", reflect.TypeOf((*MockErc20Keeper)(nil).GetTokenPair), ctx, tokenOrDenom)
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

// HookOutgoingRefund mocks base method.
func (m *MockErc20Keeper) HookOutgoingRefund(ctx types.Context, moduleName string, txID uint64, sender types.AccAddress, totalCoin types.Coin) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HookOutgoingRefund", ctx, moduleName, txID, sender, totalCoin)
	ret0, _ := ret[0].(error)
	return ret0
}

// HookOutgoingRefund indicates an expected call of HookOutgoingRefund.
func (mr *MockErc20KeeperMockRecorder) HookOutgoingRefund(ctx, moduleName, txID, sender, totalCoin any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HookOutgoingRefund", reflect.TypeOf((*MockErc20Keeper)(nil).HookOutgoingRefund), ctx, moduleName, txID, sender, totalCoin)
}

// IsOriginOrConvertedDenom mocks base method.
func (m *MockErc20Keeper) IsOriginOrConvertedDenom(ctx types.Context, denom string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsOriginOrConvertedDenom", ctx, denom)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsOriginOrConvertedDenom indicates an expected call of IsOriginOrConvertedDenom.
func (mr *MockErc20KeeperMockRecorder) IsOriginOrConvertedDenom(ctx, denom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsOriginOrConvertedDenom", reflect.TypeOf((*MockErc20Keeper)(nil).IsOriginOrConvertedDenom), ctx, denom)
}

// RefundLiquidity mocks base method.
func (m *MockErc20Keeper) RefundLiquidity(ctx types.Context, from types.AccAddress, coin types.Coin) (types.Coin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefundLiquidity", ctx, from, coin)
	ret0, _ := ret[0].(types.Coin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefundLiquidity indicates an expected call of RefundLiquidity.
func (mr *MockErc20KeeperMockRecorder) RefundLiquidity(ctx, from, coin any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefundLiquidity", reflect.TypeOf((*MockErc20Keeper)(nil).RefundLiquidity), ctx, from, coin)
}

// SetOutgoingTransferRelation mocks base method.
func (m *MockErc20Keeper) SetOutgoingTransferRelation(ctx types.Context, moduleName string, txID uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOutgoingTransferRelation", ctx, moduleName, txID)
}

// SetOutgoingTransferRelation indicates an expected call of SetOutgoingTransferRelation.
func (mr *MockErc20KeeperMockRecorder) SetOutgoingTransferRelation(ctx, moduleName, txID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOutgoingTransferRelation", reflect.TypeOf((*MockErc20Keeper)(nil).SetOutgoingTransferRelation), ctx, moduleName, txID)
}

// ToTargetDenom mocks base method.
func (m *MockErc20Keeper) ToTargetDenom(ctx types.Context, denom, base string, aliases []string, fxTarget types7.FxTarget) string {
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

// TransferAfter mocks base method.
func (m *MockErc20Keeper) TransferAfter(ctx types.Context, sender types.AccAddress, receive string, coin, fee types.Coin, arg5, arg6 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferAfter", ctx, sender, receive, coin, fee, arg5, arg6)
	ret0, _ := ret[0].(error)
	return ret0
}

// TransferAfter indicates an expected call of TransferAfter.
func (mr *MockErc20KeeperMockRecorder) TransferAfter(ctx, sender, receive, coin, fee, arg5, arg6 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferAfter", reflect.TypeOf((*MockErc20Keeper)(nil).TransferAfter), ctx, sender, receive, coin, fee, arg5, arg6)
}

// MockEVMKeeper is a mock of EVMKeeper interface.
type MockEVMKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockEVMKeeperMockRecorder
}

// MockEVMKeeperMockRecorder is the mock recorder for MockEVMKeeper.
type MockEVMKeeperMockRecorder struct {
	mock *MockEVMKeeper
}

// NewMockEVMKeeper creates a new mock instance.
func NewMockEVMKeeper(ctrl *gomock.Controller) *MockEVMKeeper {
	mock := &MockEVMKeeper{ctrl: ctrl}
	mock.recorder = &MockEVMKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEVMKeeper) EXPECT() *MockEVMKeeperMockRecorder {
	return m.recorder
}

// CallEVM mocks base method.
func (m *MockEVMKeeper) CallEVM(ctx types.Context, from common.Address, contract *common.Address, value *big.Int, gasLimit uint64, data []byte, commit bool) (*types6.MsgEthereumTxResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CallEVM", ctx, from, contract, value, gasLimit, data, commit)
	ret0, _ := ret[0].(*types6.MsgEthereumTxResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CallEVM indicates an expected call of CallEVM.
func (mr *MockEVMKeeperMockRecorder) CallEVM(ctx, from, contract, value, gasLimit, data, commit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CallEVM", reflect.TypeOf((*MockEVMKeeper)(nil).CallEVM), ctx, from, contract, value, gasLimit, data, commit)
}

// IsContract mocks base method.
func (m *MockEVMKeeper) IsContract(ctx types.Context, account common.Address) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsContract", ctx, account)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsContract indicates an expected call of IsContract.
func (mr *MockEVMKeeperMockRecorder) IsContract(ctx, account any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsContract", reflect.TypeOf((*MockEVMKeeper)(nil).IsContract), ctx, account)
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

// SetDenomTrace mocks base method.
func (m *MockIBCTransferKeeper) SetDenomTrace(ctx types.Context, denomTrace types5.DenomTrace) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDenomTrace", ctx, denomTrace)
}

// SetDenomTrace indicates an expected call of SetDenomTrace.
func (mr *MockIBCTransferKeeperMockRecorder) SetDenomTrace(ctx, denomTrace any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDenomTrace", reflect.TypeOf((*MockIBCTransferKeeper)(nil).SetDenomTrace), ctx, denomTrace)
}

// Transfer mocks base method.
func (m *MockIBCTransferKeeper) Transfer(goCtx context.Context, msg *types5.MsgTransfer) (*types5.MsgTransferResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transfer", goCtx, msg)
	ret0, _ := ret[0].(*types5.MsgTransferResponse)
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

// GetAccount mocks base method.
func (m *MockAccountKeeper) GetAccount(ctx types.Context, addr types.AccAddress) types0.AccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", ctx, addr)
	ret0, _ := ret[0].(types0.AccountI)
	return ret0
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockAccountKeeperMockRecorder) GetAccount(ctx, addr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccountKeeper)(nil).GetAccount), ctx, addr)
}

// GetModuleAddress mocks base method.
func (m *MockAccountKeeper) GetModuleAddress(name string) types.AccAddress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModuleAddress", name)
	ret0, _ := ret[0].(types.AccAddress)
	return ret0
}

// GetModuleAddress indicates an expected call of GetModuleAddress.
func (mr *MockAccountKeeperMockRecorder) GetModuleAddress(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModuleAddress", reflect.TypeOf((*MockAccountKeeper)(nil).GetModuleAddress), name)
}

// NewAccountWithAddress mocks base method.
func (m *MockAccountKeeper) NewAccountWithAddress(ctx types.Context, addr types.AccAddress) types0.AccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewAccountWithAddress", ctx, addr)
	ret0, _ := ret[0].(types0.AccountI)
	return ret0
}

// NewAccountWithAddress indicates an expected call of NewAccountWithAddress.
func (mr *MockAccountKeeperMockRecorder) NewAccountWithAddress(ctx, addr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewAccountWithAddress", reflect.TypeOf((*MockAccountKeeper)(nil).NewAccountWithAddress), ctx, addr)
}

// SetAccount mocks base method.
func (m *MockAccountKeeper) SetAccount(ctx types.Context, acc types0.AccountI) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetAccount", ctx, acc)
}

// SetAccount indicates an expected call of SetAccount.
func (mr *MockAccountKeeperMockRecorder) SetAccount(ctx, acc any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAccount", reflect.TypeOf((*MockAccountKeeper)(nil).SetAccount), ctx, acc)
}

// MockSubspace is a mock of Subspace interface.
type MockSubspace struct {
	ctrl     *gomock.Controller
	recorder *MockSubspaceMockRecorder
}

// MockSubspaceMockRecorder is the mock recorder for MockSubspace.
type MockSubspaceMockRecorder struct {
	mock *MockSubspace
}

// NewMockSubspace creates a new mock instance.
func NewMockSubspace(ctrl *gomock.Controller) *MockSubspace {
	mock := &MockSubspace{ctrl: ctrl}
	mock.recorder = &MockSubspaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubspace) EXPECT() *MockSubspaceMockRecorder {
	return m.recorder
}

// GetParamSet mocks base method.
func (m *MockSubspace) GetParamSet(ctx types.Context, ps types8.ParamSet) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetParamSet", ctx, ps)
}

// GetParamSet indicates an expected call of GetParamSet.
func (mr *MockSubspaceMockRecorder) GetParamSet(ctx, ps any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParamSet", reflect.TypeOf((*MockSubspace)(nil).GetParamSet), ctx, ps)
}

// HasKeyTable mocks base method.
func (m *MockSubspace) HasKeyTable() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasKeyTable")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasKeyTable indicates an expected call of HasKeyTable.
func (mr *MockSubspaceMockRecorder) HasKeyTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasKeyTable", reflect.TypeOf((*MockSubspace)(nil).HasKeyTable))
}

// WithKeyTable mocks base method.
func (m *MockSubspace) WithKeyTable(table types3.KeyTable) types3.Subspace {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithKeyTable", table)
	ret0, _ := ret[0].(types3.Subspace)
	return ret0
}

// WithKeyTable indicates an expected call of WithKeyTable.
func (mr *MockSubspaceMockRecorder) WithKeyTable(table any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithKeyTable", reflect.TypeOf((*MockSubspace)(nil).WithKeyTable), table)
}

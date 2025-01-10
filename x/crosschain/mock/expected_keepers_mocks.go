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

	math "cosmossdk.io/math"
	bytes "github.com/cometbft/cometbft/libs/bytes"
	types "github.com/cosmos/cosmos-sdk/types"
	types0 "github.com/cosmos/cosmos-sdk/x/distribution/types"
	types1 "github.com/cosmos/cosmos-sdk/x/staking/types"
	types2 "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	abi "github.com/ethereum/go-ethereum/accounts/abi"
	common "github.com/ethereum/go-ethereum/common"
	types3 "github.com/evmos/ethermint/x/evm/types"
	contract "github.com/pundiai/fx-core/v8/contract"
	types4 "github.com/pundiai/fx-core/v8/x/erc20/types"
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

// GetDelegation mocks base method.
func (m *MockStakingKeeper) GetDelegation(ctx context.Context, delAddr types.AccAddress, valAddr types.ValAddress) (types1.Delegation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDelegation", ctx, delAddr, valAddr)
	ret0, _ := ret[0].(types1.Delegation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDelegation indicates an expected call of GetDelegation.
func (mr *MockStakingKeeperMockRecorder) GetDelegation(ctx, delAddr, valAddr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDelegation", reflect.TypeOf((*MockStakingKeeper)(nil).GetDelegation), ctx, delAddr, valAddr)
}

// GetUnbondingDelegation mocks base method.
func (m *MockStakingKeeper) GetUnbondingDelegation(ctx context.Context, delAddr types.AccAddress, valAddr types.ValAddress) (types1.UnbondingDelegation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnbondingDelegation", ctx, delAddr, valAddr)
	ret0, _ := ret[0].(types1.UnbondingDelegation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnbondingDelegation indicates an expected call of GetUnbondingDelegation.
func (mr *MockStakingKeeperMockRecorder) GetUnbondingDelegation(ctx, delAddr, valAddr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnbondingDelegation", reflect.TypeOf((*MockStakingKeeper)(nil).GetUnbondingDelegation), ctx, delAddr, valAddr)
}

// GetValidator mocks base method.
func (m *MockStakingKeeper) GetValidator(ctx context.Context, addr types.ValAddress) (types1.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidator", ctx, addr)
	ret0, _ := ret[0].(types1.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidator indicates an expected call of GetValidator.
func (mr *MockStakingKeeperMockRecorder) GetValidator(ctx, addr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidator", reflect.TypeOf((*MockStakingKeeper)(nil).GetValidator), ctx, addr)
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
func (m *MockStakingMsgServer) BeginRedelegate(goCtx context.Context, msg *types1.MsgBeginRedelegate) (*types1.MsgBeginRedelegateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginRedelegate", goCtx, msg)
	ret0, _ := ret[0].(*types1.MsgBeginRedelegateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginRedelegate indicates an expected call of BeginRedelegate.
func (mr *MockStakingMsgServerMockRecorder) BeginRedelegate(goCtx, msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginRedelegate", reflect.TypeOf((*MockStakingMsgServer)(nil).BeginRedelegate), goCtx, msg)
}

// Delegate mocks base method.
func (m *MockStakingMsgServer) Delegate(goCtx context.Context, msg *types1.MsgDelegate) (*types1.MsgDelegateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delegate", goCtx, msg)
	ret0, _ := ret[0].(*types1.MsgDelegateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delegate indicates an expected call of Delegate.
func (mr *MockStakingMsgServerMockRecorder) Delegate(goCtx, msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delegate", reflect.TypeOf((*MockStakingMsgServer)(nil).Delegate), goCtx, msg)
}

// Undelegate mocks base method.
func (m *MockStakingMsgServer) Undelegate(goCtx context.Context, msg *types1.MsgUndelegate) (*types1.MsgUndelegateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Undelegate", goCtx, msg)
	ret0, _ := ret[0].(*types1.MsgUndelegateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Undelegate indicates an expected call of Undelegate.
func (mr *MockStakingMsgServerMockRecorder) Undelegate(goCtx, msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Undelegate", reflect.TypeOf((*MockStakingMsgServer)(nil).Undelegate), goCtx, msg)
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
func (m *MockDistributionMsgServer) WithdrawDelegatorReward(goCtx context.Context, msg *types0.MsgWithdrawDelegatorReward) (*types0.MsgWithdrawDelegatorRewardResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithdrawDelegatorReward", goCtx, msg)
	ret0, _ := ret[0].(*types0.MsgWithdrawDelegatorRewardResponse)
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
func (m *MockBankKeeper) BurnCoins(ctx context.Context, name string, amt types.Coins) error {
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
func (m *MockBankKeeper) GetAllBalances(ctx context.Context, addr types.AccAddress) types.Coins {
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
func (m *MockBankKeeper) GetSupply(ctx context.Context, denom string) types.Coin {
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
func (m *MockBankKeeper) MintCoins(ctx context.Context, name string, amt types.Coins) error {
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
func (m *MockBankKeeper) SendCoins(ctx context.Context, fromAddr, toAddr types.AccAddress, amt types.Coins) error {
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
func (m *MockBankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr types.AccAddress, recipientModule string, amt types.Coins) error {
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
func (m *MockBankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr types.AccAddress, amt types.Coins) error {
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

// AddBridgeToken mocks base method.
func (m *MockErc20Keeper) AddBridgeToken(ctx context.Context, baseDenom, chainName, contract string, isNative bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBridgeToken", ctx, baseDenom, chainName, contract, isNative)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBridgeToken indicates an expected call of AddBridgeToken.
func (mr *MockErc20KeeperMockRecorder) AddBridgeToken(ctx, baseDenom, chainName, contract, isNative any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBridgeToken", reflect.TypeOf((*MockErc20Keeper)(nil).AddBridgeToken), ctx, baseDenom, chainName, contract, isNative)
}

// BaseCoinToEvm mocks base method.
func (m *MockErc20Keeper) BaseCoinToEvm(ctx context.Context, caller contract.Caller, holder common.Address, coin types.Coin) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BaseCoinToEvm", ctx, caller, holder, coin)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BaseCoinToEvm indicates an expected call of BaseCoinToEvm.
func (mr *MockErc20KeeperMockRecorder) BaseCoinToEvm(ctx, caller, holder, coin any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BaseCoinToEvm", reflect.TypeOf((*MockErc20Keeper)(nil).BaseCoinToEvm), ctx, caller, holder, coin)
}

// DeleteCache mocks base method.
func (m *MockErc20Keeper) DeleteCache(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCache", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCache indicates an expected call of DeleteCache.
func (mr *MockErc20KeeperMockRecorder) DeleteCache(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCache", reflect.TypeOf((*MockErc20Keeper)(nil).DeleteCache), ctx, key)
}

// GetBaseDenom mocks base method.
func (m *MockErc20Keeper) GetBaseDenom(ctx context.Context, token string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBaseDenom", ctx, token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBaseDenom indicates an expected call of GetBaseDenom.
func (mr *MockErc20KeeperMockRecorder) GetBaseDenom(ctx, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBaseDenom", reflect.TypeOf((*MockErc20Keeper)(nil).GetBaseDenom), ctx, token)
}

// GetBridgeToken mocks base method.
func (m *MockErc20Keeper) GetBridgeToken(ctx context.Context, chainName, baseDenom string) (types4.BridgeToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBridgeToken", ctx, chainName, baseDenom)
	ret0, _ := ret[0].(types4.BridgeToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBridgeToken indicates an expected call of GetBridgeToken.
func (mr *MockErc20KeeperMockRecorder) GetBridgeToken(ctx, chainName, baseDenom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBridgeToken", reflect.TypeOf((*MockErc20Keeper)(nil).GetBridgeToken), ctx, chainName, baseDenom)
}

// GetBridgeTokens mocks base method.
func (m *MockErc20Keeper) GetBridgeTokens(ctx context.Context, chainName string) ([]types4.BridgeToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBridgeTokens", ctx, chainName)
	ret0, _ := ret[0].([]types4.BridgeToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBridgeTokens indicates an expected call of GetBridgeTokens.
func (mr *MockErc20KeeperMockRecorder) GetBridgeTokens(ctx, chainName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBridgeTokens", reflect.TypeOf((*MockErc20Keeper)(nil).GetBridgeTokens), ctx, chainName)
}

// GetCache mocks base method.
func (m *MockErc20Keeper) GetCache(ctx context.Context, key string) (math.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCache", ctx, key)
	ret0, _ := ret[0].(math.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCache indicates an expected call of GetCache.
func (mr *MockErc20KeeperMockRecorder) GetCache(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCache", reflect.TypeOf((*MockErc20Keeper)(nil).GetCache), ctx, key)
}

// GetERC20Token mocks base method.
func (m *MockErc20Keeper) GetERC20Token(ctx context.Context, baseDenom string) (types4.ERC20Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetERC20Token", ctx, baseDenom)
	ret0, _ := ret[0].(types4.ERC20Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetERC20Token indicates an expected call of GetERC20Token.
func (mr *MockErc20KeeperMockRecorder) GetERC20Token(ctx, baseDenom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetERC20Token", reflect.TypeOf((*MockErc20Keeper)(nil).GetERC20Token), ctx, baseDenom)
}

// GetIBCToken mocks base method.
func (m *MockErc20Keeper) GetIBCToken(ctx context.Context, channel, baseDenom string) (types4.IBCToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIBCToken", ctx, channel, baseDenom)
	ret0, _ := ret[0].(types4.IBCToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIBCToken indicates an expected call of GetIBCToken.
func (mr *MockErc20KeeperMockRecorder) GetIBCToken(ctx, channel, baseDenom any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIBCToken", reflect.TypeOf((*MockErc20Keeper)(nil).GetIBCToken), ctx, channel, baseDenom)
}

// HasCache mocks base method.
func (m *MockErc20Keeper) HasCache(ctx context.Context, key string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasCache", ctx, key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasCache indicates an expected call of HasCache.
func (mr *MockErc20KeeperMockRecorder) HasCache(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasCache", reflect.TypeOf((*MockErc20Keeper)(nil).HasCache), ctx, key)
}

// ReSetCache mocks base method.
func (m *MockErc20Keeper) ReSetCache(ctx context.Context, oldKey, newKey string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReSetCache", ctx, oldKey, newKey)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReSetCache indicates an expected call of ReSetCache.
func (mr *MockErc20KeeperMockRecorder) ReSetCache(ctx, oldKey, newKey any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReSetCache", reflect.TypeOf((*MockErc20Keeper)(nil).ReSetCache), ctx, oldKey, newKey)
}

// SetCache mocks base method.
func (m *MockErc20Keeper) SetCache(ctx context.Context, key string, amount math.Int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCache", ctx, key, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCache indicates an expected call of SetCache.
func (mr *MockErc20KeeperMockRecorder) SetCache(ctx, key, amount any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCache", reflect.TypeOf((*MockErc20Keeper)(nil).SetCache), ctx, key, amount)
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

// ApplyContract mocks base method.
func (m *MockEVMKeeper) ApplyContract(ctx context.Context, from, contract common.Address, value *big.Int, abi abi.ABI, method string, constructorData ...any) (*types3.MsgEthereumTxResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, from, contract, value, abi, method}
	for _, a := range constructorData {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ApplyContract", varargs...)
	ret0, _ := ret[0].(*types3.MsgEthereumTxResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplyContract indicates an expected call of ApplyContract.
func (mr *MockEVMKeeperMockRecorder) ApplyContract(ctx, from, contract, value, abi, method any, constructorData ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, from, contract, value, abi, method}, constructorData...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyContract", reflect.TypeOf((*MockEVMKeeper)(nil).ApplyContract), varargs...)
}

// ExecuteEVM mocks base method.
func (m *MockEVMKeeper) ExecuteEVM(ctx types.Context, from common.Address, contract *common.Address, value *big.Int, gasLimit uint64, data []byte) (*types3.MsgEthereumTxResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecuteEVM", ctx, from, contract, value, gasLimit, data)
	ret0, _ := ret[0].(*types3.MsgEthereumTxResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecuteEVM indicates an expected call of ExecuteEVM.
func (mr *MockEVMKeeperMockRecorder) ExecuteEVM(ctx, from, contract, value, gasLimit, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteEVM", reflect.TypeOf((*MockEVMKeeper)(nil).ExecuteEVM), ctx, from, contract, value, gasLimit, data)
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

// QueryContract mocks base method.
func (m *MockEVMKeeper) QueryContract(ctx context.Context, from, contract common.Address, abi abi.ABI, method string, res any, args ...any) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, from, contract, abi, method, res}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryContract", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// QueryContract indicates an expected call of QueryContract.
func (mr *MockEVMKeeperMockRecorder) QueryContract(ctx, from, contract, abi, method, res any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, from, contract, abi, method, res}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContract", reflect.TypeOf((*MockEVMKeeper)(nil).QueryContract), varargs...)
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

// GetDenomTrace mocks base method.
func (m *MockIBCTransferKeeper) GetDenomTrace(ctx types.Context, denomTraceHash bytes.HexBytes) (types2.DenomTrace, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDenomTrace", ctx, denomTraceHash)
	ret0, _ := ret[0].(types2.DenomTrace)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetDenomTrace indicates an expected call of GetDenomTrace.
func (mr *MockIBCTransferKeeperMockRecorder) GetDenomTrace(ctx, denomTraceHash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDenomTrace", reflect.TypeOf((*MockIBCTransferKeeper)(nil).GetDenomTrace), ctx, denomTraceHash)
}

// SetDenomTrace mocks base method.
func (m *MockIBCTransferKeeper) SetDenomTrace(ctx types.Context, denomTrace types2.DenomTrace) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDenomTrace", ctx, denomTrace)
}

// SetDenomTrace indicates an expected call of SetDenomTrace.
func (mr *MockIBCTransferKeeperMockRecorder) SetDenomTrace(ctx, denomTrace any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDenomTrace", reflect.TypeOf((*MockIBCTransferKeeper)(nil).SetDenomTrace), ctx, denomTrace)
}

// Transfer mocks base method.
func (m *MockIBCTransferKeeper) Transfer(ctx context.Context, msg *types2.MsgTransfer) (*types2.MsgTransferResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transfer", ctx, msg)
	ret0, _ := ret[0].(*types2.MsgTransferResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Transfer indicates an expected call of Transfer.
func (mr *MockIBCTransferKeeperMockRecorder) Transfer(ctx, msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transfer", reflect.TypeOf((*MockIBCTransferKeeper)(nil).Transfer), ctx, msg)
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
func (m *MockAccountKeeper) GetAccount(ctx context.Context, addr types.AccAddress) types.AccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", ctx, addr)
	ret0, _ := ret[0].(types.AccountI)
	return ret0
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockAccountKeeperMockRecorder) GetAccount(ctx, addr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccountKeeper)(nil).GetAccount), ctx, addr)
}

// GetModuleAccount mocks base method.
func (m *MockAccountKeeper) GetModuleAccount(ctx context.Context, moduleName string) types.ModuleAccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModuleAccount", ctx, moduleName)
	ret0, _ := ret[0].(types.ModuleAccountI)
	return ret0
}

// GetModuleAccount indicates an expected call of GetModuleAccount.
func (mr *MockAccountKeeperMockRecorder) GetModuleAccount(ctx, moduleName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModuleAccount", reflect.TypeOf((*MockAccountKeeper)(nil).GetModuleAccount), ctx, moduleName)
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
func (m *MockAccountKeeper) NewAccountWithAddress(ctx context.Context, addr types.AccAddress) types.AccountI {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewAccountWithAddress", ctx, addr)
	ret0, _ := ret[0].(types.AccountI)
	return ret0
}

// NewAccountWithAddress indicates an expected call of NewAccountWithAddress.
func (mr *MockAccountKeeperMockRecorder) NewAccountWithAddress(ctx, addr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewAccountWithAddress", reflect.TypeOf((*MockAccountKeeper)(nil).NewAccountWithAddress), ctx, addr)
}

// SetAccount mocks base method.
func (m *MockAccountKeeper) SetAccount(ctx context.Context, acc types.AccountI) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetAccount", ctx, acc)
}

// SetAccount indicates an expected call of SetAccount.
func (mr *MockAccountKeeperMockRecorder) SetAccount(ctx, acc any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAccount", reflect.TypeOf((*MockAccountKeeper)(nil).SetAccount), ctx, acc)
}

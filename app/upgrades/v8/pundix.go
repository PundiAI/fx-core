package v8

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/pundiai/fx-core/v8/app/keepers"
	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	bsctypes "github.com/pundiai/fx-core/v8/x/bsc/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20keeper "github.com/pundiai/fx-core/v8/x/erc20/keeper"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

type Pundix struct {
	cdc               codec.Codec
	accountKeeper     authkeeper.AccountKeeper
	erc20Keeper       erc20keeper.Keeper
	bankKeeper        bankkeeper.Keeper
	ibcTransferKeeper ibctransferkeeper.Keeper
	erc20TokenKeeper  contract.ERC20TokenKeeper
	storeKey          storetypes.StoreKey
}

func NewPundix(cdc codec.Codec, app *keepers.AppKeepers) *Pundix {
	return &Pundix{
		cdc:               cdc,
		accountKeeper:     app.AccountKeeper,
		erc20Keeper:       app.Erc20Keeper,
		bankKeeper:        app.BankKeeper,
		ibcTransferKeeper: app.IBCTransferKeeper,
		erc20TokenKeeper:  contract.NewERC20TokenKeeper(app.EvmKeeper),
		storeKey:          app.GetKey(ibctransfertypes.StoreKey),
	}
}

func (m *Pundix) Migrate(ctx sdk.Context) error {
	pundixGenesisPath := path.Join(fxtypes.GetDefaultNodeHome(), "config/pundix_genesis.json")
	appState, err := ReadGenesisState(pundixGenesisPath)
	if err != nil {
		return err
	}

	authRaw, ok := appState[types.ModuleName]
	if !ok || len(authRaw) == 0 {
		return sdkerrors.ErrNotFound.Wrap("auth genesis")
	}
	if err = m.migrateAuth(ctx, appState[types.ModuleName]); err != nil {
		return err
	}

	// todo migrate other data
	bankGenesis, ok := appState[banktypes.ModuleName]
	if !ok || len(bankGenesis) == 0 {
		return sdkerrors.ErrNotFound.Wrap("bank genesis")
	}
	return m.migrateBank(ctx, bankGenesis)
}

func (m *Pundix) migrateAuth(ctx sdk.Context, authRaw json.RawMessage) error {
	var authGenesis types.GenesisState
	if err := m.cdc.UnmarshalJSON(authRaw, &authGenesis); err != nil {
		return err
	}
	genesisAccounts, err := types.UnpackAccounts(authGenesis.Accounts)
	if err != nil {
		return err
	}
	for _, genAcc := range genesisAccounts {
		pubKey := genAcc.GetPubKey()
		if pubKey == nil {
			continue
		}
		baseAcc, ok := genAcc.(*types.BaseAccount)
		if !ok {
			continue
		}
		accAddr, err := sdk.GetFromBech32(baseAcc.Address, "px")
		if err != nil {
			return err
		}
		if m.accountKeeper.HasAccount(ctx, accAddr) {
			continue
		}
		newAcc := m.accountKeeper.NewAccountWithAddress(ctx, accAddr)
		if err = newAcc.SetPubKey(pubKey); err != nil {
			return err
		}
		m.accountKeeper.SetAccount(ctx, newAcc)
	}
	return nil
}

func (m *Pundix) migrateBank(ctx sdk.Context, bankRaw json.RawMessage) error {
	var bankGenesis banktypes.GenesisState
	if err := m.cdc.UnmarshalJSON(bankRaw, &bankGenesis); err != nil {
		return err
	}
	if err := m.migratePUNDIX(ctx, bankGenesis.Balances, bankGenesis.Supply); err != nil {
		return err
	}
	return m.migratePURSE(ctx, bankGenesis.Balances, bankGenesis.Supply)
}

func (m *Pundix) migratePUNDIX(ctx sdk.Context, balances []banktypes.Balance, supply sdk.Coins) error {
	baseDenom := "pundix"
	bridgeToken, err := m.erc20Keeper.GetBridgeToken(ctx, ethtypes.ModuleName, baseDenom)
	if err != nil {
		return err
	}
	pundixGenesisAmount := sdkmath.NewInt(3400).Mul(sdkmath.NewInt(1e18))
	pundixDenomHash := sha256.Sum256([]byte(fmt.Sprintf("%s/channel-0/%s", ibctransfertypes.ModuleName, bridgeToken.BridgeDenom())))
	pundixIBCDenom := fmt.Sprintf("%s/%X", ibctransfertypes.DenomPrefix, pundixDenomHash[:])
	escrowAddr := ibctransfertypes.GetEscrowAddress(ibctransfertypes.ModuleName, "channel-0")

	pundixSupply := supply.AmountOf(pundixIBCDenom)
	ibcTotalEscrow := m.ibcTransferKeeper.GetTotalEscrowForDenom(ctx, bridgeToken.BridgeDenom())
	ibcChannelBalance := m.bankKeeper.GetBalance(ctx, escrowAddr, bridgeToken.BridgeDenom())

	if ibcChannelBalance.Amount.Add(pundixGenesisAmount).LT(pundixSupply) ||
		ibcTotalEscrow.Amount.LT(ibcChannelBalance.Amount) {
		return sdkerrors.ErrInvalidCoins.Wrap("pundix ibc amount not match")
	}

	// remove ibc channel pundix amount
	if err = m.bankKeeper.SendCoinsFromAccountToModule(ctx, escrowAddr, crosschaintypes.ModuleName, sdk.NewCoins(ibcChannelBalance)); err != nil {
		return err
	}
	totalEscrow := m.ibcTransferKeeper.GetTotalEscrowForDenom(ctx, bridgeToken.BridgeDenom())
	totalEscrow = totalEscrow.Sub(ibcChannelBalance)
	m.ibcTransferKeeper.SetTotalEscrowForDenom(ctx, totalEscrow)

	// pundix chain will burn pundix token if slash validator
	uncountedAmount := pundixSupply.Sub(ibcChannelBalance.Amount)
	uncountedCoin := sdk.NewCoins(sdk.NewCoin(ibcChannelBalance.Denom, uncountedAmount.Abs()))
	if uncountedAmount.IsPositive() {
		if err = m.bankKeeper.MintCoins(ctx, crosschaintypes.ModuleName, uncountedCoin); err != nil {
			return err
		}
	} else if uncountedAmount.IsNegative() {
		if err = m.bankKeeper.BurnCoins(ctx, crosschaintypes.ModuleName, uncountedCoin); err != nil {
			return err
		}
	}

	for _, bal := range balances {
		bech32Addr, err := sdk.GetFromBech32(bal.Address, "px")
		if err != nil {
			return err
		}
		account := m.accountKeeper.GetAccount(ctx, bech32Addr)
		if ma, ok := account.(sdk.ModuleAccountI); ok {
			// todo migrate staking and reward
			ctx.Logger().Info("migrate PUNDIX skip module account", "module", ma.GetName(), "address", bal.Address, "balance", bal.Coins.String())
			continue
		}
		baesPundixCoin := sdk.NewCoins(sdk.NewCoin(baseDenom, bal.Coins.AmountOf(pundixIBCDenom)))
		if err = m.bankKeeper.MintCoins(ctx, crosschaintypes.ModuleName, baesPundixCoin); err != nil {
			return err
		}
		if err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, crosschaintypes.ModuleName, bech32Addr, baesPundixCoin); err != nil {
			return err
		}
	}
	return nil
}

func (m *Pundix) migratePURSE(ctx sdk.Context, balances []banktypes.Balance, supply sdk.Coins) error {
	baseDenom := "purse"

	erc20Token, err := m.erc20Keeper.GetERC20Token(ctx, baseDenom)
	if err != nil {
		return err
	}
	ibcToken, err := m.erc20Keeper.GetIBCToken(ctx, "channel-0", baseDenom)
	if err != nil {
		return err
	}
	bscBridgeToken, err := m.erc20Keeper.GetBridgeToken(ctx, bsctypes.ModuleName, baseDenom)
	if err != nil {
		return err
	}

	ibcPurseSupply := m.bankKeeper.GetSupply(ctx, ibcToken.IbcDenom)
	// bsc bridge denom is origin denom of pundix chain
	pundixPurseSupply := supply.AmountOf(bscBridgeToken.BridgeDenom())
	pxEscrowAddr, pxEscrowPurseAmount, err := pxEscrowAddrAndAmount(balances, bscBridgeToken.BridgeDenom())
	if err != nil {
		return err
	}

	if err = m.lockFxcorePurse(ctx, erc20Token, bscBridgeToken, ibcToken, ibcPurseSupply.Amount, pxEscrowPurseAmount); err != nil {
		return err
	}
	if err = m.migratePundixPurse(ctx, erc20Token, bscBridgeToken, balances, pundixPurseSupply.Sub(pxEscrowPurseAmount), pxEscrowAddr); err != nil {
		return err
	}
	if err = m.updatePurseErc20Owner(ctx, erc20Token); err != nil {
		return err
	}
	return m.removeDeprecatedPurse(ctx, baseDenom, "channel-0", ibcToken)
}

func pxEscrowAddrAndAmount(balances []banktypes.Balance, denom string) (string, sdkmath.Int, error) {
	pxEscrowAddr, err := GetPxChannelEscrowAddr()
	if err != nil {
		return "", sdkmath.Int{}, err
	}
	// bsc bridge denom is origin denom of pundix chain
	for _, bal := range balances {
		if bal.Address == pxEscrowAddr {
			return pxEscrowAddr, bal.Coins.AmountOf(denom), nil
		}
	}
	return pxEscrowAddr, sdkmath.ZeroInt(), nil
}

func GetPxChannelEscrowAddr() (string, error) {
	return bech32.ConvertAndEncode("px", ibctransfertypes.GetEscrowAddress(ibctransfertypes.ModuleName, "channel-0"))
}

func (m *Pundix) lockFxcorePurse(
	ctx sdk.Context,
	erc20Token erc20types.ERC20Token,
	bscBridgeToken erc20types.BridgeToken,
	ibcToken erc20types.IBCToken,
	ibcDenomPurseSupply, pxEscrowPurseAmount sdkmath.Int,
) error {
	erc20ModuleAddress := common.BytesToAddress(types.NewModuleAddress(erc20types.ModuleName).Bytes())
	erc20PurseTotalSupply, err := m.erc20TokenKeeper.TotalSupply(ctx, erc20Token.GetERC20Contract())
	if err != nil {
		return err
	}
	// cosmosPurseBal contain burned bsc bridge token
	cosmosPurseBal := pxEscrowPurseAmount.Sub(sdkmath.NewIntFromBigInt(erc20PurseTotalSupply))
	if _, err = m.erc20TokenKeeper.Mint(ctx, erc20Token.GetERC20Contract(), erc20ModuleAddress, erc20ModuleAddress, cosmosPurseBal.BigInt()); err != nil {
		return err
	}
	// make up the balance of base purse
	if err = m.bankKeeper.MintCoins(ctx, erc20types.ModuleName, sdk.NewCoins(sdk.NewCoin(erc20Token.GetDenom(), cosmosPurseBal))); err != nil {
		return err
	}

	// convert ibc purse to base purse, exclude erc20 module
	m.bankKeeper.IterateAllBalances(ctx, func(addr sdk.AccAddress, coin sdk.Coin) (stop bool) {
		if coin.GetDenom() != ibcToken.GetIbcDenom() {
			return false
		}
		account := m.accountKeeper.GetAccount(ctx, addr)
		if ma, ok := account.(sdk.ModuleAccountI); ok && ma.GetName() == erc20types.ModuleName {
			return false
		}

		if err = m.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, erc20types.ModuleName, sdk.NewCoins(coin)); err != nil {
			return true
		}
		basePurse := sdk.NewCoins(sdk.NewCoin(erc20Token.GetDenom(), coin.Amount))
		if err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, erc20types.ModuleName, addr, basePurse); err != nil {
			return true
		}
		return false
	})
	if err != nil {
		return err
	}

	// mint bsc bridge token to module
	bscBridgeAmount := pxEscrowPurseAmount.Sub(ibcDenomPurseSupply)
	if err = m.bankKeeper.MintCoins(ctx, bsctypes.ModuleName, sdk.NewCoins(sdk.NewCoin(bscBridgeToken.BridgeDenom(), bscBridgeAmount))); err != nil {
		return err
	}

	// burn all ibc purse
	return m.bankKeeper.BurnCoins(ctx, erc20types.ModuleName, sdk.NewCoins(sdk.NewCoin(ibcToken.GetIbcDenom(), ibcDenomPurseSupply)))
}

func (m *Pundix) migratePundixPurse(
	ctx sdk.Context,
	erc20Token erc20types.ERC20Token,
	bridgeToken erc20types.BridgeToken,
	purseBals []banktypes.Balance,
	supplyExcludeEscrow sdkmath.Int,
	pxEscrowAddr string,
) error {
	erc20ModuleAddress := common.BytesToAddress(types.NewModuleAddress(erc20types.ModuleName).Bytes())
	if _, err := m.erc20TokenKeeper.Mint(ctx, erc20Token.GetERC20Contract(), erc20ModuleAddress, erc20ModuleAddress, supplyExcludeEscrow.BigInt()); err != nil {
		return err
	}
	if err := m.bankKeeper.MintCoins(ctx, erc20types.ModuleName, sdk.NewCoins(sdk.NewCoin(erc20Token.Denom, supplyExcludeEscrow))); err != nil {
		return err
	}
	for _, bal := range purseBals {
		if bal.Address == pxEscrowAddr {
			continue
		}
		bech32Addr, err := sdk.GetFromBech32(bal.Address, "px")
		if err != nil {
			return err
		}
		account := m.accountKeeper.GetAccount(ctx, bech32Addr)
		if ma, ok := account.(sdk.ModuleAccountI); ok {
			// todo migrate staking and reward
			ctx.Logger().Info("migrate PURSE skip module account", "module", ma.GetName(), "address", bal.Address, "balance", bal.Coins.String())
			continue
		}
		purseBalance := sdk.NewCoin(erc20Token.Denom, bal.Coins.AmountOf(bridgeToken.BridgeDenom()))
		if err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, erc20types.ModuleName, bech32Addr, sdk.NewCoins(purseBalance)); err != nil {
			return err
		}
	}
	return nil
}

func (m *Pundix) updatePurseErc20Owner(ctx sdk.Context, erc20Token erc20types.ERC20Token) error {
	erc20ModuleHexAddress := common.BytesToAddress(types.NewModuleAddress(erc20types.ModuleName).Bytes())
	newOwner := common.BytesToAddress(types.NewModuleAddress(evmtypes.ModuleName))
	if _, err := m.erc20TokenKeeper.TransferOwnership(ctx, erc20Token.GetERC20Contract(), erc20ModuleHexAddress, newOwner); err != nil {
		return err
	}
	erc20Token.ContractOwner = erc20types.OWNER_EXTERNAL
	return m.erc20Keeper.ERC20Token.Set(ctx, erc20Token.Denom, erc20Token)
}

func (m *Pundix) removeDeprecatedPurse(ctx sdk.Context, baseDenom, channelId string, ibcToken erc20types.IBCToken) error {
	// remove ibc token
	if err := m.erc20Keeper.DenomIndex.Remove(ctx, ibcToken.IbcDenom); err != nil {
		return err
	}
	key := collections.Join(baseDenom, channelId)
	if err := m.erc20Keeper.IBCToken.Remove(ctx, key); err != nil {
		return err
	}

	// remove denom trace
	hexHash := strings.TrimPrefix(ibcToken.IbcDenom, ibctransfertypes.DenomPrefix+"/")
	hash, err := ibctransfertypes.ParseHexHash(hexHash)
	if err != nil {
		return err
	}
	kvStore := prefix.NewStore(ctx.KVStore(m.storeKey), ibctransfertypes.DenomTraceKey)
	kvStore.Delete(hash)
	return nil
}

func ReadGenesisState(genesisPath string) (map[string]json.RawMessage, error) {
	genesisFile, err := os.ReadFile(genesisPath)
	if err != nil {
		return nil, err
	}
	genesisState := make(map[string]json.RawMessage)
	if err = tmjson.Unmarshal(genesisFile, &genesisState); err != nil {
		return nil, err
	}
	appState := make(map[string]json.RawMessage)
	appStateBz := genesisState["app_state"]
	err = tmjson.Unmarshal(appStateBz, &appState)
	return appState, err
}

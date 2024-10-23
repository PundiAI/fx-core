package cmd

import (
	"fmt"
	"math/big"
	"os"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tendermintos "github.com/cometbft/cometbft/libs/os"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/v8/app"
	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/server"
	fxtypes "github.com/functionx/fx-core/v8/types"
)

type output struct {
	Delegates       map[string]sdkmath.Int `json:"delegates,omitempty"`
	DenomHolders    map[string]sdkmath.Int `json:"denom_holders,omitempty"`
	ContractHolders map[string]sdkmath.Int `json:"contract_holders,omitempty"`
}

func exportDelegatesCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-delegates",
		Short: "Export all delegates and holders",
		Example: fmt.Sprintf(`$ %s export-delegates --contract-addr=<token address> --height=<height> 
--output=out.json --home=<snapshot path>`, version.AppName),
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			clientCtx := client.GetClientContextFromCmd(cmd)

			contractAddr := serverCtx.Viper.GetString("contract-addr")
			denom := serverCtx.Viper.GetString("denom")
			height := serverCtx.Viper.GetInt64("height")

			db, err := server.NewDatabase(serverCtx.Config)
			if err != nil {
				return err
			}

			defer db.Close()

			myApp, ctx, err := buildApp(db.AppDB(), height)
			if err != nil {
				return err
			}

			delegateData := allDelegates(myApp.CommitMultiStore().GetKVStore(myApp.GetKey(stakingtypes.StoreKey)), clientCtx.Codec)

			denomHolders, contractHolders := allHolder(ctx, myApp, contractAddr, denom)

			outData := output{Delegates: delegateData, DenomHolders: denomHolders, ContractHolders: contractHolders}

			encoded, err := tmjson.Marshal(outData)
			if err != nil {
				return err
			}

			out := serverCtx.Viper.GetString("output")

			return tendermintos.WriteFile(out, encoded, os.ModePerm)
		},
	}
	cmd.PersistentFlags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().String("contract-addr", "", "query contract holders, if empty, it will not query contract holders")
	cmd.Flags().String("denom", fxtypes.DefaultDenom, "query denom holders, if empty, it will not query denom holders")
	cmd.Flags().Int64("height", 0, "height to query")
	cmd.Flags().String("output", "out.json", "location of the exported data file")
	return cmd
}

func allDelegates(kvStore storetypes.KVStore, codec codec.Codec) map[string]sdkmath.Int {
	iterator := storetypes.KVStorePrefixIterator(kvStore, stakingtypes.DelegationKey)
	defer iterator.Close()

	delegations := make(map[string]sdkmath.Int)
	for ; iterator.Valid(); iterator.Next() {
		delegation := stakingtypes.MustUnmarshalDelegation(codec, iterator.Value())
		value, ok := delegations[delegation.DelegatorAddress]
		if ok {
			value = value.Add(delegation.Shares.TruncateInt())
		} else {
			value = delegation.Shares.TruncateInt()
		}
		delegations[delegation.DelegatorAddress] = value
	}
	return delegations
}

func allHolder(ctx sdk.Context, myApp *app.App, contractAddrStr, denom string) (map[string]sdkmath.Int, map[string]sdkmath.Int) {
	denomHolder := map[string]sdkmath.Int{}
	contractHolders := map[string]sdkmath.Int{}

	validators, err := myApp.StakingKeeper.GetValidators(ctx, 1)
	if err != nil {
		panic(err)
	}
	consAddr, err := validators[0].GetConsAddr()
	if err != nil {
		panic(err)
	}
	ctx = ctx.WithProposer(consAddr)
	contractAddr := common.HexToAddress(contractAddrStr)

	myApp.AccountKeeper.IterateAccounts(ctx, func(account sdk.AccountI) (stop bool) {
		queryContractBalance(myApp, ctx, contractAddr, common.Address(account.GetAddress()), contractHolders)
		queryDenomBalance(myApp, ctx, account, denom, denomHolder)
		return false
	})
	return denomHolder, contractHolders
}

func queryDenomBalance(myApp *app.App, ctx sdk.Context, account sdk.AccountI, denom string, holder map[string]sdkmath.Int) {
	if len(denom) == 0 {
		return
	}
	balance := myApp.BankKeeper.GetBalance(ctx, account.GetAddress(), denom)
	if balance.IsZero() {
		return
	}
	holder[account.GetAddress().String()] = balance.Amount
}

func queryContractBalance(myApp *app.App, ctx sdk.Context, contractAddr, address common.Address, holders map[string]sdkmath.Int) {
	if contract.IsZeroEthAddress(contractAddr) {
		return
	}
	balance, err := contract.NewERC20TokenKeeper(myApp.EvmKeeper).BalanceOf(ctx, contractAddr, address)
	if err != nil {
		panic(err)
	}
	if balance.Cmp(big.NewInt(0)) == 0 {
		return
	}
	holders[address.Hex()] = sdkmath.NewIntFromBigInt(balance)
}

func buildApp(db dbm.DB, height int64) (*app.App, sdk.Context, error) {
	myApp := app.New(log.NewNopLogger(), db, nil,
		false, map[int64]bool{}, "", app.EmptyAppOptions{})

	if err := myApp.LoadLatestVersion(); err != nil {
		return nil, sdk.Context{}, errors.Wrap(err, "failed to load latest version")
	}
	var multiStore storetypes.CacheMultiStore

	if height > 0 {
		var err error
		multiStore, err = myApp.CommitMultiStore().CacheMultiStoreWithVersion(height)
		if err != nil {
			return nil, sdk.Context{}, errors.Wrapf(err, "failed to load version %d", height)
		}
	} else {
		multiStore = myApp.CommitMultiStore().CacheMultiStore()
	}

	ctx := myApp.NewUncachedContext(false,
		tmproto.Header{Height: myApp.LastBlockHeight()}).WithMultiStore(multiStore)

	return myApp, ctx, nil
}

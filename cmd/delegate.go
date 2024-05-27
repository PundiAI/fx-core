package cmd

import (
	"fmt"
	"math/big"
	"os"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tendermintos "github.com/tendermint/tendermint/libs/os"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/contract"
	"github.com/functionx/fx-core/v7/server"
	fxtypes "github.com/functionx/fx-core/v7/types"
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

			db, err := server.NewDatabase(serverCtx.Config, clientCtx.Codec)
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

func allDelegates(kvStore types.KVStore, codec codec.Codec) map[string]sdkmath.Int {
	iterator := sdk.KVStorePrefixIterator(kvStore, stakingtypes.DelegationKey)
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

func allHolder(ctx sdk.Context, myApp *app.App, contractAddrStr string, denom string) (map[string]sdkmath.Int, map[string]sdkmath.Int) {
	denomHolder := map[string]sdkmath.Int{}
	contractHolders := map[string]sdkmath.Int{}

	consAddr, err := myApp.StakingKeeper.GetValidators(ctx, 1)[0].GetConsAddr()
	if err != nil {
		panic(err)
	}
	ctx = ctx.WithProposer(consAddr)
	contractAddr := common.HexToAddress(contractAddrStr)

	myApp.AccountKeeper.IterateAccounts(ctx, func(account authtypes.AccountI) (stop bool) {
		queryContractBalance(myApp, ctx, contractAddr, common.Address(account.GetAddress()), contractHolders)
		queryDenomBalance(myApp, ctx, account, denom, denomHolder)
		return false
	})
	return denomHolder, contractHolders
}

func queryDenomBalance(myApp *app.App, ctx sdk.Context, account authtypes.AccountI, denom string, holder map[string]sdkmath.Int) {
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

	var balanceRes struct{ Value *big.Int }
	err := myApp.EvmKeeper.QueryContract(ctx, contractAddr, contractAddr, contract.GetFIP20().ABI, "balanceOf", &balanceRes, address)
	if err != nil {
		panic(err)
	}
	if balanceRes.Value.Cmp(big.NewInt(0)) == 0 {
		return
	}
	holders[address.Hex()] = sdkmath.NewIntFromBigInt(balanceRes.Value)
}

func buildApp(db dbm.DB, height int64) (*app.App, sdk.Context, error) {
	myApp := app.New(log.NewFilter(log.NewTMLogger(os.Stdout), log.AllowAll()),
		db, nil, false, map[int64]bool{}, "", 0,
		app.MakeEncodingConfig(), app.EmptyAppOptions{})

	if err := myApp.LoadLatestVersion(); err != nil {
		return nil, sdk.Context{}, errors.Wrap(err, "failed to load latest version")
	}
	var multiStore types.CacheMultiStore

	if height > 0 {
		var err error
		multiStore, err = myApp.CommitMultiStore().CacheMultiStoreWithVersion(height)
		if err != nil {
			return nil, sdk.Context{}, errors.Wrapf(err, "failed to load version %d", height)
		}
	} else {
		multiStore = myApp.CommitMultiStore().CacheMultiStore()
	}

	ctx := myApp.NewUncachedContext(false, tmproto.Header{
		ChainID: fxtypes.ChainId(), Height: myApp.LastBlockHeight(),
	}).WithMultiStore(multiStore)

	return myApp, ctx, nil
}

package keeper_test

import (
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/evmos/ethermint/types"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
)

func (s *KeeperTestSuite) TestKeeper_SetAccount() {
	address := helpers.GenHexAddress()
	s.Nil(s.App.AccountKeeper.GetAccount(s.Ctx, address.Bytes()))

	acc := s.App.EvmKeeper.GetAccountOrEmpty(s.Ctx, address)
	s.NotNil(acc)
	acc.CodeHash = common.BytesToHash([]byte{1, 2, 3}).Bytes()
	s.NoError(s.App.EvmKeeper.SetAccount(s.Ctx, address, acc))

	account := s.App.AccountKeeper.GetAccount(s.Ctx, address.Bytes())
	ethAcc, ok := account.(ethermint.EthAccountI)
	s.True(ok)
	s.Equal(ethAcc.GetCodeHash(), common.BytesToHash(acc.CodeHash))
}

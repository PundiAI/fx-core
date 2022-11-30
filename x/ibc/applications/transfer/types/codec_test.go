package types_test

import (
	"testing"

	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	"github.com/functionx/fx-core/v3/x/ibc/applications/transfer/types"
)

func TestFullData(t *testing.T) {

	ibcPacketData := ibctransfertypes.FungibleTokenPacketData{
		Denom:    "demo",
		Amount:   "amount",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     "memo",
		//Router:   "",
		//Fee:      "",
	}

	ibcPacketDataJson := types.ModuleCdc.MustMarshalJSON(&ibcPacketData)
	t.Logf("%+v", string(ibcPacketDataJson))
	fxPacketData := &types.FungibleTokenPacketData{}
	types.ModuleCdc.MustUnmarshalJSON(ibcPacketDataJson, fxPacketData)
	t.Logf("%+v", fxPacketData)

	//var ibcPacketData = &types.FungibleTokenPacketData{}
	//types.ModuleCdc.MustUnmarshalJSON(fxPacketDataJson, ibcPacketData)

	//t.Logf("%+v", ibcPacketData)
}

func TestFxPacketDataToIBCPacketData(t *testing.T) {

	fxPacketData := types.FungibleTokenPacketData{
		Denom:    "demo",
		Amount:   "amount",
		Sender:   "sender",
		Receiver: "receiver",
		Fee:      "",
		Router:   "",
		Memo:     "memo",
	}

	fxPacketDataJson := fxPacketData.GetBytes()
	t.Logf("%s", string(types.ModuleCdc.MustMarshalJSON(&fxPacketData)))

	t.Logf("%+v", string(fxPacketDataJson))
	ibcPacketData := &ibctransfertypes.FungibleTokenPacketData{}
	types.ModuleCdc.MustUnmarshalJSON(fxPacketDataJson, ibcPacketData)
	t.Logf("%+v", fxPacketData)
}

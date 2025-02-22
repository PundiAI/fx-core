package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

var (
	testnetIBCDenomTrace = map[string]string{
		// pundix
		"ibc/4757BC3AA2C696F7083C825BD3951AE3D1631F2A272EA7AFB9B3E1CCCA8560D4": "channel-0", // purse
		// cosmoshub
		"ibc/C892C98C728A916DFD74A8A6036DF0B6C9B590C813C4CA1E94402D20F7771174": "channel-90", // uatom

		// osmosis
		"ibc/DA227B314073473C4008CDFAE7D6629F5182891AA9A34403BF9CB0F2E380D274": "channel-119", // uosmo
		"ibc/8B7A6637DC5F0921CB24D2F3A051E9797CA50EF728E72872219F0E85A1BFDDD9": "channel-119", // atom/osmo
		"ibc/530063F78415D637217C28ACE7CC7A9ADCC7DD4E7D5A20CE5E143744BE33E685": "channel-119", // stosmo/tosmo
		"ibc/A6D5D1C0A6BBA47E6BBC6C6CBEF531045599345A6CD4FC8B6447909C8D900B33": "channel-119", // tatom/tosmo
		"ibc/C2D7826927DD7FF3EA2F9D3A93889EA09667429847C23007D3583DCEB98CD0BF": "channel-119", // usdc/tosmo
		"ibc/8C8FA24B0CE2E2BE006A1D6CE7D14DD64C463D7887992832C8AE9EF4CEA6FC13": "channel-119", // wbtc/tosmo
		"ibc/3DCFA7233C0E27C07223980CB233EB3D87F64FE4D5DBB88DBC744E0088AF7927": "channel-119", // weth/tosmo
	}
	mainnetIBCDenomTrace = map[string]string{
		// pundix
		"ibc/F08B62C2C1BE9E52942617489CAB1E94537FE3849F8EEC910B142468C340EB0D": "channel-0", // purse
		// cosmoshub
		"ibc/A670D9568B3E399316EEDE40C1181B7AA4BD0695F0B37513CE9B95B977DFC12E": "channel-10", // uatom
		"ibc/19FE4199D5E206A756AA8656D6CCECA3A50AECB081CC73126060D592BD93A15B": "channel-10", // uosmo(from cosmoshub)
		"ibc/04734B17EDC1F68BA599D8387A8BB0266684F46D42F4D5C7667BF370D21D9B8E": "channel-10", // aevmos
		"ibc/4E7EFFD6C691F22FEA1C3B72DB3C275AF595A48C769CBB79289BDBD109A1F39A": "channel-10", // inj
		"ibc/B61B71A43BA8D837D866B6417CF65B391A9177B3C3391527BF88A1174C55D939": "channel-10", // uaxl
		"ibc/DE7DEEC0A8B377C9C3292884B80BDF264BEE97EF2F85DABE2E46B00A2E15C29A": "channel-10", // ukava
		"ibc/8136D7DCECF799F835497A02725092811BA22810F88A1D0678600E13CEB249BE": "channel-10", // uscrt
		"ibc/5A730E758B54FDEFAC2696BC6A24342192DF8F3D2B0B947C6D26AF53E9E943CA": "channel-10", // ustrd

		// osmosis
		"ibc/D7B22A85AB15F44A3152EBF7F2D37B6061F66FAF637E42C287FC649F1F5CA348": "channel-19", // uosmo
		"ibc/95B9D47D7890C9F3E9104DD1DA74D3CCD4890E453FE0F0A711E782CE77B13455": "channel-19", // ulvn(from osmosis)
		"ibc/BE612CFB5445AD2F56FAE496C0848FA381F967970A4CA586B1F90123AA62C4D4": "channel-19", // atom/osmo
		"ibc/1B68E41D8D074F645388189E46C0F480111936F136C06E57FF559AB052D5BF78": "channel-19", // akt/atom
		"ibc/145B926094C7649D85679D577FE8D3FD713C8586346518BF838712D4498A37E3": "channel-19", // akt/osmo
		"ibc/49ABFFC9B2823450EDF1D4548B5B6314B473E497CF454726575458E7FDF3C52B": "channel-19", // atom/qatom
		"ibc/373E4EFF2671B14802D9B3AB3EAF0345FBED921949654EB115C5072DFEB93AB0": "channel-19", // atom/statom
		"ibc/AC80A1534F33EBA754A4EFC4E6AC3672850CCD0CDE9666B635984FB41BE91DA1": "channel-19", // atom/stkatom
		"ibc/16770A91B7849DAFC413E7094E2CD66A8E908DB8550D7AD7F9B2115F7E00AC33": "channel-19", // axl/osmo
		"ibc/A283012C5CA64DE2B8C6B4979B9A99E5144AA30F707D0FE16ED3ADD0F0BEDB42": "channel-19", // cro/osmo
		"ibc/E30827509649DBDA201E98E14EDD52116124D2BCFAF70BBA5E7CA2615E086A2C": "channel-19", // dai/osmo
		"ibc/30033E46154DB10107B18AE42AFFEBD1FAF4C36AF0142F1156B0ECC1796227FE": "channel-19", // eth/osmo
		"ibc/86C7F2DD93F556650BF0992A2F61CBBA42600BD7FFD3F98CCE74C9A1994EB565": "channel-19", // evmos/osmo
		"ibc/97B0CC9611DFE37744BA8EE27B32DE9509FED9ECDFCE59434E0B549BAD67D83E": "channel-19", // inj/osmo
		"ibc/E126BB8628A2259291A668E2C666C97D60D06C59717BD2E34290A6DEF8AA8D19": "channel-19", // juno/osmo
		"ibc/1CC704EB23F1B0EAA9634F088E2F476937E3B4519782B996B63F096E69DBAF5B": "channel-19", // scrt/osmo
		"ibc/DFF778E0F743B1FE66D6018051AB5395F7CB470B19673149F1422B9DE68CD62C": "channel-19", // stars/osmo
		"ibc/69058D0A3D0E3E63DC4AA39571007277A16F968290084BD80AFFFA2831B4AD10": "channel-19", // stosmo/osmo
		"ibc/FDF9794382805B73B7182E1155E9A74EAA0CB37CFEC547B66335E21CC55829C2": "channel-19", // strd/osmo
		"ibc/D54952FCA1DAE3919B89F2D809A0400C372D5919480BA2FDC9FF4217384DA7CE": "channel-19", // usdc/nls
		"ibc/8C412C753AF39C7A050A250FE4DD3B4975AE6B714F3D922E94E0DFA5F81EBEA0": "channel-19", // usdc/osmo
		"ibc/A655D0A048B0FD73458D9E6F020ABBBC3DD39F2D3AEACD2DCD1B0EC26B5EC425": "channel-19", // wbtc/osmo

		// chihuahua
		"ibc/210FA8AD411B627A0EFAEF4580B0D61C707C2C9E138AAF5AA087B4864B1E3599": "channel-22", // uhuahua
	}

	testnetExcludeBridgeToken = map[string]bool{
		"layer20xb1efb300876f993Dc6826f09E66FEaa1bc3A735F": true, // DEGEN
	}
)

func GetIBCDenomTrace(ctx sdk.Context, denom string) (trace string, ok bool) {
	if ctx.ChainID() == fxtypes.TestnetChainId {
		trace, ok = testnetIBCDenomTrace[denom]
	} else {
		trace, ok = mainnetIBCDenomTrace[denom]
	}
	return
}

func getExcludeBridgeToken(ctx sdk.Context, denom string) bool {
	if ctx.ChainID() == fxtypes.TestnetChainId {
		return testnetExcludeBridgeToken[denom]
	}
	return false
}

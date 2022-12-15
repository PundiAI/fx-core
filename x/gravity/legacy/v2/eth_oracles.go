package v2

import (
	fxtypes "github.com/functionx/fx-core/v3/types"
)

func EthInitOracles(chainId string) []string {
	if chainId == fxtypes.MainnetChainId {
		return []string{}
	} else if chainId == fxtypes.TestnetChainId {
		return []string{
			"fx1rcac2frfhfdt3v63jkf24z0kr62qe6pa6d4ped",
			"fx1vqdmvj68gfqyyqs47plje5f06vmccpll4j96pw",
			"fx1eucw3ekl0a4fr65cpejqjz8gj87zr9et8dn20h",
			"fx1p4haqg04mvfhrum0w4tdxd5xgn3rmfmwzs48fk",
			"fx12jxpmejkmwz8rhnxms4mx48mccmqlkct3k5tp7",
			"fx137pzz666wg73awhrgfavj2lye5pkspfckekxu6",
			"fx13ma2nfjxkzjkvj398fygvxfxsss3dwjzmeys0t",
			"fx1zcgh03svy05z3yxm5pwuhms3g8ha5wglg87qw5",
			"fx1kx6djlcqfug335fgxhu5tehl8f2ah8y9zfwgqn",
			"fx177mcyuca86zjefm6ev9jw59n3wxuh2m375srsu",
			"fx1jlgmk30cn7sz3q8wumadygstdkm59npn5kwh5q",
			"fx16hzd2vwmnn6rtm2v6lrmatcd5hzjat59yrldn3",
			"fx16d7nkmnzx3lzsxtk0kqh2z0f7c0wpym887e845",
			"fx1axukfjpg9djsd6cycycup07s8jr2w2ztyha43v",
			"fx1puj4lwchadaddpqhxym83cjl63y5e7av5yzxms",
			"fx1ge55nq9p6qtysz88dw252zy8h3zlprwweqjcv0",
			"fx1h4pfgmjka8stdwmqqkgkcn4cgtxy76ygl4hddr",
			"fx10ur36g0tmtqgzvf09t9qz3u6sl2lyftlx2wm60",
			"fx1upz23pc492008xy3ymks5a4y73zp0yf0dhn9ae",
			"fx16wd68vzz6tqdapfnsuewcm8lkuze7lwen4w0x9",
		}
	} else {
		panic("invalid chainId:" + chainId)
	}
}

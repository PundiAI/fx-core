package v3

import (
	fxtypes "github.com/functionx/fx-core/v3/types"
)

func GetEthOracleAddrs(chainId string) []string {
	if chainId == fxtypes.MainnetChainId {
		return []string{
			"fx1p5kzz68xxncdhjct497kf5d560jm05c9slqct2",
			"fx1knv2jpzhjea4ttknhasl9vreppf2xmqptj4e6q",
			"fx16p2mhts39vwzfttekq3r65w05vu6helnzhqj4g",
			"fx1cmcm3dcmsk37hd4vtrtmqur0q6rk63z3q6cxtx",
			"fx1mwugh372ynjm47ne70cy6k0s4yhmv6rdk729ly",
			"fx1tu70dscm75pgwk4kp8v7mn87kg3j0da8tptfq2",
			"fx1pxd9uepcwvjv5lcpk9kr64g20jd3e3q0d9zh4z",
			"fx1krwalhaw7l6kcatv9wmrh605fhqs45ccwgpzwp",
			"fx1zgu07t9k3fhkqgnajsqhknhdrtprlttuu8jhxd",
			"fx1jz42z9y39va542k4lnxamcjxscetd2ahyflrvp",
			"fx1qvss2fj9nffkdxy5f9eqpt0wax29tcpd9xzwaq",
			"fx1vjdmv8aexvu65em3nfl7z8m9u4zzqwvhh0e6lj",
			"fx1uk6tcw6qppuk4t479afsx55xmpfzsn6d8g4hma",
			"fx15gwu4wjazahjgy924r86lfk4tm5kuk5788dkpx",
			"fx1hpce5xhahppsn4u76fkuvzg0ev0d425xls3lt6",
			"fx12ea7cg6hkf702qc32g5474a7zjdx6z4w92ttdg",
			"fx1jkne5mxqac0amudeyf4u8grex7qqmeaq7yhwfa",
			"fx10k08xgjmesrmmn3xzrjdx4scv8ky7wlrdfe3h7",
			"fx1zv782q995uwnklgmkdgyzd4hn937ln7x6de0se",
			"fx135mnzqwtasydflfeu6hjvcygf9ags2sy6w08c6",
		}
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

module github.com/functionx/fx-core

go 1.16

require (
	github.com/armon/go-metrics v0.3.6
	github.com/cosmos/cosmos-sdk v0.42.4
	github.com/cosmos/go-bip39 v1.0.0
	github.com/ethereum/go-ethereum v1.10.9
	github.com/fbsobreira/gotron-sdk v0.0.0-20211012084317-763989224068
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.0
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/kr/text v0.2.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/pelletier/go-toml v1.9.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.9
	github.com/tendermint/tm-db v0.6.4
	github.com/tyler-smith/go-bip39 v1.0.2
	google.golang.org/genproto v0.0.0-20210114201628-6edceaf6022f
	google.golang.org/grpc v1.38.0
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4

replace github.com/cosmos/cosmos-sdk => github.com/functionx/cosmos-sdk v0.42.5-0.20211125094332-429e0931442e

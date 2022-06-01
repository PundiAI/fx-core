package evm

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/rpc/client/http"
	"google.golang.org/grpc"

	"github.com/functionx/fx-core/app"
	cryptohd "github.com/functionx/fx-core/crypto/hd"
	fxtypes "github.com/functionx/fx-core/types"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/x/erc20/types"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
	gravitytypes "github.com/functionx/fx-core/x/gravity/types"
)

const (
	DefaultGRPCUrl    = "127.0.0.1:9090"
	DefaultNodeRPCUrl = "tcp://localhost:26657"
	DefaultEthUrl     = "http://127.0.0.1:8545"

	DefaultMnemonic   = "dune antenna hood magic kit blouse film video another pioneer dilemma hobby message rug sail gas culture upgrade twin flag joke people general aunt"
	DefaultHexAddress = "0x61bd2030908d658dd5a2139D2C13Af55b9138efb"
	EthHDPath         = "m/44'/60'/0'/0/0"
	CosmosHDPath      = "m/44'/118'/0'/0/0"

	EthFXTokenContract = "0x0000000000000000000000000000000000000001"
)

var (
	txModeInfo = &tx.ModeInfo{
		Sum: &tx.ModeInfo_Single_{
			Single: &tx.ModeInfo_Single{Mode: signing.SignMode_SIGN_MODE_DIRECT},
		},
	}

	ERC20ABI, _ = abi.JSON(strings.NewReader("[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"tokenName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"tokenSymbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"tokenDecimals\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"))
	ERC20Bin, _ = hex.DecodeString("60806040523480156200001157600080fd5b5060405162000ca238038062000ca2833981016040819052620000349162000208565b6200004460ff8216600a620002df565b620000509085620003d4565b600381905533600090815260046020908152604082209290925584516200007a92860190620000af565b50815162000090906001906020850190620000af565b506002805460ff191660ff92909216919091179055506200045f915050565b828054620000bd90620003f6565b90600052602060002090601f016020900481019282620000e157600085556200012c565b82601f10620000fc57805160ff19168380011785556200012c565b828001600101855582156200012c579182015b828111156200012c5782518255916020019190600101906200010f565b506200013a9291506200013e565b5090565b5b808211156200013a57600081556001016200013f565b600082601f83011262000166578081fd5b81516001600160401b038082111562000183576200018362000449565b604051601f8301601f19908116603f01168101908282118183101715620001ae57620001ae62000449565b81604052838152602092508683858801011115620001ca578485fd5b8491505b83821015620001ed5785820183015181830184015290820190620001ce565b83821115620001fe57848385830101525b9695505050505050565b600080600080608085870312156200021e578384fd5b845160208601519094506001600160401b03808211156200023d578485fd5b6200024b8883890162000155565b9450604087015191508082111562000261578384fd5b50620002708782880162000155565b925050606085015160ff8116811462000287578182fd5b939692955090935050565b80825b6001808611620002a65750620002d6565b818704821115620002bb57620002bb62000433565b80861615620002c957918102915b9490941c93800262000295565b94509492505050565b6000620002f06000198484620002f7565b9392505050565b6000826200030857506001620002f0565b816200031757506000620002f0565b81600181146200033057600281146200033b576200036f565b6001915050620002f0565b60ff8411156200034f576200034f62000433565b6001841b91508482111562000368576200036862000433565b50620002f0565b5060208310610133831016604e8410600b8410161715620003a7575081810a83811115620003a157620003a162000433565b620002f0565b620003b6848484600162000292565b808604821115620003cb57620003cb62000433565b02949350505050565b6000816000190483118215151615620003f157620003f162000433565b500290565b6002810460018216806200040b57607f821691505b602082108114156200042d57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052604160045260246000fd5b610833806200046f6000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c806342966c681161007157806342966c681461013857806370a082311461014b57806395d89b411461016b578063a0712d6814610173578063a9059cbb14610186578063dd62ed3e14610199576100a9565b806306fdde03146100ae578063095ea7b3146100cc57806318160ddd146100ef57806323b872dd14610106578063313ce56714610119575b600080fd5b6100b66101c4565b6040516100c3919061072a565b60405180910390f35b6100df6100da3660046106e9565b610252565b60405190151581526020016100c3565b6100f860035481565b6040519081526020016100c3565b6100df6101143660046106ae565b6102be565b6002546101269060ff1681565b60405160ff90911681526020016100c3565b6100df610146366004610712565b61033c565b6100f861015936600461065b565b60046020526000908152604090205481565b6100b66103d7565b6100df610181366004610712565b6103e4565b6100df6101943660046106e9565b610458565b6100f86101a736600461067c565b600560209081526000928352604080842090915290825290205481565b600080546101d1906107ac565b80601f01602080910402602001604051908101604052809291908181526020018280546101fd906107ac565b801561024a5780601f1061021f5761010080835404028352916020019161024a565b820191906000526020600020905b81548152906001019060200180831161022d57829003601f168201915b505050505081565b3360008181526005602090815260408083206001600160a01b038716808552925280832085905551919290917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925906102ad9086815260200190565b60405180910390a350600192915050565b6001600160a01b03831660009081526005602090815260408083203384529091528120548211156102ee57600080fd5b6001600160a01b038416600090815260056020908152604080832033845290915281208054849290610321908490610795565b90915550610332905084848461046e565b5060019392505050565b3360009081526004602052604081205482111561035857600080fd5b3360009081526004602052604081208054849290610377908490610795565b9250508190555081600360008282546103909190610795565b909155505060405182815260009033907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef906020015b60405180910390a35060015b919050565b600180546101d1906107ac565b3360009081526004602052604081208054839190839061040590849061077d565b92505081905550816003600082825461041e919061077d565b909155505060405182815233906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef906020016103c6565b600061046533848461046e565b50600192915050565b6001600160a01b0382166104c85760405162461bcd60e51b815260206004820152601860248201527f7472616e7366657220746f207a65726f20616464726573730000000000000000604482015260640160405180910390fd5b6001600160a01b0383166000908152600460205260409020548111156104ed57600080fd5b6001600160a01b038216600090815260046020526040902054610510828261077d565b101561051b57600080fd5b6001600160a01b0380831660009081526004602052604080822054928616825281205490916105499161077d565b6001600160a01b038516600090815260046020526040812080549293508492909190610576908490610795565b90915550506001600160a01b038316600090815260046020526040812080548492906105a390849061077d565b90915550506001600160a01b0380841660009081526004602052604080822054928716825290205482916105d69161077d565b146105f157634e487b7160e01b600052600160045260246000fd5b826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161063691815260200190565b60405180910390a350505050565b80356001600160a01b03811681146103d257600080fd5b60006020828403121561066c578081fd5b61067582610644565b9392505050565b6000806040838503121561068e578081fd5b61069783610644565b91506106a560208401610644565b90509250929050565b6000806000606084860312156106c2578081fd5b6106cb84610644565b92506106d960208501610644565b9150604084013590509250925092565b600080604083850312156106fb578182fd5b61070483610644565b946020939093013593505050565b600060208284031215610723578081fd5b5035919050565b6000602080835283518082850152825b818110156107565785810183015185820160400152820161073a565b818111156107675783604083870101525b50601f01601f1916929092016040019392505050565b60008219821115610790576107906107e7565b500190565b6000828210156107a7576107a76107e7565b500390565b6002810460018216806107c057607f821691505b602082108114156107e157634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fdfea2646970667358221220a12642861f022af61137e950b14bdece7118dc7d0ff204c0029c52361b54baeb64736f6c63430008020033")

	FIP20ABI, _ = abi.JSON(strings.NewReader("[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"target\",\"type\":\"bytes32\"}],\"name\":\"TransferCrossChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"module_\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"module\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"target\",\"type\":\"bytes32\"}],\"name\":\"transferCrossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]"))
	WFXABI, _   = abi.JSON(strings.NewReader("[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"target\",\"type\":\"bytes32\"}],\"name\":\"TransferCrossChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"module_\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"module\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"target\",\"type\":\"bytes32\"}],\"name\":\"transferCrossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address payable\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"))
)

type Client struct {
	t               *testing.T
	ctx             context.Context
	chainId         string
	evmChainId      *big.Int
	privateKey      cryptotypes.PrivKey
	encodingConfig  app.EncodingConfig
	ethClient       *ethclient.Client
	httpClient      *http.HTTP
	grpcClient      *grpc.ClientConn
	txClient        tx.ServiceClient
	authQueryClient authtypes.QueryClient
	evmQuery        evmtypes.QueryClient
	erc20Query      erc20types.QueryClient
	gravityQuery    gravitytypes.QueryClient
	crossChainQuery crosschaintypes.QueryClient
	ibcChannelQuery ibcchanneltypes.QueryClient
}

func NewClient(t *testing.T, grpcUrl, nodeUrl, ethUrl, mnemonic, hdPath string) *Client {
	ctx := context.Background()
	grpcClient, err := grpcNewClient(grpcUrl)
	require.NoError(t, err)
	httpClient, err := newHttpClient(nodeUrl)
	require.NoError(t, err)
	ethClient, err := ethNewClient(ctx, ethUrl)
	require.NoError(t, err)
	if len(mnemonic) == 0 {
		mnemonic = DefaultMnemonic
	}
	if len(hdPath) == 0 {
		hdPath = EthHDPath
	}
	privateKey, err := mnemonicToPrivKey(mnemonic, hdPath)
	require.NoError(t, err)

	status, err := httpClient.Status(ctx)
	require.NoError(t, err)

	evmChainId, err := ethClient.ChainID(ctx)
	require.NoError(t, err)

	return &Client{
		t:               t,
		ctx:             ctx,
		chainId:         status.NodeInfo.Network,
		evmChainId:      evmChainId,
		httpClient:      httpClient,
		privateKey:      privateKey,
		encodingConfig:  app.MakeEncodingConfig(),
		ethClient:       ethClient,
		grpcClient:      grpcClient,
		txClient:        tx.NewServiceClient(grpcClient),
		authQueryClient: authtypes.NewQueryClient(grpcClient),
		evmQuery:        evmtypes.NewQueryClient(grpcClient),
		erc20Query:      erc20types.NewQueryClient(grpcClient),
		gravityQuery:    gravitytypes.NewQueryClient(grpcClient),
		crossChainQuery: crosschaintypes.NewQueryClient(grpcClient),
		ibcChannelQuery: ibcchanneltypes.NewQueryClient(grpcClient),
	}
}

func (c *Client) SetKey(mnemonic, hdPath string) *Client {
	key, err := mnemonicToPrivKey(mnemonic, hdPath)
	require.NoError(c.t, err)
	c.privateKey = key
	return c
}

func (c *Client) AccAddress(mnemonic ...string) sdk.AccAddress {
	priKey := c.privateKey
	if len(mnemonic) > 0 {
		var err error
		priKey, err = mnemonicToPrivKey(mnemonic[0], EthHDPath)
		require.NoError(c.t, err)
	}
	return sdk.AccAddress(priKey.PubKey().Address())
}
func (c *Client) HexAddress(mnemonic ...string) common.Address {
	priKey := c.privateKey
	if len(mnemonic) > 0 {
		var err error
		priKey, err = mnemonicToPrivKey(mnemonic[0], EthHDPath)
		require.NoError(c.t, err)
	}
	return common.BytesToAddress(priKey.PubKey().Address())
}
func (c *Client) TransactOpts(mnemonic ...string) *bind.TransactOpts {
	priKey := c.privateKey
	if len(mnemonic) > 0 {
		var err error
		priKey, err = mnemonicToPrivKey(mnemonic[0], EthHDPath)
		require.NoError(c.t, err)
	}
	ecdsa, err := crypto.ToECDSA(priKey.Bytes())
	require.NoError(c.t, err)

	transactOpts, err := bind.NewKeyedTransactorWithChainID(ecdsa, c.evmChainId)
	require.NoError(c.t, err)

	return transactOpts
}

func (c *Client) Balance(addr common.Address) *big.Int {
	at, err := c.ethClient.BalanceAt(c.ctx, addr, nil)
	require.NoError(c.t, err)
	return at
}
func (c *Client) PendingTx(tx *ethtypes.Transaction) {
	c.t.Log("pending tx hash", tx.Hash())
	ctx, cancel := context.WithTimeout(c.ctx, 30*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, c.ethClient, tx)
	require.NoError(c.t, err)
	require.Equal(c.t, receipt.Status, ethtypes.ReceiptStatusSuccessful)
}
func (c *Client) Token(denom string) common.Address {
	pair, err := c.erc20Query.TokenPair(c.ctx, &erc20types.QueryTokenPairRequest{Token: denom})
	require.NoError(c.t, err)
	return pair.GetTokenPair().GetERC20Contract()
}
func (c *Client) Transfer(recipient common.Address, value *big.Int) common.Hash {
	c.t.Logf("transfer to %s value %s\n", recipient.String(), value.String())
	ethTx, err := dynamicFeeTx(c.ethClient, c.privateKey, &recipient, value, nil)
	require.NoError(c.t, err)

	err = c.ethClient.SendTransaction(c.ctx, ethTx)
	require.NoError(c.t, err)

	c.PendingTx(ethTx)

	return ethTx.Hash()
}

func (c *Client) BuildTx(msgList []sdk.Msg) *tx.TxRaw {
	accountResponse, err := c.authQueryClient.Account(c.ctx, &authtypes.QueryAccountRequest{Address: c.AccAddress().String()})
	require.NoError(c.t, err)
	var account authtypes.AccountI
	err = c.encodingConfig.InterfaceRegistry.UnpackAny(accountResponse.GetAccount(), &account)
	require.NoError(c.t, err)

	//c.t.Log("account ===>", account.GetAddress().String(), account.GetAccountNumber())

	txBodyMessage := make([]*types.Any, 0)
	for i := 0; i < len(msgList); i++ {
		msgAnyValue, err := types.NewAnyWithValue((msgList)[i])
		if err != nil {
			c.t.Fatal(err)
		}
		txBodyMessage = append(txBodyMessage, msgAnyValue)
	}
	txBody := &tx.TxBody{
		Messages: txBodyMessage,
	}
	pubKey, err := codectypes.NewAnyWithValue(c.privateKey.PubKey())
	require.NoError(c.t, err)
	authInfo := &tx.AuthInfo{
		SignerInfos: []*tx.SignerInfo{
			{
				PublicKey: pubKey,
				ModeInfo:  txModeInfo,
				Sequence:  account.GetSequence(),
			},
		},
		Fee: &tx.Fee{
			GasLimit: 0,
		},
	}
	txBodyBytes := mustProtoMarshal(c.t, txBody)

	signResultBytes := sign(c.t, c.privateKey, &tx.SignDoc{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: nil,
		ChainId:       c.chainId,
		AccountNumber: account.GetAccountNumber(),
	})
	simulateResponse, err := c.txClient.Simulate(context.Background(), &tx.SimulateRequest{
		Tx: &tx.Tx{
			Body:       txBody,
			AuthInfo:   authInfo,
			Signatures: [][]byte{signResultBytes},
		},
	})
	require.NoError(c.t, err)

	// adjustment gasLimit 1.3 .
	var gasLimit uint64 = simulateResponse.GasInfo.GasUsed * 13 / 10
	gasPrice, _ := sdk.NewIntFromString("4000000000000")
	gasFeeAmount := gasPrice.Mul(sdk.NewInt(int64(gasLimit)))
	authInfo.Fee = &tx.Fee{
		Amount:   sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, gasFeeAmount)),
		GasLimit: gasLimit,
	}
	txAuthInfoBytes := mustProtoMarshal(c.t, authInfo)

	signResultBytes = sign(c.t, c.privateKey, &tx.SignDoc{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		ChainId:       c.chainId,
		AccountNumber: account.GetAccountNumber(),
	})

	txRaw := &tx.TxRaw{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		Signatures:    [][]byte{signResultBytes},
	}

	return txRaw
}
func (c *Client) ConvertCoin(recipient common.Address, coin sdk.Coin) {
	msg := erc20types.NewMsgConvertCoin(coin, recipient, c.AccAddress())
	txRaw := c.BuildTx([]sdk.Msg{msg})
	broadcastTx(c.t, c.ctx, c.txClient, txRaw)
}
func (c *Client) ConvertERC20(token common.Address, amount sdk.Int, recipient sdk.AccAddress) {
	if !amount.GT(sdk.ZeroInt()) {
		return
	}
	msg := erc20types.NewMsgConvertERC20(amount, recipient, token, c.HexAddress())
	txRaw := c.BuildTx([]sdk.Msg{msg})
	broadcastTx(c.t, c.ctx, c.txClient, txRaw)
}
func (c *Client) TransferCrossChain(token common.Address, recipient string, amount, fee *big.Int, target string) common.Hash {
	c.t.Log("transfer cross chain", target)
	pack, err := FIP20ABI.Pack("transferCrossChain", recipient, amount, fee, fxtypes.StringToByte32(target))
	require.NoError(c.t, err)

	ethTx, err := dynamicFeeTx(c.ethClient, c.privateKey, &token, nil, pack)
	require.NoError(c.t, err)

	err = c.ethClient.SendTransaction(c.ctx, ethTx)
	require.NoError(c.t, err)

	c.PendingTx(ethTx)

	return ethTx.Hash()
}
func (c *Client) Deposit(wfx common.Address, amount *big.Int) common.Hash {
	pack, err := WFXABI.Pack("deposit")
	require.NoError(c.t, err)

	ethTx, err := dynamicFeeTx(c.ethClient, c.privateKey, &wfx, amount, pack)
	require.NoError(c.t, err)

	err = c.ethClient.SendTransaction(c.ctx, ethTx)
	require.NoError(c.t, err)

	c.PendingTx(ethTx)

	return ethTx.Hash()
}
func (c *Client) Withdraw(wfx, recipient common.Address, value *big.Int) common.Hash {
	pack, err := WFXABI.Pack("withdraw", recipient, value)
	require.NoError(c.t, err)

	ethTx, err := dynamicFeeTx(c.ethClient, c.privateKey, &wfx, nil, pack)
	require.NoError(c.t, err)

	err = c.ethClient.SendTransaction(c.ctx, ethTx)
	require.NoError(c.t, err)

	c.PendingTx(ethTx)

	return ethTx.Hash()
}

// EthBlockHeight assume it eth height
func (c *Client) EthBlockHeight() uint64 {
	number, err := c.ethClient.BlockNumber(c.ctx)
	require.NoError(c.t, err)
	return number
}

func (c *Client) GravityInitialize() {
	c.t.Log("gravity initialize")
	if c.gravityCheckValset() {
		return
	}
	go c.gravityValsetConfirm()

	c.gravitySetOrchestratorAddress()
	c.gravityFxOriginatedTokenClaim()

	c.gravityCheckConfirm()
}
func (c *Client) gravityId() string {
	var chainGravityId string
	abciQuery, err := c.httpClient.ABCIQuery(c.ctx, "/custom/gravity/gravityID", bytes.HexBytes{})
	require.NoError(c.t, err)
	require.Equal(c.t, abciQuery.Response.Code, uint32(0))
	err = c.encodingConfig.Amino.UnmarshalJSON(abciQuery.Response.Value, &chainGravityId)
	require.NoError(c.t, err)
	return chainGravityId
}
func (c *Client) gravityValsetConfirm() {
	gravityId := c.gravityId()
	requestParams := &gravitytypes.QueryLastPendingValsetRequestByAddrRequest{Address: c.AccAddress().String()}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		queryResponse, err := c.gravityQuery.LastPendingValsetRequestByAddr(c.ctx, requestParams)
		require.NoError(c.t, err)

		valsets := queryResponse.Valsets
		if len(valsets) <= 0 {
			continue
		}

		for _, valset := range valsets {
			checkpoint := valset.GetCheckpoint(gravityId)
			//c.t.Logf("need confirm valset: nonce:%v EthAddress:%v\n", valset.Nonce, c.HexAddress())
			ecdsa, _ := crypto.ToECDSA(c.privateKey.Bytes())
			signature, err := gravitytypes.NewEthereumSignature(checkpoint, ecdsa)
			require.NoError(c.t, err)

			txRaw := c.BuildTx([]sdk.Msg{
				&gravitytypes.MsgValsetConfirm{
					Nonce:        valset.Nonce,
					Orchestrator: c.AccAddress().String(),
					EthAddress:   c.HexAddress().String(),
					Signature:    hex.EncodeToString(signature),
				},
			})
			broadcastTx(c.t, c.ctx, c.txClient, txRaw)
		}
	}
}
func (c *Client) gravityCheckConfirm() {
	count := 0
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		resp, err := c.gravityQuery.LastPendingValsetRequestByAddr(c.ctx, &gravitytypes.QueryLastPendingValsetRequestByAddrRequest{Address: c.AccAddress().String()})
		require.NoError(c.t, err)
		if len(resp.Valsets) > 0 {
			count++
			require.True(c.t, count <= 3)
			continue
		}
		break
	}
}
func (c *Client) gravityCheckValset() bool {
	respValset, err := c.gravityQuery.CurrentValset(c.ctx, &gravitytypes.QueryCurrentValsetRequest{})
	require.NoError(c.t, err)
	for _, m := range respValset.GetValset().GetMembers() {
		if strings.EqualFold(m.EthAddress, c.HexAddress().String()) {
			return true
		}
	}
	return false
}
func (c *Client) gravitySetOrchestratorAddress() {
	c.t.Log("gravity setOrchestratorAddress")
	queryOrchestratorResponse, err := c.gravityQuery.GetDelegateKeyByOrchestrator(c.ctx,
		&gravitytypes.QueryDelegateKeyByOrchestratorRequest{OrchestratorAddress: c.AccAddress().String()})
	if queryOrchestratorResponse != nil && len(queryOrchestratorResponse.EthAddress) > 0 {
		c.t.Logf("already set orchestrator address! address:[%v], validatorAddress:[%v], ethAddress:[%v]\n",
			c.AccAddress().String(), queryOrchestratorResponse.ValidatorAddress, queryOrchestratorResponse.EthAddress)
		return
	}
	//require.NoError(c.t, err)
	if err != nil {
		if !strings.Contains(err.Error(), "No validator") {
			c.t.Fatal(err)
		}
		c.t.Logf("not found validator!!error msg:%v\n", err.Error())
	}
	msgSetOrchestratorAddress := gravitytypes.NewMsgSetOrchestratorAddress(sdk.ValAddress(c.AccAddress()), c.AccAddress(), c.HexAddress().String())
	txRaw := c.BuildTx([]sdk.Msg{msgSetOrchestratorAddress})
	broadcastTx(c.t, c.ctx, c.txClient, txRaw)
}
func (c *Client) gravityQueryFxLastEventNonce() uint64 {
	lastEventNonce, err := c.gravityQuery.LastEventNonceByAddr(c.ctx, &gravitytypes.QueryLastEventNonceByAddrRequest{Address: c.AccAddress().String()})
	if err != nil {
		c.t.Fatal(err)
	}
	return lastEventNonce.EventNonce + 1
}
func (c *Client) GravityCheckPoolTx(addr ...string) {
	sender := c.AccAddress().String()
	if len(addr) > 0 {
		sender = addr[0]
	}
	resp, err := c.gravityQuery.GetPendingSendToEth(c.ctx, &gravitytypes.QueryPendingSendToEthRequest{SenderAddress: sender})
	require.NoError(c.t, err)
	require.True(c.t, len(resp.UnbatchedTransfers) > 0)

	for _, ub := range resp.UnbatchedTransfers {
		c.t.Log("unBatch id", ub.Id)
	}
}
func (c *Client) gravityFxOriginatedTokenClaim() {
	c.t.Log("gravity setFxOriginatedTokenClaim")
	resp, err := c.gravityQuery.LastEventNonceByAddr(c.ctx, &gravitytypes.QueryLastEventNonceByAddrRequest{Address: c.AccAddress().String()})
	require.NoError(c.t, err)

	msg := &gravitytypes.MsgFxOriginatedTokenClaim{
		EventNonce:    resp.EventNonce + 1,
		BlockHeight:   c.EthBlockHeight(),
		TokenContract: EthFXTokenContract,
		Name:          "Function X",
		Symbol:        "FX",
		Decimals:      18,
		Orchestrator:  c.AccAddress().String(),
	}
	txRaw := c.BuildTx([]sdk.Msg{msg})
	broadcastTx(c.t, c.ctx, c.txClient, txRaw)
}
func (c *Client) GravitySendToTx(recipient sdk.AccAddress, token common.Address, value *big.Int, targetIBC string) {
	_, err := c.gravityQuery.ERC20ToDenom(c.ctx, &gravitytypes.QueryERC20ToDenomRequest{Erc20: token.String()})
	require.NoError(c.t, err)

	if len(targetIBC) > 0 {
		targetIBC = hex.EncodeToString([]byte(targetIBC))
	}
	depositClaimMsg := gravitytypes.NewMsgDepositClaim(c.gravityQueryFxLastEventNonce(), c.EthBlockHeight(), token.String(),
		sdk.NewIntFromBigInt(value), c.HexAddress().String(), recipient.String(), targetIBC, c.AccAddress().String())
	txRaw := c.BuildTx([]sdk.Msg{depositClaimMsg})
	broadcastTx(c.t, c.ctx, c.txClient, txRaw)
}

func (c *Client) CheckIBCChannelState(portId, channelId string) {
	_, err := c.ibcChannelQuery.ChannelClientState(c.ctx, &ibcchanneltypes.QueryChannelClientStateRequest{PortId: portId, ChannelId: channelId})
	require.NoError(c.t, err)
}

func grpcNewClient(grpcUrl string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	// http use  grpc.WithInsecure()
	opts = append(opts, grpc.WithInsecure())
	// https use grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "*.function.io"))
	//opts = append(opts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, strings.Split(grpcUrl, ":")[0])))
	return grpc.Dial(grpcUrl, opts...)
}
func newHttpClient(rpcUrl string) (*http.HTTP, error) {
	return http.New(rpcUrl, "/websocket")
}
func ethNewClient(ctx context.Context, rpcUrl string) (*ethclient.Client, error) {
	return ethclient.DialContext(ctx, rpcUrl)
}
func mnemonicToPrivKey(mnemonic, hdPath string) (cryptotypes.PrivKey, error) {
	algo := cryptohd.EthSecp256k1
	bytes, err := algo.Derive()(mnemonic, "", hdPath)
	if err != nil {
		return nil, err
	}
	privKey := algo.Generate()(bytes)
	return privKey, nil
}
func dynamicFeeTx(cli *ethclient.Client, priKey cryptotypes.PrivKey, to *common.Address, value *big.Int, data []byte) (*ethtypes.Transaction, error) {
	ctx := context.Background()
	sender := common.BytesToAddress(priKey.PubKey().Address().Bytes())

	chainId, err := cli.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	nonce, err := cli.NonceAt(ctx, sender, nil)
	if err != nil {
		return nil, err
	}
	head, err := cli.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	var gasTipCap, gasFeeCap, gasPrice *big.Int
	if head.BaseFee != nil {
		tip, err := cli.SuggestGasTipCap(ctx)
		if err != nil {
			return nil, err
		}
		gasTipCap = tip
		gasFeeCap = new(big.Int).Add(tip, new(big.Int).Mul(head.BaseFee, big.NewInt(2)))
		if gasFeeCap.Cmp(gasTipCap) < 0 {
			return nil, fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", gasFeeCap, gasTipCap)
		}
	} else {
		gasPrice, err = cli.SuggestGasPrice(ctx)
		if err != nil {
			return nil, err
		}
	}

	msg := ethereum.CallMsg{From: sender, To: to, GasPrice: gasPrice, GasTipCap: gasTipCap, GasFeeCap: gasFeeCap, Value: value, Data: data}
	gasLimit, err := cli.EstimateGas(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas needed: %v", err)
	}
	gasLimit = gasLimit * 130 / 100
	if value == nil {
		value = big.NewInt(0)
	}

	var rawTx *ethtypes.Transaction
	if gasFeeCap == nil {
		baseTx := &ethtypes.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      gasLimit,
			To:       to,
			Value:    value,
			Data:     data,
		}
		rawTx = ethtypes.NewTx(baseTx)
	} else {
		baseTx := &ethtypes.DynamicFeeTx{
			ChainID:   chainId,
			Nonce:     nonce,
			GasFeeCap: gasFeeCap,
			GasTipCap: gasTipCap,
			Gas:       gasLimit,
			To:        to,
			Value:     value,
			Data:      data,
		}
		rawTx = ethtypes.NewTx(baseTx)
	}
	signer := ethtypes.NewLondonSigner(chainId)
	signature, err := priKey.Sign(signer.Hash(rawTx).Bytes())
	if err != nil {
		return nil, err
	}
	return rawTx.WithSignature(signer, signature)
}
func broadcastTx(t *testing.T, ctx context.Context, txClient tx.ServiceClient, data *tx.TxRaw) string {
	broadcastData := mustProtoMarshal(t, data)
	broadcastTxResponse, err := txClient.BroadcastTx(ctx, &tx.BroadcastTxRequest{
		TxBytes: broadcastData,
		Mode:    tx.BroadcastMode_BROADCAST_MODE_BLOCK,
	})
	require.NoError(t, err)
	txResponse := broadcastTxResponse.TxResponse
	if txResponse.Code != 0 {
		t.Logf("broadcast tx error ===> %+v\n", txResponse.String())
	}
	require.Equal(t, txResponse.Code, uint32(0))
	t.Logf("broadcast tx success! height:%v txHash:%v gasUsed:%v\n", txResponse.Height, txResponse.TxHash, txResponse.GasUsed)
	return txResponse.TxHash
}
func sign(t *testing.T, priKey cryptotypes.PrivKey, signDoc *tx.SignDoc) []byte {
	signDataBytes := mustProtoMarshal(t, signDoc)
	signResultBytes, err := priKey.Sign(signDataBytes)
	require.NoError(t, err)

	return signResultBytes
}
func mustProtoMarshal(t *testing.T, pb proto.Message) (bytes []byte) {
	bytes, err := proto.Marshal(pb)
	require.NoError(t, err)
	return
}

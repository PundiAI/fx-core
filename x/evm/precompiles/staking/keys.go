package staking

import "github.com/ethereum/go-ethereum/common"

const (
	DelegateGas   = 100000
	UndelegateGas = 100000
	WithdrawGas   = 50000

	DelegationGas = 20000
)

const (
	DelegateMethodName   = "delegate"
	UndelegateMethodName = "undelegate"
	WithdrawMethodName   = "withdraw"

	DelegationMethodName = "delegation"
)

var StakingAddress = common.HexToAddress("0x0000000000000000000000000000000000000064")

const JsonABI = `
[
    {
        "type":"function",
        "name":"delegate",
        "inputs":[
            {
                "name":"validator",
                "type":"string"
            },
            {
                "name":"amount",
                "type":"uint256"
            }
        ],
        "outputs":[
            {
                "name":"shares",
                "type":"uint256"
            }
        ],
        "payable":true,
        "stateMutability":"payable"
    },
    {
        "type":"function",
        "name":"undelegate",
        "inputs":[
            {
                "name":"validator",
                "type":"string"
            },
            {
                "name":"shares",
                "type":"uint256"
            }
        ],
        "outputs":[
            {
                "name":"amount",
                "type":"uint256"
            },
            {
                "name":"endTime",
                "type":"uint256"
            }
        ],
        "payable":false,
        "stateMutability":"nonpayable"
    },
    {
        "type":"function",
        "name":"withdraw",
        "inputs":[
            {
                "name":"validator",
                "type":"string"
            }
        ],
        "outputs":[
            {
                "name":"reward",
                "type":"uint256"
            }
        ],
        "payable":false,
        "stateMutability":"nonpayable"
    },
    {
        "type":"function",
        "name":"delegation",
        "inputs":[
            {
                "name":"validator",
                "type":"string"
            },
            {
                "name":"delegator",
                "type":"address"
            }
        ],
        "outputs":[
            {
                "name":"delegate",
                "type":"uint256"
            }
        ],
        "payable":false,
        "stateMutability":"nonpayable"
    }
]`

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

const GrantPrivilegeSignaturePrefix = "GrantPrivilege:"

type CProcess []byte

var (
	ProcessStart CProcess = []byte{0x1}
	ProcessEnd   CProcess = []byte{0x2}
)

var DisablePKBytes = [33]byte{
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff,
}

var (
	AllowanceKey         = []byte{0x90}
	ValidatorOperatorKey = []byte{0x91}
	ConsensusPubKey      = []byte{0x92}
	ConsensusProcessKey  = []byte{0x93}
)

func GetAllowanceKey(valAddr sdk.ValAddress, owner, spender sdk.AccAddress) []byte {
	// key is of the form AllowanceKey || valAddrLen (1 byte) || valAddr || ownerAddrLen (1 byte) || ownerAddr || spenderAddrLen (1 byte) || spenderAddr
	offset := len(AllowanceKey)
	key := make([]byte, offset+3+len(valAddr)+len(owner)+len(spender))
	copy(key[0:offset], AllowanceKey)
	key[offset] = byte(len(valAddr))
	copy(key[offset+1:offset+1+len(valAddr)], valAddr.Bytes())
	key[offset+1+len(valAddr)] = byte(len(owner))
	copy(key[offset+2+len(valAddr):offset+2+len(valAddr)+len(owner)], owner.Bytes())
	key[offset+2+len(valAddr)+len(owner)] = byte(len(spender))
	copy(key[offset+3+len(valAddr)+len(owner):], spender.Bytes())

	return key
}

func GetValidatorOperatorKey(addr sdk.ValAddress) []byte {
	return append(ValidatorOperatorKey, addr...)
}

func GetConsensusPubKey(addr sdk.ValAddress) []byte {
	return append(ConsensusPubKey, addr...)
}

func GetConsensusProcessKey(process CProcess, addr sdk.ValAddress) []byte {
	return append(ConsensusProcessKey, append(process, addr...)...)
}

func AddressFromConsensusPubKey(key []byte) []byte {
	kv.AssertKeyAtLeastLength(key, 2)
	return key[1:] // remove prefix bytes
}

func AddressFromConsensusProcessKey(key []byte) []byte {
	kv.AssertKeyAtLeastLength(key, 3)
	return key[2:] // remove prefix bytes and process bytes
}

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fx/gravity/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Deprecated after upgrade v3
type Params struct {
	GravityId                      string                                 `protobuf:"bytes,1,opt,name=gravity_id,json=gravityId,proto3" json:"gravity_id,omitempty"`
	ContractSourceHash             string                                 `protobuf:"bytes,2,opt,name=contract_source_hash,json=contractSourceHash,proto3" json:"contract_source_hash,omitempty"`
	BridgeEthAddress               string                                 `protobuf:"bytes,4,opt,name=bridge_eth_address,json=bridgeEthAddress,proto3" json:"bridge_eth_address,omitempty"`
	BridgeChainId                  uint64                                 `protobuf:"varint,5,opt,name=bridge_chain_id,json=bridgeChainId,proto3" json:"bridge_chain_id,omitempty"`
	SignedValsetsWindow            uint64                                 `protobuf:"varint,6,opt,name=signed_valsets_window,json=signedValsetsWindow,proto3" json:"signed_valsets_window,omitempty"`
	SignedBatchesWindow            uint64                                 `protobuf:"varint,7,opt,name=signed_batches_window,json=signedBatchesWindow,proto3" json:"signed_batches_window,omitempty"`
	SignedClaimsWindow             uint64                                 `protobuf:"varint,8,opt,name=signed_claims_window,json=signedClaimsWindow,proto3" json:"signed_claims_window,omitempty"`
	TargetBatchTimeout             uint64                                 `protobuf:"varint,10,opt,name=target_batch_timeout,json=targetBatchTimeout,proto3" json:"target_batch_timeout,omitempty"`
	AverageBlockTime               uint64                                 `protobuf:"varint,11,opt,name=average_block_time,json=averageBlockTime,proto3" json:"average_block_time,omitempty"`
	AverageEthBlockTime            uint64                                 `protobuf:"varint,12,opt,name=average_eth_block_time,json=averageEthBlockTime,proto3" json:"average_eth_block_time,omitempty"`
	SlashFractionValset            github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,13,opt,name=slash_fraction_valset,json=slashFractionValset,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_valset"`
	SlashFractionBatch             github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,14,opt,name=slash_fraction_batch,json=slashFractionBatch,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_batch"`
	SlashFractionClaim             github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,15,opt,name=slash_fraction_claim,json=slashFractionClaim,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_claim"`
	SlashFractionConflictingClaim  github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,16,opt,name=slash_fraction_conflicting_claim,json=slashFractionConflictingClaim,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_conflicting_claim"`
	UnbondSlashingValsetsWindow    uint64                                 `protobuf:"varint,17,opt,name=unbond_slashing_valsets_window,json=unbondSlashingValsetsWindow,proto3" json:"unbond_slashing_valsets_window,omitempty"`
	IbcTransferTimeoutHeight       uint64                                 `protobuf:"varint,18,opt,name=ibc_transfer_timeout_height,json=ibcTransferTimeoutHeight,proto3" json:"ibc_transfer_timeout_height,omitempty"`
	ValsetUpdatePowerChangePercent github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,19,opt,name=valset_update_power_change_percent,json=valsetUpdatePowerChangePercent,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"valset_update_power_change_percent"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f11fb942b548d13, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetGravityId() string {
	if m != nil {
		return m.GravityId
	}
	return ""
}

func (m *Params) GetContractSourceHash() string {
	if m != nil {
		return m.ContractSourceHash
	}
	return ""
}

func (m *Params) GetBridgeEthAddress() string {
	if m != nil {
		return m.BridgeEthAddress
	}
	return ""
}

func (m *Params) GetBridgeChainId() uint64 {
	if m != nil {
		return m.BridgeChainId
	}
	return 0
}

func (m *Params) GetSignedValsetsWindow() uint64 {
	if m != nil {
		return m.SignedValsetsWindow
	}
	return 0
}

func (m *Params) GetSignedBatchesWindow() uint64 {
	if m != nil {
		return m.SignedBatchesWindow
	}
	return 0
}

func (m *Params) GetSignedClaimsWindow() uint64 {
	if m != nil {
		return m.SignedClaimsWindow
	}
	return 0
}

func (m *Params) GetTargetBatchTimeout() uint64 {
	if m != nil {
		return m.TargetBatchTimeout
	}
	return 0
}

func (m *Params) GetAverageBlockTime() uint64 {
	if m != nil {
		return m.AverageBlockTime
	}
	return 0
}

func (m *Params) GetAverageEthBlockTime() uint64 {
	if m != nil {
		return m.AverageEthBlockTime
	}
	return 0
}

func (m *Params) GetUnbondSlashingValsetsWindow() uint64 {
	if m != nil {
		return m.UnbondSlashingValsetsWindow
	}
	return 0
}

func (m *Params) GetIbcTransferTimeoutHeight() uint64 {
	if m != nil {
		return m.IbcTransferTimeoutHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "fx.gravity.v1.Params")
}

func init() { proto.RegisterFile("fx/gravity/v1/genesis.proto", fileDescriptor_1f11fb942b548d13) }

var fileDescriptor_1f11fb942b548d13 = []byte{
	// 651 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0x4d, 0x4f, 0x1b, 0x39,
	0x18, 0xc7, 0x33, 0x2b, 0x96, 0x5d, 0xbc, 0xb0, 0x50, 0x13, 0xd0, 0x14, 0xc4, 0x10, 0x71, 0x40,
	0x1c, 0x20, 0x43, 0xca, 0xad, 0x52, 0x0f, 0x4d, 0x0a, 0x82, 0x1b, 0x02, 0xda, 0x4a, 0xbd, 0xb8,
	0x1e, 0x8f, 0x67, 0xc6, 0x22, 0xb1, 0x23, 0xdb, 0x79, 0xa1, 0xa7, 0x5e, 0x7a, 0xef, 0xc7, 0xe2,
	0xc8, 0xb1, 0xaa, 0x2a, 0x54, 0xc1, 0x17, 0xa9, 0xe6, 0xb1, 0x21, 0x2f, 0xea, 0x09, 0xf5, 0x94,
	0xc8, 0xbf, 0xff, 0x4b, 0xe6, 0xf1, 0x93, 0x41, 0xeb, 0xd9, 0x30, 0xce, 0x35, 0xed, 0x0b, 0x7b,
	0x15, 0xf7, 0x1b, 0x71, 0xce, 0x25, 0x37, 0xc2, 0xd4, 0xbb, 0x5a, 0x59, 0x85, 0x17, 0xb2, 0x61,
	0xdd, 0xc3, 0x7a, 0xbf, 0xb1, 0x56, 0xcd, 0x55, 0xae, 0x80, 0xc4, 0xe5, 0x37, 0x27, 0x5a, 0x7b,
	0x3e, 0x99, 0x60, 0xaf, 0xba, 0xdc, 0xfb, 0xd7, 0x56, 0xa7, 0xd0, 0xd0, 0x9f, 0x47, 0x4c, 0x99,
	0x8e, 0x32, 0x71, 0x42, 0x0d, 0x8f, 0xfb, 0x8d, 0x84, 0x5b, 0xda, 0x88, 0x99, 0x12, 0xd2, 0xf1,
	0xad, 0x2f, 0x73, 0x68, 0xf6, 0x94, 0x6a, 0xda, 0x31, 0x78, 0x03, 0x21, 0x9f, 0x40, 0x44, 0x1a,
	0x06, 0xb5, 0x60, 0x67, 0xee, 0x6c, 0xce, 0x9f, 0x9c, 0xa4, 0x78, 0x1f, 0x55, 0x99, 0x92, 0x56,
	0x53, 0x66, 0x89, 0x51, 0x3d, 0xcd, 0x38, 0x29, 0xa8, 0x29, 0xc2, 0xbf, 0x40, 0x88, 0x1f, 0xd8,
	0x39, 0xa0, 0x63, 0x6a, 0x0a, 0xbc, 0x8b, 0x70, 0xa2, 0x45, 0x9a, 0x73, 0xc2, 0x6d, 0x41, 0x68,
	0x9a, 0x6a, 0x6e, 0x4c, 0x38, 0x03, 0xfa, 0x25, 0x47, 0x0e, 0x6d, 0xf1, 0xda, 0x9d, 0xe3, 0x6d,
	0xb4, 0xe8, 0xd5, 0xac, 0xa0, 0x42, 0x96, 0xbf, 0xe1, 0xef, 0x5a, 0xb0, 0x33, 0x73, 0xb6, 0xe0,
	0x8e, 0x5b, 0xe5, 0xe9, 0x49, 0x8a, 0x5f, 0xa0, 0x15, 0x23, 0x72, 0xc9, 0x53, 0xd2, 0xa7, 0x6d,
	0xc3, 0xad, 0x21, 0x03, 0x21, 0x53, 0x35, 0x08, 0x67, 0x41, 0xbd, 0xec, 0xe0, 0x3b, 0xc7, 0xde,
	0x03, 0x1a, 0xf3, 0x24, 0xd4, 0xb2, 0x82, 0x3f, 0x7a, 0xfe, 0x19, 0xf7, 0x34, 0x1d, 0xf3, 0x9e,
	0x7d, 0x54, 0xf5, 0x1e, 0xd6, 0xa6, 0xa2, 0xf3, 0x68, 0xf9, 0x17, 0x2c, 0xd8, 0xb1, 0x16, 0xa0,
	0x91, 0xc3, 0x52, 0x9d, 0x73, 0xeb, 0x5a, 0x88, 0x15, 0x1d, 0xae, 0x7a, 0x36, 0x44, 0xce, 0xe1,
	0x18, 0x94, 0x5c, 0x38, 0x52, 0x4e, 0x88, 0xf6, 0xb9, 0xa6, 0x39, 0x27, 0x49, 0x5b, 0xb1, 0x4b,
	0xb0, 0x84, 0xff, 0x81, 0x7e, 0xc9, 0x93, 0x66, 0x09, 0x4a, 0x03, 0x3e, 0x40, 0xab, 0x0f, 0xea,
	0x72, 0xa0, 0x63, 0x8e, 0x79, 0xf7, 0x18, 0x9e, 0x1e, 0xda, 0x62, 0x64, 0x4a, 0xd0, 0x8a, 0x69,
	0x53, 0x53, 0x90, 0xac, 0xbc, 0x1d, 0xa1, 0xa4, 0x1f, 0x5b, 0xb8, 0x50, 0x0b, 0x76, 0xe6, 0x9b,
	0xf5, 0xeb, 0xdb, 0xcd, 0xca, 0xf7, 0xdb, 0xcd, 0xed, 0x5c, 0xd8, 0xa2, 0x97, 0xd4, 0x99, 0xea,
	0xc4, 0x7e, 0x65, 0xdc, 0xc7, 0x9e, 0x49, 0x2f, 0xfd, 0xa6, 0xbd, 0xe1, 0xec, 0x6c, 0x19, 0xc2,
	0x8e, 0x7c, 0x96, 0x9b, 0x32, 0xfe, 0x88, 0xaa, 0x53, 0x1d, 0x30, 0x80, 0xf0, 0xff, 0x27, 0x55,
	0xe0, 0x89, 0x0a, 0x98, 0xd7, 0x6f, 0x1a, 0xe0, 0x52, 0xc2, 0xc5, 0x3f, 0xd0, 0x00, 0x77, 0x88,
	0x07, 0xa8, 0x36, 0xdd, 0xa0, 0x64, 0xd6, 0x16, 0xcc, 0x0a, 0x99, 0xfb, 0xb6, 0xa5, 0x27, 0xb5,
	0x6d, 0x4c, 0xb6, 0x8d, 0x52, 0x5d, 0x71, 0x0b, 0x45, 0x3d, 0x99, 0x28, 0x99, 0x12, 0xd0, 0x95,
	0x6d, 0x53, 0x8b, 0xfd, 0x0c, 0x6e, 0x77, 0xdd, 0xa9, 0xce, 0xbd, 0x68, 0x72, 0xc1, 0x5f, 0xa1,
	0x75, 0x91, 0x30, 0x62, 0x35, 0x95, 0x26, 0xe3, 0xfa, 0x61, 0xf5, 0x48, 0xc1, 0x45, 0x5e, 0xd8,
	0x10, 0x43, 0x42, 0x28, 0x12, 0x76, 0xe1, 0x15, 0x7e, 0x03, 0x8f, 0x81, 0xe3, 0x4f, 0x68, 0xcb,
	0x75, 0x92, 0x5e, 0x37, 0xa5, 0x96, 0x93, 0xae, 0x1a, 0x70, 0x5d, 0xfe, 0x11, 0x65, 0xce, 0x49,
	0x97, 0x6b, 0xc6, 0xa5, 0x0d, 0x97, 0x9f, 0xf4, 0xf8, 0x91, 0x4b, 0x7e, 0x0b, 0xc1, 0xa7, 0x65,
	0x6e, 0x0b, 0x62, 0x4f, 0x5d, 0xea, 0xcb, 0x99, 0xcf, 0x3f, 0x6a, 0x95, 0xe6, 0xd1, 0xf5, 0x5d,
	0x14, 0xdc, 0xdc, 0x45, 0xc1, 0xcf, 0xbb, 0x28, 0xf8, 0x7a, 0x1f, 0x55, 0x6e, 0xee, 0xa3, 0xca,
	0xb7, 0xfb, 0xa8, 0xf2, 0x61, 0x77, 0xac, 0x27, 0xeb, 0x49, 0x18, 0xe2, 0x30, 0xce, 0x86, 0x7b,
	0x4c, 0x69, 0x1e, 0x8f, 0xde, 0x7a, 0xd0, 0x98, 0xcc, 0xc2, 0x6b, 0xed, 0xe0, 0x57, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x5d, 0xb6, 0x9b, 0x4b, 0x6d, 0x05, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.ValsetUpdatePowerChangePercent.Size()
		i -= size
		if _, err := m.ValsetUpdatePowerChangePercent.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1
	i--
	dAtA[i] = 0x9a
	if m.IbcTransferTimeoutHeight != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.IbcTransferTimeoutHeight))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x90
	}
	if m.UnbondSlashingValsetsWindow != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.UnbondSlashingValsetsWindow))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x88
	}
	{
		size := m.SlashFractionConflictingClaim.Size()
		i -= size
		if _, err := m.SlashFractionConflictingClaim.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1
	i--
	dAtA[i] = 0x82
	{
		size := m.SlashFractionClaim.Size()
		i -= size
		if _, err := m.SlashFractionClaim.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x7a
	{
		size := m.SlashFractionBatch.Size()
		i -= size
		if _, err := m.SlashFractionBatch.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x72
	{
		size := m.SlashFractionValset.Size()
		i -= size
		if _, err := m.SlashFractionValset.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x6a
	if m.AverageEthBlockTime != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.AverageEthBlockTime))
		i--
		dAtA[i] = 0x60
	}
	if m.AverageBlockTime != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.AverageBlockTime))
		i--
		dAtA[i] = 0x58
	}
	if m.TargetBatchTimeout != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.TargetBatchTimeout))
		i--
		dAtA[i] = 0x50
	}
	if m.SignedClaimsWindow != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.SignedClaimsWindow))
		i--
		dAtA[i] = 0x40
	}
	if m.SignedBatchesWindow != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.SignedBatchesWindow))
		i--
		dAtA[i] = 0x38
	}
	if m.SignedValsetsWindow != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.SignedValsetsWindow))
		i--
		dAtA[i] = 0x30
	}
	if m.BridgeChainId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.BridgeChainId))
		i--
		dAtA[i] = 0x28
	}
	if len(m.BridgeEthAddress) > 0 {
		i -= len(m.BridgeEthAddress)
		copy(dAtA[i:], m.BridgeEthAddress)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.BridgeEthAddress)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.ContractSourceHash) > 0 {
		i -= len(m.ContractSourceHash)
		copy(dAtA[i:], m.ContractSourceHash)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.ContractSourceHash)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.GravityId) > 0 {
		i -= len(m.GravityId)
		copy(dAtA[i:], m.GravityId)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.GravityId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.GravityId)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.ContractSourceHash)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.BridgeEthAddress)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.BridgeChainId != 0 {
		n += 1 + sovGenesis(uint64(m.BridgeChainId))
	}
	if m.SignedValsetsWindow != 0 {
		n += 1 + sovGenesis(uint64(m.SignedValsetsWindow))
	}
	if m.SignedBatchesWindow != 0 {
		n += 1 + sovGenesis(uint64(m.SignedBatchesWindow))
	}
	if m.SignedClaimsWindow != 0 {
		n += 1 + sovGenesis(uint64(m.SignedClaimsWindow))
	}
	if m.TargetBatchTimeout != 0 {
		n += 1 + sovGenesis(uint64(m.TargetBatchTimeout))
	}
	if m.AverageBlockTime != 0 {
		n += 1 + sovGenesis(uint64(m.AverageBlockTime))
	}
	if m.AverageEthBlockTime != 0 {
		n += 1 + sovGenesis(uint64(m.AverageEthBlockTime))
	}
	l = m.SlashFractionValset.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.SlashFractionBatch.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.SlashFractionClaim.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.SlashFractionConflictingClaim.Size()
	n += 2 + l + sovGenesis(uint64(l))
	if m.UnbondSlashingValsetsWindow != 0 {
		n += 2 + sovGenesis(uint64(m.UnbondSlashingValsetsWindow))
	}
	if m.IbcTransferTimeoutHeight != 0 {
		n += 2 + sovGenesis(uint64(m.IbcTransferTimeoutHeight))
	}
	l = m.ValsetUpdatePowerChangePercent.Size()
	n += 2 + l + sovGenesis(uint64(l))
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GravityId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GravityId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractSourceHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContractSourceHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BridgeEthAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BridgeEthAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BridgeChainId", wireType)
			}
			m.BridgeChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BridgeChainId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignedValsetsWindow", wireType)
			}
			m.SignedValsetsWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SignedValsetsWindow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignedBatchesWindow", wireType)
			}
			m.SignedBatchesWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SignedBatchesWindow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignedClaimsWindow", wireType)
			}
			m.SignedClaimsWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SignedClaimsWindow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TargetBatchTimeout", wireType)
			}
			m.TargetBatchTimeout = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TargetBatchTimeout |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AverageBlockTime", wireType)
			}
			m.AverageBlockTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AverageBlockTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AverageEthBlockTime", wireType)
			}
			m.AverageEthBlockTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AverageEthBlockTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFractionValset", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SlashFractionValset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFractionBatch", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SlashFractionBatch.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 15:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFractionClaim", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SlashFractionClaim.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 16:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFractionConflictingClaim", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SlashFractionConflictingClaim.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 17:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondSlashingValsetsWindow", wireType)
			}
			m.UnbondSlashingValsetsWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UnbondSlashingValsetsWindow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 18:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IbcTransferTimeoutHeight", wireType)
			}
			m.IbcTransferTimeoutHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.IbcTransferTimeoutHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 19:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValsetUpdatePowerChangePercent", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ValsetUpdatePowerChangePercent.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)

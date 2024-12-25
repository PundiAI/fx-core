// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fx/gov/v1/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/durationpb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type UpdateStore struct {
	Space    string `protobuf:"bytes,1,opt,name=space,proto3" json:"space,omitempty"`
	Key      string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	OldValue string `protobuf:"bytes,3,opt,name=old_value,json=oldValue,proto3" json:"old_value,omitempty"`
	Value    string `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *UpdateStore) Reset()      { *m = UpdateStore{} }
func (*UpdateStore) ProtoMessage() {}
func (*UpdateStore) Descriptor() ([]byte, []int) {
	return fileDescriptor_a8e5d06ed1291671, []int{0}
}
func (m *UpdateStore) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UpdateStore) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UpdateStore.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UpdateStore) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateStore.Merge(m, src)
}
func (m *UpdateStore) XXX_Size() int {
	return m.Size()
}
func (m *UpdateStore) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateStore.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateStore proto.InternalMessageInfo

func (m *UpdateStore) GetSpace() string {
	if m != nil {
		return m.Space
	}
	return ""
}

func (m *UpdateStore) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *UpdateStore) GetOldValue() string {
	if m != nil {
		return m.OldValue
	}
	return ""
}

func (m *UpdateStore) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type SwitchParams struct {
	DisablePrecompiles []string `protobuf:"bytes,1,rep,name=disable_precompiles,json=disablePrecompiles,proto3" json:"disable_precompiles,omitempty"`
	DisableMsgTypes    []string `protobuf:"bytes,2,rep,name=disable_msg_types,json=disableMsgTypes,proto3" json:"disable_msg_types,omitempty"`
}

func (m *SwitchParams) Reset()         { *m = SwitchParams{} }
func (m *SwitchParams) String() string { return proto.CompactTextString(m) }
func (*SwitchParams) ProtoMessage()    {}
func (*SwitchParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_a8e5d06ed1291671, []int{1}
}
func (m *SwitchParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SwitchParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SwitchParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SwitchParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwitchParams.Merge(m, src)
}
func (m *SwitchParams) XXX_Size() int {
	return m.Size()
}
func (m *SwitchParams) XXX_DiscardUnknown() {
	xxx_messageInfo_SwitchParams.DiscardUnknown(m)
}

var xxx_messageInfo_SwitchParams proto.InternalMessageInfo

func (m *SwitchParams) GetDisablePrecompiles() []string {
	if m != nil {
		return m.DisablePrecompiles
	}
	return nil
}

func (m *SwitchParams) GetDisableMsgTypes() []string {
	if m != nil {
		return m.DisableMsgTypes
	}
	return nil
}

type CustomParams struct {
	// For EGF parameters, what percentage of deposit is required to enter the
	DepositRatio string `protobuf:"bytes,1,opt,name=deposit_ratio,json=depositRatio,proto3" json:"deposit_ratio,omitempty"`
	// Duration of the voting period.
	VotingPeriod *time.Duration `protobuf:"bytes,2,opt,name=voting_period,json=votingPeriod,proto3,stdduration" json:"voting_period,omitempty"`
	// Minimum percentage of total stake needed to vote for a result to be
	// considered valid.
	Quorum string `protobuf:"bytes,3,opt,name=quorum,proto3" json:"quorum,omitempty"`
}

func (m *CustomParams) Reset()         { *m = CustomParams{} }
func (m *CustomParams) String() string { return proto.CompactTextString(m) }
func (*CustomParams) ProtoMessage()    {}
func (*CustomParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_a8e5d06ed1291671, []int{2}
}
func (m *CustomParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CustomParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CustomParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CustomParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CustomParams.Merge(m, src)
}
func (m *CustomParams) XXX_Size() int {
	return m.Size()
}
func (m *CustomParams) XXX_DiscardUnknown() {
	xxx_messageInfo_CustomParams.DiscardUnknown(m)
}

var xxx_messageInfo_CustomParams proto.InternalMessageInfo

func (m *CustomParams) GetDepositRatio() string {
	if m != nil {
		return m.DepositRatio
	}
	return ""
}

func (m *CustomParams) GetVotingPeriod() *time.Duration {
	if m != nil {
		return m.VotingPeriod
	}
	return nil
}

func (m *CustomParams) GetQuorum() string {
	if m != nil {
		return m.Quorum
	}
	return ""
}

func init() {
	proto.RegisterType((*UpdateStore)(nil), "fx.gov.v1.UpdateStore")
	proto.RegisterType((*SwitchParams)(nil), "fx.gov.v1.SwitchParams")
	proto.RegisterType((*CustomParams)(nil), "fx.gov.v1.CustomParams")
}

func init() { proto.RegisterFile("fx/gov/v1/params.proto", fileDescriptor_a8e5d06ed1291671) }

var fileDescriptor_a8e5d06ed1291671 = []byte{
	// 425 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x92, 0x3d, 0x6f, 0x13, 0x31,
	0x18, 0xc7, 0x73, 0x4d, 0xa8, 0x88, 0x9b, 0xf2, 0x62, 0x2a, 0x74, 0x2d, 0xd2, 0xb5, 0xca, 0x50,
	0x55, 0x48, 0x3d, 0xab, 0x74, 0x63, 0x42, 0x21, 0x2b, 0x52, 0x74, 0x05, 0x06, 0x96, 0x93, 0x73,
	0x76, 0x5c, 0xab, 0x77, 0x79, 0x8c, 0x5f, 0x8e, 0xf4, 0x5b, 0x30, 0x76, 0xe4, 0x2b, 0x20, 0xf1,
	0x21, 0x18, 0x2b, 0x26, 0x36, 0x50, 0xf2, 0x45, 0xd0, 0xd9, 0xae, 0x18, 0xba, 0xf9, 0xff, 0xf2,
	0xd8, 0xf7, 0xfc, 0x74, 0xe8, 0xf9, 0x62, 0x45, 0x04, 0xb4, 0xa4, 0x3d, 0x23, 0x8a, 0x6a, 0xda,
	0x98, 0x5c, 0x69, 0xb0, 0x80, 0x87, 0x8b, 0x55, 0x2e, 0xa0, 0xcd, 0xdb, 0xb3, 0x83, 0xfd, 0x0a,
	0x4c, 0x03, 0xa6, 0xf4, 0x01, 0x09, 0x22, 0xb4, 0x0e, 0xf6, 0x04, 0x08, 0x08, 0x7e, 0x77, 0x8a,
	0x6e, 0x26, 0x00, 0x44, 0xcd, 0x89, 0x57, 0x73, 0xb7, 0x20, 0xcc, 0x69, 0x6a, 0x25, 0x2c, 0x43,
	0x3e, 0x5e, 0xa2, 0x9d, 0x0f, 0x8a, 0x51, 0xcb, 0x2f, 0x2c, 0x68, 0x8e, 0xf7, 0xd0, 0x03, 0xa3,
	0x68, 0xc5, 0xd3, 0xe4, 0x28, 0x39, 0x19, 0x16, 0x41, 0xe0, 0x27, 0xa8, 0x7f, 0xc5, 0xaf, 0xd3,
	0x2d, 0xef, 0x75, 0x47, 0xfc, 0x02, 0x0d, 0xa1, 0x66, 0x65, 0x4b, 0x6b, 0xc7, 0xd3, 0xbe, 0xf7,
	0x1f, 0x42, 0xcd, 0x3e, 0x76, 0xba, 0xbb, 0x24, 0x04, 0x83, 0x70, 0x89, 0x17, 0xaf, 0x07, 0x37,
	0xdf, 0x0e, 0x7b, 0xe3, 0x2b, 0x34, 0xba, 0xf8, 0x22, 0x6d, 0x75, 0x39, 0xf3, 0x1b, 0x62, 0x82,
	0x9e, 0x31, 0x69, 0xe8, 0xbc, 0xe6, 0xa5, 0xd2, 0xbc, 0x82, 0x46, 0xc9, 0x9a, 0x9b, 0x34, 0x39,
	0xea, 0x9f, 0x0c, 0x0b, 0x1c, 0xa3, 0xd9, 0xff, 0x04, 0xbf, 0x44, 0x4f, 0xef, 0x06, 0x1a, 0x23,
	0x4a, 0x7b, 0xad, 0xb8, 0x49, 0xb7, 0x7c, 0xfd, 0x71, 0x0c, 0xde, 0x19, 0xf1, 0xbe, 0xb3, 0xc7,
	0xdf, 0x13, 0x34, 0x7a, 0xeb, 0x8c, 0x85, 0x26, 0xbe, 0x76, 0x8e, 0x76, 0x19, 0x57, 0x60, 0xa4,
	0x2d, 0x3d, 0x85, 0xb0, 0xe6, 0xe4, 0xd1, 0xaf, 0x1f, 0xa7, 0x28, 0xc2, 0x9c, 0xf2, 0xaa, 0x18,
	0xc5, 0x52, 0xd1, 0x75, 0xf0, 0x14, 0xed, 0xb6, 0x60, 0xe5, 0x52, 0x94, 0x8a, 0x6b, 0x09, 0xcc,
	0x73, 0xd8, 0x79, 0xb5, 0x9f, 0x07, 0xb4, 0xf9, 0x1d, 0xda, 0x7c, 0x1a, 0xd1, 0x4e, 0x06, 0x37,
	0x7f, 0x0e, 0x93, 0x62, 0x14, 0xa6, 0x66, 0x7e, 0x08, 0x1f, 0xa3, 0xed, 0xcf, 0x0e, 0xb4, 0x6b,
	0x02, 0xae, 0x7b, 0x6f, 0xc6, 0x74, 0xf2, 0xe6, 0xe7, 0x3a, 0x4b, 0x6e, 0xd7, 0x59, 0xf2, 0x77,
	0x9d, 0x25, 0x5f, 0x37, 0x59, 0xef, 0x76, 0x93, 0xf5, 0x7e, 0x6f, 0xb2, 0xde, 0xa7, 0x63, 0x21,
	0xed, 0xa5, 0x9b, 0xe7, 0x15, 0x34, 0x44, 0xb9, 0x25, 0x93, 0x54, 0x92, 0xc5, 0xea, 0xb4, 0x02,
	0xcd, 0x49, 0xf8, 0x71, 0x3c, 0x8c, 0xf9, 0xb6, 0xff, 0xa0, 0xf3, 0x7f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x6c, 0x57, 0xf2, 0x32, 0x4f, 0x02, 0x00, 0x00,
}

func (m *UpdateStore) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpdateStore) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UpdateStore) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Value) > 0 {
		i -= len(m.Value)
		copy(dAtA[i:], m.Value)
		i = encodeVarintParams(dAtA, i, uint64(len(m.Value)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.OldValue) > 0 {
		i -= len(m.OldValue)
		copy(dAtA[i:], m.OldValue)
		i = encodeVarintParams(dAtA, i, uint64(len(m.OldValue)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintParams(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Space) > 0 {
		i -= len(m.Space)
		copy(dAtA[i:], m.Space)
		i = encodeVarintParams(dAtA, i, uint64(len(m.Space)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SwitchParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SwitchParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SwitchParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.DisableMsgTypes) > 0 {
		for iNdEx := len(m.DisableMsgTypes) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.DisableMsgTypes[iNdEx])
			copy(dAtA[i:], m.DisableMsgTypes[iNdEx])
			i = encodeVarintParams(dAtA, i, uint64(len(m.DisableMsgTypes[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.DisablePrecompiles) > 0 {
		for iNdEx := len(m.DisablePrecompiles) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.DisablePrecompiles[iNdEx])
			copy(dAtA[i:], m.DisablePrecompiles[iNdEx])
			i = encodeVarintParams(dAtA, i, uint64(len(m.DisablePrecompiles[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *CustomParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CustomParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CustomParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Quorum) > 0 {
		i -= len(m.Quorum)
		copy(dAtA[i:], m.Quorum)
		i = encodeVarintParams(dAtA, i, uint64(len(m.Quorum)))
		i--
		dAtA[i] = 0x1a
	}
	if m.VotingPeriod != nil {
		n1, err1 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(*m.VotingPeriod, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(*m.VotingPeriod):])
		if err1 != nil {
			return 0, err1
		}
		i -= n1
		i = encodeVarintParams(dAtA, i, uint64(n1))
		i--
		dAtA[i] = 0x12
	}
	if len(m.DepositRatio) > 0 {
		i -= len(m.DepositRatio)
		copy(dAtA[i:], m.DepositRatio)
		i = encodeVarintParams(dAtA, i, uint64(len(m.DepositRatio)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *UpdateStore) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Space)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.OldValue)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	return n
}

func (m *SwitchParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.DisablePrecompiles) > 0 {
		for _, s := range m.DisablePrecompiles {
			l = len(s)
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if len(m.DisableMsgTypes) > 0 {
		for _, s := range m.DisableMsgTypes {
			l = len(s)
			n += 1 + l + sovParams(uint64(l))
		}
	}
	return n
}

func (m *CustomParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DepositRatio)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.VotingPeriod != nil {
		l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(*m.VotingPeriod)
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.Quorum)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *UpdateStore) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: UpdateStore: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UpdateStore: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Space", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Space = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OldValue", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OldValue = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func (m *SwitchParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: SwitchParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SwitchParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DisablePrecompiles", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DisablePrecompiles = append(m.DisablePrecompiles, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DisableMsgTypes", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DisableMsgTypes = append(m.DisableMsgTypes, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func (m *CustomParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: CustomParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CustomParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DepositRatio", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DepositRatio = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VotingPeriod", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.VotingPeriod == nil {
				m.VotingPeriod = new(time.Duration)
			}
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(m.VotingPeriod, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Quorum", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Quorum = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)

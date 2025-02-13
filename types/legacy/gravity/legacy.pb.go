// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fx/gravity/v1/legacy.proto

package gravity

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

// Deprecated: after upgrade v3
type MsgRequestBatch struct {
	Sender     string                `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Denom      string                `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
	MinimumFee cosmossdk_io_math.Int `protobuf:"bytes,3,opt,name=minimum_fee,json=minimumFee,proto3,customtype=cosmossdk.io/math.Int" json:"minimum_fee"`
	FeeReceive string                `protobuf:"bytes,4,opt,name=fee_receive,json=feeReceive,proto3" json:"fee_receive,omitempty"`
	BaseFee    cosmossdk_io_math.Int `protobuf:"bytes,5,opt,name=base_fee,json=baseFee,proto3,customtype=cosmossdk.io/math.Int" json:"base_fee"`
}

func (m *MsgRequestBatch) Reset()         { *m = MsgRequestBatch{} }
func (m *MsgRequestBatch) String() string { return proto.CompactTextString(m) }
func (*MsgRequestBatch) ProtoMessage()    {}
func (*MsgRequestBatch) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef133984f5014717, []int{0}
}
func (m *MsgRequestBatch) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRequestBatch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRequestBatch.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRequestBatch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRequestBatch.Merge(m, src)
}
func (m *MsgRequestBatch) XXX_Size() int {
	return m.Size()
}
func (m *MsgRequestBatch) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRequestBatch.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRequestBatch proto.InternalMessageInfo

func (m *MsgRequestBatch) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func (m *MsgRequestBatch) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *MsgRequestBatch) GetFeeReceive() string {
	if m != nil {
		return m.FeeReceive
	}
	return ""
}

// Deprecated: after upgrade v3
type MsgSetOrchestratorAddress struct {
	Validator    string `protobuf:"bytes,1,opt,name=validator,proto3" json:"validator,omitempty"`
	Orchestrator string `protobuf:"bytes,2,opt,name=orchestrator,proto3" json:"orchestrator,omitempty"`
	EthAddress   string `protobuf:"bytes,3,opt,name=eth_address,json=ethAddress,proto3" json:"eth_address,omitempty"`
}

func (m *MsgSetOrchestratorAddress) Reset()         { *m = MsgSetOrchestratorAddress{} }
func (m *MsgSetOrchestratorAddress) String() string { return proto.CompactTextString(m) }
func (*MsgSetOrchestratorAddress) ProtoMessage()    {}
func (*MsgSetOrchestratorAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef133984f5014717, []int{1}
}
func (m *MsgSetOrchestratorAddress) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetOrchestratorAddress) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetOrchestratorAddress.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetOrchestratorAddress) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetOrchestratorAddress.Merge(m, src)
}
func (m *MsgSetOrchestratorAddress) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetOrchestratorAddress) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetOrchestratorAddress.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetOrchestratorAddress proto.InternalMessageInfo

func (m *MsgSetOrchestratorAddress) GetValidator() string {
	if m != nil {
		return m.Validator
	}
	return ""
}

func (m *MsgSetOrchestratorAddress) GetOrchestrator() string {
	if m != nil {
		return m.Orchestrator
	}
	return ""
}

func (m *MsgSetOrchestratorAddress) GetEthAddress() string {
	if m != nil {
		return m.EthAddress
	}
	return ""
}

func init() {
	proto.RegisterType((*MsgRequestBatch)(nil), "fx.gravity.v1.MsgRequestBatch")
	proto.RegisterType((*MsgSetOrchestratorAddress)(nil), "fx.gravity.v1.MsgSetOrchestratorAddress")
}

func init() { proto.RegisterFile("fx/gravity/v1/legacy.proto", fileDescriptor_ef133984f5014717) }

var fileDescriptor_ef133984f5014717 = []byte{
	// 378 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0xcf, 0x4a, 0xeb, 0x40,
	0x14, 0xc6, 0x93, 0x7b, 0x6f, 0x7b, 0x6f, 0xa7, 0x57, 0x84, 0xa1, 0x4a, 0x2c, 0x9a, 0x4a, 0x57,
	0x6e, 0xcc, 0x10, 0xdd, 0x74, 0x25, 0xd8, 0x45, 0x41, 0xb0, 0x08, 0x71, 0xe7, 0xa6, 0x4c, 0x26,
	0x27, 0xc9, 0x60, 0x93, 0xa9, 0x33, 0x93, 0xd0, 0x6e, 0x7c, 0x06, 0x1f, 0xab, 0xcb, 0x2e, 0xa5,
	0x8b, 0x22, 0xed, 0x8b, 0x48, 0xfe, 0x88, 0xba, 0x73, 0x37, 0xe7, 0x77, 0xce, 0xf9, 0xe6, 0x7c,
	0x7c, 0xa8, 0x1b, 0xce, 0x49, 0x24, 0x69, 0xce, 0xf5, 0x82, 0xe4, 0x2e, 0x99, 0x42, 0x44, 0xd9,
	0xc2, 0x99, 0x49, 0xa1, 0x05, 0xde, 0x0b, 0xe7, 0x4e, 0xdd, 0x73, 0x72, 0xb7, 0x6b, 0x33, 0xa1,
	0x12, 0xa1, 0x88, 0x4f, 0x15, 0x90, 0xdc, 0xf5, 0x41, 0x53, 0x97, 0x30, 0xc1, 0xd3, 0x6a, 0xbc,
	0xdb, 0x89, 0x44, 0x24, 0xca, 0x27, 0x29, 0x5e, 0x15, 0xed, 0xaf, 0x4d, 0xb4, 0x3f, 0x56, 0x91,
	0x07, 0x4f, 0x19, 0x28, 0x3d, 0xa4, 0x9a, 0xc5, 0xf8, 0x10, 0x35, 0x15, 0xa4, 0x01, 0x48, 0xcb,
	0x3c, 0x35, 0xcf, 0x5a, 0x5e, 0x5d, 0xe1, 0x0e, 0x6a, 0x04, 0x90, 0x8a, 0xc4, 0xfa, 0x55, 0xe2,
	0xaa, 0xc0, 0x57, 0xa8, 0x9d, 0xf0, 0x94, 0x27, 0x59, 0x32, 0x09, 0x01, 0xac, 0xdf, 0x45, 0x6f,
	0x78, 0xb2, 0xdc, 0xf4, 0x8c, 0xf5, 0xa6, 0x77, 0x50, 0x1d, 0xa5, 0x82, 0x47, 0x87, 0x0b, 0x92,
	0x50, 0x1d, 0x3b, 0x37, 0xa9, 0xf6, 0x50, 0xbd, 0x31, 0x02, 0xc0, 0x3d, 0xd4, 0x0e, 0x01, 0x26,
	0x12, 0x18, 0xf0, 0x1c, 0xac, 0x3f, 0xa5, 0x36, 0x0a, 0x01, 0xbc, 0x8a, 0xe0, 0x01, 0xfa, 0x57,
	0x78, 0x2a, 0xd5, 0x1b, 0x3f, 0x51, 0xff, 0x5b, 0x8c, 0x8f, 0x00, 0xfa, 0xcf, 0xe8, 0x68, 0xac,
	0xa2, 0x7b, 0xd0, 0x77, 0x92, 0xc5, 0xa0, 0xb4, 0xa4, 0x5a, 0xc8, 0xeb, 0x20, 0x90, 0xa0, 0x14,
	0x3e, 0x46, 0xad, 0x9c, 0x4e, 0x79, 0x50, 0xb0, 0xda, 0xe8, 0x27, 0xc0, 0x7d, 0xf4, 0x5f, 0x7c,
	0x59, 0xaa, 0x2d, 0x7f, 0x63, 0xc5, 0xe5, 0xa0, 0xe3, 0x09, 0xad, 0x04, 0x2b, 0xe7, 0x1e, 0x02,
	0x1d, 0xd7, 0x5f, 0x0c, 0x6f, 0x97, 0x5b, 0xdb, 0x5c, 0x6d, 0x6d, 0xf3, 0x6d, 0x6b, 0x9b, 0x2f,
	0x3b, 0xdb, 0x58, 0xed, 0x6c, 0xe3, 0x75, 0x67, 0x1b, 0x0f, 0x17, 0x11, 0xd7, 0x71, 0xe6, 0x3b,
	0x4c, 0x24, 0x64, 0x96, 0xa5, 0x01, 0xa7, 0x9c, 0x84, 0xf3, 0x73, 0x26, 0x24, 0x90, 0x7c, 0x40,
	0xf4, 0x62, 0x06, 0xaa, 0x4e, 0xfb, 0x23, 0x7f, 0xbf, 0x59, 0x26, 0x76, 0xf9, 0x1e, 0x00, 0x00,
	0xff, 0xff, 0xa2, 0x12, 0x0c, 0xbe, 0x14, 0x02, 0x00, 0x00,
}

func (m *MsgRequestBatch) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRequestBatch) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRequestBatch) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.BaseFee.Size()
		i -= size
		if _, err := m.BaseFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLegacy(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if len(m.FeeReceive) > 0 {
		i -= len(m.FeeReceive)
		copy(dAtA[i:], m.FeeReceive)
		i = encodeVarintLegacy(dAtA, i, uint64(len(m.FeeReceive)))
		i--
		dAtA[i] = 0x22
	}
	{
		size := m.MinimumFee.Size()
		i -= size
		if _, err := m.MinimumFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLegacy(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintLegacy(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintLegacy(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgSetOrchestratorAddress) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetOrchestratorAddress) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetOrchestratorAddress) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.EthAddress) > 0 {
		i -= len(m.EthAddress)
		copy(dAtA[i:], m.EthAddress)
		i = encodeVarintLegacy(dAtA, i, uint64(len(m.EthAddress)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Orchestrator) > 0 {
		i -= len(m.Orchestrator)
		copy(dAtA[i:], m.Orchestrator)
		i = encodeVarintLegacy(dAtA, i, uint64(len(m.Orchestrator)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Validator) > 0 {
		i -= len(m.Validator)
		copy(dAtA[i:], m.Validator)
		i = encodeVarintLegacy(dAtA, i, uint64(len(m.Validator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintLegacy(dAtA []byte, offset int, v uint64) int {
	offset -= sovLegacy(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgRequestBatch) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovLegacy(uint64(l))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovLegacy(uint64(l))
	}
	l = m.MinimumFee.Size()
	n += 1 + l + sovLegacy(uint64(l))
	l = len(m.FeeReceive)
	if l > 0 {
		n += 1 + l + sovLegacy(uint64(l))
	}
	l = m.BaseFee.Size()
	n += 1 + l + sovLegacy(uint64(l))
	return n
}

func (m *MsgSetOrchestratorAddress) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Validator)
	if l > 0 {
		n += 1 + l + sovLegacy(uint64(l))
	}
	l = len(m.Orchestrator)
	if l > 0 {
		n += 1 + l + sovLegacy(uint64(l))
	}
	l = len(m.EthAddress)
	if l > 0 {
		n += 1 + l + sovLegacy(uint64(l))
	}
	return n
}

func sovLegacy(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozLegacy(x uint64) (n int) {
	return sovLegacy(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgRequestBatch) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLegacy
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
			return fmt.Errorf("proto: MsgRequestBatch: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRequestBatch: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLegacy
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
				return ErrInvalidLengthLegacy
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLegacy
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLegacy
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
				return ErrInvalidLengthLegacy
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLegacy
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinimumFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLegacy
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
				return ErrInvalidLengthLegacy
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLegacy
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinimumFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeReceive", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLegacy
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
				return ErrInvalidLengthLegacy
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLegacy
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FeeReceive = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLegacy
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
				return ErrInvalidLengthLegacy
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLegacy
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BaseFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLegacy(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLegacy
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
func (m *MsgSetOrchestratorAddress) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLegacy
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
			return fmt.Errorf("proto: MsgSetOrchestratorAddress: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetOrchestratorAddress: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLegacy
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
				return ErrInvalidLengthLegacy
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLegacy
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Orchestrator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLegacy
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
				return ErrInvalidLengthLegacy
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLegacy
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Orchestrator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EthAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLegacy
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
				return ErrInvalidLengthLegacy
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLegacy
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EthAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLegacy(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLegacy
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
func skipLegacy(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLegacy
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
					return 0, ErrIntOverflowLegacy
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
					return 0, ErrIntOverflowLegacy
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
				return 0, ErrInvalidLengthLegacy
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupLegacy
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthLegacy
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthLegacy        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLegacy          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupLegacy = fmt.Errorf("proto: unexpected end of group")
)

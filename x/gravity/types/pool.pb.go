// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fx/gravity/v1/pool.proto

package types

import (
	fmt "fmt"
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

// IDSet represents a set of IDs
type IDSet struct {
	Ids []uint64 `protobuf:"varint,1,rep,packed,name=ids,proto3" json:"ids,omitempty"`
}

func (m *IDSet) Reset()         { *m = IDSet{} }
func (m *IDSet) String() string { return proto.CompactTextString(m) }
func (*IDSet) ProtoMessage()    {}
func (*IDSet) Descriptor() ([]byte, []int) {
	return fileDescriptor_627395d2ed3607d7, []int{0}
}
func (m *IDSet) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IDSet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IDSet.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IDSet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IDSet.Merge(m, src)
}
func (m *IDSet) XXX_Size() int {
	return m.Size()
}
func (m *IDSet) XXX_DiscardUnknown() {
	xxx_messageInfo_IDSet.DiscardUnknown(m)
}

var xxx_messageInfo_IDSet proto.InternalMessageInfo

func (m *IDSet) GetIds() []uint64 {
	if m != nil {
		return m.Ids
	}
	return nil
}

type BatchFees struct {
	TokenContract string                                 `protobuf:"bytes,1,opt,name=token_contract,json=tokenContract,proto3" json:"token_contract,omitempty"`
	TotalFees     github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,2,opt,name=total_fees,json=totalFees,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"total_fees"`
	TotalTxs      uint64                                 `protobuf:"varint,3,opt,name=total_txs,json=totalTxs,proto3" json:"total_txs,omitempty"`
}

func (m *BatchFees) Reset()         { *m = BatchFees{} }
func (m *BatchFees) String() string { return proto.CompactTextString(m) }
func (*BatchFees) ProtoMessage()    {}
func (*BatchFees) Descriptor() ([]byte, []int) {
	return fileDescriptor_627395d2ed3607d7, []int{1}
}
func (m *BatchFees) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BatchFees) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BatchFees.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BatchFees) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BatchFees.Merge(m, src)
}
func (m *BatchFees) XXX_Size() int {
	return m.Size()
}
func (m *BatchFees) XXX_DiscardUnknown() {
	xxx_messageInfo_BatchFees.DiscardUnknown(m)
}

var xxx_messageInfo_BatchFees proto.InternalMessageInfo

func (m *BatchFees) GetTokenContract() string {
	if m != nil {
		return m.TokenContract
	}
	return ""
}

func (m *BatchFees) GetTotalTxs() uint64 {
	if m != nil {
		return m.TotalTxs
	}
	return 0
}

func init() {
	proto.RegisterType((*IDSet)(nil), "fx.gravity.v1.IDSet")
	proto.RegisterType((*BatchFees)(nil), "fx.gravity.v1.BatchFees")
}

func init() { proto.RegisterFile("fx/gravity/v1/pool.proto", fileDescriptor_627395d2ed3607d7) }

var fileDescriptor_627395d2ed3607d7 = []byte{
	// 288 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0x41, 0x4a, 0xc3, 0x40,
	0x18, 0x85, 0x33, 0xb6, 0x8a, 0x1d, 0xa8, 0x48, 0x70, 0x11, 0x15, 0xa6, 0xa5, 0xa0, 0x74, 0x61,
	0x67, 0x28, 0xde, 0xa0, 0x4a, 0xa1, 0x0b, 0x37, 0xd1, 0x95, 0x9b, 0x92, 0x4e, 0x27, 0x69, 0x68,
	0x9b, 0x3f, 0x64, 0xfe, 0x86, 0xe9, 0x2d, 0x3c, 0x81, 0xe7, 0xe9, 0xb2, 0x4b, 0x71, 0x51, 0x24,
	0xb9, 0x88, 0x64, 0x1a, 0xd0, 0xd5, 0xbc, 0x79, 0xef, 0xf1, 0x7f, 0xf0, 0xa8, 0x17, 0x1a, 0x11,
	0x65, 0x41, 0x1e, 0xe3, 0x56, 0xe4, 0x43, 0x91, 0x02, 0xac, 0x78, 0x9a, 0x01, 0x82, 0xdb, 0x0e,
	0x0d, 0xaf, 0x13, 0x9e, 0x0f, 0x6f, 0xae, 0x22, 0x88, 0xc0, 0x26, 0xa2, 0x52, 0xc7, 0x52, 0xef,
	0x9a, 0x9e, 0x4e, 0x9e, 0x5f, 0x15, 0xba, 0x97, 0xb4, 0x11, 0xcf, 0xb5, 0x47, 0xba, 0x8d, 0x7e,
	0xd3, 0xaf, 0x64, 0xef, 0x93, 0xd0, 0xd6, 0x28, 0x40, 0xb9, 0x18, 0x2b, 0xa5, 0xdd, 0x3b, 0x7a,
	0x81, 0xb0, 0x54, 0xc9, 0x54, 0x42, 0x82, 0x59, 0x20, 0xd1, 0x23, 0x5d, 0xd2, 0x6f, 0xf9, 0x6d,
	0xeb, 0x3e, 0xd5, 0xa6, 0xfb, 0x42, 0x29, 0x02, 0x06, 0xab, 0x69, 0xa8, 0x94, 0xf6, 0x4e, 0xaa,
	0xca, 0x88, 0xef, 0x0e, 0x1d, 0xe7, 0xfb, 0xd0, 0xb9, 0x8f, 0x62, 0x5c, 0x6c, 0x66, 0x5c, 0xc2,
	0x5a, 0x48, 0xd0, 0x6b, 0xd0, 0xf5, 0x33, 0xd0, 0xf3, 0xa5, 0xc0, 0x6d, 0xaa, 0x34, 0x9f, 0x24,
	0xe8, 0xb7, 0xec, 0x05, 0x4b, 0xbd, 0xa5, 0xc7, 0xcf, 0x14, 0x8d, 0xf6, 0x1a, 0x5d, 0xd2, 0x6f,
	0xfa, 0xe7, 0xd6, 0x78, 0x33, 0x7a, 0x34, 0xde, 0x15, 0x8c, 0xec, 0x0b, 0x46, 0x7e, 0x0a, 0x46,
	0x3e, 0x4a, 0xe6, 0xec, 0x4b, 0xe6, 0x7c, 0x95, 0xcc, 0x79, 0x7f, 0xf8, 0x47, 0x0a, 0x37, 0x89,
	0xc4, 0x18, 0x12, 0x23, 0x42, 0x33, 0x90, 0x90, 0x29, 0xf1, 0x37, 0x98, 0x65, 0xce, 0xce, 0xec,
	0x14, 0x8f, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x4a, 0xbc, 0xfc, 0xe2, 0x4b, 0x01, 0x00, 0x00,
}

func (m *IDSet) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IDSet) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IDSet) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Ids) > 0 {
		dAtA2 := make([]byte, len(m.Ids)*10)
		var j1 int
		for _, num := range m.Ids {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		i -= j1
		copy(dAtA[i:], dAtA2[:j1])
		i = encodeVarintPool(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *BatchFees) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BatchFees) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BatchFees) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.TotalTxs != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.TotalTxs))
		i--
		dAtA[i] = 0x18
	}
	{
		size := m.TotalFees.Size()
		i -= size
		if _, err := m.TotalFees.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.TokenContract) > 0 {
		i -= len(m.TokenContract)
		copy(dAtA[i:], m.TokenContract)
		i = encodeVarintPool(dAtA, i, uint64(len(m.TokenContract)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintPool(dAtA []byte, offset int, v uint64) int {
	offset -= sovPool(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *IDSet) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Ids) > 0 {
		l = 0
		for _, e := range m.Ids {
			l += sovPool(uint64(e))
		}
		n += 1 + sovPool(uint64(l)) + l
	}
	return n
}

func (m *BatchFees) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.TokenContract)
	if l > 0 {
		n += 1 + l + sovPool(uint64(l))
	}
	l = m.TotalFees.Size()
	n += 1 + l + sovPool(uint64(l))
	if m.TotalTxs != 0 {
		n += 1 + sovPool(uint64(m.TotalTxs))
	}
	return n
}

func sovPool(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPool(x uint64) (n int) {
	return sovPool(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *IDSet) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPool
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
			return fmt.Errorf("proto: IDSet: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IDSet: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPool
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Ids = append(m.Ids, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPool
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthPool
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthPool
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.Ids) == 0 {
					m.Ids = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowPool
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Ids = append(m.Ids, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Ids", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPool(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPool
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
func (m *BatchFees) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPool
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
			return fmt.Errorf("proto: BatchFees: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BatchFees: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenContract", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenContract = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalFees", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalFees.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalTxs", wireType)
			}
			m.TotalTxs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TotalTxs |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPool(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPool
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
func skipPool(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPool
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
					return 0, ErrIntOverflowPool
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
					return 0, ErrIntOverflowPool
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
				return 0, ErrInvalidLengthPool
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPool
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPool
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPool        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPool          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPool = fmt.Errorf("proto: unexpected end of group")
)
// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fx/crosschain/v1/genesis.proto

package types

import (
	fmt "fmt"
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

// GenesisState struct
type GenesisState struct {
	Params                  *Params                  `protobuf:"bytes,1,opt,name=params,proto3" json:"params,omitempty"`
	LastObservedBlockHeight *LastObservedBlockHeight `protobuf:"bytes,2,opt,name=last_observed_block_height,json=lastObservedBlockHeight,proto3" json:"last_observed_block_height,omitempty"`
	OracleSet               []*OracleSet             `protobuf:"bytes,3,rep,name=OracleSet,proto3" json:"OracleSet,omitempty"`
	Oracle                  []*Oracle                `protobuf:"bytes,4,rep,name=oracle,proto3" json:"oracle,omitempty"`
	UnbatchedTransfers      []*OutgoingTransferTx    `protobuf:"bytes,5,rep,name=unbatched_transfers,json=unbatchedTransfers,proto3" json:"unbatched_transfers,omitempty"`
	Batches                 []*OutgoingTxBatch       `protobuf:"bytes,6,rep,name=batches,proto3" json:"batches,omitempty"`
	BridgeToken             []*BridgeToken           `protobuf:"bytes,7,rep,name=bridge_token,json=bridgeToken,proto3" json:"bridge_token,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_55adf565ef43ce68, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() *Params {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *GenesisState) GetLastObservedBlockHeight() *LastObservedBlockHeight {
	if m != nil {
		return m.LastObservedBlockHeight
	}
	return nil
}

func (m *GenesisState) GetOracleSet() []*OracleSet {
	if m != nil {
		return m.OracleSet
	}
	return nil
}

func (m *GenesisState) GetOracle() []*Oracle {
	if m != nil {
		return m.Oracle
	}
	return nil
}

func (m *GenesisState) GetUnbatchedTransfers() []*OutgoingTransferTx {
	if m != nil {
		return m.UnbatchedTransfers
	}
	return nil
}

func (m *GenesisState) GetBatches() []*OutgoingTxBatch {
	if m != nil {
		return m.Batches
	}
	return nil
}

func (m *GenesisState) GetBridgeToken() []*BridgeToken {
	if m != nil {
		return m.BridgeToken
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "fx.gravity.crosschain.v1.GenesisState")
}

func init() { proto.RegisterFile("fx/crosschain/v1/genesis.proto", fileDescriptor_55adf565ef43ce68) }

var fileDescriptor_55adf565ef43ce68 = []byte{
	// 388 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x41, 0xab, 0xd3, 0x30,
	0x1c, 0xc0, 0x57, 0x37, 0x37, 0xcc, 0x76, 0x8a, 0x07, 0xe3, 0x90, 0x52, 0x14, 0x61, 0x82, 0xb6,
	0x6c, 0x5e, 0xbc, 0x5a, 0x0f, 0x4e, 0x10, 0x26, 0xdd, 0x4e, 0x82, 0x94, 0x34, 0x4b, 0xdb, 0xb0,
	0x2e, 0x19, 0x49, 0x5a, 0xba, 0x6f, 0xe1, 0x77, 0xf0, 0xcb, 0x78, 0xdc, 0xd1, 0xa3, 0x6c, 0x5f,
	0xe4, 0xd1, 0xac, 0xdb, 0x78, 0xec, 0xf5, 0xbd, 0x77, 0x4c, 0xfe, 0xbf, 0xdf, 0x2f, 0x04, 0xfe,
	0xc0, 0x8e, 0x4b, 0x8f, 0x48, 0xa1, 0x14, 0x49, 0x31, 0xe3, 0x5e, 0x31, 0xf6, 0x12, 0xca, 0xa9,
	0x62, 0xca, 0xdd, 0x48, 0xa1, 0x05, 0x44, 0x71, 0xe9, 0x26, 0x12, 0x17, 0x4c, 0x6f, 0xdd, 0x0b,
	0xe7, 0x16, 0xe3, 0xe1, 0xab, 0x2b, 0x53, 0x6f, 0x37, 0xb4, 0xf6, 0x86, 0x2f, 0xaf, 0xa7, 0xe5,
	0x71, 0xf4, 0xfa, 0x4f, 0x07, 0x0c, 0xbe, 0x1e, 0x1f, 0x99, 0x6b, 0xac, 0x29, 0xfc, 0x04, 0xba,
	0x1b, 0x2c, 0xf1, 0x5a, 0x21, 0xcb, 0xb1, 0x46, 0xfd, 0x89, 0xe3, 0x36, 0x3d, 0xea, 0xfe, 0x30,
	0x5c, 0x50, 0xf3, 0x90, 0x83, 0x61, 0x86, 0x95, 0x0e, 0x45, 0xa4, 0xa8, 0x2c, 0xe8, 0x32, 0x8c,
	0x32, 0x41, 0x56, 0x61, 0x4a, 0x59, 0x92, 0x6a, 0xf4, 0xc4, 0xd4, 0xc6, 0xcd, 0xb5, 0xef, 0x58,
	0xe9, 0x59, 0xad, 0xfa, 0x95, 0x39, 0x35, 0x62, 0xf0, 0x22, 0xbb, 0x7b, 0x00, 0x3f, 0x83, 0x67,
	0x33, 0x89, 0x49, 0x46, 0xe7, 0x54, 0xa3, 0xb6, 0xd3, 0x1e, 0xf5, 0x27, 0x6f, 0x9a, 0xf3, 0x67,
	0x34, 0xb8, 0x58, 0xd5, 0x67, 0x85, 0x39, 0xa0, 0x8e, 0xf1, 0x9d, 0x87, 0xfc, 0xa0, 0xe6, 0xe1,
	0x2f, 0xf0, 0x3c, 0xe7, 0x11, 0xd6, 0x24, 0xa5, 0xcb, 0x50, 0x4b, 0xcc, 0x55, 0x4c, 0xa5, 0x42,
	0x4f, 0x4d, 0xe6, 0xfd, 0x3d, 0x99, 0x5c, 0x27, 0x82, 0xf1, 0x64, 0x51, 0x2b, 0x8b, 0x32, 0x80,
	0xe7, 0xd0, 0xe9, 0x52, 0xc1, 0x2f, 0xa0, 0x77, 0xbc, 0x53, 0xa8, 0x6b, 0x92, 0xef, 0x1e, 0x91,
	0x2c, 0xfd, 0x4a, 0x09, 0x4e, 0x26, 0x9c, 0x82, 0x41, 0x24, 0xd9, 0x32, 0xa1, 0xa1, 0x16, 0x2b,
	0xca, 0x51, 0xcf, 0x94, 0xde, 0x36, 0x97, 0x7c, 0x43, 0x2f, 0x2a, 0x38, 0xe8, 0x47, 0x97, 0x83,
	0xff, 0xed, 0xef, 0xde, 0xb6, 0x76, 0x7b, 0xdb, 0xfa, 0xbf, 0xb7, 0xad, 0xdf, 0x07, 0xbb, 0xb5,
	0x3b, 0xd8, 0xad, 0x7f, 0x07, 0xbb, 0xf5, 0xd3, 0x4b, 0x98, 0x4e, 0xf3, 0xc8, 0x25, 0x62, 0xed,
	0xc5, 0x39, 0x27, 0x9a, 0x09, 0x5e, 0x7a, 0x71, 0xf9, 0x81, 0x08, 0x49, 0xbd, 0x5b, 0x6b, 0x67,
	0x36, 0x32, 0xea, 0x9a, 0xbd, 0xfb, 0x78, 0x13, 0x00, 0x00, 0xff, 0xff, 0x03, 0xbd, 0xd7, 0x52,
	0xec, 0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.BridgeToken) > 0 {
		for iNdEx := len(m.BridgeToken) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.BridgeToken[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.Batches) > 0 {
		for iNdEx := len(m.Batches) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Batches[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.UnbatchedTransfers) > 0 {
		for iNdEx := len(m.UnbatchedTransfers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.UnbatchedTransfers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Oracle) > 0 {
		for iNdEx := len(m.Oracle) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Oracle[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.OracleSet) > 0 {
		for iNdEx := len(m.OracleSet) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.OracleSet[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.LastObservedBlockHeight != nil {
		{
			size, err := m.LastObservedBlockHeight.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Params != nil {
		{
			size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
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
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Params != nil {
		l = m.Params.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.LastObservedBlockHeight != nil {
		l = m.LastObservedBlockHeight.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.OracleSet) > 0 {
		for _, e := range m.OracleSet {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Oracle) > 0 {
		for _, e := range m.Oracle {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.UnbatchedTransfers) > 0 {
		for _, e := range m.UnbatchedTransfers {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Batches) > 0 {
		for _, e := range m.Batches {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.BridgeToken) > 0 {
		for _, e := range m.BridgeToken {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Params == nil {
				m.Params = &Params{}
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastObservedBlockHeight", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.LastObservedBlockHeight == nil {
				m.LastObservedBlockHeight = &LastObservedBlockHeight{}
			}
			if err := m.LastObservedBlockHeight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OracleSet", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OracleSet = append(m.OracleSet, &OracleSet{})
			if err := m.OracleSet[len(m.OracleSet)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Oracle", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Oracle = append(m.Oracle, &Oracle{})
			if err := m.Oracle[len(m.Oracle)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbatchedTransfers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UnbatchedTransfers = append(m.UnbatchedTransfers, &OutgoingTransferTx{})
			if err := m.UnbatchedTransfers[len(m.UnbatchedTransfers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Batches", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Batches = append(m.Batches, &OutgoingTxBatch{})
			if err := m.Batches[len(m.Batches)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BridgeToken", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BridgeToken = append(m.BridgeToken, &BridgeToken{})
			if err := m.BridgeToken[len(m.BridgeToken)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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

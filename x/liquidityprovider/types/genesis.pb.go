// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: em/liquidityprovider/v1beta1/genesis.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
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

type GenesisState struct {
	Accounts []GenesisAcc `protobuf:"bytes,1,rep,name=accounts,proto3" json:"accounts" yaml:"accounts"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_38424f7a0f6956c3, []int{0}
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

func (m *GenesisState) GetAccounts() []GenesisAcc {
	if m != nil {
		return m.Accounts
	}
	return nil
}

type GenesisAcc struct {
	Account  string                                   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty" yaml:"address"`
	Mintable github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,2,rep,name=mintable,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"mintable" yaml:"mintable"`
}

func (m *GenesisAcc) Reset()         { *m = GenesisAcc{} }
func (m *GenesisAcc) String() string { return proto.CompactTextString(m) }
func (*GenesisAcc) ProtoMessage()    {}
func (*GenesisAcc) Descriptor() ([]byte, []int) {
	return fileDescriptor_38424f7a0f6956c3, []int{1}
}
func (m *GenesisAcc) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisAcc) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisAcc.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisAcc) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisAcc.Merge(m, src)
}
func (m *GenesisAcc) XXX_Size() int {
	return m.Size()
}
func (m *GenesisAcc) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisAcc.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisAcc proto.InternalMessageInfo

func (m *GenesisAcc) GetAccount() string {
	if m != nil {
		return m.Account
	}
	return ""
}

func (m *GenesisAcc) GetMintable() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Mintable
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "em.liquidityprovider.v1beta1.GenesisState")
	proto.RegisterType((*GenesisAcc)(nil), "em.liquidityprovider.v1beta1.GenesisAcc")
}

func init() {
	proto.RegisterFile("em/liquidityprovider/v1beta1/genesis.proto", fileDescriptor_38424f7a0f6956c3)
}

var fileDescriptor_38424f7a0f6956c3 = []byte{
	// 349 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0xc1, 0x4a, 0xeb, 0x40,
	0x18, 0x85, 0x93, 0x7b, 0xe1, 0xb6, 0x37, 0xf7, 0xa2, 0x10, 0x04, 0x6b, 0xd1, 0x44, 0xb2, 0x0a,
	0x42, 0x67, 0xa8, 0x82, 0x82, 0xbb, 0xa6, 0x0b, 0xb7, 0x52, 0x77, 0x82, 0x8b, 0xc9, 0xe4, 0x27,
	0x0e, 0x66, 0x32, 0x35, 0x33, 0x2d, 0xc6, 0xa7, 0xf0, 0x39, 0x7c, 0x0b, 0x77, 0x5d, 0x76, 0xe9,
	0x2a, 0x4a, 0xfa, 0x06, 0x7d, 0x02, 0x69, 0x26, 0x8d, 0x8a, 0xe2, 0x6a, 0x06, 0xe6, 0x9c, 0xef,
	0x9c, 0x7f, 0x7e, 0xeb, 0x00, 0x38, 0x4e, 0xd8, 0xed, 0x84, 0x45, 0x4c, 0xe5, 0xe3, 0x4c, 0x4c,
	0x59, 0x04, 0x19, 0x9e, 0xf6, 0x43, 0x50, 0xa4, 0x8f, 0x63, 0x48, 0x41, 0x32, 0x89, 0xc6, 0x99,
	0x50, 0xc2, 0xde, 0x05, 0x8e, 0xbe, 0x68, 0x51, 0xad, 0xed, 0x6e, 0xc5, 0x22, 0x16, 0x95, 0x10,
	0xaf, 0x6e, 0xda, 0xd3, 0x75, 0xa8, 0x90, 0x5c, 0x48, 0x1c, 0x12, 0x09, 0x0d, 0x96, 0x0a, 0x96,
	0xea, 0x77, 0x8f, 0x5b, 0xff, 0xcf, 0x74, 0xc8, 0x85, 0x22, 0x0a, 0xec, 0x2b, 0xab, 0x4d, 0x28,
	0x15, 0x93, 0x54, 0xc9, 0x8e, 0xb9, 0xff, 0xdb, 0xff, 0x77, 0xe8, 0xa3, 0x9f, 0x62, 0x51, 0xed,
	0x1e, 0x50, 0x1a, 0x6c, 0xcf, 0x0a, 0xd7, 0x58, 0x16, 0xee, 0x66, 0x4e, 0x78, 0x72, 0xea, 0xad,
	0x39, 0xde, 0xa8, 0x41, 0x7a, 0x4f, 0xa6, 0x65, 0xbd, 0x3b, 0xec, 0x13, 0xab, 0x45, 0xa2, 0x28,
	0x03, 0xb9, 0x0a, 0x33, 0xfd, 0xbf, 0xc1, 0x5e, 0x59, 0xb8, 0xad, 0x81, 0x56, 0x2f, 0x0b, 0x77,
	0xa3, 0x26, 0x69, 0x8d, 0x37, 0x5a, 0xab, 0xed, 0x7b, 0xab, 0xcd, 0x59, 0xaa, 0x48, 0x98, 0x40,
	0xe7, 0x57, 0x55, 0x73, 0x07, 0xe9, 0x49, 0xd1, 0x6a, 0xd2, 0xa6, 0xdd, 0x50, 0xb0, 0x34, 0x18,
	0x7e, 0xee, 0xb5, 0x36, 0x7a, 0x8f, 0x2f, 0xae, 0x1f, 0x33, 0x75, 0x3d, 0x09, 0x11, 0x15, 0x1c,
	0xd7, 0x3f, 0xa5, 0x8f, 0x9e, 0x8c, 0x6e, 0xb0, 0xca, 0xc7, 0x20, 0x2b, 0x86, 0x1c, 0x35, 0x79,
	0xc1, 0xf9, 0xac, 0x74, 0xcc, 0x79, 0xe9, 0x98, 0xaf, 0xa5, 0x63, 0x3e, 0x2c, 0x1c, 0x63, 0xbe,
	0x70, 0x8c, 0xe7, 0x85, 0x63, 0x5c, 0x1e, 0x7f, 0xa0, 0x41, 0x8f, 0x8b, 0x14, 0x72, 0x0c, 0xbc,
	0x97, 0x40, 0x14, 0x43, 0x86, 0xef, 0xbe, 0x59, 0x74, 0x95, 0x10, 0xfe, 0xa9, 0x76, 0x71, 0xf4,
	0x16, 0x00, 0x00, 0xff, 0xff, 0x28, 0x5d, 0xd0, 0xda, 0x0d, 0x02, 0x00, 0x00,
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
	if len(m.Accounts) > 0 {
		for iNdEx := len(m.Accounts) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Accounts[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *GenesisAcc) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisAcc) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisAcc) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Mintable) > 0 {
		for iNdEx := len(m.Mintable) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Mintable[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Account) > 0 {
		i -= len(m.Account)
		copy(dAtA[i:], m.Account)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Account)))
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
	if len(m.Accounts) > 0 {
		for _, e := range m.Accounts {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *GenesisAcc) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Account)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.Mintable) > 0 {
		for _, e := range m.Mintable {
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
				return fmt.Errorf("proto: wrong wireType = %d for field Accounts", wireType)
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
			m.Accounts = append(m.Accounts, GenesisAcc{})
			if err := m.Accounts[len(m.Accounts)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *GenesisAcc) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: GenesisAcc: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisAcc: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
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
			m.Account = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Mintable", wireType)
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
			m.Mintable = append(m.Mintable, types.Coin{})
			if err := m.Mintable[len(m.Mintable)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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

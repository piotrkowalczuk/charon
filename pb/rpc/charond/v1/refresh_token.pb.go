// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/piotrkowalczuk/charon/pb/rpc/charond/v1/refresh_token.proto

package charond // import "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"
import ntypes "github.com/piotrkowalczuk/ntypes"
import qtypes "github.com/piotrkowalczuk/qtypes"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RefreshToken struct {
	Token                string               `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Notes                *ntypes.String       `protobuf:"bytes,2,opt,name=notes,proto3" json:"notes,omitempty"`
	UserId               int64                `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Revoked              bool                 `protobuf:"varint,4,opt,name=revoked,proto3" json:"revoked,omitempty"`
	ExpireAt             *timestamp.Timestamp `protobuf:"bytes,5,opt,name=expire_at,json=expireAt,proto3" json:"expire_at,omitempty"`
	LastUsedAt           *timestamp.Timestamp `protobuf:"bytes,6,opt,name=last_used_at,json=lastUsedAt,proto3" json:"last_used_at,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	CreatedBy            *ntypes.Int64        `protobuf:"bytes,8,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	UpdatedAt            *timestamp.Timestamp `protobuf:"bytes,9,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	UpdatedBy            *ntypes.Int64        `protobuf:"bytes,10,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *RefreshToken) Reset()         { *m = RefreshToken{} }
func (m *RefreshToken) String() string { return proto.CompactTextString(m) }
func (*RefreshToken) ProtoMessage()    {}
func (*RefreshToken) Descriptor() ([]byte, []int) {
	return fileDescriptor_refresh_token_1e8e54ee3f1fef8c, []int{0}
}
func (m *RefreshToken) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RefreshToken.Unmarshal(m, b)
}
func (m *RefreshToken) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RefreshToken.Marshal(b, m, deterministic)
}
func (dst *RefreshToken) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RefreshToken.Merge(dst, src)
}
func (m *RefreshToken) XXX_Size() int {
	return xxx_messageInfo_RefreshToken.Size(m)
}
func (m *RefreshToken) XXX_DiscardUnknown() {
	xxx_messageInfo_RefreshToken.DiscardUnknown(m)
}

var xxx_messageInfo_RefreshToken proto.InternalMessageInfo

func (m *RefreshToken) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *RefreshToken) GetNotes() *ntypes.String {
	if m != nil {
		return m.Notes
	}
	return nil
}

func (m *RefreshToken) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *RefreshToken) GetRevoked() bool {
	if m != nil {
		return m.Revoked
	}
	return false
}

func (m *RefreshToken) GetExpireAt() *timestamp.Timestamp {
	if m != nil {
		return m.ExpireAt
	}
	return nil
}

func (m *RefreshToken) GetLastUsedAt() *timestamp.Timestamp {
	if m != nil {
		return m.LastUsedAt
	}
	return nil
}

func (m *RefreshToken) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *RefreshToken) GetCreatedBy() *ntypes.Int64 {
	if m != nil {
		return m.CreatedBy
	}
	return nil
}

func (m *RefreshToken) GetUpdatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

func (m *RefreshToken) GetUpdatedBy() *ntypes.Int64 {
	if m != nil {
		return m.UpdatedBy
	}
	return nil
}

type RefreshTokenQuery struct {
	UserId               *qtypes.Int64     `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Notes                *qtypes.String    `protobuf:"bytes,2,opt,name=notes,proto3" json:"notes,omitempty"`
	Revoked              *ntypes.Bool      `protobuf:"bytes,3,opt,name=revoked,proto3" json:"revoked,omitempty"`
	ExpireAt             *qtypes.Timestamp `protobuf:"bytes,4,opt,name=expire_at,json=expireAt,proto3" json:"expire_at,omitempty"`
	LastUsedAt           *qtypes.Timestamp `protobuf:"bytes,5,opt,name=last_used_at,json=lastUsedAt,proto3" json:"last_used_at,omitempty"`
	CreatedAt            *qtypes.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt            *qtypes.Timestamp `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *RefreshTokenQuery) Reset()         { *m = RefreshTokenQuery{} }
func (m *RefreshTokenQuery) String() string { return proto.CompactTextString(m) }
func (*RefreshTokenQuery) ProtoMessage()    {}
func (*RefreshTokenQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_refresh_token_1e8e54ee3f1fef8c, []int{1}
}
func (m *RefreshTokenQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RefreshTokenQuery.Unmarshal(m, b)
}
func (m *RefreshTokenQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RefreshTokenQuery.Marshal(b, m, deterministic)
}
func (dst *RefreshTokenQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RefreshTokenQuery.Merge(dst, src)
}
func (m *RefreshTokenQuery) XXX_Size() int {
	return xxx_messageInfo_RefreshTokenQuery.Size(m)
}
func (m *RefreshTokenQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_RefreshTokenQuery.DiscardUnknown(m)
}

var xxx_messageInfo_RefreshTokenQuery proto.InternalMessageInfo

func (m *RefreshTokenQuery) GetUserId() *qtypes.Int64 {
	if m != nil {
		return m.UserId
	}
	return nil
}

func (m *RefreshTokenQuery) GetNotes() *qtypes.String {
	if m != nil {
		return m.Notes
	}
	return nil
}

func (m *RefreshTokenQuery) GetRevoked() *ntypes.Bool {
	if m != nil {
		return m.Revoked
	}
	return nil
}

func (m *RefreshTokenQuery) GetExpireAt() *qtypes.Timestamp {
	if m != nil {
		return m.ExpireAt
	}
	return nil
}

func (m *RefreshTokenQuery) GetLastUsedAt() *qtypes.Timestamp {
	if m != nil {
		return m.LastUsedAt
	}
	return nil
}

func (m *RefreshTokenQuery) GetCreatedAt() *qtypes.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *RefreshTokenQuery) GetUpdatedAt() *qtypes.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

type CreateRefreshTokenRequest struct {
	Notes                *ntypes.String       `protobuf:"bytes,1,opt,name=notes,proto3" json:"notes,omitempty"`
	ExpireAt             *timestamp.Timestamp `protobuf:"bytes,2,opt,name=expire_at,json=expireAt,proto3" json:"expire_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CreateRefreshTokenRequest) Reset()         { *m = CreateRefreshTokenRequest{} }
func (m *CreateRefreshTokenRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRefreshTokenRequest) ProtoMessage()    {}
func (*CreateRefreshTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_refresh_token_1e8e54ee3f1fef8c, []int{2}
}
func (m *CreateRefreshTokenRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateRefreshTokenRequest.Unmarshal(m, b)
}
func (m *CreateRefreshTokenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateRefreshTokenRequest.Marshal(b, m, deterministic)
}
func (dst *CreateRefreshTokenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateRefreshTokenRequest.Merge(dst, src)
}
func (m *CreateRefreshTokenRequest) XXX_Size() int {
	return xxx_messageInfo_CreateRefreshTokenRequest.Size(m)
}
func (m *CreateRefreshTokenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateRefreshTokenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateRefreshTokenRequest proto.InternalMessageInfo

func (m *CreateRefreshTokenRequest) GetNotes() *ntypes.String {
	if m != nil {
		return m.Notes
	}
	return nil
}

func (m *CreateRefreshTokenRequest) GetExpireAt() *timestamp.Timestamp {
	if m != nil {
		return m.ExpireAt
	}
	return nil
}

type CreateRefreshTokenResponse struct {
	RefreshToken         *RefreshToken `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *CreateRefreshTokenResponse) Reset()         { *m = CreateRefreshTokenResponse{} }
func (m *CreateRefreshTokenResponse) String() string { return proto.CompactTextString(m) }
func (*CreateRefreshTokenResponse) ProtoMessage()    {}
func (*CreateRefreshTokenResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_refresh_token_1e8e54ee3f1fef8c, []int{3}
}
func (m *CreateRefreshTokenResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateRefreshTokenResponse.Unmarshal(m, b)
}
func (m *CreateRefreshTokenResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateRefreshTokenResponse.Marshal(b, m, deterministic)
}
func (dst *CreateRefreshTokenResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateRefreshTokenResponse.Merge(dst, src)
}
func (m *CreateRefreshTokenResponse) XXX_Size() int {
	return xxx_messageInfo_CreateRefreshTokenResponse.Size(m)
}
func (m *CreateRefreshTokenResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateRefreshTokenResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateRefreshTokenResponse proto.InternalMessageInfo

func (m *CreateRefreshTokenResponse) GetRefreshToken() *RefreshToken {
	if m != nil {
		return m.RefreshToken
	}
	return nil
}

type ListRefreshTokensRequest struct {
	Offset               *ntypes.Int64      `protobuf:"bytes,1,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                *ntypes.Int64      `protobuf:"bytes,2,opt,name=limit,proto3" json:"limit,omitempty"`
	OrderBy              []*Order           `protobuf:"bytes,3,rep,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	Query                *RefreshTokenQuery `protobuf:"bytes,11,opt,name=query,proto3" json:"query,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ListRefreshTokensRequest) Reset()         { *m = ListRefreshTokensRequest{} }
func (m *ListRefreshTokensRequest) String() string { return proto.CompactTextString(m) }
func (*ListRefreshTokensRequest) ProtoMessage()    {}
func (*ListRefreshTokensRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_refresh_token_1e8e54ee3f1fef8c, []int{4}
}
func (m *ListRefreshTokensRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRefreshTokensRequest.Unmarshal(m, b)
}
func (m *ListRefreshTokensRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRefreshTokensRequest.Marshal(b, m, deterministic)
}
func (dst *ListRefreshTokensRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRefreshTokensRequest.Merge(dst, src)
}
func (m *ListRefreshTokensRequest) XXX_Size() int {
	return xxx_messageInfo_ListRefreshTokensRequest.Size(m)
}
func (m *ListRefreshTokensRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRefreshTokensRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListRefreshTokensRequest proto.InternalMessageInfo

func (m *ListRefreshTokensRequest) GetOffset() *ntypes.Int64 {
	if m != nil {
		return m.Offset
	}
	return nil
}

func (m *ListRefreshTokensRequest) GetLimit() *ntypes.Int64 {
	if m != nil {
		return m.Limit
	}
	return nil
}

func (m *ListRefreshTokensRequest) GetOrderBy() []*Order {
	if m != nil {
		return m.OrderBy
	}
	return nil
}

func (m *ListRefreshTokensRequest) GetQuery() *RefreshTokenQuery {
	if m != nil {
		return m.Query
	}
	return nil
}

type ListRefreshTokensResponse struct {
	RefreshTokens        []*RefreshToken `protobuf:"bytes,1,rep,name=refresh_tokens,json=refreshTokens,proto3" json:"refresh_tokens,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ListRefreshTokensResponse) Reset()         { *m = ListRefreshTokensResponse{} }
func (m *ListRefreshTokensResponse) String() string { return proto.CompactTextString(m) }
func (*ListRefreshTokensResponse) ProtoMessage()    {}
func (*ListRefreshTokensResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_refresh_token_1e8e54ee3f1fef8c, []int{5}
}
func (m *ListRefreshTokensResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRefreshTokensResponse.Unmarshal(m, b)
}
func (m *ListRefreshTokensResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRefreshTokensResponse.Marshal(b, m, deterministic)
}
func (dst *ListRefreshTokensResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRefreshTokensResponse.Merge(dst, src)
}
func (m *ListRefreshTokensResponse) XXX_Size() int {
	return xxx_messageInfo_ListRefreshTokensResponse.Size(m)
}
func (m *ListRefreshTokensResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRefreshTokensResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListRefreshTokensResponse proto.InternalMessageInfo

func (m *ListRefreshTokensResponse) GetRefreshTokens() []*RefreshToken {
	if m != nil {
		return m.RefreshTokens
	}
	return nil
}

type RevokeRefreshTokenRequest struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	UserId               int64    `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RevokeRefreshTokenRequest) Reset()         { *m = RevokeRefreshTokenRequest{} }
func (m *RevokeRefreshTokenRequest) String() string { return proto.CompactTextString(m) }
func (*RevokeRefreshTokenRequest) ProtoMessage()    {}
func (*RevokeRefreshTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_refresh_token_1e8e54ee3f1fef8c, []int{6}
}
func (m *RevokeRefreshTokenRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RevokeRefreshTokenRequest.Unmarshal(m, b)
}
func (m *RevokeRefreshTokenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RevokeRefreshTokenRequest.Marshal(b, m, deterministic)
}
func (dst *RevokeRefreshTokenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RevokeRefreshTokenRequest.Merge(dst, src)
}
func (m *RevokeRefreshTokenRequest) XXX_Size() int {
	return xxx_messageInfo_RevokeRefreshTokenRequest.Size(m)
}
func (m *RevokeRefreshTokenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RevokeRefreshTokenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RevokeRefreshTokenRequest proto.InternalMessageInfo

func (m *RevokeRefreshTokenRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *RevokeRefreshTokenRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

type RevokeRefreshTokenResponse struct {
	RefreshToken         *RefreshToken `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *RevokeRefreshTokenResponse) Reset()         { *m = RevokeRefreshTokenResponse{} }
func (m *RevokeRefreshTokenResponse) String() string { return proto.CompactTextString(m) }
func (*RevokeRefreshTokenResponse) ProtoMessage()    {}
func (*RevokeRefreshTokenResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_refresh_token_1e8e54ee3f1fef8c, []int{7}
}
func (m *RevokeRefreshTokenResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RevokeRefreshTokenResponse.Unmarshal(m, b)
}
func (m *RevokeRefreshTokenResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RevokeRefreshTokenResponse.Marshal(b, m, deterministic)
}
func (dst *RevokeRefreshTokenResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RevokeRefreshTokenResponse.Merge(dst, src)
}
func (m *RevokeRefreshTokenResponse) XXX_Size() int {
	return xxx_messageInfo_RevokeRefreshTokenResponse.Size(m)
}
func (m *RevokeRefreshTokenResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RevokeRefreshTokenResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RevokeRefreshTokenResponse proto.InternalMessageInfo

func (m *RevokeRefreshTokenResponse) GetRefreshToken() *RefreshToken {
	if m != nil {
		return m.RefreshToken
	}
	return nil
}

func init() {
	proto.RegisterType((*RefreshToken)(nil), "charon.rpc.charond.v1.RefreshToken")
	proto.RegisterType((*RefreshTokenQuery)(nil), "charon.rpc.charond.v1.RefreshTokenQuery")
	proto.RegisterType((*CreateRefreshTokenRequest)(nil), "charon.rpc.charond.v1.CreateRefreshTokenRequest")
	proto.RegisterType((*CreateRefreshTokenResponse)(nil), "charon.rpc.charond.v1.CreateRefreshTokenResponse")
	proto.RegisterType((*ListRefreshTokensRequest)(nil), "charon.rpc.charond.v1.ListRefreshTokensRequest")
	proto.RegisterType((*ListRefreshTokensResponse)(nil), "charon.rpc.charond.v1.ListRefreshTokensResponse")
	proto.RegisterType((*RevokeRefreshTokenRequest)(nil), "charon.rpc.charond.v1.RevokeRefreshTokenRequest")
	proto.RegisterType((*RevokeRefreshTokenResponse)(nil), "charon.rpc.charond.v1.RevokeRefreshTokenResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RefreshTokenManagerClient is the client API for RefreshTokenManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RefreshTokenManagerClient interface {
	Create(ctx context.Context, in *CreateRefreshTokenRequest, opts ...grpc.CallOption) (*CreateRefreshTokenResponse, error)
	Revoke(ctx context.Context, in *RevokeRefreshTokenRequest, opts ...grpc.CallOption) (*RevokeRefreshTokenResponse, error)
	List(ctx context.Context, in *ListRefreshTokensRequest, opts ...grpc.CallOption) (*ListRefreshTokensResponse, error)
}

type refreshTokenManagerClient struct {
	cc *grpc.ClientConn
}

func NewRefreshTokenManagerClient(cc *grpc.ClientConn) RefreshTokenManagerClient {
	return &refreshTokenManagerClient{cc}
}

func (c *refreshTokenManagerClient) Create(ctx context.Context, in *CreateRefreshTokenRequest, opts ...grpc.CallOption) (*CreateRefreshTokenResponse, error) {
	out := new(CreateRefreshTokenResponse)
	err := c.cc.Invoke(ctx, "/charon.rpc.charond.v1.RefreshTokenManager/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *refreshTokenManagerClient) Revoke(ctx context.Context, in *RevokeRefreshTokenRequest, opts ...grpc.CallOption) (*RevokeRefreshTokenResponse, error) {
	out := new(RevokeRefreshTokenResponse)
	err := c.cc.Invoke(ctx, "/charon.rpc.charond.v1.RefreshTokenManager/Revoke", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *refreshTokenManagerClient) List(ctx context.Context, in *ListRefreshTokensRequest, opts ...grpc.CallOption) (*ListRefreshTokensResponse, error) {
	out := new(ListRefreshTokensResponse)
	err := c.cc.Invoke(ctx, "/charon.rpc.charond.v1.RefreshTokenManager/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RefreshTokenManagerServer is the server API for RefreshTokenManager service.
type RefreshTokenManagerServer interface {
	Create(context.Context, *CreateRefreshTokenRequest) (*CreateRefreshTokenResponse, error)
	Revoke(context.Context, *RevokeRefreshTokenRequest) (*RevokeRefreshTokenResponse, error)
	List(context.Context, *ListRefreshTokensRequest) (*ListRefreshTokensResponse, error)
}

func RegisterRefreshTokenManagerServer(s *grpc.Server, srv RefreshTokenManagerServer) {
	s.RegisterService(&_RefreshTokenManager_serviceDesc, srv)
}

func _RefreshTokenManager_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRefreshTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RefreshTokenManagerServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/charon.rpc.charond.v1.RefreshTokenManager/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RefreshTokenManagerServer).Create(ctx, req.(*CreateRefreshTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RefreshTokenManager_Revoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevokeRefreshTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RefreshTokenManagerServer).Revoke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/charon.rpc.charond.v1.RefreshTokenManager/Revoke",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RefreshTokenManagerServer).Revoke(ctx, req.(*RevokeRefreshTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RefreshTokenManager_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRefreshTokensRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RefreshTokenManagerServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/charon.rpc.charond.v1.RefreshTokenManager/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RefreshTokenManagerServer).List(ctx, req.(*ListRefreshTokensRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RefreshTokenManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "charon.rpc.charond.v1.RefreshTokenManager",
	HandlerType: (*RefreshTokenManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _RefreshTokenManager_Create_Handler,
		},
		{
			MethodName: "Revoke",
			Handler:    _RefreshTokenManager_Revoke_Handler,
		},
		{
			MethodName: "List",
			Handler:    _RefreshTokenManager_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/piotrkowalczuk/charon/pb/rpc/charond/v1/refresh_token.proto",
}

func init() {
	proto.RegisterFile("github.com/piotrkowalczuk/charon/pb/rpc/charond/v1/refresh_token.proto", fileDescriptor_refresh_token_1e8e54ee3f1fef8c)
}

var fileDescriptor_refresh_token_1e8e54ee3f1fef8c = []byte{
	// 737 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0xdd, 0x6e, 0xd3, 0x4a,
	0x10, 0xae, 0xf3, 0x9f, 0x49, 0x5a, 0x9d, 0x6e, 0xcf, 0xd1, 0x71, 0x23, 0x24, 0x22, 0x17, 0xaa,
	0x5c, 0x20, 0x3b, 0x6d, 0x11, 0x15, 0x3f, 0x02, 0x35, 0x48, 0x88, 0x56, 0x20, 0xc0, 0x94, 0x1b,
	0x6e, 0x82, 0x63, 0x6f, 0x52, 0x2b, 0x89, 0xd7, 0xd9, 0x5d, 0x17, 0xdc, 0x07, 0xe4, 0x01, 0x78,
	0x06, 0x78, 0x0f, 0x64, 0xaf, 0xb7, 0xb1, 0x8b, 0x4d, 0x53, 0xc4, 0x55, 0xb2, 0xbb, 0xdf, 0x37,
	0xb3, 0x33, 0xdf, 0xb7, 0x63, 0x78, 0x31, 0x71, 0xf9, 0x59, 0x30, 0xd2, 0x6d, 0x32, 0x37, 0x7c,
	0x97, 0x70, 0x3a, 0x25, 0x9f, 0xad, 0x99, 0x7d, 0x11, 0x4c, 0x0d, 0xfb, 0xcc, 0xa2, 0xc4, 0x33,
	0xfc, 0x91, 0x41, 0x7d, 0x3b, 0x59, 0x39, 0xc6, 0xf9, 0x9e, 0x41, 0xf1, 0x98, 0x62, 0x76, 0x36,
	0xe4, 0x64, 0x8a, 0x3d, 0xdd, 0xa7, 0x84, 0x13, 0xf4, 0x9f, 0x38, 0xd7, 0xa9, 0x6f, 0xeb, 0x09,
	0x54, 0x3f, 0xdf, 0xeb, 0xdc, 0x9e, 0x10, 0x32, 0x99, 0x61, 0x23, 0x06, 0x8d, 0x82, 0xb1, 0xc1,
	0xdd, 0x39, 0x66, 0xdc, 0x9a, 0xfb, 0x82, 0xd7, 0x79, 0xf6, 0x07, 0xf9, 0x6d, 0x32, 0x9f, 0x93,
	0x24, 0x71, 0x67, 0x6b, 0xc1, 0x43, 0x1f, 0x33, 0x43, 0xfc, 0xc8, 0x4d, 0x4f, 0x6c, 0x7a, 0xa9,
	0x4d, 0xed, 0x6b, 0x19, 0xda, 0xa6, 0xb8, 0xfa, 0x69, 0x74, 0x73, 0xf4, 0x2f, 0x54, 0xe3, 0x12,
	0x54, 0xa5, 0xab, 0xf4, 0x9a, 0xa6, 0x58, 0xa0, 0x3b, 0x50, 0xf5, 0x08, 0xc7, 0x4c, 0x2d, 0x75,
	0x95, 0x5e, 0x6b, 0x7f, 0x43, 0x4f, 0x82, 0xbc, 0xe7, 0xd4, 0xf5, 0x26, 0xa6, 0x38, 0x44, 0xff,
	0x43, 0x3d, 0x60, 0x98, 0x0e, 0x5d, 0x47, 0x2d, 0x77, 0x95, 0x5e, 0xd9, 0xac, 0x45, 0xcb, 0x63,
	0x07, 0xa9, 0x50, 0xa7, 0xf8, 0x9c, 0x4c, 0xb1, 0xa3, 0x56, 0xba, 0x4a, 0xaf, 0x61, 0xca, 0x25,
	0x3a, 0x84, 0x26, 0xfe, 0xe2, 0xbb, 0x14, 0x0f, 0x2d, 0xae, 0x56, 0xe3, 0xe0, 0x1d, 0x5d, 0xf4,
	0x47, 0x97, 0xfd, 0xd1, 0x4f, 0x65, 0x7f, 0xcc, 0x86, 0x00, 0x1f, 0x71, 0xf4, 0x04, 0xda, 0x33,
	0x8b, 0xf1, 0x61, 0xc0, 0xb0, 0x13, 0x71, 0x6b, 0xd7, 0x72, 0x21, 0xc2, 0x7f, 0x60, 0xd8, 0x39,
	0xe2, 0xe8, 0x21, 0x80, 0x4d, 0xb1, 0xc5, 0x05, 0xb7, 0x7e, 0x2d, 0xb7, 0x99, 0xa0, 0x8f, 0x38,
	0xba, 0xb7, 0xa4, 0x8e, 0x42, 0xb5, 0x11, 0x53, 0xd7, 0x65, 0x3f, 0x8e, 0x3d, 0xfe, 0xe0, 0xfe,
	0x25, 0x7a, 0x10, 0x46, 0x89, 0x02, 0xdf, 0x91, 0x89, 0x9a, 0xd7, 0x27, 0x4a, 0xd0, 0x22, 0x91,
	0xa4, 0x8e, 0x42, 0x15, 0x72, 0x13, 0x25, 0x80, 0x41, 0xa8, 0x7d, 0x2b, 0xc1, 0x66, 0x5a, 0xc8,
	0x77, 0x01, 0xa6, 0x21, 0xda, 0x5d, 0x2a, 0xa2, 0x24, 0x01, 0x16, 0xe9, 0x00, 0x52, 0xa0, 0x5f,
	0xf4, 0x5d, 0xe4, 0xe9, 0xbb, 0xbb, 0x94, 0xb1, 0x1c, 0xe3, 0xda, 0xf2, 0x3a, 0x03, 0x42, 0x66,
	0x4b, 0x51, 0xf5, 0xb4, 0xa8, 0x95, 0x18, 0xb9, 0x29, 0x23, 0xe6, 0x69, 0x79, 0x70, 0x45, 0xcb,
	0x6a, 0x11, 0x25, 0x2d, 0x61, 0x3f, 0x23, 0x61, 0xad, 0x88, 0x92, 0x52, 0xae, 0x9f, 0xd1, 0xa2,
	0x5e, 0xc8, 0xb8, 0x94, 0x40, 0xbb, 0x80, 0xed, 0xe7, 0x31, 0x3d, 0xdd, 0x59, 0x13, 0x2f, 0x02,
	0xcc, 0xf8, 0xb2, 0x67, 0xca, 0xef, 0xde, 0x44, 0xc6, 0xe0, 0xa5, 0xd5, 0x0d, 0xae, 0x8d, 0xa1,
	0x93, 0x97, 0x9b, 0xf9, 0xc4, 0x63, 0x18, 0xbd, 0x84, 0xf5, 0xcc, 0xc4, 0x49, 0x2e, 0xb1, 0xa3,
	0xe7, 0x8e, 0x1c, 0x3d, 0x13, 0xa3, 0x4d, 0x53, 0x2b, 0xed, 0x87, 0x02, 0xea, 0x2b, 0x97, 0xf1,
	0x34, 0x84, 0xc9, 0x1a, 0xef, 0x42, 0x8d, 0x8c, 0xc7, 0x0c, 0xf3, 0x4b, 0xfb, 0x64, 0xfc, 0x97,
	0x1c, 0xa2, 0x1d, 0xa8, 0xce, 0xdc, 0xb9, 0x2b, 0x0b, 0xbc, 0x82, 0x12, 0x67, 0xe8, 0x10, 0x1a,
	0x84, 0x3a, 0x98, 0x46, 0x6e, 0x2e, 0x77, 0xcb, 0xbd, 0xd6, 0xfe, 0xad, 0x82, 0xdb, 0xbe, 0x89,
	0x60, 0x66, 0x3d, 0x46, 0x0f, 0x42, 0xf4, 0x14, 0xaa, 0x8b, 0xc8, 0xcd, 0x6a, 0x2b, 0x8e, 0xde,
	0x5b, 0xa1, 0xc6, 0xd8, 0xfd, 0xa6, 0xa0, 0x9d, 0x54, 0x1a, 0x95, 0x7f, 0x5a, 0xda, 0x04, 0xb6,
	0x73, 0xca, 0x4c, 0xda, 0x79, 0x02, 0x1b, 0x99, 0x76, 0x46, 0xa2, 0x96, 0x57, 0xed, 0xe7, 0x7a,
	0xba, 0x9f, 0x4c, 0x3b, 0x81, 0x6d, 0x33, 0x7e, 0x08, 0x79, 0xa6, 0xc9, 0x1f, 0xaf, 0xa9, 0xc1,
	0x59, 0x4a, 0x0f, 0xce, 0xc8, 0x04, 0x79, 0xb1, 0xfe, 0xb6, 0x09, 0xf6, 0xbf, 0x97, 0x60, 0x2b,
	0x7d, 0xfc, 0xda, 0xf2, 0xac, 0x09, 0xa6, 0x88, 0x40, 0x4d, 0x98, 0x10, 0xf5, 0x0b, 0x82, 0x16,
	0xbe, 0x8f, 0xce, 0xde, 0x0d, 0x18, 0xa2, 0x20, 0x6d, 0x2d, 0x4a, 0x28, 0x0a, 0x2e, 0x4c, 0x58,
	0xd8, 0xdb, 0xc2, 0x84, 0xc5, 0x1d, 0xd4, 0xd6, 0xd0, 0x14, 0x2a, 0x91, 0x2d, 0x90, 0x51, 0x40,
	0x2e, 0x7a, 0x1a, 0x9d, 0xfe, 0xea, 0x04, 0x99, 0x6c, 0xf0, 0x09, 0xba, 0x36, 0x99, 0xeb, 0xf2,
	0xf3, 0x9e, 0xc7, 0x7f, 0xab, 0x7c, 0x7c, 0x74, 0xf3, 0xcf, 0xff, 0xe3, 0xe4, 0xef, 0xa8, 0x16,
	0xcf, 0x94, 0x83, 0x9f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xbf, 0x47, 0x0d, 0x45, 0xc3, 0x08, 0x00,
	0x00,
}

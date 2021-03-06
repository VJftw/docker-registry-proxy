// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/proto/v1/configuration_api.proto

package dockerregistryproxyv1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// GetConfigurationSchemaRequest is empty.
type GetConfigurationSchemaRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetConfigurationSchemaRequest) Reset()         { *m = GetConfigurationSchemaRequest{} }
func (m *GetConfigurationSchemaRequest) String() string { return proto.CompactTextString(m) }
func (*GetConfigurationSchemaRequest) ProtoMessage()    {}
func (*GetConfigurationSchemaRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_190512c21805dc3a, []int{0}
}

func (m *GetConfigurationSchemaRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetConfigurationSchemaRequest.Unmarshal(m, b)
}
func (m *GetConfigurationSchemaRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetConfigurationSchemaRequest.Marshal(b, m, deterministic)
}
func (m *GetConfigurationSchemaRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetConfigurationSchemaRequest.Merge(m, src)
}
func (m *GetConfigurationSchemaRequest) XXX_Size() int {
	return xxx_messageInfo_GetConfigurationSchemaRequest.Size(m)
}
func (m *GetConfigurationSchemaRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetConfigurationSchemaRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetConfigurationSchemaRequest proto.InternalMessageInfo

// GetConfigurationSchema represents the configuration for a plugin.
type GetConfigurationSchemaResponse struct {
	// key is the attribute name.
	Attributes           map[string]*ConfigurationAttribute `protobuf:"bytes,1,rep,name=attributes,proto3" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                           `json:"-"`
	XXX_unrecognized     []byte                             `json:"-"`
	XXX_sizecache        int32                              `json:"-"`
}

func (m *GetConfigurationSchemaResponse) Reset()         { *m = GetConfigurationSchemaResponse{} }
func (m *GetConfigurationSchemaResponse) String() string { return proto.CompactTextString(m) }
func (*GetConfigurationSchemaResponse) ProtoMessage()    {}
func (*GetConfigurationSchemaResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_190512c21805dc3a, []int{1}
}

func (m *GetConfigurationSchemaResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetConfigurationSchemaResponse.Unmarshal(m, b)
}
func (m *GetConfigurationSchemaResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetConfigurationSchemaResponse.Marshal(b, m, deterministic)
}
func (m *GetConfigurationSchemaResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetConfigurationSchemaResponse.Merge(m, src)
}
func (m *GetConfigurationSchemaResponse) XXX_Size() int {
	return xxx_messageInfo_GetConfigurationSchemaResponse.Size(m)
}
func (m *GetConfigurationSchemaResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetConfigurationSchemaResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetConfigurationSchemaResponse proto.InternalMessageInfo

func (m *GetConfigurationSchemaResponse) GetAttributes() map[string]*ConfigurationAttribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

type ConfigureRequest struct {
	Attributes           map[string]*ConfigurationAttributeValue `protobuf:"bytes,1,rep,name=attributes,proto3" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                                `json:"-"`
	XXX_unrecognized     []byte                                  `json:"-"`
	XXX_sizecache        int32                                   `json:"-"`
}

func (m *ConfigureRequest) Reset()         { *m = ConfigureRequest{} }
func (m *ConfigureRequest) String() string { return proto.CompactTextString(m) }
func (*ConfigureRequest) ProtoMessage()    {}
func (*ConfigureRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_190512c21805dc3a, []int{2}
}

func (m *ConfigureRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigureRequest.Unmarshal(m, b)
}
func (m *ConfigureRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigureRequest.Marshal(b, m, deterministic)
}
func (m *ConfigureRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigureRequest.Merge(m, src)
}
func (m *ConfigureRequest) XXX_Size() int {
	return xxx_messageInfo_ConfigureRequest.Size(m)
}
func (m *ConfigureRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigureRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigureRequest proto.InternalMessageInfo

func (m *ConfigureRequest) GetAttributes() map[string]*ConfigurationAttributeValue {
	if m != nil {
		return m.Attributes
	}
	return nil
}

type ConfigureResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigureResponse) Reset()         { *m = ConfigureResponse{} }
func (m *ConfigureResponse) String() string { return proto.CompactTextString(m) }
func (*ConfigureResponse) ProtoMessage()    {}
func (*ConfigureResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_190512c21805dc3a, []int{3}
}

func (m *ConfigureResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigureResponse.Unmarshal(m, b)
}
func (m *ConfigureResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigureResponse.Marshal(b, m, deterministic)
}
func (m *ConfigureResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigureResponse.Merge(m, src)
}
func (m *ConfigureResponse) XXX_Size() int {
	return xxx_messageInfo_ConfigureResponse.Size(m)
}
func (m *ConfigureResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigureResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigureResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*GetConfigurationSchemaRequest)(nil), "vjftw.dockerregistryproxy.v1.GetConfigurationSchemaRequest")
	proto.RegisterType((*GetConfigurationSchemaResponse)(nil), "vjftw.dockerregistryproxy.v1.GetConfigurationSchemaResponse")
	proto.RegisterMapType((map[string]*ConfigurationAttribute)(nil), "vjftw.dockerregistryproxy.v1.GetConfigurationSchemaResponse.AttributesEntry")
	proto.RegisterType((*ConfigureRequest)(nil), "vjftw.dockerregistryproxy.v1.ConfigureRequest")
	proto.RegisterMapType((map[string]*ConfigurationAttributeValue)(nil), "vjftw.dockerregistryproxy.v1.ConfigureRequest.AttributesEntry")
	proto.RegisterType((*ConfigureResponse)(nil), "vjftw.dockerregistryproxy.v1.ConfigureResponse")
}

func init() {
	proto.RegisterFile("api/proto/v1/configuration_api.proto", fileDescriptor_190512c21805dc3a)
}

var fileDescriptor_190512c21805dc3a = []byte{
	// 381 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x49, 0x2c, 0xc8, 0xd4,
	0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x33, 0xd4, 0x4f, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0x2f,
	0x2d, 0x4a, 0x2c, 0xc9, 0xcc, 0xcf, 0x8b, 0x4f, 0x2c, 0xc8, 0xd4, 0x03, 0x4b, 0x09, 0xc9, 0x94,
	0x65, 0xa5, 0x95, 0x94, 0xeb, 0xa5, 0xe4, 0x27, 0x67, 0xa7, 0x16, 0x15, 0xa5, 0xa6, 0x67, 0x16,
	0x97, 0x14, 0x55, 0x16, 0x14, 0xe5, 0x57, 0x54, 0xea, 0x95, 0x19, 0x4a, 0x29, 0xe0, 0x36, 0x03,
	0xa2, 0x5f, 0x49, 0x9e, 0x4b, 0xd6, 0x3d, 0xb5, 0xc4, 0x19, 0x59, 0x26, 0x38, 0x39, 0x23, 0x35,
	0x37, 0x31, 0x28, 0xb5, 0xb0, 0x34, 0xb5, 0xb8, 0x44, 0xa9, 0x99, 0x89, 0x4b, 0x0e, 0x97, 0x8a,
	0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0xa1, 0x1c, 0x2e, 0xae, 0xc4, 0x92, 0x92, 0xa2, 0xcc, 0xa4,
	0xd2, 0x92, 0xd4, 0x62, 0x09, 0x46, 0x05, 0x66, 0x0d, 0x6e, 0x23, 0x1f, 0x3d, 0x7c, 0x0e, 0xd3,
	0xc3, 0x6f, 0xa2, 0x9e, 0x23, 0xdc, 0x38, 0xd7, 0xbc, 0x92, 0xa2, 0xca, 0x20, 0x24, 0xf3, 0xa5,
	0x8a, 0xb9, 0xf8, 0xd1, 0xa4, 0x85, 0x04, 0xb8, 0x98, 0xb3, 0x53, 0x2b, 0x25, 0x18, 0x15, 0x18,
	0x35, 0x38, 0x83, 0x40, 0x4c, 0x21, 0x2f, 0x2e, 0xd6, 0xb2, 0xc4, 0x9c, 0xd2, 0x54, 0x09, 0x26,
	0x05, 0x46, 0x0d, 0x6e, 0x23, 0x13, 0xfc, 0xae, 0x41, 0x71, 0x0a, 0xdc, 0xf0, 0x20, 0x88, 0x11,
	0x56, 0x4c, 0x16, 0x8c, 0x4a, 0x6f, 0x18, 0xb9, 0x04, 0x60, 0xaa, 0x52, 0xa1, 0x41, 0x23, 0x14,
	0x87, 0xc5, 0xdf, 0x76, 0xc4, 0xd9, 0x04, 0x33, 0x03, 0xaf, 0x4f, 0x2b, 0x88, 0xf1, 0xa9, 0x3f,
	0xaa, 0x4f, 0x2d, 0xc9, 0xf1, 0x69, 0x18, 0xc8, 0x00, 0x64, 0xef, 0x0a, 0x73, 0x09, 0x22, 0xb9,
	0x14, 0x12, 0x29, 0x46, 0xd3, 0x98, 0x10, 0x61, 0x00, 0xd1, 0x1f, 0xe0, 0x29, 0x34, 0x99, 0x91,
	0x4b, 0x0c, 0x7b, 0x64, 0x0a, 0x59, 0x93, 0x97, 0x04, 0xc0, 0xe1, 0x22, 0x65, 0x43, 0x49, 0xfa,
	0x11, 0xca, 0xe1, 0xe2, 0x84, 0xbb, 0x5f, 0x48, 0x8f, 0xb4, 0x28, 0x91, 0xd2, 0x27, 0x5a, 0x3d,
	0xc4, 0x36, 0xa7, 0x69, 0x8c, 0x5c, 0x0a, 0xc9, 0xf9, 0xb9, 0x78, 0xb5, 0x39, 0x89, 0xa2, 0x06,
	0x5d, 0x41, 0x66, 0x00, 0x28, 0xff, 0x05, 0x30, 0x46, 0x89, 0x62, 0xd1, 0x50, 0x66, 0xb8, 0x88,
	0x89, 0x39, 0xcc, 0x25, 0x62, 0x15, 0x93, 0x4c, 0x18, 0xd8, 0x50, 0x17, 0x2c, 0x86, 0x86, 0x19,
	0x9e, 0x82, 0x4a, 0xc7, 0x60, 0x91, 0x8e, 0x09, 0x33, 0x4c, 0x62, 0x03, 0xe7, 0x71, 0x63, 0x40,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x4c, 0x08, 0xb1, 0xea, 0x4b, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ConfigurationAPIClient is the client API for ConfigurationAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConfigurationAPIClient interface {
	// GetConfigurationSchema returns the schema for the plugin.
	GetConfigurationSchema(ctx context.Context, in *GetConfigurationSchemaRequest, opts ...grpc.CallOption) (*GetConfigurationSchemaResponse, error)
	// Configure configures the plugin.
	Configure(ctx context.Context, in *ConfigureRequest, opts ...grpc.CallOption) (*ConfigureResponse, error)
}

type configurationAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigurationAPIClient(cc grpc.ClientConnInterface) ConfigurationAPIClient {
	return &configurationAPIClient{cc}
}

func (c *configurationAPIClient) GetConfigurationSchema(ctx context.Context, in *GetConfigurationSchemaRequest, opts ...grpc.CallOption) (*GetConfigurationSchemaResponse, error) {
	out := new(GetConfigurationSchemaResponse)
	err := c.cc.Invoke(ctx, "/vjftw.dockerregistryproxy.v1.ConfigurationAPI/GetConfigurationSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configurationAPIClient) Configure(ctx context.Context, in *ConfigureRequest, opts ...grpc.CallOption) (*ConfigureResponse, error) {
	out := new(ConfigureResponse)
	err := c.cc.Invoke(ctx, "/vjftw.dockerregistryproxy.v1.ConfigurationAPI/Configure", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigurationAPIServer is the server API for ConfigurationAPI service.
type ConfigurationAPIServer interface {
	// GetConfigurationSchema returns the schema for the plugin.
	GetConfigurationSchema(context.Context, *GetConfigurationSchemaRequest) (*GetConfigurationSchemaResponse, error)
	// Configure configures the plugin.
	Configure(context.Context, *ConfigureRequest) (*ConfigureResponse, error)
}

// UnimplementedConfigurationAPIServer can be embedded to have forward compatible implementations.
type UnimplementedConfigurationAPIServer struct {
}

func (*UnimplementedConfigurationAPIServer) GetConfigurationSchema(ctx context.Context, req *GetConfigurationSchemaRequest) (*GetConfigurationSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfigurationSchema not implemented")
}
func (*UnimplementedConfigurationAPIServer) Configure(ctx context.Context, req *ConfigureRequest) (*ConfigureResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Configure not implemented")
}

func RegisterConfigurationAPIServer(s *grpc.Server, srv ConfigurationAPIServer) {
	s.RegisterService(&_ConfigurationAPI_serviceDesc, srv)
}

func _ConfigurationAPI_GetConfigurationSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConfigurationSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigurationAPIServer).GetConfigurationSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vjftw.dockerregistryproxy.v1.ConfigurationAPI/GetConfigurationSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigurationAPIServer).GetConfigurationSchema(ctx, req.(*GetConfigurationSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigurationAPI_Configure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigurationAPIServer).Configure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vjftw.dockerregistryproxy.v1.ConfigurationAPI/Configure",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigurationAPIServer).Configure(ctx, req.(*ConfigureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ConfigurationAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vjftw.dockerregistryproxy.v1.ConfigurationAPI",
	HandlerType: (*ConfigurationAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetConfigurationSchema",
			Handler:    _ConfigurationAPI_GetConfigurationSchema_Handler,
		},
		{
			MethodName: "Configure",
			Handler:    _ConfigurationAPI_Configure_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/v1/configuration_api.proto",
}

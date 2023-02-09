// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: add_service.proto

package protoservices

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type AddServiceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceName string `protobuf:"bytes,1,opt,name=serviceName,proto3" json:"serviceName,omitempty"`
}

func (x *AddServiceRequest) Reset() {
	*x = AddServiceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_add_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddServiceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddServiceRequest) ProtoMessage() {}

func (x *AddServiceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_add_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddServiceRequest.ProtoReflect.Descriptor instead.
func (*AddServiceRequest) Descriptor() ([]byte, []int) {
	return file_add_service_proto_rawDescGZIP(), []int{0}
}

func (x *AddServiceRequest) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

type AddServiceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *AddServiceResponse) Reset() {
	*x = AddServiceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_add_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddServiceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddServiceResponse) ProtoMessage() {}

func (x *AddServiceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_add_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddServiceResponse.ProtoReflect.Descriptor instead.
func (*AddServiceResponse) Descriptor() ([]byte, []int) {
	return file_add_service_proto_rawDescGZIP(), []int{1}
}

func (x *AddServiceResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_add_service_proto protoreflect.FileDescriptor

var file_add_service_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x64, 0x64, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x14, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x22, 0x35, 0x0a, 0x11, 0x41, 0x64, 0x64,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20,
	0x0a, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x22, 0x2e, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x32, 0x71, 0x0a, 0x0a, 0x41, 0x64, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x63,
	0x0a, 0x0a, 0x41, 0x64, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x27, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2e, 0x41, 0x64, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x41, 0x64, 0x64,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x28, 0x01, 0x42, 0x11, 0x5a, 0x0f, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_add_service_proto_rawDescOnce sync.Once
	file_add_service_proto_rawDescData = file_add_service_proto_rawDesc
)

func file_add_service_proto_rawDescGZIP() []byte {
	file_add_service_proto_rawDescOnce.Do(func() {
		file_add_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_add_service_proto_rawDescData)
	})
	return file_add_service_proto_rawDescData
}

var file_add_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_add_service_proto_goTypes = []interface{}{
	(*AddServiceRequest)(nil),  // 0: server.protoservices.AddServiceRequest
	(*AddServiceResponse)(nil), // 1: server.protoservices.AddServiceResponse
}
var file_add_service_proto_depIdxs = []int32{
	0, // 0: server.protoservices.AddService.AddService:input_type -> server.protoservices.AddServiceRequest
	1, // 1: server.protoservices.AddService.AddService:output_type -> server.protoservices.AddServiceResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_add_service_proto_init() }
func file_add_service_proto_init() {
	if File_add_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_add_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddServiceRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_add_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddServiceResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_add_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_add_service_proto_goTypes,
		DependencyIndexes: file_add_service_proto_depIdxs,
		MessageInfos:      file_add_service_proto_msgTypes,
	}.Build()
	File_add_service_proto = out.File
	file_add_service_proto_rawDesc = nil
	file_add_service_proto_goTypes = nil
	file_add_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AddServiceClient is the client API for AddService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AddServiceClient interface {
	AddService(ctx context.Context, opts ...grpc.CallOption) (AddService_AddServiceClient, error)
}

type addServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAddServiceClient(cc grpc.ClientConnInterface) AddServiceClient {
	return &addServiceClient{cc}
}

func (c *addServiceClient) AddService(ctx context.Context, opts ...grpc.CallOption) (AddService_AddServiceClient, error) {
	stream, err := c.cc.NewStream(ctx, &_AddService_serviceDesc.Streams[0], "/server.protoservices.AddService/AddService", opts...)
	if err != nil {
		return nil, err
	}
	x := &addServiceAddServiceClient{stream}
	return x, nil
}

type AddService_AddServiceClient interface {
	Send(*AddServiceRequest) error
	CloseAndRecv() (*AddServiceResponse, error)
	grpc.ClientStream
}

type addServiceAddServiceClient struct {
	grpc.ClientStream
}

func (x *addServiceAddServiceClient) Send(m *AddServiceRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *addServiceAddServiceClient) CloseAndRecv() (*AddServiceResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(AddServiceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AddServiceServer is the server API for AddService service.
type AddServiceServer interface {
	AddService(AddService_AddServiceServer) error
}

// UnimplementedAddServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAddServiceServer struct {
}

func (*UnimplementedAddServiceServer) AddService(AddService_AddServiceServer) error {
	return status.Errorf(codes.Unimplemented, "method AddService not implemented")
}

func RegisterAddServiceServer(s *grpc.Server, srv AddServiceServer) {
	s.RegisterService(&_AddService_serviceDesc, srv)
}

func _AddService_AddService_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AddServiceServer).AddService(&addServiceAddServiceServer{stream})
}

type AddService_AddServiceServer interface {
	SendAndClose(*AddServiceResponse) error
	Recv() (*AddServiceRequest, error)
	grpc.ServerStream
}

type addServiceAddServiceServer struct {
	grpc.ServerStream
}

func (x *addServiceAddServiceServer) SendAndClose(m *AddServiceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *addServiceAddServiceServer) Recv() (*AddServiceRequest, error) {
	m := new(AddServiceRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _AddService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "server.protoservices.AddService",
	HandlerType: (*AddServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "AddService",
			Handler:       _AddService_AddService_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "add_service.proto",
}
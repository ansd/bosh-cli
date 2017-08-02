// Code generated by protoc-gen-go.
// source: metrics.proto
// DO NOT EDIT!

/*
Package grpc_testing is a generated protocol buffer package.

It is generated from these files:
	metrics.proto

It has these top-level messages:
	GaugeResponse
	GaugeRequest
	EmptyMessage
*/
package grpc_testing

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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

// Reponse message containing the gauge name and value
type GaugeResponse struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Types that are valid to be assigned to Value:
	//	*GaugeResponse_LongValue
	//	*GaugeResponse_DoubleValue
	//	*GaugeResponse_StringValue
	Value isGaugeResponse_Value `protobuf_oneof:"value"`
}

func (m *GaugeResponse) Reset()                    { *m = GaugeResponse{} }
func (m *GaugeResponse) String() string            { return proto.CompactTextString(m) }
func (*GaugeResponse) ProtoMessage()               {}
func (*GaugeResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type isGaugeResponse_Value interface {
	isGaugeResponse_Value()
}

type GaugeResponse_LongValue struct {
	LongValue int64 `protobuf:"varint,2,opt,name=long_value,json=longValue,oneof"`
}
type GaugeResponse_DoubleValue struct {
	DoubleValue float64 `protobuf:"fixed64,3,opt,name=double_value,json=doubleValue,oneof"`
}
type GaugeResponse_StringValue struct {
	StringValue string `protobuf:"bytes,4,opt,name=string_value,json=stringValue,oneof"`
}

func (*GaugeResponse_LongValue) isGaugeResponse_Value()   {}
func (*GaugeResponse_DoubleValue) isGaugeResponse_Value() {}
func (*GaugeResponse_StringValue) isGaugeResponse_Value() {}

func (m *GaugeResponse) GetValue() isGaugeResponse_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *GaugeResponse) GetLongValue() int64 {
	if x, ok := m.GetValue().(*GaugeResponse_LongValue); ok {
		return x.LongValue
	}
	return 0
}

func (m *GaugeResponse) GetDoubleValue() float64 {
	if x, ok := m.GetValue().(*GaugeResponse_DoubleValue); ok {
		return x.DoubleValue
	}
	return 0
}

func (m *GaugeResponse) GetStringValue() string {
	if x, ok := m.GetValue().(*GaugeResponse_StringValue); ok {
		return x.StringValue
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*GaugeResponse) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _GaugeResponse_OneofMarshaler, _GaugeResponse_OneofUnmarshaler, _GaugeResponse_OneofSizer, []interface{}{
		(*GaugeResponse_LongValue)(nil),
		(*GaugeResponse_DoubleValue)(nil),
		(*GaugeResponse_StringValue)(nil),
	}
}

func _GaugeResponse_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*GaugeResponse)
	// value
	switch x := m.Value.(type) {
	case *GaugeResponse_LongValue:
		b.EncodeVarint(2<<3 | proto.WireVarint)
		b.EncodeVarint(uint64(x.LongValue))
	case *GaugeResponse_DoubleValue:
		b.EncodeVarint(3<<3 | proto.WireFixed64)
		b.EncodeFixed64(math.Float64bits(x.DoubleValue))
	case *GaugeResponse_StringValue:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.StringValue)
	case nil:
	default:
		return fmt.Errorf("GaugeResponse.Value has unexpected type %T", x)
	}
	return nil
}

func _GaugeResponse_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*GaugeResponse)
	switch tag {
	case 2: // value.long_value
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Value = &GaugeResponse_LongValue{int64(x)}
		return true, err
	case 3: // value.double_value
		if wire != proto.WireFixed64 {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeFixed64()
		m.Value = &GaugeResponse_DoubleValue{math.Float64frombits(x)}
		return true, err
	case 4: // value.string_value
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Value = &GaugeResponse_StringValue{x}
		return true, err
	default:
		return false, nil
	}
}

func _GaugeResponse_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*GaugeResponse)
	// value
	switch x := m.Value.(type) {
	case *GaugeResponse_LongValue:
		n += proto.SizeVarint(2<<3 | proto.WireVarint)
		n += proto.SizeVarint(uint64(x.LongValue))
	case *GaugeResponse_DoubleValue:
		n += proto.SizeVarint(3<<3 | proto.WireFixed64)
		n += 8
	case *GaugeResponse_StringValue:
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.StringValue)))
		n += len(x.StringValue)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Request message containing the gauge name
type GaugeRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *GaugeRequest) Reset()                    { *m = GaugeRequest{} }
func (m *GaugeRequest) String() string            { return proto.CompactTextString(m) }
func (*GaugeRequest) ProtoMessage()               {}
func (*GaugeRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type EmptyMessage struct {
}

func (m *EmptyMessage) Reset()                    { *m = EmptyMessage{} }
func (m *EmptyMessage) String() string            { return proto.CompactTextString(m) }
func (*EmptyMessage) ProtoMessage()               {}
func (*EmptyMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func init() {
	proto.RegisterType((*GaugeResponse)(nil), "grpc.testing.GaugeResponse")
	proto.RegisterType((*GaugeRequest)(nil), "grpc.testing.GaugeRequest")
	proto.RegisterType((*EmptyMessage)(nil), "grpc.testing.EmptyMessage")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for MetricsService service

type MetricsServiceClient interface {
	// Returns the values of all the gauges that are currently being maintained by
	// the service
	GetAllGauges(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (MetricsService_GetAllGaugesClient, error)
	// Returns the value of one gauge
	GetGauge(ctx context.Context, in *GaugeRequest, opts ...grpc.CallOption) (*GaugeResponse, error)
}

type metricsServiceClient struct {
	cc *grpc.ClientConn
}

func NewMetricsServiceClient(cc *grpc.ClientConn) MetricsServiceClient {
	return &metricsServiceClient{cc}
}

func (c *metricsServiceClient) GetAllGauges(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (MetricsService_GetAllGaugesClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_MetricsService_serviceDesc.Streams[0], c.cc, "/grpc.testing.MetricsService/GetAllGauges", opts...)
	if err != nil {
		return nil, err
	}
	x := &metricsServiceGetAllGaugesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MetricsService_GetAllGaugesClient interface {
	Recv() (*GaugeResponse, error)
	grpc.ClientStream
}

type metricsServiceGetAllGaugesClient struct {
	grpc.ClientStream
}

func (x *metricsServiceGetAllGaugesClient) Recv() (*GaugeResponse, error) {
	m := new(GaugeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *metricsServiceClient) GetGauge(ctx context.Context, in *GaugeRequest, opts ...grpc.CallOption) (*GaugeResponse, error) {
	out := new(GaugeResponse)
	err := grpc.Invoke(ctx, "/grpc.testing.MetricsService/GetGauge", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for MetricsService service

type MetricsServiceServer interface {
	// Returns the values of all the gauges that are currently being maintained by
	// the service
	GetAllGauges(*EmptyMessage, MetricsService_GetAllGaugesServer) error
	// Returns the value of one gauge
	GetGauge(context.Context, *GaugeRequest) (*GaugeResponse, error)
}

func RegisterMetricsServiceServer(s *grpc.Server, srv MetricsServiceServer) {
	s.RegisterService(&_MetricsService_serviceDesc, srv)
}

func _MetricsService_GetAllGauges_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EmptyMessage)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MetricsServiceServer).GetAllGauges(m, &metricsServiceGetAllGaugesServer{stream})
}

type MetricsService_GetAllGaugesServer interface {
	Send(*GaugeResponse) error
	grpc.ServerStream
}

type metricsServiceGetAllGaugesServer struct {
	grpc.ServerStream
}

func (x *metricsServiceGetAllGaugesServer) Send(m *GaugeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _MetricsService_GetGauge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GaugeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).GetGauge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.testing.MetricsService/GetGauge",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).GetGauge(ctx, req.(*GaugeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MetricsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.testing.MetricsService",
	HandlerType: (*MetricsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetGauge",
			Handler:    _MetricsService_GetGauge_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAllGauges",
			Handler:       _MetricsService_GetAllGauges_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "metrics.proto",
}

func init() { proto.RegisterFile("metrics.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 253 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4d, 0x2d, 0x29,
	0xca, 0x4c, 0x2e, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x49, 0x2f, 0x2a, 0x48, 0xd6,
	0x2b, 0x49, 0x2d, 0x2e, 0xc9, 0xcc, 0x4b, 0x57, 0x9a, 0xce, 0xc8, 0xc5, 0xeb, 0x9e, 0x58, 0x9a,
	0x9e, 0x1a, 0x94, 0x5a, 0x5c, 0x90, 0x9f, 0x57, 0x9c, 0x2a, 0x24, 0xc4, 0xc5, 0x92, 0x97, 0x98,
	0x9b, 0x2a, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0x66, 0x0b, 0xc9, 0x73, 0x71, 0xe5, 0xe4,
	0xe7, 0xa5, 0xc7, 0x97, 0x25, 0xe6, 0x94, 0xa6, 0x4a, 0x30, 0x01, 0x65, 0x98, 0x3d, 0x18, 0x82,
	0x38, 0x41, 0x62, 0x61, 0x20, 0x21, 0x21, 0x65, 0x2e, 0x9e, 0x94, 0xfc, 0xd2, 0xa4, 0x9c, 0x54,
	0xa8, 0x12, 0x66, 0xa0, 0x12, 0x46, 0xa0, 0x12, 0x6e, 0x88, 0x28, 0x5c, 0x51, 0x31, 0xd0, 0x25,
	0x70, 0x73, 0x58, 0x40, 0x36, 0x80, 0x14, 0x41, 0x44, 0xc1, 0x8a, 0x9c, 0xd8, 0xb9, 0x58, 0xc1,
	0xb2, 0x4a, 0x4a, 0x5c, 0x3c, 0x50, 0x87, 0x15, 0x96, 0x02, 0x1d, 0x8b, 0xcd, 0x5d, 0x4a, 0x7c,
	0x5c, 0x3c, 0xae, 0xb9, 0x05, 0x25, 0x95, 0xbe, 0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9, 0x46, 0x0b,
	0x18, 0xb9, 0xf8, 0x7c, 0x21, 0xbe, 0x0d, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0x15, 0xf2, 0x04,
	0x1a, 0x93, 0x5a, 0xe2, 0x98, 0x93, 0x03, 0x36, 0xac, 0x58, 0x48, 0x4a, 0x0f, 0xd9, 0xff, 0x7a,
	0xc8, 0xda, 0xa5, 0xa4, 0x51, 0xe5, 0x50, 0xc2, 0xc5, 0x80, 0x51, 0xc8, 0x99, 0x8b, 0x03, 0x68,
	0x14, 0x58, 0x14, 0xdd, 0x18, 0x64, 0x97, 0xe2, 0x35, 0x26, 0x89, 0x0d, 0x1c, 0x0b, 0xc6, 0x80,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x5e, 0x7d, 0xb2, 0xc9, 0x96, 0x01, 0x00, 0x00,
}
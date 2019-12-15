// Code generated by protoc-gen-go. DO NOT EDIT.
// source: app/proto/event.proto

package rpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type EventIdentifier struct {
	Value                int64    `protobuf:"varint,1,opt,name=Value,proto3" json:"Value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventIdentifier) Reset()         { *m = EventIdentifier{} }
func (m *EventIdentifier) String() string { return proto.CompactTextString(m) }
func (*EventIdentifier) ProtoMessage()    {}
func (*EventIdentifier) Descriptor() ([]byte, []int) {
	return fileDescriptor_935e3c3e3df46d27, []int{0}
}

func (m *EventIdentifier) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventIdentifier.Unmarshal(m, b)
}
func (m *EventIdentifier) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventIdentifier.Marshal(b, m, deterministic)
}
func (m *EventIdentifier) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventIdentifier.Merge(m, src)
}
func (m *EventIdentifier) XXX_Size() int {
	return xxx_messageInfo_EventIdentifier.Size(m)
}
func (m *EventIdentifier) XXX_DiscardUnknown() {
	xxx_messageInfo_EventIdentifier.DiscardUnknown(m)
}

var xxx_messageInfo_EventIdentifier proto.InternalMessageInfo

func (m *EventIdentifier) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type EventListQuery struct {
	UserID               int64                `protobuf:"varint,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	From                 *timestamp.Timestamp `protobuf:"bytes,2,opt,name=From,proto3" json:"From,omitempty"`
	To                   *timestamp.Timestamp `protobuf:"bytes,3,opt,name=To,proto3" json:"To,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *EventListQuery) Reset()         { *m = EventListQuery{} }
func (m *EventListQuery) String() string { return proto.CompactTextString(m) }
func (*EventListQuery) ProtoMessage()    {}
func (*EventListQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_935e3c3e3df46d27, []int{1}
}

func (m *EventListQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventListQuery.Unmarshal(m, b)
}
func (m *EventListQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventListQuery.Marshal(b, m, deterministic)
}
func (m *EventListQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventListQuery.Merge(m, src)
}
func (m *EventListQuery) XXX_Size() int {
	return xxx_messageInfo_EventListQuery.Size(m)
}
func (m *EventListQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_EventListQuery.DiscardUnknown(m)
}

var xxx_messageInfo_EventListQuery proto.InternalMessageInfo

func (m *EventListQuery) GetUserID() int64 {
	if m != nil {
		return m.UserID
	}
	return 0
}

func (m *EventListQuery) GetFrom() *timestamp.Timestamp {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *EventListQuery) GetTo() *timestamp.Timestamp {
	if m != nil {
		return m.To
	}
	return nil
}

type EventListReply struct {
	Items                []*Event `protobuf:"bytes,1,rep,name=Items,proto3" json:"Items,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventListReply) Reset()         { *m = EventListReply{} }
func (m *EventListReply) String() string { return proto.CompactTextString(m) }
func (*EventListReply) ProtoMessage()    {}
func (*EventListReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_935e3c3e3df46d27, []int{2}
}

func (m *EventListReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventListReply.Unmarshal(m, b)
}
func (m *EventListReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventListReply.Marshal(b, m, deterministic)
}
func (m *EventListReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventListReply.Merge(m, src)
}
func (m *EventListReply) XXX_Size() int {
	return xxx_messageInfo_EventListReply.Size(m)
}
func (m *EventListReply) XXX_DiscardUnknown() {
	xxx_messageInfo_EventListReply.DiscardUnknown(m)
}

var xxx_messageInfo_EventListReply proto.InternalMessageInfo

func (m *EventListReply) GetItems() []*Event {
	if m != nil {
		return m.Items
	}
	return nil
}

type Event struct {
	ID                   int64                `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Title                string               `protobuf:"bytes,2,opt,name=Title,proto3" json:"Title,omitempty"`
	Description          string               `protobuf:"bytes,3,opt,name=Description,proto3" json:"Description,omitempty"`
	Location             string               `protobuf:"bytes,4,opt,name=Location,proto3" json:"Location,omitempty"`
	StartTime            *timestamp.Timestamp `protobuf:"bytes,5,opt,name=StartTime,proto3" json:"StartTime,omitempty"`
	EndTime              *timestamp.Timestamp `protobuf:"bytes,6,opt,name=EndTime,proto3" json:"EndTime,omitempty"`
	NotifyBefore         int64                `protobuf:"varint,7,opt,name=NotifyBefore,proto3" json:"NotifyBefore,omitempty"`
	UserID               int64                `protobuf:"varint,8,opt,name=UserID,proto3" json:"UserID,omitempty"`
	CalendarID           int64                `protobuf:"varint,9,opt,name=CalendarID,proto3" json:"CalendarID,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_935e3c3e3df46d27, []int{3}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Event) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Event) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Event) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

func (m *Event) GetStartTime() *timestamp.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *Event) GetEndTime() *timestamp.Timestamp {
	if m != nil {
		return m.EndTime
	}
	return nil
}

func (m *Event) GetNotifyBefore() int64 {
	if m != nil {
		return m.NotifyBefore
	}
	return 0
}

func (m *Event) GetUserID() int64 {
	if m != nil {
		return m.UserID
	}
	return 0
}

func (m *Event) GetCalendarID() int64 {
	if m != nil {
		return m.CalendarID
	}
	return 0
}

func init() {
	proto.RegisterType((*EventIdentifier)(nil), "calendar.EventIdentifier")
	proto.RegisterType((*EventListQuery)(nil), "calendar.EventListQuery")
	proto.RegisterType((*EventListReply)(nil), "calendar.EventListReply")
	proto.RegisterType((*Event)(nil), "calendar.Event")
}

func init() { proto.RegisterFile("app/proto/event.proto", fileDescriptor_935e3c3e3df46d27) }

var fileDescriptor_935e3c3e3df46d27 = []byte{
	// 434 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xcd, 0x6a, 0xdb, 0x40,
	0x10, 0xc7, 0x91, 0x6c, 0x29, 0xf6, 0x38, 0x24, 0x30, 0xb4, 0x65, 0xab, 0x43, 0x6b, 0x04, 0xa5,
	0xc6, 0x07, 0x19, 0xdc, 0x42, 0x7b, 0xca, 0x21, 0x51, 0x0a, 0x82, 0x50, 0xa8, 0xe2, 0xf4, 0xd0,
	0xdb, 0x46, 0x1a, 0x87, 0x05, 0x59, 0xbb, 0xac, 0xd6, 0x01, 0xdf, 0xfb, 0x14, 0x7d, 0x91, 0xbe,
	0x5e, 0xd1, 0xca, 0x76, 0x64, 0xf5, 0xc3, 0xbd, 0x69, 0xe6, 0xff, 0x1b, 0xed, 0x7f, 0x3e, 0xe0,
	0x39, 0x57, 0x6a, 0xa6, 0xb4, 0x34, 0x72, 0x46, 0x8f, 0x54, 0x9a, 0xc8, 0x7e, 0xe3, 0x20, 0xe3,
	0x05, 0x95, 0x39, 0xd7, 0xc1, 0xeb, 0x07, 0x29, 0x1f, 0x0a, 0x6a, 0x98, 0xfb, 0xf5, 0x72, 0x66,
	0xc4, 0x8a, 0x2a, 0xc3, 0x57, 0xaa, 0x41, 0xc3, 0xb7, 0x70, 0x7e, 0x5d, 0x57, 0x26, 0x39, 0x95,
	0x46, 0x2c, 0x05, 0x69, 0x7c, 0x06, 0xde, 0x57, 0x5e, 0xac, 0x89, 0x39, 0x63, 0x67, 0xd2, 0x4b,
	0x9b, 0x20, 0xfc, 0xee, 0xc0, 0x99, 0x25, 0x6f, 0x44, 0x65, 0xbe, 0xac, 0x49, 0x6f, 0xf0, 0x05,
	0xf8, 0x77, 0x15, 0xe9, 0x24, 0xde, 0x92, 0xdb, 0x08, 0x23, 0xe8, 0x7f, 0xd2, 0x72, 0xc5, 0xdc,
	0xb1, 0x33, 0x19, 0xcd, 0x83, 0xa8, 0xf1, 0x10, 0xed, 0x3c, 0x44, 0x8b, 0x9d, 0x87, 0xd4, 0x72,
	0x38, 0x05, 0x77, 0x21, 0x59, 0xef, 0x28, 0xed, 0x2e, 0x64, 0xf8, 0xa1, 0xe5, 0x22, 0x25, 0x55,
	0x6c, 0xf0, 0x0d, 0x78, 0x89, 0xa1, 0x55, 0xc5, 0x9c, 0x71, 0x6f, 0x32, 0x9a, 0x9f, 0x47, 0xbb,
	0xe6, 0x23, 0x0b, 0xa6, 0x8d, 0x1a, 0xfe, 0x74, 0xc1, 0xb3, 0x09, 0x3c, 0x03, 0x77, 0x6f, 0xd9,
	0x4d, 0xe2, 0xba, 0xdf, 0x85, 0x30, 0x05, 0x59, 0xbf, 0xc3, 0xb4, 0x09, 0x70, 0x0c, 0xa3, 0x98,
	0xaa, 0x4c, 0x0b, 0x65, 0x84, 0x2c, 0xad, 0xbb, 0x61, 0xda, 0x4e, 0x61, 0x00, 0x83, 0x1b, 0x99,
	0x71, 0x2b, 0xf7, 0xad, 0xbc, 0x8f, 0xf1, 0x23, 0x0c, 0x6f, 0x0d, 0xd7, 0xa6, 0x36, 0xcf, 0xbc,
	0xa3, 0x9d, 0x3d, 0xc1, 0xf8, 0x1e, 0x4e, 0xae, 0xcb, 0xdc, 0xd6, 0xf9, 0x47, 0xeb, 0x76, 0x28,
	0x86, 0x70, 0xfa, 0x59, 0x1a, 0xb1, 0xdc, 0x5c, 0xd2, 0x52, 0x6a, 0x62, 0x27, 0xb6, 0xbb, 0x83,
	0x5c, 0x6b, 0x5d, 0x83, 0x83, 0x75, 0xbd, 0x02, 0xb8, 0xda, 0x8e, 0x2c, 0x89, 0xd9, 0xd0, 0x6a,
	0xad, 0xcc, 0xfc, 0x87, 0x0b, 0xa7, 0x76, 0x72, 0xb7, 0xa4, 0x1f, 0x45, 0x46, 0x38, 0x05, 0xff,
	0x4a, 0x13, 0x37, 0x84, 0xdd, 0x61, 0x07, 0xdd, 0x04, 0xce, 0xa1, 0x9f, 0x12, 0xcf, 0xf1, 0x65,
	0x47, 0x78, 0xba, 0xb7, 0xdf, 0x6b, 0xa6, 0xe0, 0xdf, 0xa9, 0xfc, 0xff, 0xfe, 0x7f, 0x01, 0x7e,
	0x4c, 0x05, 0x19, 0xfa, 0xd7, 0x0b, 0x7f, 0x97, 0xf0, 0x02, 0x06, 0xb5, 0xbf, 0xfa, 0x9c, 0x90,
	0x75, 0xb0, 0xfd, 0xa5, 0x07, 0x7f, 0x52, 0xec, 0xf5, 0x5d, 0x7a, 0xdf, 0x7a, 0x5a, 0x65, 0xf7,
	0xbe, 0x5d, 0xce, 0xbb, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa9, 0x79, 0x26, 0xc3, 0x91, 0x03,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EventServiceClient is the client API for EventService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EventServiceClient interface {
	Create(ctx context.Context, in *Event, opts ...grpc.CallOption) (*Event, error)
	Read(ctx context.Context, in *EventIdentifier, opts ...grpc.CallOption) (*Event, error)
	Update(ctx context.Context, in *Event, opts ...grpc.CallOption) (*Event, error)
	Delete(ctx context.Context, in *EventIdentifier, opts ...grpc.CallOption) (*EventIdentifier, error)
	ReadList(ctx context.Context, in *EventListQuery, opts ...grpc.CallOption) (*EventListReply, error)
}

type eventServiceClient struct {
	cc *grpc.ClientConn
}

func NewEventServiceClient(cc *grpc.ClientConn) EventServiceClient {
	return &eventServiceClient{cc}
}

func (c *eventServiceClient) Create(ctx context.Context, in *Event, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/calendar.EventService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) Read(ctx context.Context, in *EventIdentifier, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/calendar.EventService/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) Update(ctx context.Context, in *Event, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/calendar.EventService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) Delete(ctx context.Context, in *EventIdentifier, opts ...grpc.CallOption) (*EventIdentifier, error) {
	out := new(EventIdentifier)
	err := c.cc.Invoke(ctx, "/calendar.EventService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) ReadList(ctx context.Context, in *EventListQuery, opts ...grpc.CallOption) (*EventListReply, error) {
	out := new(EventListReply)
	err := c.cc.Invoke(ctx, "/calendar.EventService/ReadList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventServiceServer is the server API for EventService service.
type EventServiceServer interface {
	Create(context.Context, *Event) (*Event, error)
	Read(context.Context, *EventIdentifier) (*Event, error)
	Update(context.Context, *Event) (*Event, error)
	Delete(context.Context, *EventIdentifier) (*EventIdentifier, error)
	ReadList(context.Context, *EventListQuery) (*EventListReply, error)
}

// UnimplementedEventServiceServer can be embedded to have forward compatible implementations.
type UnimplementedEventServiceServer struct {
}

func (*UnimplementedEventServiceServer) Create(ctx context.Context, req *Event) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedEventServiceServer) Read(ctx context.Context, req *EventIdentifier) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (*UnimplementedEventServiceServer) Update(ctx context.Context, req *Event) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (*UnimplementedEventServiceServer) Delete(ctx context.Context, req *EventIdentifier) (*EventIdentifier, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (*UnimplementedEventServiceServer) ReadList(ctx context.Context, req *EventListQuery) (*EventListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadList not implemented")
}

func RegisterEventServiceServer(s *grpc.Server, srv EventServiceServer) {
	s.RegisterService(&_EventService_serviceDesc, srv)
}

func _EventService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calendar.EventService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).Create(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventIdentifier)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calendar.EventService/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).Read(ctx, req.(*EventIdentifier))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calendar.EventService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).Update(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventIdentifier)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calendar.EventService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).Delete(ctx, req.(*EventIdentifier))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_ReadList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventListQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).ReadList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calendar.EventService/ReadList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).ReadList(ctx, req.(*EventListQuery))
	}
	return interceptor(ctx, in, info, handler)
}

var _EventService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "calendar.EventService",
	HandlerType: (*EventServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _EventService_Create_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _EventService_Read_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _EventService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _EventService_Delete_Handler,
		},
		{
			MethodName: "ReadList",
			Handler:    _EventService_ReadList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "app/proto/event.proto",
}

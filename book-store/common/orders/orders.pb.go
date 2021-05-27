// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.16.0
// source: orders/orders.proto

package orders

import (
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

type OrderID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id []byte `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *OrderID) Reset() {
	*x = OrderID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orders_orders_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderID) ProtoMessage() {}

func (x *OrderID) ProtoReflect() protoreflect.Message {
	mi := &file_orders_orders_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderID.ProtoReflect.Descriptor instead.
func (*OrderID) Descriptor() ([]byte, []int) {
	return file_orders_orders_proto_rawDescGZIP(), []int{0}
}

func (x *OrderID) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

type Order struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          []byte `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Book        []byte `protobuf:"bytes,3,opt,name=book,proto3" json:"book,omitempty"`
}

func (x *Order) Reset() {
	*x = Order{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orders_orders_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Order) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Order) ProtoMessage() {}

func (x *Order) ProtoReflect() protoreflect.Message {
	mi := &file_orders_orders_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Order.ProtoReflect.Descriptor instead.
func (*Order) Descriptor() ([]byte, []int) {
	return file_orders_orders_proto_rawDescGZIP(), []int{1}
}

func (x *Order) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Order) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Order) GetBook() []byte {
	if x != nil {
		return x.Book
	}
	return nil
}

type OrderCreateDTO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Description string `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	Book        []byte `protobuf:"bytes,2,opt,name=book,proto3" json:"book,omitempty"`
}

func (x *OrderCreateDTO) Reset() {
	*x = OrderCreateDTO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orders_orders_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderCreateDTO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderCreateDTO) ProtoMessage() {}

func (x *OrderCreateDTO) ProtoReflect() protoreflect.Message {
	mi := &file_orders_orders_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderCreateDTO.ProtoReflect.Descriptor instead.
func (*OrderCreateDTO) Descriptor() ([]byte, []int) {
	return file_orders_orders_proto_rawDescGZIP(), []int{2}
}

func (x *OrderCreateDTO) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *OrderCreateDTO) GetBook() []byte {
	if x != nil {
		return x.Book
	}
	return nil
}

var File_orders_orders_proto protoreflect.FileDescriptor

var file_orders_orders_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x22, 0x19, 0x0a,
	0x07, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x22, 0x4d, 0x0a, 0x05, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x6f, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x62, 0x6f, 0x6f, 0x6b, 0x22, 0x46, 0x0a, 0x0e, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x54, 0x4f, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x62,
	0x6f, 0x6f, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x6f, 0x6b, 0x32,
	0x6e, 0x0a, 0x06, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x12, 0x2c, 0x0a, 0x08, 0x47, 0x65, 0x74,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x0f, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x2e, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x49, 0x44, 0x1a, 0x0d, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x2e,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x2e,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x54, 0x4f, 0x1a, 0x0d,
	0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x22, 0x00, 0x42,
	0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x56, 0x65,
	0x73, 0x6e, 0x69, 0x6e, 0x6f, 0x76, 0x69, 0x63, 0x68, 0x2f, 0x67, 0x6f, 0x2d, 0x74, 0x61, 0x73,
	0x6b, 0x73, 0x2f, 0x62, 0x6f, 0x6f, 0x6b, 0x2d, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_orders_orders_proto_rawDescOnce sync.Once
	file_orders_orders_proto_rawDescData = file_orders_orders_proto_rawDesc
)

func file_orders_orders_proto_rawDescGZIP() []byte {
	file_orders_orders_proto_rawDescOnce.Do(func() {
		file_orders_orders_proto_rawDescData = protoimpl.X.CompressGZIP(file_orders_orders_proto_rawDescData)
	})
	return file_orders_orders_proto_rawDescData
}

var file_orders_orders_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_orders_orders_proto_goTypes = []interface{}{
	(*OrderID)(nil),        // 0: orders.OrderID
	(*Order)(nil),          // 1: orders.Order
	(*OrderCreateDTO)(nil), // 2: orders.OrderCreateDTO
}
var file_orders_orders_proto_depIdxs = []int32{
	0, // 0: orders.Orders.GetOrder:input_type -> orders.OrderID
	2, // 1: orders.Orders.CreateOrder:input_type -> orders.OrderCreateDTO
	1, // 2: orders.Orders.GetOrder:output_type -> orders.Order
	1, // 3: orders.Orders.CreateOrder:output_type -> orders.Order
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_orders_orders_proto_init() }
func file_orders_orders_proto_init() {
	if File_orders_orders_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_orders_orders_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OrderID); i {
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
		file_orders_orders_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Order); i {
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
		file_orders_orders_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OrderCreateDTO); i {
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
			RawDescriptor: file_orders_orders_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_orders_orders_proto_goTypes,
		DependencyIndexes: file_orders_orders_proto_depIdxs,
		MessageInfos:      file_orders_orders_proto_msgTypes,
	}.Build()
	File_orders_orders_proto = out.File
	file_orders_orders_proto_rawDesc = nil
	file_orders_orders_proto_goTypes = nil
	file_orders_orders_proto_depIdxs = nil
}
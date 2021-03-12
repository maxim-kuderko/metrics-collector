// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.3
// source: metrics.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
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

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string   `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Time   int64    `protobuf:"varint,2,opt,name=Time,proto3" json:"Time,omitempty"`
	Hash   uint64   `protobuf:"varint,3,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Values *Values  `protobuf:"bytes,4,opt,name=Values,proto3" json:"Values,omitempty"`
	Tags   []string `protobuf:"bytes,5,rep,name=Tags,proto3" json:"Tags,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metrics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_metrics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_metrics_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Metric) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *Metric) GetHash() uint64 {
	if x != nil {
		return x.Hash
	}
	return 0
}

func (x *Metric) GetValues() *Values {
	if x != nil {
		return x.Values
	}
	return nil
}

func (x *Metric) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type Values struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count int64   `protobuf:"varint,1,opt,name=Count,proto3" json:"Count,omitempty"`
	Sum   float64 `protobuf:"fixed64,2,opt,name=Sum,proto3" json:"Sum,omitempty"`
	Min   float64 `protobuf:"fixed64,3,opt,name=Min,proto3" json:"Min,omitempty"`
	Max   float64 `protobuf:"fixed64,4,opt,name=Max,proto3" json:"Max,omitempty"`
}

func (x *Values) Reset() {
	*x = Values{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metrics_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Values) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Values) ProtoMessage() {}

func (x *Values) ProtoReflect() protoreflect.Message {
	mi := &file_metrics_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Values.ProtoReflect.Descriptor instead.
func (*Values) Descriptor() ([]byte, []int) {
	return file_metrics_proto_rawDescGZIP(), []int{1}
}

func (x *Values) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *Values) GetSum() float64 {
	if x != nil {
		return x.Sum
	}
	return 0
}

func (x *Values) GetMin() float64 {
	if x != nil {
		return x.Min
	}
	return 0
}

func (x *Values) GetMax() float64 {
	if x != nil {
		return x.Max
	}
	return 0
}

var File_metrics_proto protoreflect.FileDescriptor

var file_metrics_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x79, 0x0a, 0x06,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x69,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x48, 0x61,
	0x73, 0x68, 0x12, 0x1f, 0x0a, 0x06, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x07, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x06, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x61, 0x67, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x04, 0x54, 0x61, 0x67, 0x73, 0x22, 0x54, 0x0a, 0x06, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x53, 0x75, 0x6d, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x53, 0x75, 0x6d, 0x12, 0x10, 0x0a, 0x03, 0x4d, 0x69, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x4d, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x4d,
	0x61, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x4d, 0x61, 0x78, 0x32, 0x45, 0x0a,
	0x14, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x47, 0x72, 0x70, 0x63, 0x12, 0x2d, 0x0a, 0x04, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x07, 0x2e,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00,
	0x28, 0x01, 0x30, 0x01, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6d, 0x61, 0x78, 0x69, 0x6d, 0x2d, 0x6b, 0x75, 0x64, 0x65, 0x72, 0x6b, 0x6f,
	0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2d, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_metrics_proto_rawDescOnce sync.Once
	file_metrics_proto_rawDescData = file_metrics_proto_rawDesc
)

func file_metrics_proto_rawDescGZIP() []byte {
	file_metrics_proto_rawDescOnce.Do(func() {
		file_metrics_proto_rawDescData = protoimpl.X.CompressGZIP(file_metrics_proto_rawDescData)
	})
	return file_metrics_proto_rawDescData
}

var file_metrics_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_metrics_proto_goTypes = []interface{}{
	(*Metric)(nil),      // 0: Metric
	(*Values)(nil),      // 1: Values
	(*empty.Empty)(nil), // 2: google.protobuf.Empty
}
var file_metrics_proto_depIdxs = []int32{
	1, // 0: Metric.Values:type_name -> Values
	0, // 1: MetricsCollectorGrpc.Send:input_type -> Metric
	2, // 2: MetricsCollectorGrpc.Send:output_type -> google.protobuf.Empty
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_metrics_proto_init() }
func file_metrics_proto_init() {
	if File_metrics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_metrics_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
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
		file_metrics_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Values); i {
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
			RawDescriptor: file_metrics_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_metrics_proto_goTypes,
		DependencyIndexes: file_metrics_proto_depIdxs,
		MessageInfos:      file_metrics_proto_msgTypes,
	}.Build()
	File_metrics_proto = out.File
	file_metrics_proto_rawDesc = nil
	file_metrics_proto_goTypes = nil
	file_metrics_proto_depIdxs = nil
}
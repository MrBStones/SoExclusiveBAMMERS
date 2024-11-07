// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: stc/mutex.proto

package mutex_proto

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

type AccessRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId    string `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Timestamp int64  `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *AccessRequest) Reset() {
	*x = AccessRequest{}
	mi := &file_stc_mutex_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AccessRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessRequest) ProtoMessage() {}

func (x *AccessRequest) ProtoReflect() protoreflect.Message {
	mi := &file_stc_mutex_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessRequest.ProtoReflect.Descriptor instead.
func (*AccessRequest) Descriptor() ([]byte, []int) {
	return file_stc_mutex_proto_rawDescGZIP(), []int{0}
}

func (x *AccessRequest) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *AccessRequest) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type AccessResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Granted bool `protobuf:"varint,1,opt,name=granted,proto3" json:"granted,omitempty"`
}

func (x *AccessResponse) Reset() {
	*x = AccessResponse{}
	mi := &file_stc_mutex_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AccessResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessResponse) ProtoMessage() {}

func (x *AccessResponse) ProtoReflect() protoreflect.Message {
	mi := &file_stc_mutex_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessResponse.ProtoReflect.Descriptor instead.
func (*AccessResponse) Descriptor() ([]byte, []int) {
	return file_stc_mutex_proto_rawDescGZIP(), []int{1}
}

func (x *AccessResponse) GetGranted() bool {
	if x != nil {
		return x.Granted
	}
	return false
}

type ReleaseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId string `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (x *ReleaseRequest) Reset() {
	*x = ReleaseRequest{}
	mi := &file_stc_mutex_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReleaseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReleaseRequest) ProtoMessage() {}

func (x *ReleaseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_stc_mutex_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReleaseRequest.ProtoReflect.Descriptor instead.
func (*ReleaseRequest) Descriptor() ([]byte, []int) {
	return file_stc_mutex_proto_rawDescGZIP(), []int{2}
}

func (x *ReleaseRequest) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

type ReleaseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Acknowledged bool `protobuf:"varint,1,opt,name=acknowledged,proto3" json:"acknowledged,omitempty"`
}

func (x *ReleaseResponse) Reset() {
	*x = ReleaseResponse{}
	mi := &file_stc_mutex_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReleaseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReleaseResponse) ProtoMessage() {}

func (x *ReleaseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_stc_mutex_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReleaseResponse.ProtoReflect.Descriptor instead.
func (*ReleaseResponse) Descriptor() ([]byte, []int) {
	return file_stc_mutex_proto_rawDescGZIP(), []int{3}
}

func (x *ReleaseResponse) GetAcknowledged() bool {
	if x != nil {
		return x.Acknowledged
	}
	return false
}

var File_stc_mutex_proto protoreflect.FileDescriptor

var file_stc_mutex_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x73, 0x74, 0x63, 0x2f, 0x6d, 0x75, 0x74, 0x65, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x46, 0x0a, 0x0d, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x2a, 0x0a, 0x0e, 0x41, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x67,
	0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x67, 0x72,
	0x61, 0x6e, 0x74, 0x65, 0x64, 0x22, 0x29, 0x0a, 0x0e, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64,
	0x22, 0x35, 0x0a, 0x0f, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x61, 0x63, 0x6b, 0x6e, 0x6f, 0x77, 0x6c, 0x65, 0x64,
	0x67, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x61, 0x63, 0x6b, 0x6e, 0x6f,
	0x77, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x64, 0x32, 0x78, 0x0a, 0x0c, 0x4d, 0x75, 0x74, 0x65, 0x78,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x0d, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x0e, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x0d, 0x52,
	0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x0f, 0x2e, 0x52,
	0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e,
	0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x0d, 0x5a, 0x0b, 0x6d, 0x75, 0x74, 0x65, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_stc_mutex_proto_rawDescOnce sync.Once
	file_stc_mutex_proto_rawDescData = file_stc_mutex_proto_rawDesc
)

func file_stc_mutex_proto_rawDescGZIP() []byte {
	file_stc_mutex_proto_rawDescOnce.Do(func() {
		file_stc_mutex_proto_rawDescData = protoimpl.X.CompressGZIP(file_stc_mutex_proto_rawDescData)
	})
	return file_stc_mutex_proto_rawDescData
}

var file_stc_mutex_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_stc_mutex_proto_goTypes = []any{
	(*AccessRequest)(nil),   // 0: AccessRequest
	(*AccessResponse)(nil),  // 1: AccessResponse
	(*ReleaseRequest)(nil),  // 2: ReleaseRequest
	(*ReleaseResponse)(nil), // 3: ReleaseResponse
}
var file_stc_mutex_proto_depIdxs = []int32{
	0, // 0: MutexService.RequestAccess:input_type -> AccessRequest
	2, // 1: MutexService.ReleaseAccess:input_type -> ReleaseRequest
	1, // 2: MutexService.RequestAccess:output_type -> AccessResponse
	3, // 3: MutexService.ReleaseAccess:output_type -> ReleaseResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_stc_mutex_proto_init() }
func file_stc_mutex_proto_init() {
	if File_stc_mutex_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_stc_mutex_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_stc_mutex_proto_goTypes,
		DependencyIndexes: file_stc_mutex_proto_depIdxs,
		MessageInfos:      file_stc_mutex_proto_msgTypes,
	}.Build()
	File_stc_mutex_proto = out.File
	file_stc_mutex_proto_rawDesc = nil
	file_stc_mutex_proto_goTypes = nil
	file_stc_mutex_proto_depIdxs = nil
}

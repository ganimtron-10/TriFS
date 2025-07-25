// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.12.4
// source: internal/proto/master.proto

package protocol

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetFileWorkersRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Filename      string                 `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFileWorkersRequest) Reset() {
	*x = GetFileWorkersRequest{}
	mi := &file_internal_proto_master_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFileWorkersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileWorkersRequest) ProtoMessage() {}

func (x *GetFileWorkersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_master_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileWorkersRequest.ProtoReflect.Descriptor instead.
func (*GetFileWorkersRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_master_proto_rawDescGZIP(), []int{0}
}

func (x *GetFileWorkersRequest) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

type GetFileWorkersResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	WorkerUrls    []string               `protobuf:"bytes,1,rep,name=worker_urls,json=workerUrls,proto3" json:"worker_urls,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFileWorkersResponse) Reset() {
	*x = GetFileWorkersResponse{}
	mi := &file_internal_proto_master_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFileWorkersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileWorkersResponse) ProtoMessage() {}

func (x *GetFileWorkersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_master_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileWorkersResponse.ProtoReflect.Descriptor instead.
func (*GetFileWorkersResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_master_proto_rawDescGZIP(), []int{1}
}

func (x *GetFileWorkersResponse) GetWorkerUrls() []string {
	if x != nil {
		return x.WorkerUrls
	}
	return nil
}

type AllocateFileWorkersRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Filename      string                 `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AllocateFileWorkersRequest) Reset() {
	*x = AllocateFileWorkersRequest{}
	mi := &file_internal_proto_master_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AllocateFileWorkersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllocateFileWorkersRequest) ProtoMessage() {}

func (x *AllocateFileWorkersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_master_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllocateFileWorkersRequest.ProtoReflect.Descriptor instead.
func (*AllocateFileWorkersRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_master_proto_rawDescGZIP(), []int{2}
}

func (x *AllocateFileWorkersRequest) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

type AllocateFileWorkersResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	WorkerUrls    []string               `protobuf:"bytes,1,rep,name=worker_urls,json=workerUrls,proto3" json:"worker_urls,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AllocateFileWorkersResponse) Reset() {
	*x = AllocateFileWorkersResponse{}
	mi := &file_internal_proto_master_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AllocateFileWorkersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllocateFileWorkersResponse) ProtoMessage() {}

func (x *AllocateFileWorkersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_master_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllocateFileWorkersResponse.ProtoReflect.Descriptor instead.
func (*AllocateFileWorkersResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_master_proto_rawDescGZIP(), []int{3}
}

func (x *AllocateFileWorkersResponse) GetWorkerUrls() []string {
	if x != nil {
		return x.WorkerUrls
	}
	return nil
}

type HeartbeatRequest struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	WorkerAddress    string                 `protobuf:"bytes,1,opt,name=worker_address,json=workerAddress,proto3" json:"worker_address,omitempty"`
	HostedFileHashes []string               `protobuf:"bytes,2,rep,name=hosted_file_hashes,json=hostedFileHashes,proto3" json:"hosted_file_hashes,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *HeartbeatRequest) Reset() {
	*x = HeartbeatRequest{}
	mi := &file_internal_proto_master_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HeartbeatRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatRequest) ProtoMessage() {}

func (x *HeartbeatRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_master_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartbeatRequest.ProtoReflect.Descriptor instead.
func (*HeartbeatRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_master_proto_rawDescGZIP(), []int{4}
}

func (x *HeartbeatRequest) GetWorkerAddress() string {
	if x != nil {
		return x.WorkerAddress
	}
	return ""
}

func (x *HeartbeatRequest) GetHostedFileHashes() []string {
	if x != nil {
		return x.HostedFileHashes
	}
	return nil
}

type HeartbeatResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HeartbeatResponse) Reset() {
	*x = HeartbeatResponse{}
	mi := &file_internal_proto_master_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HeartbeatResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatResponse) ProtoMessage() {}

func (x *HeartbeatResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_master_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartbeatResponse.ProtoReflect.Descriptor instead.
func (*HeartbeatResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_master_proto_rawDescGZIP(), []int{5}
}

var File_internal_proto_master_proto protoreflect.FileDescriptor

const file_internal_proto_master_proto_rawDesc = "" +
	"\n" +
	"\x1binternal/proto/master.proto\"3\n" +
	"\x15GetFileWorkersRequest\x12\x1a\n" +
	"\bfilename\x18\x01 \x01(\tR\bfilename\"9\n" +
	"\x16GetFileWorkersResponse\x12\x1f\n" +
	"\vworker_urls\x18\x01 \x03(\tR\n" +
	"workerUrls\"8\n" +
	"\x1aAllocateFileWorkersRequest\x12\x1a\n" +
	"\bfilename\x18\x01 \x01(\tR\bfilename\">\n" +
	"\x1bAllocateFileWorkersResponse\x12\x1f\n" +
	"\vworker_urls\x18\x01 \x03(\tR\n" +
	"workerUrls\"g\n" +
	"\x10HeartbeatRequest\x12%\n" +
	"\x0eworker_address\x18\x01 \x01(\tR\rworkerAddress\x12,\n" +
	"\x12hosted_file_hashes\x18\x02 \x03(\tR\x10hostedFileHashes\"\x13\n" +
	"\x11HeartbeatResponse2\xd8\x01\n" +
	"\rMasterService\x12A\n" +
	"\x0eGetFileWorkers\x12\x16.GetFileWorkersRequest\x1a\x17.GetFileWorkersResponse\x12P\n" +
	"\x13AllocateFileWorkers\x12\x1b.AllocateFileWorkersRequest\x1a\x1c.AllocateFileWorkersResponse\x122\n" +
	"\tHeartbeat\x12\x11.HeartbeatRequest\x1a\x12.HeartbeatResponseB\rZ\v./;protocolb\x06proto3"

var (
	file_internal_proto_master_proto_rawDescOnce sync.Once
	file_internal_proto_master_proto_rawDescData []byte
)

func file_internal_proto_master_proto_rawDescGZIP() []byte {
	file_internal_proto_master_proto_rawDescOnce.Do(func() {
		file_internal_proto_master_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_internal_proto_master_proto_rawDesc), len(file_internal_proto_master_proto_rawDesc)))
	})
	return file_internal_proto_master_proto_rawDescData
}

var file_internal_proto_master_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_internal_proto_master_proto_goTypes = []any{
	(*GetFileWorkersRequest)(nil),       // 0: GetFileWorkersRequest
	(*GetFileWorkersResponse)(nil),      // 1: GetFileWorkersResponse
	(*AllocateFileWorkersRequest)(nil),  // 2: AllocateFileWorkersRequest
	(*AllocateFileWorkersResponse)(nil), // 3: AllocateFileWorkersResponse
	(*HeartbeatRequest)(nil),            // 4: HeartbeatRequest
	(*HeartbeatResponse)(nil),           // 5: HeartbeatResponse
}
var file_internal_proto_master_proto_depIdxs = []int32{
	0, // 0: MasterService.GetFileWorkers:input_type -> GetFileWorkersRequest
	2, // 1: MasterService.AllocateFileWorkers:input_type -> AllocateFileWorkersRequest
	4, // 2: MasterService.Heartbeat:input_type -> HeartbeatRequest
	1, // 3: MasterService.GetFileWorkers:output_type -> GetFileWorkersResponse
	3, // 4: MasterService.AllocateFileWorkers:output_type -> AllocateFileWorkersResponse
	5, // 5: MasterService.Heartbeat:output_type -> HeartbeatResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_proto_master_proto_init() }
func file_internal_proto_master_proto_init() {
	if File_internal_proto_master_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_internal_proto_master_proto_rawDesc), len(file_internal_proto_master_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_proto_master_proto_goTypes,
		DependencyIndexes: file_internal_proto_master_proto_depIdxs,
		MessageInfos:      file_internal_proto_master_proto_msgTypes,
	}.Build()
	File_internal_proto_master_proto = out.File
	file_internal_proto_master_proto_goTypes = nil
	file_internal_proto_master_proto_depIdxs = nil
}

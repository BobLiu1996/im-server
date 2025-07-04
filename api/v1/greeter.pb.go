// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: greeter.proto

package v1

import (
	_ "github.com/google/gnostic/openapiv3"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// The request message containing the user's name.
type ListGreeterReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListGreeterReq) Reset() {
	*x = ListGreeterReq{}
	mi := &file_greeter_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListGreeterReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListGreeterReq) ProtoMessage() {}

func (x *ListGreeterReq) ProtoReflect() protoreflect.Message {
	mi := &file_greeter_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListGreeterReq.ProtoReflect.Descriptor instead.
func (*ListGreeterReq) Descriptor() ([]byte, []int) {
	return file_greeter_proto_rawDescGZIP(), []int{0}
}

// The response message containing the greetings
type ListGreeterRsp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ret           *BaseResp              `protobuf:"bytes,1,opt,name=ret,proto3" json:"ret,omitempty"`
	Body          *ListGreeterRsp_Body   `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListGreeterRsp) Reset() {
	*x = ListGreeterRsp{}
	mi := &file_greeter_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListGreeterRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListGreeterRsp) ProtoMessage() {}

func (x *ListGreeterRsp) ProtoReflect() protoreflect.Message {
	mi := &file_greeter_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListGreeterRsp.ProtoReflect.Descriptor instead.
func (*ListGreeterRsp) Descriptor() ([]byte, []int) {
	return file_greeter_proto_rawDescGZIP(), []int{1}
}

func (x *ListGreeterRsp) GetRet() *BaseResp {
	if x != nil {
		return x.Ret
	}
	return nil
}

func (x *ListGreeterRsp) GetBody() *ListGreeterRsp_Body {
	if x != nil {
		return x.Body
	}
	return nil
}

type Greeter struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"` // Greeter Name
	Age           uint32                 `protobuf:"varint,2,opt,name=age,proto3" json:"age,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Greeter) Reset() {
	*x = Greeter{}
	mi := &file_greeter_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Greeter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Greeter) ProtoMessage() {}

func (x *Greeter) ProtoReflect() protoreflect.Message {
	mi := &file_greeter_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Greeter.ProtoReflect.Descriptor instead.
func (*Greeter) Descriptor() ([]byte, []int) {
	return file_greeter_proto_rawDescGZIP(), []int{2}
}

func (x *Greeter) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Greeter) GetAge() uint32 {
	if x != nil {
		return x.Age
	}
	return 0
}

type ListGreeterRsp_Body struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Greeters      []*Greeter             `protobuf:"bytes,1,rep,name=greeters,proto3" json:"greeters,omitempty"` // Greeters
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListGreeterRsp_Body) Reset() {
	*x = ListGreeterRsp_Body{}
	mi := &file_greeter_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListGreeterRsp_Body) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListGreeterRsp_Body) ProtoMessage() {}

func (x *ListGreeterRsp_Body) ProtoReflect() protoreflect.Message {
	mi := &file_greeter_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListGreeterRsp_Body.ProtoReflect.Descriptor instead.
func (*ListGreeterRsp_Body) Descriptor() ([]byte, []int) {
	return file_greeter_proto_rawDescGZIP(), []int{1, 0}
}

func (x *ListGreeterRsp_Body) GetGreeters() []*Greeter {
	if x != nil {
		return x.Greeters
	}
	return nil
}

var File_greeter_proto protoreflect.FileDescriptor

const file_greeter_proto_rawDesc = "" +
	"\n" +
	"\rgreeter.proto\x12\tim_server\x1a\x1cgoogle/api/annotations.proto\x1a\x1copenapi/v3/annotations.proto\x1a\fcommon.proto\"\x10\n" +
	"\x0eListGreeterReq\"\xa3\x01\n" +
	"\x0eListGreeterRsp\x12%\n" +
	"\x03ret\x18\x01 \x01(\v2\x13.im_server.BaseRespR\x03ret\x122\n" +
	"\x04body\x18\x02 \x01(\v2\x1e.im_server.ListGreeterRsp.BodyR\x04body\x1a6\n" +
	"\x04Body\x12.\n" +
	"\bgreeters\x18\x01 \x03(\v2\x12.im_server.GreeterR\bgreeters\"/\n" +
	"\aGreeter\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x10\n" +
	"\x03age\x18\x02 \x01(\rR\x03age2\x86\x01\n" +
	"\n" +
	"GreeterSvc\x12x\n" +
	"\vListGreeter\x12\x19.im_server.ListGreeterReq\x1a\x19.im_server.ListGreeterRsp\"3\xbaG\x15\x12\x13获取Greeter列表\x82\xd3\xe4\x93\x02\x15:\x01*\"\x10/v1/greeter/listB\x15Z\x13im-server/api/v1;v1b\x06proto3"

var (
	file_greeter_proto_rawDescOnce sync.Once
	file_greeter_proto_rawDescData []byte
)

func file_greeter_proto_rawDescGZIP() []byte {
	file_greeter_proto_rawDescOnce.Do(func() {
		file_greeter_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_greeter_proto_rawDesc), len(file_greeter_proto_rawDesc)))
	})
	return file_greeter_proto_rawDescData
}

var file_greeter_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_greeter_proto_goTypes = []any{
	(*ListGreeterReq)(nil),      // 0: im_server.ListGreeterReq
	(*ListGreeterRsp)(nil),      // 1: im_server.ListGreeterRsp
	(*Greeter)(nil),             // 2: im_server.Greeter
	(*ListGreeterRsp_Body)(nil), // 3: im_server.ListGreeterRsp.Body
	(*BaseResp)(nil),            // 4: im_server.BaseResp
}
var file_greeter_proto_depIdxs = []int32{
	4, // 0: im_server.ListGreeterRsp.ret:type_name -> im_server.BaseResp
	3, // 1: im_server.ListGreeterRsp.body:type_name -> im_server.ListGreeterRsp.Body
	2, // 2: im_server.ListGreeterRsp.Body.greeters:type_name -> im_server.Greeter
	0, // 3: im_server.GreeterSvc.ListGreeter:input_type -> im_server.ListGreeterReq
	1, // 4: im_server.GreeterSvc.ListGreeter:output_type -> im_server.ListGreeterRsp
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_greeter_proto_init() }
func file_greeter_proto_init() {
	if File_greeter_proto != nil {
		return
	}
	file_common_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_greeter_proto_rawDesc), len(file_greeter_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_greeter_proto_goTypes,
		DependencyIndexes: file_greeter_proto_depIdxs,
		MessageInfos:      file_greeter_proto_msgTypes,
	}.Build()
	File_greeter_proto = out.File
	file_greeter_proto_goTypes = nil
	file_greeter_proto_depIdxs = nil
}

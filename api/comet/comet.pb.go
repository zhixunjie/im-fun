// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.7.1
// source: comet/comet.proto

package comet

import (
	protocol "github.com/zhixunjie/im-fun/api/protocol"
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

type PushMsgReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserKeys []string        `protobuf:"bytes,1,rep,name=userKeys,proto3" json:"userKeys,omitempty"`
	ProtoOp  int32           `protobuf:"varint,3,opt,name=protoOp,proto3" json:"protoOp,omitempty"`
	Proto    *protocol.Proto `protobuf:"bytes,2,opt,name=proto,proto3" json:"proto,omitempty"`
}

func (x *PushMsgReq) Reset() {
	*x = PushMsgReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_comet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushMsgReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushMsgReq) ProtoMessage() {}

func (x *PushMsgReq) ProtoReflect() protoreflect.Message {
	mi := &file_comet_comet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushMsgReq.ProtoReflect.Descriptor instead.
func (*PushMsgReq) Descriptor() ([]byte, []int) {
	return file_comet_comet_proto_rawDescGZIP(), []int{0}
}

func (x *PushMsgReq) GetUserKeys() []string {
	if x != nil {
		return x.UserKeys
	}
	return nil
}

func (x *PushMsgReq) GetProtoOp() int32 {
	if x != nil {
		return x.ProtoOp
	}
	return 0
}

func (x *PushMsgReq) GetProto() *protocol.Proto {
	if x != nil {
		return x.Proto
	}
	return nil
}

type PushMsgReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PushMsgReply) Reset() {
	*x = PushMsgReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_comet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushMsgReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushMsgReply) ProtoMessage() {}

func (x *PushMsgReply) ProtoReflect() protoreflect.Message {
	mi := &file_comet_comet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushMsgReply.ProtoReflect.Descriptor instead.
func (*PushMsgReply) Descriptor() ([]byte, []int) {
	return file_comet_comet_proto_rawDescGZIP(), []int{1}
}

type BroadcastReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProtoOp int32           `protobuf:"varint,1,opt,name=protoOp,proto3" json:"protoOp,omitempty"`
	Proto   *protocol.Proto `protobuf:"bytes,2,opt,name=proto,proto3" json:"proto,omitempty"`
	Speed   int32           `protobuf:"varint,3,opt,name=speed,proto3" json:"speed,omitempty"`
}

func (x *BroadcastReq) Reset() {
	*x = BroadcastReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_comet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastReq) ProtoMessage() {}

func (x *BroadcastReq) ProtoReflect() protoreflect.Message {
	mi := &file_comet_comet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastReq.ProtoReflect.Descriptor instead.
func (*BroadcastReq) Descriptor() ([]byte, []int) {
	return file_comet_comet_proto_rawDescGZIP(), []int{2}
}

func (x *BroadcastReq) GetProtoOp() int32 {
	if x != nil {
		return x.ProtoOp
	}
	return 0
}

func (x *BroadcastReq) GetProto() *protocol.Proto {
	if x != nil {
		return x.Proto
	}
	return nil
}

func (x *BroadcastReq) GetSpeed() int32 {
	if x != nil {
		return x.Speed
	}
	return 0
}

type BroadcastReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BroadcastReply) Reset() {
	*x = BroadcastReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_comet_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastReply) ProtoMessage() {}

func (x *BroadcastReply) ProtoReflect() protoreflect.Message {
	mi := &file_comet_comet_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastReply.ProtoReflect.Descriptor instead.
func (*BroadcastReply) Descriptor() ([]byte, []int) {
	return file_comet_comet_proto_rawDescGZIP(), []int{3}
}

type BroadcastRoomReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoomId string          `protobuf:"bytes,1,opt,name=roomId,proto3" json:"roomId,omitempty"`
	Proto  *protocol.Proto `protobuf:"bytes,2,opt,name=proto,proto3" json:"proto,omitempty"`
}

func (x *BroadcastRoomReq) Reset() {
	*x = BroadcastRoomReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_comet_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastRoomReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastRoomReq) ProtoMessage() {}

func (x *BroadcastRoomReq) ProtoReflect() protoreflect.Message {
	mi := &file_comet_comet_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastRoomReq.ProtoReflect.Descriptor instead.
func (*BroadcastRoomReq) Descriptor() ([]byte, []int) {
	return file_comet_comet_proto_rawDescGZIP(), []int{4}
}

func (x *BroadcastRoomReq) GetRoomId() string {
	if x != nil {
		return x.RoomId
	}
	return ""
}

func (x *BroadcastRoomReq) GetProto() *protocol.Proto {
	if x != nil {
		return x.Proto
	}
	return nil
}

type BroadcastRoomReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BroadcastRoomReply) Reset() {
	*x = BroadcastRoomReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_comet_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastRoomReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastRoomReply) ProtoMessage() {}

func (x *BroadcastRoomReply) ProtoReflect() protoreflect.Message {
	mi := &file_comet_comet_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastRoomReply.ProtoReflect.Descriptor instead.
func (*BroadcastRoomReply) Descriptor() ([]byte, []int) {
	return file_comet_comet_proto_rawDescGZIP(), []int{5}
}

type RoomsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RoomsReq) Reset() {
	*x = RoomsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_comet_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoomsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoomsReq) ProtoMessage() {}

func (x *RoomsReq) ProtoReflect() protoreflect.Message {
	mi := &file_comet_comet_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoomsReq.ProtoReflect.Descriptor instead.
func (*RoomsReq) Descriptor() ([]byte, []int) {
	return file_comet_comet_proto_rawDescGZIP(), []int{6}
}

type RoomsReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rooms map[string]bool `protobuf:"bytes,1,rep,name=rooms,proto3" json:"rooms,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *RoomsReply) Reset() {
	*x = RoomsReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_comet_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoomsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoomsReply) ProtoMessage() {}

func (x *RoomsReply) ProtoReflect() protoreflect.Message {
	mi := &file_comet_comet_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoomsReply.ProtoReflect.Descriptor instead.
func (*RoomsReply) Descriptor() ([]byte, []int) {
	return file_comet_comet_proto_rawDescGZIP(), []int{7}
}

func (x *RoomsReply) GetRooms() map[string]bool {
	if x != nil {
		return x.Rooms
	}
	return nil
}

var File_comet_comet_proto protoreflect.FileDescriptor

var file_comet_comet_proto_rawDesc = []byte{
	0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2f, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x69, 0x6d, 0x66, 0x75, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74,
	0x1a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6f, 0x0a, 0x0a, 0x50, 0x75, 0x73,
	0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4b,
	0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4b,
	0x65, 0x79, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x4f, 0x70, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x4f, 0x70, 0x12, 0x2b, 0x0a,
	0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x69,
	0x6d, 0x66, 0x75, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0e, 0x0a, 0x0c, 0x50, 0x75,
	0x73, 0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x6b, 0x0a, 0x0c, 0x42, 0x72,
	0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65, 0x71, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x4f, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x4f, 0x70, 0x12, 0x2b, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x69, 0x6d, 0x66, 0x75, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x70, 0x65, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x73, 0x70, 0x65, 0x65, 0x64, 0x22, 0x10, 0x0a, 0x0e, 0x42, 0x72, 0x6f, 0x61, 0x64,
	0x63, 0x61, 0x73, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x57, 0x0a, 0x10, 0x42, 0x72, 0x6f,
	0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a,
	0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72,
	0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x12, 0x2b, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x69, 0x6d, 0x66, 0x75, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x14, 0x0a, 0x12, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52,
	0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x0a, 0x0a, 0x08, 0x52, 0x6f, 0x6f, 0x6d,
	0x73, 0x52, 0x65, 0x71, 0x22, 0x80, 0x01, 0x0a, 0x0a, 0x52, 0x6f, 0x6f, 0x6d, 0x73, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x38, 0x0a, 0x05, 0x72, 0x6f, 0x6f, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x22, 0x2e, 0x69, 0x6d, 0x66, 0x75, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74,
	0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2e, 0x52, 0x6f, 0x6f, 0x6d,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x72, 0x6f, 0x6f, 0x6d, 0x73, 0x1a, 0x38, 0x0a,
	0x0a, 0x52, 0x6f, 0x6f, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x95, 0x02, 0x0a, 0x05, 0x43, 0x6f, 0x6d, 0x65,
	0x74, 0x12, 0x3d, 0x0a, 0x07, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x73, 0x67, 0x12, 0x17, 0x2e, 0x69,
	0x6d, 0x66, 0x75, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x4d,
	0x73, 0x67, 0x52, 0x65, 0x71, 0x1a, 0x19, 0x2e, 0x69, 0x6d, 0x66, 0x75, 0x6e, 0x2e, 0x63, 0x6f,
	0x6d, 0x65, 0x74, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x43, 0x0a, 0x09, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x12, 0x19, 0x2e,
	0x69, 0x6d, 0x66, 0x75, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x42, 0x72, 0x6f, 0x61,
	0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x1b, 0x2e, 0x69, 0x6d, 0x66, 0x75, 0x6e,
	0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x4f, 0x0a, 0x0d, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61,
	0x73, 0x74, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x1d, 0x2e, 0x69, 0x6d, 0x66, 0x75, 0x6e, 0x2e, 0x63,
	0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x6f,
	0x6f, 0x6d, 0x52, 0x65, 0x71, 0x1a, 0x1f, 0x2e, 0x69, 0x6d, 0x66, 0x75, 0x6e, 0x2e, 0x63, 0x6f,
	0x6d, 0x65, 0x74, 0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x6f, 0x6f,
	0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x37, 0x0a, 0x05, 0x52, 0x6f, 0x6f, 0x6d, 0x73, 0x12,
	0x15, 0x2e, 0x69, 0x6d, 0x66, 0x75, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x52, 0x6f,
	0x6f, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x69, 0x6d, 0x66, 0x75, 0x6e, 0x2e, 0x63,
	0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42,
	0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x7a, 0x68,
	0x69, 0x78, 0x75, 0x6e, 0x6a, 0x69, 0x65, 0x2f, 0x69, 0x6d, 0x2d, 0x66, 0x75, 0x6e, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x3b, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_comet_comet_proto_rawDescOnce sync.Once
	file_comet_comet_proto_rawDescData = file_comet_comet_proto_rawDesc
)

func file_comet_comet_proto_rawDescGZIP() []byte {
	file_comet_comet_proto_rawDescOnce.Do(func() {
		file_comet_comet_proto_rawDescData = protoimpl.X.CompressGZIP(file_comet_comet_proto_rawDescData)
	})
	return file_comet_comet_proto_rawDescData
}

var file_comet_comet_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_comet_comet_proto_goTypes = []interface{}{
	(*PushMsgReq)(nil),         // 0: imfun.comet.PushMsgReq
	(*PushMsgReply)(nil),       // 1: imfun.comet.PushMsgReply
	(*BroadcastReq)(nil),       // 2: imfun.comet.BroadcastReq
	(*BroadcastReply)(nil),     // 3: imfun.comet.BroadcastReply
	(*BroadcastRoomReq)(nil),   // 4: imfun.comet.BroadcastRoomReq
	(*BroadcastRoomReply)(nil), // 5: imfun.comet.BroadcastRoomReply
	(*RoomsReq)(nil),           // 6: imfun.comet.RoomsReq
	(*RoomsReply)(nil),         // 7: imfun.comet.RoomsReply
	nil,                        // 8: imfun.comet.RoomsReply.RoomsEntry
	(*protocol.Proto)(nil),     // 9: imfun.protocol.Proto
}
var file_comet_comet_proto_depIdxs = []int32{
	9, // 0: imfun.comet.PushMsgReq.proto:type_name -> imfun.protocol.Proto
	9, // 1: imfun.comet.BroadcastReq.proto:type_name -> imfun.protocol.Proto
	9, // 2: imfun.comet.BroadcastRoomReq.proto:type_name -> imfun.protocol.Proto
	8, // 3: imfun.comet.RoomsReply.rooms:type_name -> imfun.comet.RoomsReply.RoomsEntry
	0, // 4: imfun.comet.Comet.PushMsg:input_type -> imfun.comet.PushMsgReq
	2, // 5: imfun.comet.Comet.Broadcast:input_type -> imfun.comet.BroadcastReq
	4, // 6: imfun.comet.Comet.BroadcastRoom:input_type -> imfun.comet.BroadcastRoomReq
	6, // 7: imfun.comet.Comet.Rooms:input_type -> imfun.comet.RoomsReq
	1, // 8: imfun.comet.Comet.PushMsg:output_type -> imfun.comet.PushMsgReply
	3, // 9: imfun.comet.Comet.Broadcast:output_type -> imfun.comet.BroadcastReply
	5, // 10: imfun.comet.Comet.BroadcastRoom:output_type -> imfun.comet.BroadcastRoomReply
	7, // 11: imfun.comet.Comet.Rooms:output_type -> imfun.comet.RoomsReply
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_comet_comet_proto_init() }
func file_comet_comet_proto_init() {
	if File_comet_comet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_comet_comet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushMsgReq); i {
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
		file_comet_comet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushMsgReply); i {
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
		file_comet_comet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastReq); i {
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
		file_comet_comet_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastReply); i {
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
		file_comet_comet_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastRoomReq); i {
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
		file_comet_comet_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastRoomReply); i {
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
		file_comet_comet_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoomsReq); i {
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
		file_comet_comet_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoomsReply); i {
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
			RawDescriptor: file_comet_comet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_comet_comet_proto_goTypes,
		DependencyIndexes: file_comet_comet_proto_depIdxs,
		MessageInfos:      file_comet_comet_proto_msgTypes,
	}.Build()
	File_comet_comet_proto = out.File
	file_comet_comet_proto_rawDesc = nil
	file_comet_comet_proto_goTypes = nil
	file_comet_comet_proto_depIdxs = nil
}

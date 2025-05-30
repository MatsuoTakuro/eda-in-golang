// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: depotpb/api.proto

package depotpb

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

type CreateShoppingListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId string       `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	Items   []*OrderItem `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *CreateShoppingListRequest) Reset() {
	*x = CreateShoppingListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depotpb_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateShoppingListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateShoppingListRequest) ProtoMessage() {}

func (x *CreateShoppingListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_depotpb_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateShoppingListRequest.ProtoReflect.Descriptor instead.
func (*CreateShoppingListRequest) Descriptor() ([]byte, []int) {
	return file_depotpb_api_proto_rawDescGZIP(), []int{0}
}

func (x *CreateShoppingListRequest) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *CreateShoppingListRequest) GetItems() []*OrderItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type CreateShoppingListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *CreateShoppingListResponse) Reset() {
	*x = CreateShoppingListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depotpb_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateShoppingListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateShoppingListResponse) ProtoMessage() {}

func (x *CreateShoppingListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_depotpb_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateShoppingListResponse.ProtoReflect.Descriptor instead.
func (*CreateShoppingListResponse) Descriptor() ([]byte, []int) {
	return file_depotpb_api_proto_rawDescGZIP(), []int{1}
}

func (x *CreateShoppingListResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CancelShoppingListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *CancelShoppingListRequest) Reset() {
	*x = CancelShoppingListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depotpb_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CancelShoppingListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelShoppingListRequest) ProtoMessage() {}

func (x *CancelShoppingListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_depotpb_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelShoppingListRequest.ProtoReflect.Descriptor instead.
func (*CancelShoppingListRequest) Descriptor() ([]byte, []int) {
	return file_depotpb_api_proto_rawDescGZIP(), []int{2}
}

func (x *CancelShoppingListRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CancelShoppingListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CancelShoppingListResponse) Reset() {
	*x = CancelShoppingListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depotpb_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CancelShoppingListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelShoppingListResponse) ProtoMessage() {}

func (x *CancelShoppingListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_depotpb_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelShoppingListResponse.ProtoReflect.Descriptor instead.
func (*CancelShoppingListResponse) Descriptor() ([]byte, []int) {
	return file_depotpb_api_proto_rawDescGZIP(), []int{3}
}

type AssignShoppingListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	BotId string `protobuf:"bytes,2,opt,name=bot_id,json=botId,proto3" json:"bot_id,omitempty"`
}

func (x *AssignShoppingListRequest) Reset() {
	*x = AssignShoppingListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depotpb_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignShoppingListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignShoppingListRequest) ProtoMessage() {}

func (x *AssignShoppingListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_depotpb_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignShoppingListRequest.ProtoReflect.Descriptor instead.
func (*AssignShoppingListRequest) Descriptor() ([]byte, []int) {
	return file_depotpb_api_proto_rawDescGZIP(), []int{4}
}

func (x *AssignShoppingListRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AssignShoppingListRequest) GetBotId() string {
	if x != nil {
		return x.BotId
	}
	return ""
}

type AssignShoppingListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AssignShoppingListResponse) Reset() {
	*x = AssignShoppingListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depotpb_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignShoppingListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignShoppingListResponse) ProtoMessage() {}

func (x *AssignShoppingListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_depotpb_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignShoppingListResponse.ProtoReflect.Descriptor instead.
func (*AssignShoppingListResponse) Descriptor() ([]byte, []int) {
	return file_depotpb_api_proto_rawDescGZIP(), []int{5}
}

type CompleteShoppingListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *CompleteShoppingListRequest) Reset() {
	*x = CompleteShoppingListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depotpb_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompleteShoppingListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteShoppingListRequest) ProtoMessage() {}

func (x *CompleteShoppingListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_depotpb_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteShoppingListRequest.ProtoReflect.Descriptor instead.
func (*CompleteShoppingListRequest) Descriptor() ([]byte, []int) {
	return file_depotpb_api_proto_rawDescGZIP(), []int{6}
}

func (x *CompleteShoppingListRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CompleteShoppingListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CompleteShoppingListResponse) Reset() {
	*x = CompleteShoppingListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depotpb_api_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompleteShoppingListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompleteShoppingListResponse) ProtoMessage() {}

func (x *CompleteShoppingListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_depotpb_api_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompleteShoppingListResponse.ProtoReflect.Descriptor instead.
func (*CompleteShoppingListResponse) Descriptor() ([]byte, []int) {
	return file_depotpb_api_proto_rawDescGZIP(), []int{7}
}

var File_depotpb_api_proto protoreflect.FileDescriptor

var file_depotpb_api_proto_rawDesc = []byte{
	0x0a, 0x11, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x70, 0x62, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x07, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x70, 0x62, 0x1a, 0x16, 0x64, 0x65,
	0x70, 0x6f, 0x74, 0x70, 0x62, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x60, 0x0a, 0x19, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68,
	0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x05,
	0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x64, 0x65,
	0x70, 0x6f, 0x74, 0x70, 0x62, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x74, 0x65, 0x6d, 0x52,
	0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x2c, 0x0a, 0x1a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x2b, 0x0a, 0x19, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x53, 0x68,
	0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x22, 0x1c, 0x0a, 0x1a, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x53, 0x68, 0x6f, 0x70, 0x70,
	0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x42, 0x0a, 0x19, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e,
	0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x15, 0x0a, 0x06,
	0x62, 0x6f, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x62, 0x6f,
	0x74, 0x49, 0x64, 0x22, 0x1c, 0x0a, 0x1a, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x53, 0x68, 0x6f,
	0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x2d, 0x0a, 0x1b, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x68, 0x6f,
	0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x22, 0x1e, 0x0a, 0x1c, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x70,
	0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x32, 0x98, 0x03, 0x0a, 0x0c, 0x44, 0x65, 0x70, 0x6f, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x5f, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x70, 0x70,
	0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x22, 0x2e, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x70,
	0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x64, 0x65,
	0x70, 0x6f, 0x74, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x70,
	0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x5f, 0x0a, 0x12, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x53, 0x68, 0x6f, 0x70,
	0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x22, 0x2e, 0x64, 0x65, 0x70, 0x6f, 0x74,
	0x70, 0x62, 0x2e, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e,
	0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x64,
	0x65, 0x70, 0x6f, 0x74, 0x70, 0x62, 0x2e, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x53, 0x68, 0x6f,
	0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x5f, 0x0a, 0x12, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x53, 0x68, 0x6f,
	0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x22, 0x2e, 0x64, 0x65, 0x70, 0x6f,
	0x74, 0x70, 0x62, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x53, 0x68, 0x6f, 0x70, 0x70, 0x69,
	0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e,
	0x64, 0x65, 0x70, 0x6f, 0x74, 0x70, 0x62, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x53, 0x68,
	0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x65, 0x0a, 0x14, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65,
	0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x24, 0x2e, 0x64,
	0x65, 0x70, 0x6f, 0x74, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x53,
	0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x25, 0x2e, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6d,
	0x70, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x78, 0x0a, 0x0b, 0x63,
	0x6f, 0x6d, 0x2e, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x70, 0x62, 0x42, 0x08, 0x41, 0x70, 0x69, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x23, 0x65, 0x64, 0x61, 0x2d, 0x69, 0x6e, 0x2d, 0x67,
	0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x2f, 0x64, 0x65, 0x70, 0x6f,
	0x74, 0x70, 0x62, 0x2f, 0x64, 0x65, 0x70, 0x6f, 0x74, 0x70, 0x62, 0xa2, 0x02, 0x03, 0x44, 0x58,
	0x58, 0xaa, 0x02, 0x07, 0x44, 0x65, 0x70, 0x6f, 0x74, 0x70, 0x62, 0xca, 0x02, 0x07, 0x44, 0x65,
	0x70, 0x6f, 0x74, 0x70, 0x62, 0xe2, 0x02, 0x13, 0x44, 0x65, 0x70, 0x6f, 0x74, 0x70, 0x62, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x07, 0x44, 0x65,
	0x70, 0x6f, 0x74, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_depotpb_api_proto_rawDescOnce sync.Once
	file_depotpb_api_proto_rawDescData = file_depotpb_api_proto_rawDesc
)

func file_depotpb_api_proto_rawDescGZIP() []byte {
	file_depotpb_api_proto_rawDescOnce.Do(func() {
		file_depotpb_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_depotpb_api_proto_rawDescData)
	})
	return file_depotpb_api_proto_rawDescData
}

var file_depotpb_api_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_depotpb_api_proto_goTypes = []any{
	(*CreateShoppingListRequest)(nil),    // 0: depotpb.CreateShoppingListRequest
	(*CreateShoppingListResponse)(nil),   // 1: depotpb.CreateShoppingListResponse
	(*CancelShoppingListRequest)(nil),    // 2: depotpb.CancelShoppingListRequest
	(*CancelShoppingListResponse)(nil),   // 3: depotpb.CancelShoppingListResponse
	(*AssignShoppingListRequest)(nil),    // 4: depotpb.AssignShoppingListRequest
	(*AssignShoppingListResponse)(nil),   // 5: depotpb.AssignShoppingListResponse
	(*CompleteShoppingListRequest)(nil),  // 6: depotpb.CompleteShoppingListRequest
	(*CompleteShoppingListResponse)(nil), // 7: depotpb.CompleteShoppingListResponse
	(*OrderItem)(nil),                    // 8: depotpb.OrderItem
}
var file_depotpb_api_proto_depIdxs = []int32{
	8, // 0: depotpb.CreateShoppingListRequest.items:type_name -> depotpb.OrderItem
	0, // 1: depotpb.DepotService.CreateShoppingList:input_type -> depotpb.CreateShoppingListRequest
	2, // 2: depotpb.DepotService.CancelShoppingList:input_type -> depotpb.CancelShoppingListRequest
	4, // 3: depotpb.DepotService.AssignShoppingList:input_type -> depotpb.AssignShoppingListRequest
	6, // 4: depotpb.DepotService.CompleteShoppingList:input_type -> depotpb.CompleteShoppingListRequest
	1, // 5: depotpb.DepotService.CreateShoppingList:output_type -> depotpb.CreateShoppingListResponse
	3, // 6: depotpb.DepotService.CancelShoppingList:output_type -> depotpb.CancelShoppingListResponse
	5, // 7: depotpb.DepotService.AssignShoppingList:output_type -> depotpb.AssignShoppingListResponse
	7, // 8: depotpb.DepotService.CompleteShoppingList:output_type -> depotpb.CompleteShoppingListResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_depotpb_api_proto_init() }
func file_depotpb_api_proto_init() {
	if File_depotpb_api_proto != nil {
		return
	}
	file_depotpb_messages_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_depotpb_api_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*CreateShoppingListRequest); i {
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
		file_depotpb_api_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CreateShoppingListResponse); i {
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
		file_depotpb_api_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*CancelShoppingListRequest); i {
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
		file_depotpb_api_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*CancelShoppingListResponse); i {
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
		file_depotpb_api_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*AssignShoppingListRequest); i {
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
		file_depotpb_api_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*AssignShoppingListResponse); i {
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
		file_depotpb_api_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*CompleteShoppingListRequest); i {
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
		file_depotpb_api_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*CompleteShoppingListResponse); i {
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
			RawDescriptor: file_depotpb_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_depotpb_api_proto_goTypes,
		DependencyIndexes: file_depotpb_api_proto_depIdxs,
		MessageInfos:      file_depotpb_api_proto_msgTypes,
	}.Build()
	File_depotpb_api_proto = out.File
	file_depotpb_api_proto_rawDesc = nil
	file_depotpb_api_proto_goTypes = nil
	file_depotpb_api_proto_depIdxs = nil
}

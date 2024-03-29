// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: protoc/paxos.proto

package protoc

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

type LogEntryStatus int32

const (
	LogEntryStatus_Idle    LogEntryStatus = 0
	LogEntryStatus_Running LogEntryStatus = 1
	LogEntryStatus_Chosen  LogEntryStatus = 2
)

// Enum value maps for LogEntryStatus.
var (
	LogEntryStatus_name = map[int32]string{
		0: "Idle",
		1: "Running",
		2: "Chosen",
	}
	LogEntryStatus_value = map[string]int32{
		"Idle":    0,
		"Running": 1,
		"Chosen":  2,
	}
)

func (x LogEntryStatus) Enum() *LogEntryStatus {
	p := new(LogEntryStatus)
	*p = x
	return p
}

func (x LogEntryStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LogEntryStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_protoc_paxos_proto_enumTypes[0].Descriptor()
}

func (LogEntryStatus) Type() protoreflect.EnumType {
	return &file_protoc_paxos_proto_enumTypes[0]
}

func (x LogEntryStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LogEntryStatus.Descriptor instead.
func (LogEntryStatus) EnumDescriptor() ([]byte, []int) {
	return file_protoc_paxos_proto_rawDescGZIP(), []int{0}
}

type PrepareRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index       int64 `protobuf:"varint,1,opt,name=Index,proto3" json:"Index,omitempty"`
	ProposalNum int64 `protobuf:"varint,2,opt,name=ProposalNum,proto3" json:"ProposalNum,omitempty"`
}

func (x *PrepareRequest) Reset() {
	*x = PrepareRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protoc_paxos_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrepareRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareRequest) ProtoMessage() {}

func (x *PrepareRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protoc_paxos_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareRequest.ProtoReflect.Descriptor instead.
func (*PrepareRequest) Descriptor() ([]byte, []int) {
	return file_protoc_paxos_proto_rawDescGZIP(), []int{0}
}

func (x *PrepareRequest) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *PrepareRequest) GetProposalNum() int64 {
	if x != nil {
		return x.ProposalNum
	}
	return 0
}

type PrepareReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AcceptedProposal int64  `protobuf:"varint,1,opt,name=AcceptedProposal,proto3" json:"AcceptedProposal,omitempty"`
	AcceptedValue    string `protobuf:"bytes,2,opt,name=AcceptedValue,proto3" json:"AcceptedValue,omitempty"`
	From             string `protobuf:"bytes,3,opt,name=From,proto3" json:"From,omitempty"`
}

func (x *PrepareReply) Reset() {
	*x = PrepareReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protoc_paxos_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrepareReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareReply) ProtoMessage() {}

func (x *PrepareReply) ProtoReflect() protoreflect.Message {
	mi := &file_protoc_paxos_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareReply.ProtoReflect.Descriptor instead.
func (*PrepareReply) Descriptor() ([]byte, []int) {
	return file_protoc_paxos_proto_rawDescGZIP(), []int{1}
}

func (x *PrepareReply) GetAcceptedProposal() int64 {
	if x != nil {
		return x.AcceptedProposal
	}
	return 0
}

func (x *PrepareReply) GetAcceptedValue() string {
	if x != nil {
		return x.AcceptedValue
	}
	return ""
}

func (x *PrepareReply) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

type AcceptRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index         int64  `protobuf:"varint,1,opt,name=Index,proto3" json:"Index,omitempty"`
	ProposalNum   int64  `protobuf:"varint,2,opt,name=ProposalNum,proto3" json:"ProposalNum,omitempty"`
	ProposalValue string `protobuf:"bytes,3,opt,name=ProposalValue,proto3" json:"ProposalValue,omitempty"`
}

func (x *AcceptRequest) Reset() {
	*x = AcceptRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protoc_paxos_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptRequest) ProtoMessage() {}

func (x *AcceptRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protoc_paxos_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptRequest.ProtoReflect.Descriptor instead.
func (*AcceptRequest) Descriptor() ([]byte, []int) {
	return file_protoc_paxos_proto_rawDescGZIP(), []int{2}
}

func (x *AcceptRequest) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *AcceptRequest) GetProposalNum() int64 {
	if x != nil {
		return x.ProposalNum
	}
	return 0
}

func (x *AcceptRequest) GetProposalValue() string {
	if x != nil {
		return x.ProposalValue
	}
	return ""
}

type AcceptReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MiniProposal int64  `protobuf:"varint,1,opt,name=MiniProposal,proto3" json:"MiniProposal,omitempty"`
	From         string `protobuf:"bytes,2,opt,name=From,proto3" json:"From,omitempty"`
}

func (x *AcceptReply) Reset() {
	*x = AcceptReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protoc_paxos_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptReply) ProtoMessage() {}

func (x *AcceptReply) ProtoReflect() protoreflect.Message {
	mi := &file_protoc_paxos_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptReply.ProtoReflect.Descriptor instead.
func (*AcceptReply) Descriptor() ([]byte, []int) {
	return file_protoc_paxos_proto_rawDescGZIP(), []int{3}
}

func (x *AcceptReply) GetMiniProposal() int64 {
	if x != nil {
		return x.MiniProposal
	}
	return 0
}

func (x *AcceptReply) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

type LearnRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index         int64  `protobuf:"varint,1,opt,name=Index,proto3" json:"Index,omitempty"`
	ProposalNum   int64  `protobuf:"varint,2,opt,name=ProposalNum,proto3" json:"ProposalNum,omitempty"`
	ProposalValue string `protobuf:"bytes,3,opt,name=ProposalValue,proto3" json:"ProposalValue,omitempty"`
	Proposer      string `protobuf:"bytes,4,opt,name=Proposer,proto3" json:"Proposer,omitempty"`
}

func (x *LearnRequest) Reset() {
	*x = LearnRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protoc_paxos_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LearnRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LearnRequest) ProtoMessage() {}

func (x *LearnRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protoc_paxos_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LearnRequest.ProtoReflect.Descriptor instead.
func (*LearnRequest) Descriptor() ([]byte, []int) {
	return file_protoc_paxos_proto_rawDescGZIP(), []int{4}
}

func (x *LearnRequest) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *LearnRequest) GetProposalNum() int64 {
	if x != nil {
		return x.ProposalNum
	}
	return 0
}

func (x *LearnRequest) GetProposalValue() string {
	if x != nil {
		return x.ProposalValue
	}
	return ""
}

func (x *LearnRequest) GetProposer() string {
	if x != nil {
		return x.Proposer
	}
	return ""
}

type LearnReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LearnReply) Reset() {
	*x = LearnReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protoc_paxos_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LearnReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LearnReply) ProtoMessage() {}

func (x *LearnReply) ProtoReflect() protoreflect.Message {
	mi := &file_protoc_paxos_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LearnReply.ProtoReflect.Descriptor instead.
func (*LearnReply) Descriptor() ([]byte, []int) {
	return file_protoc_paxos_proto_rawDescGZIP(), []int{5}
}

type ProposeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProposalValue string `protobuf:"bytes,1,opt,name=ProposalValue,proto3" json:"ProposalValue,omitempty"`
}

func (x *ProposeRequest) Reset() {
	*x = ProposeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protoc_paxos_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProposeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposeRequest) ProtoMessage() {}

func (x *ProposeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protoc_paxos_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposeRequest.ProtoReflect.Descriptor instead.
func (*ProposeRequest) Descriptor() ([]byte, []int) {
	return file_protoc_paxos_proto_rawDescGZIP(), []int{6}
}

func (x *ProposeRequest) GetProposalValue() string {
	if x != nil {
		return x.ProposalValue
	}
	return ""
}

type ProposeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ProposeReply) Reset() {
	*x = ProposeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protoc_paxos_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProposeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposeReply) ProtoMessage() {}

func (x *ProposeReply) ProtoReflect() protoreflect.Message {
	mi := &file_protoc_paxos_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposeReply.ProtoReflect.Descriptor instead.
func (*ProposeReply) Descriptor() ([]byte, []int) {
	return file_protoc_paxos_proto_rawDescGZIP(), []int{7}
}

type LogEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MiniProposal     int64          `protobuf:"varint,1,opt,name=MiniProposal,proto3" json:"MiniProposal,omitempty"`
	AcceptedProposal int64          `protobuf:"varint,2,opt,name=AcceptedProposal,proto3" json:"AcceptedProposal,omitempty"`
	AcceptedValue    string         `protobuf:"bytes,3,opt,name=AcceptedValue,proto3" json:"AcceptedValue,omitempty"`
	Status           LogEntryStatus `protobuf:"varint,4,opt,name=Status,proto3,enum=protoc.LogEntryStatus" json:"Status,omitempty"`
}

func (x *LogEntry) Reset() {
	*x = LogEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protoc_paxos_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogEntry) ProtoMessage() {}

func (x *LogEntry) ProtoReflect() protoreflect.Message {
	mi := &file_protoc_paxos_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogEntry.ProtoReflect.Descriptor instead.
func (*LogEntry) Descriptor() ([]byte, []int) {
	return file_protoc_paxos_proto_rawDescGZIP(), []int{8}
}

func (x *LogEntry) GetMiniProposal() int64 {
	if x != nil {
		return x.MiniProposal
	}
	return 0
}

func (x *LogEntry) GetAcceptedProposal() int64 {
	if x != nil {
		return x.AcceptedProposal
	}
	return 0
}

func (x *LogEntry) GetAcceptedValue() string {
	if x != nil {
		return x.AcceptedValue
	}
	return ""
}

func (x *LogEntry) GetStatus() LogEntryStatus {
	if x != nil {
		return x.Status
	}
	return LogEntryStatus_Idle
}

var File_protoc_paxos_proto protoreflect.FileDescriptor

var file_protoc_paxos_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2f, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x22, 0x48, 0x0a, 0x0e,
	0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c,
	0x4e, 0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x50, 0x72, 0x6f, 0x70, 0x6f,
	0x73, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x22, 0x74, 0x0a, 0x0c, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72,
	0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x2a, 0x0a, 0x10, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74,
	0x65, 0x64, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x10, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73,
	0x61, 0x6c, 0x12, 0x24, 0x0a, 0x0d, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x41, 0x63, 0x63, 0x65, 0x70,
	0x74, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x22, 0x6d, 0x0a, 0x0d,
	0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x4e,
	0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73,
	0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x12, 0x24, 0x0a, 0x0d, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61,
	0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x50, 0x72,
	0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x45, 0x0a, 0x0b, 0x41,
	0x63, 0x63, 0x65, 0x70, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x22, 0x0a, 0x0c, 0x4d, 0x69,
	0x6e, 0x69, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0c, 0x4d, 0x69, 0x6e, 0x69, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x12, 0x12,
	0x0a, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x46, 0x72,
	0x6f, 0x6d, 0x22, 0x88, 0x01, 0x0a, 0x0c, 0x4c, 0x65, 0x61, 0x72, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x72, 0x6f,
	0x70, 0x6f, 0x73, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b,
	0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x12, 0x24, 0x0a, 0x0d, 0x50,
	0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x22, 0x0c, 0x0a,
	0x0a, 0x4c, 0x65, 0x61, 0x72, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x36, 0x0a, 0x0e, 0x50,
	0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a,
	0x0d, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x22, 0x0e, 0x0a, 0x0c, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x22, 0xb0, 0x01, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x22, 0x0a, 0x0c, 0x4d, 0x69, 0x6e, 0x69, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x4d, 0x69, 0x6e, 0x69, 0x50, 0x72, 0x6f, 0x70,
	0x6f, 0x73, 0x61, 0x6c, 0x12, 0x2a, 0x0a, 0x10, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64,
	0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10,
	0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c,
	0x12, 0x24, 0x0a, 0x0d, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65,
	0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2e, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2e,
	0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2a, 0x33, 0x0a, 0x0e, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x64, 0x6c, 0x65,
	0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x10, 0x01, 0x12,
	0x0a, 0x0a, 0x06, 0x43, 0x68, 0x6f, 0x73, 0x65, 0x6e, 0x10, 0x02, 0x32, 0x7d, 0x0a, 0x08, 0x41,
	0x63, 0x63, 0x65, 0x70, 0x74, 0x6f, 0x72, 0x12, 0x39, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x70, 0x61,
	0x72, 0x65, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2e, 0x50, 0x72, 0x65, 0x70,
	0x61, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x22, 0x00, 0x12, 0x36, 0x0a, 0x06, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x12, 0x15, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2e, 0x41, 0x63, 0x63,
	0x65, 0x70, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x32, 0x45, 0x0a, 0x08, 0x50, 0x72,
	0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x12, 0x39, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73,
	0x65, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x6f,
	0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22,
	0x00, 0x32, 0x3e, 0x0a, 0x07, 0x4c, 0x65, 0x61, 0x72, 0x6e, 0x65, 0x72, 0x12, 0x33, 0x0a, 0x05,
	0x4c, 0x65, 0x61, 0x72, 0x6e, 0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2e, 0x4c,
	0x65, 0x61, 0x72, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x2e, 0x4c, 0x65, 0x61, 0x72, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22,
	0x00, 0x42, 0x20, 0x5a, 0x1e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x66, 0x75, 0x64, 0x75, 0x74, 0x65, 0x2f, 0x70, 0x61, 0x78, 0x6f, 0x73, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protoc_paxos_proto_rawDescOnce sync.Once
	file_protoc_paxos_proto_rawDescData = file_protoc_paxos_proto_rawDesc
)

func file_protoc_paxos_proto_rawDescGZIP() []byte {
	file_protoc_paxos_proto_rawDescOnce.Do(func() {
		file_protoc_paxos_proto_rawDescData = protoimpl.X.CompressGZIP(file_protoc_paxos_proto_rawDescData)
	})
	return file_protoc_paxos_proto_rawDescData
}

var file_protoc_paxos_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_protoc_paxos_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_protoc_paxos_proto_goTypes = []interface{}{
	(LogEntryStatus)(0),    // 0: protoc.LogEntryStatus
	(*PrepareRequest)(nil), // 1: protoc.PrepareRequest
	(*PrepareReply)(nil),   // 2: protoc.PrepareReply
	(*AcceptRequest)(nil),  // 3: protoc.AcceptRequest
	(*AcceptReply)(nil),    // 4: protoc.AcceptReply
	(*LearnRequest)(nil),   // 5: protoc.LearnRequest
	(*LearnReply)(nil),     // 6: protoc.LearnReply
	(*ProposeRequest)(nil), // 7: protoc.ProposeRequest
	(*ProposeReply)(nil),   // 8: protoc.ProposeReply
	(*LogEntry)(nil),       // 9: protoc.LogEntry
}
var file_protoc_paxos_proto_depIdxs = []int32{
	0, // 0: protoc.LogEntry.Status:type_name -> protoc.LogEntryStatus
	1, // 1: protoc.Acceptor.Prepare:input_type -> protoc.PrepareRequest
	3, // 2: protoc.Acceptor.Accept:input_type -> protoc.AcceptRequest
	7, // 3: protoc.Proposer.Propose:input_type -> protoc.ProposeRequest
	5, // 4: protoc.Learner.Learn:input_type -> protoc.LearnRequest
	2, // 5: protoc.Acceptor.Prepare:output_type -> protoc.PrepareReply
	4, // 6: protoc.Acceptor.Accept:output_type -> protoc.AcceptReply
	8, // 7: protoc.Proposer.Propose:output_type -> protoc.ProposeReply
	6, // 8: protoc.Learner.Learn:output_type -> protoc.LearnReply
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_protoc_paxos_proto_init() }
func file_protoc_paxos_proto_init() {
	if File_protoc_paxos_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protoc_paxos_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrepareRequest); i {
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
		file_protoc_paxos_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrepareReply); i {
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
		file_protoc_paxos_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptRequest); i {
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
		file_protoc_paxos_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptReply); i {
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
		file_protoc_paxos_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LearnRequest); i {
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
		file_protoc_paxos_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LearnReply); i {
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
		file_protoc_paxos_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProposeRequest); i {
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
		file_protoc_paxos_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProposeReply); i {
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
		file_protoc_paxos_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogEntry); i {
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
			RawDescriptor: file_protoc_paxos_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   3,
		},
		GoTypes:           file_protoc_paxos_proto_goTypes,
		DependencyIndexes: file_protoc_paxos_proto_depIdxs,
		EnumInfos:         file_protoc_paxos_proto_enumTypes,
		MessageInfos:      file_protoc_paxos_proto_msgTypes,
	}.Build()
	File_protoc_paxos_proto = out.File
	file_protoc_paxos_proto_rawDesc = nil
	file_protoc_paxos_proto_goTypes = nil
	file_protoc_paxos_proto_depIdxs = nil
}

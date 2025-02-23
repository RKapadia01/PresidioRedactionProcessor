// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: presidio.proto

package presidioclient

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Represents a request to analyze text for PII entities.
type PresidioAnalyzerRequest struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Text           string                 `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Language       string                 `protobuf:"bytes,2,opt,name=language,proto3" json:"language,omitempty"`
	ScoreThreshold float64                `protobuf:"fixed64,3,opt,name=score_threshold,json=scoreThreshold,proto3" json:"score_threshold,omitempty"`
	Entities       []string               `protobuf:"bytes,4,rep,name=entities,proto3" json:"entities,omitempty"`
	Context        []string               `protobuf:"bytes,5,rep,name=context,proto3" json:"context,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *PresidioAnalyzerRequest) Reset() {
	*x = PresidioAnalyzerRequest{}
	mi := &file_presidio_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PresidioAnalyzerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PresidioAnalyzerRequest) ProtoMessage() {}

func (x *PresidioAnalyzerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_presidio_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PresidioAnalyzerRequest.ProtoReflect.Descriptor instead.
func (*PresidioAnalyzerRequest) Descriptor() ([]byte, []int) {
	return file_presidio_proto_rawDescGZIP(), []int{0}
}

func (x *PresidioAnalyzerRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *PresidioAnalyzerRequest) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

func (x *PresidioAnalyzerRequest) GetScoreThreshold() float64 {
	if x != nil {
		return x.ScoreThreshold
	}
	return 0
}

func (x *PresidioAnalyzerRequest) GetEntities() []string {
	if x != nil {
		return x.Entities
	}
	return nil
}

func (x *PresidioAnalyzerRequest) GetContext() []string {
	if x != nil {
		return x.Context
	}
	return nil
}

// Represents the outcome of an analysis operation.
type PresidioAnalyzerResponses struct {
	state           protoimpl.MessageState      `protogen:"open.v1"`
	AnalyzerResults []*PresidioAnalyzerResponse `protobuf:"bytes,1,rep,name=analyzer_results,json=analyzerResults,proto3" json:"analyzer_results,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *PresidioAnalyzerResponses) Reset() {
	*x = PresidioAnalyzerResponses{}
	mi := &file_presidio_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PresidioAnalyzerResponses) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PresidioAnalyzerResponses) ProtoMessage() {}

func (x *PresidioAnalyzerResponses) ProtoReflect() protoreflect.Message {
	mi := &file_presidio_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PresidioAnalyzerResponses.ProtoReflect.Descriptor instead.
func (*PresidioAnalyzerResponses) Descriptor() ([]byte, []int) {
	return file_presidio_proto_rawDescGZIP(), []int{1}
}

func (x *PresidioAnalyzerResponses) GetAnalyzerResults() []*PresidioAnalyzerResponse {
	if x != nil {
		return x.AnalyzerResults
	}
	return nil
}

// Represents the individual outcome of an analysis operation.
type PresidioAnalyzerResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Start         int32                  `protobuf:"varint,1,opt,name=start,proto3" json:"start,omitempty"`
	End           int32                  `protobuf:"varint,2,opt,name=end,proto3" json:"end,omitempty"`
	Score         float64                `protobuf:"fixed64,3,opt,name=score,proto3" json:"score,omitempty"`
	EntityType    string                 `protobuf:"bytes,4,opt,name=entity_type,json=entityType,proto3" json:"entity_type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PresidioAnalyzerResponse) Reset() {
	*x = PresidioAnalyzerResponse{}
	mi := &file_presidio_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PresidioAnalyzerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PresidioAnalyzerResponse) ProtoMessage() {}

func (x *PresidioAnalyzerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_presidio_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PresidioAnalyzerResponse.ProtoReflect.Descriptor instead.
func (*PresidioAnalyzerResponse) Descriptor() ([]byte, []int) {
	return file_presidio_proto_rawDescGZIP(), []int{2}
}

func (x *PresidioAnalyzerResponse) GetStart() int32 {
	if x != nil {
		return x.Start
	}
	return 0
}

func (x *PresidioAnalyzerResponse) GetEnd() int32 {
	if x != nil {
		return x.End
	}
	return 0
}

func (x *PresidioAnalyzerResponse) GetScore() float64 {
	if x != nil {
		return x.Score
	}
	return 0
}

func (x *PresidioAnalyzerResponse) GetEntityType() string {
	if x != nil {
		return x.EntityType
	}
	return ""
}

// Represents a request to anonymize text.
type PresidioAnonymizerRequest struct {
	state           protoimpl.MessageState         `protogen:"open.v1"`
	Text            string                         `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Anonymizers     map[string]*PresidioAnonymizer `protobuf:"bytes,2,rep,name=anonymizers,proto3" json:"anonymizers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	AnalyzerResults []*PresidioAnalyzerResponse    `protobuf:"bytes,3,rep,name=analyzer_results,json=analyzerResults,proto3" json:"analyzer_results,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *PresidioAnonymizerRequest) Reset() {
	*x = PresidioAnonymizerRequest{}
	mi := &file_presidio_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PresidioAnonymizerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PresidioAnonymizerRequest) ProtoMessage() {}

func (x *PresidioAnonymizerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_presidio_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PresidioAnonymizerRequest.ProtoReflect.Descriptor instead.
func (*PresidioAnonymizerRequest) Descriptor() ([]byte, []int) {
	return file_presidio_proto_rawDescGZIP(), []int{3}
}

func (x *PresidioAnonymizerRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *PresidioAnonymizerRequest) GetAnonymizers() map[string]*PresidioAnonymizer {
	if x != nil {
		return x.Anonymizers
	}
	return nil
}

func (x *PresidioAnonymizerRequest) GetAnalyzerResults() []*PresidioAnalyzerResponse {
	if x != nil {
		return x.AnalyzerResults
	}
	return nil
}

// Defines the configuration of a particular anonymizer.
type PresidioAnonymizer struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Type          string                 `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	NewValue      string                 `protobuf:"bytes,2,opt,name=new_value,json=newValue,proto3" json:"new_value,omitempty"`
	MaskingChar   string                 `protobuf:"bytes,3,opt,name=masking_char,json=maskingChar,proto3" json:"masking_char,omitempty"`
	CharsToMask   int32                  `protobuf:"varint,4,opt,name=chars_to_mask,json=charsToMask,proto3" json:"chars_to_mask,omitempty"`
	FromEnd       bool                   `protobuf:"varint,5,opt,name=from_end,json=fromEnd,proto3" json:"from_end,omitempty"`
	HashType      string                 `protobuf:"bytes,6,opt,name=hash_type,json=hashType,proto3" json:"hash_type,omitempty"`
	Key           string                 `protobuf:"bytes,7,opt,name=key,proto3" json:"key,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PresidioAnonymizer) Reset() {
	*x = PresidioAnonymizer{}
	mi := &file_presidio_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PresidioAnonymizer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PresidioAnonymizer) ProtoMessage() {}

func (x *PresidioAnonymizer) ProtoReflect() protoreflect.Message {
	mi := &file_presidio_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PresidioAnonymizer.ProtoReflect.Descriptor instead.
func (*PresidioAnonymizer) Descriptor() ([]byte, []int) {
	return file_presidio_proto_rawDescGZIP(), []int{4}
}

func (x *PresidioAnonymizer) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *PresidioAnonymizer) GetNewValue() string {
	if x != nil {
		return x.NewValue
	}
	return ""
}

func (x *PresidioAnonymizer) GetMaskingChar() string {
	if x != nil {
		return x.MaskingChar
	}
	return ""
}

func (x *PresidioAnonymizer) GetCharsToMask() int32 {
	if x != nil {
		return x.CharsToMask
	}
	return 0
}

func (x *PresidioAnonymizer) GetFromEnd() bool {
	if x != nil {
		return x.FromEnd
	}
	return false
}

func (x *PresidioAnonymizer) GetHashType() string {
	if x != nil {
		return x.HashType
	}
	return ""
}

func (x *PresidioAnonymizer) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

// Represents the outcome of an anonymization operation.
type PresidioAnonymizerResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Operation     string                 `protobuf:"bytes,1,opt,name=operation,proto3" json:"operation,omitempty"`
	EntityType    string                 `protobuf:"bytes,2,opt,name=entity_type,json=entityType,proto3" json:"entity_type,omitempty"`
	Start         int32                  `protobuf:"varint,3,opt,name=start,proto3" json:"start,omitempty"`
	End           int32                  `protobuf:"varint,4,opt,name=end,proto3" json:"end,omitempty"`
	Text          string                 `protobuf:"bytes,5,opt,name=text,proto3" json:"text,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PresidioAnonymizerResponse) Reset() {
	*x = PresidioAnonymizerResponse{}
	mi := &file_presidio_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PresidioAnonymizerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PresidioAnonymizerResponse) ProtoMessage() {}

func (x *PresidioAnonymizerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_presidio_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PresidioAnonymizerResponse.ProtoReflect.Descriptor instead.
func (*PresidioAnonymizerResponse) Descriptor() ([]byte, []int) {
	return file_presidio_proto_rawDescGZIP(), []int{5}
}

func (x *PresidioAnonymizerResponse) GetOperation() string {
	if x != nil {
		return x.Operation
	}
	return ""
}

func (x *PresidioAnonymizerResponse) GetEntityType() string {
	if x != nil {
		return x.EntityType
	}
	return ""
}

func (x *PresidioAnonymizerResponse) GetStart() int32 {
	if x != nil {
		return x.Start
	}
	return 0
}

func (x *PresidioAnonymizerResponse) GetEnd() int32 {
	if x != nil {
		return x.End
	}
	return 0
}

func (x *PresidioAnonymizerResponse) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

// Represents a request to analyze and anonymize text.
type PresidioAnalyzerAnomymizerRequest struct {
	state          protoimpl.MessageState         `protogen:"open.v1"`
	Text           string                         `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Language       string                         `protobuf:"bytes,2,opt,name=language,proto3" json:"language,omitempty"`
	ScoreThreshold float64                        `protobuf:"fixed64,3,opt,name=score_threshold,json=scoreThreshold,proto3" json:"score_threshold,omitempty"`
	Entities       []string                       `protobuf:"bytes,4,rep,name=entities,proto3" json:"entities,omitempty"`
	Context        []string                       `protobuf:"bytes,5,rep,name=context,proto3" json:"context,omitempty"`
	Anonymizers    map[string]*PresidioAnonymizer `protobuf:"bytes,6,rep,name=anonymizers,proto3" json:"anonymizers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *PresidioAnalyzerAnomymizerRequest) Reset() {
	*x = PresidioAnalyzerAnomymizerRequest{}
	mi := &file_presidio_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PresidioAnalyzerAnomymizerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PresidioAnalyzerAnomymizerRequest) ProtoMessage() {}

func (x *PresidioAnalyzerAnomymizerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_presidio_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PresidioAnalyzerAnomymizerRequest.ProtoReflect.Descriptor instead.
func (*PresidioAnalyzerAnomymizerRequest) Descriptor() ([]byte, []int) {
	return file_presidio_proto_rawDescGZIP(), []int{6}
}

func (x *PresidioAnalyzerAnomymizerRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *PresidioAnalyzerAnomymizerRequest) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

func (x *PresidioAnalyzerAnomymizerRequest) GetScoreThreshold() float64 {
	if x != nil {
		return x.ScoreThreshold
	}
	return 0
}

func (x *PresidioAnalyzerAnomymizerRequest) GetEntities() []string {
	if x != nil {
		return x.Entities
	}
	return nil
}

func (x *PresidioAnalyzerAnomymizerRequest) GetContext() []string {
	if x != nil {
		return x.Context
	}
	return nil
}

func (x *PresidioAnalyzerAnomymizerRequest) GetAnonymizers() map[string]*PresidioAnonymizer {
	if x != nil {
		return x.Anonymizers
	}
	return nil
}

var File_presidio_proto protoreflect.FileDescriptor

var file_presidio_proto_rawDesc = string([]byte{
	0x0a, 0x0e, 0x70, 0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xa8, 0x01, 0x0a, 0x17, 0x50, 0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x41, 0x6e, 0x61,
	0x6c, 0x79, 0x7a, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x27, 0x0a, 0x0f,
	0x73, 0x63, 0x6f, 0x72, 0x65, 0x5f, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0e, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x54, 0x68, 0x72, 0x65,
	0x73, 0x68, 0x6f, 0x6c, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65,
	0x73, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x05, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x22, 0x61, 0x0a, 0x19, 0x50,
	0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x12, 0x44, 0x0a, 0x10, 0x61, 0x6e, 0x61, 0x6c,
	0x79, 0x7a, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x50, 0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x41, 0x6e, 0x61,
	0x6c, 0x79, 0x7a, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x0f, 0x61,
	0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x22, 0x79,
	0x0a, 0x18, 0x50, 0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x65,
	0x6e, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x22, 0x99, 0x02, 0x0a, 0x19, 0x50, 0x72,
	0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x41, 0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x4d, 0x0a, 0x0b, 0x61,
	0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x2b, 0x2e, 0x50, 0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x41, 0x6e, 0x6f, 0x6e, 0x79,
	0x6d, 0x69, 0x7a, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x6e, 0x6f,
	0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x61,
	0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72, 0x73, 0x12, 0x44, 0x0a, 0x10, 0x61, 0x6e,
	0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x50, 0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x41,
	0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52,
	0x0f, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73,
	0x1a, 0x53, 0x0a, 0x10, 0x41, 0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x29, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x50, 0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f,
	0x41, 0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xd6, 0x01, 0x0a, 0x12, 0x50, 0x72, 0x65, 0x73, 0x69, 0x64,
	0x69, 0x6f, 0x41, 0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x65, 0x77, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x65, 0x77, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21, 0x0a,
	0x0c, 0x6d, 0x61, 0x73, 0x6b, 0x69, 0x6e, 0x67, 0x5f, 0x63, 0x68, 0x61, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x61, 0x73, 0x6b, 0x69, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x72,
	0x12, 0x22, 0x0a, 0x0d, 0x63, 0x68, 0x61, 0x72, 0x73, 0x5f, 0x74, 0x6f, 0x5f, 0x6d, 0x61, 0x73,
	0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x72, 0x73, 0x54, 0x6f,
	0x4d, 0x61, 0x73, 0x6b, 0x12, 0x19, 0x0a, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x65, 0x6e, 0x64,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x66, 0x72, 0x6f, 0x6d, 0x45, 0x6e, 0x64, 0x12,
	0x1b, 0x0a, 0x09, 0x68, 0x61, 0x73, 0x68, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x68, 0x61, 0x73, 0x68, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x97,
	0x01, 0x0a, 0x1a, 0x50, 0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x41, 0x6e, 0x6f, 0x6e, 0x79,
	0x6d, 0x69, 0x7a, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x65, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x22, 0xde, 0x02, 0x0a, 0x21, 0x50, 0x72, 0x65,
	0x73, 0x69, 0x64, 0x69, 0x6f, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72, 0x41, 0x6e, 0x6f,
	0x6d, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65,
	0x78, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x27,
	0x0a, 0x0f, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x5f, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0e, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x54, 0x68,
	0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x69, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x69, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12, 0x55, 0x0a,
	0x0b, 0x61, 0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72, 0x73, 0x18, 0x06, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x33, 0x2e, 0x50, 0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x41, 0x6e, 0x61,
	0x6c, 0x79, 0x7a, 0x65, 0x72, 0x41, 0x6e, 0x6f, 0x6d, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65,
	0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x61, 0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69,
	0x7a, 0x65, 0x72, 0x73, 0x1a, 0x53, 0x0a, 0x10, 0x41, 0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a,
	0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x29, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x50, 0x72, 0x65, 0x73,
	0x69, 0x64, 0x69, 0x6f, 0x41, 0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x81, 0x02, 0x0a, 0x1a, 0x50, 0x72,
	0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x52, 0x65, 0x64, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x41, 0x0a, 0x07, 0x41, 0x6e, 0x61, 0x6c,
	0x79, 0x7a, 0x65, 0x12, 0x18, 0x2e, 0x50, 0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x41, 0x6e,
	0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e,
	0x50, 0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x09, 0x41,
	0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x12, 0x1a, 0x2e, 0x50, 0x72, 0x65, 0x73, 0x69,
	0x64, 0x69, 0x6f, 0x41, 0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x50, 0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x41,
	0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x58, 0x0a, 0x13, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x41, 0x6e,
	0x64, 0x41, 0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x12, 0x22, 0x2e, 0x50, 0x72, 0x65,
	0x73, 0x69, 0x64, 0x69, 0x6f, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72, 0x41, 0x6e, 0x6f,
	0x6d, 0x79, 0x6d, 0x69, 0x7a, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b,
	0x2e, 0x50, 0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x41, 0x6e, 0x6f, 0x6e, 0x79, 0x6d, 0x69,
	0x7a, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x4a, 0x5a,
	0x48, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x52, 0x4b, 0x61, 0x70,
	0x61, 0x64, 0x69, 0x61, 0x30, 0x31, 0x2f, 0x50, 0x72, 0x65, 0x73, 0x69, 0x64, 0x69, 0x6f, 0x52,
	0x65, 0x64, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f,
	0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x65, 0x73, 0x69,
	0x64, 0x69, 0x6f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_presidio_proto_rawDescOnce sync.Once
	file_presidio_proto_rawDescData []byte
)

func file_presidio_proto_rawDescGZIP() []byte {
	file_presidio_proto_rawDescOnce.Do(func() {
		file_presidio_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_presidio_proto_rawDesc), len(file_presidio_proto_rawDesc)))
	})
	return file_presidio_proto_rawDescData
}

var file_presidio_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_presidio_proto_goTypes = []any{
	(*PresidioAnalyzerRequest)(nil),           // 0: PresidioAnalyzerRequest
	(*PresidioAnalyzerResponses)(nil),         // 1: PresidioAnalyzerResponses
	(*PresidioAnalyzerResponse)(nil),          // 2: PresidioAnalyzerResponse
	(*PresidioAnonymizerRequest)(nil),         // 3: PresidioAnonymizerRequest
	(*PresidioAnonymizer)(nil),                // 4: PresidioAnonymizer
	(*PresidioAnonymizerResponse)(nil),        // 5: PresidioAnonymizerResponse
	(*PresidioAnalyzerAnomymizerRequest)(nil), // 6: PresidioAnalyzerAnomymizerRequest
	nil, // 7: PresidioAnonymizerRequest.AnonymizersEntry
	nil, // 8: PresidioAnalyzerAnomymizerRequest.AnonymizersEntry
}
var file_presidio_proto_depIdxs = []int32{
	2, // 0: PresidioAnalyzerResponses.analyzer_results:type_name -> PresidioAnalyzerResponse
	7, // 1: PresidioAnonymizerRequest.anonymizers:type_name -> PresidioAnonymizerRequest.AnonymizersEntry
	2, // 2: PresidioAnonymizerRequest.analyzer_results:type_name -> PresidioAnalyzerResponse
	8, // 3: PresidioAnalyzerAnomymizerRequest.anonymizers:type_name -> PresidioAnalyzerAnomymizerRequest.AnonymizersEntry
	4, // 4: PresidioAnonymizerRequest.AnonymizersEntry.value:type_name -> PresidioAnonymizer
	4, // 5: PresidioAnalyzerAnomymizerRequest.AnonymizersEntry.value:type_name -> PresidioAnonymizer
	0, // 6: PresidioRedactionProcessor.Analyze:input_type -> PresidioAnalyzerRequest
	3, // 7: PresidioRedactionProcessor.Anonymize:input_type -> PresidioAnonymizerRequest
	6, // 8: PresidioRedactionProcessor.AnalyzeAndAnonymize:input_type -> PresidioAnalyzerAnomymizerRequest
	1, // 9: PresidioRedactionProcessor.Analyze:output_type -> PresidioAnalyzerResponses
	5, // 10: PresidioRedactionProcessor.Anonymize:output_type -> PresidioAnonymizerResponse
	5, // 11: PresidioRedactionProcessor.AnalyzeAndAnonymize:output_type -> PresidioAnonymizerResponse
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_presidio_proto_init() }
func file_presidio_proto_init() {
	if File_presidio_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_presidio_proto_rawDesc), len(file_presidio_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_presidio_proto_goTypes,
		DependencyIndexes: file_presidio_proto_depIdxs,
		MessageInfos:      file_presidio_proto_msgTypes,
	}.Build()
	File_presidio_proto = out.File
	file_presidio_proto_goTypes = nil
	file_presidio_proto_depIdxs = nil
}

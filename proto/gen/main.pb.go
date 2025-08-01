// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: main.proto

package grpcapipb

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type DeleteTeachersConfirmation struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	DeletedIds    []string               `protobuf:"bytes,2,rep,name=deleted_ids,json=deletedIds,proto3" json:"deleted_ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteTeachersConfirmation) Reset() {
	*x = DeleteTeachersConfirmation{}
	mi := &file_main_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteTeachersConfirmation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteTeachersConfirmation) ProtoMessage() {}

func (x *DeleteTeachersConfirmation) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteTeachersConfirmation.ProtoReflect.Descriptor instead.
func (*DeleteTeachersConfirmation) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{0}
}

func (x *DeleteTeachersConfirmation) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *DeleteTeachersConfirmation) GetDeletedIds() []string {
	if x != nil {
		return x.DeletedIds
	}
	return nil
}

type TeacherId struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TeacherId) Reset() {
	*x = TeacherId{}
	mi := &file_main_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TeacherId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TeacherId) ProtoMessage() {}

func (x *TeacherId) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TeacherId.ProtoReflect.Descriptor instead.
func (*TeacherId) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{1}
}

func (x *TeacherId) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type TeacherIds struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ids           []*TeacherId           `protobuf:"bytes,1,rep,name=ids,proto3" json:"ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TeacherIds) Reset() {
	*x = TeacherIds{}
	mi := &file_main_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TeacherIds) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TeacherIds) ProtoMessage() {}

func (x *TeacherIds) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TeacherIds.ProtoReflect.Descriptor instead.
func (*TeacherIds) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{2}
}

func (x *TeacherIds) GetIds() []*TeacherId {
	if x != nil {
		return x.Ids
	}
	return nil
}

type GetTeachersRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Teacher       *Teacher               `protobuf:"bytes,1,opt,name=teacher,proto3" json:"teacher,omitempty"`
	SortBy        []*SortField           `protobuf:"bytes,2,rep,name=sort_by,json=sortBy,proto3" json:"sort_by,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetTeachersRequest) Reset() {
	*x = GetTeachersRequest{}
	mi := &file_main_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTeachersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTeachersRequest) ProtoMessage() {}

func (x *GetTeachersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTeachersRequest.ProtoReflect.Descriptor instead.
func (*GetTeachersRequest) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{3}
}

func (x *GetTeachersRequest) GetTeacher() *Teacher {
	if x != nil {
		return x.Teacher
	}
	return nil
}

func (x *GetTeachersRequest) GetSortBy() []*SortField {
	if x != nil {
		return x.SortBy
	}
	return nil
}

type Teacher struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FirstName     string                 `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName      string                 `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Email         string                 `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	Class         string                 `protobuf:"bytes,5,opt,name=class,proto3" json:"class,omitempty"`
	Subject       string                 `protobuf:"bytes,6,opt,name=subject,proto3" json:"subject,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Teacher) Reset() {
	*x = Teacher{}
	mi := &file_main_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Teacher) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Teacher) ProtoMessage() {}

func (x *Teacher) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Teacher.ProtoReflect.Descriptor instead.
func (*Teacher) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{4}
}

func (x *Teacher) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Teacher) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *Teacher) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *Teacher) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *Teacher) GetClass() string {
	if x != nil {
		return x.Class
	}
	return ""
}

func (x *Teacher) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

type Teachers struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Teachers      []*Teacher             `protobuf:"bytes,1,rep,name=teachers,proto3" json:"teachers,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Teachers) Reset() {
	*x = Teachers{}
	mi := &file_main_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Teachers) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Teachers) ProtoMessage() {}

func (x *Teachers) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Teachers.ProtoReflect.Descriptor instead.
func (*Teachers) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{5}
}

func (x *Teachers) GetTeachers() []*Teacher {
	if x != nil {
		return x.Teachers
	}
	return nil
}

type StudentCount struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        bool                   `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	StudentCount  int32                  `protobuf:"varint,2,opt,name=student_count,json=studentCount,proto3" json:"student_count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StudentCount) Reset() {
	*x = StudentCount{}
	mi := &file_main_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StudentCount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StudentCount) ProtoMessage() {}

func (x *StudentCount) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StudentCount.ProtoReflect.Descriptor instead.
func (*StudentCount) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{6}
}

func (x *StudentCount) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

func (x *StudentCount) GetStudentCount() int32 {
	if x != nil {
		return x.StudentCount
	}
	return 0
}

var File_main_proto protoreflect.FileDescriptor

const file_main_proto_rawDesc = "" +
	"\n" +
	"\n" +
	"main.proto\x12\x04main\x1a\x17validate/validate.proto\x1a\x0estudents.proto\"U\n" +
	"\x1aDeleteTeachersConfirmation\x12\x16\n" +
	"\x06status\x18\x01 \x01(\tR\x06status\x12\x1f\n" +
	"\vdeleted_ids\x18\x02 \x03(\tR\n" +
	"deletedIds\"\x1b\n" +
	"\tTeacherId\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"/\n" +
	"\n" +
	"TeacherIds\x12!\n" +
	"\x03ids\x18\x01 \x03(\v2\x0f.main.TeacherIdR\x03ids\"g\n" +
	"\x12GetTeachersRequest\x12'\n" +
	"\ateacher\x18\x01 \x01(\v2\r.main.TeacherR\ateacher\x12(\n" +
	"\asort_by\x18\x02 \x03(\v2\x0f.main.SortFieldR\x06sortBy\"\x83\x02\n" +
	"\aTeacher\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x124\n" +
	"\n" +
	"first_name\x18\x02 \x01(\tB\x15\xfaB\x12r\x10\x10\x012\f^[A-Za-z ]*$R\tfirstName\x122\n" +
	"\tlast_name\x18\x03 \x01(\tB\x15\xfaB\x12r\x10\x10\x012\f^[A-Za-z ]*$R\blastName\x12\x1d\n" +
	"\x05email\x18\x04 \x01(\tB\a\xfaB\x04r\x02`\x01R\x05email\x12+\n" +
	"\x05class\x18\x05 \x01(\tB\x15\xfaB\x12r\x10\x10\x012\f^[1-9][A-Z]$R\x05class\x122\n" +
	"\asubject\x18\x06 \x01(\tB\x18\xfaB\x15r\x13\x10\x012\x0f^[A-Za-z0-9 ]*$R\asubject\"5\n" +
	"\bTeachers\x12)\n" +
	"\bteachers\x18\x01 \x03(\v2\r.main.TeacherR\bteachers\"K\n" +
	"\fStudentCount\x12\x16\n" +
	"\x06status\x18\x01 \x01(\bR\x06status\x12#\n" +
	"\rstudent_count\x18\x02 \x01(\x05R\fstudentCount2\xf4\x02\n" +
	"\x0eTeacherService\x127\n" +
	"\vGetTeachers\x12\x18.main.GetTeachersRequest\x1a\x0e.main.Teachers\x12-\n" +
	"\vAddTeachers\x12\x0e.main.Teachers\x1a\x0e.main.Teachers\x120\n" +
	"\x0eUpdateTeachers\x12\x0e.main.Teachers\x1a\x0e.main.Teachers\x12D\n" +
	"\x0eDeleteTeachers\x12\x10.main.TeacherIds\x1a .main.DeleteTeachersConfirmation\x12<\n" +
	"\x19GetStudentsByClassTeacher\x12\x0f.main.TeacherId\x1a\x0e.main.Students\x12D\n" +
	"\x1dGetStudentCountByClassTeacher\x12\x0f.main.TeacherId\x1a\x12.main.StudentCountB\x16Z\x14/proto/gen;grpcapipbb\x06proto3"

var (
	file_main_proto_rawDescOnce sync.Once
	file_main_proto_rawDescData []byte
)

func file_main_proto_rawDescGZIP() []byte {
	file_main_proto_rawDescOnce.Do(func() {
		file_main_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_main_proto_rawDesc), len(file_main_proto_rawDesc)))
	})
	return file_main_proto_rawDescData
}

var file_main_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_main_proto_goTypes = []any{
	(*DeleteTeachersConfirmation)(nil), // 0: main.DeleteTeachersConfirmation
	(*TeacherId)(nil),                  // 1: main.TeacherId
	(*TeacherIds)(nil),                 // 2: main.TeacherIds
	(*GetTeachersRequest)(nil),         // 3: main.GetTeachersRequest
	(*Teacher)(nil),                    // 4: main.Teacher
	(*Teachers)(nil),                   // 5: main.Teachers
	(*StudentCount)(nil),               // 6: main.StudentCount
	(*SortField)(nil),                  // 7: main.SortField
	(*Students)(nil),                   // 8: main.Students
}
var file_main_proto_depIdxs = []int32{
	1,  // 0: main.TeacherIds.ids:type_name -> main.TeacherId
	4,  // 1: main.GetTeachersRequest.teacher:type_name -> main.Teacher
	7,  // 2: main.GetTeachersRequest.sort_by:type_name -> main.SortField
	4,  // 3: main.Teachers.teachers:type_name -> main.Teacher
	3,  // 4: main.TeacherService.GetTeachers:input_type -> main.GetTeachersRequest
	5,  // 5: main.TeacherService.AddTeachers:input_type -> main.Teachers
	5,  // 6: main.TeacherService.UpdateTeachers:input_type -> main.Teachers
	2,  // 7: main.TeacherService.DeleteTeachers:input_type -> main.TeacherIds
	1,  // 8: main.TeacherService.GetStudentsByClassTeacher:input_type -> main.TeacherId
	1,  // 9: main.TeacherService.GetStudentCountByClassTeacher:input_type -> main.TeacherId
	5,  // 10: main.TeacherService.GetTeachers:output_type -> main.Teachers
	5,  // 11: main.TeacherService.AddTeachers:output_type -> main.Teachers
	5,  // 12: main.TeacherService.UpdateTeachers:output_type -> main.Teachers
	0,  // 13: main.TeacherService.DeleteTeachers:output_type -> main.DeleteTeachersConfirmation
	8,  // 14: main.TeacherService.GetStudentsByClassTeacher:output_type -> main.Students
	6,  // 15: main.TeacherService.GetStudentCountByClassTeacher:output_type -> main.StudentCount
	10, // [10:16] is the sub-list for method output_type
	4,  // [4:10] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_main_proto_init() }
func file_main_proto_init() {
	if File_main_proto != nil {
		return
	}
	file_students_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_main_proto_rawDesc), len(file_main_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_main_proto_goTypes,
		DependencyIndexes: file_main_proto_depIdxs,
		MessageInfos:      file_main_proto_msgTypes,
	}.Build()
	File_main_proto = out.File
	file_main_proto_goTypes = nil
	file_main_proto_depIdxs = nil
}

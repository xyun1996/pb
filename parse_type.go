package main

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

func parseType(msg *desc.MessageDescriptor) *objType {
	obj := objType{}.newtype().(*objType)
	obj.Title = msg.GetName()
	for _, f := range msg.GetFields() {
		parseField(obj, f)
	}
	return obj
}

func parseField(obj *objType, field *desc.FieldDescriptor) {
	prop := field.GetName()
	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if !field.IsMap() && field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
			arr := arrType{}.newtype().(*arrType)
			arr.Items = parseType(field.GetMessageType())
			obj.Properties[prop] = arr
			return
		}
		if field.IsMap() {
			objT := objType{}.newtype().(*objType)
			objT.AddProp = parseMap(field)
			obj.Properties[prop] = objT
			return
		}
		obj.Properties[prop] = parseType(field.GetMessageType())
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		obj.Properties[prop] = stringType{}.newtype()
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		obj.Properties[prop] = boolType{}.newtype()
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
	case descriptorpb.FieldDescriptorProto_TYPE_INT32:
		obj.Properties[prop] = int32Type{}.newtype()
	case descriptorpb.FieldDescriptorProto_TYPE_INT64:
		obj.Properties[prop] = int64Type{}.newtype()
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		obj.Properties[prop] = uint64Type{}.newtype()
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32:
		obj.Properties[prop] = uint32Type{}.newtype()
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
	default:
		panic(fmt.Sprintf("err type :%v", field))
	}
}

func parseMap(field *desc.FieldDescriptor) *objType {
	obj := objType{}.newtype().(*objType)
	typ := field.GetMapValueType()
	parseField(obj, typ)
	return obj
}

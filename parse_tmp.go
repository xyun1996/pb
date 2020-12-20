package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
)

var (
	tmpobj = "tmpobj.proto"
	tmpmsg = "tmpmsg.proto"
)

type pbField struct {
	Rule    string `json:"rule"`
	Type    string `json:"type"`
	KeyType string `json:"keytype,omitempty"`
	Name    string `json:"name"`
	Id      int32  `json:"id"`
}

type pbMsg struct {
	Name   string     `json:"name"`
	Fields []*pbField `json:"fields"`
}

func parseAllV2() {
	p := protoparse.Parser{}
	files := []string{"message/login.proto"} //getFiles("object", "message")
	dss, err := p.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	for _, ds := range dss {
		msg := ds.FindMessage("proto_csmsg.CS_Login")
		if msg == nil {
			continue
		}

		fmt.Println(msg)
	}
}

var (
	msgFile *desc.FileDescriptor
	objFile *desc.FileDescriptor
)

func msg2Json(msg *desc.MessageDescriptor) {
	for _, f := range msg.GetFields() {
		if f.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			name := f.GetMessageType().GetFullyQualifiedName()
			pkg := f.GetMessageType().GetFile().GetPackage()
			var objMsg *desc.MessageDescriptor
			if pkg == "proto_object" {
				objMsg = objFile.FindMessage(name)
				if objMsg == nil {
					panic("objMsg == nil")
				}
			} else if pkg == "proto_message" {
				objMsg = msgFile.FindMessage(name)
				if objMsg == nil {
					panic("objMsg == nil")
				}
			}
			fmt.Println(objMsg.GetFullyQualifiedName())
			fmt.Printf("\n")
			msg2Json(objMsg)
		}
	}
}

func ParseV1() {
	p := protoparse.Parser{}
	result, err := p.ParseFiles("message/message_id.proto")
	if err != nil {
		fmt.Printf("ParseFiles err:%v", err)
		return
	}
	pkgName := result[0].GetPackage()
	fmt.Println(pkgName)
	pbMsgs := make([]*pbMsg, 0)
	for _, msgD := range result[0].GetMessageTypes() {
		pbFields := make([]*pbField, 0)
		for _, field := range msgD.GetFields() {

			fieldType := field.GetType().String()

			if fieldType == "TYPE_MESSAGE" && !field.IsMap() {
				fieldType = field.AsFieldDescriptorProto().GetTypeName()
			}

			fieldStruct := &pbField{
				Rule: field.GetLabel().String(),
				Type: fieldType,
				Name: field.GetName(),
				Id:   field.GetNumber(),
			}

			if field.IsMap() {
				keyType, valType := getMapKVType(field.GetMessageType())
				fieldStruct.KeyType = keyType
				fieldStruct.Type = valType
			}
			pbFields = append(pbFields, fieldStruct)
		}
		msgStruct := &pbMsg{
			Name:   msgD.GetName(),
			Fields: pbFields,
		}
		pbMsgs = append(pbMsgs, msgStruct)
	}
	jsonData, err := json.MarshalIndent(pbMsgs, "", "\t")
	if err != nil {
		fmt.Printf("json.Marshal err:%v", err)
		return
	}
	file, err := os.OpenFile("output.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Printf("OpenFile err:%v", err)
		return
	}
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Printf("Write data err:%v", err)
		return
	}

	file.Close()
}

func getMapKVType(field *desc.MessageDescriptor) (string, string) {
	fields := field.GetFields()
	if len(fields) < 2 {
		return "", ""
	}
	return fields[0].AsFieldDescriptorProto().GetTypeName(), fields[1].AsFieldDescriptorProto().GetTypeName()
}

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
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

func main() {
	p := protoparse.Parser{}
	result, err := p.ParseFiles("message/login.proto")
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
	jsonData, err := json.Marshal(pbMsgs)
	if err != nil {
		fmt.Printf("json.Marshal err:%v", err)
		return
	}
	file, err := os.OpenFile("output.txt", os.O_CREATE, os.ModePerm)
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
	//for _, eum := range result[0].GetEnumTypes(){
	//	for _, value := range eum.GetValues(){
	//		fmt.Printf("%s = %d\n",value.GetName(), value.GetNumber())
	//	}
	//}
}

func getMessageType(field int) string {
	return ""
}

func getMapKVType(field *desc.MessageDescriptor) (string, string) {
	fields := field.GetFields()
	if len(fields) < 2 {
		return "", ""
	}
	return fields[0].AsFieldDescriptorProto().GetTypeName(), fields[1].AsFieldDescriptorProto().GetTypeName()
}

func test(msg proto.Message) {
	msg.ProtoMessage()
}

package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/jhump/protoreflect/desc/protoparse"
)

func Init() {
	mergeFiles(getFiles("message"), tmpmsg)
	mergeFiles(getFiles("object"), tmpobj)
	dss, err := protoparse.Parser{}.ParseFiles(tmpmsg, tmpobj)
	if err != nil {
		panic(err)
	}
	msgFile = dss[0]
	objFile = dss[1]
}

func Msg2Json(name string) []byte {
	msg := msgFile.FindMessage(name)
	obj := parseType(msg)
	b, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		panic(err)
	}
	return b
}

func SaveMsg2Json(name string) {
	msg := msgFile.FindMessage(name)
	obj := parseType(msg)
	fs, err := os.OpenFile("message.json", os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		panic(err)
	}
	fs.Write(b)
	fs.Close()
}

func ParseIdMap() {
	ds, err := protoparse.Parser{}.ParseFiles(*idmap)
	if err != nil {
		panic(err)
	}
	reqMap := make(map[int]string)
	for _, enum := range ds[0].GetEnumTypes()[0].GetValues() {
		if enum.GetNumber() == 0 {
			continue
		}
		if !strings.HasPrefix(enum.GetName(), "MSG_CS") {
			continue
		}
		reqMap[int(enum.GetNumber())] = enum.GetName()
	}
	fs, err := os.OpenFile(idmapFile, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(reqMap, "", "\t")
	if err != nil {
		panic(err)
	}
	fs.Write(b)
	fs.Close()
}

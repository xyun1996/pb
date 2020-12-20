package main

import "math"

type jsonType interface {
	newtype() jsonType
}

type int32Type struct {
	Type    string `json:"type"`
	Minimum int32  `json:"minimum"`
	Maximum int32  `json:"maximum"`
}

type uint32Type struct {
	Type    string `json:"type"`
	Minimum uint32 `json:"minimum"`
	Maximum uint32 `json:"maximum"`
}

type uint64Type struct {
	Type    string `json:"type"`
	Minimum uint64 `json:"minimum"`
	Maximum uint64 `json:"maximum"`
}

type int64Type struct {
	Type    string `json:"type"`
	Minimum int64  `json:"minimum"`
	Maximum int64  `json:"maximum"`
}

type objType struct {
	Title      string              `json:"title,omitempty"`
	Type       string              `json:"type"`
	AddProp    jsonType            `json:"additionalProperties,omitempty"`
	Properties map[string]jsonType `json:"properties,omitempty"`
}

type mapType struct {
	Title string `json:"title"`
}

type arrType struct {
	Type  string   `json:"type"`
	Items jsonType `json:"items,omitempty"`
}

type boolType struct {
	Type string `json:"type"`
}

type stringType struct {
	Type string `json:"type"`
}

func (m int32Type) newtype() jsonType {
	return &int32Type{
		Type:    "integer",
		Maximum: math.MaxInt32,
		Minimum: math.MinInt32,
	}
}

func (m uint64Type) newtype() jsonType {
	return &uint64Type{
		Type:    "integer",
		Maximum: math.MaxUint64,
	}
}

func (m uint32Type) newtype() jsonType {
	return &uint32Type{
		Type:    "integer",
		Maximum: math.MaxUint32,
	}
}

func (m int64Type) newtype() jsonType {
	return &int64Type{
		Type:    "integer",
		Maximum: math.MaxInt64,
		Minimum: math.MinInt64,
	}
}

func (m boolType) newtype() jsonType {
	return &boolType{
		Type: "boolean",
	}
}

func (m stringType) newtype() jsonType {
	return &boolType{
		Type: "string",
	}
}

func (m objType) newtype() jsonType {
	return &objType{
		Type:       "object",
		Properties: map[string]jsonType{},
	}
}

func (m arrType) newtype() jsonType {
	return &arrType{
		Type: "array",
	}
}

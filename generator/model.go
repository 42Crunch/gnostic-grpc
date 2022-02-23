package generator

import (
	surface "github.com/google/gnostic/surface"
)

type Field struct {
	*surface.Field
}

type Type struct {
	*surface.Type
	Fields []*Field `protobuf:"bytes,5,rep,name=fields,proto3" json:"fields,omitempty"`
}

type Model struct {
	*surface.Model
	Types []*Type `protobuf:"bytes,2,rep,name=types,proto3" json:"types,omitempty"` // the types used by the API
}

package generator

import (
	"fmt"
	"strings"

	surface "github.com/google/gnostic/surface"
)

type Field struct {
	*surface.Field
	ParamFormat        *string
	ParamPattern       *string
	ParamMinLength     *int64
	ParamMaxLength     *int64
	ParamMinimum       *float64
	ParamMaximum       *float64
	ParamRequired      *bool
	ParamExampleString *string
	ParamExampleInt    *int64
}

func String(s string) *string {
	return &s
}

func Int64(i int64) *int64 {
	return &i
}

func Float64(i float64) *float64 {
	return &i
}

func (f *Field) toTags() map[string]string {
	tags := make(map[string]string)
	if f.ParamPattern != nil {
		pattern := strings.Replace(*f.ParamPattern, "\n", "\\n", -1)
		pattern = strings.Replace(pattern, "\r", "\\r", -1)
		tags["pattern"] = pattern
	}
	if f.ParamFormat != nil {
		tags["format"] = *f.ParamFormat
	}
	if f.ParamMinimum != nil {
		tags["minimum"] = fmt.Sprintf("%f", *f.ParamMinimum)
	}
	if f.ParamMaximum != nil {
		tags["maximum"] = fmt.Sprintf("%f", *f.ParamMaximum)
	}
	if f.ParamMinLength != nil {
		tags["minLength"] = fmt.Sprintf("%d", *f.ParamMinLength)
	}
	if f.ParamMaxLength != nil {
		tags["maxLength"] = fmt.Sprintf("%d", *f.ParamMaxLength)
	}

	return tags
}

func Bool(i bool) *bool {
	return &i
}

type Type struct {
	*surface.Type
	Fields        []*Field `protobuf:"bytes,5,rep,name=fields,proto3" json:"fields,omitempty"`
	ParamRequired []string
}

type Model struct {
	*surface.Model
	Types []*Type `protobuf:"bytes,2,rep,name=types,proto3" json:"types,omitempty"` // the types used by the API
}

func (t *Type) GetField(name string) *Field {
	for _, field := range t.Fields {
		if field.GetName() == name {
			return field
		}
	}

	return nil
}

package json2struct

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/dapings/examples/go-programing-tour-2023/tour/internal/word"
)

const (
	TYPE_MAP_STRING_INTERFACE = "map[string]interface{}"
	TYPE_INTERFACE            = "[]interface{}"
)

type Parse struct {
	Source     map[string]interface{}
	StructTag  string
	StructName string
	Output     Output
	Children   Output
}

type Output []string

func (op *Output) appendSegment(format, title string, args ...interface{}) {
	var s []interface{}
	s = append(s, word.UnderscoreToUpperCamelCase(title))
	if len(args) != 0 {
		s = append(s, args...)
		format = "    " + format
	}
	*op = append(*op, fmt.Sprintf(format, s...))
}

func (op *Output) appendSuffix() {
	*op = append(*op, "}\n")
}

func NewParser(s string) (*Parse, error) {
	source := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &source); err != nil {
		return nil, err
	}

	return &Parse{
		Source:     source,
		StructTag:  "type %s struct {",
		StructName: "tour",
	}, nil
}

func (p *Parse) JSON2Struct() string {
	p.Output.appendSegment(p.StructTag, p.StructName)
	for parentName, parentVal := range p.Source {
		valType := reflect.TypeOf(parentVal).String()
		if valType == TYPE_INTERFACE {
			p.toParentList(parentName, parentVal.([]interface{}), true)
		} else {
			var fields Fields
			fields.appendSegment(parentName, FieldSegment{
				Format: "%s",
				FieldValues: []FieldValue{
					{CamelCase: false, Val: valType},
				},
			})
			p.Output.appendSegment("%s %s", fields[0].Name, fields[0].Type)
		}
	}
	p.Output.appendSuffix()
	return strings.Join(append(p.Output, p.Children...), "\n")
}

func (p *Parse) toParentList(parentName string, parentVals []interface{}, isTop bool) {
	var fields Fields
	for _, val := range parentVals {
		valType := reflect.TypeOf(val).String()
		if valType == TYPE_MAP_STRING_INTERFACE {
			fields = append(fields, p.handleParentTypeMapInterface(val.(map[string]interface{}))...)
			p.Children.appendSegment(p.StructTag, parentName)
			for _, field := range fields.removeDuplicate() {
				p.Children.appendSegment("%s %s", field.Name, field.Type)
			}
			p.Children.appendSuffix()
			if isTop {
				valType = word.UnderscoreToUpperCamelCase(parentName)
			}
		}

		if isTop {
			p.Output.appendSegment("%s %s%s", parentName, "[]", valType)
		}
		break
	}
}

func (p *Parse) handleParentTypeMapInterface(vals map[string]interface{}) Fields {
	var fields Fields
	for fieldName, fieldVal := range vals {
		fieldValType := reflect.TypeOf(fieldVal).String()
		fieldSegment := FieldSegment{
			Format: "%s",
			FieldValues: []FieldValue{
				{CamelCase: true, Val: fieldValType},
			},
		}
		switch fieldValType {
		case TYPE_INTERFACE:
			fieldSegment = p.handleTypeInterface(fieldName, fieldVal.([]interface{}))
		case TYPE_MAP_STRING_INTERFACE:
			fieldSegment = p.handleTypeMapInterface(fieldName, fieldVal.(map[string]interface{}))
		}
		fields.appendSegment(fieldName, fieldSegment)
	}
	return fields
}

func (p *Parse) handleTypeInterface(fieldName string, fieldVals []interface{}) FieldSegment {
	fieldSegment := FieldSegment{
		Format: "%s%s",
		FieldValues: []FieldValue{
			{CamelCase: false, Val: "[]"},
			{CamelCase: true, Val: fieldName},
		},
	}
	p.toParentList(fieldName, fieldVals, false)
	return fieldSegment
}

func (p *Parse) handleTypeMapInterface(fieldName string, fieldVals map[string]interface{}) FieldSegment {
	fieldSegment := FieldSegment{
		Format: "%s",
		FieldValues: []FieldValue{
			{CamelCase: true, Val: fieldName},
		},
	}
	p.toChildrenStruct(fieldName, fieldVals)
	return fieldSegment
}

func (p *Parse) toChildrenStruct(parentName string, fieldVals map[string]interface{}) {
	p.Children.appendSegment(p.StructTag, parentName)
	for fieldName, fieldVal := range fieldVals {
		p.Children.appendSegment("%s %s", fieldName, reflect.TypeOf(fieldVal).String())
	}
	p.Children.appendSuffix()
}

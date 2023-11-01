package json2struct

import (
	"fmt"

	"github.com/dapings/examples/go-programing-tour-2023/tour/internal/word"
)

type FieldSegment struct {
	Format      string
	FieldValues []FieldValue
}

type FieldValue struct {
	CamelCase bool
	Val       string
}

type Field struct {
	Name string
	Type string
}

type Fields []*Field

func (f *Fields) appendSegment(name string, segment FieldSegment) {
	var s []interface{}
	for _, fv := range segment.FieldValues {
		val := fv.Val
		if fv.CamelCase {
			val = word.UnderscoreToUpperCamelCase(val)
		}

		s = append(s, val)
	}
	*f = append(*f, &Field{
		Name: word.UnderscoreToUpperCamelCase(name),
		Type: fmt.Sprintf(segment.Format, s...),
	})
}

func (f *Fields) removeDuplicate() Fields {
	fields := Fields{}
	dupMap := make(map[string]bool)
	for _, entry := range *f {
		if _, exist := dupMap[entry.Name]; !exist {
			dupMap[entry.Name] = true
			fields = append(fields, entry)
		}
	}

	return fields
}

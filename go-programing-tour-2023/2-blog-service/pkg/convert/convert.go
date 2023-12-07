package convert

import (
	"strconv"
)

// 类型转换

type StrTo string

func (t StrTo) String() string {
	return string(t)
}

func (t StrTo) Int() (int, error) {
	return strconv.Atoi(t.String())
}

func (t StrTo) MustInt() int {
	v, _ := t.Int()
	return v
}

func (t StrTo) UInt32() (uint32, error) {
	v, err := t.Int()
	return uint32(v), err
}

func (t StrTo) MustUInt32() uint32 {
	v, _ := t.UInt32()
	return v
}

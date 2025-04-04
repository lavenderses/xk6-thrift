package thrift

import (
	"context"
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
)

type TEnum struct {
	value int32
}

func NewTEnum(v int32) TEnum {
	return TEnum{value: v}
}

func (p TEnum) Equals(other *TValue) bool {
	o, ok := (*other).(TEnum)
	if !ok {
		return false
	}
	return p.value == o.value
}

// See [Thrift IDL protocol spec]
//
//	<field> ::= <field-begin> <field-data> <field-end>
//	<field-data> ::= I32
//
// [Thrift IDL protocol spec]: https://github.com/apache/thrift/blob/eec0b584e657e4250e22f3fd492858d632e2aa7b/doc/specs/thrift-protocol-spec.md
func (p TEnum) WriteFieldData(cxt context.Context, oprot thrift.TProtocol) error {
	err := oprot.WriteI32(cxt, p.value)
	if err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.id (1) field write error: ", p), err)
	}
	return nil
}

func (p TEnum) TType() thrift.TType {
	return thrift.I32
}

func ReadEnum(cxt context.Context, iproto thrift.TProtocol) (TValue, error) {
	v, err := iproto.ReadI32(cxt)
	if err != nil {
		return nil, thrift.PrependError("error while reading i32 field", err)
	}

	res := NewTEnum(v)
	return res, nil
}

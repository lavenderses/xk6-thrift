package thrift

import (
	"context"
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
)

type TBool struct {
	value bool
}

func NewTBool(v bool) TBool {
	return TBool{value: v}
}

func (p TBool) Equals(other *TValue) bool {
	o, ok := (*other).(TBool)
	if !ok {
		return false
	}
	return p.value == o.value
}

// See [Thrift IDL protocol spec]
//
//	<field> ::= <field-begin> <field-data> <field-end>
//	<field-data> ::= BOOL
//
// [Thrift IDL protocol spec]: https://github.com/apache/thrift/blob/eec0b584e657e4250e22f3fd492858d632e2aa7b/doc/specs/thrift-protocol-spec.md
func (p TBool) WriteFieldData(cxt context.Context, oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteBool(cxt, p.value); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T.id (1) field write error: ", p), err)
		return
	}
	return
}

func (p TBool) TType() thrift.TType {
	return thrift.BOOL
}

func ReadBool(cxt context.Context, iprot thrift.TProtocol) (TValue, error) {
	v, err := iprot.ReadBool(cxt)
	if err != nil {
		return nil, thrift.PrependError("error while reading boolean: ", err)
	}

	res := NewTBool(v)
	return res, nil
}

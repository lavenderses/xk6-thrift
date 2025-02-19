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
//   <field> ::= <field-begin> <field-data> <field-end>
//   <field-data> ::= BOOL
//
// [Thrift IDL protocol spec]: https://github.com/apache/thrift/blob/eec0b584e657e4250e22f3fd492858d632e2aa7b/doc/specs/thrift-protocol-spec.md
func (p TBool) WriteField(cxt context.Context, oprot thrift.TProtocol, fid int16, fname string) (err error) {
	if err = oprot.WriteFieldBegin(cxt, fname, thrift.BOOL, fid); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write field begin error %d:%s: ", p, fid, fname), err)
		return
	}
	if err = oprot.WriteBool(cxt, p.value); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T.id (1) field write error: ", p), err)
		return
	}
	if err = oprot.WriteFieldEnd(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write field end error %d:%s: ", p, fid, fname), err)
		return
	}
	return
}

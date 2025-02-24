package thrift

import (
	"context"
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
)

type TMap struct {
	value     map[TValue]TValue
	keyType   thrift.TType
	valueType thrift.TType
}

func NewTMap(v *map[TValue]TValue) *TMap {
	return &TMap{value: *v}
}

func (p *TMap) Equals(other *TValue) bool {
	o, ok := (*other).(*TMap)
	if !ok {
		return false
	}
	if len(p.value) != len(o.value) {
		return false
	}
	for pk, pv := range p.value {
		ov := o.value[pk]
		if !pv.Equals(&ov) {
			return false
		}
	}
	return true
}

// See [Thrift IDL protocol spec]
//
//   <field>      ::= <field-begin> <map> <field-end>
//   <map>        ::= <map-begin> <field-data>* <map-end>
//   <field-data> ::= I8 | I16 | I32 | I64 | DOUBLE | STRING | BINARY
//                    <struct> | <map> | <list> | <set>
//
// [Thrift IDL protocol spec]: https://github.com/apache/thrift/blob/eec0b584e657e4250e22f3fd492858d632e2aa7b/doc/specs/thrift-protocol-spec.md
func (p *TMap) WriteField(cxt context.Context, oprot thrift.TProtocol, fid int16, fname string) (err error) {
	if err = oprot.WriteFieldBegin(cxt, fname, thrift.MAP, fid); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write field begin error %d:%s: ", p, fid, fname), err)
		return
	}
	if err = oprot.WriteMapBegin(cxt, p.keyType, p.valueType, len(p.value)); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write map begin error %d:%s", p, fid, fname), err)
		return
	}

	for k, v := range p.value {
		if err = p.writeFieldDataKeyValue(cxt, oprot, k, v); err != nil {
			err = thrift.PrependError(fmt.Sprintf("%T write key value (field %d) error", p, fid), err)
			return
		}
	}

	if err = oprot.WriteMapEnd(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write map end error %d:%s", p, fid, fname), err)
		return
	}
	if err = oprot.WriteFieldEnd(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write field end error %d:%s: ", p, fid, fname), err)
		return
	}

	return
}

// See the above spec.
//
//   <map>        ::= <map-begin> <field-data>* <map-end>
//   <field-data> ::= I8 | I16 | I32 | I64 | DOUBLE | STRING | BINARY
//                    <struct> | <map> | <list> | <set>
func (p *TMap) writeFieldDataKeyValue(cxt context.Context, oprot thrift.TProtocol, k, v TValue) (err error) {
	if err = p.writeFieldData(cxt, oprot, k); err != nil {
		return
	}
	if err = p.writeFieldData(cxt, oprot, v); err != nil {
		return
	}
	return
}

// See the above spec.
//
//   <field-data> ::= I8 | I16 | I32 | I64 | DOUBLE | STRING | BINARY
//                    <struct> | <map> | <list> | <set>
func (p *TMap) writeFieldData(cxt context.Context, oprot thrift.TProtocol, value TValue) (err error) {
	if o, ok := value.(TString); ok {
		err = oprot.WriteString(cxt, o.value)
	} else if o, ok := value.(TBool); ok {
		err = oprot.WriteBool(cxt, o.value)
	} else {
		return
	}
	return
}

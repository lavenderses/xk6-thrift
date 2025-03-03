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
//	<field>      ::= <field-begin> <map> <field-end>
//	<map>        ::= <map-begin> <field-data>* <map-end>
//	<field-data> ::= I8 | I16 | I32 | I64 | DOUBLE | STRING | BINARY
//	                 <struct> | <map> | <list> | <set>
//
// [Thrift IDL protocol spec]: https://github.com/apache/thrift/blob/eec0b584e657e4250e22f3fd492858d632e2aa7b/doc/specs/thrift-protocol-spec.md
func (p *TMap) WriteFieldData(cxt context.Context, oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteMapBegin(cxt, p.keyType, p.valueType, len(p.value)); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write map begin error", p), err)
		return
	}

	for k, v := range p.value {
		//	<map>        ::= <map-begin> <field-data>* <map-end>
		//	<field-data> ::= I8 | I16 | I32 | I64 | DOUBLE | STRING | BINARY
		//			<struct> | <map> | <list> | <set>
		if err = k.WriteFieldData(cxt, oprot); err != nil {
			err = thrift.PrependError(fmt.Sprintf("%T write key error", p), err)
			return
		}
		if err = v.WriteFieldData(cxt, oprot); err != nil {
			err = thrift.PrependError(fmt.Sprintf("%T write value error", p), err)
			return
		}
	}

	if err = oprot.WriteMapEnd(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write map end error", p), err)
		return
	}

	return
}

func (p *TMap) TType() thrift.TType {
	return thrift.MAP
}

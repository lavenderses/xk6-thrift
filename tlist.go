package thrift

import (
	"context"
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
)

type TList struct {
	value     []TValue
	valueType thrift.TType
}

func NewTList(v *[]TValue, valueType thrift.TType) *TList {
	return &TList{value: *v, valueType: valueType}
}

func (p *TList) Equals(other *TValue) bool {
	o, ok := (*other).(*TList)
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
//	<field>          ::= <field-begin> <list> <field-end>
//	<list>           ::= <list-begin> <field-data>* <list-end>
//	<list-begin>     ::= <list-elem-type> <list-size>
//	<list-elem-type> ::= <field-type>
//	<list-size>      ::= I32
//	<field-data>     ::= I8 | I16 | I32 | I64 | DOUBLE | STRING | BINARY
//			<struct> | <map> | <list> | <set>
//
// [Thrift IDL protocol spec]: https://github.com/apache/thrift/blob/eec0b584e657e4250e22f3fd492858d632e2aa7b/doc/specs/thrift-protocol-spec.md
func (p *TList) WriteFieldData(cxt context.Context, oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteListBegin(cxt, p.valueType, len(p.value)); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write list begin error", p), err)
		return
	}

	for _, v := range p.value {
		// if err = v.WriteField(cxt, oprot, fid, fname); err != nil {
		// 	err = thrift.PrependError(fmt.Sprintf("%T write list field data (field %d) error", p, fid), err)
		// 	return
		// }
		if err = p.writeFieldData(cxt, oprot, v); err != nil {
			err = thrift.PrependError(fmt.Sprintf("%T write list field data error", p), err)
			return
		}
	}

	if err = oprot.WriteListEnd(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write list end error", p), err)
		return
	}
	return
}

// See the above spec.
//
//	<field-data>     ::= I8 | I16 | I32 | I64 | DOUBLE | STRING | BINARY
//			<struct> | <map> | <list> | <set>
func (p *TList) writeFieldData(cxt context.Context, oprot thrift.TProtocol, value TValue) (err error) {
	if o, ok := value.(TString); ok {
		err = oprot.WriteString(cxt, o.value)
	} else if o, ok := value.(TBool); ok {
		err = oprot.WriteBool(cxt, o.value)
	} else if o, ok := value.(*TMap); ok {
		err = o.WriteFieldData(cxt, oprot)
	} else if o, ok := value.(*TStruct); ok {
		err = o.WriteFieldData(cxt, oprot)
	}
	return
}

func (p *TList) TType() thrift.TType {
	return thrift.LIST
}

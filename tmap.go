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

func NewTMap(keyType, valueType thrift.TType, v *map[TValue]TValue) *TMap {
	return &TMap{value: *v, keyType: keyType, valueType: valueType}
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
		err = thrift.PrependError(fmt.Sprintf("%T write map begin error: ", p), err)
		return
	}

	for k, v := range p.value {
		//	<map>        ::= <map-begin> <field-data>* <map-end>
		//	<field-data> ::= I8 | I16 | I32 | I64 | DOUBLE | STRING | BINARY
		//			<struct> | <map> | <list> | <set>
		if err = k.WriteFieldData(cxt, oprot); err != nil {
			err = thrift.PrependError(fmt.Sprintf("%T write key error: ", p), err)
			return
		}
		if err = v.WriteFieldData(cxt, oprot); err != nil {
			err = thrift.PrependError(fmt.Sprintf("%T write value error: ", p), err)
			return
		}
	}

	if err = oprot.WriteMapEnd(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write map end error: ", p), err)
		return
	}

	return
}

func (p *TMap) TType() thrift.TType {
	return thrift.MAP
}

func ReadMap(cxt context.Context, iproto thrift.TProtocol) (TValue, error) {
	keyType, valueType, size, err := iproto.ReadMapBegin(cxt)
	if err != nil {
		return nil, thrift.PrependError("error while reading map field: ", err)
	}

	tmap := make(map[TValue]TValue)
	for i := 0; i < size; i++ {
		if err = readFeidlDataList(cxt, iproto, &tmap, keyType, valueType); err != nil {
			return nil, thrift.PrependError("error while reading map: ", err)
		}
	}

	res := NewTMap(keyType, valueType, &tmap)
	return res, nil
}

func readFeidlDataList(cxt context.Context, iprot thrift.TProtocol, tmap *map[TValue]TValue, ktype, vtype thrift.TType) error {
	var key, value TValue
	var err error
	if key, err = ReadContainerData(ktype, cxt, iprot); err != nil {
		return err
	}
	if value, err = ReadContainerData(vtype, cxt, iprot); err != nil {
		return err
	}

	(*tmap)[key] = value
	return nil
}

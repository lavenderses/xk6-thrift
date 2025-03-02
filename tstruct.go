package thrift

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/apache/thrift/lib/go/thrift"
)

type TStructField struct {
	id   int16
	name string
}

type TStruct struct {
	value map[TStructField]TValue
}

func NewTStructField(id int16, name string) *TStructField {
	return &TStructField{id: id, name: name}
}

func NewTStruct(value *map[TStructField]TValue) *TStruct {
	return &TStruct{value: *value}
}

func (p *TStruct) Equals(other *TValue) bool {
	o, ok := (*other).(*TStruct)
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

// see [Thrift protocol spec @ 1a31d90 (v0.21.0)].
//
//	<field> ::= <field-begin> <field-data> <field-end>
//	<field-data> ::= I8 | I16 | I32 | I64 | DOUBLE | STRING | BINARY
//			<struct> | <map> | <list> | <set>
//	<struct> ::= <struct-begin> <field>* <field-stop> <struct-end>
//
// [Thrift protocol spec @ 1a31d90 (v0.21.0)]: https://github.com/apache/thrift/blob/1a31d9051d35b732a5fce258955ef95f576694ba/doc/specs/thrift-protocol-spec.md (v0.21.0)
func (p *TStruct) WriteFieldData(cxt context.Context, oprot thrift.TProtocol) (err error) {
	structName := "dummy"
	if err = oprot.WriteStructBegin(cxt, structName); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
		return
	}

	// write struct fields recursively
	for _, f := range slices.SortedFunc(maps.Keys(p.value), func(a, b TStructField) int {
		return int(a.id) - int(b.id)
	}) {
		v := p.value[f]
		ttype := v.TType()

		if err = oprot.WriteFieldBegin(cxt, f.name, ttype, f.id); err != nil {
			err = thrift.PrependError(fmt.Sprintf("%T write field begin error %d:%s", p, f.id, f.name), err)
			return
		}
		if err = v.WriteFieldData(cxt, oprot); err != nil {
			return
		}
		if err = oprot.WriteFieldEnd(cxt); err != nil {
			err = thrift.PrependError(fmt.Sprintf("%T write field end error %d:%s", p, f.id, f.name), err)
			return
		}
	}

	if err = oprot.WriteFieldStop(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write struct stop error: ", p), err)
		return
	}
	if err = oprot.WriteStructEnd(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
		return
	}
	return
}

func (p *TStruct) TType() thrift.TType {
	return thrift.STRUCT
}

func ReadStruct(cxt context.Context, iprot thrift.TProtocol) (TValue, error) {
	fieldName, err := iprot.ReadStructBegin(cxt)
	if err != nil {
		return nil, thrift.PrependError(fmt.Sprintf("error while struct begin (%s)", fieldName), err)
	}

	tvalue := make(map[TStructField]TValue)
	for {
		fname, ftype, fid, err := iprot.ReadFieldBegin(cxt)
		if err != nil {
			return nil, thrift.PrependError(fmt.Sprintf("error while reading field begin (%d:%s)", fid, fname), err)
		}
		if ftype == thrift.STOP {
			break
		}

		var tv TValue
		tv, err = ReadContainerData(ftype, cxt, iprot)
		if err != nil {
			return nil, err
		}

		err = iprot.ReadFieldEnd(cxt)
		if err != nil {
			return nil, err
		}

		// TODO: somehow fname gecomes an empty string
		tvalue[*NewTStructField(fid, fname)] = tv
	}

	err = iprot.ReadStructEnd(cxt)
	if err != nil {
		return nil, thrift.PrependError(fmt.Sprintf("error while reading struct end"), err)
	}

	res := NewTStruct(&tvalue)
	return res, nil
}

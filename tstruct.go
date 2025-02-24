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
//		<field-data> ::= I8 | I16 | I32 | I64 | DOUBLE | STRING | BINARY
//	<struct> | <map> | <list> | <set>
//	<struct> ::= <struct-begin> <field>* <field-stop> <struct-end>
//
// [Thrift protocol spec @ 1a31d90 (v0.21.0)]: https://github.com/apache/thrift/blob/1a31d9051d35b732a5fce258955ef95f576694ba/doc/specs/thrift-protocol-spec.md (v0.21.0)
func (p *TStruct) WriteField(cxt context.Context, oprot thrift.TProtocol, fid int16, fname string) (err error) {
	if err = oprot.WriteFieldBegin(cxt, fname, thrift.STRUCT, fid); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write field begin error %d:%s: ", p, fid, fname), err)
		return
	}
	if err = oprot.WriteStructBegin(cxt, fname); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write struct (%d, %s) begin error: ", p, fid, fname), err)
		return
	}

	// write struct fields recursively
	for _, f := range slices.SortedFunc(maps.Keys(p.value), func(a, b TStructField) int {
		return int(a.id) - int(b.id)
	}) {
		err = p.value[f].WriteField(cxt, oprot, f.id, f.name)
		if err != nil {
			return
		}
	}

	if err = oprot.WriteFieldStop(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write struct (%d, %s) stop error: ", p, fid, fname), err)
		return
	}
	if err = oprot.WriteStructEnd(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write struct (%d, %s) end error: ", p, fid, fname), err)
		return
	}
	if err = oprot.WriteFieldEnd(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write field end error %d:%s: ", p, fid, fname), err)
		return
	}
	return
}

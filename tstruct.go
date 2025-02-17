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

func (p *TStruct) WriteField(cxt context.Context, oprot thrift.TProtocol, fid int16, fname string) (err error) {
	if err = oprot.WriteFieldBegin(cxt, fname, thrift.STRUCT, fid); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write field begin error %d:%s: ", p, fid, fname), err)
		return
	}
	if err = oprot.WriteStructBegin(cxt, fname); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write struct (%d, %s) begin error: ", p, fid, fname), err)
	}

	// write struct fields recursively
	for _, f := range slices.SortedFunc(maps.Keys(p.value), func(a, b TStructField) int {
		return int(a.id) - int(b.id)
	}){
		p.value[f].WriteField(cxt, oprot, f.id, f.name)
	}

	if err = oprot.WriteFieldStop(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write struct (%d, %s) stop error: ", p, fid, fname), err)
	}
	if err = oprot.WriteStructEnd(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write struct (%d, %s) end error: ", p, fid, fname), err)
	}
	if err = oprot.WriteFieldEnd(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write field end error %d:%s: ", p, fid, fname), err)
		return
	}
	return
}
